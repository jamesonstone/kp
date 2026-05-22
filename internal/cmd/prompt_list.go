package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newPromptListCommand(opts *runtimeOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List prompts",
		RunE: func(cmd *cobra.Command, args []string) error {
			reg, _, err := loadRegistry(opts)
			if err != nil {
				return err
			}
			for _, p := range reg.List() {
				if opts.Verbose {
					fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\n", p.Name, p.Label, sourceString(p.Source))
				} else {
					fmt.Fprintln(cmd.OutOrStdout(), p.Name)
				}
			}
			return nil
		},
	}
}
