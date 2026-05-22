package cmd

import (
	"errors"

	"github.com/jamesonstone/kp/internal/prompt"
	"github.com/spf13/cobra"
)

func newPromptRMCommand(opts *runtimeOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "rm <name>",
		Short: "Remove a user prompt",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reg, _, err := loadRegistry(opts)
			if err != nil {
				return err
			}
			if err := reg.Remove(args[0]); err != nil {
				if errors.Is(err, prompt.ErrBuiltIn) || errors.Is(err, prompt.ErrNotFound) {
					return userErr(err)
				}
				return configErr(err)
			}
			return nil
		},
	}
}
