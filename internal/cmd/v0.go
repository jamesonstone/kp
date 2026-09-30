package cmd

import (
	"github.com/jamesonstone/kp/internal/prompt"
	"github.com/spf13/cobra"
)

// newV0Command serves the legacy prompt commands under `kp v0 <command>` so
// their former root names are free for reuse. Behavior matches the old root
// commands: print to stdout and copy with verification, honoring --print and
// --copy.
func (a *app) newV0Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "v0",
		Short: "Legacy prompt commands",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	cmd.PersistentFlags().BoolVar(&a.copyOnly, "copy", false, "copy without printing")
	cmd.PersistentFlags().BoolVar(&a.printOnly, "print", false, "print without clipboard side effects")

	prompts, err := prompt.V0BuiltIns()
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

func isV0Prompt(name string) bool {
	prompts, err := prompt.V0BuiltIns()
	if err != nil {
		return false
	}
	for _, p := range prompts {
		if p.Name == name {
			return true
		}
	}
	return false
}
