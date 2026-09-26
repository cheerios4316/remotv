const { test } = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const vm = require("node:vm");

function extension(settings = { serverUri: "localhost:8080/", deviceName: "Living Room" }) {
  const state = { settings, requests: [], menus: [], optionsOpened: 0, permitted: true };
  const browser = {
    storage: {
      local: { get: async () => state.settings },
      onChanged: { addListener: (listener) => { state.changed = listener; } },
    },
    menus: {
      create: (menu) => state.menus.push(menu),
      update: async (_id, menu) => state.menus.push(menu),
      onClicked: { addListener: (listener) => { state.click = listener; } },
    },
    permissions: { contains: async () => state.permitted },
    runtime: { openOptionsPage: async () => { state.optionsOpened++; } },
    browserAction: {
      setBadgeText: async ({ text }) => { state.badge = text; },
      setBadgeBackgroundColor: async () => {},
      setTitle: async ({ title }) => { state.title = title; },
      onClicked: { addListener: () => {} },
    },
  };
  const context = vm.createContext({
    browser, URL, AbortController, setTimeout, clearTimeout, console,
    fetch: async (url, options) => {
      state.requests.push({ url, options });
      if (state.networkError) throw new Error("Network unavailable");
      return { ok: !state.httpError, status: state.httpError || 200 };
    },
  });
  for (const file of ["settings.js", "background.js"]) {
    vm.runInContext(fs.readFileSync(path.join(__dirname, "..", file), "utf8"), context);
  }
  state.api = context.Remotv;
  return state;
}

test("posts the clicked tab URL with the exact command body and encoded device name", async () => {
  const state = extension();
  await state.click({ menuItemId: "remotv-open", pageUrl: "https://wrong.example/" }, { url: "https://example.com/?x=1#part", active: false });
  assert.equal(state.requests.length, 1);
  const { url, options } = state.requests[0];
  assert.equal(url, "http://localhost:8080/remote/Living%20Room");
  assert.equal(options.method, "POST");
  assert.equal(options.headers["Content-Type"], "application/json");
  assert.deepEqual(JSON.parse(options.body), { action: "browser", args: "https://example.com/?x=1#part" });
  assert.equal(state.badge, "OK");
});

test("opens configuration without sending when settings or server permission are missing", async () => {
  for (const configured of [false, true]) {
    const state = configured ? extension() : extension({});
    state.permitted = false;
    await state.click({ menuItemId: "remotv-open" }, { url: "https://example.com" });
    assert.equal(state.optionsOpened, 1);
    assert.equal(state.requests.length, 0);
  }
});

test("reports HTTP and network failures", async () => {
  for (const failure of ["httpError", "networkError"]) {
    const state = extension();
    state[failure] = failure === "httpError" ? 404 : true;
    await state.click({ menuItemId: "remotv-open" }, { url: "https://example.com" });
    assert.equal(state.badge, "!");
    assert.match(state.title, failure === "httpError" ? /404/ : /Network unavailable/);
  }
});

test("updates the menu when settings change", async () => {
  const state = extension();
  await new Promise(setImmediate);
  assert.equal(state.menus.at(-1).title, "Open on Living Room");
  state.settings = { serverUri: "localhost:8080", deviceName: "Office" };
  state.changed({}, "local");
  await new Promise(setImmediate);
  assert.equal(state.menus.at(-1).title, "Open on Office");
});

test("normalizes server addresses and rejects invalid configuration", () => {
  const { api } = extension();
  assert.equal(api.normalizeSettings({ serverUri: " https://example.com/base/ ", deviceName: " TV " }).serverUri, "https://example.com/base");
  assert.equal(api.serverPermission("http://localhost:8080"), "http://localhost/*");
  for (const serverUri of ["", "ftp://example.com", "http://user:pass@example.com", "https://example.com?x=1"]) {
    assert.throws(() => api.normalizeSettings({ serverUri, deviceName: "TV" }));
  }
  for (const deviceName of ["", "..", "a/b", "a\\b"]) {
    assert.throws(() => api.normalizeSettings({ serverUri: "localhost:8080", deviceName }));
  }
});
