## Why

任务列表会持续累积已完成和失败的历史记录，用户只能逐条点击"删除"，清理几十条旧记录需要反复确认、耗时且易误删。需要一个一键清空已结束记录的入口。

## What Changes

- 新增后端接口 `POST /api/tasks/clear-finished`，批量删除数据库中状态为已完成（status=2）和失败（status=3）的任务记录。
- 在 `TaskRepository` 中新增按状态批量删除的方法，`TaskService` 暴露对应的 `ClearFinishedTasks`。
- 在 Wails 绑定中新增 `ClearFinishedTasks()` 方法，使桌面端与 HTTP 模式能力对等（前端 API 层会自动在两种实现间切换，缺一不可）。
- 前端任务列表工具栏新增"清空已结束"按钮：二次确认后调用接口，提示删除条数并刷新列表。

范围与约束：

- **不删除磁盘文件**：只清除数据库记录，已下载的视频文件保留在下载目录中。
- **不触碰未结束任务**：待执行（0）、执行中（1）、已暂停（4）的记录一律保留，因此不存在中断正在下载的任务的问题。
- 不做按状态筛选清空、不做软删除、不做清空历史回滚。

## Capabilities

### New Capabilities

- 无

### Modified Capabilities

- `http-api`: 新增"清空已结束任务接口"需求，定义 `POST /api/tasks/clear-finished` 的路由、响应体与行为。
- `task-storage`: 新增"按状态批量删除任务"需求，定义仓储层清除已完成/失败记录的行为。
- `wails-bindings`: 在"任务增删改查绑定"需求中补充 `ClearFinishedTasks()` 绑定，保持与 REST API 能力对等。
- `task-management-ui`: 新增"清空已结束任务按钮"需求，定义按钮位置、确认交互与结果提示。

## Impact

- **后端**：`internal/repository/task.go`（新增批量删除方法）、`internal/service/task.go`（新增 `ClearFinishedTasks`）、`internal/api/task.go`（新增 handler）、`main.go`（注册路由）。
- **桌面端绑定**：`app.go`（新增 `ClearFinishedTasks` 方法，需重新生成 wails 绑定）。
- **前端**：`web/src/api/task.ts`、`web/src/api/wails.ts`、`web/src/api/index.ts`（类型与实现已通过 `createApi` 代理统一，无需改动 index.ts）、`web/src/components/TaskList.vue`（工具栏按钮与处理函数）。
- **数据**：SQLite `tasks` 表中 status=2/3 的行被物理删除，不可恢复；下载目录文件不受影响。
- **无新增依赖**，无破坏性接口变更。
