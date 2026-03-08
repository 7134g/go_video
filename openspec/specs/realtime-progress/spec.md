## ADDED Requirements

### Requirement: WebSocket 连接管理
系统 SHALL 在任务页面建立 WebSocket 连接获取实时进度。

#### Scenario: 建立连接
- **WHEN** 用户进入任务管理页面
- **THEN** 系统自动连接 WS /api/tasks/progress

#### Scenario: 断线重连
- **WHEN** WebSocket 连接断开
- **THEN** 系统自动尝试重新连接

### Requirement: 进度展示
系统 SHALL 实时展示正在执行任务的下载进度。

#### Scenario: 显示进度条
- **WHEN** 收到任务进度数据
- **THEN** 系统在对应任务行显示进度百分比和进度条

#### Scenario: 任务完成
- **WHEN** 任务进度达到 100%
- **THEN** 系统更新任务状态为已完成
