package kp

import "embed"

// PromptFS contains the built-in prompt markdown files shipped with kp.
// Root prompts live in prompts/; legacy prompts served by `kp v0` live in
// prompts/v0/.
//
//go:embed prompts/*.md prompts/v0/*.md
var PromptFS embed.FS
