package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode"

	"github.com/jamesonstone/kp/internal/picker"
)

func TestRootLauncherListsCurrentActions(t *testing.T) {
	fake := &fakeClipboard{}
	var got []picker.Item
	stdout, stderr, err := executeTestCommand(t,
		withClipboard(fake),
		withPicker(t, func(items []picker.Item) (string, error) {
			got = items
			return "prompt:review", nil
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	wantIDs := []string{"prompt:merge", "prompt:review", "prompt:ship", "command:init", "command:find-port", "command:help"}
	if strings.Join(itemIDs(got), ",") != strings.Join(wantIDs, ",") {
		t.Fatalf("launcher items = %v, want %v", itemIDs(got), wantIDs)
	}

	review := got[1]
	if review.Title != "Independent branch/PR review" || review.Command != "kp review" || review.Group != "prompts" {
		t.Fatalf("review item = %+v", review)
	}
	if got[3].Title != "Init" || got[3].Command != "kp init" || got[3].Group != "commands" {
		t.Fatalf("init item = %+v", got[3])
	}

	if stdout == "" || stdout != review.Preview {
		t.Fatalf("stdout = %q, want the review preview body", stdout)
	}
	if !strings.Contains(stderr, "Prompt \"review\" copied to clipboard.") {
		t.Fatalf("stderr = %q", stderr)
	}
	if fake.copied != stdout || fake.verified != stdout {
		t.Fatalf("clipboard copied=%q verified=%q stdout=%q", fake.copied, fake.verified, stdout)
	}
}

func TestLauncherItemsHaveNoEmoji(t *testing.T) {
	var got []picker.Item
	_, _, _ = executeTestCommand(t, withPicker(t, func(items []picker.Item) (string, error) {
		got = items
		return "", errPickerCanceled
	}))
	for _, item := range got {
		for _, text := range []string{item.Title, item.Command, item.Note} {
			for _, r := range text {
				if r > unicode.MaxLatin1 && r != '·' {
					t.Fatalf("item %q contains non-text glyph %q in %q", item.ID, r, text)
				}
			}
		}
	}
}

func TestLauncherMarksUserPromptsAsSecondaryDetail(t *testing.T) {
	configDir := t.TempDir()
	writeUserPrompt(t, configDir, "standup", "---\nlabel: Daily standup\n---\nSummarize.\n")

	var got []picker.Item
	_, _, err := executeTestCommandWithConfig(t, configDir,
		withPicker(t, func(items []picker.Item) (string, error) {
			got = items
			return "", errPickerCanceled
		}),
	)
	if ExitCode(err) != ExitCancel {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	for _, item := range got {
		if item.ID == "prompt:standup" {
			if item.Title != "Daily standup" || item.Command != "kp standup" || item.Note != "user prompt" {
				t.Fatalf("user item = %+v", item)
			}
			return
		}
	}
	t.Fatalf("user prompt missing from launcher: %v", itemIDs(got))
}

func TestRootLauncherRunsInit(t *testing.T) {
	fake := &fakeClipboard{}
	input := strings.Join([]string{
		"Ship the init command.",
		"kp is a local macOS prompt CLI.",
		"Clipboard verification stays exact.",
		"Do not add network calls.",
		"kp init prints and copies the prompt.",
	}, "\n") + "\n"
	stdout, _, err := executeTestCommand(t,
		withClipboard(fake),
		withStdin(input),
		withPicker(t, selectItem("command:init")),
	)
	if err != nil {
		t.Fatal(err)
	}
	want := renderInitPrompt(initAnswers{
		objective:   "Ship the init command.",
		context:     "kp is a local macOS prompt CLI.",
		invariants:  "Clipboard verification stays exact.",
		constraints: "Do not add network calls.",
		done:        "kp init prints and copies the prompt.",
	})
	if stdout != want || fake.copied != stdout {
		t.Fatalf("stdout = %q, clipboard = %q, want %q", stdout, fake.copied, want)
	}
}

func TestRootLauncherShowsStaticHelp(t *testing.T) {
	stdout, _, err := executeTestCommand(t, withPicker(t, selectItem("command:help")))
	if err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"Usage", "Prompt Commands", "kp new <name>", "kp v0 <command>", "kp --version"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout missing %q:\n%s", text, stdout)
		}
	}
}

func TestRootLauncherCancelUsesFarewell(t *testing.T) {
	fake := &fakeClipboard{}
	stdout, stderr, err := executeTestCommand(t,
		withClipboard(fake),
		withPicker(t, func([]picker.Item) (string, error) { return "", errPickerCanceled }),
	)
	if got := ExitCode(err); got != ExitCancel {
		t.Fatalf("ExitCode = %d, want %d", got, ExitCancel)
	}
	if stdout != "" || stderr != "" || fake.copied != "" {
		t.Fatalf("stdout=%q stderr=%q copied=%q, want no side effects", stdout, stderr, fake.copied)
	}
	assertPickerFarewell(t, ExitMessage(err))
}

func TestRootLauncherRejectsUnknownSelection(t *testing.T) {
	fake := &fakeClipboard{}
	_, _, err := executeTestCommand(t, withClipboard(fake), withPicker(t, selectItem("prompt:clarify")))
	if ExitCode(err) != ExitUser || !strings.Contains(err.Error(), "invalid picker selection") {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if fake.copied != "" {
		t.Fatalf("clipboard copied %q", fake.copied)
	}
}

func TestPickerWithoutTerminalExitsConfig(t *testing.T) {
	fake := &fakeClipboard{}
	_, _, err := executeTestCommand(t,
		withClipboard(fake),
		withPicker(t, func([]picker.Item) (string, error) {
			return "", fmt.Errorf("%w: open /dev/tty: device not configured", picker.ErrNoTerminal)
		}),
	)
	if ExitCode(err) != ExitConfig {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if !strings.Contains(err.Error(), "kp --help") || !strings.Contains(err.Error(), "--no-fzf") {
		t.Fatalf("err = %v", err)
	}
	if fake.copied != "" {
		t.Fatalf("clipboard copied %q", fake.copied)
	}
}

func TestPickerOperationalFailureExitsUser(t *testing.T) {
	_, _, err := executeTestCommand(t, withPicker(t, func([]picker.Item) (string, error) {
		return "", errors.New("render failed")
	}))
	if ExitCode(err) != ExitUser || ExitMessage(err) != "render failed" {
		t.Fatalf("ExitCode = %d, message = %q", ExitCode(err), ExitMessage(err))
	}
}

func itemIDs(items []picker.Item) []string {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}
	return ids
}

func writeUserPrompt(t *testing.T, configDir, name, content string) {
	t.Helper()
	dir := filepath.Join(configDir, "prompts")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name+".md"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
