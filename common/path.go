package common

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetPathFromPath(fileName string) (string, error) {

	homeDir, err := os.UserHomeDir()
	path := os.Getenv("TEXT_CLIPPER_PATH")
	if path == "" {
		path = homeDir
	}

	if err != nil {
		return "", fmt.Errorf("unable to find home directory: %w", err)
	}

	defaultPath := filepath.Join(path, ".text-clipper", fileName)

	return defaultPath, nil
}
