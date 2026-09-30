package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesonstone/kp/internal/prompt"
)

var v0Names = []string{"agent-handoff", "chat-handoff", "clarify", "continue", "goal", "parentthread", "plan", "pr", "punchlist", "status"}

func TestV0ServesEveryLegacyPromptUnchanged(t *testing.T) {
	for _, name := range v0Names {
		t.Run(name, func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join("..", "..", "prompts", "v0", name+".md"))
			if err != nil {
				t.Fatal(err)
			}
			doc, err := prompt.ParseDocument(name, content)
			if err != nil {
				t.Fatal(err)
			}

			stdout, stderr, err := executeTestCommand(t, "v0", name, "--print")
			if err != nil {
				t.Fatal(err)
			}
			if stdout != doc.Body || stderr != "" {
				t.Fatalf("stdout matches body = %v, stderr = %q", stdout == doc.Body, stderr)
			}
		})
	}
}

func TestV0CopiesAndVerifiesLikeRootPrompts(t *testing.T) {
	fake := &fakeClipboard{}
	stdout, stderr, err := executeTestCommand(t, "v0", "goal", withClipboard(fake))
	if err != nil {
		t.Fatal(err)
	}
	if stdout == "" || fake.copied != stdout || fake.verified != stdout {
		t.Fatalf("copied=%q verified=%q stdout=%q", fake.copied, fake.verified, stdout)
	}
	if !strings.Contains(stderr, "Prompt \"goal\" copied to clipboard.") {
		t.Fatalf("stderr = %q", stderr)
	}

	fake = &fakeClipboard{}
	stdout, _, err = executeTestCommand(t, "v0", "--copy", "goal", withClipboard(fake))
	if err != nil {
		t.Fatal(err)
	}
	if stdout != "" || fake.copied == "" {
		t.Fatalf("--copy stdout=%q copied=%q", stdout, fake.copied)
	}
}

func TestRootNoLongerResolvesLegacyNames(t *testing.T) {
	for _, name := range v0Names {
		_, _, err := executeTestCommand(t, name, "--print")
		if ExitCode(err) != ExitUser {
			t.Fatalf("%s: ExitCode = %d, err = %v", name, ExitCode(err), err)
		}
		if !strings.Contains(err.Error(), "kp v0 "+name) {
			t.Fatalf("%s: err = %v, want a pointer to kp v0", name, err)
		}
	}
}

func TestUserPromptCanReuseLegacyName(t *testing.T) {
	configDir := t.TempDir()
	writeUserPrompt(t, configDir, "clarify", "My own clarify.\n")

	stdout, _, err := executeTestCommandWithConfig(t, configDir, "clarify", "--print")
	if err != nil {
		t.Fatal(err)
	}
	if stdout != "My own clarify.\n" {
		t.Fatalf("stdout = %q", stdout)
	}

	stdout, _, err = executeTestCommandWithConfig(t, configDir, "v0", "clarify", "--print")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(stdout, "Clarify before implementing.") {
		t.Fatalf("v0 stdout = %q, want the built-in legacy prompt", stdout)
	}
}

func TestV0RejectsUnknownCommand(t *testing.T) {
	_, _, err := executeTestCommand(t, "v0", "review", "--print")
	if ExitCode(err) != ExitUser || !strings.Contains(err.Error(), `unknown command "review"`) {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}

func TestV0IsReserved(t *testing.T) {
	_, _, err := executeTestCommand(t, "new", "v0",
		withEditor(func(string, []string, string) error { return nil }),
	)
	if ExitCode(err) != ExitUser || !strings.Contains(err.Error(), "reserved") {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}
