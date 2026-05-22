package main

import (
	"os"

	kp "github.com/jamesonstone/kp"
	"github.com/jamesonstone/kp/internal/cmd"
)

var (
	version = "0.0.0-dev"
	commit  = "unknown"
)

func main() {
	os.Exit(cmd.Execute(cmd.Options{Version: version, Commit: commit, BuiltinFS: kp.BuiltinPromptsFS}))
}
