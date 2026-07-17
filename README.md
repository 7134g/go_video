# 介绍
本项目为了下载在线播放的m3u8视频和mp4视频而设计的

# 启动
双击启动 desktop.exe


## 两种模式

### 1. 自己抓包连接添加进任务列表

### 2. 拦截器自动解析任务
#### 启动准备
1. 双击 install_cert.exe 生成并安装拦截器功能所需的本地证书，用于解析https
2. 参考 chrome_ext/README.md 中内容将浏览器扩展组件安装正确

完成上述两点就可以正常使用该拦截器功能

***每次使用时需要将刚刚加入浏览器扩展里代理打开，不使用时候关闭***
***若当前已经启动了 desktop.exe 请关闭重新打开即可***

拦截器功能正常使用条件
1. 证书已经安装
2. 浏览器扩展里的 "go_video tab tagger" 代理已经开启
3. desktop.exe 运行中


## 注意
1. 若在使用拦截器功能时候，关闭了desktop.exe，并且没有将 "go_video tab tagger" 浏览器扩展的代理关闭，会导致浏览器无法上网。
    - 解决方式，只需要在 "go_video tab tagger"  点击关闭代理
2. 卸载ca证书执行 uninstall_cert.exe
3. 卸载chrome扩展，需要自行去chrome扩展管理里操作


=======================================================================================================================

# 调试与开发

## 启动前端
- `cd web && npm run dev`

## 启动后端服务
- `go build . && ./go_video`

## 编译证书安装工具
- `go build -o install_cert.exe ./tool/install-cert`

## 编译证书卸载工具
- `go build -o install_cert.exe ./tool/uninstall-cert`


## 一次性编译打包上述所有
- `release.sh`




