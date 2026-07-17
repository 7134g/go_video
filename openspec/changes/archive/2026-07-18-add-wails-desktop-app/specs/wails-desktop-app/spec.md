## ADDED Requirements

### Requirement: Wails 应用启动
系统 SHALL 在 desktop 构建标签下提供 Wails 应用入口，启动时创建主窗口并加载 Vue 前端。

#### Scenario: 桌面应用启动成功
- **WHEN** 用户运行 desktop 构建的二进制文件
- **THEN** 系统创建一个窗口（默认 1200x800），加载 Vue 前端页面
- **AND** 系统使用系统原生 WebView 渲染前端内容

#### Scenario: Wails 应用仅绑定桌面方法
- **WHEN** Wails 应用启动
- **THEN** 系统仅暴露 Wails 绑定方法（Go 函数调用），不启动 Gin HTTP 服务器
- **AND** 系统不监听任何 TCP 端口

### Requirement: 应用生命周期管理
系统 SHALL 在窗口关闭时触发应用退出流程，停止所有下载任务并保存数据库状态。

#### Scenario: 窗口关闭时优雅退出
- **WHEN** 用户关闭桌面窗口
- **THEN** 系统停止所有运行中的下载任务
- **AND** 系统将内存中的任务状态持久化到数据库
- **AND** 系统正常退出（exit code 0）

### Requirement: 窗口标题与尺寸
系统 SHALL 设置窗口标题为应用名称，提供合理的默认窗口尺寸。

#### Scenario: 窗口配置
- **WHEN** Wails 主窗口创建
- **THEN** 窗口标题 SHALL 为"视频下载器"
- **AND** 窗口默认尺寸 SHALL 为 1200x800 像素
- **AND** 窗口最小尺寸 SHALL 为 800x600 像素
