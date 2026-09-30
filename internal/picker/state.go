package picker

import (
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/x/ansi"
)

// model holds all picker state. Transitions are plain methods so navigation,
// filtering, and scrolling can be tested without a terminal.
type model struct {
	title string
	items []Item
	st    styles

	longestTitle int // display width of the longest title, for layout

	visible    []int // indices into items that match the query
	cursor     int   // position within visible
	listTop    int   // first list row in view
	previewTop int   // first preview body line in view
	query      string
	filtering  bool

	width  int
	height int

	selected string
	quitting bool
}

func newModel(title string, items []Item, st styles) model {
	m := model{title: title, items: cleanItems(items), st: st, width: 80, height: 24}
	for _, item := range m.items {
		m.longestTitle = max(m.longestTitle, ansi.StringWidth(item.Title))
	}
	m.applyFilter()
	return m
}

func (m *model) current() (Item, bool) {
	if len(m.visible) == 0 {
		return Item{}, false
	}
	return m.items[m.visible[m.cursor]], true
}

// move shifts the selection by delta, wrapping at both ends.
func (m *model) move(delta int) {
	n := len(m.visible)
	if n == 0 {
		return
	}
	m.cursor = ((m.cursor+delta)%n + n) % n
	m.previewTop = 0
	m.clampList()
}

// jump selects an absolute position, clamped to the visible range.
func (m *model) jump(pos int) {
	if len(m.visible) == 0 {
		return
	}
	m.cursor = min(max(pos, 0), len(m.visible)-1)
	m.previewTop = 0
	m.clampList()
}

// setQuery refilters and keeps the current item selected when it still matches.
func (m *model) setQuery(query string) {
	keep, hadCurrent := m.current()
	m.query = query
	m.applyFilter()
	m.cursor = 0
	if hadCurrent {
		for pos, idx := range m.visible {
			if m.items[idx].ID == keep.ID {
				m.cursor = pos
				break
			}
		}
	}
	m.previewTop = 0
	m.listTop = 0
	m.clampList()
}

func (m *model) appendQuery(text string) {
	m.filtering = true
	m.setQuery(m.query + text)
}

func (m *model) backspace() {
	if m.query == "" {
		m.filtering = false
		return
	}
	_, size := utf8.DecodeLastRuneInString(m.query)
	m.setQuery(m.query[:len(m.query)-size])
}

func (m *model) clearFilter() {
	m.filtering = false
	m.setQuery("")
}

func (m *model) applyFilter() {
	terms := strings.Fields(strings.ToLower(m.query))
	m.visible = make([]int, 0, len(m.items))
	for i, item := range m.items {
		if matches(item, terms) {
			m.visible = append(m.visible, i)
		}
	}
}

// matches reports whether every term appears in the item's title or detail.
// Substring terms keep results predictable and in their original order.
func matches(item Item, terms []string) bool {
	haystack := strings.ToLower(item.Title + " " + item.Detail)
	for _, term := range terms {
		if !strings.Contains(haystack, term) {
			return false
		}
	}
	return true
}

// listRows maps list rows to visible positions. A -1 row is a spacer between
// item groups.
func (m *model) listRows() []int {
	rows := make([]int, 0, len(m.visible)+2)
	for pos, idx := range m.visible {
		if pos > 0 && m.items[idx].Group != m.items[m.visible[pos-1]].Group {
			rows = append(rows, -1)
		}
		rows = append(rows, pos)
	}
	return rows
}

// clampList scrolls the list so the selected row stays in view.
func (m *model) clampList() {
	height := m.layout().bodyHeight
	rows := m.listRows()
	row := 0
	for i, pos := range rows {
		if pos == m.cursor {
			row = i
			break
		}
	}
	if row < m.listTop {
		m.listTop = row
	}
	if row >= m.listTop+height {
		m.listTop = row - height + 1
	}
	m.listTop = min(max(m.listTop, 0), max(len(rows)-height, 0))
}

// scrollPreview moves the preview body by delta lines within its bounds.
func (m *model) scrollPreview(delta int) {
	m.previewTop = min(max(m.previewTop+delta, 0), m.maxPreviewTop())
}

func (m *model) maxPreviewTop() int {
	l := m.layout()
	item, ok := m.current()
	if !ok || l.previewWidth == 0 {
		return 0
	}
	return max(len(wrapText(item.Preview, l.previewWidth))-l.previewBodyHeight(), 0)
}

func (m *model) resize(width, height int) {
	m.width, m.height = width, height
	m.clampList()
	m.scrollPreview(0)
}
