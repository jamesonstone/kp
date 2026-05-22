package config

import (
	"errors"
	"os"
	"path/filepath"
)

func PromptsDir() (string, error) {
	xdg := os.Getenv("XDG_CONFIG_HOME")
	base := xdg
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", errors.New("cannot resolve home directory")
		}
		base = filepath.Join(home, ".config")
	}
	dir := filepath.Join(base, "kp", "prompts")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return dir, nil
}
