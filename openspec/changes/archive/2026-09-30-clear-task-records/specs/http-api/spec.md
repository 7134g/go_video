## ADDED Requirements

### Requirement: 清空已结束任务接口
系统 SHALL 提供 `POST /api/tasks/clear-finished` 接口，用于删除所有状态为已完成（status=2）和失败（status=3）的任务记录。该接口 SHALL 仅删除数据库记录，不得删除磁盘上已下载的文件；状态为待执行（0）、执行中（1）、已暂停（4）的任务记录 SHALL 保持不变。

#### Scenario: 成功清空已结束任务
- **WHEN** 用户发送 `POST /api/tasks/clear-finished` 请求，且数据库中存在已完成或失败的任务
- **THEN** 系统删除这些任务记录，返回 200 状态码和形如 `{"deleted": N}` 的响应体，其中 N 为实际删除的记录数

#### Scenario: 没有可清空的任务
- **WHEN** 用户发送 `POST /api/tasks/clear-finished` 请求，且数据库中不存在已完成或失败的任务
- **THEN** 系统返回 200 状态码和 `{"deleted": 0}`，不产生任何删除操作

#### Scenario: 未结束任务被保留
- **WHEN** 用户发送 `POST /api/tasks/clear-finished` 请求，数据库中同时存在执行中、待执行、已暂停的任务
- **THEN** 这些任务记录在清空后仍然存在，正在进行的下载不受影响

#### Scenario: 下载文件不被删除
- **WHEN** 用户发送 `POST /api/tasks/clear-finished` 请求
- **THEN** 系统不调用任何文件删除逻辑，下载目录中的视频文件保持原样
