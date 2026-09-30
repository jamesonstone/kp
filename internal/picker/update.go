package picker

import (
	tea "charm.land/bubbletea/v2"
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
	case tea.KeyPressMsg:
		if m.handleKey(msg) {
			return m, tea.Quit
		}
	case tea.PasteMsg:
		m.appendQuery(oneLine(msg.Content))
	}
	return m, nil
}

// handleKey applies one key press and reports whether the picker is done.
//
// Outside filter mode j/k move and any other printable key starts a filter,
// matching the previous fzf launcher where typing narrowed the list. In filter
// mode every printable key edits the query.
func (m *model) handleKey(msg tea.KeyPressMsg) bool {
	half := max(m.layout().previewBodyHeight()/2, 1)

	switch msg.String() {
	case "ctrl+c":
		m.quitting = true
	case "esc":
		if m.filtering {
			m.clearFilter()
			return false
		}
		m.quitting = true
	case "enter":
		item, ok := m.current()
		if !ok {
			return false
		}
		m.selected = item.ID
		m.quitting = true
	case "up", "shift+tab", "ctrl+p":
		m.move(-1)
	case "down", "tab", "ctrl+n":
		m.move(1)
	case "home":
		m.jump(0)
	case "end":
		m.jump(len(m.visible) - 1)
	case "pgdown", "ctrl+d":
		m.scrollPreview(half)
	case "pgup", "ctrl+u":
		m.scrollPreview(-half)
	case "shift+down":
		m.scrollPreview(1)
	case "shift+up":
		m.scrollPreview(-1)
	case "backspace":
		m.backspace()
	case "space":
		if m.filtering {
			m.appendQuery(" ")
		}
	default:
		if msg.Text == "" || msg.Mod&(tea.ModCtrl|tea.ModAlt|tea.ModMeta|tea.ModSuper) != 0 {
			return false
		}
		m.typeRunes(msg.Text)
	}
	return m.quitting
}

func (m *model) typeRunes(text string) {
	if m.filtering || len(text) != 1 {
		m.appendQuery(text)
		return
	}
	switch text {
	case "j":
		m.move(1)
	case "k":
		m.move(-1)
	case "/":
		m.filtering = true
	default:
		m.appendQuery(text)
	}
}
