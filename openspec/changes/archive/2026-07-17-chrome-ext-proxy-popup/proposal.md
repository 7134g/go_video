## Why

当前 Chrome 扩展已支持通过 `chrome.proxy.settings` 设置代理，但用户必须打开 options 页面才能开关代理，操作路径长。增加 Popup 弹窗后，点击图标即可一键开关代理，提升使用体验。

## What Changes

- **新增 `popup.html` + `popup.js`**：点击扩展图标弹出面板，显示代理连接状态、当前代理地址、Tab 注入数统计，以及一键开关按钮
- **更新 `manifest.json`**：添加 `action.default_popup` 字段指向 popup.html
- **更新 `background.js`**：添加消息响应（popup 查询/切换代理）、图标 badge 显示 ON/OFF
- **微调 `options.html`/`options.js`**：无功能性变更，仅协调文案

## Capabilities

### New Capabilities
- `proxy-popup`: Chrome 扩展 popup 弹窗，提供代理状态查看和快速开关能力

### Modified Capabilities

无

## Impact

仅影响 `chrome_ext/` 目录下的文件：
- `manifest.json` — 新增 `action.default_popup`
- `background.js` — 新增消息监听和 badge 更新逻辑
- `options.html` — 无功能变化
- `options.js` — 无变化
- `popup.html` — 新增
- `popup.js` — 新增

后端（Go 代码）无影响。
