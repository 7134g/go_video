## ADDED Requirements

### Requirement: 拦截器配置字段
系统配置SHALL包含拦截器启用状态和代理地址字段。

#### Scenario: 读取拦截器配置
- **WHEN** 系统加载配置
- **THEN** 配置包含interceptor_enabled和proxy_address字段

#### Scenario: 更新拦截器配置
- **WHEN** 用户保存配置
- **THEN** 系统持久化拦截器相关配置
