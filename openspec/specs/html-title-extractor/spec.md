## ADDED Requirements

### Requirement: 从WebTree缓存中提取标题
系统 SHALL 将每个Tab的HTTP响应体缓存到 WebTree（按 tabId + URL 索引），拦截到视频URL时扫描同Tab下的缓存内容查找 `<title>` 标签提取标题。

#### Scenario: 从HTML响应提取标题
- **WHEN** 同Tab下存在HTML页面（包含 `<title>` 标签）
- **THEN** 系统提取标题文本作为视频任务名称

#### Scenario: 无标题时使用时间戳
- **WHEN** 同Tab下没有HTML页面或不含 `<title>` 标签
- **THEN** 系统使用当前时间戳（毫秒）作为任务标题
