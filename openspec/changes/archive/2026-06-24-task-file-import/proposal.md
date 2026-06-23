## Why

用户需批量导入一组视频任务，手动逐个创建效率太低。程序启动时自动读取 `task.txt` 文件并导入，可以快速建立任务队列，省去重复操作。

## What Changes

- 程序启动完成时（DB 初始化、CA 检查之后），读取工作目录下的 `task.txt` 文件
- 按行解析：奇数行为任务名称，偶数行为 URL，每两行组成一个任务
- 根据 URL 后缀（`.m3u8` / `.mp4`）自动判断任务 Type，其余后缀按 M3U8 处理
- 使用配置中的默认 Header 写入数据库
- `task.txt` 解析完毕后自动删除，避免下次启动重复导入

## Capabilities

### New Capabilities
- `task-file-import`: 启动时从 `task.txt` 文件批量导入任务到数据库

### Modified Capabilities
<!-- 无现有规格需要变更 -->

## Impact

- `main.go`：启动流程中新增读取解析步骤
- `internal/service/`：可能新增解析逻辑或直接在 `main.go` 中内联处理
- `model.Task` / `repository`：复用现有数据库写入路径
