package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesonstone/kp/internal/prompt"
)

const cursorEnvSetupBody = "/env-setup Set up a new development environment for this codebase. Follow the “Set up a new environment” workflow in the env-setup skill. Run the relevant applications and demonstrate that the environment works end to end.\n"

func TestCursorEnvSetupPrintsApprovedPrompt(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("..", "..", "prompts", "cursor", "env-setup.md"))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := prompt.ParseDocument("env-setup", content)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Body != cursorEnvSetupBody {
		t.Fatalf("prompt body = %q", doc.Body)
	}

	stdout, stderr, err := executeTestCommand(t, "cursor", "env-setup", "--print")
	if err != nil {
		t.Fatal(err)
	}
	if stdout != cursorEnvSetupBody || stderr != "" {
		t.Fatalf("stdout = %q, stderr = %q", stdout, stderr)
	}
	if strings.Contains(stdout, "---") {
		t.Fatalf("stdout includes frontmatter")
	}
}

func TestCursorCopiesAndVerifies(t *testing.T) {
	fake := &fakeClipboard{}
	stdout, stderr, err := executeTestCommand(t, "cursor", "env-setup", withClipboard(fake))
	if err != nil {
		t.Fatal(err)
	}
	if stdout != cursorEnvSetupBody || fake.copied != stdout || fake.verified != stdout {
		t.Fatalf("copied=%q verified=%q stdout=%q", fake.copied, fake.verified, stdout)
	}
	if !strings.Contains(stderr, "Prompt \"env-setup\" copied to clipboard.") {
		t.Fatalf("stderr = %q", stderr)
	}

	fake = &fakeClipboard{}
	stdout, _, err = executeTestCommand(t, "cursor", "--copy", "env-setup", withClipboard(fake))
	if err != nil {
		t.Fatal(err)
	}
	if stdout != "" || fake.copied != cursorEnvSetupBody {
		t.Fatalf("--copy stdout=%q copied=%q", stdout, fake.copied)
	}
}

func TestCursorHelpListsEnvSetup(t *testing.T) {
	stdout, _, err := executeTestCommand(t, "cursor")
	if err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"Cursor-specific prompts", "env-setup", "Set up a Cursor cloud environment"} {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout missing %q:\n%s", text, stdout)
		}
	}
}

func TestCursorRejectsUnknownCommand(t *testing.T) {
	_, _, err := executeTestCommand(t, "cursor", "review", "--print")
	if ExitCode(err) != ExitUser || !strings.Contains(err.Error(), `unknown command "review"`) {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}

func TestCursorIsReserved(t *testing.T) {
	_, _, err := executeTestCommand(t, "new", "cursor",
		withEditor(func(string, []string, string) error { return nil }),
	)
	if ExitCode(err) != ExitUser || !strings.Contains(err.Error(), "reserved") {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}

func TestLauncherRunsCursorEnvSetup(t *testing.T) {
	fake := &fakeClipboard{}
	stdout, stderr, err := executeTestCommand(t,
		withClipboard(fake),
		withPicker(t, selectItem("command:cursor:env-setup")),
	)
	if err != nil {
		t.Fatal(err)
	}
	if stdout != cursorEnvSetupBody || fake.copied != stdout || fake.verified != stdout {
		t.Fatalf("copied=%q verified=%q stdout=%q", fake.copied, fake.verified, stdout)
	}
	if !strings.Contains(stderr, "Prompt \"env-setup\" copied to clipboard.") {
		t.Fatalf("stderr = %q", stderr)
	}
}
