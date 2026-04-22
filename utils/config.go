package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const appName = "GetDevice"

// GetEnvPath returns the full path to the .env file in the user's
// config directory (e.g. ~/Library/Application Support/GetDevice/.env on macOS,
// %AppData%/GetDevice/.env on Windows). It creates the directory if needed.
//
// If the config directory cannot be determined, it falls back to ".env"
// in the current working directory (development mode).
func GetEnvPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Printf("Warning: cannot determine user config dir: %v, falling back to .env in CWD", err)
		return ".env"
	}

	appDir := filepath.Join(configDir, appName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		log.Printf("Warning: cannot create config dir %s: %v, falling back to .env in CWD", appDir, err)
		return ".env"
	}

	return filepath.Join(appDir, ".env")
}

// MigrateEnvIfNeeded copies an existing .env file from legacy locations
// to the new config directory path, preserving existing DEVICE_ID and
// CLIENT_KEY values. It only migrates if the destination does not already exist.
//
// Legacy locations checked (in order):
//  1. ".env" in the current working directory
//  2. ".env" next to the running executable
func MigrateEnvIfNeeded(newEnvPath string) {
	// If the new path already exists, nothing to migrate
	if _, err := os.Stat(newEnvPath); err == nil {
		return
	}

	// Candidate legacy paths
	candidates := []string{}

	// 1. CWD/.env
	candidates = append(candidates, ".env")

	// 2. Next to the executable
	if exePath, err := os.Executable(); err == nil {
		if resolved, err := filepath.EvalSymlinks(exePath); err == nil {
			candidates = append(candidates, filepath.Join(filepath.Dir(resolved), ".env"))
		}
	}

	for _, candidate := range candidates {
		abs, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}

		// Don't migrate from the same path
		newAbs, _ := filepath.Abs(newEnvPath)
		if abs == newAbs {
			continue
		}

		if _, err := os.Stat(abs); err == nil {
			if err := copyFile(abs, newEnvPath); err != nil {
				log.Printf("Warning: failed to migrate .env from %s: %v", abs, err)
			} else {
				log.Printf("Migrated .env from %s to %s", abs, newEnvPath)
			}
			return
		}
	}
}

// copyFile copies a file from src to dst, preserving permissions.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("stat source: %w", err)
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("create destination: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("copy data: %w", err)
	}

	return nil
}
