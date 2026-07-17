## Context

当前 Chrome 扩展（MV3）已有：
- `declarativeNetRequest` 注入 `X-Tab-Id` 请求头
- `chrome.proxy.settings` 设置固定代理（`fixed_servers`）
- Options 页面配置代理地址和开关

缺少的是 Popup 弹窗交互——用户必须打开 options 页面才能开关代理，且无法在工具栏快速查看代理状态。

## Goals / Non-Goals

**Goals:**
- 点击扩展图标弹出面板，显示代理连接状态和地址
- 一键开关代理，无需打开 options 页面
- 图标 badge 上显示代理开关状态（ON/OFF）
- 保持与现有 `chrome.storage.local` + `background.js` 的兼容

**Non-Goals:**
- 不支持多 Profile 管理
- 不支持 PAC 或规则切换
- 不在 Popup 中编辑代理地址（地址配置仍在 options 页）

## Decisions

**1. Popup 与 Background 通信方式：`chrome.runtime.sendMessage`**

Popup 是短暂存在的（关闭即销毁），不能直接依赖它维护状态。因此 popup 通过 `sendMessage` 向 background service worker 发起请求：
- `getStatus` → 返回当前代理状态 + tab 规则数
- `toggleProxy` → 切换代理开关，返回新状态

选择消息机制而非 `storage.local` 监听，因为消息是 request/response 模式，popup 可以拿到即时返回值更新 UI，而 `storage.onChanged` 在 popup 场景下时序不可靠。

**2. 图标 Badge：`chrome.action.setBadgeText`**

Background 在 `applyProxy()` 时同步设置 badge 文字（`ON`/`OFF`）和颜色（绿/灰），这样用户即使不打开 popup 也能一眼看到代理状态。

**3. Popup 关闭时不清理任何状态**

Popup 是纯展示+控制界面，不持有状态。代理开关的持久化由 `applyProxy()` 统一管理（写入 `chrome.storage.local` + `chrome.proxy.settings`），popup 只做读取和触发。

## Risks / Trade-offs

- **Popup 短暂生命周期**：MV3 的 popup 在失焦后立即销毁，用户无法在 popup 中查看长时间的实时状态。解决：badge 承担持续可见的状态指示。
- **Badge 文字在 Chrome 关闭后丢失**：`onStartup` 和 `onInstalled` 监听器已存在，启动时同步触发 `applyProxy()`，badge 会自动恢复。
