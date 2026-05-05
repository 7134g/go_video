## ADDED Requirements

### Requirement: 停止所有任务
系统 SHALL 提供 `StopAll` 方法，停止所有当前在内存中的任务（包括排队中和执行中的任务）。

#### Scenario: 停止所有进行中的任务
- **WHEN** 调用 StopAll 且存在多个运行中的任务
- **THEN** 系统逐个调用每个任务的 cancel 函数，取消其 context，广播"已停止所有任务"

#### Scenario: 停止时包含排队任务
- **WHEN** 调用 StopAll 且 taskQueue 中还有等待执行的任务
- **THEN** 这些任务被取消 context，当 dispatcher 取出它们时检测到已取消并跳过执行

#### Scenario: 无运行任务时停止
- **WHEN** 调用 StopAll 且当前没有任何任务
- **THEN** 系统静默返回，不做任何操作

### Requirement: StopAll HTTP 端点
系统 SHALL 提供 `POST /api/tasks/stop-all` 接口用于停止所有进行中的任务。

#### Scenario: 成功停止所有任务
- **WHEN** 用户发送 `POST /api/tasks/stop-all` 请求
- **THEN** 系统调用 StopAll 停止所有任务，返回 200 状态码

#### Scenario: 无任务时调用
- **WHEN** 用户发送 `POST /api/tasks/stop-all` 请求且无进行中任务
- **THEN** 系统返回 200 状态码（幂等操作）
