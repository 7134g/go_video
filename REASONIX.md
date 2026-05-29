# REASONIX.md

## 技术栈

- **Go 1.24** — module `go_video`
- **Gin v1.9.1** — HTTP 框架
- **GORM + SQLite** — 任务持久化
- **Vue 3 + Element Plus + Vite** — 前端 SPA，位于 `web/`
- **google/martian** — HTTPS MITM 代理
- **gorilla/websocket** — 下载进度实时推送

## 目录结构

- `main.go` — 入口；嵌入 `web/dist`，注册 Gin 路由，启动 MITM 代理
- `cmd/proxy/main.go` — CA 证书安装工具（各平台实现）
- `internal/api/` — HTTP 处理器：任务 CRUD、配置、WebSocket、ffmpeg
- `internal/service/` — TaskService / ConfigService（编排层）
- `internal/controller/` — DownloadController（m3u8/mp4 下载、进度、广播）
- `internal/repository/` — GORM/SQLite 任务库 + JSON 配置库
- `internal/downloader/` — 基于信号量的按域名并发池
- `internal/ffmpeg/` — ffmpeg 检测、下载、平台判断
- `internal/model/` — Task / Config 结构体
- `pkg/m3u8/` — M3U8 解析、AES-128-CBC 解密、Go 原生 remux、ffmpeg 合并
- `pkg/proxy/` — HTTPS MITM 代理、视频 URL 嗅探、请求通道、拦截器测试
- `chrome_ext/` — Chrome MV3 扩展（注入 `X-Tab-Id` 请求头）
- `web/` — Vue 3 + TypeScript + Vite 前端

## 常用命令

| 操作 | 命令 |
|---|---|
| 构建全部 | `bash build.sh`（前端 + Go 二进制 + 证书工具 + ffmpeg） |
| 前端开发 | `cd web && npm run dev` |
| 运行 | `go build . && ./go_video`（需 CA 证书 + ffmpeg） |
| 测试 | `go test ./...` |
| 前端构建 | `cd web && npm run build`（`vue-tsc -b && vite build` 类型检查） |
| 证书工具 | `go build -o build/install_cert.exe ./cmd/proxy` |

## 约定

- 前端使用 `vue-tsc -b && vite build` 构建（带类型检查）
- Go 测试分布在 `pkg/m3u8/`、`pkg/proxy/`、`internal/controller/`
- 无 lint 配置（无 `.golangci.yml`、无 Makefile）
- 配置存 `config.json`（已 gitignore），任务库 `video.db`（SQLite，启动时自动迁移）

## 注意

- **CA 证书必须先安装**。`InitCa()` 检查系统信任存储，未安装则告知各平台安装方式
- **`web/dist/` 是构建产物**。编辑 `web/src/`，不要直接改 `web/dist/`
- **`ffmpeg` 放项目根目录**。m3u8 合并优先用 Go 原生 remux，失败时回退到 ffmpeg
- **`config.json` 和 `video.db` 都是运行时文件**，已 gitignore
- **MITM 代理默认监听 `127.0.0.1:9999`**；上游代理地址 `127.0.0.1:7890`
- **Go 构建需要 GCC**。SQLite 通过 CGo 调用，Windows 需 TDM-GCC / MinGW
