# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此仓库中工作提供指引。

**规则：本项目所有文档、注释、回复一律使用中文。**

## 构建 / 运行

```bash
# 开发：前后端分离
cd web && npm run dev             # Vite 开发服务器，将 /api 和 WebSocket 反代到 :8080（见 web/vite.config.ts）
go build . && ./go_video          # 主程序（在 :8080 提供嵌入的 Vue SPA + API）

# 一次性发布构建（Windows 目标，产物在 build/）：web 构建 + 主程序 + 证书安装器 + 拷贝 ffmpeg.exe
./build.sh

# 其他常用命令
go build ./cmd/proxy && ./proxy   # CA 证书安装工具（cmd/proxy 仅 Windows，依赖 golang.org/x/sys/windows；需管理员权限）
go test ./...                     # 运行所有测试（主要在 pkg/m3u8/）
```

前端位于 `web/`（Vue 3 + Element Plus + Vite + TypeScript），通过 `//go:embed web/dist` 嵌入到 Go 二进制文件中。`web/dist` 在 `.gitignore` 中，构建 Go 二进制前必须先 `npm run build`，否则 `embed` 失败。

## 架构

这是一个视频下载器，支持 **m3u8/HLS 流**（含 AES-128 加密）和 **直接 MP4 下载**。内置 HTTPS MITM 代理，可拦截浏览器流量并自动捕获视频 URL。

### 分层结构

```
HTTP API (internal/api) → Service (internal/service) → Controller (internal/controller) → Downloader (internal/downloader)
```

- **`internal/api/`** — Gin HTTP 处理器：任务增删改查、配置获取/更新、WebSocket 进度推送。
- **`internal/service/`** — `TaskService` 连接 API 与控制器（管理生命周期和 Header 合并）。`ConfigService` 管理应用配置和 MITM 代理生命周期。
- **`internal/controller/`** — 单例 `DownloadController`，持有内存中的任务表；通过缓冲 channel `taskQueue` 派发任务，`taskSem` 信号量按 `MaxConcurrentTasks` 限流；调度 m3u8/mp4 下载，跟踪进度，并向 WebSocket 监听者广播消息。
- **`internal/repository/`** — 通过 GORM/SQLite 持久化任务，通过 JSON 文件（`config.json`）持久化配置。
- **`internal/downloader/`** — `Pool` 基于信号量的 `Group` 机制，按域名限制并发下载数。
- **`pkg/m3u8/`** — M3U8 播放列表解析器（支持主播放列表、AES-128 密钥、字节范围），AES-128-CBC 解密，IV 推导，密钥缓存，以及 ffmpeg 合并。
- **`pkg/proxy/`** — 基于 `github.com/google/martian` 的 HTTPS MITM 代理。`Server.ModifyResponse` 拦截每个响应并按 `Content-Encoding` 解压（gzip / zstd），`GetVideo` 通过 URL 后缀（`.m3u8` / `.mp4`）识别视频；命中后经 `Collector` 投递到 channel，由 `ConfigService.doTask` 去重持久化。HTML 响应额外存入全局 `WebTree`（按 `X-Tab-Id` 分桶），供后续提取 `<title>` 作为任务名。
- **`cmd/proxy/`** — 独立的 CA 证书安装工具。

### 启动流程

1. `repository.InitDB()` — 打开 SQLite 并自动迁移 `Task` 表
2. `InitCa()` — 检查 CA 证书（`ca.crt`）是否已安装到系统信任存储区（Windows: `certutil -verifystore Root`；macOS: `security find-certificate`），未安装则 **panic**。`cmd/proxy` 安装器目前仅 Windows 实现；macOS 需手动 `security add-trusted-cert` 或借助 `proxy.InstallCert`。
3. 从 `config.json` 加载配置，应用到 `DownloadController`
4. `importTaskFile()` — 读取工作目录下的 `task.txt`，按行解析（奇数行为任务名称，偶数行为 URL），根据 URL 后缀（`.mp4` / `.m3u8`）判定 Type，使用配置的 `default_headers` 创建任务写入数据库，URL 已存在则跳过。处理完毕后删除 `task.txt`，文件不存在则静默跳过。
5. 若 `interceptor_enabled` 为 true，在 goroutine 中启动 MITM 代理
6. Gin 服务器在 `:8080` 启动，提供嵌入的 Vue SPA 和 API 路由

### 关键行为

- **Header 合并**：任务级 Header 覆盖配置中的默认 Header；两者均以 `http.Header`（map[string][]string 格式）进行 JSON 序列化。
- **MP4 断点续传**：根据本地文件大小，使用 `Range` 请求头继续未完成的下载。
- **M3U8 断点续传**：跳过磁盘上已存在的分段文件。
- **并发控制**：任务级并发由 `MaxConcurrentTasks` 控制；分段级并发由 `MaxSegmentWorkers` 按域名限制。
- **代理任务去重**：`doTask` 在插入前检查 URL 是否已存在；若已存在但任务更新，则更新名称和 Header。
- **`HasExactlyOneHttp` 过滤**：MITM 仅捕获字符串中恰好包含一个 `http(s)://` 的 URL — 用于排除埋点像素 / 跳转链接里嵌套的次级 URL。
- **WebTree / `X-Tab-Id`**：MITM 按 `X-Tab-Id` 请求头分桶记录每个 Tab 出现过的 URL，并从 HTML 响应中即时提取 `<title>` 后只保留标题字符串（不缓存 body），整体由 `sync.Mutex` 串行化、按 LRU 限容（默认 64 个 Tab）。`POST /api/tasks/update-title` 复用此缓存重新刷新已有任务的名称。
- **任务状态枚举** (`model.TaskStatus`)：Pending=0, Running=1, Completed=2, Failed=3, Paused=4。启动时 `repo.ResetStatus()` 会把残留的 Running 置回 Pending。
- **task.txt 批量导入**：`importTaskFile(cfg.DefaultHeaders)` 在启动时读取 `task.txt`，按行解析（跳过空行，奇数行为名称、偶数行为 URL），根据后缀判定 Type（`.mp4` → `"mp4"`，其余 → `"m3u8"`），URL 去重后写入数据库，完成后删除文件。无文件则静默跳过。单行失败不阻断其余任务。
- **ffmpeg**：合并 M3U8 分段优先使用**纯 Go** remux（`pkg/m3u8/remux.go` 的 `MergeFilesNative`，基于 `github.com/yapingcat/gomedia` 做 TS→MP4 容器转换，不依赖外部二进制）；失败时自动回退到项目根目录下的 `ffmpeg` 二进制（`MergeFilesFfmpeg`）。若 ffmpeg 也不存在则提示用户下载到当前目录。

### 配置项（config.json）

| 字段 | 说明 | 默认值 |
|------|------|--------|
| `max_concurrent_tasks` | 最大并行任务数 | 3 |
| `max_segment_workers` | 每域名最大分段下载并发 | 5 |
| `download_dir` | 下载目录 | `./downloads` |
| `max_consecutive_errors` | 连续错误容忍数 | 10 |
| `default_headers` | 全局默认 HTTP 请求头 | 预置 `user_agent`（Chrome UA） |
| `interceptor_enabled` | 是否启用代理拦截 | false |
| `agent_address` | 代理监听地址 | `127.0.0.1:9999` |
| `vpn_address` | 上游 HTTP 代理地址 | `127.0.0.1:7890` |
| `gin_mode` | Gin 框架模式（main.go: 空字符串走 release，非空走 debug） | `release` |
| `ffmpeg_prompt_declined` | 用户已拒绝启动时下载 ffmpeg，置 true 后不再追问 | false |


# karpathy

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

---

**These guidelines are working if:** fewer unnecessary changes in diffs, fewer rewrites due to overcomplication, and clarifying questions come before implementation rather than after mistakes.