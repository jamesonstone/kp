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

const approvedPlanSHA256 = "d36013bd21e8167bd41fa112fe2cb04ac8465cdf3f3f5568236ad7cb9cab64c2"

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
		"treat it as the intent contract",
		"Consume it unchanged rather than rediscovering or redefining the objective",
		"obtain an explicit revised user decision before changing the objective",
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
