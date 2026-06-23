## 1. 核心实现

- [x] 1.1 在 main.go 中新增 `importTaskFile` 函数：读取 `task.txt`，按行分割并跳过空行，每两行组成一个任务（名称 + URL）
- [x] 1.2 根据 URL 后缀判定 Type：`.mp4` → `"mp4"`，其余 → `"m3u8"`
- [x] 1.3 构造 `model.Task`（Status=Pending，Header=配置中的 `default_headers`），通过 `repository.Create` 写入数据库；写入前检查 URL 是否已存在，存在则跳过
- [x] 1.4 导入完成后删除 `task.txt`，删除失败仅输出 warning 日志

## 2. 启动集成

- [x] 2.1 在 `main()` 中 `controller.GetController().ApplyConfig(...)` 之后、`svr.Init()` 之前调用 `importTaskFile`
- [x] 2.2 编译并手动测试：准备 `task.txt`，启动程序，验证任务写入数据库且文件被删除；再验证无文件时正常启动
