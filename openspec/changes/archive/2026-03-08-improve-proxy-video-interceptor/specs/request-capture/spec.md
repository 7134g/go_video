## ADDED Requirements

### Requirement: 捕获完整的HTTP请求信息
系统必须从拦截的请求中捕获URL、方法、请求头和请求体。

#### Scenario: 捕获视频请求
- **WHEN** 检测到视频URL
- **THEN** 系统捕获URL、方法、所有请求头和请求体

### Requirement: 保存请求头用于下载
系统必须将所有HTTP请求头存储为键值对。

#### Scenario: 存储请求头
- **WHEN** 捕获请求时
- **THEN** 所有请求头以map[string]string格式存储
