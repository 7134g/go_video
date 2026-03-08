## ADDED Requirements

### Requirement: Task 数据模型
系统 SHALL 定义 Task 模型，包含 id、name、url、header、type、status、created_at、updated_at 字段。

#### Scenario: 创建 Task 记录
- **WHEN** 调用创建方法并传入有效的 Task 数据
- **THEN** 系统在数据库中创建记录并返回包含自增 id 的 Task 对象

### Requirement: 数据库初始化
系统 SHALL 在启动时自动创建 SQLite 数据库文件 video.db 并迁移 Task 表结构。

#### Scenario: 首次启动
- **WHEN** 应用首次启动且 video.db 不存在
- **THEN** 系统创建 video.db 文件并自动创建 task 表

### Requirement: 查询未完成任务
系统 SHALL 提供方法查询所有状态为待执行(status=0)的任务。

#### Scenario: 存在未完成任务
- **WHEN** 调用查询未完成任务方法且数据库中存在待执行任务
- **THEN** 系统返回所有 status=0 的任务列表
