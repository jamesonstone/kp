// Package picker is kp's two-pane interactive selector: a list on the left
// and a live preview on the right. It owns rendering, keys, scrolling, and
// resizing only; callers supply items and act on the selected ID.
package picker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Item is one selectable row.
type Item struct {
	ID      string // returned on selection
	Title   string // primary label in the list
	Command string // how to run it directly, e.g. "kp review"; shown in the list and preview
	Note    string // optional tag shown in the preview, e.g. "user prompt"
	Group   string // section heading; consecutive items share one heading
	Preview string // plain text shown in the preview pane
}

var (
	// ErrCanceled reports that the user left without selecting.
	ErrCanceled = errors.New("picker canceled")
	// ErrNoTerminal reports that no controlling terminal is available.
	ErrNoTerminal = errors.New("interactive picker requires a terminal")
)

// Run shows the picker on the controlling terminal and returns the selected
// item's ID. It draws on /dev/tty rather than stdout so stdout stays clean
// for the selected command's output.
func Run(ctx context.Context, title string, items []Item) (string, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrNoTerminal, err)
	}
	defer tty.Close()

	program := tea.NewProgram(newModel(title, items, newStyles()),
		tea.WithContext(ctx),
		tea.WithInput(tty),
		tea.WithOutput(tty),
	)
	final, err := program.Run()
	if err != nil {
		if errors.Is(err, tea.ErrProgramKilled) || errors.Is(err, tea.ErrInterrupted) {
			return "", fmt.Errorf("%w: %v", ErrCanceled, err)
		}
		return "", err
	}
	if result, ok := final.(model); ok && result.selected != "" {
		return result.selected, nil
	}
	return "", ErrCanceled
}

type styles struct {
	badge       lipgloss.Style // "kp" title badge
	pointer     lipgloss.Style // selection marker
	selected    lipgloss.Style // selected title
	command     lipgloss.Style // command name, e.g. "review" in "kp review"
	commandDim  lipgloss.Style // command prefix and arguments
	selectedCmd lipgloss.Style
	section     lipgloss.Style // group headings
	previewHead lipgloss.Style // preview title
	note        lipgloss.Style // tags, counts, scroll position
	key         lipgloss.Style // footer keys
	faint       lipgloss.Style
	cursor      lipgloss.Style
	mdHeading   lipgloss.Style
	mdMarker    lipgloss.Style
	mdCode      lipgloss.Style
}

// newStyles uses only the terminal's basic ANSI palette so every theme
// shades these colors for its own background, light or dark: magenta for the
// brand and selection, cyan for commands and keys, blue for structure, yellow
// for tags. Bubble Tea drops color entirely under NO_COLOR.
func newStyles() styles {
	style := lipgloss.NewStyle
	return styles{
		badge:       style().Bold(true).Reverse(true).Foreground(lipgloss.Magenta),
		pointer:     style().Bold(true).Foreground(lipgloss.Magenta),
		selected:    style().Bold(true).Foreground(lipgloss.Magenta),
		command:     style().Foreground(lipgloss.Cyan),
		commandDim:  style().Faint(true),
		selectedCmd: style().Bold(true).Foreground(lipgloss.Cyan),
		section:     style().Bold(true).Foreground(lipgloss.Blue),
		previewHead: style().Bold(true).Foreground(lipgloss.Magenta),
		note:        style().Foreground(lipgloss.Yellow),
		key:         style().Bold(true).Foreground(lipgloss.Cyan),
		faint:       style().Faint(true),
		cursor:      style().Reverse(true),
		mdHeading:   style().Bold(true).Foreground(lipgloss.Blue),
		mdMarker:    style().Foreground(lipgloss.Magenta),
		mdCode:      style().Foreground(lipgloss.Cyan),
	}
}

func cleanItems(items []Item) []Item {
	cleaned := make([]Item, len(items))
	for i, item := range items {
		item.Title = oneLine(item.Title)
		item.Command = oneLine(item.Command)
		item.Note = oneLine(item.Note)
		item.Group = oneLine(item.Group)
		cleaned[i] = item
	}
	return cleaned
}

func oneLine(text string) string {
	return strings.Join(strings.Fields(sanitize(text)), " ")
}
