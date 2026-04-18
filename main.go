package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"getdevice-api/handlers"
	"getdevice-api/middleware"
	"getdevice-api/services"
)

//go:embed all:frontend/src
var assets embed.FS

// httpPort is the port the background HTTP server listens on
var httpPort string

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Get port from environment or use default
	httpPort = os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	// Initialize services
	envPath := ".env"
	deviceService := services.NewDeviceService(envPath)

	// Start HTTP API server in background goroutine
	go startHTTPServer(deviceService)

	// Create Wails app
	app := NewApp(deviceService)

	err := wails.Run(&options.App{
		Title:  "GetDevice",
		Width:  520,
		Height: 680,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatal("Error starting Wails app:", err)
	}
}

// startHTTPServer runs the existing REST API in the background
func startHTTPServer(deviceService *services.DeviceService) {
	deviceHandler := handlers.NewDeviceHandler(deviceService)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", deviceHandler.HealthCheck)
	mux.HandleFunc("/getdevice", deviceHandler.GetDevice)

	addr := ":" + httpPort
	fmt.Printf("HTTP API server running on port %s\n", httpPort)
	fmt.Printf("Endpoints:\n")
	fmt.Printf("  - GET http://localhost:%s/health\n", httpPort)
	fmt.Printf("  - GET http://localhost:%s/getdevice\n", httpPort)

	if err := http.ListenAndServe(addr, middleware.CORS(mux)); err != nil {
		log.Printf("HTTP server error: %v", err)
	}
}
