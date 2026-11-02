package cmd

import (
	"fmt"

	"github.com/jamesonstone/kp/internal/prompt"
	"github.com/spf13/cobra"
)

// newCursorCommand serves Cursor-specific prompts under `kp cursor`.
func (a *app) newCursorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cursor",
		Short: "Cursor-specific prompts",
		Long:  "Cursor-specific prompts for Cursor cloud environments.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.PersistentFlags().BoolVar(&a.copyOnly, "copy", false, "copy without printing")
	cmd.PersistentFlags().BoolVar(&a.printOnly, "print", false, "print without clipboard side effects")

	prompts, err := prompt.CursorBuiltIns()
	if err != nil {
		cmd.RunE = func(*cobra.Command, []string) error { return err }
		return cmd
	}
	for _, p := range prompts {
		cmd.AddCommand(&cobra.Command{
			Use:   p.Name,
			Short: p.Label,
			Args:  cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error {
				return a.emitPrompt(p)
			},
		})
	}
	return cmd
}

func (a *app) emitCursorPrompt(name string) error {
	prompts, err := prompt.CursorBuiltIns()
	if err != nil {
		return NewExitError(ExitConfig, err)
	}
	for _, p := range prompts {
		if p.Name == name {
			return a.emitPrompt(p)
		}
	}
	return NewExitError(ExitUser, fmt.Errorf("unknown cursor prompt %q", name))
}
