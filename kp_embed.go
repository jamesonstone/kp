package kp

import "embed"

//go:embed prompts/*.md
var BuiltinPromptsFS embed.FS
