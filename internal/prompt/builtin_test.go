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

	if len(prompts) != 10 {
		t.Fatalf("len(prompts) = %d, want 10", len(prompts))
	}
	wantNames := []string{
		"agent-handoff",
		"chat-handoff",
		"clarify",
		"continue",
		"goal",
		"merge",
		"parentthread",
		"plan",
		"pr",
		"punchlist",
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

func TestBuiltIn_SourceFilesAreUsablePromptDocuments(t *testing.T) {
	for _, name := range []string{"agent-handoff.md", "chat-handoff.md", "clarify.md", "continue.md", "goal.md", "merge.md", "parentthread.md", "plan.md", "pr.md", "punchlist.md"} {
		t.Run(name, func(t *testing.T) {
			got, err := os.ReadFile(filepath.Join("..", "..", "prompts", name))
			if err != nil {
				t.Fatal(err)
			}
			doc, err := ParseDocument(strings.TrimSuffix(name, ".md"), got)
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
