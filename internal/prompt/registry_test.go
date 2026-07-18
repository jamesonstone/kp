package prompt

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRegistry_ListSorted(t *testing.T) {
	reg := newTestRegistry(t)

	prompts := reg.List()
	got := make([]string, len(prompts))
	for i, prompt := range prompts {
		got[i] = prompt.Name
	}
	want := []string{"clarify", "continue", "handoff", "parentthread", "pr"}
	if len(got) != len(want) {
		t.Fatalf("names = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("names = %v, want %v", got, want)
		}
	}
}

func TestRegistry_UserOverrides(t *testing.T) {
	dir := t.TempDir()
	writePrompt(t, dir, "handoff.md", "---\nlabel: User handoff\n---\nuser body")

	reg, err := NewRegistry(dir)
	if err != nil {
		t.Fatal(err)
	}
	p, err := reg.Get("handoff")
	if err != nil {
		t.Fatal(err)
	}

	if p.Source != SourceUser {
		t.Fatalf("Source = %v, want SourceUser", p.Source)
	}
	if p.Label != "User handoff" {
		t.Fatalf("Label = %q", p.Label)
	}
	if p.Body != "user body" {
		t.Fatalf("Body = %q", p.Body)
	}
}

func TestRegistry_AddRejectsEmpty(t *testing.T) {
	reg := newTestRegistry(t)

	_, err := reg.Add("x", "")
	if !errors.Is(err, ErrEmpty) {
		t.Fatalf("Add error = %v, want ErrEmpty", err)
	}
}

func TestRegistry_AddRejectsCollision(t *testing.T) {
	reg := newTestRegistry(t)

	_, err := reg.Add("handoff", "body")
	if !errors.Is(err, ErrExists) {
		t.Fatalf("Add error = %v, want ErrExists", err)
	}
}

func TestRegistry_AddCreatesUserPrompt(t *testing.T) {
	dir := t.TempDir()
	reg, err := NewRegistry(dir)
	if err != nil {
		t.Fatal(err)
	}

	p, err := reg.Add("custom", "---\nlabel: Custom\n---\nbody")
	if err != nil {
		t.Fatal(err)
	}

	if p.Source != SourceUser {
		t.Fatalf("Source = %v, want SourceUser", p.Source)
	}
	if p.FilePath == "" {
		t.Fatal("FilePath is empty")
	}
	if p.Body != "body" {
		t.Fatalf("Body = %q", p.Body)
	}
	if _, err := os.Stat(filepath.Join(dir, "custom.md")); err != nil {
		t.Fatal(err)
	}
}

func TestRegistry_RemoveBuiltinFails(t *testing.T) {
	reg := newTestRegistry(t)

	err := reg.Remove("handoff")
	if !errors.Is(err, ErrBuiltIn) {
		t.Fatalf("Remove error = %v, want ErrBuiltIn", err)
	}
}

func TestRegistry_RemoveUserPrompt(t *testing.T) {
	dir := t.TempDir()
	writePrompt(t, dir, "custom.md", "body")
	reg, err := NewRegistry(dir)
	if err != nil {
		t.Fatal(err)
	}

	if err := reg.Remove("custom"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "custom.md")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("stat removed file error = %v, want os.ErrNotExist", err)
	}
	if _, err := reg.Get("custom"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get removed prompt error = %v, want ErrNotFound", err)
	}
}

func TestRegistry_RemoveUserOverrideRestoresBuiltin(t *testing.T) {
	dir := t.TempDir()
	writePrompt(t, dir, "clarify.md", "user body")
	reg, err := NewRegistry(dir)
	if err != nil {
		t.Fatal(err)
	}

	if err := reg.Remove("clarify"); err != nil {
		t.Fatal(err)
	}
	p, err := reg.Get("clarify")
	if err != nil {
		t.Fatal(err)
	}
	if p.Source != SourceBuiltIn {
		t.Fatalf("Source = %v, want SourceBuiltIn", p.Source)
	}
}

func TestRegistry_PromoteToUserCopiesSource(t *testing.T) {
	dir := t.TempDir()
	reg, err := NewRegistry(dir)
	if err != nil {
		t.Fatal(err)
	}

	p, err := reg.PromoteToUser("clarify")
	if err != nil {
		t.Fatal(err)
	}
	if p.Source != SourceUser {
		t.Fatalf("Source = %v, want SourceUser", p.Source)
	}

	got, err := os.ReadFile(filepath.Join(dir, "clarify.md"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := loadBuiltInSource("clarify")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatal("promoted source does not match built-in source")
	}
}

func TestRegistry_RejectsInvalidUserFileName(t *testing.T) {
	dir := t.TempDir()
	writePrompt(t, dir, "Bad.md", "body")

	_, err := NewRegistry(dir)
	if !errors.Is(err, ErrInvalidName) {
		t.Fatalf("NewRegistry error = %v, want ErrInvalidName", err)
	}
}

func TestRegistry_ReportsMalformedUserFrontmatterPath(t *testing.T) {
	dir := t.TempDir()
	writePrompt(t, dir, "custom.md", "---\nlabel: [\n---\nbody")

	_, err := NewRegistry(dir)
	if err == nil {
		t.Fatal("NewRegistry error = nil, want frontmatter error")
	}
	if !strings.Contains(err.Error(), filepath.Join(dir, "custom.md")) {
		t.Fatalf("error = %q, want path", err.Error())
	}
}

func newTestRegistry(t *testing.T) Registry {
	t.Helper()
	reg, err := NewRegistry(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return reg
}

func writePrompt(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
