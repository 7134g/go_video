## ADDED Requirements

### Requirement: 拦截器配置字段
系统配置 SHALL 包含拦截器启用状态（`interceptor_enabled`）、代理监听地址（`agent_address`）和上游代理地址（`vpn_address`）字段。

#### Scenario: 读取拦截器配置
- **WHEN** 系统加载配置
- **THEN** 配置包含 interceptor_enabled、agent_address 和 vpn_address 字段

#### Scenario: 更新拦截器配置
- **WHEN** 用户保存配置
- **THEN** 系统持久化拦截器相关配置，并根据 interceptor_enabled 启停代理服务

#### Scenario: 默认值
- **WHEN** 用户未配置拦截器相关字段
- **THEN** interceptor_enabled 默认为 false，agent_address 默认为 127.0.0.1:9999，vpn_address 默认为 127.0.0.1:7890
