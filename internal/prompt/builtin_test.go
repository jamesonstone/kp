package prompt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestBuiltIn_LoadsApprovedPrompts(t *testing.T) {
	prompts, err := loadBuiltInsFromDefault()
	if err != nil {
		t.Fatal(err)
	}

	wantNames := []string{"merge", "review", "ship"}
	if len(prompts) != len(wantNames) {
		t.Fatalf("len(prompts) = %d, want %d", len(prompts), len(wantNames))
	}
	for i, want := range wantNames {
		if prompts[i].Name != want {
			t.Fatalf("prompts[%d].Name = %q, want %q", i, prompts[i].Name, want)
		}
	}
	for _, p := range prompts {
		if p.Source != SourceBuiltIn {
			t.Fatalf("%s Source = %v, want SourceBuiltIn", p.Name, p.Source)
		}
		if p.FilePath != "" {
			t.Fatalf("%s FilePath = %q, want empty", p.Name, p.FilePath)
		}
		if p.Body == "" {
			t.Fatalf("%s Body is empty", p.Name)
		}
		if p.Body[:3] == "---" {
			t.Fatalf("%s Body includes frontmatter", p.Name)
		}
	}
}

func TestBuiltIn_V0LoadsLegacyPrompts(t *testing.T) {
	prompts, err := V0BuiltIns()
	if err != nil {
		t.Fatal(err)
	}
	wantNames := []string{"agent-handoff", "chat-handoff", "clarify", "continue", "goal", "parentthread", "plan", "pr", "punchlist", "status"}
	if len(prompts) != len(wantNames) {
		t.Fatalf("len(prompts) = %d, want %d", len(prompts), len(wantNames))
	}
	for i, want := range wantNames {
		if prompts[i].Name != want || prompts[i].Source != SourceBuiltIn || prompts[i].Body == "" {
			t.Fatalf("prompts[%d] = %+v, want built-in %q", i, prompts[i], want)
		}
	}
}

func TestBuiltIn_SourceFilesAreUsablePromptDocuments(t *testing.T) {
	for _, name := range []string{"v0/agent-handoff.md", "v0/chat-handoff.md", "v0/clarify.md", "v0/continue.md", "v0/goal.md", "merge.md", "v0/parentthread.md", "v0/plan.md", "v0/pr.md", "v0/punchlist.md", "review.md", "ship.md", "v0/status.md"} {
		t.Run(name, func(t *testing.T) {
			got, err := os.ReadFile(filepath.Join("..", "..", "prompts", name))
			if err != nil {
				t.Fatal(err)
			}
			doc, err := ParseDocument(strings.TrimSuffix(filepath.Base(name), ".md"), got)
			if err != nil {
				t.Fatal(err)
			}
			if doc.Label == "" {
				t.Fatalf("%s label is empty", name)
			}
			if strings.TrimSpace(doc.Body) == "" {
				t.Fatalf("%s body is empty", name)
			}
		})
	}
}

func TestBuiltIn_RejectsInvalidName(t *testing.T) {
	_, err := loadBuiltIns(fstest.MapFS{
		"Bad.md": {Data: []byte("body")},
	})
	if err == nil {
		t.Fatal("loadBuiltIns error = nil, want invalid name")
	}
}
