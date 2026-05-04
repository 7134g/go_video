## 为什么需要这个改动

当前用户只能通过"开始下载"按钮批量启动所有待执行任务，无法单独启动某一个任务，操作粒度太粗。同时 `internal/controller` 中的任务生命周期管理代码职责不够清晰，操作（添加、启动、停止、删除）缺少统一的前端日志输出，用户无法在 WebSocket 日志区看到操作反馈。

## 改什么

- **新增单个任务启动按钮**：在 `TaskList.vue` 操作列增加"启动"按钮，对未开始/已暂停/失败状态的任务可点击启动
- **新增单个任务启动 API**：`POST /api/tasks/start-one`，请求体 `{"id": N}`，用于启动指定任务
- **新增「添加并启动」功能**：在 controller/service 层提供一步完成添加+启动的能力
- **重构 controller 任务生命周期方法**：统一整理 `DownloadController` 中的任务操作方法，包括添加任务、添加并启动、停止任务、删除任务、启动所有任务
- **操作日志广播**：上述所有操作均通过 `BroadcastMessage` 输出到前端 WebSocket 日志区
- **统一 API 风格**：所有写操作接口统一使用 `POST` + JSON body 传参，去掉 URL 路径参数（`:id`），与 `POST /api/tasks` 和 `POST /api/tasks/start` 保持一致

## 能力

### 新增能力
- `single-task-start`: 用户可以通过 API 和 UI 单独启动某一个任务，而非只能批量启动
- `task-lifecycle`: 下载控制器的任务生命周期管理（添加、添加并启动、停止、删除、批量启动），所有操作通过 broadcast 输出日志

### 修改的能力
- `download-controller`: 任务操作方法需要统一重构，并在操作时调用 BroadcastMessage 输出
- `task-management-ui`: 操作列新增"启动"按钮
- `http-api`: 新增 `POST /api/tasks/start-one` 单个任务启动接口

## 影响范围

- `web/src/components/TaskList.vue` — 操作列新增单个启动按钮
- `web/src/api/task.ts` — 所有接口改为 POST + body 传参
- `internal/api/task.go` — 新增 `StartOne` handler，现有 handler 改为从 body 解析 id
- `main.go` — 路由全部改为无路径参数的 POST 形式
- `internal/service/task.go` — 新增 `StartTask`、`AddAndStart` 方法
- `internal/controller/controller.go` — 重构任务生命周期方法，增加 broadcast 调用
