## ADDED Requirements

### Requirement: 单个任务启动接口
系统 SHALL 提供 `POST /api/tasks/start-one` 接口用于启动指定的单个任务，任务 ID 通过请求体 `{"id": N}` 传递。

#### Scenario: 成功启动单个任务
- **WHEN** 用户发送 `POST /api/tasks/start-one` 请求，body 为 `{"id": N}`，且任务状态为待执行、失败或已暂停
- **THEN** 系统将任务加入下载控制器并开始下载，返回 200 状态码

#### Scenario: 任务不存在
- **WHEN** 用户发送 `POST /api/tasks/start-one` 请求，且任务不存在
- **THEN** 系统返回 404 状态码和错误信息

#### Scenario: 任务已在运行
- **WHEN** 用户发送 `POST /api/tasks/start-one` 请求，且任务已在执行中
- **THEN** 系统返回 400 状态码和错误信息

### Requirement: 单个任务启动按钮
系统 SHALL 在任务列表操作列提供单个任务的"启动"按钮。

#### Scenario: 按钮可见性
- **WHEN** 任务状态为待执行、失败或已暂停
- **THEN** 操作列显示"启动"按钮

#### Scenario: 点击启动
- **WHEN** 用户点击某个任务的"启动"按钮
- **THEN** 系统调用 `POST /api/tasks/start-one` 并刷新任务列表
