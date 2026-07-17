## Context

当前项目使用 Go 标准库 `log` 在少数几个地方直接输出到控制台（`internal/api/middleware.go` 等），没有统一的日志框架。需要引入日志系统覆盖所有模块：API、Service、Controller、Downloader、Proxy。

## Goals / Non-Goals

**Goals:**
- 所有模块通过统一 logger 记录日志
- 日志写入 `logs/` 目录，按天分割文件
- 自动清理超过 30 天的日志文件
- 支持日志级别（debug/info/warn/error）
- config.json 可配日志级别
- 开发时同时输出到终端和文件

**Non-Goals:**
- 不引入第三方日志库（使用 Go 标准库 `log/slog`，Go 1.21+）
- 不做日志轮转的压缩归档
- 不添加远程日志收集

## Decisions

1. **使用 `log/slog` 而非第三方库**
   - Go 1.21 标准库，零依赖
   - 结构化日志，支持 JSON 和 text 格式
   - 性能优于 `log` 包
   - 项目中 Go 版本 `go.mod` 为 1.21+，无需升级

2. **按天分割 + 30 天清理**
   - 每天一个文件，文件名 `logs/2006-01-02.log`
   - 使用 `os.Stat` 检查 mtime，在每次写日志时惰性清理；或在初始化时启动一个 goroutine 每天定时清理一次
   - 选择惰性清理（写日志前检查），没有额外 goroutine 开销

3. **全局 Logger 变量**
   - `internal/logger` 包暴露 `Init()` 和全局 `Logger` 变量
   - 各包通过 `logger.Info(...)` / `logger.Error(...)` 调用
   - 避免依赖注入的改造成本

4. **双写终端和文件**
   - `slog.Handler` 支持 `io.MultiWriter`
   - 开发时同时输出彩色文本到 stderr 和 JSON 到文件
   - 生产（无终端）时只写文件

## Risks / Trade-offs

- 惰性清理可能在日志流量低时延迟删除过期文件 → 可接受，最晚 1 天内会清理
- 全局变量不利于测试隔离 → 测试时可替换为 `slog.DiscardHandler`
- 文件日志不支持并发切割 → 同步写入通过 `slog.Handler` 内部锁保证安全
