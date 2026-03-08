## ADDED Requirements

### Requirement: 添加任务接口
系统 SHALL 提供 POST /api/tasks 接口用于添加下载任务。

#### Scenario: 成功添加任务
- **WHEN** 用户发送 POST /api/tasks 请求，包含 name、url、type 字段
- **THEN** 系统创建任务记录并返回 201 状态码和任务详情

#### Scenario: 缺少必填字段
- **WHEN** 用户发送 POST /api/tasks 请求，缺少 name 或 url 字段
- **THEN** 系统返回 400 状态码和错误信息

### Requirement: 删除任务接口
系统 SHALL 提供 DELETE /api/tasks/:id 接口用于删除任务。

#### Scenario: 成功删除任务
- **WHEN** 用户发送 DELETE /api/tasks/:id 请求，且任务存在
- **THEN** 系统删除任务并返回 200 状态码

#### Scenario: 任务不存在
- **WHEN** 用户发送 DELETE /api/tasks/:id 请求，且任务不存在
- **THEN** 系统返回 404 状态码

### Requirement: 修改任务接口
系统 SHALL 提供 PUT /api/tasks/:id 接口用于修改任务。

#### Scenario: 成功修改任务
- **WHEN** 用户发送 PUT /api/tasks/:id 请求，包含要修改的字段
- **THEN** 系统更新任务并返回 200 状态码和更新后的任务详情

### Requirement: 执行任务接口
系统 SHALL 提供 POST /api/tasks/start 接口用于批量执行未完成任务。

#### Scenario: 成功启动任务执行
- **WHEN** 用户发送 POST /api/tasks/start 请求
- **THEN** 系统将所有状态为待执行的任务添加到下载控制器，返回 200 状态码和启动的任务数量

### Requirement: WebSocket 进度推送接口
系统 SHALL 提供 GET /api/tasks/progress WebSocket 接口用于实时推送下载进度。

#### Scenario: 建立 WebSocket 连接
- **WHEN** 前端发起 WebSocket 连接到 /api/tasks/progress
- **THEN** 系统建立连接并开始定时推送进度

#### Scenario: 定时推送进度
- **WHEN** WebSocket 连接已建立
- **THEN** 系统每秒从下载控制器获取进度并推送给客户端（任务ID、已下载大小、总大小、百分比）

#### Scenario: 连接断开
- **WHEN** 客户端断开 WebSocket 连接
- **THEN** 系统停止该连接的定时推送

### Requirement: 配置管理 API 端点
系统 SHALL 在现有 API 路由中注册配置管理端点。

#### Scenario: 路由注册
- **WHEN** 服务器启动
- **THEN** 系统注册 GET /api/config 和 PUT /api/config 路由
