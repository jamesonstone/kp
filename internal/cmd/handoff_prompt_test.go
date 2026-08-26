package cmd

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"
)

func TestHandoffPromptsPrintApprovedContracts(t *testing.T) {
	tests := []struct {
		name     string
		wantHash string
		required []string
	}{
		{
			name:     "agent-handoff",
			wantHash: "31f584b505fc1f50cfa66ca66297e5cc96867d1c1df86eb5aac985ed381282c5",
			required: []string{
				"lossless, zero-context handoff",
				"## Origin phase 1: clarify",
				"## Origin phase 2: emit",
				"## Authority and Safety Boundaries",
				"Repository and Workspace State",
				"Preserve native evidence states",
				"Do not ask the user to attach",
				"emit immediately",
				"Do not repeat completed work.",
				"Proceed with the hydrated task?",
			},
		},
		{
			name:     "chat-handoff",
			wantHash: "9e350823259e3856c63a987a92428772f69872c11a1c640b91d53e91b8cfd063",
			required: []string{
				"zero-context handoff",
				"## Origin phase 1: clarify",
				"## Origin phase 2: emit",
				"## Source Map",
				"Write `UNKNOWN` for missing evidence",
				"## Destination Protocol",
				"Proceed with the hydrated task?",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := executeTestCommand(t, tt.name, "--print")
			if err != nil {
				t.Fatal(err)
			}
			gotHash := fmt.Sprintf("%x", sha256.Sum256([]byte(stdout)))
			if gotHash != tt.wantHash {
				t.Fatalf("output hash = %s, want %s", gotHash, tt.wantHash)
			}
			for _, want := range tt.required {
				if !strings.Contains(stdout, want) {
					t.Errorf("stdout missing %q", want)
				}
			}
			if strings.Contains(stdout, "---") {
				t.Fatal("stdout includes frontmatter")
			}
			if stderr != "" {
				t.Fatalf("stderr = %q", stderr)
			}
		})
	}
}

func TestLegacyHandoffPromptIsRemoved(t *testing.T) {
	_, _, err := executeTestCommand(t, "handoff", "--print")
	if ExitCode(err) != ExitUser {
		t.Fatalf("ExitCode = %d, err = %v", ExitCode(err), err)
	}
}
