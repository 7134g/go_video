## MODIFIED Requirements

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

### Requirement: 添加并启动任务
系统 SHALL 提供 `AddAndStart` 方法，将任务先加入 tasks map，再存储回调并入队到 taskQueue。

#### Scenario: 添加并启动
- **WHEN** 调用 `AddAndStart(id, name, url, headerJSON, taskType, callback)`
- **THEN** 系统执行 AddTask 逻辑加入 tasks map，存储 callback 到 DTask，通过 taskQueue 入队，广播"任务已添加并启动: <name>"
