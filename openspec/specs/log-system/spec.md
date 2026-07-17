## ADDED Requirements

### Requirement: 按天分割日志文件
系统 SHALL 将日志写入程序所在目录下的 `logs/` 文件夹中，按天分割为单独的日志文件，文件名格式为 `YYYY-MM-DD.log`。

#### Scenario: 每天创建一个新的日志文件
- **WHEN** 当前日期变更后首次写入日志
- **THEN** 系统创建一个新的 `logs/YYYY-MM-DD.log` 文件

#### Scenario: 同一天的日志写入同一文件
- **WHEN** 在同一天内多次写入日志
- **THEN** 所有日志追加到当天的同一个日志文件中

### Requirement: 30 天日志自动清理
系统 MUST 自动删除超过 30 天的日志文件，以 `logs/` 目录下文件的最后修改时间为准。

#### Scenario: 启动时清理过期日志
- **WHEN** 程序启动时
- **THEN** 系统遍历 `logs/` 目录，删除 mtime 超过 30 天的文件

#### Scenario: 写入日志时惰性清理
- **WHEN** 程序运行时写入日志
- **THEN** 系统检查是否已有过期文件，如有则清理

#### Scenario: 保留 30 天内的日志
- **WHEN** 日志文件修改时间在 30 天内
- **THEN** 系统不删除该文件

### Requirement: 日志级别配置
系统 SHALL 支持通过 `config.json` 的 `log_level` 字段配置日志级别，可选值：`debug`、`info`、`warn`、`error`，默认为 `info`。

#### Scenario: 使用默认日志级别
- **WHEN** 用户未配置 `log_level`
- **THEN** 系统使用 `info` 级别

#### Scenario: 配置日志级别
- **WHEN** 用户在 `config.json` 中设置 `log_level: "debug"`
- **THEN** 系统输出 debug 及以上级别的日志

### Requirement: 终端与文件双写
系统 SHALL 在开发时同时输出日志到终端（彩色文本）和日志文件（JSON 格式）。

#### Scenario: 开发模式双写
- **WHEN** 程序有可用的终端（stderr）
- **THEN** 日志同时输出到 stderr 和当日日志文件

#### Scenario: 后台运行仅写文件
- **WHEN** 程序无可用终端（如后台服务）
- **THEN** 日志仅写入文件

### Requirement: 各模块接入日志
系统 SHALL 在以下模块中记录关键操作和错误：API 处理器、TaskService、ConfigService、DownloadController、Downloader Pool、MITM Proxy、M3U8 解析/解密。

#### Scenario: 下载任务状态变更记录
- **WHEN** 下载任务启动、完成、失败或暂停
- **THEN** 系统记录包含任务 ID、URL、状态的事件日志

#### Scenario: 代理拦截记录
- **WHEN** MITM 代理拦截到视频 URL
- **THEN** 系统记录已拦截的 URL 和目标类型

#### Scenario: 错误日志
- **WHEN** 下载失败、网络异常、解析错误
- **THEN** 系统记录包含错误详情和上下文的 error 级别日志
