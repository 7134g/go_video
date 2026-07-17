## 1. 更新 Manifest

- [x] 1.1 在 `manifest.json` 中添加 `action` 字段，设置 `default_popup: "popup.html"`

## 2. 增强 Background

- [x] 2.1 在 `background.js` 中添加 `chrome.runtime.onMessage` 监听，处理 `getProxyStatus` 和 `toggleProxy` 消息
- [x] 2.2 在 `applyProxy()` 中添加 `chrome.action.setBadgeText` 和 `setBadgeBackgroundColor` 调用，代理开启时 badge 显示绿色 ON，关闭时显示灰色 OFF

## 3. 创建 Popup

- [x] 3.1 创建 `popup.html`：状态指示灯、代理地址显示、开关按钮、Tab 统计、打开设置入口
- [x] 3.2 创建 `popup.js`：打开时发送 `getProxyStatus` 获取状态并渲染 UI，开关按钮发送 `toggleProxy` 切换代理
