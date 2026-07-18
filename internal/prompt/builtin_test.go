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

	if len(prompts) != 5 {
		t.Fatalf("len(prompts) = %d, want 5", len(prompts))
	}
	if prompts[0].Name != "clarify" || prompts[1].Name != "continue" || prompts[2].Name != "handoff" || prompts[3].Name != "parentthread" || prompts[4].Name != "pr" {
		t.Fatalf("prompt names = %q, %q, %q, %q, %q", prompts[0].Name, prompts[1].Name, prompts[2].Name, prompts[3].Name, prompts[4].Name)
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
	for _, name := range []string{"clarify.md", "continue.md", "handoff.md", "parentthread.md", "pr.md"} {
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
