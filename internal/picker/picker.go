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
	Detail  string // secondary metadata shown under the preview title
	Group   string // consecutive items with different groups get a spacer row
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
	brand    lipgloss.Style
	pointer  lipgloss.Style
	selected lipgloss.Style
	heading  lipgloss.Style
	faint    lipgloss.Style
	cursor   lipgloss.Style
}

// newStyles uses one basic ANSI accent so the terminal theme picks a shade
// that suits its own background; everything else is bold, faint, or default.
// Bubble Tea downsamples or drops color for the detected terminal profile.
func newStyles() styles {
	accent := lipgloss.Magenta
	return styles{
		brand:    lipgloss.NewStyle().Bold(true).Foreground(accent),
		pointer:  lipgloss.NewStyle().Foreground(accent),
		selected: lipgloss.NewStyle().Bold(true).Foreground(accent),
		heading:  lipgloss.NewStyle().Bold(true),
		faint:    lipgloss.NewStyle().Faint(true),
		cursor:   lipgloss.NewStyle().Reverse(true),
	}
}

func cleanItems(items []Item) []Item {
	cleaned := make([]Item, len(items))
	for i, item := range items {
		item.Title = oneLine(item.Title)
		item.Detail = oneLine(item.Detail)
		cleaned[i] = item
	}
	return cleaned
}

func oneLine(text string) string {
	return strings.Join(strings.Fields(sanitize(text)), " ")
}
