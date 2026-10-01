package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func configPath(extensions ...string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolving home directory: %w", err)
	}
	parts := append([]string{home, ".gatewright"}, extensions...)
	return filepath.Join(parts...), nil
}
