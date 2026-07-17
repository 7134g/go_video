## 1. Logger 包（internal/logger）

- [x] 1.1 创建 `internal/logger/logger.go`，实现按天分割的文件 writer：`NewDateRotatingWriter(dir string)`，在 `logs/` 目录下按 `2006-01-02.log` 格式创建文件，当日切换时自动关闭旧文件打开新文件
- [x] 1.2 实现 `cleanOldLogs(dir string, days int)` 函数，遍历 `logs/` 目录删除 mtime 超过 30 天的文件；启动时和每次写日志前惰性调用
- [x] 1.3 实现 `Init(level slog.Level, logDir string)` 函数，创建 `slog.TextHandler`（终端，彩色）和 `slog.JSONHandler`（文件）的 `io.MultiWriter` 组合，`Init` 后替换全局 `slog.SetDefault`
- [x] 1.4 实现 `InitWithLevel(levelStr string, logDir string)` 包装函数，解析 `debug/info/warn/error` 字符串为 `slog.Level`

## 2. Config 集成

- [x] 2.1 在 `model.Config` 中增加 `LogLevel string \`json:"log_level"\`` 字段，`DefaultConfig` 中默认 `""`（等价于 info）
- [x] 2.2 在 `shared.go` 的 `initShared` 中调用 `logger.InitWithLevel(cfg.LogLevel, pwd+"/"+logger.LogDir)`，确保日志系统在 Gin 启动前初始化

## 3. 替换各模块 log 调用

- [x] 3.1 `internal/api/` 中无 `log` 调用，无需修改
- [x] 3.2 替换 `internal/service/config.go` 中的 `fmt.Print*` 为 `slog`
- [x] 3.3 替换 `internal/controller/controller.go` 和 `m3u8.go` 中的 `log` 调用为 `slog`
- [x] 3.4 替换 `pkg/proxy/collector.go` 和 `server.go` 中的 `fmt.Print*` 为 `slog`
- [x] 3.5 `pkg/m3u8/ffmpeg.go` 仅含注释掉的 `fmt.Println`，无需修改
- [x] 3.6 `pkg/downloader/pool.go` 中无 `log` 调用，无需修改

## 4. 替换 main.go 和 shared.go 中的 log 调用

- [x] 4.1 替换 `main.go` 中的 `log.Fatal`/`log.Println` 为 `slog.Error`/`slog.Info` + `os.Exit(1)`
- [x] 4.2 替换 `shared.go` 中 `importTaskFile` 的 `log.Print*` 为 `slog`；`InitCa` 和 `initShared` 前置 `log.Fatal` 保留（日志系统初始化前执行）

## 5. 验证

- [x] 5.1 编译项目：`go build ./...` ✓
- [x] 5.2 运行测试：`go test ./...` ✓
- [x] 5.3 确认 `logs/` 目录按预期创建日志文件 ✓
