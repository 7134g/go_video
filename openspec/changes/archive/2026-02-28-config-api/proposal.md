## Why

当前下载控制器的配置（并发数、下载目录等）是硬编码的，用户无法在运行时调整。需要提供配置管理接口，让用户能够根据网络环境和系统资源动态调整下载参数。

## What Changes

- 新增 GET /api/config 接口，获取当前系统配置
- 新增 PUT /api/config 接口，更新系统配置
- 新增配置持久化存储（保存到文件或数据库）
- 下载控制器读取配置并应用到下载任务

## Capabilities

### New Capabilities

- `config-management`: 系统配置的获取、设置和持久化，包括：
  - 并发任务数（同时执行的下载任务数量）
  - 单个任务同时下载的 m3u8 分片数
  - 下载存储目录
  - 连续错误最大数（达到后任务标记为失败）
  - 默认请求头

### Modified Capabilities

- `http-api`: 新增配置相关的 API 端点（GET/PUT /api/config）

## Impact

- `internal/api/`: 新增 config handler
- `internal/controller/controller.go`: 从配置读取参数而非硬编码
- `internal/model/`: 新增 Config 模型
- `internal/service/`: 新增配置服务层
- 配置文件或数据库表：持久化存储配置
