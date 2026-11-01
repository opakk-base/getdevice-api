package main

import (
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestApp_GetConfigPath(t *testing.T) {
	dummyPath := filepath.Join("testdir", ".env")
	app := NewApp(nil, "8080", "secret123", dummyPath, nil, true)

	if got := app.GetConfigPath(); got != dummyPath {
		t.Fatalf("expected config path %q, got %q", dummyPath, got)
	}
}

func TestApp_VerifyPassword(t *testing.T) {
	app := NewApp(nil, "8080", "mypassword", ".env", nil, true)

	if !app.VerifyPassword("mypassword") {
		t.Errorf("expected VerifyPassword to succeed for valid password")
	}
	if app.VerifyPassword("wrongpassword") {
		t.Errorf("expected VerifyPassword to fail for incorrect password")
	}
}

func TestApp_WailsServiceRegistration(t *testing.T) {
	app := NewApp(nil, "8080", "12345678", ".env", nil, true)
	service := application.NewService(app)
	if service.Instance() != app {
		t.Fatal("expected service instance to match app")
	}
}
