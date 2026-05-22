package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPromptsDir_XDGOverride(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	dir, err := PromptsDir()
	if err != nil {
		t.Fatalf("PromptsDir error: %v", err)
	}
	if want := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "kp", "prompts"); dir != want {
		t.Fatalf("got %q want %q", dir, want)
	}
}

func TestPromptsDir_DefaultsToHome(t *testing.T) {
	h := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", h)
	dir, err := PromptsDir()
	if err != nil {
		t.Fatalf("PromptsDir error: %v", err)
	}
	if want := filepath.Join(h, ".config", "kp", "prompts"); dir != want {
		t.Fatalf("got %q want %q", dir, want)
	}
}
