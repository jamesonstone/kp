package prompt

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Document struct {
	Label string
	Body  string
}

type frontmatter struct {
	Label string `yaml:"label"`
}

func ParseDocument(name string, content []byte) (Document, error) {
	text := string(content)
	doc := Document{
		Label: defaultLabel(name),
		Body:  text,
	}

	if !strings.HasPrefix(text, "---\n") && !strings.HasPrefix(text, "---\r\n") {
		return doc, nil
	}

	afterOpen, openLen := trimOpeningDelimiter(text)
	if openLen == 0 {
		return doc, nil
	}

	yamlText, body, ok := splitFrontmatter(afterOpen)
	if !ok {
		return Document{}, fmt.Errorf("parse frontmatter for %q: missing closing delimiter", name)
	}

	var meta frontmatter
	if err := yaml.Unmarshal([]byte(yamlText), &meta); err != nil {
		return Document{}, fmt.Errorf("parse frontmatter for %q: %w", name, err)
	}

	if strings.TrimSpace(meta.Label) != "" {
		doc.Label = meta.Label
	}
	doc.Body = body
	return doc, nil
}

func trimOpeningDelimiter(text string) (string, int) {
	switch {
	case strings.HasPrefix(text, "---\r\n"):
		return text[5:], 5
	case strings.HasPrefix(text, "---\n"):
		return text[4:], 4
	default:
		return text, 0
	}
}

func splitFrontmatter(text string) (string, string, bool) {
	offset := 0
	for {
		line, nextOffset, ok := nextLine(text, offset)
		if strings.TrimSuffix(line, "\r") == "---" {
			return text[:offset], text[nextOffset:], true
		}
		if !ok {
			return "", "", false
		}
		offset = nextOffset
	}
}

func nextLine(text string, offset int) (line string, nextOffset int, ok bool) {
	if offset >= len(text) {
		return "", len(text), false
	}

	rest := text[offset:]
	newline := strings.IndexByte(rest, '\n')
	if newline == -1 {
		return rest, len(text), false
	}

	end := offset + newline
	return text[offset:end], end + 1, true
}
