package prompt

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	kp "github.com/jamesonstone/kp"
)

func loadBuiltInsFromDefault() ([]Prompt, error) {
	fsys, err := fs.Sub(kp.PromptFS, "prompts")
	if err != nil {
		return nil, fmt.Errorf("open built-in prompts: %w", err)
	}
	return loadBuiltIns(fsys)
}

func BuiltIns() ([]Prompt, error) {
	return loadBuiltInsFromDefault()
}

func loadBuiltInSource(name string) ([]byte, error) {
	if err := ValidateName(name); err != nil {
		return nil, err
	}

	content, err := kp.PromptFS.ReadFile("prompts/" + name + ".md")
	if err != nil {
		return nil, fmt.Errorf("read built-in prompt %q: %w", name, err)
	}
	return content, nil
}

func loadBuiltIns(fsys fs.FS) ([]Prompt, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("read built-in prompts: %w", err)
	}

	prompts := make([]Prompt, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".md")
		if err := ValidateName(name); err != nil {
			return nil, fmt.Errorf("load built-in prompt %q: %w", entry.Name(), err)
		}

		content, err := fs.ReadFile(fsys, entry.Name())
		if err != nil {
			return nil, fmt.Errorf("read built-in prompt %q: %w", entry.Name(), err)
		}

		doc, err := ParseDocument(name, content)
		if err != nil {
			return nil, err
		}

		prompts = append(prompts, Prompt{
			Name:   name,
			Label:  doc.Label,
			Source: SourceBuiltIn,
			Body:   doc.Body,
		})
	}

	sort.Slice(prompts, func(i, j int) bool {
		return prompts[i].Name < prompts[j].Name
	})
	return prompts, nil
}
