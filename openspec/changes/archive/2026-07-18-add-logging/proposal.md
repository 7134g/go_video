## Why

程序目前没有任何日志记录，运行时缺乏可见性。下载失败、代理异常、配置错误等问题发生时无法排查。需要通过日志系统记录关键操作和异常，方便问题定位和运行监控。

## What Changes

- 引入基于 `log/slog` 的日志系统，按天分割日志文件
- 日志写入程序所在目录下的 `logs/` 文件夹
- 自动清理超过 30 天的旧日志文件
- Go 主程序、MITM 代理、下载器等各模块接入日志
- 保留控制台输出（开发时双写终端和文件）
- 新增日志级别配置项（config.json 可选），默认 info

## Capabilities

### New Capabilities
- `log-system`: 按天分割的日志文件写入，30 天自动清理，支持 slog 结构化日志

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- 新增依赖：无（使用 Go 标准库 `log/slog`，Go 1.21+ 内置）
- `internal/` 各层需注入或通过全局 logger 记录日志
- `config.json` 增加可选字段 `log_level`
- 主程序启动时初始化日志系统
