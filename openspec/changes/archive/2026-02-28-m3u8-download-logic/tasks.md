## 1. 重构 common/m3u8 包

- [x] 1.1 重构 M3u8 结构体，添加解密相关字段
- [x] 1.2 实现 AES-128-CBC 解密函数 `Decrypt(data, key, iv []byte) ([]byte, error)`
- [x] 1.3 添加 Key 下载和缓存功能
- [x] 1.4 优化 URL 解析，支持相对路径转绝对路径

## 2. 重构 controller 层 M3U8 下载逻辑

- [x] 2.1 修改 `parseM3u8()` 使用 `common/m3u8.ParseM3u8Data()`
- [x] 2.2 实现 Master Playlist 检测和最高码率选择
- [x] 2.3 实现加密密钥预下载和缓存
- [x] 2.4 修改 `downloadSegment()` 支持解密

## 3. 分片管理和合并

- [x] 3.1 确保分片按 `%05d.ts` 格式命名
- [x] 3.2 实现下载完成后自动调用合并功能
- [x] 3.3 添加 ffmpeg 可用性检测，选择合适的合并方式

## 4. 测试验证

- [x] 4.1 测试普通 M3U8 下载
- [x] 4.2 测试 Master Playlist 处理
- [x] 4.3 测试加密分片下载和解密
