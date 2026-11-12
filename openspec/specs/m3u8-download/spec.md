# M3U8 Download Capability

## Requirements

### Requirement: M3U8 文件解析
系统 SHALL 使用 `go_video/pkg/m3u8.ParseM3u8Data()` 解析 M3U8 文件内容。

#### Scenario: 解析标准 M3U8 文件
- **WHEN** 提供有效的 M3U8 URL
- **THEN** 系统返回包含所有分片信息的 M3u8 结构体

#### Scenario: 解析无效 M3U8 文件
- **WHEN** M3U8 文件缺少 #EXTM3U 头
- **THEN** 系统返回解析错误

### Requirement: Master Playlist 处理
系统 SHALL 检测 Master Playlist 并自动选择最高码率的媒体播放列表。

#### Scenario: 处理 Master Playlist
- **WHEN** M3U8 文件包含 EXT-X-STREAM-INF 标签
- **THEN** 系统选择 BANDWIDTH 最大的流并递归获取其媒体播放列表

#### Scenario: 处理普通媒体播放列表
- **WHEN** M3U8 文件不包含 EXT-X-STREAM-INF 标签
- **THEN** 系统直接解析分片列表

### Requirement: 分片 URL 解析
系统 SHALL 将分片的相对 URL 转换为绝对 URL。

#### Scenario: 相对路径分片
- **WHEN** 分片 URI 为相对路径（如 "segment001.ts"）
- **THEN** 系统基于 M3U8 文件 URL 构建完整的分片 URL

#### Scenario: 绝对路径分片
- **WHEN** 分片 URI 为绝对 URL
- **THEN** 系统直接使用该 URL

### Requirement: 按域名并发分片下载
系统 SHALL 通过 `downloader.Pool` 按域名限制并发下载分片数（`MaxSegmentWorkers`）。

#### Scenario: 并发下载
- **WHEN** 开始下载包含 100 个分片的 M3U8
- **THEN** 系统通过 Group 提交分片任务，同域名下同时下载不超过 MaxSegmentWorkers 个分片

#### Scenario: 连续错误中断
- **WHEN** 连续下载失败次数达到 MaxConsecutiveErrors
- **THEN** 系统停止下载并返回错误

### Requirement: 分片文件命名
系统 SHALL 按顺序命名分片文件，格式为 `%06d.ts`。

#### Scenario: 分片存储
- **WHEN** 下载第 1 个分片
- **THEN** 保存为 `000000.ts`

#### Scenario: 分片目录
- **WHEN** 开始下载任务
- **THEN** 在 `<downloadDir>/<taskName>/` 目录下存储所有分片

### Requirement: 断点续传
系统 SHALL 支持分片级别的断点续传。

#### Scenario: 跳过已下载分片
- **WHEN** 分片文件已存在
- **THEN** 跳过该分片的下载

### Requirement: AES-128 加密支持
系统 SHALL 支持 AES-128 加密分片的下载和解密。

#### Scenario: 下载加密密钥
- **WHEN** M3U8 包含 EXT-X-KEY 标签且 METHOD=AES-128
- **THEN** 系统下载 Key 文件并缓存在 KeyCache 中

#### Scenario: 解密分片
- **WHEN** 分片关联了加密密钥
- **THEN** 系统使用 AES-128-CBC 解密分片内容（IV 基于 MediaSequence + KeyIndex 推导）

### Requirement: 分片合并
系统 SHALL 在所有分片下载完成后使用 ffmpeg 合并为单个视频文件。Windows 使用 `ffmpeg.exe`，其他平台使用 `ffmpeg`。

#### Scenario: 使用 ffmpeg 合并
- **WHEN** 所有分片下载完成
- **THEN** 调用项目根目录下的 ffmpeg 合并分片为 MP4，合并后删除分片目录

### Requirement: 进度报告
系统 SHALL 报告下载进度。

#### Scenario: 分片进度
- **WHEN** 完成一个分片下载
- **THEN** 更新进度为 (已完成分片数/总分片数)

### Requirement: FFmpeg 直连下载模式

系统 SHALL 在配置项 `ffmpeg_download_enabled` 为 `true` 时，将 M3U8 任务交由 ffmpeg 直接下载；为 `false` 时 MUST 保持现有的分片下载与合并流程不变。

#### Scenario: 开关开启时执行直连下载

- **WHEN** 任务类型为 m3u8 且 `ffmpeg_download_enabled` 为 `true`
- **THEN** 系统以任务 URL 作为 ffmpeg 输入，使用 `-c copy` 直接输出到 `<download_dir>/<任务名>.mp4`，并跳过分片下载、解密与合并阶段

#### Scenario: 开关关闭时保持原有流程

- **WHEN** 任务类型为 m3u8 且 `ffmpeg_download_enabled` 为 `false`
- **THEN** 系统使用现有的分片下载、AES-128 解密和 ffmpeg 合并流程

#### Scenario: 附加请求头

- **WHEN** 任务 Header 非空
- **THEN** 系统将任务 Header 以 `Key: Value\r\n` 拼接后作为 ffmpeg `-headers` 参数传入

#### Scenario: 使用默认请求头

- **WHEN** 任务 Header 为空
- **THEN** 系统使用配置中的 `default_headers` 拼装 ffmpeg `-headers` 参数

#### Scenario: 附加 HTTP 代理

- **WHEN** `vpn_status` 为 `true` 且 `vpn_address` 非空
- **THEN** 系统向 ffmpeg 传入 `-http_proxy http://<vpn_address>` 参数

#### Scenario: 不附加 HTTP 代理

- **WHEN** `vpn_status` 为 `false` 或 `vpn_address` 为空
- **THEN** 系统不向 ffmpeg 传入 `-http_proxy` 参数

#### Scenario: ffmpeg 下载失败

- **WHEN** ffmpeg 进程以非零退出码结束或启动失败
- **THEN** 系统使用 `slog` 记录错误信息与 ffmpeg stderr，并将任务置为失败终态

#### Scenario: 停止直连下载任务

- **WHEN** 用户停止一个正在由 ffmpeg 直连下载的任务
- **THEN** 系统取消任务 context 以终止 ffmpeg 进程，并按 `context.Canceled` 处理该任务

#### Scenario: 直连模式进度

- **WHEN** 直连下载开始与结束
- **THEN** 系统在开始时将进度置为 0/1，成功后置为 1/1，不报告分片粒度的中间进度
