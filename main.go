package main

import (
	"embed"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"getdevice-api/services"
)

//go:embed all:frontend/src
var assets embed.FS

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize services
	envPath := ".env"
	deviceService := services.NewDeviceService(envPath)

	// Create Wails app (server starts automatically in OnStartup)
	app := NewApp(deviceService, port)

	err := wails.Run(&options.App{
		Title:  "GetDevice",
		Width:  520,
		Height: 780,
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
