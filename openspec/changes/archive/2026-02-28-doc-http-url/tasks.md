## 1. 项目初始化

- [x] 1.1 添加 gin 和 gorm 依赖到 go.mod
- [x] 1.2 创建 internal 目录结构 (api, model, service, repository, controller)

## 2. 数据层实现

- [x] 2.1 创建 Task 模型 (internal/model/task.go)
- [x] 2.2 初始化 SQLite 数据库连接和自动迁移 (internal/repository/db.go)
- [x] 2.3 实现 TaskRepository CRUD 方法 (internal/repository/task.go)

## 3. 下载控制器实现

- [x] 3.1 定义 DTask 结构体，包含任务名、url、header、type、进度信息
- [x] 3.2 实现控制器添加任务方法，解析 url 和 header
- [x] 3.3 实现 MP4 断点续传下载逻辑
- [x] 3.4 实现 m3u8 分片下载逻辑（检查已下载分片）
- [x] 3.5 实现 ffmpeg 合并分片功能
- [x] 3.6 实现并发调度，为每个 DTask 启动 goroutine
- [x] 3.7 实现下载进度记录和查询

## 4. 业务层实现

- [x] 4.1 实现 TaskService 添加任务逻辑
- [x] 4.2 实现 TaskService 删除任务逻辑
- [x] 4.3 实现 TaskService 修改任务逻辑
- [x] 4.4 实现 TaskService 启动任务逻辑，调用下载控制器

## 5. API 层实现

- [x] 5.1 实现 POST /api/tasks 添加任务接口
- [x] 5.2 实现 DELETE /api/tasks/:id 删除任务接口
- [x] 5.3 实现 PUT /api/tasks/:id 修改任务接口
- [x] 5.4 实现 POST /api/tasks/start 执行任务接口
- [x] 5.5 实现 WebSocket /api/tasks/progress 进度推送接口（定时从控制器获取进度）

## 6. 集成与启动

- [x] 6.1 在 main.go 中初始化数据库和路由
- [x] 6.2 启动 Gin HTTP 服务
