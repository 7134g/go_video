## ADDED Requirements

### Requirement: 配置管理 API 端点
系统 SHALL 在现有 API 路由中注册配置管理端点。

#### Scenario: 路由注册
- **WHEN** 服务器启动
- **THEN** 系统注册 GET /api/config 和 PUT /api/config 路由
