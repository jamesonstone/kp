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

const approvedPunchlistSHA256 = "ca6314ad7670ac08d6acb3b95f58df221ced4fc0120c119f6d29b5d3c0e0a533"

func TestPunchlistPromptPrintsApprovedInstructions(t *testing.T) {
	stdout, stderr, err := executeTestCommand(t, "punchlist", "--print")
	if err != nil {
		t.Fatal(err)
	}

	want := approvedPunchlistBody(t)
	if stdout != want {
		t.Fatalf("stdout does not match prompts/punchlist.md body")
	}
	sum := sha256.Sum256([]byte(stdout))
	got := hex.EncodeToString(sum[:])
	if got != approvedPunchlistSHA256 {
		t.Fatalf("punchlist prompt hash = %s, want %s", got, approvedPunchlistSHA256)
	}
	if strings.Contains(stdout, "---") {
		t.Fatalf("stdout includes frontmatter")
	}
	if stderr != "" {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestPunchlistPromptRequiresControlLoopContract(t *testing.T) {
	stdout, _, err := executeTestCommand(t, "punchlist", "--print")
	if err != nil {
		t.Fatal(err)
	}

	required := []string{
		"# Purpose",
		"# Runtime contract",
		"# Control loop",
		"# Invariants",
		"# Cluster before fixing",
		"# Clarification gate",
		"# Worklane strategy",
		"# Engineering notes",
		"# Status handling",
		"# Initial execution",
		"## Punch List state",
		"## Proposed clusters",
		"## Already addressed / awaiting validation",
		"Confidence: NN%",
		"implemented is not merged",
		"Never overwrite, delete, fabricate, or silently dismiss",
		"Do not assume GitHub",
		"mandated delivery-consent confirmation",
		"exact authorized PR set",
		"work-lane-gating",
	}
	for _, text := range required {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout missing %q", text)
		}
	}
}

func approvedPunchlistBody(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("..", "..", "prompts", "punchlist.md"))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := prompt.ParseDocument("punchlist", content)
	if err != nil {
		t.Fatal(err)
	}
	return doc.Body
}
