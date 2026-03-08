## Context

现有common/proxy包功能单一，无法完整捕获HTTP请求信息。需要重新设计一个代理拦截系统，确保捕获的视频URL可以被后续下载模块正确使用。

约束：
- 代理监听127.0.0.1:10888
- 支持Go标准库net/http
- 需要并发安全

## Goals / Non-Goals

**Goals:**
- 完整捕获HTTP请求（URL、Headers、Body、Method）
- 检测视频URL（.m3u8、.mp4）
- 提取HTML页面标题
- 生成可用于下载的VideoTask结构体
- 提供清晰的拦截器接口

**Non-Goals:**
- 不做请求/响应修改
- 不做视频下载（仅捕获信息）

## Decisions

### 1. 使用HTTP代理而非透明代理
**选择**: 实现标准HTTP代理协议
**原因**: 简单、无需系统权限、易于测试
**替代方案**: 透明代理需要iptables/系统配置，复杂度高

### 2. 拦截器模式
**选择**: 使用责任链模式，每个拦截器处理特定任务
**原因**: 职责分离、易扩展、可测试
**结构**:
- VideoDetector: 检测视频URL
- RequestCapture: 捕获请求信息
- TitleExtractor: 提取HTML标题

### 3. VideoTask结构设计
**选择**: 包含完整HTTP请求信息
```go
type VideoTask struct {
    URL     string
    Method  string
    Headers map[string]string
    Body    []byte
    Title   string
}
```
**原因**: 确保下载时可以完整复现原始请求

### 4. 并发安全
**选择**: 使用channel传递VideoTask
**原因**: Go惯用法、无需显式锁、天然并发安全

### 5. HTTPS中间人拦截
**选择**: 使用自签名CA证书进行HTTPS拦截
**原因**: 可以解密HTTPS流量，捕获加密请求中的视频URL
**实现**:
- 使用martian/mitm库生成CA证书
- 用户安装CA证书到系统信任列表
- 代理动态为每个域名签发证书
**参考**: common/proxy/cert.go的现有实现

## Risks / Trade-offs

**[风险] 内存占用** → 限制Body大小，视频请求通常无Body或很小
**[风险] 代理性能** → 仅拦截目标URL，其他请求快速转发
**[风险] CA证书安全** → 用户需妥善保管私钥文件，避免泄露
**[权衡] 需要手动安装证书** → 用户首次使用需安装CA证书到系统，但之后可拦截所有HTTPS流量
