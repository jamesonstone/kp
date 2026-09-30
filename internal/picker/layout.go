package picker

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

const (
	marginX       = 2  // blank columns at the left and right edges
	chromeRows    = 5  // top pad, header, blank, blank, footer
	compactRows   = 8  // below this height only header and footer remain
	dividerWidth  = 5  // "  │  "
	minListWidth  = 18 // narrowest useful list column
	minTitleWidth = 14 // below this the command column is hidden
	columnGap     = 3  // between title and command in the list
	minPreview    = 28 // below this the preview is hidden
	maxPreview    = 88 // comfortable reading measure on wide terminals
	previewHeader = 3  // title, detail, blank
	tabWidth      = 4
)

// layout is the pane geometry for one terminal size.
type layout struct {
	compact      bool // no padding rows on very short terminals
	listWidth    int
	commandWidth int // 0 hides the command column
	previewWidth int // 0 hides the preview
	bodyHeight   int
}

// contentWidth spans the panes actually drawn, so header and footer align
// with the content instead of the far terminal edge.
func (l layout) titleWidth() int {
	if l.commandWidth == 0 {
		return max(l.listWidth-2, 0)
	}
	return l.listWidth - 2 - columnGap - l.commandWidth
}

func (l layout) contentWidth() int {
	if l.previewWidth == 0 {
		return l.listWidth
	}
	return l.listWidth + dividerWidth + l.previewWidth
}

func (l layout) previewBodyHeight() int {
	return max(l.bodyHeight-previewHeader, 1)
}

// computeLayout sizes the list to its longest title plus the command column,
// capped at just over half the width, and gives the rest to the preview.
// Narrow terminals drop the preview rather than squeezing both panes, and
// drop the command column before titles get too short to read.
func computeLayout(width, height, longestTitle, longestCommand int) layout {
	inner := max(width-2*marginX, 1)
	l := layout{bodyHeight: max(height-chromeRows, 1)}
	if height < compactRows {
		l.compact = true
		l.bodyHeight = max(height-2, 1)
	}

	want := longestTitle + 2
	if longestCommand > 0 {
		want += columnGap + longestCommand
	}
	list := min(max(want, minListWidth), inner*11/20)
	preview := inner - list - dividerWidth
	if preview < minPreview {
		list = inner
	} else {
		l.previewWidth = min(preview, maxPreview)
	}
	l.listWidth = list
	if longestCommand > 0 && list-2-columnGap-longestCommand >= min(minTitleWidth, longestTitle) {
		l.commandWidth = longestCommand
	}
	return l
}

func (m *model) layout() layout {
	return computeLayout(m.width, m.height, m.longestTitle, m.longestCommand)
}

// wrapText turns preview content into display lines no wider than width.
// Continuation lines keep the indentation of list items and indented text.
func wrapText(text string, width int) []string {
	if width <= 0 {
		return nil
	}
	var out []string
	for _, line := range strings.Split(sanitize(text), "\n") {
		line = strings.TrimRight(line, " ")
		if ansi.StringWidth(line) <= width {
			out = append(out, line)
			continue
		}
		indent := hangingIndent(line)
		if indent >= width/2 {
			indent = 0
		}
		wrapped := ansi.Wrap(line, width, "")
		first, others, _ := strings.Cut(wrapped, "\n")
		out = append(out, first)
		// Rewrap the untouched remainder so hyphen or hard breaks gain no space.
		rest := strings.ReplaceAll(others, "\n", " ")
		if strings.HasPrefix(line, first) {
			rest = line[len(first):]
		}
		rest = strings.TrimLeft(rest, " ")
		if rest == "" {
			continue
		}
		pad := strings.Repeat(" ", indent)
		for _, wrapped := range strings.Split(ansi.Wrap(rest, width-indent, ""), "\n") {
			out = append(out, pad+wrapped)
		}
	}
	for len(out) > 0 && out[len(out)-1] == "" {
		out = out[:len(out)-1]
	}
	return out
}

// hangingIndent returns the column where wrapped text should resume: after
// leading spaces and any list marker such as "- ", "* ", or "12. ".
func hangingIndent(line string) int {
	trimmed := strings.TrimLeft(line, " ")
	indent := len(line) - len(trimmed)
	for _, marker := range []string{"- ", "* ", "+ ", "> "} {
		if strings.HasPrefix(trimmed, marker) {
			return indent + len(marker)
		}
	}
	digits := 0
	for digits < len(trimmed) && trimmed[digits] >= '0' && trimmed[digits] <= '9' {
		digits++
	}
	if digits > 0 && strings.HasPrefix(trimmed[digits:], ". ") {
		return indent + digits + 2
	}
	return indent
}

// sanitize expands tabs and drops control characters so user prompt files
// cannot inject terminal escape sequences into the picker.
func sanitize(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	var b strings.Builder
	b.Grow(len(text))
	for _, r := range text {
		switch {
		case r == '\n':
			b.WriteRune(r)
		case r == '\t':
			b.WriteString(strings.Repeat(" ", tabWidth))
		case r < 0x20 || r == 0x7f || (r >= 0x80 && r < 0xa0):
			continue
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
