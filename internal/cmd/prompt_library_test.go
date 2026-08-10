package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewCreatesPromptWithEditor(t *testing.T) {
	configDir := t.TempDir()
	_, _, err := executeTestCommandWithConfig(t, configDir,
		"new", "custom",
		withEditor(func(_ string, _ []string, path string) error {
			return os.WriteFile(path, []byte("body"), 0o600)
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(filepath.Join(configDir, "prompts", "custom.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "body" {
		t.Fatalf("prompt file = %q", got)
	}
}

func TestNewRejectsBuiltinCollision(t *testing.T) {
	_, _, err := executeTestCommand(t,
		"new", "handoff",
		withEditor(func(string, []string, string) error { return nil }),
	)
	if ExitCode(err) != ExitUser {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}

func TestNewRejectsInvalidName(t *testing.T) {
	_, _, err := executeTestCommand(t,
		"new", "Bad Name",
		withEditor(func(string, []string, string) error { return nil }),
	)
	if ExitCode(err) != ExitUser {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}

func TestNewRejectsReservedScaffoldName(t *testing.T) {
	_, _, err := executeTestCommand(t,
		"new", "scaffold",
		withEditor(func(string, []string, string) error { return nil }),
	)
	if ExitCode(err) != ExitUser {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}

func TestNewEmptyDeletesStub(t *testing.T) {
	configDir := t.TempDir()
	_, _, err := executeTestCommandWithConfig(t, configDir,
		"new", "custom",
		withEditor(func(string, []string, string) error { return nil }),
	)
	if ExitCode(err) != ExitUser {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if _, err := os.Stat(filepath.Join(configDir, "prompts", "custom.md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("stub stat error = %v, want os.ErrNotExist", err)
	}
}

func TestNewEditorCancelDeletesStub(t *testing.T) {
	configDir := t.TempDir()
	_, _, err := executeTestCommandWithConfig(t, configDir,
		"new", "custom",
		withEditor(func(string, []string, string) error { return errPickerCanceled }),
	)
	if ExitCode(err) != ExitCancel {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if _, err := os.Stat(filepath.Join(configDir, "prompts", "custom.md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("stub stat error = %v, want os.ErrNotExist", err)
	}
}

func TestEditPromotesBuiltin(t *testing.T) {
	configDir := t.TempDir()
	_, stderr, err := executeTestCommandWithConfig(t, configDir,
		"edit", "clarify",
		withEditor(func(string, []string, string) error { return nil }),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stderr, "promoted built-in \"clarify\"") {
		t.Fatalf("stderr = %q", stderr)
	}
	got, err := os.ReadFile(filepath.Join(configDir, "prompts", "clarify.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(got), "---\nlabel: Clarify before implementing") {
		t.Fatalf("promoted file = %q", got)
	}
}

func TestEditMissingPrompt(t *testing.T) {
	_, _, err := executeTestCommand(t,
		"edit", "missing",
		withEditor(func(string, []string, string) error { return nil }),
	)
	if ExitCode(err) != ExitUser {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}

func TestRMUserPrompt(t *testing.T) {
	configDir := t.TempDir()
	promptsDir := filepath.Join(configDir, "prompts")
	if err := os.MkdirAll(promptsDir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(promptsDir, "custom.md"), []byte("body"), 0o600); err != nil {
		t.Fatal(err)
	}

	_, stderr, err := executeTestCommandWithConfig(t, configDir, "rm", "custom")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stderr, "removed custom") {
		t.Fatalf("stderr = %q", stderr)
	}
	if _, err := os.Stat(filepath.Join(promptsDir, "custom.md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("prompt stat error = %v, want os.ErrNotExist", err)
	}
}

func TestRMBuiltinFails(t *testing.T) {
	_, _, err := executeTestCommand(t, "rm", "handoff")
	if ExitCode(err) != ExitUser {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}

func TestMissingEditorExitsConfig(t *testing.T) {
	_, _, err := executeTestCommand(t,
		"new", "custom",
		func(opts *Options) {
			opts.Getenv = func(key string) string {
				if key == "KP_EDITOR" {
					return "missing-editor"
				}
				return ""
			}
			opts.LookPath = func(name string) (string, error) {
				return "", errors.New("not found")
			}
		},
	)
	if ExitCode(err) != ExitConfig {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}

func TestConfigFailureExitsConfig(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "not-dir")
	if err := os.WriteFile(configPath, []byte("file"), 0o600); err != nil {
		t.Fatal(err)
	}

	_, _, err := executeTestCommandWithConfig(t, configPath, "list")
	if ExitCode(err) != ExitConfig {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}

func TestPromptNotFoundExitsUser(t *testing.T) {
	_, _, err := executeTestCommand(t, "nonexistent")
	if ExitCode(err) != ExitUser {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}

func TestPromptInvalidNameExitsUser(t *testing.T) {
	_, _, err := executeTestCommand(t, "Bad Name")
	if ExitCode(err) != ExitUser {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}

func TestPromptReservedNameExitsUser(t *testing.T) {
	_, _, err := executeTestCommand(t, "prompt")
	if ExitCode(err) != ExitUser {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}
