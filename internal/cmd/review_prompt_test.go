package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesonstone/kp/internal/prompt"
)

const approvedReviewSHA256 = "9089c22d09d920182bfa4318eb4b2deaf89df130cea68ab44d264da45938b036"

func TestReviewPromptPrintsApprovedInstructions(t *testing.T) {
	stdout, stderr, err := executeTestCommand(t, "review", "--print")
	if err != nil {
		t.Fatal(err)
	}

	want := approvedReviewBody(t)
	if stdout != want {
		t.Fatalf("stdout does not match prompts/review.md body")
	}
	sum := sha256.Sum256([]byte(stdout))
	got := hex.EncodeToString(sum[:])
	if got != approvedReviewSHA256 {
		t.Fatalf("review prompt hash = %s, want %s", got, approvedReviewSHA256)
	}
	if strings.Contains(stdout, "---") {
		t.Fatalf("stdout includes frontmatter")
	}
	if stderr != "" {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestReviewPromptRequiresCorrectnessContract(t *testing.T) {
	stdout, _, err := executeTestCommand(t, "review", "--print")
	if err != nil {
		t.Fatal(err)
	}

	required := []string{
		"Perform an independent review of the changes in this branch/PR against the target branch",
		"Treat the existing implementation as untrusted",
		"Do not assume the approach is correct because another agent wrote it or because tests currently pass",
		"Focus primarily on correctness",
		"incorrect behavior or logic",
		"edge cases and boundary conditions",
		"error handling and failure modes",
		"concurrency, ordering, idempotency, and retry issues where applicable",
		"state consistency and transactional correctness",
		"security or data-integrity problems",
		"violations of existing contracts, invariants, or repository conventions",
		"missing or insufficient tests",
		"unnecessary complexity that creates correctness risk",
		"Trace important execution paths rather than reviewing only the diff syntactically",
		"Run the relevant tests, linters, type checks, and other repository validation available to you",
		"Do not make changes yet",
		"Return findings ordered by severity",
		"the smallest appropriate fix",
		"test coverage that should prove the fix",
		"Do not manufacture findings",
		"If the implementation is correct, say so explicitly",
		"identify any remaining validation gaps or residual risks",
		"review the merge-base diff against its remote copy",
		"read the change's stated intent",
		"behavior that differs from the change's stated intent",
		"without modifying tracked files or external state",
		"If a check cannot run, say which one and why",
		"Search broadly, then confirm each finding before reporting it",
		"a scratch test in a temporary copy of the repository",
		"List concerns you cannot confirm as open questions rather than dropping them",
		"Report problems this change introduces or makes reachable",
		"list pre-existing issues you notice separately as follow-ups",
		"Skip style or preference nits unless a repository rule requires them or they create correctness risk",
		"P0 for data loss, an exploitable security flaw, or an outage",
	}
	for _, text := range required {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout missing %q", text)
		}
	}
}

func approvedReviewBody(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("..", "..", "prompts", "review.md"))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := prompt.ParseDocument("review", content)
	if err != nil {
		t.Fatal(err)
	}
	return doc.Body
}
