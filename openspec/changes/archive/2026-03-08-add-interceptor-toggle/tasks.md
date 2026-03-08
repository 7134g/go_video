## 实现任务

### 后端任务

- [x] 扩展Config模型，增加InterceptorEnabled和ProxyAddress字段
- [x] 实现证书检查功能CheckCertInstalled
- [x] 实现证书安装功能InstallCert
- [x] 在配置服务中增加拦截器启停逻辑
- [x] 实现拦截器单例管理和生命周期控制

### 前端任务

- [x] 在ConfigDialog.vue增加拦截器开关
- [x] 在ConfigDialog.vue增加代理地址输入框
- [x] 更新Config类型定义

### 测试任务

- [x] 测试拦截器启停功能
- [x] 测试证书检查和安装
- [x] 测试配置持久化
