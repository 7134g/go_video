## ADDED Requirements

### Requirement: 控制器添加任务
系统 SHALL 提供 DTask 结构体和添加任务方法，支持 MP4 和 m3u8 两种格式。

#### Scenario: 添加 MP4 任务
- **WHEN** 调用添加任务方法，传入 url、header、任务名，type=mp4
- **THEN** 系统解析 url 和 header，构造 DTask 并加入任务队列，广播操作日志

#### Scenario: 添加 m3u8 任务
- **WHEN** 调用添加任务方法，传入 url、header、任务名，type=m3u8
- **THEN** 系统解析 url 和 header，构造 DTask 并加入任务队列，广播操作日志

### Requirement: 并发调度下载
系统 SHALL 通过任务队列和信号量机制调度下载，同时运行的任务数不超过 `MaxConcurrentTasks`。

#### Scenario: 批量启动下载
- **WHEN** 调用启动方法
- **THEN** 系统将所有任务入队到 taskQueue，由后台 dispatcher 按 `MaxConcurrentTasks` 限制逐个调度执行，广播启动数量

### Requirement: MP4 断点续传下载
系统 SHALL 支持 MP4 文件的断点续传。

#### Scenario: 文件已存在且未完成
- **WHEN** 下载目录存在 任务名.mp4 文件且文件大小小于远程文件
- **THEN** 系统使用 Range 请求从断点处继续下载，追加到文件末尾

#### Scenario: 文件不存在
- **WHEN** 下载目录不存在 任务名.mp4 文件
- **THEN** 系统创建新文件并从头下载

### Requirement: m3u8 分片下载
系统 SHALL 支持 m3u8 格式的分片下载和断点续传。

#### Scenario: 文件夹已存在
- **WHEN** 下载目录存在 任务名 文件夹
- **THEN** 系统检查已下载分片，仅下载缺失的分片

#### Scenario: 文件夹不存在
- **WHEN** 下载目录不存在 任务名 文件夹
- **THEN** 系统创建文件夹并下载所有分片

### Requirement: m3u8 分片合并
系统 SHALL 在所有分片下载完成后使用 ffmpeg 合并。

#### Scenario: 分片下载完成
- **WHEN** 所有分片下载完成
- **THEN** 系统调用 ffmpeg 合并分片为单个视频文件

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
系统 SHALL 提供 `AddAndStart` 组合方法，将任务加入内存并发送到执行队列。

#### Scenario: 添加并启动
- **WHEN** 调用 AddAndStart 一步完成添加和启动
- **THEN** 系统添加任务到 tasks map，存储回调，通过 taskQueue 入队，广播"任务已添加并启动: <name>"

### Requirement: 下载进度查询
系统 SHALL 提供 API 接口查询下载进度。

#### Scenario: 查询进度
- **WHEN** 前端请求下载进度接口
- **THEN** 系统返回各任务的下载进度信息（已下载/总大小）

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
