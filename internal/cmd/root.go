package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/jamesonstone/kp/internal/clipboard"
	"github.com/jamesonstone/kp/internal/prompt"
	"github.com/spf13/cobra"
)

const (
	ExitSuccess = 0
	ExitUser    = 1
	ExitSystem  = 2
	ExitConfig  = 3
	ExitCancel  = 130
)

// ExitError carries a stable process exit code to the executable boundary.
type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string {
	if e == nil || e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

func (e *ExitError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func NewExitError(code int, err error) error {
	if err == nil {
		return nil
	}
	return &ExitError{Code: code, Err: err}
}

func ExitCode(err error) int {
	if err == nil {
		return ExitSuccess
	}

	var exitErr *ExitError
	if errors.As(err, &exitErr) && exitErr.Code != 0 {
		return exitErr.Code
	}

	return ExitUser
}

type Options struct {
	Version          string
	Commit           string
	Stdin            io.Reader
	Stdout           io.Writer
	Stderr           io.Writer
	RegistryFactory  func(userDir string) (prompt.Registry, error)
	ClipboardFactory func() clipboard.Clipboard
	LookPath         func(file string) (string, error)
	FZFRunner        func(prompts []prompt.Prompt) (string, error)
	LauncherRunner   func(items []LauncherItem) (string, error)
	PortLookup       func(port int) ([]PortProcess, error)
	PortStop         func(pid int, force bool) error
	Getenv           func(key string) string
	EditorRunner     func(name string, args []string, path string) error
}

func NewRoot(opts Options) *cobra.Command {
	app := newApp(opts)

	cmd := &cobra.Command{
		Use:           "kp",
		Short:         "CLI prompt utilities",
		Long:          "Low-friction prompt utilities for printing and copying reusable prompts.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       fmt.Sprintf("version=%s commit=%s", app.version, app.commit),
		Args:          cobra.ArbitraryArgs,
		RunE:          app.runRoot,
	}
	cmd.SetIn(app.stdin)
	cmd.SetOut(app.stdout)
	cmd.SetErr(app.stderr)

	cmd.PersistentFlags().StringVar(&app.configDir, "config", "", "config root override")
	cmd.PersistentFlags().BoolVar(&app.verbose, "verbose", false, "emit verbose output")
	cmd.PersistentFlags().BoolVar(&app.noFzf, "no-fzf", false, "use numbered picker instead of fzf")
	cmd.Flags().BoolVar(&app.copyOnly, "copy", false, "copy without printing")
	cmd.Flags().BoolVar(&app.printOnly, "print", false, "print without clipboard side effects")

	cmd.AddCommand(app.newListCommand())
	cmd.AddCommand(app.newNewCommand())
	cmd.AddCommand(app.newEditCommand())
	cmd.AddCommand(app.newRMCommand())
	cmd.AddCommand(app.newFindPortCommand())
	cmd.AddCommand(app.newScaffoldCommand())
	configureRootHelp(cmd)

	return cmd
}

type app struct {
	version          string
	commit           string
	stdin            io.Reader
	inputReader      *bufio.Reader
	stdout           io.Writer
	stderr           io.Writer
	registryFactory  func(userDir string) (prompt.Registry, error)
	clipboardFactory func() clipboard.Clipboard
	lookPath         func(file string) (string, error)
	fzfRunner        func(prompts []prompt.Prompt) (string, error)
	launcherRunner   func(items []LauncherItem) (string, error)
	portLookup       func(context.Context, int) ([]PortProcess, error)
	portStopper      func(context.Context, PortProcess, bool) error
	getenv           func(key string) string
	editorRunner     func(name string, args []string, path string) error

	configDir string
	verbose   bool
	copyOnly  bool
	printOnly bool
	noFzf     bool
	listPlain bool

	scaffoldDir    string
	scaffoldDryRun bool
	scaffoldForce  bool

	portCopy  string
	portStop  bool
	portForce bool
}

func newApp(opts Options) *app {
	if opts.Stdout == nil {
		opts.Stdout = io.Discard
	}
	if opts.Stderr == nil {
		opts.Stderr = io.Discard
	}
	if opts.Stdin == nil {
		opts.Stdin = strings.NewReader("")
	}
	if opts.Version == "" {
		opts.Version = "dev"
	}
	if opts.Commit == "" {
		opts.Commit = "unknown"
	}
	if opts.RegistryFactory == nil {
		opts.RegistryFactory = prompt.NewRegistry
	}
	if opts.ClipboardFactory == nil {
		opts.ClipboardFactory = clipboard.New
	}
	if opts.LookPath == nil {
		opts.LookPath = exec.LookPath
	}
	if opts.Getenv == nil {
		opts.Getenv = os.Getenv
	}
	portLookup := lookupPortProcesses
	if opts.PortLookup != nil {
		portLookup = func(_ context.Context, port int) ([]PortProcess, error) {
			return opts.PortLookup(port)
		}
	}
	portStopper := stopPortProcess
	if opts.PortStop != nil {
		portStopper = func(_ context.Context, process PortProcess, force bool) error {
			return opts.PortStop(process.PID, force)
		}
	}

	return &app{
		version:          opts.Version,
		commit:           opts.Commit,
		stdin:            opts.Stdin,
		stdout:           opts.Stdout,
		stderr:           opts.Stderr,
		inputReader:      bufio.NewReader(opts.Stdin),
		registryFactory:  opts.RegistryFactory,
		clipboardFactory: opts.ClipboardFactory,
		lookPath:         opts.LookPath,
		fzfRunner:        opts.FZFRunner,
		launcherRunner:   opts.LauncherRunner,
		portLookup:       portLookup,
		portStopper:      portStopper,
		getenv:           opts.Getenv,
		editorRunner:     opts.EditorRunner,
	}
}

func (a *app) runRoot(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return a.runLauncher(cmd)
	}
	if len(args) != 1 {
		return NewExitError(ExitUser, fmt.Errorf("expected one prompt name"))
	}
	return a.runPrompt(args[0])
}
