package main

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v3/pkg/application"

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

type BackgroundSettings struct {
	RunInBackground bool `json:"runInBackground"`
	StartAtLogin    bool `json:"startAtLogin"`
}

// AppInfo holds metadata for the About page
type AppInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Author      string `json:"author"`
	Description string `json:"description"`
	License     string `json:"license"`
	GitHub      string `json:"github"`
}

// App struct provides Wails bindings to the frontend
type App struct {
	service         *services.DeviceService
	server          *http.Server
	port            string
	running         bool
	exposedFields   map[string]bool
	password        string // API access password
	envPath         string
	autostart       *application.AutostartManager
	runInBackground bool
	mu              sync.Mutex
}

// NewApp creates a new App instance with all fields exposed by default
func NewApp(service *services.DeviceService, port string, password string, envPath string, autostart *application.AutostartManager, runInBackground bool) *App {
	exposed := make(map[string]bool)
	for _, f := range allFields {
		exposed[f] = f != "client_key"
	}
	if password == "" {
		password = "12345678"
	}
	return &App{
		service:         service,
		port:            port,
		exposedFields:   exposed,
		password:        password,
		envPath:         envPath,
		autostart:       autostart,
		runInBackground: runInBackground,
	}
}

func (a *App) GetBackgroundSettings() (BackgroundSettings, error) {
	a.mu.Lock()
	runInBackground := a.runInBackground
	a.mu.Unlock()

	startAtLogin := false
	if a.autostart != nil {
		var err error
		startAtLogin, err = a.autostart.IsEnabled()
		if err != nil {
			return BackgroundSettings{}, err
		}
	}
	return BackgroundSettings{RunInBackground: runInBackground, StartAtLogin: startAtLogin}, nil
}

func (a *App) SetRunInBackground(enabled bool) error {
	value := strconv.FormatBool(enabled)
	if err := a.saveEnvKey("RUN_IN_BACKGROUND", value); err != nil {
		return err
	}
	a.mu.Lock()
	a.runInBackground = enabled
	a.mu.Unlock()
	return nil
}

func (a *App) SetStartAtLogin(enabled bool) error {
	if a.autostart == nil {
		return fmt.Errorf("automatic startup is unavailable")
	}
	if enabled {
		return a.autostart.EnableWithOptions(application.AutostartOptions{
			Identifier: "com.opakk.pegasus",
			Arguments:  []string{"--background"},
		})
	}
	return a.autostart.Disable()
}

func (a *App) runInBackgroundEnabled() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.runInBackground
}

// ServiceStartup starts the API when Pegasus launches.
func (a *App) ServiceStartup(context.Context, application.ServiceOptions) error {
	return a.StartServer()
}

// ServiceShutdown stops the API when Pegasus quits from its tray menu.
func (a *App) ServiceShutdown() error {
	if a.GetServerStatus().Running {
		return a.StopServer()
	}
	return nil
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

	for k := range a.exposedFields {
		a.exposedFields[k] = false
	}
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
	deviceHandler := handlers.NewDeviceHandler(a.service)

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
		Addr:    "127.0.0.1:" + a.port,
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

// ---------- App Security ----------

// VerifyPassword checks the provided password against the current app password
// using constant-time comparison.
func (a *App) VerifyPassword(password string) bool {
	a.mu.Lock()
	current := a.password
	a.mu.Unlock()

	return len(password) == len(current) &&
		subtle.ConstantTimeCompare([]byte(password), []byte(current)) == 1
}

// ChangePassword verifies the current password and, if valid, sets a new one
// and persists it to .env.
func (a *App) ChangePassword(currentPassword, newPassword string) error {
	newPassword = strings.TrimSpace(newPassword)
	if len(newPassword) < 4 {
		return fmt.Errorf("new password must be at least 4 characters")
	}

	if !a.VerifyPassword(currentPassword) {
		return fmt.Errorf("current password is incorrect")
	}

	a.mu.Lock()
	a.password = newPassword
	a.mu.Unlock()

	return a.saveEnvKey("PASSWORD", newPassword)
}

// GetConfigPath returns the absolute path to the .env configuration file
func (a *App) GetConfigPath() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.envPath
}

// saveEnvKey reads the .env file, sets a key, and writes it back
func (a *App) saveEnvKey(key, value string) error {
	env, err := godotenv.Read(a.envPath)
	if err != nil {
		env = make(map[string]string)
	}

	env[key] = value

	var content strings.Builder
	content.WriteString("# Device Configuration\n")
	content.WriteString("# Auto-generated - Do not edit manually unless needed\n\n")

	for k, v := range env {
		content.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}

	return os.WriteFile(a.envPath, []byte(content.String()), 0644)
}

// ---------- About ----------

// GetAppInfo returns app metadata for the About page
func (a *App) GetAppInfo() AppInfo {
	return AppInfo{
		Name:        "Pegasus",
		Version:     "1.4.0",
		Author:      "opakk",
		Description: "A native desktop application that exposes device information via a REST API on a single port.",
		License:     "MIT",
		GitHub:      "https://github.com/opakk-base/getdevice-api",
	}
}
