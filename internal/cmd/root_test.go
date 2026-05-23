package cmd

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jamesonstone/kp/internal/clipboard"
	"github.com/jamesonstone/kp/internal/prompt"
)

func TestListPlain(t *testing.T) {
	stdout, _, err := executeTestCommand(t, "list", "--plain")
	if err != nil {
		t.Fatal(err)
	}

	if stdout != "clarify\ninstructions\nparentthread\n" {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestListVerbose(t *testing.T) {
	stdout, _, err := executeTestCommand(t, "list", "--verbose")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(stdout, "clarify\tClarify before implementing\tbuiltin\n") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stdout, "instructions\tCoding agent instructions\tbuiltin\n") {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestRootShowsHelpWithoutSideEffects(t *testing.T) {
	var registryCalled bool
	var clipboardCalled bool
	stdout, _, err := executeTestCommand(t, func(opts *Options) {
		opts.RegistryFactory = func(string) (prompt.Registry, error) {
			registryCalled = true
			return nil, nil
		}
		opts.ClipboardFactory = func() clipboard.Clipboard {
			clipboardCalled = true
			return nil
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{
		"Low-friction prompt utilities",
		"Usage",
		"kp <prompt>",
		"Prompt Commands",
		"kp clarify",
		"kp instructions",
		"Prompt Library",
		"kp list",
		"kp list --plain",
		"Repo Setup",
		"kp scaffold",
		"Utilities",
	}
	for _, text := range expected {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout missing %q:\n%s", text, stdout)
		}
	}
	if strings.Index(stdout, "Prompt Commands") > strings.Index(stdout, "Prompt Library") {
		t.Fatalf("prompt command section appears after library section:\n%s", stdout)
	}
	if strings.Index(stdout, "kp clarify") > strings.Index(stdout, "kp list") {
		t.Fatalf("direct prompt commands appear after list commands:\n%s", stdout)
	}
	if strings.Contains(stdout, "kp"+" prompt") {
		t.Fatalf("stdout = %q", stdout)
	}
	if registryCalled || clipboardCalled {
		t.Fatalf("registryCalled=%v clipboardCalled=%v", registryCalled, clipboardCalled)
	}
}

func TestRootHelpUsesKitStyleWhenTerminal(t *testing.T) {
	previous := terminalWriterCheck
	terminalWriterCheck = func(io.Writer) bool { return true }
	t.Cleanup(func() {
		terminalWriterCheck = previous
	})

	stdout, _, err := executeTestCommand(t)
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"\x1b[1;37m🚀 Usage\x1b[0m",
		"\x1b[1;37m🧠 Prompt Commands\x1b[0m",
		"\x1b[1;37m🧰 Prompt Library\x1b[0m",
		"\x1b[1;37m🏗️ Repo Setup\x1b[0m",
		"\x1b[1;37m🛠️ Utilities\x1b[0m",
		"\x1b[1;37m⚙️ Flags\x1b[0m",
	}
	for _, text := range expected {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout missing %q:\n%s", text, stdout)
		}
	}
}

func TestPromptPrint(t *testing.T) {
	stdout, stderr, err := executeTestCommand(t, "clarify", "--print")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(stdout, "Clarify before implementing.") {
		t.Fatalf("stdout = %q", stdout)
	}
	if strings.Contains(stdout, "---") {
		t.Fatalf("stdout includes frontmatter: %q", stdout)
	}
	if strings.Contains(stderr, "copied to clipboard") {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestPromptCopy(t *testing.T) {
	fake := &fakeClipboard{}
	stdout, stderr, err := executeTestCommand(t, "clarify", "--copy", withClipboard(fake))
	if err != nil {
		t.Fatal(err)
	}

	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if fake.copied == "" {
		t.Fatal("Copy was not called")
	}
	if fake.verified == "" {
		t.Fatal("Verify was not called")
	}
	if fake.pasted {
		t.Fatal("Paste was called for --copy")
	}
	if !strings.Contains(stderr, "✅ 📋 Prompt \"clarify\" copied to clipboard.") {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestPromptDefaultShowsClipboardInstructionsWithSpacing(t *testing.T) {
	fake := &fakeClipboard{}
	stdout, stderr, err := executeTestCommand(t, "clarify", withClipboard(fake))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(stdout, "Clarify before implementing.") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stderr, "✅ 📋 Prompt \"clarify\" copied to clipboard.") {
		t.Fatalf("stderr = %q", stderr)
	}
	if !strings.Contains(stderr, "🧾 Full prompt content is printed to stdout below.\n\n") {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestPromptDefaultPrintsAndCopiesWithoutPaste(t *testing.T) {
	fake := &fakeClipboard{}
	stdout, _, err := executeTestCommand(t, "clarify", withClipboard(fake))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(stdout, "Clarify before implementing.") {
		t.Fatalf("stdout = %q", stdout)
	}
	if fake.copied == "" {
		t.Fatal("Copy was not called")
	}
	if fake.verified == "" {
		t.Fatal("Verify was not called")
	}
	if fake.pasted {
		t.Fatal("Paste was called")
	}
}

func TestPromptCopyVerifyFailureExitsSystem(t *testing.T) {
	fake := &fakeClipboard{verifyErr: clipboard.ErrVerifyFailed}
	_, _, err := executeTestCommand(t, "clarify", "--copy", withClipboard(fake))
	if ExitCode(err) != ExitSystem {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if fake.pasted {
		t.Fatal("Paste was called after verification failure")
	}
}

func TestPromptDefaultVerifyFailureExitsSystem(t *testing.T) {
	fake := &fakeClipboard{verifyErr: clipboard.ErrVerifyFailed}
	_, _, err := executeTestCommand(t, "clarify", withClipboard(fake))
	if ExitCode(err) != ExitSystem {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if fake.pasted {
		t.Fatal("Paste was called after verification failure")
	}
}

func TestPromptVerboseLogsToStderr(t *testing.T) {
	fake := &fakeClipboard{}
	_, stderr, err := executeTestCommand(t, "clarify", "--copy", "--verbose", withClipboard(fake))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(stderr, "event=copy name=clarify bytes=") {
		t.Fatalf("stderr = %q", stderr)
	}
}

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

func TestPickerFZFSelection(t *testing.T) {
	fake := &fakeClipboard{}
	_, _, err := executeTestCommand(t,
		"list",
		withClipboard(fake),
		withFZF("clarify", nil),
	)
	if err != nil {
		t.Fatal(err)
	}
	if fake.copied == "" {
		t.Fatal("Copy was not called after fzf selection")
	}
	if fake.pasted {
		t.Fatal("Paste was called after fzf selection")
	}
}

func TestPickerFZFCancel(t *testing.T) {
	fake := &fakeClipboard{}
	_, _, err := executeTestCommand(t,
		"list",
		withClipboard(fake),
		withFZF("", errPickerCanceled),
	)
	if ExitCode(err) != ExitCancel {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if fake.copied != "" || fake.pasted {
		t.Fatalf("clipboard side effects copied=%q pasted=%v", fake.copied, fake.pasted)
	}
}

func TestPickerFZFMissing(t *testing.T) {
	fake := &fakeClipboard{}
	_, _, err := executeTestCommand(t,
		"list",
		withClipboard(fake),
		func(opts *Options) {
			opts.LookPath = func(string) (string, error) {
				return "", errors.New("not found")
			}
		},
	)
	if ExitCode(err) != ExitConfig {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if !strings.Contains(err.Error(), "brew install fzf") || !strings.Contains(err.Error(), "--no-fzf") {
		t.Fatalf("err = %v", err)
	}
	if fake.copied != "" || fake.pasted {
		t.Fatalf("clipboard side effects copied=%q pasted=%v", fake.copied, fake.pasted)
	}
}

func TestPickerNoFZFValidSelection(t *testing.T) {
	fake := &fakeClipboard{}
	_, stderr, err := executeTestCommand(t,
		"list",
		"--no-fzf",
		withStdin("1\n"),
		withClipboard(fake),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stderr, "1\tclarify\tClarify before implementing\n") {
		t.Fatalf("stderr = %q", stderr)
	}
	if fake.copied == "" {
		t.Fatal("Copy was not called after numbered selection")
	}
	if fake.pasted {
		t.Fatal("Paste was called after numbered selection")
	}
}

func TestPickerNoFZFInvalidSelection(t *testing.T) {
	fake := &fakeClipboard{}
	_, _, err := executeTestCommand(t,
		"list",
		"--no-fzf",
		withStdin("abc\n"),
		withClipboard(fake),
	)
	if ExitCode(err) != ExitUser {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if fake.copied != "" || fake.pasted {
		t.Fatalf("clipboard side effects copied=%q pasted=%v", fake.copied, fake.pasted)
	}
}

func TestPickerNoFZFOutOfRange(t *testing.T) {
	fake := &fakeClipboard{}
	_, _, err := executeTestCommand(t,
		"list",
		"--no-fzf",
		withStdin("99\n"),
		withClipboard(fake),
	)
	if ExitCode(err) != ExitUser {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if fake.copied != "" || fake.pasted {
		t.Fatalf("clipboard side effects copied=%q pasted=%v", fake.copied, fake.pasted)
	}
}

func TestPickerNoFZFCancel(t *testing.T) {
	fake := &fakeClipboard{}
	_, _, err := executeTestCommand(t,
		"list",
		"--no-fzf",
		withStdin(""),
		withClipboard(fake),
	)
	if ExitCode(err) != ExitCancel {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
	if fake.copied != "" || fake.pasted {
		t.Fatalf("clipboard side effects copied=%q pasted=%v", fake.copied, fake.pasted)
	}
}

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
		"new", "instructions",
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
	_, _, err := executeTestCommand(t, "rm", "instructions")
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

func TestHelpDoesNotLoadRegistryOrClipboard(t *testing.T) {
	var registryCalled bool
	var clipboardCalled bool
	_, _, err := executeTestCommand(t, "--help", func(opts *Options) {
		opts.RegistryFactory = func(string) (prompt.Registry, error) {
			registryCalled = true
			return nil, nil
		}
		opts.ClipboardFactory = func() clipboard.Clipboard {
			clipboardCalled = true
			return nil
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if registryCalled || clipboardCalled {
		t.Fatalf("registryCalled=%v clipboardCalled=%v", registryCalled, clipboardCalled)
	}
}

func executeTestCommand(t *testing.T, args ...any) (string, string, error) {
	t.Helper()
	return executeTestCommandWithConfig(t, t.TempDir(), args...)
}

func executeTestCommandWithConfig(t *testing.T, configDir string, args ...any) (string, string, error) {
	t.Helper()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	opts := Options{
		Version: "test",
		Commit:  "abc123",
		Stdout:  &stdout,
		Stderr:  &stderr,
	}

	cmdArgs := make([]string, 0, len(args)+2)
	cmdArgs = append(cmdArgs, "--config", configDir)
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			cmdArgs = append(cmdArgs, v)
		case func(*Options):
			v(&opts)
		default:
			t.Fatalf("unsupported test arg %T", arg)
		}
	}

	cmd := NewRoot(opts)
	cmd.SetArgs(cmdArgs)
	err := cmd.Execute()
	return stdout.String(), stderr.String(), err
}

func withClipboard(fake *fakeClipboard) func(*Options) {
	return func(opts *Options) {
		opts.ClipboardFactory = func() clipboard.Clipboard {
			return fake
		}
	}
}

func withStdin(input string) func(*Options) {
	return func(opts *Options) {
		opts.Stdin = strings.NewReader(input)
	}
}

func withFZF(selection string, runErr error) func(*Options) {
	return func(opts *Options) {
		opts.LookPath = func(string) (string, error) {
			return "/usr/local/bin/fzf", nil
		}
		opts.FZFRunner = func(prompts []prompt.Prompt) (string, error) {
			if len(prompts) < 2 {
				return "", errors.New("unexpected prompt count")
			}
			for _, p := range prompts {
				if strings.Contains(p.Body, "---") {
					return "", errors.New("frontmatter leaked into picker body")
				}
			}
			return selection, runErr
		}
	}
}

func withEditor(run func(name string, args []string, path string) error) func(*Options) {
	return func(opts *Options) {
		opts.Getenv = func(key string) string {
			if key == "KP_EDITOR" {
				return "test-editor"
			}
			return ""
		}
		opts.LookPath = func(name string) (string, error) {
			if name == "test-editor" {
				return "/bin/test-editor", nil
			}
			return "", errors.New("unexpected lookup")
		}
		opts.EditorRunner = run
	}
}

type fakeClipboard struct {
	copied    string
	verified  string
	pasted    bool
	verifyErr error
}

func (f *fakeClipboard) Copy(body string) error {
	f.copied = body
	return nil
}

func (f *fakeClipboard) Read() (string, error) {
	return f.copied, nil
}

func (f *fakeClipboard) Verify(expected string, _ time.Duration) error {
	f.verified = expected
	if f.verifyErr != nil {
		return f.verifyErr
	}
	return nil
}
