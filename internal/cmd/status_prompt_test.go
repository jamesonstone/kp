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

const approvedStatusSHA256 = "fc5f19f044293dbcd5ee2b962a676bea158d3ce85566ad805303437aa0acdd1d"

func TestStatusPromptPrintsApprovedInstructions(t *testing.T) {
	stdout, stderr, err := executeTestCommand(t, "status", "--print")
	if err != nil {
		t.Fatal(err)
	}

	want := approvedStatusBody(t)
	if stdout != want {
		t.Fatalf("stdout does not match prompts/status.md body")
	}
	sum := sha256.Sum256([]byte(stdout))
	got := hex.EncodeToString(sum[:])
	if got != approvedStatusSHA256 {
		t.Fatalf("status prompt hash = %s, want %s", got, approvedStatusSHA256)
	}
	if strings.Contains(stdout, "---") {
		t.Fatalf("stdout includes frontmatter")
	}
	if stderr != "" {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestStatusPromptRequiresCompletionAuditContract(t *testing.T) {
	stdout, _, err := executeTestCommand(t, "status", "--print")
	if err != nil {
		t.Fatal(err)
	}

	required := []string{
		"# Purpose",
		"# Runtime contract",
		"# Audit scope",
		"# Classification",
		"# Required answers",
		"# Completion statement",
		"# Initial execution",
		"## Thread work inventory",
		"## Item classification",
		"## Required answers",
		"## Completion outcome",
		"comprehensive completion audit",
		"Do not infer completion merely because implementation work occurred",
		"Verify the actual current repository and GitHub state",
		"**DONE**",
		"**OPEN**",
		"**PARTIAL**",
		"**BLOCKED**",
		"**NOT NEEDED**",
		"Is all work from this thread complete?",
		"**THREAD COMPLETE:**",
		"prioritized completion checklist",
		"Search for omissions rather than trying to justify completion",
	}
	for _, text := range required {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout missing %q", text)
		}
	}
}

func approvedStatusBody(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("..", "..", "prompts", "status.md"))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := prompt.ParseDocument("status", content)
	if err != nil {
		t.Fatal(err)
	}
	return doc.Body
}
