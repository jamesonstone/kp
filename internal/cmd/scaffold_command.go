package cmd

import (
	"fmt"

	"github.com/jamesonstone/kp/internal/scaffold"
	"github.com/spf13/cobra"
)

func (a *app) newScaffoldCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scaffold",
		Short: "Create repo support files",
		Long:  "Create repo support files such as agent instructions, review config, PR template, and local env files without creating direct Kit project state.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runScaffold()
		},
	}
	cmd.Flags().StringVar(&a.scaffoldDir, "dir", "", "target directory (default current directory)")
	cmd.Flags().BoolVar(&a.scaffoldDryRun, "dry-run", false, "show planned scaffold actions without writing files")
	cmd.Flags().BoolVar(&a.scaffoldForce, "force", false, "overwrite existing scaffold files except .gitignore")
	return cmd
}

func (a *app) runScaffold() error {
	results, err := scaffold.Run(scaffold.Options{
		Dir:    a.scaffoldDir,
		DryRun: a.scaffoldDryRun,
		Force:  a.scaffoldForce,
	})
	if err != nil {
		return NewExitError(ExitConfig, err)
	}

	for _, result := range results {
		fmt.Fprintf(a.stdout, "%s\t%s\n", result.Action, result.Path)
	}
	return nil
}
