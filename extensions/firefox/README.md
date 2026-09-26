# Remotv for Firefox

## Build

Requires Node.js 20 or newer and npm. From the repository root:

```sh
npm --prefix extensions/firefox ci
npm --prefix extensions/firefox run build
```

Or run `make build-firefox`. This validates the extension and creates
`extensions/firefox/dist/remotv-unsigned.xpi`. Tests, dependencies, and build
configuration are excluded. The XPI is a ZIP archive with the manifest at its root.

## Install permanently in standard Firefox

1. Download the latest signed XPI file from the releases page.
2. In Firefox, open `about:addons`, click the gear menu, choose **Install Add-on
   From File…**, and select the signed XPI. It remains installed after restarting.

This extension requires Firefox 140 or newer. Its installation disclosure declares
browsing activity because clicking the menu sends the selected tab's URL to your
configured server. There is no analytics service.

## Load and configure

1. Click the toolbar button to open settings. You can also use **Preferences** from `about:addons`.
2. Enter the server URI, including its HTTP port, and the exact connected device name. For example: `http://192.168.1.10:8080` and `living-room`.
3. Save and allow access to your server when Firefox prompts.
4. Right-click a tab or a page and select **Open on living-room**.

Device names are URL-encoded. A URI without a scheme defaults to HTTP. Settings are stored locally in Firefox. The menu updates when the device name changes. Before configuration, the menu opens settings.

The toolbar badge shows `OK` when the server accepts the command or `!` on failure; hover over it for details. Requests time out after 10 seconds. Acceptance means the server accepted the request, not confirmation that the remote browser opened it.

Permissions use Firefox's [activeTab and host permissions](https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/manifest.json/permissions). Server access is requested on save through [optional permissions](https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/manifest.json/optional_permissions); no server CORS change is required. Previously granted server permissions remain until revoked in Firefox.

## Check

From the repository root, run `node --test extensions/firefox/tests/*.test.cjs`.

For a manual smoke test, configure a connected device, send a foreground and background tab, change the device name and check the menu, then try an unavailable server and check the failure badge.
