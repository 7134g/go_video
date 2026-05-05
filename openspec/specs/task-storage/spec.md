## ADDED Requirements

### Requirement: Task 数据模型
系统 SHALL 定义 Task 模型，包含 id、name、url、header、type、status、created_at、updated_at 字段。

#### Scenario: 创建 Task 记录
- **WHEN** 调用创建方法并传入有效的 Task 数据
- **THEN** 系统在数据库中创建记录并返回包含自增 id 的 Task 对象

### Requirement: 数据库初始化
系统 SHALL 在启动时通过 GORM 自动迁移 Task 表结构（SQLite）。

#### Scenario: 首次启动
- **WHEN** 应用首次启动且数据库文件不存在
- **THEN** 系统创建 SQLite 数据库文件并自动创建 tasks 表

### Requirement: 查询未完成任务
系统 SHALL 提供方法查询所有状态为待执行(status=0)或已暂停(status=4)的任务。

#### Scenario: 存在未完成任务
- **WHEN** 调用查询未完成任务方法且数据库中存在待执行或已暂停任务
- **THEN** 系统返回所有 status IN (0, 4) 的任务列表

### Requirement: 按URL查询任务
系统 SHALL 提供按 URL 查询任务的方法，用于代理拦截去重。

#### Scenario: 代理检查重复URL
- **WHEN** 代理拦截到视频URL，调用 GetByURL
- **THEN** 若URL已存在且非运行中状态，则更新名称和Header；若不存在则创建新任务
