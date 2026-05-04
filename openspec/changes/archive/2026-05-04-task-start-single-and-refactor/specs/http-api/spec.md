## ADDED Requirements

### Requirement: 单个任务启动接口
系统 SHALL 提供 `POST /api/tasks/start-one` 接口用于启动指定的单个任务，任务 ID 通过请求体 `{"id": N}` 传递。

#### Scenario: 成功启动单个任务
- **WHEN** 用户发送 `POST /api/tasks/start-one` 请求，body 为 `{"id": N}`，且任务存在且不在执行中
- **THEN** 系统启动该任务并返回 200 状态码

#### Scenario: 任务不存在
- **WHEN** 用户发送 `POST /api/tasks/start-one` 请求，且任务不存在
- **THEN** 系统返回 404 状态码和错误信息

#### Scenario: 任务已在执行中
- **WHEN** 用户发送 `POST /api/tasks/start-one` 请求，且任务状态为执行中
- **THEN** 系统返回 400 状态码和错误信息

## MODIFIED Requirements

### Requirement: 删除任务接口
系统 SHALL 提供 `POST /api/tasks/delete` 接口用于删除任务，任务 ID 通过请求体 `{"id": N}` 传递。

#### Scenario: 成功删除任务
- **WHEN** 用户发送 `POST /api/tasks/delete` 请求，body 为 `{"id": N}`，且任务存在
- **THEN** 系统删除任务并返回 200 状态码

#### Scenario: 任务不存在
- **WHEN** 用户发送 `POST /api/tasks/delete` 请求，且任务不存在
- **THEN** 系统返回 404 状态码

### Requirement: 修改任务接口
系统 SHALL 提供 `POST /api/tasks/update` 接口用于修改任务，任务 ID 通过请求体 `{"id": N, ...}` 传递。

#### Scenario: 成功修改任务
- **WHEN** 用户发送 `POST /api/tasks/update` 请求，body 包含 id 和要修改的字段
- **THEN** 系统更新任务并返回 200 状态码和更新后的任务详情
