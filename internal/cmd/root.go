package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jamesonstone/kp/internal/clipboard"
	"github.com/jamesonstone/kp/internal/config"
	"github.com/jamesonstone/kp/internal/prompt"
	"github.com/jamesonstone/kp/internal/scaffold"
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
	cmd.AddCommand(app.newScaffoldCommand())
	configureRootHelp(cmd)

	return cmd
}

type app struct {
	version          string
	commit           string
	stdin            io.Reader
	stdout           io.Writer
	stderr           io.Writer
	registryFactory  func(userDir string) (prompt.Registry, error)
	clipboardFactory func() clipboard.Clipboard
	lookPath         func(file string) (string, error)
	fzfRunner        func(prompts []prompt.Prompt) (string, error)
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

	return &app{
		version:          opts.Version,
		commit:           opts.Commit,
		stdin:            opts.Stdin,
		stdout:           opts.Stdout,
		stderr:           opts.Stderr,
		registryFactory:  opts.RegistryFactory,
		clipboardFactory: opts.ClipboardFactory,
		lookPath:         opts.LookPath,
		fzfRunner:        opts.FZFRunner,
		getenv:           opts.Getenv,
		editorRunner:     opts.EditorRunner,
	}
}

func (a *app) runRoot(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}
	if len(args) != 1 {
		return NewExitError(ExitUser, fmt.Errorf("expected one prompt name"))
	}
	return a.runPrompt(args[0])
}

func (a *app) runPicker() error {
	reg, err := a.loadRegistry()
	if err != nil {
		return err
	}

	name := ""
	if a.noFzf {
		name, err = a.pickNumbered(reg.List())
	} else {
		name, err = a.pickFZF(reg.List())
	}
	if err != nil {
		return err
	}

	return a.runPrompt(name)
}

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

func (a *app) pickFZF(prompts []prompt.Prompt) (string, error) {
	if _, err := a.lookPath("fzf"); err != nil {
		return "", NewExitError(ExitConfig, errors.New("fzf not found; install fzf via 'brew install fzf' or use --no-fzf"))
	}

	if a.fzfRunner != nil {
		name, err := a.fzfRunner(prompts)
		if err != nil {
			return "", mapPickerError(err)
		}
		return name, nil
	}

	name, err := runFZF(prompts)
	if err != nil {
		return "", mapPickerError(err)
	}
	return name, nil
}

func (a *app) pickNumbered(prompts []prompt.Prompt) (string, error) {
	for i, p := range prompts {
		fmt.Fprintf(a.stderr, "%d\t%s\t%s\n", i+1, p.Name, p.Label)
	}

	line, err := bufio.NewReader(a.stdin).ReadString('\n')
	if err != nil && !(errors.Is(err, io.EOF) && line != "") {
		if errors.Is(err, io.EOF) {
			return "", NewExitError(ExitCancel, errors.New("picker cancelled"))
		}
		return "", NewExitError(ExitUser, err)
	}

	choiceText := strings.TrimSpace(line)
	if choiceText == "" {
		return "", NewExitError(ExitCancel, errors.New("picker cancelled"))
	}

	choice, err := strconv.Atoi(choiceText)
	if err != nil {
		return "", NewExitError(ExitUser, fmt.Errorf("invalid selection %q", choiceText))
	}
	if choice < 1 || choice > len(prompts) {
		return "", NewExitError(ExitUser, fmt.Errorf("selection %d out of range", choice))
	}

	return prompts[choice-1].Name, nil
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

type editorCommand struct {
	name string
	args []string
}

func (a *app) resolveEditor() (editorCommand, error) {
	editorText := a.getenv("KP_EDITOR")
	if strings.TrimSpace(editorText) == "" {
		editorText = a.getenv("EDITOR")
	}
	if strings.TrimSpace(editorText) == "" {
		editorText = "vi"
	}

	fields := strings.Fields(editorText)
	if len(fields) == 0 {
		fields = []string{"vi"}
	}

	name := fields[0]
	if _, err := a.lookPath(name); err != nil {
		return editorCommand{}, NewExitError(ExitConfig, fmt.Errorf("editor %q not found: %w", name, err))
	}
	return editorCommand{name: name, args: fields[1:]}, nil
}

func (a *app) runEditor(editor editorCommand, path string) error {
	if a.editorRunner != nil {
		return mapEditorError(a.editorRunner(editor.name, editor.args, path))
	}

	args := append(append([]string{}, editor.args...), path)
	cmd := exec.Command(editor.name, args...)
	cmd.Stdin = a.stdin
	cmd.Stdout = a.stdout
	cmd.Stderr = a.stderr
	return mapEditorError(cmd.Run())
}

func mapEditorError(err error) error {
	if err == nil {
		return nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == ExitCancel {
		return NewExitError(ExitCancel, err)
	}
	if errors.Is(err, errPickerCanceled) {
		return NewExitError(ExitCancel, err)
	}
	return NewExitError(ExitConfig, err)
}

var errPickerCanceled = errors.New("picker cancelled")

func mapPickerError(err error) error {
	if errors.Is(err, errPickerCanceled) {
		return NewExitError(ExitCancel, err)
	}
	return NewExitError(ExitUser, err)
}

func runFZF(prompts []prompt.Prompt) (string, error) {
	previewDir, err := os.MkdirTemp("", "kp-preview-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(previewDir)

	var input strings.Builder
	for _, p := range prompts {
		if err := os.WriteFile(filepath.Join(previewDir, p.Name), []byte(p.Body), 0o600); err != nil {
			return "", err
		}
		fmt.Fprintf(&input, "%s\t%s\t%s\n", promptIcon(p), p.Name, p.Label)
	}

	cmd := exec.Command(
		"fzf",
		"--height", "60%",
		"--reverse",
		"--cycle",
		"--prompt", "📋 kp list › ",
		"--pointer", "👉",
		"--header", "enter: copy prompt · tab/shift-tab: cycle · esc: cancel",
		"--delimiter", "\t",
		"--with-nth", "1,2,3",
		"--nth", "2,3",
		"--bind", "tab:down,btab:up",
		"--preview", "cat "+shellQuote(previewDir)+"/{2}",
	)
	cmd.Stdin = strings.NewReader(input.String())
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %v", errPickerCanceled, err)
	}

	selected := strings.TrimSpace(stdout.String())
	if selected == "" {
		return "", errPickerCanceled
	}
	_, rest, ok := strings.Cut(selected, "\t")
	if !ok {
		return "", fmt.Errorf("invalid picker selection")
	}
	name, _, _ := strings.Cut(rest, "\t")
	return name, nil
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func promptIcon(p prompt.Prompt) string {
	if p.Source == prompt.SourceUser {
		return "📝"
	}
	return "📦"
}
