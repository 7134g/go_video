## Context

当前进度值通过 WebSocket ticker 每秒轮询 `GetAllProgress()` 推送全部任务的进度数组，存在最长 1 秒延迟且每秒全量序列化。目标改为在 `Progress.AddDone()` / `IncrementDone()` 被调用时主动推单条 `ProgressInfo`，实现即时更新。

## Goals / Non-Goals

**Goals:**
- 在 `AddDone`/`IncrementDone` 每次更新计数后立即推送该任务的 `ProgressInfo`
- 任务启动时（`AddAndStart` / `StartTask` / `StartAll`）推送一次该任务的初始 `ProgressInfo`（done=0,total=0）
- WebSocket 连接在建立时发送一次全量快照，后续仅靠事件推送增量更新
- 复用现有的全局广播 pub/sub 模式，与 `BroadcastMessage` 平行
- 前端单次任务进度更新时不需全量替换列表

**Non-Goals:**
- 不修改 `ProgressInfo` 结构体或新增字段（`Status` 等留待后续）
- 不涉及任务状态变更的广播（生命周期事件仍走 `BroadcastMessage`）
- 不改变 `SetTotal` 的行为（Total 的变更不在本需求范围内）

## Decisions

1. **Progress 新增字段，通过全局 BroadcastProgress 广播**
   - `Progress` 增加 `taskID uint` 和 `taskName string` 字段，创建 `DTask` 时赋值
   - 新增全局 `BroadcastProgress(ProgressInfo)` 函数，与 `BroadcastMessage` 平行
   - `AddDone`/`IncrementDone` 在**释放写锁后**调用 `BroadcastProgress`，避免持有锁时调用外部代码
   - **理由**: 匹配现有代码风格（全局广播 + 非阻塞 channel）；不需要在 Progress 中存储函数引用；方便单元测试

2. **移除 WebSocket ticker，改为连接时快照 + 事件推送**
   - 连接建立后立刻发送一次 `GetAllProgress()` 全量快照（用作初始状态）
   - 移除每秒 ticker
   - 新增 `progressCh` channel 接收 `ProgressInfo` 广播，每次收到后直接 WriteMessage 单个 JSON 对象
   - **理由**: 避免冗余推送；前端按消息类型（数组/对象）区分即可；无进度变化时零推送

3. **前端兼容单条 ProgressInfo**
   - `ws.onmessage` 中检测：`Array.isArray(data)` → 全量快照（现有逻辑）；`data.id !== undefined` → 单条进度更新
   - 单条更新时只修改对应 `progress.value[data.id]` 和 `taskProgressList.value` 中对应项
   - **理由**: 后端推单条，前端不需要全量替换

4. **Progress 内构造 ProgressInfo 时不做额外查询**
   - `AddDone`/`IncrementDone` 持有 `p.Done` / `p.Total` / `p.taskID` / `p.taskName` / `p.Type` 的最新值，直接在方法内构造 `ProgressInfo` 后广播
   - Percent 在 Progress 内计算，保持与 `GetAllProgress` 一致的算法

## Risks / Trade-offs

- [Connect 时快照] → 如果连接刚建立时大量任务正在频繁更新进度，快照可能在推送到达之前被稍旧的进度覆盖，但下一次更新会立即补上，影响极小
- [非阻塞丢弃] → `BroadcastProgress` 沿用 `BroadcastMessage` 的非阻塞丢弃策略，极端高并发下进度更新可能丢失，但下一个 IncrementDone 会补上；相比当前每秒全量推的延迟已大幅改善
- [暂无修改 SetTotal] → 如果某个任务 Total 在运行中变更（如动态 m3u8），进度百分比可能不准确，但当前代码也如此，属于后续改进范围
