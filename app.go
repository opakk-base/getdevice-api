package main

import (
	"context"

	"getdevice-api/services"
)

// App struct provides Wails bindings to the frontend
type App struct {
	ctx     context.Context
	service *services.DeviceService
}

// NewApp creates a new App instance
func NewApp(service *services.DeviceService) *App {
	return &App{
		service: service,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetDeviceInfo returns the current device information.
// This is bound to the frontend and callable from JavaScript.
func (a *App) GetDeviceInfo() (*services.DeviceInfo, error) {
	return a.service.GetDeviceInfo()
}

// GetPort returns the HTTP server port
func (a *App) GetPort() string {
	return httpPort
}
