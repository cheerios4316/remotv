"use strict";

const form = document.querySelector("#settings");
const serverInput = document.querySelector("#server-uri");
const deviceInput = document.querySelector("#device-name");
const status = document.querySelector("#status");

function showStatus(message, failed = false) {
  status.textContent = message;
  status.dataset.error = String(failed);
}

form.addEventListener("submit", async (event) => {
  event.preventDefault();
  try {
    const settings = Remotv.normalizeSettings({ serverUri: serverInput.value, deviceName: deviceInput.value });
    // Request directly from the submit gesture, before any other asynchronous work.
    const granted = await browser.permissions.request({ origins: [Remotv.serverPermission(settings.serverUri)] });
    if (!granted) throw new Error("Server access was not granted. Allow access to save these settings.");
    await browser.storage.local.set(settings);
    serverInput.value = settings.serverUri;
    deviceInput.value = settings.deviceName;
    showStatus(`Saved. Right-click a tab or page and choose “Open on ${settings.deviceName}”.`);
  } catch (error) {
    showStatus(error.message, true);
  }
});

browser.storage.local.get(["serverUri", "deviceName"]).then((settings) => {
  serverInput.value = settings.serverUri || "";
  deviceInput.value = settings.deviceName || "";
}).catch((error) => showStatus(error.message, true));
