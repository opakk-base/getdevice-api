package utils

import (
	"net"
)

// GetLocalIP returns the local IP address of the machine
func GetLocalIP() string {
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

// GetMACAddress returns the MAC address of the first available network interface
func GetMACAddress() string {
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
