## ADDED Requirements

### Requirement: Wails 项目初始化
系统 SHALL 提供 `wails.json` 配置文件，定义桌面应用名称、前端路径、Go 模块等 Wails 项目元数据。

#### Scenario: wails.json 配置
- **WHEN** Wails 项目配置存在
- **THEN** 配置文件 SHALL 指定应用名称为 "go_video"
- **AND** 前端构建输出路径 SHALL 指向 `web/dist`
- **AND** Go 入口文件 SHALL 为 `desktop/main.go`

### Requirement: 桌面构建脚本
系统 SHALL 提供 `build_desktop.sh` 构建脚本，先编译前端再调用 Wails 打包。

#### Scenario: 成功构建桌面应用
- **WHEN** 执行 `build_desktop.sh`
- **THEN** 脚本先执行 `cd web && npm run build`
- **THEN** 脚本执行 `wails build` 生成桌面应用二进制文件
- **AND** 产物输出到 `build/desktop/` 目录

#### Scenario: 前端构建失败时终止
- **WHEN** `npm run build` 失败
- **THEN** 脚本终止并返回错误码，不执行 `wails build`

### Requirement: Go 模块依赖隔离
系统 SHALL 确保 Wails 依赖仅被 desktop 构建标签下的文件引用，默认 `go build` 不引入 Wails 依赖。

#### Scenario: web 构建不引入 Wails
- **WHEN** 执行 `go build .`（不带 `-tags desktop`）
- **THEN** 构建过程不链接 Wails 库
- **AND** 不要求 CGo 环境

#### Scenario: desktop 构建引入 Wails
- **WHEN** 执行 `go build -tags desktop`
- **THEN** 编译包含 `desktop/` 下的 Wails 入口文件
- **AND** 链接 Wails 和 WebView2 依赖

### Requirement: 资源嵌入
系统 SHALL 在 desktop 构建中将前端静态资源嵌入二进制，使桌面应用可独立运行。

#### Scenario: 前端资源绑定到桌面二进制
- **WHEN** Wails 打包完成
- **THEN** 生成的二进制文件包含 `web/dist` 所有内容
- **AND** 桌面应用离线可运行，无需额外文件

### Requirement: 平台特定打包配置
系统 SHALL 支持 Windows 和 macOS 平台的图标和打包配置。

#### Scenario: Windows 打包
- **WHEN** 在 Windows 平台构建 desktop 版本
- **THEN** 系统产出 `go_video.exe`
- **AND** 应用使用指定的 `.ico` 图标文件

#### Scenario: macOS 打包
- **WHEN** 在 macOS 平台构建 desktop 版本
- **THEN** 系统产出 `go_video.app` 应用包
- **AND** 应用使用指定的 `.icns` 图标文件
