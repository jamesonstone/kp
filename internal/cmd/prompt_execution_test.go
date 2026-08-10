package cmd

import (
	"strings"
	"testing"

	"github.com/jamesonstone/kp/internal/clipboard"
)

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

func TestPRPromptPrintsApprovedInstructions(t *testing.T) {
	stdout, stderr, err := executeTestCommand(t, "pr", "--print")
	if err != nil {
		t.Fatal(err)
	}

	want := "create issues, branches, and pull requests for this work, in all project-repositories effected as per our repository and kit-defined rulesets."
	if strings.TrimSpace(stdout) != want {
		t.Fatalf("stdout = %q, want %q", strings.TrimSpace(stdout), want)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestContinuePromptPrintsApprovedInstructions(t *testing.T) {
	stdout, stderr, err := executeTestCommand(t, "continue", "--print")
	if err != nil {
		t.Fatal(err)
	}

	want := "Resolve all issues autonomously and continue until the goal is fully complete; ask permission only before large-scale deletion or deleting sensitive files.\n"
	if stdout != want {
		t.Fatalf("stdout = %q, want %q", stdout, want)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestHandoffPromptIncludesAcceptanceCriteriaWithoutWordLimit(t *testing.T) {
	stdout, stderr, err := executeTestCommand(t, "handoff", "--print")
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"4. Acceptance Criteria",
		"5. Resource Links",
		"Acceptance Criteria: List success criteria in bullet-list form",
		"Each criterion must be binary verifiable.",
		"Remove any sentence where deletion doesn't degrade the agent's ability to execute correctly.",
		"audit the conversation for ambiguities, contradictions, missing requirements, hidden assumptions",
		"Assume the coding agent has not seen the original conversation.",
		"Include in-scope, out-of-scope, and deferred future work",
		"error handling, logging and observability, security and authorization, compatibility, rollback, operator visibility",
		"List decisions already made as implementation constraints",
	}
	for _, text := range expected {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout missing %q:\n%s", text, stdout)
		}
	}
	if strings.Contains(stdout, "650") || strings.Contains(stdout, "word limit") {
		t.Fatalf("stdout includes removed word-limit language:\n%s", stdout)
	}
	if stderr != "" {
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
