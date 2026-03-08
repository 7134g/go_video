## ADDED Requirements

### Requirement: 任务列表展示
系统 SHALL 展示所有任务的列表，包含任务名称、URL、类型、状态、创建时间。

#### Scenario: 加载任务列表
- **WHEN** 用户访问任务管理页面
- **THEN** 系统从 API 获取任务列表并以表格形式展示

#### Scenario: 空列表提示
- **WHEN** 没有任何任务
- **THEN** 系统显示"暂无任务"提示

### Requirement: 创建任务
系统 SHALL 提供表单让用户创建新任务，包含名称、URL、请求头、类型字段。

#### Scenario: 打开创建对话框
- **WHEN** 用户点击"新建任务"按钮
- **THEN** 系统弹出任务创建表单对话框

#### Scenario: 提交创建表单
- **WHEN** 用户填写必填字段并提交
- **THEN** 系统调用 POST /api/tasks 创建任务并刷新列表

#### Scenario: 表单验证失败
- **WHEN** 用户未填写必填字段就提交
- **THEN** 系统显示验证错误提示

### Requirement: 编辑任务
系统 SHALL 允许用户编辑待执行状态的任务。

#### Scenario: 打开编辑对话框
- **WHEN** 用户点击任务行的"编辑"按钮
- **THEN** 系统弹出预填充当前值的编辑表单

#### Scenario: 提交编辑
- **WHEN** 用户修改字段并提交
- **THEN** 系统调用 PUT /api/tasks/:id 更新任务

### Requirement: 删除任务
系统 SHALL 允许用户删除任务。

#### Scenario: 确认删除
- **WHEN** 用户点击"删除"按钮
- **THEN** 系统显示确认对话框

#### Scenario: 执行删除
- **WHEN** 用户确认删除
- **THEN** 系统调用 DELETE /api/tasks/:id 并从列表移除

### Requirement: 启动任务
系统 SHALL 提供按钮启动所有待执行任务。

#### Scenario: 启动任务
- **WHEN** 用户点击"开始下载"按钮
- **THEN** 系统调用 POST /api/tasks/start 启动任务
