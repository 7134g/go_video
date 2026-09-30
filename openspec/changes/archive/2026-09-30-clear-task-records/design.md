## Context

任务记录全部持久化在 SQLite 的 `tasks` 表中，前端每次通过 `GET /api/tasks` 或 Wails 的 `GetTasks()` 拉取列表。当前删除能力只有单条 `POST /api/tasks/delete`（`repository.Delete(id)`），删除已结束的历史记录需要逐条操作。

下载控制器 `DownloadController` 只在内存中持有正在执行的任务；任务进入终态时 `finishCallback` 会调用 `ctrl.RemoveTask(id)` 把任务移出内存表，因此已完成/失败的任务在内存中并不存在。

前端有一套双实现 API 层：`web/src/api/index.ts` 用 Proxy 在运行时检测 `window.go`，Wails 桌面端走 `web/src/api/wails.ts` 的绑定包装，Web 端走 `web/src/api/task.ts` 的 axios。新增的 API 方法必须在两处都实现，否则在另一种模式下会抛错。

## Goals / Non-Goals

**Goals:**

- 一条命令清除全部已完成和失败的任务记录。
- 保持 HTTP 与 Wails 两种运行模式能力对等。
- 不触碰未结束的任务，不中断任何正在进行的下载。
- 不动磁盘上的视频文件。

**Non-Goals:**

- 不做按状态筛选清空（如"只清失败"）、不做按时间范围清空。
- 不做软删除、回收站或撤销。
- 不删除下载目录中的文件（与 `redownload` 的语义明确区分）。
- 不做清空进度条 / 批量分批（记录量为几十到几百条，单条 DELETE 足够）。

## Decisions

### 1. 删除放在数据库层一次性完成，而不是查列表后逐条删

在 `TaskRepository` 新增：

```go
func (r *TaskRepository) DeleteFinished() (int64, error)
```

内部执行 `DB.Where("status IN ?", []model.TaskStatus{TaskStatusCompleted, TaskStatusFailed}).Delete(&model.Task{})`，返回 `RowsAffected`。

- **为什么**：与已有的 `ResetStatus` / `UpdateAllByStatus` 一致，都是单条条件批量语句，风格统一且无 N+1。
- **备选**：`GetAll()` 后过滤再逐条 `Delete(id)` —— 记录多时产生 N 次事务，且可能与并发状态变更竞争，放弃。

### 2. 状态集合在服务层显式定义为"已结束"，而不是复用 `GetPending` 的反集

`TaskService.ClearFinishedTasks()` 直接调用仓储方法。状态集合写死为 Completed + Failed。

- **为什么**："已结束"是业务语义，Pending/Running/Paused 都还有后续动作，不能删。用反集（`NOT IN (0,1,4)`）会让以后新增状态时默认被静默删除。
- **注意**：不加 `Running` 的运行时守卫。因为筛选条件本身已经排除 Running，加守卫只是多余的分支。

### 3. 端点命名 `POST /api/tasks/clear-finished`

- **为什么**：现有路由全部是 `POST /api/tasks/<动作>` 形式（`delete`、`start-one`、`stop-all`），保持一致。名字里带 `finished` 而非裸 `clear`，让破坏范围在 URL 上自解释。
- **备选**：`DELETE /api/tasks/finished` —— 语义更 RESTful，但与现有全 POST 风格不一致，前端 request 封装也要改判断，放弃。

### 4. 响应体 `{"deleted": N}`

- **为什么**：与 `POST /api/tasks/start` 返回 `{"started": count}` 保持一致，前端可直接用于提示文案"已清空 N 条记录"。

### 5. Wails 绑定返回删除条数

`func (a *App) ClearFinishedTasks() (int64, error)`，前端 `wails.ts` 包装成 `{ data: { deleted } }`，与 HTTP 版本同构。

- **为什么不返回 `error` 之外的包装结构**：Wails 绑定返回值经 JSON 序列化后与 axios 的 `response.data` 对齐即可，无需额外类型。

### 6. 前端按钮放在工具栏，带二次确认

在 `TaskList.vue` 工具栏"暂停全部"之后新增 `<el-button type="danger" plain :loading="clearing" @click="handleClearFinished">清空已结束</el-button>`，处理函数沿用现有 `ElMessageBox.confirm` + `ElMessage.success` + `loadTasks()` 模式。

- **为什么用 `plain` 而非实心 danger**：工具栏已有"暂停全部"用 warning 实心；清空是低频破坏性操作，弱化视觉权重降低误点概率。
- **确认文案需说明不删文件**：用户对"清空"的直觉预期包含删文件，必须在对话框里显式澄清。

## Risks / Trade-offs

- [误删已完成记录后无法恢复] → 二次确认对话框明确写出删除范围与条数不可恢复；不删除任何文件，重新下载成本可控。
- [用户以为"清空"会连文件一起删，随后手动去下载目录删文件] → 确认对话框和成功提示中明确"文件已保留"。
- [代理拦截（MITM）在清空的同时写入新任务] → 删除条件是 status 快照，新写入的任务是 Pending，不在删除范围内；最多出现本次未删净的已完成记录，用户再次点击即可。
- [Wails 绑定未重新生成导致桌面端按钮报错] → 实施步骤中包含重新生成绑定（`wails generate module` 或按现有构建流程），并在桌面端模式下验证按钮可用。
- [清空大量记录时阻塞请求] → SQLite 单条 DELETE 在几百行量级为毫秒级，不做异步任务化。

## Migration Plan

无数据迁移。接口为纯新增，旧版本前端不会调用。回滚方式为移除按钮与路由，已删除的记录无法恢复。

## Open Questions

无。
