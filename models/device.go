package models

import "time"

// DeviceInfo represents the device information response
type DeviceInfo struct {
	DeviceID     string `json:"device_id"`
	DeviceName   string `json:"device_name"`
	ClientKey    string `json:"client_key"`
	Hostname     string `json:"hostname"`
	OS           string `json:"os"`
	Arch         string `json:"arch"`
	MACAddress   string `json:"mac_address"`
	IPAddress    string `json:"ip_address"`
	Timestamp    string `json:"timestamp"`
}

// Response represents the API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Success   bool   `json:"success"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// NewHealthResponse creates a new health check response
func NewHealthResponse() *HealthResponse {
	return &HealthResponse{
		Success:   true,
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
