package prompt

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	kp "github.com/jamesonstone/kp"
)

func TestRegistry_BuiltIn(t *testing.T) {
	r, err := NewRegistry(kp.BuiltinPromptsFS, t.TempDir())
	if err != nil {
		t.Fatalf("NewRegistry error: %v", err)
	}
	if got := len(r.List()); got != 3 {
		t.Fatalf("got %d prompts want 3", got)
	}
}

func TestRegistry_UserOverrides(t *testing.T) {
	d := t.TempDir()
	if err := os.WriteFile(filepath.Join(d, "instructions.md"), []byte("user body"), 0o600); err != nil {
		t.Fatal(err)
	}
	r, err := NewRegistry(kp.BuiltinPromptsFS, d)
	if err != nil {
		t.Fatal(err)
	}
	p, err := r.Get("instructions")
	if err != nil {
		t.Fatal(err)
	}
	if p.Source != SourceUser || p.Body != "user body" {
		t.Fatalf("override failed: %#v", p)
	}
}

func TestRegistry_AddRejectsEmpty(t *testing.T) {
	r, err := NewRegistry(kp.BuiltinPromptsFS, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	_, err = r.Add("x", "")
	if !errors.Is(err, ErrEmpty) {
		t.Fatalf("got %v", err)
	}
}

func TestRegistry_AddRejectsInvalidName(t *testing.T) {
	r, err := NewRegistry(kp.BuiltinPromptsFS, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := r.Add("Bad Name", "body"); err == nil {
		t.Fatal("expected error")
	}
}

func TestRegistry_RemoveBuiltinFails(t *testing.T) {
	r, err := NewRegistry(kp.BuiltinPromptsFS, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Remove("instructions"); !errors.Is(err, ErrBuiltIn) {
		t.Fatalf("got %v", err)
	}
}

func TestRegistry_PromoteToUser_CopiesBody(t *testing.T) {
	d := t.TempDir()
	r, err := NewRegistry(kp.BuiltinPromptsFS, d)
	if err != nil {
		t.Fatal(err)
	}
	before, err := r.Get("clarify")
	if err != nil {
		t.Fatal(err)
	}
	after, err := r.PromoteToUser("clarify")
	if err != nil {
		t.Fatal(err)
	}
	if after.Source != SourceUser || after.Body != before.Body {
		t.Fatalf("bad promote: %#v", after)
	}
	if _, err := os.Stat(filepath.Join(d, "clarify.md")); err != nil {
		t.Fatal(err)
	}
}

func TestNameValidation_RejectsInvalid(t *testing.T) {
	for _, n := range []string{"Foo", "1abc", "foo bar"} {
		if err := ValidateName(n); err == nil {
			t.Fatalf("expected invalid for %q", n)
		}
	}
}
