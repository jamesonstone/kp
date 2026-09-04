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

const approvedStatusSHA256 = "04d3d25f29feeafb63e486f4187d0aebd5a80eab8c3aa3c2a2ff27a101c7e6e0"

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

func TestStatusPromptRequiresSuccinctAuditContract(t *testing.T) {
	stdout, _, err := executeTestCommand(t, "status", "--print")
	if err != nil {
		t.Fatal(err)
	}

	required := []string{
		"Audit whether this thread is actually complete",
		"Be skeptical",
		"Verify current repository and GitHub state",
		"Do not infer completion from implementation",
		"## Summary",
		"Two to four short sentences",
		"## Remaining",
		"numbered list of remaining actions only",
		"One line each",
		"**THREAD COMPLETE:** nothing left from this thread.",
		"Do not emit inventories, item classifications, required-answer sections",
	}
	for _, text := range required {
		if !strings.Contains(stdout, text) {
			t.Fatalf("stdout missing %q", text)
		}
	}

	forbidden := []string{
		"## Thread work inventory",
		"## Item classification",
		"## Required answers",
		"## Completion outcome",
		"**DONE**",
		"**OPEN**",
		"**PARTIAL**",
		"**BLOCKED**",
		"**NOT NEEDED**",
		"Is all work from this thread complete?",
	}
	for _, text := range forbidden {
		if strings.Contains(stdout, text) {
			t.Fatalf("stdout unexpectedly contains %q", text)
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
