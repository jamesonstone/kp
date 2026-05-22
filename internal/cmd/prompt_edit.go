package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jamesonstone/kp/internal/prompt"
	"github.com/spf13/cobra"
)

func newPromptEditCommand(opts *runtimeOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "edit <name>",
		Short: "Edit a prompt",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			reg, dir, err := loadRegistry(opts)
			if err != nil {
				return err
			}
			p, err := reg.Get(name)
			if err != nil {
				if errors.Is(err, prompt.ErrNotFound) {
					return userErr(err)
				}
				return configErr(err)
			}
			path := p.FilePath
			if p.Source == prompt.SourceBuiltIn {
				promoted, err := reg.PromoteToUser(name)
				if err != nil {
					return configErr(err)
				}
				path = promoted.FilePath
				fmt.Fprintf(os.Stderr, "promoted built-in '%s' to %s\n", name, path)
			}
			if path == "" {
				path = filepath.Join(dir, name+".md")
			}
			return openEditor(path)
		},
	}
}
