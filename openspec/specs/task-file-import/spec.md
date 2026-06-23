## Purpose

TBD — 启动时从 `task.txt` 批量导入任务到数据库。

## ADDED Requirements

### Requirement: 启动时导入 task.txt
系统 SHALL 在程序启动完成（DB 初始化、CA 检查之后）读取工作目录下的 `task.txt` 文件，按行解析并批量创建任务。

#### Scenario: 正常导入多个任务
- **WHEN** `task.txt` 包含偶数行文本，奇数行为任务名称、偶数行为 URL
- **THEN** 系统按每两行一组解析，为每组创建 `model.Task` 并写入数据库，处理完毕后删除 `task.txt`

#### Scenario: 文件不存在
- **WHEN** 工作目录下不存在 `task.txt`
- **THEN** 系统静默跳过，不输出任何日志，继续正常启动

#### Scenario: 存在空行
- **WHEN** `task.txt` 包含空行
- **THEN** 系统跳过空行，仅处理非空行组成的名称/URL 对

### Requirement: URL 后缀类型判定
系统 SHALL 根据 URL 后缀判断任务类型：以 `.mp4` 结尾判定为 MP4 类型，其余一律判定为 M3U8 类型。

#### Scenario: mp4 URL
- **WHEN** URL 以 `.mp4` 结尾
- **THEN** 任务 Type 设置为 MP4

#### Scenario: m3u8 URL
- **WHEN** URL 以 `.m3u8` 结尾
- **THEN** 任务 Type 设置为 M3U8

#### Scenario: 其他后缀
- **WHEN** URL 不含可识别的视频后缀
- **THEN** 任务 Type 默认为 M3U8

### Requirement: 使用默认 Header
系统 SHALL 为导入的任务使用配置中的默认 Header，Status 设置为 Pending。

#### Scenario: 导入任务带默认 Header
- **WHEN** 从 `task.txt` 解析出任务
- **THEN** 任务的 Header 字段填充为配置中的 `default_headers`，Status 为 Pending

### Requirement: 导入完成后删除文件
系统 SHALL 在成功解析并导入所有任务后删除 `task.txt` 文件。

#### Scenario: 成功导入后删除
- **WHEN** `task.txt` 中所有任务成功写入数据库
- **THEN** 系统删除 `task.txt` 文件

#### Scenario: 删除失败不阻塞启动
- **WHEN** 任务导入成功但删除文件失败（如权限不足）
- **THEN** 系统输出 warning 日志，不中断程序启动

### Requirement: 单行异常不阻断其余任务
系统 SHALL 在某一行解析或插入失败时记录错误日志并继续处理后续行。

#### Scenario: URL 已存在
- **WHEN** 某任务的 URL 在数据库中已存在
- **THEN** 系统跳过该任务，记录日志，继续处理后续任务
