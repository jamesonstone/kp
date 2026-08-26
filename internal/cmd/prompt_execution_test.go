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

	want := "Use the Kit workflow implementation path and create a new worklane for this work: a new issue, canonical worktree, `GH-<issue-number>` branch, and ready pull request in every affected project repository."
	if strings.TrimSpace(stdout) != want {
		t.Fatalf("stdout = %q, want %q", strings.TrimSpace(stdout), want)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestMergePromptPrintsApprovedInstructions(t *testing.T) {
	stdout, stderr, err := executeTestCommand(t, "merge", "--print")
	if err != nil {
		t.Fatal(err)
	}

	// The approved body stores "~" in place of a backtick so inline code spans
	// survive a Go raw string literal.
	want := strings.ReplaceAll(`Coordinate the already-scoped pull requests using the conversation and current repository context. Derive the exact PR set, heads, bases, dependencies, permitted merge methods, deployment targets, acceptance gates, and authorization boundaries; do not ask the user to restate discoverable facts or expand scope.

## Analysis and approval

1. Load every applicable repository-local merge, infrastructure, orchestration, and completion rule.
2. Build the dependency and deployment graph from authoritative evidence. Classify each node as ~MERGE_READY~, ~BLOCKED~, or ~UNKNOWN~; missing, stale, pending, conflicted, or unattributable evidence never passes.
3. Record each node's repository, PR, exact head/base, method, dependencies, infrastructure effects, recovery, and acceptance signal. Keep merge, CI, deployment, runtime, and production acceptance distinct.
4. Before the first merge or infrastructure mutation, present one consolidated approval request for the exact current frontier and all known infrastructure effects unless equivalent exact approval from the conversation remains valid. Never delete, destroy, purge, destructively replace, or remove infrastructure; stop and isolate that work for separate explicit authorization.
5. Run one complete preflight immediately before each consequential mutation. Repeat only when a material fact changes: head/base, reviews/checks, policy, actor, dependency, deployment target/effect, approval, or acceptance window.

## Execution

- One primary coordinator owns the graph, authority, wave selection, recovery, and final acceptance.
- When the host supports explicit model selection, use lower-cost or lower-capability agents only for exact, bounded ~MERGE_READY~ merges and deployment monitoring. Keep graph changes, repair decisions, recovery, and acceptance with the coordinator.
- Parallelize nodes only when both source and deployment effects are independent. Serialize shared bases, services, environments, databases, migrations, queues, and acceptance gates.
- Merge only the authorized ready frontier with repository-permitted methods and required queues. Never bypass policy, switch identity, force-push, weaken a gate, or explicitly delete PR branches.
- Monitor with event-driven waits or bounded backoff; do not emit or repeat unchanged polling.
- Reconcile failures autonomously within the approved scope: diagnose once, apply only authorized in-lane repair, rerun affected evidence, and refresh that node and its dependents. Do not retry blindly or introduce a new PR, infrastructure effect, target, method, or authority boundary.
- A changed head returns to ~UNKNOWN~ and requires fresh current-head evidence and authorization. Continue independent valid nodes; stop when recovery cannot make progress safely.

## Final response

Use repository-required status vocabulary; otherwise emit:

~Status: SUCCESS | PARTIAL | BLOCKED | FAILURE~

~Result:~ one to three sentences stating the exact outcome and evidence boundary.

~Next steps:~ ~none~ or at most three specific, copy-ready sentences.

Do not include a chronological work log, repeated checks, unchanged polling, or routine command details.
`, "~", "`")
	if stdout != want {
		t.Fatalf("stdout = %q, want %q", stdout, want)
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
