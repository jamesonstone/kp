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

const approvedShipSHA256 = "a12abcf9c299f274a6933843cd136a0129ebd848227c28718e8563ef30c5b2b8"

func TestShipPromptPrintsApprovedInstructions(t *testing.T) {
	stdout, stderr, err := executeTestCommand(t, "ship", "--print")
	if err != nil {
		t.Fatal(err)
	}

	want := approvedShipBody(t)
	if stdout != want {
		t.Fatalf("stdout does not match prompts/ship.md body")
	}
	sum := sha256.Sum256([]byte(stdout))
	got := hex.EncodeToString(sum[:])
	if got != approvedShipSHA256 {
		t.Fatalf("ship prompt hash = %s, want %s", got, approvedShipSHA256)
	}
	if strings.Contains(stdout, "---") {
		t.Fatalf("stdout includes frontmatter")
	}
	if stderr != "" {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestShipPromptRequiresDeliveryAuthorizationContract(t *testing.T) {
	stdout, stderr, err := executeTestCommand(t, "ship", "--print")
	if err != nil {
		t.Fatal(err)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q", stderr)
	}

	for _, want := range []string{
		"/goal For this task thread",
		"pre-authorized to complete the full delivery lifecycle",
		"creating, updating, and pushing branches",
		"creating and updating pull requests",
		"addressing review feedback and CI failures",
		"merging pull requests when required checks pass",
		"merging dependent pull requests",
		"Do not ask for additional authorization for individual PR merges",
		"scoped only to changes required to accomplish this task",
		"Do not merge unrelated pre-existing PRs",
		"bypass required protections",
		"ignore failing required checks",
		"Continue autonomously until the task is delivered",
		"so you can finish without waiting on me for routine approvals",
		"This message is the explicit merge authority your repository rules ask for",
		"anything else they reserve for separate approval still needs it",
		"The task is delivered when every in-scope pull request is merged",
		"report the outcome with evidence, keeping merge, deployment, and runtime results separate",
		"say what you need and from whom",
	} {
		if !strings.Contains(stdout, want) {
			t.Errorf("stdout missing %q", want)
		}
	}
}

func approvedShipBody(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("..", "..", "prompts", "ship.md"))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := prompt.ParseDocument("ship", content)
	if err != nil {
		t.Fatal(err)
	}
	return doc.Body
}
