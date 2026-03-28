# getdevice-api

A simple Go API that returns comprehensive device information including device ID, hostname, OS details, network information, and more.

## Features

- **Auto-generated IDs**: Automatically generates UUID for device_id if not configured
- **Auto-generated Client Key**: Creates a random client key if not provided
- **Network Information**: Retrieves MAC address and local IP address
- **System Information**: Returns OS name, architecture, and hostname
- **Environment Configuration**: Use `.env` file to customize device settings
- **Health Check Endpoint**: Simple health check endpoint for monitoring
- **RESTful API**: Clean JSON response format

## Installation

### Prerequisites

- Go 1.21 or higher
- Git

### Clone the Repository

```bash
git clone <repository-url>
cd getdevice-api
```

### Install Dependencies

```bash
go mod download
```

### Build the Application

```bash
go build -o getdevice-api
```

### Run the Application

```bash
./getdevice-api
```

Or run directly:

```bash
go run main.go
```

## Usage

### Configuration

Create a `.env` file in the project root:

```env
# Device Configuration
DEVICE_ID=          # Optional: Leave empty to auto-generate UUID
DEVICE_NAME=my-device
CLIENT_KEY=         # Optional: Leave empty to auto-generate client key

# Server Configuration
PORT=8080
```

### Run with Custom Configuration

```bash
# Using .env file
./getdevice-api

# With environment variables
DEVICE_ID=my-uuid DEVICE_NAME=my-device PORT=3000 ./getdevice-api
```

## API Documentation

### Health Check

**Endpoint:** `GET /health`

Returns the current health status of the API.

**Response:**
```json
{
  "success": true,
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Get Device Information

**Endpoint:** `GET /getdevice`

Returns comprehensive device information.

**Response:**
```json
{
  "success": true,
  "data": {
    "device_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "device_name": "my-device",
    "client_key": "abc123def456...",
    "hostname": "myhost",
    "os": "linux",
    "arch": "amd64",
    "mac_address": "00:11:22:33:44:55",
    "ip_address": "192.168.1.100",
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `success` | boolean | Whether the request was successful |
| `data.device_id` | string | Unique device identifier (UUID) |
| `data.device_name` | string | Name of the device |
| `data.client_key` | string | Unique client identifier |
| `data.hostname` | string | System hostname |
| `data.os` | string | Operating system name (lowercase) |
| `data.arch` | string | System architecture |
| `data.mac_address` | string | MAC address of the primary network interface |
| `data.ip_address` | string | Local IP address |
| `data.timestamp` | string | ISO8601 formatted timestamp |

## Example Requests

### Using cURL

```bash
# Health check
curl http://localhost:8080/health

# Get device information
curl http://localhost:8080/getdevice

# With verbose output
curl -v http://localhost:8080/getdevice
```

### Using JavaScript/Fetch

```javascript
// Health check
fetch('http://localhost:8080/health')
  .then(res => res.json())
  .then(data => console.log(data));

// Get device information
fetch('http://localhost:8080/getdevice')
  .then(res => res.json())
  .then(data => console.log(data));
```

### Using Python/Requests

```python
import requests

# Health check
response = requests.get('http://localhost:8080/health')
print(response.json())

# Get device information
response = requests.get('http://localhost:8080/getdevice')
print(response.json())
```

## Project Structure

```
getdevice-api/
├── main.go              # Application entry point
├── go.mod               # Go module definition
├── go.sum               # Go dependencies checksum
├── .env.example         # Example environment variables
├── .gitignore           # Git ignore rules
├── README.md            # This file
├── handlers/            # HTTP request handlers
│   └── device.go
├── models/              # Data models
│   └── device.go
├── services/            # Business logic
│   ├── device.go
│   └── id_generator.go
└── utils/               # Utility functions
    ├── env.go
    └── network.go
```

## Response Example

### Full Device Info Response

```json
{
  "success": true,
  "data": {
    "device_id": "550e8400-e29b-41d4-a716-446655440000",
    "device_name": "production-server",
    "client_key": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
    "hostname": "web-server-01",
    "os": "linux",
    "arch": "amd64",
    "mac_address": "D4:BE:D9:12:34:56",
    "ip_address": "10.0.0.15",
    "timestamp": "2024-01-15T10:30:45Z"
  }
}
```

## Error Handling

The API returns proper HTTP status codes:

- `200 OK`: Successful request
- `500 Internal Server Error`: Server error (if any)

Error responses (if any) will follow this format:

```json
{
  "success": false,
  "error": "Error message"
}
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.

## Author

Created as a simple device information API service.
