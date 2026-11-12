## ADDED Requirements

### Requirement: 获取系统配置
系统 SHALL 提供 GET /api/config 接口返回当前系统配置。

#### Scenario: 成功获取配置
- **WHEN** 用户发送 GET /api/config 请求
- **THEN** 系统返回 200 状态码和完整配置 JSON（包含 max_concurrent_tasks、max_segment_workers、download_dir、max_consecutive_errors、default_headers、interceptor_enabled、agent_address、vpn_address、vpn_status、ffmpeg_download_enabled、gin_mode）

### Requirement: 更新系统配置
系统 SHALL 提供 PUT /api/config 接口用于更新系统配置。

#### Scenario: 成功更新配置
- **WHEN** 用户发送 PUT /api/config 请求，包含要更新的配置字段
- **THEN** 系统更新配置、持久化到 config.json、应用配置到 DownloadController、处理拦截器启停、返回 200 状态码和更新后的完整配置

#### Scenario: 部分更新配置
- **WHEN** 用户发送 PUT /api/config 请求，只包含部分字段
- **THEN** 系统只更新传入的字段，其他字段保持不变

#### Scenario: 更新 FFmpeg 直连下载开关
- **WHEN** 用户发送 PUT /api/config 请求，包含 `ffmpeg_download_enabled` 布尔字段
- **THEN** 系统更新该字段、持久化到 config.json，并将新值应用到 DownloadController，后续启动的 m3u8 任务按新值选择下载实现

#### Scenario: 未配置 FFmpeg 直连下载开关
- **WHEN** 配置文件或请求中不存在 `ffmpeg_download_enabled` 字段
- **THEN** 系统按默认值 `false` 处理，使用现有分片下载流程

#### Scenario: 无效配置值
- **WHEN** 用户发送 PUT /api/config 请求，包含无效值（如负数的并发数）
- **THEN** 系统返回 400 状态码和错误信息
