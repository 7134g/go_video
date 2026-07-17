## Why

当前项目只能通过浏览器访问 `127.0.0.1:8080` 使用，缺乏原生桌面应用体验。新增 Wails 打包方式可以让用户获得独立的桌面应用（窗口管理、系统托盘、原生菜单），同时保留现有 web 方式不变，两种部署形态共存。

## What Changes

- 新增 Wails 项目配置，将 Vue 3 前端与 Go 后端打包为桌面应用
- 新增 Wails 绑定层，将现有 REST API 调用替换为 Go 方法直接调用（`Call` 绑定），WebSocket 进度推送替换为 Wails 事件系统（`EventsEmit`/`EventsOn`）
- 现有 Gin HTTP 服务端保留不变，web 方式继续可用
- 新增构建脚本，生成 Windows/macOS 桌面安装包
- 启动时根据运行模式（web / desktop）选择不同的前端通信通道

## Capabilities

### New Capabilities
- `wails-desktop-app`: Wails 桌面应用壳，包含窗口管理、应用生命周期、系统托盘
- `wails-bindings`: Go ↔ JS 方法绑定层，替代 REST API 和 WebSocket，提供任务增删改查、配置读写、进度推送、ffmpeg 状态查询等接口
- `wails-build`: Wails 构建与打包流程，生成平台原生安装包（Windows .exe/msi、macOS .app）

### Modified Capabilities
<!-- 现有 spec 级别无变更，仅新增桌面入口，核心功能不变 -->

## Impact

- **新增依赖**: `github.com/wailsapp/wails/v2` (Go), `@wailsapp/runtime` (JS)
- **新增文件**: `wails.json` (Wails 配置)、`app.go` (Wails 应用入口)、`app_desktop.go` (桌面模式启动)、`build_desktop.sh` (桌面端构建脚本)
- **修改文件**: `main.go` (增加模式分支)、前端 `web/src/` (axios 调用改为 Wails 运行时调用，抽离 api 层)
- **保留不变**: `internal/` 全部核心逻辑、`pkg/` 全部公共包、`cmd/proxy/` 证书工具、现有 `build.sh`
