# Pegasus

A native desktop application + REST API that returns comprehensive device information including device ID, hostname, OS details, network information, and more. Built with Go and [Wails](https://wails.io/).

## Features

- **Native Desktop App**: Server control panel with tabbed UI (macOS, Windows, Linux)
- **Background HTTP API**: REST API runs on a configurable port alongside the desktop UI
- **Field Exposure Control**: Choose which fields are included in the API response via Settings
- **App Lock**: The app requires a password (default `12345678`) to unlock on launch, changeable in Settings
- **Start/Stop Server**: Toggle the HTTP server on/off without closing the app
- **Custom Port**: Change the server port on the fly — auto-restarts if running
- **System Tray**: Choose whether closing the window hides Pegasus in the tray or exits the app
- **Login Startup**: Optionally launch Pegasus automatically when the user signs in
- **Persistent Device ID**: Auto-generates UUID v4 on first run, saves to `.env` for persistence
- **Persistent Client Key**: Auto-generates SHA256 key on first run, remains same across restarts
- **Persistent Settings**: Password, tray preference, and device identity are stored in the user config directory
- **Single Binary**: Everything bundled into one executable — no separate frontend/backend ports

## Install

Release packages are published on the [GitHub Releases](https://github.com/opakk-base/getdevice-api/releases) page when a new version tag is pushed.

## Quick Start

### Prerequisites

- Go 1.26+
- [Wails v3 CLI](https://wails.io/): `go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0`
- Xcode Command Line Tools (macOS): `xcode-select --install`
- GTK4 and WebKitGTK 6.0 development packages (Linux)

### 1. Clone & Install

```bash
git clone https://github.com/opakk-base/getdevice-api.git
cd getdevice-api
go mod download
```

### 2. Run from source

```bash
go run .
```

This opens the native window. The HTTP API also runs in the background.

### 3. Build for Production

```bash
wails3 build
```

Output: `bin/Pegasus` (or `bin/Pegasus.exe` on Windows)

### 4. Run

```bash
./bin/Pegasus
```

You can override the port via environment variable:

```bash
PORT=9090 ./bin/Pegasus
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
- **API links** — open `http://localhost:8080/getdevice` or `http://localhost:8080/health`

### Settings Tab

- **App Password** — change the password required to unlock the app on launch (default `12345678`)
- **Exposed Fields** — checkboxes to control which fields are included in the `/getdevice` API response (device_id, device_name, client_key, hostname, os, arch, mac_address, ip_address, timestamp)
- **Background** — choose whether closing the window keeps Pegasus in the tray and whether Pegasus starts when you log in

Use the tray menu to reopen the window or quit Pegasus. The background settings also control whether Pegasus starts when you log in. On Windows, the app is built as a GUI executable and does not open a command prompt. It runs in the logged-in user's session; it is not installed as a Windows Service.

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
│  2. Pegasus System Tray + Wails Window      │
│     ├── Device tab (server controls)        │
│     ├── Settings tab (fields, behavior)     │
│     ├── About tab (app info)                │
│     └── Go bindings (IPC, no HTTP)          │
└─────────────────────────────────────────────┘
```

- The **desktop UI** calls Go functions directly via Wails v3 bindings (IPC, not HTTP)
- The **HTTP API** runs in a background goroutine on the configured port
- The `/getdevice` endpoint only returns fields that are enabled in Settings
- External services can call `/getdevice` and `/health` as before

## Configuration

### .env File

Pegasus stores `.env` under the user's config directory: `~/Library/Application Support/Pegasus` on macOS, `%AppData%\Pegasus` on Windows, and `~/.config/Pegasus` on Linux. Existing `GetDevice` config is copied on first launch.

For a manual configuration, create `.env` in that directory:

```env
# Device Configuration
# Leave DEVICE_ID and CLIENT_KEY empty to auto-generate on first run
DEVICE_ID=
DEVICE_NAME=my-device
CLIENT_KEY=

# Server Configuration
PORT=8080

# App lock password (default 12345678 if unset)
PASSWORD=12345678

# Keep Pegasus in the tray when the window is closed
RUN_IN_BACKGROUND=true
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
├── main.go                 # Wails app, system tray, and close-to-hide behavior
├── app.go                  # Wails bindings (server lifecycle, settings, field filtering)
├── wails.json              # Wails project configuration (app metadata)
├── Makefile                # Wails v3 build commands
├── Taskfile.yml            # Wails v3 task entry points
├── build/
│   ├── appicon.png         # Source icon for all platforms
│   ├── darwin/             # Info.plist + entitlements.plist (macOS bundle)
│   ├── windows/
│   │   └── installer/      # Pegasus NSIS installer script
│   └── linux/              # nfpm.yaml, .desktop entry, Dockerfile, icon scripts
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
│   ├── config.go           # Config dir resolution & .env migration
│   ├── env.go              # Environment utilities
│   └── network.go          # Network utilities
├── .github/workflows/release.yml  # Tag-triggered packages for all 3 OSes
├── .env.example            # Environment template
├── go.mod                  # Go module
└── README.md               # This file
```

## Dependencies

```go
require (
    github.com/google/uuid v1.6.0           // UUID generation
    github.com/joho/godotenv v1.5.1         // .env file loader
    github.com/wailsapp/wails/v3 v3.0.0     // Desktop app framework and system tray
)
```

## Troubleshooting

### Port Already in Use

```bash
# Use a different port
PORT=9090 ./bin/Pegasus

# Or kill the process on the port
lsof -ti:8080 | xargs kill -9
```

You can also change the port from the app UI on the Device tab. Linux tray visibility depends on the desktop environment providing a StatusNotifier tray host.

### Wails Not Found

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0

# Verify all dependencies
wails3 doctor
```

### .env File Not Found

Warning is normal — the app will create `.env` on first run with auto-generated IDs.

### Window Closed but Server Still Running

When background mode is enabled, closing the window hides Pegasus and keeps the HTTP server running. Open it again from the tray icon. Choose **Quit** in the tray menu to stop the server and exit Pegasus.

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
