package kp

import "embed"

// PromptFS contains the built-in prompt markdown files shipped with kp.
//
//go:embed prompts/*.md
var PromptFS embed.FS
