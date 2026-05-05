## 1. Controller 核心重构

- [x] 1.1 DTask 新增 `callback` 字段（TaskCallback 类型），存储任务完成后的回调
- [x] 1.2 在 `GetController()` 中启动 `dispatch()` 后台 goroutine，循环从 taskQueue 取任务，通过 taskSem 控制并发
- [x] 1.3 dispatch：获取槽位后启动 goroutine 执行 runTask，执行前检查 ctx.Done() 跳过已取消任务
- [x] 1.4 重构 `AddAndStart`：存储 callback 到 task，通过 `taskQueue <- task` 入队，移除直接 `go runTask()`
- [x] 1.5 重构 `StartTask`：存储 callback 到 task，通过 `taskQueue <- task` 入队，移除信号量获取和 goroutine 启动
- [x] 1.6 重构 `StartAll`：遍历 tasks map，存储 callback 到每个 task，通过 `taskQueue <- task` 入队，移除信号量获取和 goroutine 启动
- [x] 1.7 新增 `StopAll()` 方法：遍历 tasks map，对每个 task 调用 `cancel()`
- [x] 1.8 清理 `runningCount` 字段（未使用）

## 2. Service 层适配

- [x] 2.1 确认 `StartTasks`、`StartTask`、`AddAndStart`、`RetryTask` 的回调与 dispatcher 兼容（回调已包含 RemoveTask + 状态更新，无需改动）
- [x] 2.2 新增 `PauseAllTasks()` 方法：调用 `ctrl.StopAll()` 并将所有运行中任务的状态更新为 Paused

## 3. API 层

- [x] 3.1 新增 `POST /api/tasks/stop-all` 端点：调用 `TaskService.PauseAllTasks()`
- [x] 3.2 在 Gin router 中注册 `/api/tasks/stop-all` 路由

## 4. 验证

- [x] 4.1 `go build .` 编译通过
- [x] 4.2 验证 AddAndStart（代理拦截路径）受 MaxConcurrentTasks 限制
- [x] 4.3 验证 StartAll 批量启动后在队列中排队而非立即全部执行
- [x] 4.4 验证 StopAll 后所有任务停止，无 goroutine 泄漏
- [x] 4.5 验证任务完成后自动从队列移出，槽位释放给下一个任务
