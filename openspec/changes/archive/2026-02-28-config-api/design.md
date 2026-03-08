## Context

当前下载控制器 `DownloadController` 的配置是硬编码的（如 `downloadDir: "./downloads"`）。系统缺乏运行时配置能力，用户无法根据实际环境调整下载参数。

现有架构采用分层设计：api → service → repository，配置管理应遵循相同模式。

## Goals / Non-Goals

**Goals:**
- 提供 RESTful API 获取和设置系统配置
- 配置持久化到本地 JSON 文件
- 下载控制器从配置读取参数
- 配置变更立即生效（无需重启）

**Non-Goals:**
- 不支持多用户/多租户配置
- 不支持配置版本历史
- 不支持配置热加载监听文件变化

## Decisions

### 1. 配置存储方式：JSON 文件

**选择**: 使用本地 JSON 文件 (`config.json`)

**理由**:
- 项目已使用 SQLite 存储任务，配置数据结构简单，无需数据库
- JSON 文件便于手动编辑和备份
- 启动时加载一次，内存中缓存

**备选方案**: SQLite 表 - 过于复杂，配置项少且不需要查询

### 2. 配置结构

```go
type Config struct {
    MaxConcurrentTasks  int               `json:"max_concurrent_tasks"`
    MaxSegmentWorkers   int               `json:"max_segment_workers"`
    DownloadDir         string            `json:"download_dir"`
    MaxConsecutiveErrors int              `json:"max_consecutive_errors"`
    DefaultHeaders      map[string]string `json:"default_headers"`
}
```

**默认值**:
- `max_concurrent_tasks`: 3
- `max_segment_workers`: 5
- `download_dir`: "./downloads"
- `max_consecutive_errors`: 10
- `default_headers`: {}

### 3. API 设计

- `GET /api/config` - 返回完整配置
- `PUT /api/config` - 部分更新配置（只更新传入的字段）

### 4. 配置服务单例模式

配置服务使用单例，与 `DownloadController` 保持一致，确保全局配置一致性。

## Risks / Trade-offs

**[并发写入冲突]** → 使用 sync.RWMutex 保护配置读写

**[配置文件损坏]** → 启动时校验，无效则使用默认值

**[路径不存在]** → 设置 download_dir 时自动创建目录
