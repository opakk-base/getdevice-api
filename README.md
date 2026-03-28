# getdevice-api

A simple Go API that returns comprehensive device information including device ID, hostname, OS details, network information, and more.

## ✨ Features

- **Persistent Device ID**: Auto-generates UUID v4 on first run, saves to `.env` for persistence
- **Persistent Client Key**: Auto-generates SHA256 key on first run, remains same across restarts
- **Network Information**: Retrieves MAC address and local IP address
- **System Information**: Returns OS name, architecture, and hostname
- **Environment Configuration**: Use `.env` file to customize device settings
- **Health Check Endpoint**: Simple health check endpoint for monitoring
- **RESTful API**: Clean JSON response format

## 🚀 Quick Start

### 1. Clone & Install

```bash
git clone https://github.com/opakk-base/getdevice-api.git
cd getdevice-api
go mod download
```

### 2. Run

```bash
go run main.go
```

### 3. Test

```bash
curl http://localhost:8080/getdevice
```

## 📦 Installation

### Prerequisites

- Go 1.21 or higher
- Git

### Build

```bash
go build -o getdevice-api
./getdevice-api
```

### Run Directly

```bash
go run main.go
```

## ⚙️ Configuration

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

Leave `DEVICE_ID` and `CLIENT_KEY` empty - they will be:
1. **Auto-generated** on first run
2. **Auto-saved** to `.env` file
3. **Persistent** across restarts

```bash
# First run - generates IDs
go run main.go

# Check .env - IDs are saved!
cat .env

# Restart - same IDs!
go run main.go
```

### Manual Configuration

Set your own values:

```env
DEVICE_ID=my-custom-uuid-123
CLIENT_KEY=my-custom-key-456
DEVICE_NAME=production-server
PORT=8080
```

## 📡 API Documentation

### 1. GET /getdevice

Returns comprehensive device information.

**Request:**
```bash
curl http://localhost:8080/getdevice
```

**Response:**
```json
{
  "success": true,
  "data": {
    "device_id": "3fac064c-f7ef-4bad-812d-15607a6c61ef",
    "device_name": "my-device",
    "client_key": "d8edd98a85d248633276b463415419b41f12611c393c957109e32e70b123d422",
    "hostname": "VM-6-91-opencloudos",
    "os": "linux",
    "arch": "amd64",
    "mac_address": "52:54:00:25:b7:d9",
    "ip_address": "10.11.6.91",
    "timestamp": "2026-03-28T06:40:53Z"
  }
}
```

### 2. GET /health

Health check endpoint.

**Request:**
```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "success": true,
  "status": "healthy",
  "timestamp": "2026-03-28T06:40:53Z"
}
```

## 🔍 Response Fields

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

## 💡 How Persistent IDs Work

### First Run

1. Check `.env` for `DEVICE_ID` and `CLIENT_KEY`
2. Both are **empty** → Generate new UUID and client key
3. **Save** to `.env` file automatically
4. Return generated IDs

### Subsequent Runs

1. Check `.env` for `DEVICE_ID` and `CLIENT_KEY`
2. Both are **present** → Read from `.env`
3. Return saved IDs (**SAME as first run!**)

### Example

```bash
# First run
$ go run main.go
$ curl http://localhost:8080/getdevice | jq .data.device_id
"3fac064c-f7ef-4bad-812d-15607a6c61ef"

# Check .env
$ cat .env | grep DEVICE_ID
DEVICE_ID=3fac064c-f7ef-4bad-812d-15607a6c61ef

# Restart
$ go run main.go
$ curl http://localhost:8080/getdevice | jq .data.device_id
"3fac064c-f7ef-4bad-812d-15607a6c61ef"  # ← SAME!
```

## 🛠️ Development

### Project Structure

```
getdevice-api/
├── main.go                 # Entry point
├── handlers/
│   └── device.go          # HTTP handlers
├── models/
│   └── device.go          # Response models
├── services/
│   ├── device.go          # Device info service
│   └── id_generator.go    # ID generation & persistence
├── utils/
│   ├── env.go             # Environment utilities
│   └── network.go         # Network utilities
├── .env.example           # Environment template
├── go.mod                 # Go module
└── README.md              # This file
```

### Dependencies

```go
require (
    github.com/google/uuid v1.6.0      // UUID generation
    github.com/joho/godotenv v1.5.1    // .env file loader
)
```

### Run Tests

```bash
go test ./...
```

## 📊 Use Cases

### 1. Device Fingerprinting

Identify unique devices in your system:

```json
{
  "device_id": "3fac064c-f7ef-4bad-812d-15607a6c61ef",
  "client_key": "d8edd98a85d248633276b463415419b41f12611c393c957109e32e70b123d422"
}
```

### 2. License Validation

Use `device_id` and `client_key` for software licensing:

```bash
# Send to license server
curl -X POST https://license.example.com/validate \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "3fac064c-f7ef-4bad-812d-15607a6c61ef",
    "client_key": "d8edd98a85d248633276b463415419b41f12611c393c957109e32e70b123d422"
  }'
```

### 3. Analytics Tracking

Track which devices are using your application:

```json
{
  "device_id": "3fac064c-f7ef-4bad-812d-15607a6c61ef",
  "os": "linux",
  "arch": "amd64",
  "timestamp": "2026-03-28T06:40:53Z"
}
```

### 4. Multi-Device Management

Manage multiple installations:

```bash
# Get all devices
curl https://api.example.com/devices

# Response
[
  {
    "device_id": "3fac064c-f7ef-4bad-812d-15607a6c61ef",
    "device_name": "production-server",
    "last_seen": "2026-03-28T06:40:53Z"
  },
  {
    "device_id": "abc-123-def",
    "device_name": "staging-server",
    "last_seen": "2026-03-28T05:30:00Z"
  }
]
```

## 🔧 Troubleshooting

### Port Already in Use

```bash
# Kill process on port 8080
lsof -ti:8080 | xargs kill -9

# Or use different port
PORT=3000 go run main.go
```

### .env File Not Found

Warning is normal - app will create `.env` on first run.

### Device ID Changes After Restart

Make sure `.env` file is writable and IDs are saved:

```bash
# Check .env exists
ls -la .env

# Check IDs are saved
cat .env | grep DEVICE_ID
```

### Permission Denied

```bash
# Make .env writable
chmod 644 .env

# Or run with sudo (not recommended)
sudo go run main.go
```

## 📝 License

MIT License - feel free to use in your projects!

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 📧 Support

- **Issues:** https://github.com/opakk-base/getdevice-api/issues
- **Repository:** https://github.com/opakk-base/getdevice-api

---

**Built with ❤️ using Go**
