package prompt

import "testing"

func TestFrontmatter_WithLabel(t *testing.T) {
	label, body, err := ParseFrontmatter("x", "---\nlabel: X\n---\nbody")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if label != "X" || body != "body" {
		t.Fatalf("got (%q,%q)", label, body)
	}
}

func TestFrontmatter_None(t *testing.T) {
	label, body, err := ParseFrontmatter("foo-bar", "body")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if label != "Foo Bar" || body != "body" {
		t.Fatalf("got (%q,%q)", label, body)
	}
}

func TestFrontmatter_MalformedYAML(t *testing.T) {
	_, _, err := ParseFrontmatter("x", "---\nlabel: [\n---\nbody")
	if err == nil {
		t.Fatal("expected error")
	}
}
