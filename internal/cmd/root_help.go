package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/jamesonstone/kp/internal/prompt"
	"github.com/spf13/cobra"
)

type helpRow struct {
	command string
	summary string
}

func configureRootHelp(root *cobra.Command) {
	defaultHelp := root.HelpFunc()
	root.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd == root {
			_ = renderRootHelp(cmd)
			return
		}
		defaultHelp(cmd, args)
	})
}

func renderRootHelp(cmd *cobra.Command) error {
	out := cmd.OutOrStdout()
	style := styleForWriter(out)

	if _, err := fmt.Fprintln(out, strings.TrimRight(cmd.Long, "\n")); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(out); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(out, style.title("🚀", "Usage")); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "  %s                 open launcher\n", cmd.CommandPath()); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "  %s <prompt>        print and copy prompt\n", cmd.CommandPath()); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "  %s [command]       run command\n", cmd.CommandPath()); err != nil {
		return err
	}

	promptRows := builtInPromptRows(cmd.CommandPath())
	if err := renderHelpSection(out, style, "🧠", "Prompt Commands", promptRows); err != nil {
		return err
	}

	libraryRows := []helpRow{
		{command: cmd.CommandPath() + " list", summary: "Pick a prompt interactively"},
		{command: cmd.CommandPath() + " list --plain", summary: "Print prompt names"},
		{command: cmd.CommandPath() + " list --verbose", summary: "Print name, label, and source"},
		{command: cmd.CommandPath() + " new <name>", summary: "Create a user prompt"},
		{command: cmd.CommandPath() + " edit <name>", summary: "Edit or promote a prompt"},
		{command: cmd.CommandPath() + " rm <name>", summary: "Remove a user prompt"},
	}
	if err := renderHelpSection(out, style, "🧰", "Prompt Library", libraryRows); err != nil {
		return err
	}

	portRows := []helpRow{
		{command: cmd.CommandPath() + " find-port <port>", summary: "Inspect a process on a port"},
		{command: cmd.CommandPath() + " port-find <port>", summary: "Alias for find-port"},
	}
	if err := renderHelpSection(out, style, "🔍", "Port Tools", portRows); err != nil {
		return err
	}

	setupRows := []helpRow{
		{command: cmd.CommandPath() + " scaffold", summary: "Create repo support files"},
		{command: cmd.CommandPath() + " scaffold --dry-run", summary: "Preview scaffold actions"},
	}
	if err := renderHelpSection(out, style, "🏗️", "Repo Setup", setupRows); err != nil {
		return err
	}

	utilityRows := []helpRow{
		{command: cmd.CommandPath() + " --version", summary: "Show version metadata"},
		{command: cmd.CommandPath() + " help <command>", summary: "Show command help"},
		{command: cmd.CommandPath() + " completion", summary: "Generate shell completion"},
	}
	if err := renderHelpSection(out, style, "🛠️", "Utilities", utilityRows); err != nil {
		return err
	}

	flags := strings.TrimRight(cmd.Flags().FlagUsages(), "\n")
	if strings.TrimSpace(flags) != "" {
		if _, err := fmt.Fprintln(out); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(out, style.title("⚙️", "Flags")); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(out, flags); err != nil {
			return err
		}
	}

	_, err := fmt.Fprintf(out, "\n%s \"%s [command] --help\" for more information about a command.\n",
		moreInfoLabel(style.enabled), cmd.CommandPath())
	return err
}

func builtInPromptRows(commandPath string) []helpRow {
	builtIns, err := prompt.BuiltIns()
	if err != nil {
		return []helpRow{
			{command: commandPath + " clarify", summary: "Clarify before implementing"},
			{command: commandPath + " agent-handoff", summary: "Agent-to-agent handoff"},
			{command: commandPath + " chat-handoff", summary: "Chat-to-agent handoff"},
			{command: commandPath + " pr", summary: "Pull request workflow"},
		}
	}

	rows := make([]helpRow, 0, len(builtIns))
	for _, p := range builtIns {
		rows = append(rows, helpRow{
			command: commandPath + " " + p.Name,
			summary: p.Label,
		})
	}
	return rows
}

func renderHelpSection(out io.Writer, style humanOutputStyle, emoji string, title string, rows []helpRow) error {
	if len(rows) == 0 {
		return nil
	}

	if _, err := fmt.Fprintln(out); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(out, style.title(emoji, title)); err != nil {
		return err
	}

	padding := helpRowPadding(rows)
	for _, row := range rows {
		if _, err := fmt.Fprintf(out, "  %s %s\n", padRight(row.command, padding), row.summary); err != nil {
			return err
		}
	}
	return nil
}

func helpRowPadding(rows []helpRow) int {
	width := 0
	for _, row := range rows {
		if len(row.command) > width {
			width = len(row.command)
		}
	}
	return width + 2
}

func padRight(value string, width int) string {
	if len(value) >= width {
		return value
	}
	return value + strings.Repeat(" ", width-len(value))
}

func moreInfoLabel(enabled bool) string {
	if enabled {
		return ansiWhiteBold + "🔎 Use" + ansiReset
	}
	return "Use"
}
