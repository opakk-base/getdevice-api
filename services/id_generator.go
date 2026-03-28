package services

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"
)

// IDGenerator provides utilities for generating IDs
type IDGenerator struct{}

// NewIDGenerator creates a new IDGenerator instance
func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

// GenerateUUID generates a new UUID v4
func (g *IDGenerator) GenerateUUID() string {
	return uuid.Must(uuid.NewRandom()).String()
}

// GenerateClientKey generates a random client key
func (g *IDGenerator) GenerateClientKey(length int) string {
	if length <= 0 {
		length = 32
	}
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

// GenerateDeviceID generates a device ID from config or auto-generates
func (g *IDGenerator) GenerateDeviceID(configValue string) string {
	if configValue != "" {
		return configValue
	}
	return g.GenerateUUID()
}

// GenerateClientKeyFromConfig generates a client key from config or auto-generates
func (g *IDGenerator) GenerateClientKeyFromConfig(configValue string) string {
	if configValue != "" {
		return configValue
	}
	return g.GenerateClientKey(32)
}
