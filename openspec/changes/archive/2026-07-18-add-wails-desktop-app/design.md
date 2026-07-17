## Context

当前架构为 Go（Gin）后端 + Vue 3 前端，前端通过 `//go:embed` 嵌入二进制，在 `:8080` 提供完整 Web 应用。Wails 是将 Go 后端与 Web 前端打包为原生桌面应用的工具，它用系统 WebView 渲染前端，通过 Go 方法绑定替代 HTTP 调用。

关键约束：**保留 web 模式不变**，桌面模式作为可选构建目标。两种模式共享所有核心业务逻辑（`internal/`、`pkg/`）。

## Goals / Non-Goals

**Goals:**
- 新增 Wails 桌面应用入口，与现有 web 模式共存
- 前端 API 调用层抽象，支持 axios（web）和 Wails 绑定（桌面）两种后端通信方式
- 进度推送从 WebSocket 迁移到 Wails Events 事件系统
- 新增桌面构建脚本，产出 Windows/macOS 安装包

**Non-Goals:**
- 不修改任何核心下载逻辑（`internal/controller`、`internal/downloader`、`pkg/m3u8`）
- 不修改 MITM 代理功能
- 不删除或替换 web 模式
- 不引入系统托盘、自动更新等高级桌面功能（首版仅窗口应用）
- 不修改数据库 schema

## Decisions

### 1. Wails v2（非 v3）

**理由**: Wails v2 稳定（`github.com/wailsapp/wails/v2`），文档完善，社区成熟。v3 仍在 alpha 阶段。项目 Go 版本 1.24 完全兼容 v2。

### 2. 双模式通过构建标签切换

在 `main.go` 中通过 `//go:build` 标签区分入口：
- **web 模式**（默认 `go build`）：当前行为不变，Gin 服务器在 `:8080` 启动
- **desktop 模式**（`go build -tags desktop`）：Wails 接管生命周期，不启动 Gin 服务器

web 版本 `go build .` 不需要 Wails 依赖，避免强制引入 CGo/Wails 依赖。

### 3. 前端 API 层抽象：`apiClient` 接口

Vue 端引入 `apiClient` 接口，提供 `startTask`、`getTasks`、`getConfig` 等方法：
- **web 模式**: `httpClient` 实现，使用现有 axios 调用 `/api/*`
- **desktop 模式**: `wailsClient` 实现，使用 `window.go.main.App.<Method>()` 绑定调用

在 `main.ts` 中检测运行环境（检查 `window.go` 是否存在），自动选择实现。组件代码不变，只改 `import` 来源。

### 4. 事件推送：Wails Events 替代 WebSocket

Wails 内置事件系统（`runtime.EventsEmit` / `EventsOn`）天然支持 Go → JS 单向推送，与当前进度 WebSocket 功能对等：
- Go 端：`controller.BroadcastMessage` → `runtime.EventsEmit(ctx, "task:progress", data)`
- JS 端：`runtime.EventsOn("task:progress", callback)` 替代 WebSocket `onmessage`

这避免了在 Wails 内部再起 WebSocket 连接的冗余。

### 5. 目录结构：`desktop/` 子目录

Wails 特定文件放在 `desktop/` 下，避免污染根目录：
- `desktop/main.go` — Wails 应用入口（`func main()` with build tag）
- `desktop/app.go` — Wails `App` 结构体，包含所有绑定方法
- `desktop/wails.json` — Wails 项目配置
- `desktop/build.sh` — 桌面构建脚本

### 6. 绑定方法设计：一对一映射 API 路由

Wails 绑定方法直接映射现有 API handler 的核心逻辑，避免重复代码。`app.go` 中的方法调用 `internal/service` 和 `internal/controller`，与 Gin handler 共享底层实现。

## Risks / Trade-offs

- **CGo 依赖**: Wails v2 Windows 需要 CGo + gcc（mingw-w64），增加了 Windows 开发环境搭建复杂度 → 默认 `go build` 不带 `desktop` tag 不触发 CGo，web 构建不受影响
- **WebView2 运行时**: Windows 桌面版依赖系统 WebView2（Win10+ 已内置，Win7 需额外安装）→ 当前目标用户均为 Win10+，无需额外处理
- **前端条件编译**: JS 端没有编译时条件，需要通过运行时检测 `window.go` 选择 API 实现 → 不会引入未使用的 Wails runtime import（Vite tree-shaking 可处理）
- **双模式维护成本**: 两套入口 + 两套 API 调用层增加代码量 → API 层仅薄封装，核心逻辑不重复；若维护负担过大可后续合并
