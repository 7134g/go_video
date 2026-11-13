async function refreshStatus() {
  const resp = await chrome.runtime.sendMessage({ action: "getProxyStatus" });
  render(resp);
}

function render(state) {
  const indicator = document.getElementById("indicator");
  const statusText = document.getElementById("statusText");
  const address = document.getElementById("address");
  const toggleBtn = document.getElementById("toggleBtn");
  const ruleStat = document.getElementById("ruleStat");

  if (state.enabled) {
    indicator.className = "indicator on";
    statusText.textContent = "代理已连接";
    address.textContent = `${state.host}:${state.port}`;
    toggleBtn.textContent = "关闭代理";
  } else {
    indicator.className = "indicator off";
    statusText.textContent = "代理未连接";
    address.textContent = `${state.host}:${state.port}`;
    toggleBtn.textContent = "开启代理";
  }

  ruleStat.textContent = `正在为 ${state.ruleCount} 个 tab 注入 X-Tab-Id`;
}

document.getElementById("toggleBtn").addEventListener("click", async () => {
  try {
    const resp = await chrome.runtime.sendMessage({ action: "toggleProxy" });
    render(resp);
  } catch (e) {
    document.getElementById("statusText").textContent = "操作失败";
  }
});

document.getElementById("settingsLink").addEventListener("click", () => {
  chrome.runtime.openOptionsPage();
});

document.addEventListener("DOMContentLoaded", refreshStatus);
