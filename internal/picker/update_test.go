package picker

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func testItems() []Item {
	return []Item{
		{ID: "merge", Title: "Merge", Command: "kp merge", Group: "prompts", Preview: "merge body"},
		{ID: "review", Title: "Review", Command: "kp review", Group: "prompts", Preview: "review body"},
		{ID: "ship", Title: "Ship", Command: "kp ship", Group: "prompts", Preview: "ship body"},
		{ID: "init", Title: "Init", Command: "kp init", Group: "commands", Preview: "init body"},
		{ID: "help", Title: "Help", Command: "kp --help", Group: "commands", Preview: "help body"},
	}
}

func newTestModel(items []Item) model {
	m := newModel("kp", items, newStyles())
	m.resize(100, 30)
	return m
}

func key(code rune, mod tea.KeyMod) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: code, Mod: mod}
}

func text(s string) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: []rune(s)[0], Text: s}
}

// press feeds messages through Update, as Bubble Tea would, and reports
// whether the last one asked the program to quit.
func press(t *testing.T, m model, msgs ...tea.Msg) (model, bool) {
	t.Helper()
	quit := false
	for _, msg := range msgs {
		next, cmd := m.Update(msg)
		m = next.(model)
		quit = cmd != nil
	}
	return m, quit
}

func selectedID(m model) string {
	item, _ := m.current()
	return item.ID
}

func TestNavigationKeysMoveAndWrap(t *testing.T) {
	tests := []struct {
		name string
		msgs []tea.Msg
		want string
	}{
		{"j moves down", []tea.Msg{text("j")}, "review"},
		{"down arrow", []tea.Msg{key(tea.KeyDown, 0)}, "review"},
		{"tab", []tea.Msg{key(tea.KeyTab, 0)}, "review"},
		{"ctrl+n", []tea.Msg{key('n', tea.ModCtrl)}, "review"},
		{"k wraps to last", []tea.Msg{text("k")}, "help"},
		{"up wraps to last", []tea.Msg{key(tea.KeyUp, 0)}, "help"},
		{"shift+tab", []tea.Msg{key(tea.KeyTab, tea.ModShift)}, "help"},
		{"down wraps to first", []tea.Msg{key(tea.KeyEnd, 0), text("j")}, "merge"},
		{"end then home", []tea.Msg{key(tea.KeyEnd, 0), key(tea.KeyHome, 0)}, "merge"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, quit := press(t, newTestModel(testItems()), tt.msgs...)
			if quit {
				t.Fatal("navigation quit the picker")
			}
			if got := selectedID(m); got != tt.want {
				t.Fatalf("selected = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEnterSelectsCurrentItem(t *testing.T) {
	m, quit := press(t, newTestModel(testItems()), text("j"), text("j"), key(tea.KeyEnter, 0))
	if !quit || m.selected != "ship" {
		t.Fatalf("quit=%v selected=%q", quit, m.selected)
	}
	if m.render() != "" {
		t.Fatal("picker should clear its view after selection")
	}
}

func TestEscAndCtrlCCancel(t *testing.T) {
	for _, msg := range []tea.Msg{key(tea.KeyEscape, 0), key('c', tea.ModCtrl)} {
		m, quit := press(t, newTestModel(testItems()), msg)
		if !quit || m.selected != "" {
			t.Fatalf("%v: quit=%v selected=%q", msg, quit, m.selected)
		}
	}
}

func TestTypingFiltersAndEscClearsBeforeCanceling(t *testing.T) {
	m, quit := press(t, newTestModel(testItems()), text("r"), text("e"), text("v"))
	if quit || !m.filtering || m.query != "rev" {
		t.Fatalf("quit=%v filtering=%v query=%q", quit, m.filtering, m.query)
	}
	if len(m.visible) != 1 || selectedID(m) != "review" {
		t.Fatalf("visible=%v selected=%q", m.visible, selectedID(m))
	}

	m, quit = press(t, m, key(tea.KeyEscape, 0))
	if quit || m.filtering || m.query != "" || len(m.visible) != len(m.items) {
		t.Fatalf("after esc: quit=%v filtering=%v query=%q visible=%d", quit, m.filtering, m.query, len(m.visible))
	}
	if selectedID(m) != "review" {
		t.Fatalf("clearing the filter lost the selection: %q", selectedID(m))
	}

	if _, quit = press(t, m, key(tea.KeyEscape, 0)); !quit {
		t.Fatal("second esc should cancel")
	}
}

func TestFilterModeTypesJAndK(t *testing.T) {
	space := tea.KeyPressMsg{Code: tea.KeySpace, Text: " "}
	m, _ := press(t, newTestModel(testItems()), text("/"), text("k"), text("p"), space, text("i"), text("n"))
	if m.query != "kp in" {
		t.Fatalf("query = %q, want %q", m.query, "kp in")
	}
	if got := selectedID(m); got != "init" || len(m.visible) != 1 {
		t.Fatalf("selected=%q visible=%d", got, len(m.visible))
	}
}

func TestSpaceOnlyTypesWhileFiltering(t *testing.T) {
	m, _ := press(t, newTestModel(testItems()), tea.KeyPressMsg{Code: tea.KeySpace, Text: " "})
	if m.filtering || m.query != "" {
		t.Fatalf("space started a filter: filtering=%v query=%q", m.filtering, m.query)
	}
}

func TestBackspaceEditsThenLeavesFilter(t *testing.T) {
	m, _ := press(t, newTestModel(testItems()), text("s"), text("h"))
	m, _ = press(t, m, key(tea.KeyBackspace, 0))
	if m.query != "s" {
		t.Fatalf("query = %q", m.query)
	}
	m, _ = press(t, m, key(tea.KeyBackspace, 0), key(tea.KeyBackspace, 0))
	if m.filtering || m.query != "" {
		t.Fatalf("filtering=%v query=%q", m.filtering, m.query)
	}
}

func TestEnterWithNoMatchesDoesNothing(t *testing.T) {
	m, quit := press(t, newTestModel(testItems()), text("z"), text("z"), key(tea.KeyEnter, 0))
	if quit || m.selected != "" || len(m.visible) != 0 {
		t.Fatalf("quit=%v selected=%q visible=%d", quit, m.selected, len(m.visible))
	}
	if !strings.Contains(m.render(), "No matches") {
		t.Fatal("empty filter result should say so")
	}
}

func TestPasteAppendsToQuery(t *testing.T) {
	m, _ := press(t, newTestModel(testItems()), tea.PasteMsg{Content: "sh\nip"})
	if m.query != "sh ip" || !m.filtering {
		t.Fatalf("query=%q filtering=%v", m.query, m.filtering)
	}
}

func TestModifiedKeysDoNotType(t *testing.T) {
	m, _ := press(t, newTestModel(testItems()), tea.KeyPressMsg{Code: 'x', Text: "x", Mod: tea.ModAlt})
	if m.query != "" {
		t.Fatalf("alt+x typed into query: %q", m.query)
	}
}
