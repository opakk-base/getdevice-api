package services

import (
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"getdevice-api/models"
)

// DeviceService handles device information logic
type DeviceService struct{}

var (
	deviceServiceInstance *DeviceService
	once                  sync.Once
)

// GetDeviceService returns the singleton DeviceService instance
func GetDeviceService() *DeviceService {
	once.Do(func() {
		deviceServiceInstance = &DeviceService{}
	})
	return deviceServiceInstance
}

// GetDeviceName returns the device name from environment or hostname
func (s *DeviceService) GetDeviceName() string {
	// Try to get from environment first
	if name := os.Getenv("DEVICE_NAME"); name != "" {
		return name
	}
	return "unknown-device"
}

// GetOSInfo returns the OS and architecture
func (s *DeviceService) GetOSInfo() (string, string) {
	return runtime.GOOS, runtime.GOARCH
}

// GetHostname returns the system hostname
func (s *DeviceService) GetHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

// GetLocalIP returns the local IP address
func (s *DeviceService) GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return "unknown"
}

// GetMACAddress returns the MAC address of the first non-loopback interface
func (s *DeviceService) GetMACAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "unknown"
	}
	for _, iface := range interfaces {
		// Skip loopback and down interfaces
		if iface.Flags&net.FlagLoopback == 0 &&
			iface.Flags&net.FlagUp != 0 &&
			iface.HardwareAddr != nil {
			return iface.HardwareAddr.String()
		}
	}
	return "unknown"
}

// BuildDeviceInfo builds the complete device info structure
func (s *DeviceService) BuildDeviceInfo(deviceID, clientKey string) *models.DeviceInfo {
	hostname := s.GetHostname()
	osName, arch := s.GetOSInfo()
	ipAddress := s.GetLocalIP()
	macAddress := s.GetMACAddress()

	return &models.DeviceInfo{
		DeviceID:     deviceID,
		DeviceName:   s.GetDeviceName(),
		ClientKey:    clientKey,
		Hostname:     hostname,
		OS:           strings.ToLower(osName),
		Arch:         arch,
		MACAddress:   strings.ToUpper(macAddress),
		IPAddress:    ipAddress,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}
}
