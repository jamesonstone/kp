package cmd

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type Options struct {
	Version   string
	Commit    string
	BuiltinFS embed.FS
}

type runtimeOptions struct {
	Options
	ConfigDir string
	Verbose   bool
}

type exitError struct {
	Code int
	Err  error
}

func (e *exitError) Error() string { return e.Err.Error() }

func Execute(opts Options) int {
	return executeWithArgs(opts, os.Args[1:], os.Stdout, os.Stderr)
}

func executeWithArgs(opts Options, args []string, stdout, stderr io.Writer) int {
	r := &runtimeOptions{Options: opts}
	root := &cobra.Command{
		Use:   "kp",
		Short: "CLI prompt utilities",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SilenceUsage = true
	root.SilenceErrors = true
	root.Version = opts.Version + " (" + opts.Commit + ")"
	root.SetVersionTemplate("{{.Version}}\n")
	root.PersistentFlags().StringVar(&r.ConfigDir, "config", "", "override config directory")
	root.PersistentFlags().BoolVar(&r.Verbose, "verbose", false, "enable verbose output")
	root.AddCommand(newPromptCommand(r))
	root.SetArgs(args)

	if err := root.Execute(); err != nil {
		var ex *exitError
		if errors.As(err, &ex) {
			fmt.Fprintln(stderr, ex.Err)
			return ex.Code
		}
		if strings.Contains(err.Error(), "^C") {
			return 130
		}
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func userErr(err error) error   { return &exitError{Code: 1, Err: err} }
func systemErr(err error) error { return &exitError{Code: 2, Err: err} }
func configErr(err error) error { return &exitError{Code: 3, Err: err} }
