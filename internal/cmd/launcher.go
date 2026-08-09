package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/jamesonstone/kp/internal/prompt"
	"github.com/spf13/cobra"
)

type LauncherItem struct {
	ID          string
	Emoji       string
	Title       string
	Command     string
	Description string
	Preview     string
}

func (a *app) runLauncher(cmd *cobra.Command) error {
	if a.launcherRunner == nil {
		if _, err := a.lookPath("fzf"); err != nil {
			return NewExitError(ExitConfig, errors.New("fzf not found; install fzf via 'brew install fzf' or run 'kp --help'"))
		}
	}

	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}

	items := buildLauncherItems(cmd.CommandPath(), reg.List())
	selection := ""
	if a.launcherRunner != nil {
		selection, err = a.launcherRunner(items)
	} else {
		selection, err = runLauncherFZF(items)
	}
	if err != nil {
		return mapPickerError(err)
	}

	if !launcherHasSelection(items, selection) {
		return NewExitError(ExitUser, fmt.Errorf("invalid launcher selection %q", selection))
	}
	return a.runLauncherSelection(cmd, selection)
}

func buildLauncherItems(commandPath string, prompts []prompt.Prompt) []LauncherItem {
	items := make([]LauncherItem, 0, len(prompts)+2)
	for _, p := range prompts {
		items = append(items, LauncherItem{
			ID:          "prompt:" + p.Name,
			Emoji:       launcherPromptEmoji(p),
			Title:       p.Label,
			Command:     commandPath + " " + p.Name,
			Description: "Print and copy prompt",
			Preview:     p.Body,
		})
	}

	items = append(items,
		LauncherItem{
			ID:          "command:find-port",
			Emoji:       "🔍",
			Title:       "Find port",
			Command:     commandPath + " find-port <port>",
			Description: "Inspect a port and act on the process",
			Preview:     "Search TCP and UDP listeners on a port, inspect the matching process details, copy values, or stop the process after confirmation.",
		},
		LauncherItem{
			ID:          "command:help",
			Emoji:       "❓",
			Title:       "Help",
			Command:     commandPath + " --help",
			Description: "Show all commands",
			Preview:     "Show every command, including prompt management, repo scaffolding, and version information.",
		},
	)
	return items
}

func (a *app) runLauncherSelection(cmd *cobra.Command, selection string) error {
	if name, ok := strings.CutPrefix(selection, "prompt:"); ok {
		return a.runPrompt(name)
	}

	switch selection {
	case "command:find-port":
		return a.runFindPort("")
	case "command:help":
		return cmd.Help()
	default:
		return NewExitError(ExitUser, fmt.Errorf("invalid launcher selection %q", selection))
	}
}

func launcherHasSelection(items []LauncherItem, selection string) bool {
	for _, item := range items {
		if item.ID == selection {
			return true
		}
	}
	return false
}

func runLauncherFZF(items []LauncherItem) (string, error) {
	previewDir, err := os.MkdirTemp("", "kp-launcher-preview-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(previewDir)

	var input strings.Builder
	displayRows := launcherDisplayRows(items)
	for _, item := range items {
		if err := os.WriteFile(filepath.Join(previewDir, item.ID), []byte(item.Preview), 0o600); err != nil {
			return "", err
		}
		fmt.Fprintf(&input, "%s\t%s\t%s\t%s\t%s\n", item.ID, displayRows[item.ID], item.Title, item.Command, item.Description)
	}

	cmd := exec.Command("fzf", launcherFZFArgs(previewDir)...)
	cmd.Stdin = strings.NewReader(input.String())
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %v", errPickerCanceled, err)
	}

	selected := strings.TrimSpace(stdout.String())
	if selected == "" {
		return "", errPickerCanceled
	}
	id, _, _ := strings.Cut(selected, "\t")
	if id == "" {
		return "", fmt.Errorf("invalid launcher selection")
	}
	return id, nil
}

func launcherFZFArgs(previewDir string) []string {
	return []string{
		"--height", "70%",
		"--reverse",
		"--cycle",
		"--info", "hidden",
		"--no-separator",
		"--no-hscroll",
		"--prompt", "🎛️ kp › ",
		"--pointer", "👉",
		"--header", "j/k or arrows: move · enter: select · esc: close",
		"--delimiter", "\t",
		"--with-nth", "2",
		"--nth", "2,3,4,5",
		"--bind", "j:down,k:up,tab:down,btab:up",
		"--preview", "cat " + shellQuote(previewDir) + "/{1}",
		"--preview-window", "right,55%,wrap",
	}
}

func launcherDisplayRows(items []LauncherItem) map[string]string {
	titleWidth := displayWidth("Item")
	for _, item := range items {
		titleWidth = max(titleWidth, displayWidth(item.Title))
	}

	rows := make(map[string]string, len(items))
	for _, item := range items {
		rows[item.ID] = fmt.Sprintf(
			"%s  %s  %s",
			padDisplay(item.Emoji, 2),
			padDisplay(item.Title, titleWidth),
			item.Command,
		)
	}
	return rows
}

func padDisplay(value string, width int) string {
	padding := width - displayWidth(value)
	if padding <= 0 {
		return value
	}
	return value + strings.Repeat(" ", padding)
}

func displayWidth(value string) int {
	width := 0
	for _, r := range value {
		switch {
		case r == '\uFE0F':
			continue
		case unicode.Is(unicode.Mn, r):
			continue
		case r < 0x20:
			continue
		case r >= 0x1F300 && r <= 0x1FAFF:
			width += 2
		case r >= 0x2600 && r <= 0x27BF:
			width += 2
		default:
			width++
		}
	}
	return width
}

func launcherPromptEmoji(p prompt.Prompt) string {
	if p.Source == prompt.SourceUser {
		return "📝"
	}
	return "🧠"
}
