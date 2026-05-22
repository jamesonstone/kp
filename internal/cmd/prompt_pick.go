package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jamesonstone/kp/internal/clipboard"
	"github.com/jamesonstone/kp/internal/prompt"
	"github.com/spf13/cobra"
)

func runPromptPick(opts *runtimeOptions, cmd *cobra.Command, args []string) error {
	reg, _, err := loadRegistry(opts)
	if err != nil {
		return err
	}
	copyOnly, _ := cmd.Flags().GetBool("copy")
	printOnly, _ := cmd.Flags().GetBool("print")
	noFZF, _ := cmd.Flags().GetBool("no-fzf")
	if copyOnly && printOnly {
		return userErr(errors.New("--copy and --print are mutually exclusive"))
	}

	name := ""
	if len(args) == 1 {
		name = args[0]
	} else {
		name, err = pickNameInteractive(reg, noFZF)
		if err != nil {
			return err
		}
	}

	p, err := reg.Get(name)
	if err != nil {
		if errors.Is(err, prompt.ErrNotFound) {
			return userErr(err)
		}
		return configErr(err)
	}
	if printOnly {
		fmt.Fprint(cmd.OutOrStdout(), p.Body)
		return nil
	}
	cb := clipboard.New()
	if copyOnly {
		if err := cb.Copy(p.Body); err != nil {
			return systemErr(err)
		}
		if err := cb.Verify(p.Body, clipboard.DefaultVerifyTimeout); err != nil {
			return systemErr(err)
		}
		return nil
	}
	if err := cb.CopyAndPaste(p.Body); err != nil {
		return systemErr(err)
	}
	fmt.Fprintf(os.Stderr, "✅ %s\n", p.Name)
	return nil
}

func pickNameInteractive(reg prompt.Registry, noFZF bool) (string, error) {
	prompts := reg.List()
	if len(prompts) == 0 {
		return "", userErr(errors.New("no prompts available"))
	}
	if noFZF {
		return pickNameFallback(prompts)
	}
	if _, err := exec.LookPath("fzf"); err != nil {
		return "", configErr(errors.New("install fzf via 'brew install fzf' or use --no-fzf"))
	}

	var b strings.Builder
	for _, p := range prompts {
		clean := strings.ReplaceAll(strings.ReplaceAll(p.Body, "\n", " "), "\t", " ")
		fmt.Fprintf(&b, "%s\t%s\t%s\t%s\n", p.Name, p.Label, sourceString(p.Source), clean)
	}

	c := exec.Command("fzf", "--delimiter=\t", "--with-nth=1,2,3", "--preview", "echo {4}")
	c.Stdin = strings.NewReader(b.String())
	out, err := c.Output()
	if err != nil {
		var ex *exec.ExitError
		if errors.As(err, &ex) && ex.ExitCode() == 130 {
			return "", &exitError{Code: 130, Err: errors.New("selection canceled")}
		}
		return "", systemErr(err)
	}
	name := parseSelectedName(string(out))
	if name == "" {
		return "", userErr(errors.New("no prompt selected"))
	}
	return name, nil
}

func pickNameFallback(prompts []prompt.Prompt) (string, error) {
	for i, p := range prompts {
		fmt.Fprintf(os.Stderr, "%d) %s\n", i+1, p.Name)
	}
	fmt.Fprint(os.Stderr, "Select prompt number: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", userErr(err)
	}
	line = strings.TrimSpace(line)
	n, err := strconv.Atoi(line)
	if err != nil || n < 1 || n > len(prompts) {
		return "", userErr(errors.New("invalid selection"))
	}
	return prompts[n-1].Name, nil
}
