# GetDevice

A native desktop application + REST API that returns comprehensive device information including device ID, hostname, OS details, network information, and more. Built with Go and [Wails](https://wails.io/).

## Features

- **Native Desktop App**: Opens a native window showing device info (macOS, Windows, Linux)
- **Background HTTP API**: REST API runs on a single port alongside the desktop UI
- **Persistent Device ID**: Auto-generates UUID v4 on first run, saves to `.env` for persistence
- **Persistent Client Key**: Auto-generates SHA256 key on first run, remains same across restarts
- **Network Information**: Retrieves MAC address and local IP address
- **System Information**: Returns OS name, architecture, and hostname
- **Single Binary**: Everything bundled into one executable — no separate frontend/backend ports
- **Copy to Clipboard**: One-click JSON copy from the desktop UI

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

### 5. Test the API

```bash
curl http://localhost:8080/getdevice
```

## Architecture

```
┌─────────────────────────────────────────────┐
│              Single Binary                  │
│                                             │
│  1. Background HTTP Server (:8080)          │
│     ├── GET /health      (API)              │
│     ├── GET /getdevice   (API)              │
│     └── CORS middleware                     │
│                                             │
│  2. Wails Native Window                     │
│     ├── Device info dashboard               │
│     ├── Go bindings (no HTTP round-trip)    │
│     └── Copy JSON to clipboard              │
└─────────────────────────────────────────────┘
```

- The **desktop UI** calls Go functions directly via Wails bindings (IPC, not HTTP)
- The **HTTP API** runs in a background goroutine on the configured port
- External services can still call `/getdevice` and `/health` as before

## Configuration

### .env File

Create a `.env` file in the project root:

```env
# Device Configuration
# Leave DEVICE_ID and CLIENT_KEY empty to auto-generate on first run
DEVICE_ID=
DEVICE_NAME=my-device
CLIENT_KEY=

# Server Configuration
PORT=8080
```

### Auto-Generation (Recommended)

Leave `DEVICE_ID` and `CLIENT_KEY` empty — they will be:
1. **Auto-generated** on first run
2. **Auto-saved** to `.env` file
3. **Persistent** across restarts

## API Documentation

### GET /getdevice

Returns comprehensive device information.

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
├── main.go                 # Entry point (Wails app + HTTP server)
├── app.go                  # Wails-bound App struct
├── wails.json              # Wails project configuration
├── frontend/
│   └── src/
│       ├── index.html      # Desktop UI
│       ├── style.css       # Styling
│       └── main.js         # Frontend logic (calls Go bindings)
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

### Wails Not Found

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verify
wails doctor
```

### .env File Not Found

Warning is normal — app will create `.env` on first run with auto-generated IDs.

## License

MIT License

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## Support

- **Issues:** https://github.com/opakk-base/getdevice-api/issues
- **Repository:** https://github.com/opakk-base/getdevice-api
