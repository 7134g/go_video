## 架构决策

### 拦截器生命周期管理
拦截器由配置服务管理，当配置更新时触发启停逻辑。使用单例模式确保只有一个拦截器实例运行。

### CA证书管理
使用pkg/proxy包现有的CA证书加载功能，新增证书检查和安装模块。证书文件路径通过配置指定。

## 技术方案

### 前端实现
在ConfigDialog.vue中增加：
- el-switch组件控制拦截器启用状态
- el-input组件配置代理地址
- 配置项位于现有表单项之后

### 后端实现

**Config模型扩展** (internal/model/config.go):
```go
type Config struct {
    // 现有字段...
    InterceptorEnabled bool   `json:"interceptor_enabled"`
    ProxyAddress       string `json:"proxy_address"`
}
```

**配置服务** (internal/service/config.go):
- UpdateConfig方法中增加拦截器启停逻辑
- 启用时调用proxy.StartServer
- 禁用时调用proxy.StopServer
- 启动前调用证书检查和安装

**证书管理** (pkg/proxy/cert.go):
- CheckCertInstalled(): 检查证书是否已安装
- InstallCert(): 安装CA证书到系统

## 实现细节

### 拦截器启动流程
1. 检查InterceptorEnabled是否为true
2. 调用CheckCertInstalled检查证书
3. 若未安装则调用InstallCert
4. 使用ProxyAddress启动proxy.Server

### 证书检查逻辑
- Windows: 检查证书存储
- macOS/Linux: 检查系统证书目录

### 错误处理
- 证书安装失败返回错误，不启动拦截器
- 拦截器启动失败回滚配置状态
