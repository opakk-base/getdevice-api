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
	"getdevice-api/utils"
)

//go:embed all:frontend/src
var assets embed.FS

func main() {
	// Resolve config path in user's config directory
	// (e.g. ~/Library/Application Support/GetDevice/.env on macOS)
	envPath := utils.GetEnvPath()
	log.Printf("Using config file: %s", envPath)

	// Migrate .env from legacy locations (CWD or next to executable) if needed
	utils.MigrateEnvIfNeeded(envPath)

	// Load environment variables from the resolved path
	if err := godotenv.Load(envPath); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get close behavior from environment or use default
	closeBehavior := os.Getenv("CLOSE_BEHAVIOR")
	if closeBehavior != "minimize" {
		closeBehavior = "exit"
	}

	// Initialize services
	deviceService := services.NewDeviceService(envPath)

	// Create Wails app (server starts automatically in OnStartup)
	app := NewApp(deviceService, port, closeBehavior, envPath)

	err := wails.Run(&options.App{
		Title:            "GetDevice",
		Width:            420,
		Height:           220,
		BackgroundColour: &options.RGBA{R: 15, G: 17, B: 23, A: 255},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:    app.startup,
		OnBeforeClose: app.beforeClose,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatal("Error starting Wails app:", err)
	}
}
