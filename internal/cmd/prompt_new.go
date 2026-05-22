package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jamesonstone/kp/internal/prompt"
	"github.com/spf13/cobra"
)

func newPromptNewCommand(opts *runtimeOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "new <name>",
		Short: "Create a new prompt",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if err := prompt.ValidateName(name); err != nil {
				return userErr(err)
			}
			reg, dir, err := loadRegistry(opts)
			if err != nil {
				return err
			}
			if _, err := reg.Get(name); err == nil {
				return userErr(fmt.Errorf("prompt '%s' already exists; use 'edit' instead", name))
			} else if !errors.Is(err, prompt.ErrNotFound) {
				return configErr(err)
			}
			path := filepath.Join(dir, name+".md")
			if err := os.WriteFile(path, []byte(""), 0o600); err != nil {
				return configErr(err)
			}
			if err := openEditor(path); err != nil {
				return err
			}
			body, err := os.ReadFile(path)
			if err != nil {
				return configErr(err)
			}
			if strings.TrimSpace(string(body)) == "" {
				_ = os.Remove(path)
				return userErr(errors.New("prompt body is empty"))
			}
			return nil
		},
	}
}
