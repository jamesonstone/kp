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
	want := strings.ReplaceAll(`Coordinate the already-scoped pull requests using the conversation and current repository context. Derive the exact PR set, heads, bases, dependencies, merge methods, deployment targets, acceptance gates, and authorization boundaries. Do not ask the user to restate discoverable facts or expand scope.

## Analysis and approval

1. Before any mutation, load ~docs/agents/GUARDRAILS.md~ and applicable lane, merge, infrastructure, orchestration, testing, and completion rules. Complete read-only safety reconnaissance.
2. For repository or delivery work, verify that the lane rule's exact consent question was answered for the current lane. If not, ask it verbatim and wait; generic approval is not lane consent.
3. Build the dependency/deployment graph from authoritative evidence. Derive edge direction only from an explicit base relationship, producer/consumer contract (for example, a backend endpoint a frontend PR calls), or stated prerequisite — never from shared files or proximity alone — and never merge a consumer ahead of its producer. Classify every node as ~MERGE_READY~, ~BLOCKED~, or ~UNKNOWN~; missing, stale, pending, conflicted, unattributable, or open-review-feedback evidence never passes — route open feedback through ~pr-feedback-repair~ before reclassifying.
4. Record each node's repository, PR, exact head/base, method, dependencies, infrastructure effects, recovery, and acceptance signal. Keep merge, CI, deployment, runtime, and production acceptance distinct.
5. Before the first merge or infrastructure mutation, present one consolidated approval request for the exact frontier and known effects unless equivalent exact approval remains valid.
6. Never delete, destroy, purge, remove, or destructively replace infrastructure. Isolate such work for separate explicit authorization.
7. Run one complete preflight immediately before each consequential mutation. Repeat only after material change to head/base, policy, actor, checks/reviews, dependencies, target/effect, approval, or acceptance window.

## Deadline validation budget

Use operational correctness, not exhaustive correctness.

- Do not bypass required protections, reviews, queues, or hosted checks.
- Reuse fresh exact-head evidence; do not rerun an unchanged local or hosted check.
- Before merge, require only exact identity/policy/readiness evidence.
- After deployment, verify only:
  1. exact source, image digest, task definition, and workflow identity;
  2. stable healthy service/task counts and health endpoint;
  3. required runtime configuration and unchanged secret bindings;
  4. one smallest focused live assertion directly proving the changed behavior.
- Stop testing immediately when that focused assertion passes.
- Do not run full production aggregates, broad end-to-end suites, UI/browser automation, restart journeys, migrations, backfills, or adjacent workflows unless explicitly required by repository policy or named by the user.
- Do not create unrelated fixtures, vendor orders, receiving events, secondary-user flows, or external side effects.
- Record excluded suites as ~NOT_RUN_BY_INSTRUCTION~; they are not acceptance blockers under this deadline scope.
- If the focused assertion fails, diagnose once, repair only within authorized scope, and rerun only the affected evidence. Never broaden testing automatically.

## Execution

- One coordinator owns the graph, authority, wave selection, recovery, and acceptance.
- Use lower-cost agents only for exact bounded ~MERGE_READY~ merges or monitoring when explicit model selection is supported. Keep graph changes, repairs, recovery, and acceptance with the coordinator.
- Parallelize only source and deployment-independent nodes. Serialize shared bases, services, environments, databases, migrations, queues, and gates.
- Merge only the authorized ready frontier using permitted methods and required queues. Never bypass policy, switch identity, force-push, weaken gates, or explicitly delete PR branches.
- Monitor with event-driven waits or bounded backoff; do not repeat unchanged polling.
- Reconcile failures autonomously inside approved scope. Do not retry blindly or introduce a new PR, target, method, infrastructure effect, or authority boundary.
- A changed head returns to ~UNKNOWN~ and requires fresh exact-head evidence and authorization. Continue independent valid nodes; stop only when safe recovery cannot progress.

## Acceptance

A node is accepted when its authorized merge completes, configured deployment succeeds, exact runtime identity and health are verified, and its single focused operational assertion passes.

## Final response

Use repository-required status vocabulary; otherwise emit:

~Status: SUCCESS | PARTIAL | BLOCKED | FAILURE~

~Result:~ one to three sentences stating the exact outcome and evidence boundary.

~Next steps:~ ~none~ or at most three specific, copy-ready sentences.

Do not include a chronological log, repeated checks, unchanged polling, routine commands, or unrequested testing detail.
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
