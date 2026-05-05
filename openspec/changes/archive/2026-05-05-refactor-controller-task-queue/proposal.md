## Why

`AddAndStart` 直接 `go runTask()` 绕过了 `taskSem` 信号量，任何走该路径的调用都会不受 `MaxConcurrentTasks` 限制。此外 `taskQueue` channel 和 `runningCount` 字段已被声明但从未使用，`StopAll` 功能缺失。

## What Changes

- 激活并使用已有的 `taskQueue` channel 作为任务执行队列
- 引入后台 dispatcher goroutine，从 `taskQueue` 取任务，通过 `taskSem` 控制并发
- `AddAndStart`、`StartTask`、`StartAll` 改为将任务入队而非直接运行
- 新增 `StopAll` 方法，停止所有进行中的任务
- **BREAKING**: 移除不再需要的 `runningCount` 字段（原本未使用），`taskSem` 由 dispatcher 统一管理

## Capabilities

### New Capabilities
- `task-queue`: 基于 channel 的任务执行队列，后台 dispatcher 按 `MaxConcurrentTasks` 调度任务执行
- `stop-all-tasks`: 一次性停止所有进行中的下载任务

### Modified Capabilities
- `download-controller`: 并发调度需求变更——从"并发执行所有任务"改为"通过队列按上限调度"。任务完成后自动从队列移除
- `task-lifecycle`: 启动任务的行为从"立即启动 goroutine"改为"入队等待调度"

## Impact

- `internal/controller/controller.go` — 重构 `StartAll`、`StartTask`、`AddAndStart`，新增 dispatcher 和 `StopAll`
- `internal/service/task.go` — `StartTasks`、`RetryTask` 回调中调用 `RemoveTask` 的逻辑可能需要调整（任务由 dispatcher 自动移除）
- `internal/controller/dtask.go` — 可能需要新增 `Status` 字段区分排队/执行中状态
- `internal/api/task.go` — 可能需要新增 `StopAll` HTTP 端点