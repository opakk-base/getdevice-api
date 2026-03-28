package services

import (
	"os"
	"runtime"
	"time"

	"getdevice-api/utils"
)

// DeviceService handles device information collection
type DeviceService struct {
	idGenerator *IDGenerator
	envPath     string
}

// NewDeviceService creates a new device service
func NewDeviceService(envPath string) *DeviceService {
	return &DeviceService{
		idGenerator: NewIDGenerator(envPath),
		envPath:     envPath,
	}
}

// DeviceInfo represents device information
type DeviceInfo struct {
	DeviceID    string `json:"device_id"`
	DeviceName  string `json:"device_name"`
	ClientKey   string `json:"client_key"`
	Hostname    string `json:"hostname"`
	OS          string `json:"os"`
	Arch        string `json:"arch"`
	MACAddress  string `json:"mac_address"`
	IPAddress   string `json:"ip_address"`
	Timestamp   string `json:"timestamp"`
}

// GetDeviceInfo collects and returns device information
func (s *DeviceService) GetDeviceInfo() (*DeviceInfo, error) {
	// Get or generate device ID
	deviceID, err := s.idGenerator.GetOrCreateDeviceID()
	if err != nil {
		return nil, err
	}

	// Get or generate client key
	clientKey, err := s.idGenerator.GetOrCreateClientKey()
	if err != nil {
		return nil, err
	}

	// Get device name from env or use hostname
	deviceName := os.Getenv("DEVICE_NAME")
	if deviceName == "" {
		deviceName, _ = os.Hostname()
	}

	// Get hostname
	hostname, _ := os.Hostname()

	// Get network info
	macAddress := utils.GetMACAddress()
	ipAddress := utils.GetLocalIP()

	// Create device info
	info := &DeviceInfo{
		DeviceID:   deviceID,
		DeviceName: deviceName,
		ClientKey:  clientKey,
		Hostname:   hostname,
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		MACAddress: macAddress,
		IPAddress:  ipAddress,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}

	return info, nil
}
