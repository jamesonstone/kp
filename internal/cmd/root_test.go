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

	if stdout != "agent-handoff\nchat-handoff\nclarify\ncontinue\ngoal\nmerge\nparentthread\nplan\npr\npunchlist\nship\nstatus\n" {
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
	if !strings.Contains(stdout, "continue\tContinue autonomously\tbuiltin\n") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stdout, "goal\tConstruct an executable goal\tbuiltin\n") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stdout, "agent-handoff\tAgent-to-agent handoff\tbuiltin\n") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stdout, "chat-handoff\tChat-to-agent handoff\tbuiltin\n") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stdout, "merge\tContext-aware PR merge and deployment\tbuiltin\n") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stdout, "plan\tDrive plan to implementation-ready\tbuiltin\n") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stdout, "pr\tPull request workflow\tbuiltin\n") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stdout, "punchlist\tPunch list control loop\tbuiltin\n") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stdout, "ship\tPre-authorize task delivery\tbuiltin\n") {
		t.Fatalf("stdout = %q", stdout)
	}
	if !strings.Contains(stdout, "status\tThread completion audit\tbuiltin\n") {
		t.Fatalf("stdout = %q", stdout)
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
		"kp clarify",
		"kp continue",
		"kp goal",
		"kp agent-handoff",
		"kp chat-handoff",
		"kp merge",
		"kp plan",
		"kp pr",
		"kp punchlist",
		"kp ship",
		"kp status",
		"Prompt Library",
		"kp list",
		"kp list --plain",
		"Port Tools",
		"kp find-port <port>",
		"kp port-find <port>",
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
	if strings.Contains(stdout, "kp handoff") {
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
