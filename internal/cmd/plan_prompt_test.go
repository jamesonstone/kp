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

const approvedPlanSHA256 = "5d36e5ce4e46b70c46df111cac6a3bfe5e7af3f2b6bba81ec72504c974e00d1a"

func TestPlanPromptPrintsApprovedInstructions(t *testing.T) {
	stdout, stderr, err := executeTestCommand(t, "plan", "--print")
	if err != nil {
		t.Fatal(err)
	}

	want := approvedPlanBody(t)
	if stdout != want {
		t.Fatalf("stdout does not match prompts/plan.md body")
	}
	sum := sha256.Sum256([]byte(stdout))
	got := hex.EncodeToString(sum[:])
	if got != approvedPlanSHA256 {
		t.Fatalf("plan prompt hash = %s, want %s", got, approvedPlanSHA256)
	}
	if strings.Contains(stdout, "---") {
		t.Fatalf("stdout includes frontmatter")
	}
	if stderr != "" {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestPlanPromptRequiresConvergenceContract(t *testing.T) {
	stdout, _, err := executeTestCommand(t, "plan", "--print")
	if err != nil {
		t.Fatal(err)
	}

	required := []string{
		"Drive the current implementation plan to a practical, evidence-backed local maximum",
		"Stay in planning mode",
		"do not edit files, create delivery artifacts, deploy, or execute the plan",
		"simplest decision-complete plan",
		"Do not maximize plan length or complexity",
		"Research progressively",
		"Treat historical plans and prior conclusions as hypotheses",
		"Stop researching when additional retrieval would improve wording",
		"smallest consolidated batch of numbered questions",
		"Silently iterate",
		"zero unresolved material questions",
		"at least 95% evidence-backed goal coverage",
		"one complete adversarial pass produces no new material recommendation",
		"a final verification pass also produces no new material recommendation",
		"Do not inflate confidence",
		"Return one complete replacement plan",
		"local-maximum audit",
	}
	for _, text := range required {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout missing %q", text)
		}
	}
}

func approvedPlanBody(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("..", "..", "prompts", "plan.md"))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := prompt.ParseDocument("plan", content)
	if err != nil {
		t.Fatal(err)
	}
	return doc.Body
}
