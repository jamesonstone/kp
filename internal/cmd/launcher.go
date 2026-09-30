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

// buildLauncherItems lists the prompts, including the constructed task
// prompt in name order, then the few tool commands that make sense to start
// from the launcher. Legacy `kp v0` prompts are not listed.
func buildLauncherItems(commandPath string, prompts []prompt.Prompt) []picker.Item {
	initItem := picker.Item{
		ID:      "command:task",
		Title:   "Start a task from scratch",
		Command: commandPath + " task",
		Group:   "prompts",
		Preview: "Construct a prompt for a new conversation, or a new task inside an existing thread.\n\nOpens the template in your editor ($KP_EDITOR, $EDITOR, nvim, or vi) with sections for:\n\n- Objective\n- Known business/domain context\n- Invariants\n- Constraints\n- Definition of done\n\nWrite as much as you like under each heading, then save and quit. The finished prompt is printed and copied. An empty file cancels.",
	}

	items := make([]picker.Item, 0, len(prompts)+3)
	for _, p := range prompts {
		if initItem.ID != "" && p.Name > "task" {
			items = append(items, initItem)
			initItem.ID = ""
		}
		item := promptItem(commandPath, p)
		item.ID = "prompt:" + p.Name
		items = append(items, item)
	}
	if initItem.ID != "" {
		items = append(items, initItem)
	}

	return append(items,
		picker.Item{
			ID:      "command:find-port",
			Title:   "Find port",
			Command: commandPath + " find-port <port>",
			Group:   "commands",
			Preview: "Inspect a port and act on the process.\n\nSearch TCP and UDP listeners on a port, inspect the matching process details, copy values, or stop the process after confirmation.",
		},
		picker.Item{
			ID:      "command:help",
			Title:   "Help",
			Command: commandPath + " --help",
			Group:   "commands",
			Preview: "Show all commands.\n\nEvery command, including prompt management, repo scaffolding, legacy v0 prompts, and version information.",
		},
	)
}

// promptItem renders a prompt as a picker row with its direct command.
func promptItem(commandPath string, p prompt.Prompt) picker.Item {
	item := picker.Item{
		ID:      p.Name,
		Title:   p.Label,
		Command: commandPath + " " + p.Name,
		Group:   "prompts",
		Preview: p.Body,
	}
	if p.Source == prompt.SourceUser {
		item.Note = "user prompt"
	}
	return item
}

func (a *app) runLauncherSelection(cmd *cobra.Command, selection string) error {
	if name, ok := strings.CutPrefix(selection, "prompt:"); ok {
		return a.runPrompt(name)
	}

	switch selection {
	case "command:task":
		return a.runInit(false)
	case "command:find-port":
		return a.runFindPort(cmd.Context(), "")
	case "command:help":
		return cmd.Help()
	default:
		return NewExitError(ExitUser, fmt.Errorf("invalid launcher selection %q", selection))
	}
}
