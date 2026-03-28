package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// IDGenerator handles device ID and client key generation
type IDGenerator struct {
	envPath string
}

// NewIDGenerator creates a new ID generator
func NewIDGenerator(envPath string) *IDGenerator {
	return &IDGenerator{
		envPath: envPath,
	}
}

// GenerateUUID generates a new UUID v4
func (g *IDGenerator) GenerateUUID() string {
	return uuid.New().String()
}

// GenerateClientKey generates a random client key using SHA256
func (g *IDGenerator) GenerateClientKey() string {
	uuid := g.GenerateUUID()
	hash := sha256.Sum256([]byte(uuid))
	return hex.EncodeToString(hash[:])
}

// GetOrCreateDeviceID gets device ID from env or generates and saves it
func (g *IDGenerator) GetOrCreateDeviceID() (string, error) {
	// Load .env
	env, err := godotenv.Read(g.envPath)
	if err != nil {
		// If .env doesn't exist, create it
		env = make(map[string]string)
	}

	deviceID := strings.TrimSpace(env["DEVICE_ID"])

	// Generate if empty
	if deviceID == "" {
		deviceID = g.GenerateUUID()
		env["DEVICE_ID"] = deviceID
		
		// Save to .env
		if err := g.saveEnv(env); err != nil {
			return deviceID, fmt.Errorf("failed to save device_id: %w", err)
		}
	}

	return deviceID, nil
}

// GetOrCreateClientKey gets client key from env or generates and saves it
func (g *IDGenerator) GetOrCreateClientKey() (string, error) {
	// Load .env
	env, err := godotenv.Read(g.envPath)
	if err != nil {
		env = make(map[string]string)
	}

	clientKey := strings.TrimSpace(env["CLIENT_KEY"])

	// Generate if empty
	if clientKey == "" {
		clientKey = g.GenerateClientKey()
		env["CLIENT_KEY"] = clientKey
		
		// Save to .env
		if err := g.saveEnv(env); err != nil {
			return clientKey, fmt.Errorf("failed to save client_key: %w", err)
		}
	}

	return clientKey, nil
}

// saveEnv saves environment variables to .env file
func (g *IDGenerator) saveEnv(env map[string]string) error {
	// Create .env content
	var content strings.Builder
	
	// Write header
	content.WriteString("# Device Configuration\n")
	content.WriteString("# Auto-generated - Do not edit manually unless needed\n\n")
	
	// Write all env vars
	for key, value := range env {
		content.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	}
	
	// Write to file
	return os.WriteFile(g.envPath, []byte(content.String()), 0644)
}
