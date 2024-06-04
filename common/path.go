package common

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetPathFromDefaultPath(fileName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to find home directory: %w", err)
	}
	defaultPath := filepath.Join(homeDir, ".text-clipper", fileName)

	return defaultPath, nil
}
