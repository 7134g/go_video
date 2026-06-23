## Why

当前进度值仅通过 WebSocket ticker 每秒轮询 `GetAllProgress()` 推送，存在最长 1 秒的延迟，且每秒都会序列化全部任务（含无进度变化的 idle 任务），效率不高。改为在 `Progress.AddDone()` / `Progress.IncrementDone()` 被调用时主动推送进度消息，实现即时更新。

## What Changes

- `Progress` 结构体新增 `TaskID` 字段，使其能够生成带上下文的进度消息
- 新增全局进度广播机制 `BroadcastProgress(info ProgressInfo)`，与已有的 `BroadcastMessage` 平行
- `Progress.AddDone()` / `Progress.IncrementDone()` 在更新计数后调用 `BroadcastProgress` 推送 `ProgressInfo`
- WebSocket 处理程序新增监听进度广播通道，同时保留 ticker 作为连接初始快照和兜底
- `GetAllProgress()` 保留，只在 WebSocket 连接建立时发送一次初始快照

## Capabilities

### New Capabilities
- `progress-push`: 在 `AddDone`/`IncrementDone` 被调用时实时推送 `ProgressInfo` 到 WebSocket 客户端

### Modified Capabilities
<!-- 无已有 spec 需要修改 -->

## Impact

- `internal/controller/dtask.go` — `Progress` 结构体增加 `TaskID`；`AddDone`/`IncrementDone` 增加广播调用
- `internal/controller/broadcast.go` — 新增 `ProgressListeners` / `BroadcastProgress` / 订阅管理
- `internal/controller/controller.go` — 创建 `DTask` 时设置 `Progress.TaskID`
- `internal/api/websocket.go` — 新增从进度广播通道接收消息并写入 WS 的分支；连接建立时发送一次 `GetAllProgress()` 替代每秒轮询
