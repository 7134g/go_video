## Why

M3U8 下载逻辑与 MP4 下载逻辑存在本质差异。当前 `internal/controller/m3u8.go` 中的实现过于简化，没有充分利用 `common/m3u8` 包提供的完整解析能力（如 Master Playlist 选择、加密分片处理、字节范围请求等）。需要重构 M3U8 下载逻辑，使其正确处理各种 M3U8 格式。

## What Changes

- 重构 M3U8 下载流程，使用 `common/m3u8.ParseM3u8Data()` 替代简单的行解析
- 支持 Master Playlist：自动选择最高码率或允许用户选择
- 支持加密分片：下载 Key 文件并解密 AES-128 加密的分片
- 支持字节范围请求（EXT-X-BYTERANGE）
- 分片按顺序命名（00001.ts, 00002.ts...）存放在独立目录
- 下载完成后调用 `common/m3u8` 的合并功能

## Capabilities

### New Capabilities
- `m3u8-download`: M3U8 视频流下载能力，包括解析、分片下载、解密、合并完整流程

### Modified Capabilities
无

## Impact

- `internal/controller/m3u8.go`: 重写下载逻辑
- `internal/controller/controller.go`: 可能需要调整控制器配置
- `common/m3u8/`: 可能需要导出更多方法或添加解密功能
