## Why

当前项目已有视频下载的核心功能（m3u8/mp4 解析和下载），但缺少一个可管理的任务系统。用户需要通过 HTTP API 来添加、管理和执行下载任务，而不是手动调用代码。

## What Changes

- 新增基于 Gin 框架的 HTTP API 服务
- 新增 SQLite 数据库存储任务信息
- 新增任务管理的 CRUD 接口
- 新增下载任务执行控制器

## Capabilities

### New Capabilities

- `http-api`: 基于 Gin 框架的 HTTP 服务，提供任务管理的 RESTful 接口
- `task-storage`: SQLite 数据库存储层，管理 task 表的 CRUD 操作
- `download-controller`: 下载任务执行控制器，批量执行未完成的下载任务

### Modified Capabilities

（无）

## Impact

- 新增依赖：gin-gonic/gin, gorm (SQLite driver)
- 新增文件：video.db (SQLite 数据库文件)
- 入口文件 main.go 需要启动 HTTP 服务
- 需要与现有的 common/ways 下载逻辑集成
