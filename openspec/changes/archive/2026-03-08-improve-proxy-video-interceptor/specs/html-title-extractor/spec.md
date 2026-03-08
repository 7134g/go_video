## ADDED Requirements

### Requirement: 从HTML响应中提取标题
系统必须在URL以.html结尾时从HTML响应中提取页面标题。

#### Scenario: 提取HTML标题
- **WHEN** 响应为HTML内容
- **THEN** 系统从<title>标签中提取文本

#### Scenario: 跳过非HTML响应
- **WHEN** 响应不是HTML
- **THEN** 系统不尝试提取标题
