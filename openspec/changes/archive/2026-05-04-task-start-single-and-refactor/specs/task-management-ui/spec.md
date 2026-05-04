## ADDED Requirements

### Requirement: 单个任务启动按钮
系统 SHALL 在任务列表操作列根据任务状态显示"启动"按钮。

#### Scenario: 待执行任务显示启动按钮
- **WHEN** 任务状态为待执行
- **THEN** 操作列显示"启动"按钮

#### Scenario: 失败任务显示启动按钮
- **WHEN** 任务状态为失败
- **THEN** 操作列显示"启动"按钮和"重试"按钮

#### Scenario: 已暂停任务显示启动按钮
- **WHEN** 任务状态为已暂停
- **THEN** 操作列显示"启动"按钮和"重试"按钮

#### Scenario: 点击启动单个任务
- **WHEN** 用户点击某任务的"启动"按钮
- **THEN** 系统调用 `POST /api/tasks/start-one`（body `{"id": id}`），显示成功消息并刷新列表

## MODIFIED Requirements

### Requirement: 启动任务
系统 SHALL 提供按钮启动所有待执行任务，并提供按钮单独启动某个任务。

#### Scenario: 批量启动任务
- **WHEN** 用户点击"开始下载"按钮
- **THEN** 系统调用 POST /api/tasks/start 批量启动任务

#### Scenario: 单个启动任务
- **WHEN** 用户点击任务行的"启动"按钮
- **THEN** 系统调用 POST /api/tasks/start-one（body `{"id": id}`）启动该任务

### Requirement: 编辑任务
系统 SHALL 允许用户编辑待执行状态的任务。

#### Scenario: 打开编辑对话框
- **WHEN** 用户点击任务行的"编辑"按钮
- **THEN** 系统弹出预填充当前值的编辑表单

#### Scenario: 提交编辑
- **WHEN** 用户修改字段并提交
- **THEN** 系统调用 POST /api/tasks/update（body 含 id 和修改字段）更新任务

### Requirement: 删除任务
系统 SHALL 允许用户删除任务。

#### Scenario: 确认删除
- **WHEN** 用户点击"删除"按钮
- **THEN** 系统显示确认对话框

#### Scenario: 执行删除
- **WHEN** 用户确认删除
- **THEN** 系统调用 POST /api/tasks/delete（body `{"id": id}`）并从列表移除
