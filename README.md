# Remotv

Remotv is a software used to open URLs on different devices.

The software is composed of:

- a server, that accepts TCP clients and HTTP requests
- n clients, that connect to the TCP server and receive messages
- a browser extension, that performs the HTTP request to the server

The HTTP requests includes the desired device name and the desired URL.

Upon receiving the HTTP request, the server checks to see if a connected client matches the requested device name, and forwards the payload to it.

The client, upon receiving a message, checks its configuration to see if the command is known, and if so, runs it.

## Installation

1. Download the client, the server and the .xpi browser extension from the releases page
2. Run the server on your desired machine, with the adequate flags (described below)
3. Run the client on all clients you want, with the adequate flags and config (also described below). Optionally, you can create a systemd service that runs the client in the background. An example is provided in the repository
4. Open your browser, navigate to `about:addons`
5. Click the gear button and select `Install Add-on from file`, then select the downloaded .xpi file

### Client configuration
#### Command flags
| Flag | Shorthand | Default | Value type | Description |
| --- | --- | --- | --- | ---|
| `device-name` | `D` | `my-device` | `string` | Device display name
| `uri` | `U` | `localhost` | `string` | TCP server URI
| `port` | `P` | `9000` | `int` | TCP server port
| `config` | `C` | `config.json` | `string` | Config JSON file path

#### Structure for `config.json`
```json
{
    "browser": "<browser bin>" // default: 'firefox'
}
```

### Server configuration
#### Command flags
| Flag | Shorthand | Default | Value type | Description |
| --- | --- | --- | --- | ---|
| `tcp` | `T` | `9000` | `int` | TCP listen port
| `http` | `H` | `8080` | `int` | HTTP listen port

## Manual build

### Server + Client

Requirements:
go, make, docker, docker compose

1. Run `make build` or `make build-windows`. Alternatively, run `make build-server` and `make build-client` (both commands have the `-windows` equivalent)
2. Run the server on the desired machine and the client on all desired clients.

### Browser extension
1. CD into `extensions/firefox`
2. Run `npm run build`
3. Open Firefox, navigate to `about:debugging`
4. Navigate to "This Firefox", then click "Load Temporary Add-on"
5. Load the built unsigned .xpi file, found in the `dist` folder


## Future features

Currently only the `browser` command is supported. In the future, all commands will be
configured in the configuration JSON. A reasonable default JSON config will be provided.

Proper config handling will come soon: for now it's just a fun experiment :)

The browser extension was almost entirely built by Codex as a POC: it will be reworked soon if major problems arise.