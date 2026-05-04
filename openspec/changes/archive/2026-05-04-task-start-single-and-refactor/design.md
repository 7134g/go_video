## 背景

当前系统已具备任务的基本 CRUD 和批量启动能力，但缺少单个任务启动功能。`DownloadController` 中任务生命周期方法（AddTask、RemoveTask、PauseTask、StartAll、StartTask）已经存在，但命名不统一（PauseTask 实际是停止），且都没有调用 `BroadcastMessage` 向前端推送操作日志。用户在前端看不到操作反馈。

调用链: `API handler → TaskService → DownloadController`

## 目标 / 非目标

**目标:**
- 提供单个任务启动的 API 和 UI
- 统一 controller 任务生命周期方法命名
- 所有任务操作向前端广播日志
- 新增「添加并启动」组合方法

**非目标:**
- 不改变下载引擎逻辑（m3u8/mp4 下载流程不变）
- 不改变 WebSocket 进度推送机制
- 不涉及任务状态机变更

## 决策

### 1. API 风格统一: 全部 POST + body，去掉 URL 路径参数

所有写操作接口统一使用 `POST` + JSON body 传参，不再使用 URL 路径参数（`:id`）。与已有的 `POST /api/tasks` 和 `POST /api/tasks/start` 风格一致。

路由变更对照：

| 原路由 | 新路由 | 请求体 |
|--------|--------|--------|
| `DELETE /:id` | `POST /delete` | `{"id": N}` |
| `PUT /:id` | `POST /update` | `{"id": N, "name": "...", ...}` |
| `POST /:id/pause` | `POST /pause` | `{"id": N}` |
| `POST /:id/retry` | `POST /retry` | `{"id": N}` |
| (新增) | `POST /start-one` | `{"id": N}` |
| `POST ""` (create) | 不变 | — |
| `POST /start` (batch) | 不变 | — |
| `GET ""` (list) | 不变 | — |
| `GET /progress` (WS) | 不变 | — |

### 2. Controller 方法统一命名

| 方法 | 原名 | 说明 |
|------|------|------|
| `AddTask` | 不变 | 将任务加入内存 map |
| `AddAndStart` | **新增** | 添加任务并立即启动 |
| `StartTask` | 不变 | 启动指定任务 |
| `StartAll` | 不变 | 并发启动所有内存中的任务 |
| `StopTask` | `PauseTask` | 重命名，语义更准确（cancel context） |
| `RemoveTask` | 不变 | 从内存 map 中删除 |

`PauseTask` 改名为 `StopTask` 的理由：当前实现是调用 `task.cancel()` 取消 context，本质是停止而非暂停（暂停通常意味着可恢复）。下游 service 层收到 `context.Canceled` 后将数据库状态设为 "已暂停"，但那是数据库层面的语义，controller 层的方法名应反映实际操作。

### 3. Broadcast 调用位置

放在 controller 层的方法中，不在 service 层。因为：
- `broadcast.go` 与 controller 在同一包内，无需引入额外依赖
- Service 层职责是协调 controller 和 repository，不应关心前端日志细节

每个操作方法在关键节点调用 `BroadcastMessage(taskID, msg)`：
- `AddTask` → "任务已添加: <name>"
- `AddAndStart` → "任务已添加并启动: <name>"
- `StartTask` → "任务已启动: <name>"
- `StartAll` → "已启动 N 个任务"
- `StopTask` → "任务已停止: <name>"
- `RemoveTask` → "任务已删除: <name>"

### 4. Service 层新增方法

```go
// StartTask 启动单个任务（含 header 合并）
func (s *TaskService) StartTask(id uint) error

// AddAndStart 添加任务并立即启动（用于代理拦截等场景）
func (s *TaskService) AddAndStart(task *model.Task) error
```

`StartTask` 与现有 `StartTasks`（批量）平级，逻辑从数据库读取任务、合并 header、调用 controller 的 AddTask + StartTask。

### 5. 前端按钮显示逻辑

操作列按状态展示不同按钮：

| 状态 | 显示的按钮 |
|------|-----------|
| 待执行 (0) | 启动、编辑、删除 |
| 执行中 (1) | 暂停 |
| 完成 (2) | 删除 |
| 失败 (3) | 启动、重试、删除 |
| 已暂停 (4) | 启动、重试、删除 |

新增 `handleStartOne(id)` 方法调用 `POST /api/tasks/start-one`，请求体 `{"id": id}`。

### 6. 前端 API 客户端

`task.ts` 新增: `startOne: (id: number) => request.post<{ started: number }>('/tasks/start-one', { id })`

## 风险 / 权衡

- **并发安全**: `StartTask` 会起 goroutine 执行 `runTask`，与 `StartAll` 共享 `c.tasks` map，已有 `sync.RWMutex` 保护，无需额外处理
- **重复启动**: 如果对一个已在运行的任务调用 `StartTask`，当前实现会再起一个 goroutine（重复下载）。Service 层应在启动前检查任务状态，只允许 status 为待执行/失败/已暂停的任务启动
- **broadcast channel 满**: `BroadcastMessage` 使用非阻塞 `select default`，当日志 channel 满时消息会被丢弃，不会阻塞下载流程。这是可接受的行为
