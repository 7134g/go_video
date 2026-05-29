const DEFAULTS = { proxyEnabled: false, proxyHost: "127.0.0.1", proxyPort: 9999 };

async function load() {
  const cfg = { ...DEFAULTS, ...(await chrome.storage.local.get(DEFAULTS)) };
  document.getElementById("proxyEnabled").checked = cfg.proxyEnabled;
  document.getElementById("proxyHost").value = cfg.proxyHost;
  document.getElementById("proxyPort").value = cfg.proxyPort;

  const rules = await chrome.declarativeNetRequest.getSessionRules();
  document.getElementById("ruleCount").textContent = rules.length;
}

async function save() {
  await chrome.storage.local.set({
    proxyEnabled: document.getElementById("proxyEnabled").checked,
    proxyHost: document.getElementById("proxyHost").value.trim() || "127.0.0.1",
    proxyPort: Number(document.getElementById("proxyPort").value) || 9999
  });
  const msg = document.getElementById("msg");
  msg.textContent = "已保存";
  setTimeout(() => { msg.textContent = ""; }, 2000);
}

document.getElementById("save").addEventListener("click", save);
document.addEventListener("DOMContentLoaded", load);
