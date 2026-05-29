# go_video tab tagger

为 `go_video` 配套的 Chrome / Edge / Chromium 扩展。

## 它做什么

- 给每个浏览器 tab 的所有出站 HTTP(S) 请求自动加 `X-Tab-Id: <tab.id>` 请求头。
  go_video 的 MITM 代理(`pkg/proxy/server.go`)据此把请求按 tab 分桶,从 HTML `<title>` 推导视频任务名(否则任务名就是时间戳)。
- (可选)把 Chrome 的全局代理切到 `127.0.0.1:9999`(`agent_address` 默认值),省去手动配代理。

## 安装

1. 打开 `chrome://extensions/`(Edge 用 `edge://extensions/`),右上角开启**开发者模式**。
2. 点**加载已解压的扩展程序**,选择本 `chrome_ext/` 目录。
3. 安装后,可在扩展详情页的**扩展程序选项**里:
   - 勾选**接管 Chrome 浏览器代理**(默认关闭),启用后扩展会写入 Chrome 代理设置。
   - 修改代理主机 / 端口,需与 `config.json` 的 `agent_address` 一致。

## 注意

- 必须先运行 `install_cert(.exe)` 把 go_video 的 CA 装进系统(以及 Linux 的 NSS 库),否则 Chrome 会因证书不受信任而拒绝代理后的 HTTPS。
- "接管代理"会覆盖 Chrome 的全局代理设置。取消勾选并保存后,扩展会清空该设置,Chrome 回到系统默认 / 直连。
- 扩展用 `declarativeNetRequest` session rules 实现注入头,Chrome 重启后规则会基于 onStartup 钩子自动重建。
