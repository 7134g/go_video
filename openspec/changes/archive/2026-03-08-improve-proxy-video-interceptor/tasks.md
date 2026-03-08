## 1. 数据结构定义

- [x] 1.1 创建pkg/models/video_task.go定义VideoTask结构体
- [x] 1.2 添加必要的字段（URL、Method、Headers、Body、Title）

## 2. 代理服务器核心

- [x] 2.1 创建pkg/proxy/server.go实现HTTP代理服务器
- [x] 2.2 实现监听127.0.0.1:10888
- [x] 2.3 实现请求转发逻辑
- [x] 2.4 创建pkg/proxy/cert.go实现CA证书生成和加载
- [x] 2.5 集成HTTPS中间人拦截支持

## 3. 拦截器实现

- [x] 3.1 创建pkg/proxy/interceptor.go定义拦截器接口
- [x] 3.2 实现VideoDetector检测.m3u8和.mp4 URL
- [x] 3.3 实现RequestCapture捕获完整请求信息
- [x] 3.4 实现TitleExtractor从HTML提取标题

## 4. 任务收集

- [x] 4.1 创建pkg/proxy/collector.go实现并发安全的任务收集器
- [x] 4.2 使用channel传递VideoTask

## 5. 集成测试

- [x] 5.1 测试代理服务器启动和转发
- [x] 5.2 测试视频URL检测和捕获
- [x] 5.3 测试HTML标题提取
