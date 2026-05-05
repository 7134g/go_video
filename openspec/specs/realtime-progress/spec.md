## ADDED Requirements

### Requirement: WebSocket 连接管理
系统 SHALL 在任务页面建立 WebSocket 连接获取实时进度和广播消息。

#### Scenario: 建立连接
- **WHEN** 用户进入任务管理页面
- **THEN** 系统自动连接 WS /api/tasks/progress

#### Scenario: 断线重连
- **WHEN** WebSocket 连接断开
- **THEN** 系统 3 秒后自动尝试重新连接

### Requirement: 定时进度推送
系统 SHALL 每秒定时推送所有正在执行任务的进度数据。

#### Scenario: 定时推送
- **WHEN** WebSocket 连接已建立
- **THEN** 系统每秒推送一次所有任务的进度数组，包含每个任务的 id、name、type、downloaded、total、segment_done、segment_all、percent 字段

#### Scenario: 任务完成
- **WHEN** 任务进度达到 100%
- **THEN** 系统在后续的进度推送中不再包含该任务

### Requirement: 广播消息推送
系统 SHALL 在控制器发出广播消息时（如任务状态变更）通过 WebSocket 实时推送给客户端。

#### Scenario: 状态变更推送
- **WHEN** 控制器调用 BroadcastMessage 发出状态变更消息
- **THEN** 系统通过所有已连接 WebSocket 的消息监听通道实时推送该消息

### Requirement: 进度展示
系统 SHALL 实时展示正在执行任务的下载进度。

#### Scenario: 显示进度条
- **WHEN** 收到任务进度数据
- **THEN** 系统在对应任务行和右侧进度面板中显示进度百分比和进度条

#### Scenario: 显示分段进度
- **WHEN** 收到任务进度数据
- **THEN** 右侧进度面板显示分段完成进度（已完成分段数/总分段数）

### Requirement: 状态刷新
系统 SHALL 在 WebSocket 连接状态下定时刷新任务列表以反映最新状态。

#### Scenario: 连接状态刷新
- **WHEN** WebSocket 处于已连接状态
- **THEN** 系统每 5 秒自动刷新任务列表
