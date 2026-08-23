package main

import (
	"fmt"
	"os"

	kpcmd "github.com/jamesonstone/kp/internal/cmd"
)

var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	root := kpcmd.NewRoot(kpcmd.Options{
		Version: version,
		Commit:  commit,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	})

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, kpcmd.ExitMessage(err))
		os.Exit(kpcmd.ExitCode(err))
	}
}
