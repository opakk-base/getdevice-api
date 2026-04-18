# GetDevice

A native desktop application + REST API that returns comprehensive device information including device ID, hostname, OS details, network information, and more. Built with Go and [Wails](https://wails.io/).

## Features

- **Native Desktop App**: Server control panel with tabbed UI (macOS, Windows, Linux)
- **Background HTTP API**: REST API runs on a configurable port alongside the desktop UI
- **Field Exposure Control**: Choose which fields are included in the API response via Settings
- **Start/Stop Server**: Toggle the HTTP server on/off without closing the app
- **Custom Port**: Change the server port on the fly — auto-restarts if running
- **Close Behavior**: Choose to exit the app or minimize to tray when closing the window
- **Persistent Device ID**: Auto-generates UUID v4 on first run, saves to `.env` for persistence
- **Persistent Client Key**: Auto-generates SHA256 key on first run, remains same across restarts
- **Persistent Settings**: Close behavior and other settings are saved to `.env`
- **Single Binary**: Everything bundled into one executable — no separate frontend/backend ports

## Quick Start

### Prerequisites

- Go 1.22+
- [Wails CLI](https://wails.io/docs/gettingstarted/installation): `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- Xcode Command Line Tools (macOS): `xcode-select --install`
- Node.js 15+ (for Wails build tooling)

### 1. Clone & Install

```bash
git clone https://github.com/opakk-base/getdevice-api.git
cd getdevice-api
go mod download
```

### 2. Development Mode

```bash
wails dev
```

This opens the native window with hot-reload. The HTTP API also runs in the background.

### 3. Build for Production

```bash
wails build
```

Output: `build/bin/GetDevice` (or `GetDevice.app` on macOS)

### 4. Run

```bash
# macOS
./build/bin/getdevice-api.app/Contents/MacOS/GetDevice

# Linux/Windows
./build/bin/GetDevice
```

You can override the port via environment variable:

```bash
PORT=9090 ./build/bin/getdevice-api.app/Contents/MacOS/GetDevice
```

### 5. Test the API

```bash
curl http://localhost:8080/getdevice
curl http://localhost:8080/health
```

## App UI

The app has three tabs:

### Device Tab (Main)

Server control panel:
- **Port input** — set the HTTP server port (auto-restarts on change)
- **Start/Stop button** — toggle the HTTP server on/off
- **Status indicator** — shows whether the server is running and on which port

### Settings Tab

- **Exposed Fields** — checkboxes to control which fields are included in the `/getdevice` API response (device_id, device_name, client_key, hostname, os, arch, mac_address, ip_address, timestamp)
- **Close Behavior** — choose what happens when you close the window:
  - *Exit application* — quits the app and stops the server
  - *Minimize to tray* — hides the window but keeps the HTTP server running in the background; click the dock icon to show the window again

### About Tab

App name, version, author, license, and a link to the GitHub repository.

## Architecture

```
┌─────────────────────────────────────────────┐
│              Single Binary                  │
│                                             │
│  1. Background HTTP Server (:8080)          │
│     ├── GET /health      (API)              │
│     ├── GET /getdevice   (API, filtered)    │
│     └── CORS middleware                     │
│                                             │
│  2. Wails Native Window                     │
│     ├── Device tab (server controls)        │
│     ├── Settings tab (fields, behavior)     │
│     ├── About tab (app info)                │
│     └── Go bindings (IPC, no HTTP)          │
└─────────────────────────────────────────────┘
```

- The **desktop UI** calls Go functions directly via Wails bindings (IPC, not HTTP)
- The **HTTP API** runs in a background goroutine on the configured port
- The `/getdevice` endpoint only returns fields that are enabled in Settings
- External services can call `/getdevice` and `/health` as before

## Configuration

### .env File

Create a `.env` file in the project root (or let the app auto-generate one on first run):

```env
# Device Configuration
# Leave DEVICE_ID and CLIENT_KEY empty to auto-generate on first run
DEVICE_ID=
DEVICE_NAME=my-device
CLIENT_KEY=

# Server Configuration
PORT=8080

# App Settings
# Close behavior: "exit" (quit app) or "minimize" (hide window, server keeps running)
CLOSE_BEHAVIOR=exit
```

### Auto-Generation (Recommended)

Leave `DEVICE_ID` and `CLIENT_KEY` empty — they will be:
1. **Auto-generated** on first run
2. **Auto-saved** to `.env` file
3. **Persistent** across restarts

## API Documentation

### GET /getdevice

Returns device information. Only fields enabled in Settings are included in the response.

```bash
curl http://localhost:8080/getdevice
```

```json
{
  "success": true,
  "data": {
    "device_id": "3fac064c-f7ef-4bad-812d-15607a6c61ef",
    "device_name": "my-device",
    "client_key": "d8edd98a85d248633276b463415419b41f12611c393c957109e32e70b123d422",
    "hostname": "my-mac.local",
    "os": "darwin",
    "arch": "arm64",
    "mac_address": "D4:BE:D9:12:34:56",
    "ip_address": "192.168.1.100",
    "timestamp": "2026-04-18T06:40:53Z"
  }
}
```

If some fields are unchecked in Settings, they will be omitted from the `data` object.

### GET /health

```bash
curl http://localhost:8080/health
```

```json
{
  "success": true,
  "status": "healthy",
  "timestamp": "2026-04-18T06:40:53Z"
}
```

## Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `device_id` | string | Unique device identifier (UUID v4) |
| `device_name` | string | Device name (from .env or hostname) |
| `client_key` | string | Client authentication key (SHA256 hash) |
| `hostname` | string | System hostname |
| `os` | string | Operating system (darwin/linux/windows) |
| `arch` | string | Architecture (amd64/arm64) |
| `mac_address` | string | MAC address of primary network interface |
| `ip_address` | string | Local IP address |
| `timestamp` | string | ISO8601 timestamp (UTC) |

## Project Structure

```
getdevice-api/
├── main.go                 # Entry point (Wails app config, OnBeforeClose)
├── app.go                  # Wails bindings (server lifecycle, settings, field filtering)
├── wails.json              # Wails project configuration
├── frontend/
│   └── src/
│       ├── index.html      # UI layout (3 tabs: Device, Settings, About)
│       ├── style.css       # Styling
│       └── main.js         # Frontend logic (tab switching, server controls, settings)
├── handlers/
│   └── device.go           # HTTP handlers
├── middleware/
│   └── cors.go             # CORS middleware
├── models/
│   └── device.go           # Response models
├── services/
│   ├── device.go           # Device info service
│   └── id_generator.go     # ID generation & persistence
├── utils/
│   ├── env.go              # Environment utilities
│   └── network.go          # Network utilities
├── .env.example            # Environment template
├── go.mod                  # Go module
└── README.md               # This file
```

## Dependencies

```go
require (
    github.com/google/uuid v1.6.0           // UUID generation
    github.com/joho/godotenv v1.5.1         // .env file loader
    github.com/wailsapp/wails/v2 v2.12.0    // Desktop app framework
)
```

## Troubleshooting

### Port Already in Use

```bash
# Use a different port
PORT=9090 ./build/bin/getdevice-api.app/Contents/MacOS/GetDevice

# Or kill the process on the port
lsof -ti:8080 | xargs kill -9
```

You can also change the port from the app UI on the Device tab.

### Wails Not Found

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verify all dependencies
wails doctor
```

### .env File Not Found

Warning is normal — the app will create `.env` on first run with auto-generated IDs.

### Window Closed but Server Still Running

If you set Close Behavior to "Minimize to tray" in Settings, closing the window hides it instead of quitting. The HTTP server keeps running. Click the dock icon (macOS) to show the window again, or quit from the app menu.

## License

MIT License

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## Support

- **Issues:** https://github.com/opakk-base/getdevice-api/issues
- **Repository:** https://github.com/opakk-base/getdevice-api
