## ADDED Requirements

### Requirement: 操作日志广播
系统 SHALL 在任务生命周期操作时通过 `BroadcastMessage` 向前端推送日志。

#### Scenario: 添加任务时广播
- **WHEN** AddTask 执行成功
- **THEN** 系统广播"任务已添加: <name>"

#### Scenario: 启动任务时广播
- **WHEN** StartTask 或 StartAll 执行成功
- **THEN** 系统广播启动信息（单任务: "任务已启动: <name>"，批量: "已启动 N 个任务"）

#### Scenario: 停止任务时广播
- **WHEN** StopTask 执行成功
- **THEN** 系统广播"任务已停止: <name>"

#### Scenario: 删除任务时广播
- **WHEN** RemoveTask 执行成功
- **THEN** 系统广播"任务已删除: <name>"

### Requirement: 添加并启动任务
系统 SHALL 提供 `AddAndStart` 组合方法。

#### Scenario: 代理拦截创建任务
- **WHEN** 代理拦截到视频 URL 后调用 AddAndStart
- **THEN** 系统添加任务并立即启动下载，广播"任务已添加并启动: <name>"

## MODIFIED Requirements

### Requirement: 控制器添加任务
系统 SHALL 提供 DTask 结构体和添加任务方法，支持 MP4 和 m3u8 两种格式。

#### Scenario: 添加 MP4 任务
- **WHEN** 调用添加任务方法，传入 url、header、任务名，type=mp4
- **THEN** 系统解析 url 和 header，构造 DTask 并加入任务队列，广播操作日志

#### Scenario: 添加 m3u8 任务
- **WHEN** 调用添加任务方法，传入 url、header、任务名，type=m3u8
- **THEN** 系统解析 url 和 header，构造 DTask 并加入任务队列，广播操作日志

### Requirement: 并发调度下载
系统 SHALL 并发执行所有 DTask 任务。

#### Scenario: 批量启动下载
- **WHEN** 调用启动方法
- **THEN** 系统为每个 DTask 启动 goroutine 并发执行下载，广播启动数量

## REMOVED Requirements

### Requirement: PauseTask 方法
**Reason**: `PauseTask` 方法重命名为 `StopTask`，语义更准确（实际是取消 context 停止下载）
**Migration**: 所有调用 `PauseTask` 的地方改为 `StopTask`
