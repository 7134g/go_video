## 1. Controller 重构

- [x] 1.1 `PauseTask` 重命名为 `StopTask`，同步更新 `internal/service/task.go` 中调用
- [x] 1.2 `AddTask` 增加 `BroadcastMessage(id, "任务已添加: "+name)` 调用
- [x] 1.3 `StartTask` 增加 `BroadcastMessage(id, "任务已启动: "+task.Name)` 调用
- [x] 1.4 `StartAll` 增加 `BroadcastMessage(0, "已启动 N 个任务")` 调用
- [x] 1.5 `StopTask` 增加 `BroadcastMessage(id, "任务已停止: "+task.Name)` 调用
- [x] 1.6 `RemoveTask` 增加 `BroadcastMessage(id, "任务已删除")` 调用（需先获取 task 名称）
- [x] 1.7 新增 `AddAndStart(id, name, url, headerJSON, taskType) error` 方法：调用 AddTask 后立即 StartTask，广播"任务已添加并启动: <name>"

## 2. Service 层

- [x] 2.1 新增 `StartTask(id uint) error`：从 repo 获取任务 → 合并 header → AddTask → StartTask，启动前检查状态
- [x] 2.2 新增 `AddAndStart(task *model.Task) error`：先 Create 入库，再 AddTask + StartTask
- [x] 2.3 `PauseTask` 调用改为 `StopTask`

## 3. API 路由与 Handler

- [x] 3.1 `main.go` 路由变更：`DELETE /:id` → `POST /delete`，`PUT /:id` → `POST /update`，`POST /:id/pause` → `POST /pause`，`POST /:id/retry` → `POST /retry`，新增 `POST /start-one`
- [x] 3.2 `Delete` handler 改为从 body `{"id": N}` 解析 id
- [x] 3.3 `Update` handler 改为从 body `{"id": N, ...}` 解析 id
- [x] 3.4 `Pause` handler 改为从 body `{"id": N}` 解析 id
- [x] 3.5 `Retry` handler 改为从 body `{"id": N}` 解析 id
- [x] 3.6 新增 `StartOne` handler：从 body `{"id": N}` 解析，调用 `svc.StartTask(id)`

## 4. 前端

- [x] 4.1 `web/src/api/task.ts` 更新所有接口：delete/update/pause/retry 改为 POST + body，新增 startOne
- [x] 4.2 `TaskList.vue` 操作列新增"启动"按钮，待执行/失败/已暂停状态显示
- [x] 4.3 `handlePause`/`handleRetry`/`handleDelete`/`handleEdit` 调用方式适配新 API
- [x] 4.4 新增 `handleStartOne(id)` 方法调用 `taskApi.startOne(id)`

## 5. 验证

- [x] 5.1 `go build .` 编译通过
- [x] 5.2 前端 `npm run build` 构建通过（web/ 目录下）
- [ ] 5.3 手动测试：创建任务 → 单启 → 暂停 → 重试 → 删除，检查 WebSocket 日志区有操作输出
