## Context

`DownloadController` 当前有 `taskSem`（信号量）和 `taskQueue`（channel）两个字段，但 `AddAndStart` 直接 `go runTask()` 绕过了信号量。`taskQueue` 未被使用，`runningCount` 也未维护。`StopAll` 功能缺失。

任务启动路径有 4 条：`AddAndStart`（代理拦截）、`StartTask`（单任务启动）、`StartAll`（批量启动）、`RetryTask`（重试，内部调用 `StartAll`）。除 `AddAndStart` 外，其余都通过 `taskSem` 控制并发，但方式是在 goroutine 中阻塞等待信号量——N 个任务就创建 N 个 goroutine，不够优雅。

## Goals / Non-Goals

**Goals:**
- 所有任务启动路径统一通过 `taskQueue` channel 入队
- 后台 dispatcher goroutine 从队列取任务，通过 `taskSem` 控制并发
- 新增 `StopAll` 方法，停止所有进行中的任务
- 任务完成（成功/失败/取消）后自动从 tasks map 移除

**Non-Goals:**
- 不改变 HTTP API 接口签名（只新增 `StopAll` 端点，不删除现有端点）
- 不改变 Service 层的回调模式
- 不修改 `downloader.Pool` 的子任务并发控制（已由 `MaxSegmentWorkers` 管理）

## Decisions

### 1. Dispatcher 模式：单 goroutine 消费 channel

```
taskQueue (chan) → dispatch() goroutine → taskSem acquire → go runTask()
```

`dispatch()` 在 `GetController()` 初始化时启动，循环读取 `taskQueue`：
- 取到任务后，阻塞等待 `taskSem <- struct{}{}`（获取得一个并发槽位）
- 获取槽位后，启动 goroutine 执行 `runTask`，goroutine 结束时释放槽位

**替代方案考虑**：在调用处直接 `go func() { taskSem<-...; runTask() }()`。被否决——那只是把信号量挪到外层，没有真正解决 `AddAndStart` 绕过的问题，且没有统一入口。

### 2. 任务启动统一：存储 callback 到 DTask，入队而非直接执行

`StartTask`、`StartAll`、`AddAndStart` 不再启动 goroutine，改为：
1. 将 callback 存储到 `DTask.callback`
2. 通过 `taskQueue <- task` 入队

Dispatcher 取出任务后直接调用 `c.runTask(task, task.callback)`。

**替代方案考虑**：callback 通过 channel 传递（如 `chan struct{task *DTask; cb TaskCallback}`）。被否决——无必要增加中间类型，callback 仅由任务自身使用。

### 3. StopAll 实现：遍历 tasks map，逐个 cancel

```go
func (c *DownloadController) StopAll() {
    c.mu.RLock()
    defer c.mu.RUnlock()
    for _, t := range c.tasks {
        t.cancel()
    }
}
```

对于队列中尚未执行的任务：`dispatch()` 在启动 `runTask` 前检查 `ctx.Done()`，若已取消则直接调用 callback 返回 `context.Canceled` 并跳过下载。

**替代方案考虑**：清空 `taskQueue` channel。被否决——channel 不支持批量清空，逐个 drain 可能阻塞 dispatcher。

### 4. 任务移除：由 Service 层 callback 处理，保持不变

任务完成后，Service 层传入的 callback 负责 `RemoveTask` + 更新 DB 状态。Controller 不主动移除任务，保持职责边界。

## Risks / Trade-offs

- [Dispatcher 是单 goroutine] → `taskSem` 获取是阻塞的，dispatcher 等待槽位时整个队列暂停出队。这是期望行为——如果槽位满了，不应该再取下一个任务。释放槽位后 dispatcher 立即取下一个任务。
- [channel 无法移除已入队任务] → 已入队但未执行的任务被 StopAll 取消后，将由 dispatcher 的 `ctx.Done()` 检查快速跳过。跳过的 goroutine 开销极小。
- [`StartAll` 现有 goroutine 模式被替换] → `StartAll` 的行为从 "N goroutine 争抢槽位" 变为 "N 任务顺序入队"。功能等价，但内存占用更低。
