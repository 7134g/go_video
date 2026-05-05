## ADDED Requirements

### Requirement: 任务列表展示
系统 SHALL 展示所有任务的列表，包含任务名称、URL、类型、状态、创建时间，支持状态筛选。

#### Scenario: 加载任务列表
- **WHEN** 用户访问任务管理页面
- **THEN** 系统从 API 获取任务列表并以表格形式展示

#### Scenario: 状态筛选
- **WHEN** 用户通过状态下拉框选择某个状态
- **THEN** 系统只展示该状态的任务

#### Scenario: 空列表提示
- **WHEN** 没有任何任务
- **THEN** 系统显示"暂无任务"提示

#### Scenario: URL 复制
- **WHEN** 用户双击任务行的 URL
- **THEN** 系统将 URL 复制到剪贴板并显示成功提示

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
系统 SHALL 允许用户编辑待执行、已完成、失败或已暂停状态的任务。

#### Scenario: 打开编辑对话框
- **WHEN** 用户点击任务行的"编辑"按钮
- **THEN** 系统弹出预填充当前值的编辑表单

#### Scenario: 提交编辑
- **WHEN** 用户修改字段并提交
- **THEN** 系统调用 POST /api/tasks/update（body 含 id 和修改字段）更新任务

### Requirement: 删除任务
系统 SHALL 允许用户删除非执行中的任务。

#### Scenario: 确认删除
- **WHEN** 用户点击"删除"按钮
- **THEN** 系统显示确认对话框

#### Scenario: 执行删除
- **WHEN** 用户确认删除
- **THEN** 系统调用 POST /api/tasks/delete（body `{"id": id}`）并从列表移除

### Requirement: 批量启动任务
系统 SHALL 提供"开始下载"按钮批量启动所有待执行和已暂停任务。

#### Scenario: 批量启动任务
- **WHEN** 用户点击"开始下载"按钮
- **THEN** 系统调用 POST /api/tasks/start 批量启动任务，显示启动数量

### Requirement: 单个任务操作按钮
系统 SHALL 在任务列表操作列根据任务状态显示对应的操作按钮。

#### Scenario: 待执行任务的操作按钮
- **WHEN** 任务状态为待执行
- **THEN** 操作列显示"启动"、"更新标题"和"编辑"按钮

#### Scenario: 执行中任务的操作按钮
- **WHEN** 任务状态为执行中
- **THEN** 操作列显示"暂停"按钮，隐藏"删除"按钮

#### Scenario: 失败或已暂停任务的操作按钮
- **WHEN** 任务状态为失败或已暂停
- **THEN** 操作列显示"重试"和"编辑"按钮

#### Scenario: 已完成任务的操作按钮
- **WHEN** 任务状态为已完成
- **THEN** 操作列显示"编辑"按钮

### Requirement: 暂停全部按钮
系统 SHALL 提供"暂停全部"按钮用于一键暂停所有进行中的任务。

#### Scenario: 点击暂停全部
- **WHEN** 用户点击"暂停全部"按钮
- **THEN** 系统弹出确认对话框，确认后调用 POST /api/tasks/stop-all 暂停所有进行中任务，显示成功消息并刷新列表

### Requirement: 实时进度面板
系统 SHALL 在页面右侧面板展示所有正在执行任务的实时下载进度，包含百分比进度条和分段进度（已完成/总分段数）。

#### Scenario: 显示进度面板
- **WHEN** 存在正在进行中的任务
- **THEN** 右侧面板显示每个任务的名称、进度百分比和分段进度

#### Scenario: 无进行中任务
- **WHEN** 没有进行中的任务
- **THEN** 进度面板显示"暂无进行中的任务"

### Requirement: WebSocket 日志面板
系统 SHALL 在页面右侧底部面板展示 WebSocket 连接的实时日志，包含连接状态标签和原始消息数据。

#### Scenario: 显示连接状态
- **WHEN** WebSocket 连接成功或断开
- **THEN** 面板头部显示绿色"已连接"或红色"未连接"标签

#### Scenario: 显示消息日志
- **WHEN** 收到 WebSocket 消息
- **THEN** 日志面板显示带时间戳的原始 JSON 数据，最多保留 100 条

#### Scenario: 清空日志
- **WHEN** 用户点击"清空"按钮
- **THEN** 日志面板清空所有已记录的消息
