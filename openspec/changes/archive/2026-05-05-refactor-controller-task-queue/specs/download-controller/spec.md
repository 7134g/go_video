## MODIFIED Requirements

### Requirement: 并发调度下载
系统 SHALL 通过任务队列和信号量机制调度下载，同时运行的任务数不超过 `MaxConcurrentTasks`。

#### Scenario: 批量启动下载
- **WHEN** 调用启动方法
- **THEN** 系统将所有任务入队到 taskQueue，由后台 dispatcher 按 `MaxConcurrentTasks` 限制逐个调度执行，广播启动数量

### Requirement: 添加并启动任务
系统 SHALL 提供 `AddAndStart` 组合方法，将任务加入内存并发送到执行队列。

#### Scenario: 添加并启动
- **WHEN** 调用 AddAndStart 一步完成添加和启动
- **THEN** 系统添加任务到 tasks map，存储回调，通过 taskQueue 入队，广播"任务已添加并启动: <name>"

## ADDED Requirements

### Requirement: Dispatcher 后台调度
系统 SHALL 在控制器初始化时启动一个后台 dispatcher goroutine，持续从 taskQueue 消费任务并通过 taskSem 控制并发。

#### Scenario: Dispatcher 正常调度
- **WHEN** taskQueue 中有待执行任务且运行数未达上限
- **THEN** dispatcher 获取槽位并启动 runTask goroutine

#### Scenario: Dispatcher 等待槽位
- **WHEN** taskQueue 中有待执行任务但运行数已达上限
- **THEN** dispatcher 阻塞在 taskSem 直到有槽位释放

### Requirement: 停止所有任务
系统 SHALL 提供 `StopAll` 方法停止 tasks map 中所有任务的 context。

#### Scenario: 停止所有任务
- **WHEN** 调用 `StopAll()`
- **THEN** 系统遍历 tasks map 调用每个任务的 cancel()，广播操作日志
