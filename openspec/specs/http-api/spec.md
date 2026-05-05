## ADDED Requirements

### Requirement: 添加任务接口
系统 SHALL 提供 POST /api/tasks 接口用于添加下载任务。

#### Scenario: 成功添加任务
- **WHEN** 用户发送 POST /api/tasks 请求，包含 name、url、type 字段
- **THEN** 系统创建任务记录并返回 201 状态码和任务详情

#### Scenario: 缺少必填字段
- **WHEN** 用户发送 POST /api/tasks 请求，缺少 name 或 url 字段
- **THEN** 系统返回 400 状态码和错误信息

### Requirement: 任务列表接口
系统 SHALL 提供 GET /api/tasks 接口用于查询任务列表，支持 status 查询参数筛选。

#### Scenario: 获取全部任务
- **WHEN** 用户发送 GET /api/tasks 请求
- **THEN** 系统返回 200 状态码和所有任务列表

#### Scenario: 按状态筛选
- **WHEN** 用户发送 GET /api/tasks?status=N 请求
- **THEN** 系统返回 200 状态码和对应状态的任务列表

### Requirement: 单个任务启动接口
系统 SHALL 提供 `POST /api/tasks/start-one` 接口用于启动指定的单个任务，任务 ID 通过请求体 `{"id": N}` 传递。

#### Scenario: 成功启动单个任务
- **WHEN** 用户发送 `POST /api/tasks/start-one` 请求，body 为 `{"id": N}`，且任务状态为待执行、失败或已暂停
- **THEN** 系统将任务加入执行队列并返回 200 状态码

#### Scenario: 任务不存在或无法启动
- **WHEN** 用户发送 `POST /api/tasks/start-one` 请求，且任务不存在或状态不允许启动（执行中或已完成）
- **THEN** 系统返回 400 状态码和错误信息

### Requirement: 暂停单个任务接口
系统 SHALL 提供 `POST /api/tasks/pause` 接口用于暂停指定的任务，任务 ID 通过请求体 `{"id": N}` 传递。

#### Scenario: 成功暂停任务
- **WHEN** 用户发送 `POST /api/tasks/pause` 请求，body 为 `{"id": N}`，且任务正在执行中
- **THEN** 系统取消任务上下文并返回 200 状态码

#### Scenario: 任务不存在
- **WHEN** 用户发送 `POST /api/tasks/pause` 请求，且任务不存在
- **THEN** 系统返回 404 状态码

### Requirement: 重试任务接口
系统 SHALL 提供 `POST /api/tasks/retry` 接口用于重试失败或已暂停的任务。

#### Scenario: 成功重试
- **WHEN** 用户发送 `POST /api/tasks/retry` 请求，body 为 `{"id": N}`，且任务状态为失败或已暂停
- **THEN** 系统将任务重新加入执行队列并返回 200 状态码

#### Scenario: 任务状态不允许重试
- **WHEN** 用户发送 `POST /api/tasks/retry` 请求，且任务状态不是失败或已暂停
- **THEN** 系统返回 400 状态码

### Requirement: 暂停全部任务接口
系统 SHALL 提供 `POST /api/tasks/stop-all` 接口用于暂停所有进行中的任务。

#### Scenario: 成功暂停全部
- **WHEN** 用户发送 `POST /api/tasks/stop-all` 请求
- **THEN** 系统停止所有进行中的任务、更新数据库状态为已暂停，并返回 200 状态码

### Requirement: 更新标题接口
系统 SHALL 提供 `POST /api/tasks/update-title` 接口用于从 WebTree 缓存中获取 URL 对应的 HTML 标题并更新任务名称。

#### Scenario: 成功更新标题
- **WHEN** 用户发送 `POST /api/tasks/update-title` 请求，body 为 `{"id": N}`，且 WebTree 中存在对应标题
- **THEN** 系统更新任务名称并返回 200 状态码和更新后的任务

#### Scenario: 未找到标题
- **WHEN** WebTree 中不存在该 URL 对应的标题
- **THEN** 系统返回 200 状态码，提示未找到标题

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

### Requirement: 执行任务接口
系统 SHALL 提供 POST /api/tasks/start 接口用于批量执行未完成任务（待执行和已暂停）。

#### Scenario: 成功启动任务执行
- **WHEN** 用户发送 POST /api/tasks/start 请求
- **THEN** 系统将所有状态为待执行或已暂停的任务加入到执行队列，返回 200 状态码和启动的任务数量

### Requirement: WebSocket 进度推送接口
系统 SHALL 提供 GET /api/tasks/progress WebSocket 接口用于实时推送下载进度和广播消息。

#### Scenario: 建立 WebSocket 连接
- **WHEN** 前端发起 WebSocket 连接到 /api/tasks/progress
- **THEN** 系统建立连接并注册消息监听器

#### Scenario: 定时推送进度
- **WHEN** WebSocket 连接已建立
- **THEN** 系统每秒从下载控制器获取所有任务进度并推送给客户端

#### Scenario: 广播消息推送
- **WHEN** 控制器发出广播消息（如任务状态变更）
- **THEN** 系统通过消息监听通道实时推送给所有已连接的客户端

#### Scenario: 连接断开
- **WHEN** 客户端断开 WebSocket 连接
- **THEN** 系统停止该连接的定时推送并移除消息监听器

### Requirement: 配置管理 API 端点
系统 SHALL 在现有 API 路由中注册配置管理端点。

#### Scenario: 路由注册
- **WHEN** 服务器启动
- **THEN** 系统注册 GET /api/config 和 PUT /api/config 路由
