package handlers

import (
	"encoding/json"
	"net/http"

	"getdevice-api/models"
	"getdevice-api/services"
)

// DeviceHandler handles device-related HTTP requests
type DeviceHandler struct{}

// NewDeviceHandler creates a new DeviceHandler instance
func NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{}
}

// GetDevice handles the /getdevice endpoint
func (h *DeviceHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get services
	idGen := services.NewIDGenerator()
	devService := services.GetDeviceService()

	// Generate or use configured device ID and client key
	deviceID := idGen.GenerateDeviceID("")
	clientKey := idGen.GenerateClientKeyFromConfig("")

	// Build device info
	deviceInfo := devService.BuildDeviceInfo(deviceID, clientKey)

	// Create response
	response := models.Response{
		Success: true,
		Data:    deviceInfo,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HealthCheck handles the /health endpoint
func (h *DeviceHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := models.NewHealthResponse()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
