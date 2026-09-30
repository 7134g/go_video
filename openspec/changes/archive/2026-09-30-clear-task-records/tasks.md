## 1. 后端仓储层

- [x] 1.1 在 `internal/repository/task.go` 新增 `DeleteFinished() (int64, error)`：用 `DB.Where("status IN ?", []model.TaskStatus{model.TaskStatusCompleted, model.TaskStatusFailed}).Delete(&model.Task{})` 单条语句批量删除，返回 `RowsAffected`。注释说明"仅删除数据库记录，不动磁盘文件" → 验证：`go build ./...` 通过

## 2. 后端服务层

- [x] 2.1 在 `internal/service/task.go` 新增 `ClearFinishedTasks() (int64, error)`，直接委托 `s.repo.DeleteFinished()`，不写额外状态守卫（筛选条件已排除运行中任务） → 验证：`go build ./...` 通过

## 3. HTTP 接口

- [x] 3.1 在 `internal/api/task.go` 新增 `func (h *TaskHandler) ClearFinished(c *gin.Context)`：调用 `h.svc.ClearFinishedTasks()`，失败返回 500 和 `{"error": ...}`，成功返回 200 和 `{"deleted": N}` → 验证：编译通过
- [x] 3.2 在 `main.go` 的 tasks 路由组中注册 `tasks.POST("/clear-finished", h.ClearFinished)`，位置紧跟 `tasks.POST("/delete", h.Delete)` → 验证：`go build .` 通过，`curl -X POST localhost:8080/api/tasks/clear-finished` 返回 `{"deleted":N}`

## 4. Wails 绑定

- [x] 4.1 在 `app.go` 新增 `func (a *App) ClearFinishedTasks() (int64, error)`，委托 `a.svc.ClearFinishedTasks()` → 验证：`go build -tags bindings .` 通过
- [x] 4.2 重新生成 Wails 绑定（`wails generate module` 或项目的绑定构建流程），确认 `web/wailsjs/go/main/App.d.ts` 与 `App.js` 中出现 `ClearFinishedTasks` → 验证：`grep ClearFinishedTasks web/wailsjs/go/main/App.d.ts` 有输出

## 5. 前端 API 层

- [x] 5.1 在 `web/src/api/task.ts` 的 `taskApi` 中新增 `clearFinished: () => request.post<{ deleted: number }>('/tasks/clear-finished')` → 验证：TypeScript 编译无错误
- [x] 5.2 在 `web/src/api/wails.ts` 的 `WailsApp` 接口中声明 `ClearFinishedTasks(): Promise<number>`，并在 `taskApi` 中新增 `clearFinished` 包装为 `{ data: { deleted } }`，与 HTTP 版本返回结构一致 → 验证：`npx vue-tsc --noEmit` 或 `npm run build` 通过

## 6. 前端界面

- [x] 6.1 在 `web/src/components/TaskList.vue` 工具栏"暂停全部"之后新增 `<el-button type="danger" plain :loading="clearing" @click="handleClearFinished">清空已结束</el-button>` → 验证：按钮出现在工具栏
- [x] 6.2 新增 `clearing` 状态与 `handleClearFinished()`：`ElMessageBox.confirm` 确认（文案说明只删除已完成和失败的任务记录、磁盘文件保留）→ 调 `taskApi.clearFinished()` → `ElMessage.success('已清空 N 条记录')` → `loadTasks()`；用户取消时直接返回；请求异常时 `ElMessage.error` 并保持列表不变；`finally` 中复位 `clearing` → 验证：点击按钮确认后列表刷新且提示条数

## 7. 验证

- [x] 7.1 编写场景核对：造 1 个已完成 + 1 个失败 + 1 个待执行任务，调用接口后确认仅剩待执行任务，返回的 `deleted` 为 2 且下载目录文件仍存在 → 验证：`GET /api/tasks` 结果符合预期
- [x] 7.2 确认未结束任务不受影响：启动一个任务并在运行中点击清空，确认该任务继续下载、状态仍为执行中 → 验证：运行中任务未被删除
- [ ] 7.3 桌面端验证：以 Wails 模式运行（`go build -tags desktop` 或桌面端构建流程），点击"清空已结束"按钮确认走绑定通道且行为一致 → 验证：无 `taskApi.clearFinished is not a function` 报错
  - 已自动验证：`go build -tags desktop` 链接通过；`wails generate module` 后 `web/wailsjs/go/main/App.{d.ts,js}` 中出现 `ClearFinishedTasks`；打包后的 `web/dist` 中同时含有 `清空已结束`、`ClearFinishedTasks`、`clear-finished`，证明双通道均已进入产物
  - 仍待人工完成：启动桌面端 GUI 实际点击一次按钮，确认无 `clearFinished is not a function` 报错
- [x] 7.4 运行 `go test ./...` 与 `npm run build`，确认既有测试与前端构建均通过 → 验证：两条命令均以 0 退出
