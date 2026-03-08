## ADDED Requirements

### Requirement: 生成自签名CA证书
系统必须能够生成用于HTTPS拦截的自签名CA证书和私钥。

#### Scenario: 生成证书文件
- **WHEN** 首次启动或证书不存在时
- **THEN** 系统生成CA证书和私钥文件

### Requirement: 加载CA证书进行HTTPS拦截
系统必须加载CA证书和私钥，用于动态签发域名证书。

#### Scenario: 拦截HTTPS请求
- **WHEN** 客户端发起HTTPS请求
- **THEN** 代理使用CA证书为目标域名签发临时证书并解密流量

### Requirement: 支持HTTPS流量透明转发
系统必须在解密后将HTTPS请求转发到目标服务器，并加密响应返回客户端。

#### Scenario: HTTPS请求被转发
- **WHEN** 解密HTTPS请求后
- **THEN** 代理转发请求到目标服务器并返回加密响应
