package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jamesonstone/kp/internal/clipboard"
	"github.com/jamesonstone/kp/internal/prompt"
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

func withLauncher(run func(items []LauncherItem) (string, error)) func(*Options) {
	return func(opts *Options) {
		opts.LauncherRunner = run
	}
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
