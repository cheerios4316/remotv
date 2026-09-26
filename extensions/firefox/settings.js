"use strict";

globalThis.Remotv = {
  normalizeSettings({ serverUri = "", deviceName = "" } = {}) {
    serverUri = serverUri.trim();
    deviceName = deviceName.trim();
    if (!serverUri || !deviceName) {
      throw new Error("Enter a server URI and device name.");
    }
    if (!serverUri.includes("://")) serverUri = `http://${serverUri}`;
    let server;
    try {
      server = new URL(serverUri);
    } catch {
      throw new Error("Enter a valid server URI, such as http://192.168.1.10:8080.");
    }
    if (!["http:", "https:"].includes(server.protocol) ||
        server.username || server.password || server.search || server.hash) {
      throw new Error("Use an HTTP or HTTPS server URI without credentials, a query, or a fragment.");
    }
    if (deviceName === "." || deviceName === ".." || /[\/\\\r\n]/.test(deviceName)) {
      throw new Error("The device name cannot contain slashes or line breaks, or be . or ...");
    }
    return { serverUri: server.href.replace(/\/+$/, ""), deviceName };
  },

  serverPermission(serverUri) {
    const server = new URL(serverUri);
    // Firefox host match patterns omit ports.
    return `${server.protocol}//${server.hostname}/*`;
  },
};
