## ADDED Requirements

### Requirement: 捕获完整的HTTP请求信息
系统必须从拦截的请求中捕获URL、方法、请求头和请求体。

#### Scenario: 捕获视频请求
- **WHEN** 检测到视频URL
- **THEN** 系统捕获URL、方法、请求头（JSON字符串格式）和请求体

### Requirement: 请求头序列化
系统必须将 HTTP 请求头以 `http.Header`（map[string][]string）格式读取，JSON 序列化为字符串后存储。

#### Scenario: 存储请求头
- **WHEN** 捕获请求时
- **THEN** 请求头通过 `json.Marshal` 序列化为 JSON 字符串存入 VideoTask.Headers 字段
