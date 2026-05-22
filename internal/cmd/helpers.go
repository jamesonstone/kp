package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jamesonstone/kp/internal/config"
	"github.com/jamesonstone/kp/internal/prompt"
)

func loadRegistry(opts *runtimeOptions) (prompt.Registry, string, error) {
	var dir string
	var err error
	if opts.ConfigDir != "" {
		dir = filepath.Join(opts.ConfigDir, "kp", "prompts")
		if err = os.MkdirAll(dir, 0o700); err != nil {
			return nil, "", configErr(err)
		}
	} else {
		dir, err = config.PromptsDir()
		if err != nil {
			return nil, "", configErr(err)
		}
	}
	reg, err := prompt.NewRegistry(opts.BuiltinFS, dir)
	if err != nil {
		return nil, "", configErr(err)
	}
	return reg, dir, nil
}

func printVerbose(opts *runtimeOptions, format string, args ...any) {
	if opts.Verbose {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}

func sourceString(s prompt.Source) string {
	if s == prompt.SourceUser {
		return "user"
	}
	return "builtin"
}

func parseSelectedName(line string) string {
	parts := strings.SplitN(line, "\t", 2)
	return strings.TrimSpace(parts[0])
}
