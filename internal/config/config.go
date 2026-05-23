package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	appDirName     = "kp"
	promptsDirName = "prompts"
	dirMode        = 0o700
)

type Options struct {
	ConfigDir string
}

type Paths struct {
	ConfigRoot string
	PromptsDir string
}

func Resolve(opts Options) (Paths, error) {
	root, err := configRoot(opts.ConfigDir)
	if err != nil {
		return Paths{}, err
	}

	return Paths{
		ConfigRoot: root,
		PromptsDir: filepath.Join(root, promptsDirName),
	}, nil
}

func Ensure(opts Options) (Paths, error) {
	paths, err := Resolve(opts)
	if err != nil {
		return Paths{}, err
	}

	if err := os.MkdirAll(paths.ConfigRoot, dirMode); err != nil {
		return Paths{}, fmt.Errorf("create config dir %q: %w", paths.ConfigRoot, err)
	}
	if err := os.MkdirAll(paths.PromptsDir, dirMode); err != nil {
		return Paths{}, fmt.Errorf("create prompts dir %q: %w", paths.PromptsDir, err)
	}

	return paths, nil
}

func configRoot(override string) (string, error) {
	if override != "" {
		return absolute(override)
	}

	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return absolute(filepath.Join(xdg, appDirName))
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	if home == "" {
		return "", fmt.Errorf("resolve home dir: empty home")
	}

	return absolute(filepath.Join(home, ".config", appDirName))
}

func absolute(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve absolute path %q: %w", path, err)
	}
	return abs, nil
}
