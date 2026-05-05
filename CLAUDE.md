# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此仓库中工作提供指引。

## 构建 / 运行

```bash
go build . && ./go_video          # 主程序（在 :8080 提供 Web 界面）
go build ./cmd/proxy && ./proxy   # CA 证书安装工具（需管理员权限）
go test ./...                     # 运行所有测试
```

前端位于 `web/`（Vue + Element Plus + Vite），通过 `//go:embed web/dist` 嵌入到 Go 二进制文件中。构建 Go 二进制之前，需先在 `web/` 下执行 `npm run build` 构建前端。

## 架构

这是一个视频下载器，支持 **m3u8/HLS 流**（含 AES-128 加密）和 **直接 MP4 下载**。内置 HTTPS MITM 代理，可拦截浏览器流量并自动捕获视频 URL。

### 分层结构

```
HTTP API (internal/api) → Service (internal/service) → Controller (internal/controller) → Downloader (internal/downloader)
```

- **`internal/api/`** — Gin HTTP 处理器：任务增删改查、配置获取/更新、WebSocket 进度推送。
- **`internal/service/`** — `TaskService` 连接 API 与控制器（管理生命周期和 Header 合并）。`ConfigService` 管理应用配置和 MITM 代理生命周期。
- **`internal/controller/`** — 单例 `DownloadController`，持有内存中的任务表，调度 m3u8/mp4 下载，跟踪进度，并向 WebSocket 监听者广播进度消息。
- **`internal/repository/`** — 通过 GORM/SQLite 持久化任务，通过 JSON 文件（`config.json`）持久化配置。
- **`internal/downloader/`** — `Pool` 基于信号量的 `Group` 机制，按域名限制并发下载数。
- **`pkg/m3u8/`** — M3U8 播放列表解析器（支持主播放列表、AES-128 密钥、字节范围），AES-128-CBC 解密，IV 推导，密钥缓存，以及 ffmpeg 合并。
- **`pkg/proxy/`** — 基于 `github.com/google/martian` 的 HTTPS MITM 代理。`Server.ModifyRequest` 拦截每个请求；`VideoDetector` 匹配 `.m3u8`/`.mp4` URL；捕获到的任务通过 channel 传给 `ConfigService.doTask` 进行持久化。
- **`cmd/proxy/`** — 独立的 CA 证书安装工具。

### 启动流程

1. `repository.InitDB()` — 打开 SQLite 并自动迁移 `Task` 表
2. `InitCa()` — 检查 CA 证书是否已安装到系统信任存储区（未安装则 panic）
3. 从 `config.json` 加载配置，应用到 `DownloadController`
4. 若 `interceptor_enabled` 为 true，在 goroutine 中启动 MITM 代理
5. Gin 服务器在 `:8080` 启动，提供嵌入的 Vue SPA 和 API 路由

### 关键行为

- **Header 合并**：任务级 Header 覆盖配置中的默认 Header；两者均以 `http.Header`（map[string][]string 格式）进行 JSON 序列化。
- **MP4 断点续传**：根据本地文件大小，使用 `Range` 请求头继续未完成的下载。
- **M3U8 断点续传**：跳过磁盘上已存在的分段文件。
- **并发控制**：任务级并发由 `MaxConcurrentTasks` 控制；分段级并发由 `MaxSegmentWorkers` 按域名限制。
- **代理任务去重**：`doTask` 在插入前检查 URL 是否已存在；若已存在但任务更新，则更新名称和 Header。
- **`pkg/m3u8/m3u8.go`** 中包含一个旧的本地 `parse()` 函数；实际使用的是 `method.go` 中的 `ParseM3u8Data`。
- **ffmpeg**：项目根目录下的 `ffmpeg.exe` 用于合并 M3U8 分段文件，下载完成后自动调用。

### 配置项（config.json）

| 字段 | 说明 | 默认值 |
|------|------|--------|
| `max_concurrent_tasks` | 最大并行任务数 | 3 |
| `max_segment_workers` | 每域名最大分段下载并发 | 10 |
| `download_dir` | 下载目录 | `./downloads` |
| `max_consecutive_errors` | 连续错误容忍数 | 10 |
| `default_headers` | 全局默认 HTTP 请求头 | {} |
| `interceptor_enabled` | 是否启用代理拦截 | true |
| `agent_address` | 代理监听地址 | `127.0.0.1:9999` |
| `vpn_address` | 上游 HTTP 代理地址 | `127.0.0.1:7890` |


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