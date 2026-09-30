package picker

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

const (
	pointer  = "›"
	ellipsis = "…"
)

func (m model) View() tea.View {
	view := tea.NewView(m.render())
	view.AltScreen = true
	return view
}

func (m model) render() string {
	if m.quitting {
		return ""
	}

	l := m.layout()
	margin := strings.Repeat(" ", marginX)
	list := m.listLines(l)
	preview := m.previewLines(l)
	divider := "  " + m.st.faint.Render("│") + "  "

	var lines []string
	if !l.compact {
		lines = append(lines, "")
	}
	lines = append(lines, margin+m.header(l.contentWidth()))
	if !l.compact {
		lines = append(lines, "")
	}
	for row := range l.bodyHeight {
		line := margin + padRight(list[row], l.listWidth)
		if l.previewWidth > 0 {
			line += divider + preview[row]
		}
		lines = append(lines, line)
	}
	if !l.compact {
		lines = append(lines, "")
	}
	lines = append(lines, margin+m.footer(l.contentWidth()))

	// Guard degenerate sizes: never draw past the terminal edges.
	lines = lines[:min(len(lines), max(m.height, 1))]
	for i, line := range lines {
		lines[i] = ansi.Truncate(line, m.width, "")
	}
	return strings.Join(lines, "\n")
}

func (m model) header(width int) string {
	left := m.st.badge.Render(" " + m.title + " ")
	if !m.filtering {
		return left
	}
	right := m.st.note.Render(fmt.Sprintf("%d/%d", len(m.visible), len(m.items)))
	prefix := left + "  " + m.st.key.Render("/") + " "
	// Keep the end of a long query, where the user is typing, in view.
	room := width - ansi.StringWidth(prefix) - ansi.StringWidth(right) - 2
	query := ansi.TruncateLeft(m.query, max(ansi.StringWidth(m.query)-room, 0), ellipsis)
	return spread(prefix+query+m.st.cursor.Render(" "), right, width)
}

func (m model) listLines(l layout) []string {
	lines := make([]string, l.bodyHeight)
	if len(m.visible) == 0 {
		lines[0] = "  " + m.st.faint.Render("No matches")
		return lines
	}

	rows := m.listRows()
	for row := range lines {
		i := m.listTop + row
		if i >= len(rows) {
			break
		}
		r := rows[i]
		switch {
		case r.heading != "":
			lines[row] = "  " + m.st.section.Render(ansi.Truncate(r.heading, l.listWidth-2, ellipsis))
		case r.pos >= 0:
			lines[row] = m.itemLine(l, m.items[m.visible[r.pos]], r.pos == m.cursor)
		}
	}
	return lines
}

// itemLine renders one list entry: pointer, title, and the aligned command.
func (m model) itemLine(l layout, item Item, selected bool) string {
	title := ansi.Truncate(item.Title, l.titleWidth(), ellipsis)
	lead := "  "
	if selected {
		lead = m.st.pointer.Render(pointer) + " "
		title = m.st.selected.Render(title)
	}
	if l.commandWidth == 0 || item.Command == "" {
		return lead + title
	}
	return lead + padRight(title, l.titleWidth()+columnGap) + m.renderCommand(item.Command, selected)
}

// renderCommand highlights the part people recognize, the subcommand name,
// and dims the "kp" prefix and any argument placeholders.
func (m model) renderCommand(command string, selected bool) string {
	fields := strings.Fields(command)
	if len(fields) < 2 {
		return m.st.command.Render(command)
	}
	name := m.st.command
	if selected {
		name = m.st.selectedCmd
	}
	out := m.st.commandDim.Render(fields[0]) + " " + name.Render(fields[1])
	if len(fields) > 2 {
		out += " " + m.st.commandDim.Render(strings.Join(fields[2:], " "))
	}
	return out
}

func (m model) previewLines(l layout) []string {
	lines := make([]string, l.bodyHeight)
	item, ok := m.current()
	if !ok || l.previewWidth == 0 {
		return lines
	}

	meta := m.renderCommand(item.Command, false)
	if item.Note != "" {
		meta += m.st.faint.Render("  ·  ") + m.st.note.Render(item.Note)
	}
	header := []string{
		m.st.previewHead.Render(ansi.Truncate(item.Title, l.previewWidth, ellipsis)),
		ansi.Truncate(meta, l.previewWidth, ellipsis),
		"",
	}
	copy(lines, header)

	body := tintLines(wrapText(item.Preview, l.previewWidth), m.st)
	for row := previewHeader; row < l.bodyHeight; row++ {
		i := m.previewTop + row - previewHeader
		if i >= len(body) {
			break
		}
		lines[row] = body[i]
	}
	return lines
}

type hint struct{ key, label string }

// footer shows key hints in quiet type. The scroll hint and position only
// appear when the preview overflows. When space runs out, whole hints drop
// from the end instead of being cut mid-word.
func (m model) footer(width int) string {
	hints := []hint{{"↑↓", "move"}, {"enter", "select"}, {"esc", "quit"}, {"/", "filter"}}
	if m.filtering {
		hints = []hint{{"↑↓", "move"}, {"enter", "select"}, {"esc", "clear"}}
	}

	right := ""
	if maxTop := m.maxPreviewTop(); maxTop > 0 {
		hints = append(hints, hint{"ctrl-d/u", "scroll"})
		right = m.st.note.Render(fmt.Sprintf("%d%%", m.previewTop*100/maxTop))
	}

	room := width - ansi.StringWidth(right) - 2
	left := m.renderHints(hints)
	for len(hints) > 1 && ansi.StringWidth(left) > room {
		hints = hints[:len(hints)-1]
		left = m.renderHints(hints)
	}
	return spread(left, right, width)
}

func (m model) renderHints(hints []hint) string {
	parts := make([]string, len(hints))
	for i, h := range hints {
		parts[i] = m.st.key.Render(h.key) + " " + m.st.faint.Render(h.label)
	}
	return strings.Join(parts, "   ")
}

// spread places left and right at the edges of width.
func spread(left, right string, width int) string {
	if right == "" {
		return left
	}
	gap := width - ansi.StringWidth(left) - ansi.StringWidth(right)
	return left + strings.Repeat(" ", max(gap, 1)) + right
}

func padRight(text string, width int) string {
	return text + strings.Repeat(" ", max(width-ansi.StringWidth(text), 0))
}
