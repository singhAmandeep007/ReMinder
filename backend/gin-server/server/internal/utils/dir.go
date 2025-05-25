package utils

import (
	"os"
	"path/filepath"
)

func ResolvePathFromProjectRoot(envFile string) string {
	// First try to find go.mod (development mode)
	if projectRoot := findProjectRoot(); projectRoot != "" {
		return filepath.Join(projectRoot, envFile)
	}

	// Fallback to executable directory (production mode)
	execPath, err := os.Executable()
	if err != nil {
		// Last resort: current working directory
		return envFile
	}

	execDir := filepath.Dir(execPath)
	return filepath.Join(execDir, envFile)
}

func findProjectRoot() string {
	currentDir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return currentDir
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			return "" // Not found
		}
		currentDir = parent
	}
}
