## ADDED Requirements

### Requirement: 任务增删改查绑定
系统 SHALL 提供 Go 方法绑定，支持前端通过 `Call` 调用任务的增删改查操作，功能与 REST API 对等。

#### Scenario: 添加任务
- **WHEN** 前端调用 `AddTask(name, url, type, headers)` 绑定方法
- **THEN** 系统创建任务记录并返回任务详情
- **AND** 响应格式与 `POST /api/tasks` 一致

#### Scenario: 获取任务列表
- **WHEN** 前端调用 `GetTasks(status)` 绑定方法
- **THEN** 系统返回对应状态的任务列表
- **AND** 响应格式与 `GET /api/tasks?status=N` 一致

#### Scenario: 删除任务
- **WHEN** 前端调用 `DeleteTask(id)` 绑定方法
- **THEN** 系统删除对应任务并返回结果

#### Scenario: 更新任务
- **WHEN** 前端调用 `UpdateTask(id, fields)` 绑定方法
- **THEN** 系统更新任务字段并返回更新后的任务详情

### Requirement: 任务控制绑定
系统 SHALL 提供 Go 方法绑定，支持前端控制任务执行（启动、暂停、重试、全部停止）。

#### Scenario: 启动单个任务
- **WHEN** 前端调用 `StartOneTask(id)` 绑定方法
- **THEN** 系统将指定任务加入执行队列并返回结果

#### Scenario: 暂停任务
- **WHEN** 前端调用 `PauseTask(id)` 绑定方法
- **THEN** 系统取消该任务的下载上下文并更新状态

#### Scenario: 重试任务
- **WHEN** 前端调用 `RetryTask(id)` 绑定方法
- **THEN** 系统将失败或已暂停的任务重新入队

#### Scenario: 停止全部任务
- **WHEN** 前端调用 `StopAllTasks()` 绑定方法
- **THEN** 系统停止所有运行中的任务并更新数据库

#### Scenario: 批量启动任务
- **WHEN** 前端调用 `StartAllTasks()` 绑定方法
- **THEN** 系统将所有待执行和已暂停的任务加入执行队列

### Requirement: 配置管理绑定
系统 SHALL 提供 Go 方法绑定，支持前端读写应用配置。

#### Scenario: 获取配置
- **WHEN** 前端调用 `GetConfig()` 绑定方法
- **THEN** 系统返回当前应用配置（下载目录、并发数、默认请求头等）

#### Scenario: 更新配置
- **WHEN** 前端调用 `UpdateConfig(config)` 绑定方法
- **THEN** 系统保存配置到 `config.json` 并应用新配置到下载控制器

### Requirement: 进度事件推送
系统 SHALL 通过 Wails Events 系统向前端实时推送下载进度和广播消息，替代 WebSocket 通道。

#### Scenario: 控制器广播推送至前端
- **WHEN** 下载控制器广播进度更新或状态消息
- **THEN** 系统通过 `runtime.EventsEmit` 发送事件到前端
- **AND** 事件名称为 `task:progress` 或 `task:broadcast`

#### Scenario: 前端接收进度事件
- **WHEN** 前端通过 `runtime.EventsOn("task:progress", callback)` 订阅
- **THEN** 前端接收到进度数据并更新对应任务的进度展示

#### Scenario: 前端接收广播消息事件
- **WHEN** 前端通过 `runtime.EventsOn("task:broadcast", callback)` 订阅
- **THEN** 前端接收到广播消息（任务添加、暂停、删除等操作日志）

#### Scenario: 初始进度快照
- **WHEN** 前端启动并完成订阅
- **THEN** 系统通过绑定方法 `GetAllProgress()` 返回当前所有任务的进度快照

### Requirement: ffmpeg 状态绑定
系统 SHALL 提供 ffmpeg 状态查询和下载的绑定方法。

#### Scenario: 查询 ffmpeg 状态
- **WHEN** 前端调用 `GetFfmpegStatus()` 绑定方法
- **THEN** 系统返回 ffmpeg 是否已安装的信息

#### Scenario: 下载 ffmpeg
- **WHEN** 前端调用 `DownloadFfmpeg()` 绑定方法
- **THEN** 系统下载 ffmpeg 并返回结果

### Requirement: 更新标题绑定
系统 SHALL 提供方法用于从 WebTree 缓存获取 HTML 标题并更新任务名称。

#### Scenario: 更新任务标题
- **WHEN** 前端调用 `UpdateTaskTitle(id)` 绑定方法
- **THEN** 系统从 WebTree 查找 URL 对应的标题并更新任务名称
