package prompt

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type frontmatter struct {
	Label string `yaml:"label"`
}

func ParseFrontmatter(name, raw string) (label, body string, err error) {
	if !strings.HasPrefix(raw, "---\n") {
		return DefaultLabel(name), raw, nil
	}
	rest := raw[len("---\n"):]
	idx := strings.Index(rest, "\n---\n")
	if idx < 0 {
		return "", "", fmt.Errorf("malformed frontmatter")
	}
	meta := rest[:idx]
	body = rest[idx+len("\n---\n"):]
	var fm frontmatter
	if err := yaml.Unmarshal([]byte(meta), &fm); err != nil {
		return "", "", err
	}
	if strings.TrimSpace(fm.Label) == "" {
		fm.Label = DefaultLabel(name)
	}
	return fm.Label, body, nil
}
