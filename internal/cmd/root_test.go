package cmd

import (
	"io"
	"strings"
	"testing"

	"github.com/jamesonstone/kp/internal/clipboard"
	"github.com/jamesonstone/kp/internal/prompt"
)

func TestListPlain(t *testing.T) {
	stdout, _, err := executeTestCommand(t, "list", "--plain")
	if err != nil {
		t.Fatal(err)
	}

	if stdout != "merge\nreview\nship\n" {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestListVerbose(t *testing.T) {
	stdout, _, err := executeTestCommand(t, "list", "--verbose")
	if err != nil {
		t.Fatal(err)
	}

	want := "merge\tContext-aware PR merge and deployment\tbuiltin\n" +
		"review\tIndependent branch/PR review\tbuiltin\n" +
		"ship\tPre-authorize task delivery\tbuiltin\n"
	if stdout != want {
		t.Fatalf("stdout = %q, want %q", stdout, want)
	}
}

func TestRootHelpShowsHelpWithoutSideEffects(t *testing.T) {
	var registryCalled bool
	var clipboardCalled bool
	stdout, _, err := executeTestCommand(t, "--help", func(opts *Options) {
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
		"kp merge",
		"kp review",
		"kp ship",
		"Prompt Library",
		"kp task",
		"kp task --output-only",
		"kp list",
		"kp list --plain",
		"Port Tools",
		"kp find-port <port>",
		"kp port-find <port>",
		"Repo Setup",
		"kp scaffold",
		"Legacy v0 Prompts",
		"kp v0 <command>",
		"kp v0 --help",
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
	if strings.Index(stdout, "kp merge") > strings.Index(stdout, "kp list") {
		t.Fatalf("direct prompt commands appear after list commands:\n%s", stdout)
	}
	if strings.Contains(stdout, "kp"+" prompt") {
		t.Fatalf("stdout = %q", stdout)
	}
	if strings.Index(stdout, "Legacy v0 Prompts") < strings.Index(stdout, "Repo Setup") {
		t.Fatalf("legacy section appears before current commands:\n%s", stdout)
	}
	if strings.Contains(stdout, "kp clarify ") || strings.Contains(stdout, "kp handoff") {
		t.Fatalf("stdout includes removed legacy command: %q", stdout)
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

	stdout, _, err := executeTestCommand(t, "--help")
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"\x1b[1;37m🚀 Usage\x1b[0m",
		"\x1b[1;37m🧠 Prompt Commands\x1b[0m",
		"\x1b[1;37m🧰 Prompt Library\x1b[0m",
		"\x1b[1;37m🔍 Port Tools\x1b[0m",
		"\x1b[1;37m🏗️ Repo Setup\x1b[0m",
		"\x1b[1;37m🗄️ Legacy v0 Prompts\x1b[0m",
		"\x1b[1;37m🛠️ Utilities\x1b[0m",
		"\x1b[1;37m⚙️ Flags\x1b[0m",
	}
	for _, text := range expected {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout missing %q:\n%s", text, stdout)
		}
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
