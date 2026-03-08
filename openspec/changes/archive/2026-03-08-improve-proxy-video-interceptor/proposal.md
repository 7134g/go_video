## Why

现有的common/proxy包功能单一，缺乏完整的请求信息捕获能力。为了确保视频URL可以正常下载，需要保存完整的请求头和请求体，以便后续HTTP下载时能够正确复现原始请求。

## What Changes

- 重新设计proxy包架构，提供清晰的拦截器接口
- 完整捕获视频请求的所有信息（URL、Headers、Body、Method）
- 支持从HTML页面提取标题信息
- 生成包含完整下载信息的VideoTask结构体
- 提供可配置的URL过滤规则（.m3u8、.mp4等）
- 添加并发安全的任务收集机制

## Capabilities

### New Capabilities
- `http-proxy-interceptor`: HTTP代理服务器，拦截127.0.0.1:10888的所有请求
- `https-mitm-interceptor`: HTTPS中间人拦截，使用自签名CA证书解密HTTPS流量
- `video-url-detector`: 检测和过滤视频URL（.m3u8、.mp4等格式）
- `request-capture`: 完整捕获HTTP请求信息（headers、body、method）
- `html-title-extractor`: 从HTML响应中提取页面标题
- `video-task-builder`: 构建包含完整下载信息的VideoTask结构体

### Modified Capabilities
<!-- 无现有能力需要修改 -->

## Impact

- 新增 `pkg/proxy/` 包（替代 `common/proxy`）
- 新增 `pkg/models/` 包用于定义VideoTask等数据结构
- 可能影响现有调用common/proxy的代码
- 需要确保代理服务的性能和稳定性
