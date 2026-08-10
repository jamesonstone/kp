package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jamesonstone/kp/internal/clipboard"
	"github.com/jamesonstone/kp/internal/config"
	"github.com/jamesonstone/kp/internal/prompt"
	"github.com/spf13/cobra"
)

func (a *app) newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Select a prompt",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			reg, err := a.loadRegistry()
			if err != nil {
				return err
			}
			if !a.listPlain && !a.verbose {
				return a.runPicker()
			}
			for _, p := range reg.List() {
				if a.verbose {
					fmt.Fprintf(a.stdout, "%s\t%s\t%s\n", p.Name, p.Label, p.Source.String())
					continue
				}
				fmt.Fprintln(a.stdout, p.Name)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&a.listPlain, "plain", false, "print prompt names without interactive selection")
	return cmd
}

func (a *app) newNewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "new <name>",
		Short: "Create a user prompt",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runNew(args[0])
		},
	}
}

func (a *app) newEditCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "edit <name>",
		Short: "Edit a user prompt",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runEdit(args[0])
		},
	}
}

func (a *app) newRMCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "rm <name>",
		Short: "Remove a user prompt",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.runRM(args[0])
		},
	}
}

func (a *app) runPrompt(name string) error {
	if err := prompt.ValidateName(name); err != nil {
		return NewExitError(ExitUser, err)
	}

	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}
	p, err := reg.Get(name)
	if err != nil {
		return mapPromptError(err)
	}

	if a.printOnly {
		fmt.Fprint(a.stdout, p.Body)
		return nil
	}

	cb := a.clipboardFactory()
	if err := cb.Copy(p.Body); err != nil {
		return NewExitError(ExitSystem, err)
	}
	if err := cb.Verify(p.Body, clipboard.DefaultVerifyTimeout); err != nil {
		return NewExitError(ExitSystem, err)
	}
	fmt.Fprintf(a.stderr, "✅ 📋 Prompt %q copied to clipboard.\n", p.Name)
	if !a.copyOnly {
		fmt.Fprintln(a.stderr, "🧾 Full prompt content is printed to stdout below.")
		fmt.Fprintln(a.stderr)
		fmt.Fprint(a.stdout, p.Body)
	}
	a.logVerbose("event=copy name=%s bytes=%d", p.Name, len(p.Body))
	return nil
}

func (a *app) runNew(name string) error {
	if err := prompt.ValidateName(name); err != nil {
		return NewExitError(ExitUser, err)
	}

	editor, err := a.resolveEditor()
	if err != nil {
		return err
	}

	paths, reg, err := a.loadPathsAndRegistry()
	if err != nil {
		return err
	}
	if _, err := reg.Get(name); err == nil {
		return NewExitError(ExitUser, fmt.Errorf("%w: %s", prompt.ErrExists, name))
	} else if !errors.Is(err, prompt.ErrNotFound) {
		return mapPromptError(err)
	}

	path := filepath.Join(paths.PromptsDir, name+".md")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		return NewExitError(ExitConfig, fmt.Errorf("write prompt %q: %w", path, err))
	}

	if err := a.runEditor(editor, path); err != nil {
		_ = os.Remove(path)
		return err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return NewExitError(ExitConfig, fmt.Errorf("read prompt %q: %w", path, err))
	}
	doc, err := prompt.ParseDocument(name, content)
	if err != nil {
		_ = os.Remove(path)
		return mapPromptError(err)
	}
	if strings.TrimSpace(doc.Body) == "" {
		_ = os.Remove(path)
		return NewExitError(ExitUser, fmt.Errorf("%w: %s", prompt.ErrEmpty, name))
	}
	return nil
}

func (a *app) runEdit(name string) error {
	if err := prompt.ValidateName(name); err != nil {
		return NewExitError(ExitUser, err)
	}

	editor, err := a.resolveEditor()
	if err != nil {
		return err
	}

	_, reg, err := a.loadPathsAndRegistry()
	if err != nil {
		return err
	}
	p, err := reg.Get(name)
	if err != nil {
		return mapPromptError(err)
	}
	if p.Source == prompt.SourceBuiltIn {
		p, err = reg.PromoteToUser(name)
		if err != nil {
			return mapPromptError(err)
		}
		fmt.Fprintf(a.stderr, "promoted built-in %q to %s\n", name, p.FilePath)
	}

	return a.runEditor(editor, p.FilePath)
}

func (a *app) runRM(name string) error {
	if err := prompt.ValidateName(name); err != nil {
		return NewExitError(ExitUser, err)
	}

	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}
	if err := reg.Remove(name); err != nil {
		return mapPromptError(err)
	}
	fmt.Fprintf(a.stderr, "removed %s\n", name)
	return nil
}

func (a *app) loadRegistry() (prompt.Registry, error) {
	_, reg, err := a.loadPathsAndRegistry()
	return reg, err
}

func (a *app) loadPathsAndRegistry() (config.Paths, prompt.Registry, error) {
	paths, err := config.Ensure(config.Options{ConfigDir: a.configDir})
	if err != nil {
		return config.Paths{}, nil, NewExitError(ExitConfig, err)
	}

	reg, err := a.registryFactory(paths.PromptsDir)
	if err != nil {
		return config.Paths{}, nil, mapPromptError(err)
	}
	return paths, reg, nil
}

func (a *app) logVerbose(format string, args ...any) {
	if !a.verbose {
		return
	}
	fmt.Fprintf(a.stderr, format+"\n", args...)
}

func mapPromptError(err error) error {
	switch {
	case errors.Is(err, prompt.ErrInvalidName),
		errors.Is(err, prompt.ErrReservedName),
		errors.Is(err, prompt.ErrNotFound),
		errors.Is(err, prompt.ErrExists),
		errors.Is(err, prompt.ErrEmpty),
		errors.Is(err, prompt.ErrBuiltIn):
		return NewExitError(ExitUser, err)
	default:
		return err
	}
}
