package cmd

import "github.com/spf13/cobra"

func newPromptCommand(opts *runtimeOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prompt [name]",
		Short: "Prompt utilities",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPromptPick(opts, cmd, args)
		},
	}
	cmd.Flags().Bool("copy", false, "copy to clipboard without pasting")
	cmd.Flags().Bool("print", false, "print prompt body to stdout")
	cmd.Flags().Bool("no-fzf", false, "use fallback numbered picker")
	cmd.AddCommand(newPromptListCommand(opts))
	cmd.AddCommand(newPromptNewCommand(opts))
	cmd.AddCommand(newPromptEditCommand(opts))
	cmd.AddCommand(newPromptRMCommand(opts))
	return cmd
}
