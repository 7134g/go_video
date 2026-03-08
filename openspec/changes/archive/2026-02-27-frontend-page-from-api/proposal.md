## Why

后端已完成视频下载任务管理 API，需要配套的前端管理界面，让用户能够可视化地创建、管理下载任务并实时查看下载进度。

## What Changes

- 新增 Vue3 + Element Plus 前端项目
- 实现任务列表页面（展示所有任务及状态）
- 实现任务创建/编辑表单
- 实现 WebSocket 实时进度展示
- 支持任务的增删改查操作

## Capabilities

### New Capabilities

- `task-management-ui`: 任务管理界面，包含任务列表、创建、编辑、删除功能
- `realtime-progress`: WebSocket 实时进度展示组件

### Modified Capabilities

（无，这是新增前端项目）

## Impact

- 新增 `web/` 目录存放前端代码
- 需要配置后端 CORS 或代理
- 依赖：Vue 3、Element Plus、axios、WebSocket API
