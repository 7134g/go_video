## Context

后端已有完整的任务管理 API（Gin 框架，端口 8080），需要构建配套的前端界面。当前项目无前端代码，需从零搭建 Vue3 + Element Plus 项目。

后端 API：
- `GET /api/tasks` - 获取任务列表
- `POST /api/tasks` - 创建任务
- `PUT /api/tasks/:id` - 更新任务
- `DELETE /api/tasks/:id` - 删除任务
- `POST /api/tasks/start` - 启动待执行任务
- `WS /api/tasks/progress` - 实时进度推送

## Goals / Non-Goals

**Goals:**
- 提供任务的增删改查界面
- 实时展示下载进度
- 简洁易用的单页应用

**Non-Goals:**
- 用户认证/权限管理
- 多语言支持
- 移动端适配

## Decisions

### 1. 项目结构
采用 Vite + Vue3 + TypeScript，放置于 `web/` 目录。

**理由**：Vite 构建速度快，TypeScript 提供类型安全。

### 2. UI 框架
使用 Element Plus。

**理由**：用户指定，且与 Vue3 兼容性好，组件丰富。

### 3. 状态管理
不使用 Vuex/Pinia，直接用组合式 API 管理状态。

**理由**：应用简单，单页面无需复杂状态管理。

### 4. API 请求
使用 axios，封装统一的请求模块。

### 5. WebSocket 进度
使用原生 WebSocket API，在任务列表页面建立连接，实时更新进度。

## Risks / Trade-offs

- **CORS 问题** → 开发时使用 Vite 代理，生产环境后端配置 CORS
- **WebSocket 断连** → 实现自动重连机制
