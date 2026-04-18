package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"getdevice-api/handlers"
	"getdevice-api/middleware"
	"getdevice-api/services"
)

// allFields is the ordered list of all device info field keys
var allFields = []string{
	"device_id",
	"device_name",
	"client_key",
	"hostname",
	"os",
	"arch",
	"mac_address",
	"ip_address",
	"timestamp",
}

// ServerStatus represents the current state of the HTTP server
type ServerStatus struct {
	Running bool   `json:"running"`
	Port    string `json:"port"`
}

// App struct provides Wails bindings to the frontend
type App struct {
	ctx           context.Context
	service       *services.DeviceService
	server        *http.Server
	port          string
	running       bool
	exposedFields map[string]bool
	mu            sync.Mutex
}

// NewApp creates a new App instance with all fields exposed by default
func NewApp(service *services.DeviceService, port string) *App {
	exposed := make(map[string]bool)
	for _, f := range allFields {
		exposed[f] = true
	}
	return &App{
		service:       service,
		port:          port,
		exposedFields: exposed,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Auto-start the HTTP server on launch
	a.StartServer()
}

// ---------- Device Info ----------

// GetDeviceInfo returns the full device information
func (a *App) GetDeviceInfo() (*services.DeviceInfo, error) {
	return a.service.GetDeviceInfo()
}

// GetFilteredDeviceInfo returns device info with only exposed fields
func (a *App) GetFilteredDeviceInfo() (map[string]interface{}, error) {
	info, err := a.service.GetDeviceInfo()
	if err != nil {
		return nil, err
	}
	return a.filterFields(info), nil
}

// filterFields converts DeviceInfo to a map containing only exposed fields
func (a *App) filterFields(info *services.DeviceInfo) map[string]interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Marshal to JSON then unmarshal to map to get all fields dynamically
	raw, _ := json.Marshal(info)
	var full map[string]interface{}
	json.Unmarshal(raw, &full)

	filtered := make(map[string]interface{})
	for key, val := range full {
		if a.exposedFields[key] {
			filtered[key] = val
		}
	}
	return filtered
}

// ---------- Field Exposure ----------

// SetExposedFields updates which fields are exposed in the API response
func (a *App) SetExposedFields(fields []string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Reset all to false
	for k := range a.exposedFields {
		a.exposedFields[k] = false
	}
	// Enable only the provided fields
	for _, f := range fields {
		a.exposedFields[f] = true
	}
}

// GetExposedFields returns the list of currently exposed field names
func (a *App) GetExposedFields() []string {
	a.mu.Lock()
	defer a.mu.Unlock()

	var result []string
	for _, f := range allFields {
		if a.exposedFields[f] {
			result = append(result, f)
		}
	}
	return result
}

// ---------- Server Lifecycle ----------

// StartServer starts the HTTP API server on the current port
func (a *App) StartServer() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.running {
		return fmt.Errorf("server is already running")
	}

	mux := http.NewServeMux()

	// Create a handler that uses the app's field filter
	deviceHandler := handlers.NewDeviceHandler(a.service)

	// Filtered /getdevice endpoint — respects exposed fields
	mux.HandleFunc("/getdevice", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		info, err := a.service.GetDeviceInfo()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filtered := a.filterFields(info)
		response := map[string]interface{}{
			"success": true,
			"data":    filtered,
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})

	mux.HandleFunc("/health", deviceHandler.HealthCheck)

	a.server = &http.Server{
		Addr:    ":" + a.port,
		Handler: middleware.CORS(mux),
	}

	a.running = true

	go func() {
		fmt.Printf("HTTP API server started on port %s\n", a.port)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
			a.mu.Lock()
			a.running = false
			a.mu.Unlock()
		}
	}()

	return nil
}

// StopServer gracefully stops the HTTP API server
func (a *App) StopServer() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.running || a.server == nil {
		return fmt.Errorf("server is not running")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := a.server.Shutdown(ctx)
	a.running = false
	a.server = nil
	fmt.Printf("HTTP API server stopped\n")
	return err
}

// SetPort changes the server port. If the server is running, it restarts on the new port.
func (a *App) SetPort(port string) error {
	// Validate port
	p, err := strconv.Atoi(port)
	if err != nil || p < 1 || p > 65535 {
		return fmt.Errorf("invalid port: must be 1-65535")
	}

	wasRunning := false
	a.mu.Lock()
	wasRunning = a.running
	a.mu.Unlock()

	if wasRunning {
		a.StopServer()
	}

	a.mu.Lock()
	a.port = port
	a.mu.Unlock()

	if wasRunning {
		return a.StartServer()
	}
	return nil
}

// GetPort returns the current HTTP server port
func (a *App) GetPort() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.port
}

// GetServerStatus returns the current server status
func (a *App) GetServerStatus() ServerStatus {
	a.mu.Lock()
	defer a.mu.Unlock()
	return ServerStatus{
		Running: a.running,
		Port:    a.port,
	}
}
