## 1. 配置模型与存储

- [x] 1.1 创建 `internal/model/config.go`，定义 Config 结构体
- [x] 1.2 创建 `internal/repository/config.go`，实现配置文件读写（JSON）

## 2. 配置服务层

- [x] 2.1 创建 `internal/service/config.go`，实现配置服务单例
- [x] 2.2 实现 GetConfig 方法返回当前配置
- [x] 2.3 实现 UpdateConfig 方法更新并持久化配置
- [x] 2.4 实现配置校验逻辑（正数校验、路径有效性）

## 3. API 接口

- [x] 3.1 创建 `internal/api/config.go`，实现 ConfigHandler
- [x] 3.2 实现 GET /api/config 处理函数
- [x] 3.3 实现 PUT /api/config 处理函数
- [x] 3.4 在 main.go 中注册配置路由

## 4. 集成下载控制器

- [x] 4.1 修改 DownloadController 从配置服务读取 downloadDir
- [x] 4.2 修改 m3u8 下载逻辑使用 max_segment_workers 配置
- [x] 4.3 修改任务启动逻辑使用 max_concurrent_tasks 配置
- [x] 4.4 实现连续错误计数和 max_consecutive_errors 检查
