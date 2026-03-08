## Context

当前 `internal/controller/m3u8.go` 使用简单的行解析获取分片 URL，没有利用 `common/m3u8` 包的完整解析能力。现有代码无法处理：
- Master Playlist（多码率选择）
- AES-128 加密分片
- 字节范围请求（EXT-X-BYTERANGE）

`common/m3u8` 包已提供完整的 M3U8 解析结构和合并功能，需要在 controller 层正确调用。

## Goals / Non-Goals

**Goals:**
- 使用 `common/m3u8.ParseM3u8Data()` 解析 M3U8 文件
- 支持 Master Playlist 自动选择最高码率
- 支持 AES-128 加密分片的下载和解密
- 保持现有的并发下载和断点续传能力
- 下载完成后自动合并分片

**Non-Goals:**
- 不支持实时流（Live HLS）
- 不支持 SAMPLE-AES 加密
- 不修改 MP4 下载逻辑

## Decisions

### 1. 解析流程
使用 `common/m3u8.ParseM3u8Data()` 替代手动行解析。如果返回 Master Playlist，递归获取最高码率的媒体播放列表。

**理由**: 复用已有代码，支持完整的 M3U8 规范。

### 2. 加密处理
在 `common/m3u8` 包中新增 `Decrypt()` 函数，使用 AES-128-CBC 解密。Key 文件在下载分片前预先获取并缓存。

**理由**: 解密逻辑与解析逻辑放在同一包中，保持内聚。

### 3. 分片存储
分片命名格式 `%05d.ts`，存放在 `<downloadDir>/<taskName>/` 目录。解密后的分片直接覆盖原文件。

**理由**: 与现有逻辑一致，便于断点续传判断。

### 4. 合并策略
下载完成后调用 `common/m3u8.MergeFilesFfmpeg()`（优先）或 `MergeFiles()`。

**理由**: ffmpeg 合并更可靠，处理编码问题更好。

## Risks / Trade-offs

- [风险] Key 服务器不可用 → 重试机制 + 错误提示
- [风险] 大文件内存占用 → 流式处理，不一次性加载
- [权衡] 解密后覆盖原文件 → 节省磁盘空间，但无法保留原始加密分片
