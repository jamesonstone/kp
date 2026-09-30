package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jamesonstone/kp/internal/clipboard"
	"github.com/jamesonstone/kp/internal/picker"
)

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

// withPicker replaces the interactive picker with run, failing the test if
// any preview still carries prompt frontmatter.
func withPicker(t *testing.T, run func(items []picker.Item) (string, error)) func(*Options) {
	return func(opts *Options) {
		opts.PickerRunner = func(items []picker.Item) (string, error) {
			for _, item := range items {
				if strings.HasPrefix(item.Preview, "---") {
					t.Fatalf("frontmatter leaked into preview for %q", item.ID)
				}
			}
			return run(items)
		}
	}
}

func selectItem(id string) func([]picker.Item) (string, error) {
	return func([]picker.Item) (string, error) { return id, nil }
}

func withPortLookup(processes []PortProcess, lookupErr error) func(*Options) {
	return func(opts *Options) {
		opts.PortLookup = func(port int) ([]PortProcess, error) {
			if lookupErr != nil {
				return nil, lookupErr
			}
			cloned := make([]PortProcess, len(processes))
			copy(cloned, processes)
			for i := range cloned {
				cloned[i].Sockets = append([]string(nil), cloned[i].Sockets...)
				cloned[i].Notes = append([]string(nil), cloned[i].Notes...)
			}
			return cloned, nil
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
