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

const approvedReviewSHA256 = "ab9233c17c44e46658d68f82d90e229319806d5085bbde7d315e7ae159234999"

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
