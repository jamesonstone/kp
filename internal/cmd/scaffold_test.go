package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesonstone/kp/internal/clipboard"
	"github.com/jamesonstone/kp/internal/prompt"
)

func TestScaffoldCommandCreatesApprovedFilesWithoutSideEffects(t *testing.T) {
	target := filepath.Join(t.TempDir(), "repo")
	var registryCalled bool
	var clipboardCalled bool
	stdout, _, err := executeTestCommand(t,
		"scaffold",
		"--dir", target,
		func(opts *Options) {
			opts.RegistryFactory = func(string) (prompt.Registry, error) {
				registryCalled = true
				return nil, nil
			}
			opts.ClipboardFactory = func() clipboard.Clipboard {
				clipboardCalled = true
				return nil
			}
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout, "create\tAGENTS.md\n") {
		t.Fatalf("stdout = %q", stdout)
	}
	for _, path := range []string{
		".env",
		".envrc",
		".coderabbit.yaml",
		".github/pull_request_template.md",
		".github/copilot-instructions.md",
		"AGENTS.md",
		"CLAUDE.md",
		"docs/agents/README.md",
		"docs/references/external-systems.md",
	} {
		if _, err := os.Stat(filepath.Join(target, filepath.FromSlash(path))); err != nil {
			t.Fatalf("expected scaffold file %s: %v", path, err)
		}
	}
	for _, path := range []string{
		".kit.yaml",
		".kit",
		"docs/CONSTITUTION.md",
		"docs/PROJECT_PROGRESS_SUMMARY.md",
		"docs/specs",
		"docs/notes",
	} {
		if _, err := os.Stat(filepath.Join(target, filepath.FromSlash(path))); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("excluded path %s exists or stat failed unexpectedly: %v", path, err)
		}
	}
	if registryCalled || clipboardCalled {
		t.Fatalf("registryCalled=%v clipboardCalled=%v", registryCalled, clipboardCalled)
	}
}

func TestScaffoldCommandDryRunDoesNotWrite(t *testing.T) {
	target := t.TempDir()
	stdout, _, err := executeTestCommand(t, "scaffold", "--dir", target, "--dry-run")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout, "would-create\tAGENTS.md\n") {
		t.Fatalf("stdout = %q", stdout)
	}
	if _, err := os.Stat(filepath.Join(target, "AGENTS.md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("AGENTS.md stat error = %v, want os.ErrNotExist", err)
	}
}

func TestScaffoldCommandForceOverwritesFiles(t *testing.T) {
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "AGENTS.md"), []byte("custom"), 0o644); err != nil {
		t.Fatal(err)
	}
	stdout, _, err := executeTestCommand(t, "scaffold", "--dir", target, "--force")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout, "update\tAGENTS.md\n") {
		t.Fatalf("stdout = %q", stdout)
	}
	got, err := os.ReadFile(filepath.Join(target, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) == "custom" || !strings.Contains(string(got), "# AGENTS") {
		t.Fatalf("AGENTS.md = %q", got)
	}
}
