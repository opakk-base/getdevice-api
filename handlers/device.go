package handlers

import (
	"encoding/json"
	"net/http"

	"getdevice-api/models"
	"getdevice-api/services"
)

// DeviceHandler handles device-related HTTP requests
type DeviceHandler struct {
	deviceService *services.DeviceService
}

// NewDeviceHandler creates a new DeviceHandler instance
func NewDeviceHandler(deviceService *services.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
	}
}

// GetDevice handles the /getdevice endpoint
func (h *DeviceHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get device info
	deviceInfo, err := h.deviceService.GetDeviceInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
