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
	items := make([]LauncherItem, 0, len(prompts)+7)
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
			ID:          "command:list",
			Emoji:       "📚",
			Title:       "Prompt picker",
			Command:     commandPath + " list",
			Description: "Browse prompts only",
			Preview:     "Open the prompt-only picker. Selecting a prompt prints it and copies it to the clipboard after verification.",
		},
		LauncherItem{
			ID:          "command:new",
			Emoji:       "✨",
			Title:       "New prompt",
			Command:     commandPath + " new <name>",
			Description: "Show creation help",
			Preview:     "Create a user prompt file and open it in the configured editor.",
		},
		LauncherItem{
			ID:          "command:edit",
			Emoji:       "✏️",
			Title:       "Edit prompt",
			Command:     commandPath + " edit <name>",
			Description: "Show edit help",
			Preview:     "Edit a user prompt. Built-ins are promoted to the user prompt directory before editing.",
		},
		LauncherItem{
			ID:          "command:rm",
			Emoji:       "🗑️",
			Title:       "Remove prompt",
			Command:     commandPath + " rm <name>",
			Description: "Show removal help",
			Preview:     "Remove a user prompt. Built-in prompts cannot be deleted.",
		},
		LauncherItem{
			ID:          "command:scaffold",
			Emoji:       "🏗️",
			Title:       "Repo scaffold",
			Command:     commandPath + " scaffold",
			Description: "Show scaffold help",
			Preview:     "Create repo support files such as agent instructions, review config, PR template, and local env files.",
		},
		LauncherItem{
			ID:          "command:version",
			Emoji:       "🛠️",
			Title:       "Version",
			Command:     commandPath + " --version",
			Description: "Print version metadata",
			Preview:     "Print the current version and commit metadata.",
		},
		LauncherItem{
			ID:          "command:help",
			Emoji:       "❓",
			Title:       "Help",
			Command:     commandPath + " --help",
			Description: "Show help menu",
			Preview:     "Show the non-interactive grouped help menu.",
		},
	)
	return items
}

func (a *app) runLauncherSelection(cmd *cobra.Command, selection string) error {
	if name, ok := strings.CutPrefix(selection, "prompt:"); ok {
		return a.runPrompt(name)
	}

	switch selection {
	case "command:list":
		return a.runPicker()
	case "command:new":
		return showCommandHelp(cmd, "new")
	case "command:edit":
		return showCommandHelp(cmd, "edit")
	case "command:rm":
		return showCommandHelp(cmd, "rm")
	case "command:scaffold":
		return showCommandHelp(cmd, "scaffold")
	case "command:version":
		_, err := fmt.Fprintln(a.stdout, cmd.Version)
		return err
	case "command:help":
		return cmd.Help()
	default:
		return NewExitError(ExitUser, fmt.Errorf("invalid launcher selection %q", selection))
	}
}

func showCommandHelp(root *cobra.Command, name string) error {
	cmd, _, err := root.Find([]string{name})
	if err != nil {
		return NewExitError(ExitUser, err)
	}
	return cmd.Help()
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

	cmd := exec.Command(
		"fzf",
		"--height", "70%",
		"--reverse",
		"--cycle",
		"--prompt", "🎛️ kp › ",
		"--pointer", "👉",
		"--header", "enter: select · tab/shift-tab/arrows: move · esc: cancel",
		"--delimiter", "\t",
		"--with-nth", "2",
		"--nth", "2,3,4,5",
		"--bind", "tab:down,btab:up",
		"--preview", "cat "+shellQuote(previewDir)+"/{1}",
	)
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

func launcherDisplayRows(items []LauncherItem) map[string]string {
	titleWidth := displayWidth("Item")
	commandWidth := displayWidth("Command")
	for _, item := range items {
		titleWidth = max(titleWidth, displayWidth(item.Title))
		commandWidth = max(commandWidth, displayWidth(item.Command))
	}

	rows := make(map[string]string, len(items))
	for _, item := range items {
		rows[item.ID] = fmt.Sprintf(
			"%s  %s  %s  %s",
			padDisplay(item.Emoji, 2),
			padDisplay(item.Title, titleWidth),
			padDisplay(item.Command, commandWidth),
			item.Description,
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
