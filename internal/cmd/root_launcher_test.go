package cmd

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesonstone/kp/internal/clipboard"
	"github.com/jamesonstone/kp/internal/prompt"
)

func TestRootLaunchesInteractiveSelectorForPrompt(t *testing.T) {
	fake := &fakeClipboard{}
	var sawClarify bool
	var sawFindPort bool
	var sawHelp bool
	stdout, stderr, err := executeTestCommand(t,
		withClipboard(fake),
		withLauncher(func(items []LauncherItem) (string, error) {
			for _, item := range items {
				switch item.ID {
				case "prompt:clarify":
					sawClarify = item.Emoji == "🧠" &&
						item.Command == "kp clarify" &&
						item.Description == "Print and copy prompt"
				case "command:find-port":
					sawFindPort = item.Emoji == "🔍" &&
						item.Command == "kp find-port <port>" &&
						item.Description == "Inspect a port and act on the process"
				case "command:help":
					sawHelp = item.Emoji == "❓" &&
						item.Command == "kp --help" &&
						item.Description == "Show all commands"
				}
			}
			return "prompt:clarify", nil
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !sawClarify || !sawFindPort || !sawHelp {
		t.Fatalf("launcher items missing expected entries: sawClarify=%v sawFindPort=%v sawHelp=%v", sawClarify, sawFindPort, sawHelp)
	}
	if !strings.HasPrefix(stdout, "Clarify before implementing.") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stderr, "✅ 📋 Prompt \"clarify\" copied to clipboard.") {
		t.Fatalf("stderr = %q", stderr)
	}
	if fake.copied == "" || fake.verified == "" {
		t.Fatalf("clipboard copied=%q verified=%q", fake.copied, fake.verified)
	}
}

func TestRootLauncherShowsStaticHelp(t *testing.T) {
	secondaryCommands := map[string]bool{
		"command:list":     false,
		"command:new":      false,
		"command:edit":     false,
		"command:rm":       false,
		"command:scaffold": false,
		"command:version":  false,
	}
	stdout, _, err := executeTestCommand(t,
		withLauncher(func(items []LauncherItem) (string, error) {
			for _, item := range items {
				if _, hidden := secondaryCommands[item.ID]; hidden {
					secondaryCommands[item.ID] = true
				}
			}
			return "command:help", nil
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	for id, visible := range secondaryCommands {
		if visible {
			t.Fatalf("secondary command %q should be hidden from the launcher", id)
		}
	}
	for _, text := range []string{"Usage", "Prompt Commands", "kp new <name>", "kp edit <name>", "kp rm <name>", "kp scaffold", "kp --version"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout missing %q:\n%s", text, stdout)
		}
	}
}

func TestLauncherDisplayRowsAlignColumns(t *testing.T) {
	items := []LauncherItem{
		{
			ID:          "prompt:clarify",
			Emoji:       "🧠",
			Title:       "Clarify before implementing",
			Command:     "kp clarify",
			Description: "Print and copy prompt",
		},
		{
			ID:          "command:new",
			Emoji:       "✨",
			Title:       "New prompt",
			Command:     "kp new <name>",
			Description: "Show creation help",
		},
		{
			ID:          "command:version",
			Emoji:       "🛠️",
			Title:       "Version",
			Command:     "kp --version",
			Description: "Print version metadata",
		},
	}

	rows := launcherDisplayRows(items)
	clarify := rows["prompt:clarify"]
	newPrompt := rows["command:new"]
	version := rows["command:version"]

	commandColumn := displayColumn(t, clarify, "kp clarify")
	rowsToAlign := map[string]string{
		"new":     newPrompt,
		"version": version,
	}
	for name, row := range rowsToAlign {
		if got := displayColumn(t, row, "kp "); got != commandColumn {
			t.Fatalf("%s command column = %d, want %d\nclarify: %q\nrow: %q", name, got, commandColumn, clarify, row)
		}
	}
	for id, row := range rows {
		for _, item := range items {
			if strings.Contains(row, item.Description) {
				t.Fatalf("%s row includes description %q: %q", id, item.Description, row)
			}
		}
	}
}

func TestLauncherFZFArgsUseConciseWrappingLayoutAndVimNavigation(t *testing.T) {
	args := launcherFZFArgs("/tmp/kp preview")
	joined := strings.Join(args, " ")

	for _, expected := range []string{
		"--bind j:down,k:up,tab:down,btab:up",
		"--header j/k or arrows: move · enter: select · esc: close",
		"--info hidden",
		"--no-separator",
		"--no-hscroll",
		"--preview-window right,55%,wrap",
	} {
		if !strings.Contains(joined, expected) {
			t.Fatalf("args missing %q: %q", expected, joined)
		}
	}
}

func TestLauncherFZFDoesNotStartWhenContextIsCanceled(t *testing.T) {
	tempDir := t.TempDir()
	markerPath := filepath.Join(tempDir, "started")
	fzfPath := filepath.Join(tempDir, "fzf")
	script := "#!/bin/sh\nprintf started > \"$KP_MARKER\"\nprintf 'prompt:test\\n'\n"
	if err := os.WriteFile(fzfPath, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("KP_MARKER", markerPath)
	t.Setenv("PATH", tempDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := runLauncherFZF(ctx, []LauncherItem{{ID: "prompt:test", Title: "Test", Preview: "test"}})
	if err == nil {
		t.Fatal("runLauncherFZF returned nil error for canceled context")
	}
	if _, statErr := os.Stat(markerPath); !os.IsNotExist(statErr) {
		t.Fatalf("fzf started after cancellation: %v", statErr)
	}
}

func displayColumn(t *testing.T, row string, needle string) int {
	t.Helper()

	index := strings.Index(row, needle)
	if index < 0 {
		t.Fatalf("row %q missing %q", row, needle)
	}
	return displayWidth(row[:index])
}

func TestRootLauncherFZFMissingDoesNotLoadRegistryOrClipboard(t *testing.T) {
	var registryCalled bool
	var clipboardCalled bool
	_, _, err := executeTestCommand(t, func(opts *Options) {
		opts.RegistryFactory = func(string) (prompt.Registry, error) {
			registryCalled = true
			return nil, nil
		}
		opts.ClipboardFactory = func() clipboard.Clipboard {
			clipboardCalled = true
			return nil
		}
		opts.LookPath = func(string) (string, error) {
			return "", errors.New("not found")
		}
	})
	if ExitCode(err) != ExitConfig {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if !strings.Contains(err.Error(), "brew install fzf") || !strings.Contains(err.Error(), "kp --help") {
		t.Fatalf("err = %v", err)
	}
	if registryCalled || clipboardCalled {
		t.Fatalf("registryCalled=%v clipboardCalled=%v", registryCalled, clipboardCalled)
	}
}
