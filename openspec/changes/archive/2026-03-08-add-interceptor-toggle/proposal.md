## Why

当前系统缺少对视频拦截器的控制能力，用户无法灵活地启用或禁用拦截功能，也无法配置代理地址。需要提供一个前端开关来控制拦截器的启停，并在启动时自动检查和安装CA证书，确保HTTPS拦截正常工作。

## What Changes

- 前端配置界面增加拦截器开关（启用/禁用）
- 前端配置界面增加代理地址配置项
- 后端Config模型增加拦截器相关字段
- 后端配置API支持拦截器配置的读取和更新
- 拦截器启用时自动调用pkg/proxy/server.go启动代理服务
- 启动前自动检查CA证书是否已安装，未安装则自动安装

## Capabilities

### New Capabilities
- `interceptor-control`: 拦截器启停控制，包括前端开关UI、后端启停逻辑、与pkg/proxy/server.go的集成
- `proxy-configuration`: 代理地址配置管理，包括前端配置表单、后端配置存储
- `certificate-management`: CA证书自动检查和安装，确保HTTPS拦截正常工作

### Modified Capabilities
- `system-configuration`: 现有系统配置功能需要扩展，增加拦截器相关配置项

## Impact

**前端影响**:
- `web/src/components/ConfigDialog.vue`: 增加拦截器开关和代理地址配置项
- `web/src/api/config.ts`: Config类型定义需要增加新字段

**后端影响**:
- `internal/model/config.go`: Config结构体增加拦截器相关字段
- `internal/service/config.go`: 配置服务需要处理拦截器启停逻辑
- `internal/api/config.go`: 配置API需要支持新字段的读写
- `pkg/proxy/server.go`: 需要被配置服务调用以启动/停止拦截器

**依赖**:
- 需要CA证书管理相关功能（检查、安装）
- 需要确保pkg/proxy包的CA证书文件路径可配置
