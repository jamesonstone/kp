package picker

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestComputeLayout(t *testing.T) {
	tests := []struct {
		name                     string
		width, height, long, cmd int
		wantList, wantPrev       int
		wantCmd, wantBody        int
	}{
		{"fits longest title", 120, 30, 30, 0, 32, 79, 0, 25},
		{"short titles use minimum", 100, 30, 4, 0, minListWidth, 73, 0, 25},
		{"list capped just over half", 100, 30, 80, 0, 52, 39, 0, 25},
		{"preview capped for reading", 240, 30, 30, 0, 32, maxPreview, 0, 25},
		{"narrow terminal hides preview", 50, 20, 30, 0, 46, 0, 0, 15},
		{"tiny height is compact", 80, 3, 10, 0, 18, 53, 0, 1},
		{"command column fits", 120, 30, 37, 19, 61, 50, 19, 25},
		{"list-only keeps commands", 60, 20, 37, 19, 56, 0, 19, 15},
		{"tight list hides commands", 40, 20, 37, 19, 36, 0, 0, 15},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := computeLayout(tt.width, tt.height, tt.long, tt.cmd)
			if l.listWidth != tt.wantList || l.previewWidth != tt.wantPrev || l.commandWidth != tt.wantCmd || l.bodyHeight != tt.wantBody {
				t.Fatalf("layout = %+v, want list=%d preview=%d cmd=%d body=%d", l, tt.wantList, tt.wantPrev, tt.wantCmd, tt.wantBody)
			}
			if l.contentWidth() > tt.width-2*marginX {
				t.Fatalf("content width %d exceeds inner width %d", l.contentWidth(), tt.width-2*marginX)
			}
			if l.commandWidth > 0 && l.titleWidth() < min(minTitleWidth, tt.long) {
				t.Fatalf("title column %d narrower than %d", l.titleWidth(), minTitleWidth)
			}
		})
	}
}

func TestRenderFitsEveryTerminalSize(t *testing.T) {
	items := testItems()
	items[0].Title = "A very long prompt label that will not fit in a narrow list column at all"
	items[0].Preview = strings.Repeat("Long preview text with words and a verylongunbrokentokenthatmustbehardwrapped. ", 40)
	longQuery := tea.PasteMsg{Content: strings.Repeat("query ", 20)}
	for _, size := range [][2]int{{3, 2}, {10, 4}, {30, 7}, {40, 10}, {60, 16}, {80, 24}, {120, 40}, {200, 50}} {
		for _, filtered := range []bool{false, true} {
			name := fmt.Sprintf("%dx%d filtered=%v", size[0], size[1], filtered)
			m, _ := press(t, newTestModel(items), tea.WindowSizeMsg{Width: size[0], Height: size[1]})
			if filtered {
				m, _ = press(t, m, longQuery)
			}
			lines := strings.Split(m.render(), "\n")
			if len(lines) != size[1] {
				t.Fatalf("%s: rendered %d lines, want %d", name, len(lines), size[1])
			}
			for i, line := range lines {
				if w := ansi.StringWidth(line); w > size[0] {
					t.Fatalf("%s: line %d is %d cells wide: %q", name, i, w, ansi.Strip(line))
				}
			}
		}
	}
}

func TestResizeKeepsSelectionVisible(t *testing.T) {
	var items []Item
	for i := range 30 {
		items = append(items, Item{ID: fmt.Sprint(i), Title: fmt.Sprintf("Item %02d", i)})
	}
	m := newTestModel(items)
	m.jump(25)
	m, _ = press(t, m, tea.WindowSizeMsg{Width: 80, Height: 12})
	body := m.layout().bodyHeight
	if m.listTop > 25 || 25 >= m.listTop+body {
		t.Fatalf("cursor row 25 outside list window [%d,%d)", m.listTop, m.listTop+body)
	}
	if !strings.Contains(ansi.Strip(m.render()), "› Item 25") {
		t.Fatal("selected row not rendered after resize")
	}
}

func TestGroupsRenderUnderHeadings(t *testing.T) {
	m := newTestModel(testItems())
	var got []string
	for _, r := range m.listRows() {
		switch {
		case r.heading != "":
			got = append(got, r.heading)
		case r.pos < 0:
			got = append(got, "")
		default:
			got = append(got, fmt.Sprint(r.pos))
		}
	}
	want := []string{"prompts", "0", "1", "2", "", "commands", "3", "4"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("rows = %q, want %q", got, want)
	}
}

func TestScrollingUpRevealsGroupHeading(t *testing.T) {
	m := newTestModel(testItems())
	m, _ = press(t, m, tea.WindowSizeMsg{Width: 100, Height: chromeRows + 3})
	m.jump(4)
	m.jump(3)
	rows := m.listRows()
	if rows[m.listTop].heading != "commands" {
		t.Fatalf("list top row = %+v, want the commands heading", rows[m.listTop])
	}
}

func TestPreviewScrollClampsAndResetsOnMove(t *testing.T) {
	items := testItems()
	items[0].Preview = strings.Repeat("line\n", 100)
	m := newTestModel(items)
	maxTop := m.maxPreviewTop()
	if maxTop == 0 {
		t.Fatal("expected overflowing preview")
	}
	m, _ = press(t, m, key(tea.KeyPgDown, 0))
	if m.previewTop == 0 {
		t.Fatal("pgdown did not scroll")
	}
	for range 20 {
		m, _ = press(t, m, key('d', tea.ModCtrl))
	}
	if m.previewTop != maxTop {
		t.Fatalf("previewTop = %d, want clamp at %d", m.previewTop, maxTop)
	}
	if !strings.Contains(ansi.Strip(m.render()), "100%") {
		t.Fatal("scroll position should show when the preview overflows")
	}
	m, _ = press(t, m, text("j"))
	if m.previewTop != 0 {
		t.Fatalf("previewTop = %d after moving, want 0", m.previewTop)
	}
	if strings.Contains(ansi.Strip(m.render()), "%") {
		t.Fatal("scroll position should hide when the preview fits")
	}
}

func TestWrapTextKeepsHangingIndent(t *testing.T) {
	got := wrapText("- one two three four five\n  12. alpha beta gamma delta", 14)
	want := []string{
		"- one two",
		"  three four",
		"  five",
		"  12. alpha",
		"      beta",
		"      gamma",
		"      delta",
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("wrapText =\n%q\nwant\n%q", got, want)
	}
}

func TestWrapTextAddsNoSpaceAtHyphenBreaks(t *testing.T) {
	got := wrapText("- alpha beta-gamma delta epsilon", 16)
	want := []string{"- alpha beta-", "  gamma delta", "  epsilon"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("wrapText = %q, want %q", got, want)
	}
}

func TestLongQueryKeepsItsTypedEndVisible(t *testing.T) {
	m, _ := press(t, newTestModel(testItems()), tea.PasteMsg{Content: strings.Repeat("x", 200) + "tail"})
	header := ansi.Strip(m.header(m.layout().contentWidth()))
	if !strings.Contains(header, "tail") || !strings.Contains(header, ellipsis) {
		t.Fatalf("header = %q", header)
	}
}

func TestWrapTextHardWrapsLongTokens(t *testing.T) {
	for _, line := range wrapText(strings.Repeat("x", 50), 12) {
		if ansi.StringWidth(line) > 12 {
			t.Fatalf("line %q exceeds width", line)
		}
	}
}

func TestSanitizeDropsControlSequences(t *testing.T) {
	got := sanitize("a\tb\r\nc\x1b[31mred\x07")
	if got != "a    b\nc[31mred" {
		t.Fatalf("sanitize = %q", got)
	}
	m := newTestModel([]Item{{ID: "x", Title: "Bad\x1b]0;title\x07\nlabel"}})
	if m.items[0].Title != "Bad]0;title label" {
		t.Fatalf("title = %q", m.items[0].Title)
	}
}

func TestFooterDropsWholeHintsWhenNarrow(t *testing.T) {
	m := newTestModel(testItems())
	footer := ansi.Strip(m.footer(24))
	if strings.Contains(footer, ellipsis) || ansi.StringWidth(footer) > 24 {
		t.Fatalf("footer = %q", footer)
	}
	if !strings.HasPrefix(footer, "↑↓ move") {
		t.Fatalf("footer lost its first hint: %q", footer)
	}
}

func TestRenderShowsPointerHeadingsAndAlignedCommands(t *testing.T) {
	m := newTestModel(testItems())
	out := ansi.Strip(m.render())
	for _, want := range []string{" kp ", "prompts", "commands", "› Merge", "kp merge", "merge body", "↑↓ move", "enter select", "esc quit"} {
		if !strings.Contains(out, want) {
			t.Fatalf("render missing %q:\n%s", want, out)
		}
	}
	col := -1
	for _, line := range strings.Split(out, "\n") {
		if !strings.Contains(line, "│") {
			continue
		}
		list, _, _ := strings.Cut(line, "│")
		if i := strings.Index(list, "kp "); i >= 0 {
			i = ansi.StringWidth(list[:i])
			if col >= 0 && i != col {
				t.Fatalf("command column misaligned at %d, want %d: %q", i, col, line)
			}
			col = i
		}
	}
	if col < 0 {
		t.Fatal("list shows no command column")
	}
}
