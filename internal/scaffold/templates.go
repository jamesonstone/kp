package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed templates/**
var templateFS embed.FS

func readTemplate(path string) (string, error) {
	content, err := fs.ReadFile(templateFS, "templates/"+path)
	if err != nil {
		return "", fmt.Errorf("read scaffold template %q: %w", path, err)
	}
	return string(content), nil
}

const gitignoreTemplate = `# Kit local generated environment, cache, and scratch artifacts
.env
.envrc
.kit/runs/
.kit/state.json
.kit/cache/
.kit/tmp/
.kit/temp/
.kit/*.tmp
.kit/*.lock
`

var gitignorePatterns = []string{
	".env",
	".envrc",
	".kit/runs/",
	".kit/state.json",
	".kit/cache/",
	".kit/tmp/",
	".kit/temp/",
	".kit/*.tmp",
	".kit/*.lock",
}
