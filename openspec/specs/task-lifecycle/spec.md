## ADDED Requirements

### Requirement: 添加任务
系统 SHALL 提供 `AddTask` 方法将任务加入内存 map，并向前端广播操作日志。

#### Scenario: 添加任务成功
- **WHEN** 调用 `AddTask(id, name, url, headerJSON, taskType)`
- **THEN** 系统构造 DTask 并加入 tasks map，调用 `BroadcastMessage(id, "任务已添加: <name>")`

### Requirement: 添加并启动任务
系统 SHALL 提供 `AddAndStart` 方法，将任务先加入 tasks map，再存储回调并入队到 taskQueue。

#### Scenario: 添加并启动
- **WHEN** 调用 `AddAndStart(id, name, url, headerJSON, taskType, callback)`
- **THEN** 系统执行 AddTask 逻辑加入 tasks map，存储 callback 到 DTask，通过 taskQueue 入队，广播"任务已添加并启动: <name>"

### Requirement: 启动单个任务
系统 SHALL 提供 `StartTask` 方法将指定任务加入执行队列。

#### Scenario: 启动任务成功
- **WHEN** 调用 `StartTask(id, callback)`
- **THEN** 系统找到该任务，存储回调到 DTask，通过 taskQueue 入队，广播"任务已启动: <name>"

#### Scenario: 任务不存在
- **WHEN** 调用 `StartTask(id)` 但 id 不在 tasks map 中
- **THEN** 系统返回错误

### Requirement: 启动所有任务
系统 SHALL 提供 `StartAll` 方法将所有内存中的任务加入执行队列。

#### Scenario: 批量启动
- **WHEN** 调用 `StartAll(callback)`
- **THEN** 系统遍历 tasks map，为每个任务存储回调并入队到 taskQueue，广播"已启动 N 个任务"

### Requirement: 停止任务
系统 SHALL 提供 `StopTask` 方法取消指定任务的下载。

#### Scenario: 停止任务成功
- **WHEN** 调用 `StopTask(id)`
- **THEN** 系统取消该任务的 context，广播"任务已停止: <name>"

#### Scenario: 任务不在内存中
- **WHEN** 调用 `StopTask(id)` 但任务不在内存中
- **THEN** 系统从数据库查找该任务并更新状态为已暂停

### Requirement: 删除任务
系统 SHALL 提供 `RemoveTask` 方法从内存中移除任务。

#### Scenario: 删除任务成功
- **WHEN** 调用 `RemoveTask(id)`
- **THEN** 系统从 tasks map 中删除该任务，广播"任务已删除: <name>"
