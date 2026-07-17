## ADDED Requirements

### Requirement: 显示代理状态

Popup 打开时 SHALL 显示当前代理连接状态（已连接/未连接）和代理地址（host:port）。

#### Scenario: 代理已开启时显示绿色状态
- **WHEN** popup 打开且 `storage.local.proxyEnabled === true`
- **THEN** 状态指示灯显示绿色，文案显示"代理已连接"，下方显示当前代理地址

#### Scenario: 代理已关闭时显示灰色状态
- **WHEN** popup 打开且 `storage.local.proxyEnabled === false`
- **THEN** 状态指示灯显示灰色，文案显示"代理未连接"，下方显示代理地址（灰色）

### Requirement: 一键开关代理

Popup 中 SHALL 提供开关按钮，点击后调用 background 切换代理状态。

#### Scenario: 点击关闭代理
- **WHEN** 用户点击"关闭代理"按钮（当前代理已开启）
- **THEN** background 执行 `applyProxy(false)`，清除 `chrome.proxy.settings`
- **THEN** badge 文字更新为 OFF，颜色变灰
- **THEN** popup 内状态更新为"代理未连接"

#### Scenario: 点击开启代理
- **WHEN** 用户点击"开启代理"按钮（当前代理已关闭）
- **THEN** background 执行 `applyProxy(true)`，设置 `chrome.proxy.settings` 为 `fixed_servers`
- **THEN** badge 文字更新为 ON，颜色变绿
- **THEN** popup 内状态更新为"代理已连接"

### Requirement: 显示 Tab 注入统计

Popup 中 SHALL 显示当前正在被注入 X-Tab-Id 的 tab 数量。

#### Scenario: 正常显示 tab 数量
- **WHEN** popup 打开
- **THEN** 显示"正在为 N 个 tab 注入 X-Tab-Id"，N 来自 `declarativeNetRequest.getSessionRules()` 的规则数

#### Scenario: 无 tab 时显示 0
- **WHEN** popup 打开且没有活跃 tab 规则
- **THEN** 显示"正在为 0 个 tab 注入 X-Tab-Id"

### Requirement: 跳转到设置页

Popup 中 SHALL 提供"打开设置"入口，点击后打开 options 页面。

#### Scenario: 点击打开设置
- **WHEN** 用户点击"打开设置"
- **THEN** 调用 `chrome.runtime.openOptionsPage()` 打开 options 页
- **THEN** popup 自动关闭

### Requirement: 图标 Badge 状态指示

扩展图标上 SHALL 显示 badge 文字指示代理开关状态。

#### Scenario: 代理开启时 badge 显示 ON
- **WHEN** `applyProxy(true)` 执行完毕
- **THEN** `chrome.action.setBadgeText({ text: "ON" })`，背景色为绿色

#### Scenario: 代理关闭时 badge 显示 OFF
- **WHEN** `applyProxy(false)` 执行完毕
- **THEN** `chrome.action.setBadgeText({ text: "OFF" })`，背景色为灰色

#### Scenario: 浏览器启动时 badge 恢复
- **WHEN** `chrome.runtime.onStartup` 触发
- **THEN** `applyProxy()` 被调用，badge 恢复为正确的 ON/OFF 状态

### Requirement: Popup 与 Background 通信

Popup SHALL 通过 `chrome.runtime.sendMessage` 与 background service worker 通信获取/切换状态。

#### Scenario: popup 请求状态
- **WHEN** popup 发送 `{ action: "getProxyStatus" }` 消息
- **THEN** background 返回 `{ enabled, host, port, ruleCount }`

#### Scenario: popup 请求切换代理
- **WHEN** popup 发送 `{ action: "toggleProxy" }` 消息
- **THEN** background 执行 `applyProxy(!enabled)` 并返回新状态 `{ enabled, host, port, ruleCount }`
