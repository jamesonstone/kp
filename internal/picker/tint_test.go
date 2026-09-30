package picker

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestTintKeepsTextAndWidth(t *testing.T) {
	lines := []string{"# Heading", "- item with `code", "  span` done", "", "```", "fenced", "```", "12. step"}
	got := tintLines(lines, newStyles())
	for i := range lines {
		if ansi.Strip(got[i]) != lines[i] {
			t.Fatalf("line %d text changed: %q -> %q", i, lines[i], ansi.Strip(got[i]))
		}
	}
	if got[0] == lines[0] || got[5] == lines[5] {
		t.Fatal("heading and fenced code should be styled")
	}
	if !strings.Contains(got[2], "span") || got[2] == lines[2] {
		t.Fatalf("code span should continue across the wrapped line: %q", got[2])
	}
}

func TestListMarker(t *testing.T) {
	for line, want := range map[string]int{"- a": 2, "  12. b": 6, "plain": 0, "    continued": 0} {
		if got := listMarker(line); got != want {
			t.Fatalf("listMarker(%q) = %d, want %d", line, got, want)
		}
	}
}
