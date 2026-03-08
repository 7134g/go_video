## ADDED Requirements

### Requirement: 拦截器启停控制
系统SHALL提供拦截器的启动和停止功能，允许用户通过配置控制拦截器状态。

#### Scenario: 启用拦截器
- **WHEN** 用户在配置中启用拦截器
- **THEN** 系统调用pkg/proxy/server.go启动代理服务

#### Scenario: 禁用拦截器
- **WHEN** 用户在配置中禁用拦截器
- **THEN** 系统停止正在运行的代理服务

#### Scenario: 拦截器状态持久化
- **WHEN** 用户修改拦截器启用状态
- **THEN** 系统保存配置到存储中

### Requirement: 前端拦截器开关
前端配置界面SHALL提供拦截器开关控件。

#### Scenario: 显示拦截器开关
- **WHEN** 用户打开配置对话框
- **THEN** 界面显示拦截器启用/禁用开关

#### Scenario: 切换拦截器状态
- **WHEN** 用户点击拦截器开关
- **THEN** 系统更新配置并应用新状态
