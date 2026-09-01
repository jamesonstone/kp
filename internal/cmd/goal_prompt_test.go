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

const approvedGoalSHA256 = "ac318360ff5f848bf8fc9d673ebb22515f4aa13d6bfe3ec1aeda1e2043a8e308"

func TestGoalPromptPrintsApprovedInstructions(t *testing.T) {
	stdout, stderr, err := executeTestCommand(t, "goal", "--print")
	if err != nil {
		t.Fatal(err)
	}

	want := approvedGoalBody(t)
	if stdout != want {
		t.Fatalf("stdout does not match prompts/goal.md body")
	}
	sum := sha256.Sum256([]byte(stdout))
	got := hex.EncodeToString(sum[:])
	if got != approvedGoalSHA256 {
		t.Fatalf("goal prompt hash = %s, want %s", got, approvedGoalSHA256)
	}
	if strings.Contains(stdout, "---") {
		t.Fatalf("stdout includes frontmatter")
	}
	if stderr != "" {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestGoalPromptRequiresGoalConstructionContracts(t *testing.T) {
	stdout, stderr, err := executeTestCommand(t, "goal", "--print")
	if err != nil {
		t.Fatal(err)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q", stderr)
	}

	for _, want := range []string{
		"Construct an executable `/goal`",
		"one accumulated goal model",
		"before asking the user",
		"smallest useful numbered batch",
		"circular dependencies",
		"bootstrap work rather than a blocking precondition",
		"what PASS proves and what it must not claim",
		"no known unresolved material question remains",
		"one adversarial review",
		"Ask the user to confirm the complete `/goal`",
	} {
		if !strings.Contains(stdout, want) {
			t.Errorf("stdout missing %q", want)
		}
	}
}

func approvedGoalBody(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("..", "..", "prompts", "goal.md"))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := prompt.ParseDocument("goal", content)
	if err != nil {
		t.Fatal(err)
	}
	return doc.Body
}
