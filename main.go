package main

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"

	"getdevice-api/services"
	"getdevice-api/utils"
)

//go:embed frontend/src build/appicon.png
var assets embed.FS

func main() {
	envPath := utils.GetEnvPath()
	log.Printf("Using config file: %s", envPath)
	utils.MigrateEnvIfNeeded(envPath)
	if err := godotenv.Load(envPath); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	port := utils.GetEnv("PORT", "8080")
	password := utils.GetEnv("PASSWORD", "12345678")
	deviceService := services.NewDeviceService(envPath)
	app := application.New(application.Options{
		Name:        "Pegasus",
		Description: "Local device-information API server",
		Icon:        appIcon(),
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(mustSubFS(assets, "frontend/src")),
		},
		Windows: application.WindowsOptions{DisableQuitOnLastWindowClosed: true},
		Mac:     application.MacOptions{ActivationPolicy: application.ActivationPolicyAccessory},
	})

	runInBackground := !strings.EqualFold(utils.GetEnv("RUN_IN_BACKGROUND", "true"), "false")
	api := NewApp(deviceService, port, password, envPath, app.Autostart, runInBackground)
	app.RegisterService(application.NewService(api))

	backgroundLaunch := len(os.Args) > 1 && os.Args[1] == "--background"
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Pegasus",
		Width:            480,
		Height:           620,
		Hidden:           backgroundLaunch,
		BackgroundColour: application.NewRGBA(15, 17, 23, 255),
	})
	quitting := false
	window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		if quitting {
			return
		}
		if api.runInBackgroundEnabled() {
			window.Hide()
			event.Cancel()
			return
		}
		quitting = true
		event.Cancel()
		app.Quit()
	})

	tray := app.SystemTray.New()
	tray.SetIcon(appIcon())
	tray.SetLabel("Pegasus")
	tray.SetTooltip("Pegasus — local device API")
	menu := app.NewMenu()
	menu.Add("Open Pegasus").OnClick(func(*application.Context) {
		window.Show()
		window.Focus()
	})
	menu.AddSeparator()
	menu.Add("Quit").OnClick(func(*application.Context) {
		quitting = true
		app.Quit()
	})
	tray.SetMenu(menu)
	tray.OnClick(func() {
		window.Show()
		window.Focus()
	})
	if err := app.Run(); err != nil {
		log.Fatal("Error starting Pegasus:", err)
	}
}

func appIcon() []byte {
	icon, err := assets.ReadFile("build/appicon.png")
	if err != nil {
		log.Fatal(err)
	}
	return icon
}

func mustSubFS(fsys fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		log.Fatal(err)
	}
	return sub
}
