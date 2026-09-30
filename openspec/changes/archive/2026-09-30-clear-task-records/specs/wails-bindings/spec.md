## ADDED Requirements

### Requirement: 清空已结束任务绑定
系统 SHALL 提供 `ClearFinishedTasks()` 绑定方法，用于删除所有已完成和失败的任务记录，功能与 `POST /api/tasks/clear-finished` 对等。前端 API 层在桌面端运行时自动路由到该绑定方法。

#### Scenario: 清空已结束任务
- **WHEN** 前端调用 `ClearFinishedTasks()` 绑定方法
- **THEN** 系统删除状态为已完成或失败的任务记录，并返回被删除的记录数
- **AND** 未结束的任务记录保持存在

#### Scenario: 返回格式与 HTTP 一致
- **WHEN** 前端调用 `ClearFinishedTasks()` 绑定方法
- **THEN** 返回结果可在前端包装为 `{ data: { deleted: N } }`，与 `POST /api/tasks/clear-finished` 的响应体结构一致
