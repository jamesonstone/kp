package cmd

import (
	"fmt"
	"strings"

	"github.com/jamesonstone/kp/internal/picker"
	"github.com/jamesonstone/kp/internal/prompt"
	"github.com/spf13/cobra"
)

func (a *app) runLauncher(cmd *cobra.Command) error {
	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}

	items := buildLauncherItems(cmd.Root().CommandPath(), reg.List())
	selection, err := a.pick(cmd, items)
	if err != nil {
		return err
	}
	return a.runLauncherSelection(cmd, selection)
}

// buildLauncherItems lists root prompts, then the few commands that make sense
// to start from the launcher. Legacy `kp v0` prompts are not listed.
func buildLauncherItems(commandPath string, prompts []prompt.Prompt) []picker.Item {
	items := make([]picker.Item, 0, len(prompts)+3)
	for _, p := range prompts {
		item := promptItem(commandPath, p)
		item.ID = "prompt:" + p.Name
		items = append(items, item)
	}

	return append(items,
		picker.Item{
			ID:      "command:init",
			Title:   "Init",
			Detail:  commandPath + " init",
			Group:   "commands",
			Preview: "Construct a coding-agent prompt.\n\nAsk for objective, context, invariants, constraints, and definition of done. On a TTY, enter continues and Shift+Enter inserts a newline. The generated prompt is printed and copied.",
		},
		picker.Item{
			ID:      "command:find-port",
			Title:   "Find port",
			Detail:  commandPath + " find-port <port>",
			Group:   "commands",
			Preview: "Inspect a port and act on the process.\n\nSearch TCP and UDP listeners on a port, inspect the matching process details, copy values, or stop the process after confirmation.",
		},
		picker.Item{
			ID:      "command:help",
			Title:   "Help",
			Detail:  commandPath + " --help",
			Group:   "commands",
			Preview: "Show all commands.\n\nEvery command, including prompt management, repo scaffolding, legacy v0 prompts, and version information.",
		},
	)
}

// promptItem renders a prompt as a picker row. The command stays secondary
// metadata in the preview rather than a competing list column.
func promptItem(commandPath string, p prompt.Prompt) picker.Item {
	detail := commandPath + " " + p.Name
	if p.Source == prompt.SourceUser {
		detail += " · user prompt"
	}
	return picker.Item{
		ID:      p.Name,
		Title:   p.Label,
		Detail:  detail,
		Group:   "prompts",
		Preview: p.Body,
	}
}

func (a *app) runLauncherSelection(cmd *cobra.Command, selection string) error {
	if name, ok := strings.CutPrefix(selection, "prompt:"); ok {
		return a.runPrompt(name)
	}

	switch selection {
	case "command:init":
		return a.runInit(false)
	case "command:find-port":
		return a.runFindPort(cmd.Context(), "")
	case "command:help":
		return cmd.Help()
	default:
		return NewExitError(ExitUser, fmt.Errorf("invalid launcher selection %q", selection))
	}
}
