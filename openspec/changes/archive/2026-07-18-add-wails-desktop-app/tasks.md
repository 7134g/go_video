## 1. 环境准备

- [x] 1.1 安装 Wails CLI：`go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- [x] 1.2 验证环境：`wails doctor` 检查 Go、Node、WebView2 依赖是否就绪

## 2. Wails 项目结构

- [x] 2.1 创建 `desktop/` 目录（调整为根目录文件 + build tag 隔离）
- [x] 2.2 创建 `wails.json`，配置应用名、前端路径、入口文件
- [x] 2.3 创建 `app.go`（`//go:build desktop`），定义 `App` 结构体及 Wails 生命周期钩子（`startup`、`shutdown`）
- [x] 2.4 创建 `main_desktop.go`（`//go:build desktop`），调用 `wails.Run()` 启动应用
- [x] 2.5 创建 `shared.go` 提取公共启动逻辑，`main.go` 添加 `//go:build !desktop` 保留 web 模式

## 3. Go 绑定方法实现

- [x] 3.1 在 `app.go` 中实现任务增删改查绑定方法（`AddTask`、`GetTasks`、`DeleteTask`、`UpdateTask`），复用 `internal/service` 和 `internal/repository`
- [x] 3.2 实现任务控制绑定方法（`StartOneTask`、`PauseTask`、`RetryTask`、`StopAllTasks`、`StartAllTasks`），复用 `internal/controller`
- [x] 3.3 实现配置管理绑定方法（`GetConfig`、`UpdateConfig`），复用 `internal/service/ConfigService`
- [x] 3.4 实现 ffmpeg 状态绑定方法（`GetFfmpegStatus`、`DownloadFfmpeg`）
- [x] 3.5 实现 `UpdateTaskTitle` 绑定方法，复用 WebTree 标题缓存逻辑

## 4. 事件推送替代 WebSocket

- [x] 4.1 在 `app.go` 中实现 `bridgeMessages` goroutine，将控制器广播转为 `runtime.EventsEmit("task:broadcast", ...)`
- [x] 4.2 实现 `bridgeProgress` goroutine，将 `ProgressInfo` 转为 `runtime.EventsEmit("task:progress", ...)`
- [x] 4.3 实现 `GetAllProgress` 绑定方法，返回内存中所有任务的初始进度快照

## 5. 前端 API 层抽象

- [x] 5.1 在 `web/src/api/` 创建 `index.ts` + `wails.ts`，`taskApi`/`configApi` 的现有形状作为隐式接口
- [x] 5.2 保留现有 `api/task.ts` + `api/config.ts`（axios 调用）作为 http 实现
- [x] 5.3 创建 `wailsClient.ts`（`api/wails.ts`），使用 `window.go.main.App.<Method>()` 实现 Wails 绑定调用
- [x] 5.4 创建 `api/index.ts`，运行时检测 `window.go` 自动选择 http 或 wails 实现导出
- [x] 5.5 替换前端组件中的 `import`，从 `'../api/task'`/`'../api/config'` 改为 `'../api'` 统一导入

## 6. Wails 事件订阅（前端）

- [x] 6.1 创建 `web/src/composables/useWailsEvents.ts`，封装 `EventsOn("task:progress")` 和 `EventsOn("task:broadcast")` 订阅逻辑
- [x] 6.2 在 TaskList.vue 中集成 `useWailsEvents`，web 模式保持 WebSocket 不变
- [x] 6.3 实现组件卸载时的事件取消订阅（composable 内部 `onUnmounted` 调 `EventsOff`）

## 7. 应用图标与资源

- [x] 7.1 ~~准备 `appicon.png`~~ 跳过：视觉资源，需用户自行提供后配置到 `wails.json`
- [x] 7.2 在 `wails.json` 中保留图标路径配置项（用户自行放入图标后取消注释）

## 8. 构建脚本与验证

- [x] 8.1 创建 `build_desktop.sh`：先 `npm run build`，再 `wails build -tags desktop`
- [x] 8.2 在 `build.sh`（现有）中增加 desktop 构建步骤（`--desktop` 参数）
- [x] 8.3 更新 `web/vite.config.ts`，`base: './'` 确保 Wails 相对路径正确
- [x] 8.4 验证 web 模式不受影响：`go build .` 成功且不引入 Wails 依赖
- [x] 8.5 验证 `go build -tags desktop .` 编译成功
