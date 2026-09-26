"use strict";

const MENU_ID = "remotv-open";
browser.menus.create({
  id: MENU_ID,
  title: "Configure Remotv…",
  contexts: ["tab", "page"],
});

async function updateMenu() {
  const { deviceName } = await browser.storage.local.get("deviceName");
  await browser.menus.update(MENU_ID, {
    title: deviceName ? `Open on ${deviceName.replace(/&/g, "&&")}` : "Configure Remotv…",
  });
}

async function showStatus(message, failed = false) {
  await browser.browserAction.setBadgeText({ text: failed ? "!" : "OK" });
  await browser.browserAction.setBadgeBackgroundColor({ color: failed ? "#b42318" : "#167647" });
  await browser.browserAction.setTitle({ title: `${message} — click for settings` });
}

async function openOnDevice(info, tab) {
  if (info.menuItemId !== MENU_ID) return;
  try {
    const saved = await browser.storage.local.get(["serverUri", "deviceName"]);
    if (!saved.serverUri || !saved.deviceName) {
      await browser.runtime.openOptionsPage();
      return;
    }
    const { serverUri, deviceName } = Remotv.normalizeSettings(saved);
    if (!await browser.permissions.contains({ origins: [Remotv.serverPermission(serverUri)] })) {
      await browser.runtime.openOptionsPage();
      throw new Error("Save your settings to grant access to the server.");
    }
    if (!tab?.url) throw new Error("Firefox did not provide this tab's URL.");
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 10000);
    try {
      const response = await fetch(`${serverUri}/remote/${encodeURIComponent(deviceName)}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ action: "browser", args: tab.url }),
        credentials: "omit",
        redirect: "error",
        signal: controller.signal,
      });
      if (!response.ok) throw new Error(`Server returned HTTP ${response.status}.`);
    } finally {
      clearTimeout(timeout);
    }
    await showStatus(`Sent to ${deviceName}`);
  } catch (error) {
    await showStatus(error.name === "AbortError" ? "Server request timed out." : error.message, true);
  }
}

browser.menus.onClicked.addListener(openOnDevice);
browser.browserAction.onClicked.addListener(() => browser.runtime.openOptionsPage());
browser.storage.onChanged.addListener((_changes, area) => {
  if (area === "local") updateMenu().catch(console.error);
});
updateMenu().catch(console.error);
