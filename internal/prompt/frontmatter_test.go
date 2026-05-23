package prompt

import (
	"errors"
	"testing"
)

func TestFrontmatter_WithLabel(t *testing.T) {
	doc, err := ParseDocument("clarify", []byte("---\nlabel: Clarify before implementing\n---\nbody"))
	if err != nil {
		t.Fatal(err)
	}

	if doc.Label != "Clarify before implementing" {
		t.Fatalf("Label = %q", doc.Label)
	}
	if doc.Body != "body" {
		t.Fatalf("Body = %q", doc.Body)
	}
}

func TestFrontmatter_None(t *testing.T) {
	doc, err := ParseDocument("foo-bar", []byte("body"))
	if err != nil {
		t.Fatal(err)
	}

	if doc.Label != "Foo Bar" {
		t.Fatalf("Label = %q", doc.Label)
	}
	if doc.Body != "body" {
		t.Fatalf("Body = %q", doc.Body)
	}
}

func TestFrontmatter_MalformedYAML(t *testing.T) {
	_, err := ParseDocument("bad", []byte("---\nlabel: [\n---\nbody"))
	if err == nil {
		t.Fatal("ParseDocument error = nil, want YAML parse error")
	}
}

func TestFrontmatter_MissingClosingDelimiter(t *testing.T) {
	_, err := ParseDocument("bad", []byte("---\nlabel: X\nbody"))
	if err == nil {
		t.Fatal("ParseDocument error = nil, want missing delimiter error")
	}
}

func TestFrontmatter_LabelDefaultsWhenMissing(t *testing.T) {
	doc, err := ParseDocument("foo-bar", []byte("---\nother: value\n---\nbody"))
	if err != nil {
		t.Fatal(err)
	}

	if doc.Label != "Foo Bar" {
		t.Fatalf("Label = %q", doc.Label)
	}
	if doc.Body != "body" {
		t.Fatalf("Body = %q", doc.Body)
	}
}

func TestFrontmatter_EmptyBody(t *testing.T) {
	doc, err := ParseDocument("empty", []byte("---\nlabel: Empty\n---\n"))
	if err != nil {
		t.Fatal(err)
	}

	if doc.Body != "" {
		t.Fatalf("Body = %q, want empty", doc.Body)
	}
}

func TestFrontmatter_LiteralDelimiterBodyWithoutOpeningFrontmatter(t *testing.T) {
	doc, err := ParseDocument("literal", []byte("--- body\ntext"))
	if err != nil {
		t.Fatal(err)
	}

	if doc.Body != "--- body\ntext" {
		t.Fatalf("Body = %q", doc.Body)
	}
}

func TestNameValidation_AcceptsValid(t *testing.T) {
	for _, name := range []string{"a", "abc", "foo-1", "foo-bar"} {
		if err := ValidateName(name); err != nil {
			t.Fatalf("ValidateName(%q) error = %v", name, err)
		}
	}
}

func TestNameValidation_RejectsInvalid(t *testing.T) {
	for _, name := range []string{"", "Foo", "1abc", "foo bar", "foo_bar", "foo.bar", "foo/bar"} {
		err := ValidateName(name)
		if !errors.Is(err, ErrInvalidName) {
			t.Fatalf("ValidateName(%q) error = %v, want ErrInvalidName", name, err)
		}
	}
}

func TestNameValidation_RejectsReserved(t *testing.T) {
	for _, name := range reservedNameList() {
		err := ValidateName(name)
		if !errors.Is(err, ErrReservedName) {
			t.Fatalf("ValidateName(%q) error = %v, want ErrReservedName", name, err)
		}
	}
}
