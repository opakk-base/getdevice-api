package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"getdevice-api/handlers"
	"getdevice-api/middleware"
	"getdevice-api/services"
)

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

	// Create handlers
	deviceHandler := handlers.NewDeviceHandler(deviceService)

	// Set up routes
	mux := http.NewServeMux()

	mux.HandleFunc("/health", deviceHandler.HealthCheck)
	mux.HandleFunc("/getdevice", deviceHandler.GetDevice)

	// Start server
	addr := ":" + port
	fmt.Printf("Starting server on port %s...\n", port)
	fmt.Printf("Endpoints:\n")
	fmt.Printf("  - GET /health\n")
	fmt.Printf("  - GET /getdevice\n")
	fmt.Printf("\nPress Ctrl+C to stop\n")

	if err := http.ListenAndServe(addr, middleware.CORS(mux)); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
