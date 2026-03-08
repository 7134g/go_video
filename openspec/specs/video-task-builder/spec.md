## ADDED Requirements

### Requirement: 从捕获的数据构建VideoTask
系统必须构建包含URL、方法、请求头、请求体和标题的VideoTask结构体。

#### Scenario: 创建VideoTask
- **WHEN** 捕获到视频请求
- **THEN** 系统创建包含所有捕获信息的VideoTask

### Requirement: 线程安全的任务收集
系统必须提供并发安全的机制来收集VideoTask实例。

#### Scenario: 并发收集多个任务
- **WHEN** 同时捕获多个视频请求
- **THEN** 所有任务都被收集且无数据竞争
