## Context

当前项目是一个视频下载工具，已实现 m3u8/mp4 的解析和下载核心功能（位于 `common/ways` 和 `common/tool`）。现需要添加 HTTP API 层，让用户通过接口管理下载任务。

## Goals / Non-Goals

**Goals:**
- 提供 RESTful API 管理下载任务
- 使用 SQLite 持久化任务数据
- 支持批量执行未完成任务

**Non-Goals:**
- 不实现用户认证/授权
- 不实现任务优先级或调度策略
- 不实现下载进度的实时推送

## Decisions

### 1. Web 框架：Gin

选择 Gin 而非标准库 net/http：
- Gin 提供路由分组、中间件、参数绑定等开箱即用
- 社区成熟，文档完善

### 2. ORM：GORM + SQLite

选择 GORM 而非原生 SQL：
- 自动迁移表结构
- 简化 CRUD 操作
- SQLite 无需额外部署，数据库文件 `video.db` 存于项目根目录

### 3. 项目结构

```
internal/
  api/          # HTTP handlers
  model/        # 数据模型 (Task)
  service/      # 业务逻辑
  repository/   # 数据访问层
```

### 4. Task 表结构

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| name | string | 任务名（必填） |
| url | string | 下载地址（必填） |
| header | string | 自定义请求头（JSON） |
| type | string | 任务类型：mp4/m3u8 |
| status | int | 状态：0=待执行, 1=执行中, 2=完成, 3=失败 |
| created_at | time | 创建时间 |
| updated_at | time | 更新时间 |

## Risks / Trade-offs

- **并发下载** → 初期不限制并发数，后续可加队列控制
- **大文件下载阻塞** → 下载在 goroutine 中异步执行，不阻塞 API 响应
- **SQLite 并发写入** → 单写多读场景下 SQLite 足够，高并发场景需考虑换数据库
