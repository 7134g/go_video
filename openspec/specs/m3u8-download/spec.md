# M3U8 Download Capability

## Requirements

### Requirement: M3U8 文件解析
系统 SHALL 使用 `common/m3u8.ParseM3u8Data()` 解析 M3U8 文件内容，支持 VOD 和 EVENT 类型播放列表。

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

### Requirement: 并发分片下载
系统 SHALL 并发下载多个分片，并限制最大并发数。

#### Scenario: 并发下载
- **WHEN** 开始下载包含 100 个分片的 M3U8
- **THEN** 系统同时下载不超过配置的最大并发数的分片

#### Scenario: 连续错误中断
- **WHEN** 连续下载失败次数达到阈值
- **THEN** 系统停止下载并返回错误

### Requirement: 分片文件命名
系统 SHALL 按顺序命名分片文件，格式为 `%05d.ts`。

#### Scenario: 分片存储
- **WHEN** 下载第 1 个分片
- **THEN** 保存为 `00000.ts`

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
- **THEN** 系统下载 Key 文件并缓存

#### Scenario: 解密分片
- **WHEN** 分片关联了加密密钥
- **THEN** 系统使用 AES-128-CBC 解密分片内容

### Requirement: 分片合并
系统 SHALL 在所有分片下载完成后自动合并为单个视频文件。

#### Scenario: 使用 ffmpeg 合并
- **WHEN** ffmpeg 可用
- **THEN** 使用 ffmpeg concat 方式合并分片为 MP4

#### Scenario: 直接合并
- **WHEN** ffmpeg 不可用
- **THEN** 直接拼接 TS 分片为 MP4 文件

### Requirement: 进度报告
系统 SHALL 报告下载进度。

#### Scenario: 分片进度
- **WHEN** 完成一个分片下载
- **THEN** 更新进度为 (已完成分片数/总分片数)
