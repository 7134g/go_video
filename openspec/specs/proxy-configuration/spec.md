## ADDED Requirements

### Requirement: 代理地址配置
系统SHALL允许用户配置代理服务器监听地址。

#### Scenario: 配置代理地址
- **WHEN** 用户在配置中设置代理地址
- **THEN** 系统保存代理地址配置

#### Scenario: 使用配置的代理地址
- **WHEN** 拦截器启动
- **THEN** 系统使用配置的地址启动代理服务

#### Scenario: 默认代理地址
- **WHEN** 用户未配置代理地址
- **THEN** 系统使用默认地址127.0.0.1:9999
