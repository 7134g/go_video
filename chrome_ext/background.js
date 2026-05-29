// background service worker:
// 1. 为每个 tab 维护一条 declarativeNetRequest session rule,
//    给该 tab 所有出站请求加 X-Tab-Id: <tabId> 请求头。
//    MV3 不允许在 action 里动态读 tabId,只能"每 tab 一条规则"。
// 2. 可选地把 Chrome 代理切到本地 127.0.0.1:9999(从 options 页配置)。

const ALL_RESOURCE_TYPES = [
  "main_frame", "sub_frame", "stylesheet", "script", "image",
  "font", "object", "xmlhttprequest", "ping", "csp_report",
  "media", "websocket", "webtransport", "webbundle", "other"
];

const DEFAULTS = { proxyEnabled: false, proxyHost: "127.0.0.1", proxyPort: 9999 };

function tabRule(tabId) {
  return {
    id: tabId,
    priority: 1,
    condition: {
      tabIds: [tabId],
      resourceTypes: ALL_RESOURCE_TYPES
    },
    action: {
      type: "modifyHeaders",
      requestHeaders: [
        { header: "X-Tab-Id", operation: "set", value: String(tabId) }
      ]
    }
  };
}

async function syncAllTabs() {
  const tabs = await chrome.tabs.query({});
  const existing = await chrome.declarativeNetRequest.getSessionRules();
  await chrome.declarativeNetRequest.updateSessionRules({
    removeRuleIds: existing.map(r => r.id),
    addRules: tabs.filter(t => t.id !== undefined && t.id >= 0).map(t => tabRule(t.id))
  });
}

async function addTabRule(tabId) {
  if (tabId === undefined || tabId < 0) return;
  await chrome.declarativeNetRequest.updateSessionRules({
    removeRuleIds: [tabId],
    addRules: [tabRule(tabId)]
  });
}

async function removeTabRule(tabId) {
  await chrome.declarativeNetRequest.updateSessionRules({
    removeRuleIds: [tabId]
  });
}

async function applyProxy() {
  const cfg = { ...DEFAULTS, ...(await chrome.storage.local.get(DEFAULTS)) };
  if (cfg.proxyEnabled) {
    await chrome.proxy.settings.set({
      value: {
        mode: "fixed_servers",
        rules: {
          singleProxy: { scheme: "http", host: cfg.proxyHost, port: Number(cfg.proxyPort) },
          bypassList: ["localhost", "127.0.0.1"]
        }
      },
      scope: "regular"
    });
  } else {
    await chrome.proxy.settings.clear({ scope: "regular" });
  }
}

chrome.runtime.onInstalled.addListener(async () => {
  await syncAllTabs();
  await applyProxy();
});

chrome.runtime.onStartup.addListener(async () => {
  await syncAllTabs();
  await applyProxy();
});

chrome.tabs.onCreated.addListener(tab => {
  addTabRule(tab.id);
});

chrome.tabs.onRemoved.addListener(tabId => {
  removeTabRule(tabId);
});

chrome.storage.onChanged.addListener((changes, area) => {
  if (area !== "local") return;
  if ("proxyEnabled" in changes || "proxyHost" in changes || "proxyPort" in changes) {
    applyProxy();
  }
});
