# go_video tab tagger

为 `go_video` 配套的 Chrome / Edge / Chromium 扩展。

## 它做什么

- 为每个浏览器 tab 的所有出站 HTTP(S) 请求自动注入 `X-Tab-Id: <tab.id>` 请求头。go_video 的 MITM 代理（`pkg/proxy/server.go`）据此把请求按 tab 分桶，从 HTML `<title>` 推导视频任务名（否则任务名就是时间戳）。
- 扩展使用 MV3 `declarativeNetRequest` session rules 实现注入：每个存活的 tab 对应一条规则。Chrome 重启后 session rules 自动清空，扩展通过 `chrome.runtime.onStartup` 钩子重建。
- 可通过弹窗（popup）和选项页两种方式管理代理连接状态：
  - **弹窗**（点击工具栏图标）— 查看/切换代理开关状态、查看当前注入 tab 数量、跳转到设置页。
  - **选项页**（扩展详情 → 扩展程序选项）— 勾选**接管 Chrome 浏览器代理**（默认关闭），启用后扩展会写入 Chrome 代理设置；修改代理主机/端口使其与 `config.json` 的 `agent_address` 一致。

## 安装

1. 打开 `chrome://extensions/`（Edge 用 `edge://extensions/`），右上角开启**开发者模式**。
2. 点**加载已解压的扩展程序**，选择本 `chrome_ext/` 目录。
3. 安装后即可在工具栏看到扩展图标，点击可查看状态和快速切换代理。

## 注意

- 必须先运行 `install_cert(.exe)` 把 go_video 的 CA 装进系统（以及 Linux 的 NSS 库），否则 Chrome 会因证书不受信任而拒绝代理后的 HTTPS。
- "接管代理"会覆盖 Chrome 的全局代理设置。取消勾选并保存后，扩展会清空该设置，Chrome 回到系统默认/直连。
