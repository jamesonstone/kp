package prompt

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const userPromptMode = 0o600

type Registry interface {
	List() []Prompt
	Get(name string) (Prompt, error)
	Add(name, body string) (Prompt, error)
	Remove(name string) error
	PromoteToUser(name string) (Prompt, error)
}

type fileRegistry struct {
	userDir        string
	prompts        map[string]Prompt
	builtInSources map[string][]byte
}

func NewRegistry(userDir string) (Registry, error) {
	reg := &fileRegistry{
		userDir:        userDir,
		prompts:        make(map[string]Prompt),
		builtInSources: make(map[string][]byte),
	}

	builtIns, err := loadBuiltInsFromDefault()
	if err != nil {
		return nil, err
	}
	for _, p := range builtIns {
		reg.prompts[p.Name] = p
		source, err := loadBuiltInSource(p.Name)
		if err != nil {
			return nil, err
		}
		reg.builtInSources[p.Name] = source
	}

	if err := reg.loadUsers(); err != nil {
		return nil, err
	}

	return reg, nil
}

func (r *fileRegistry) List() []Prompt {
	prompts := make([]Prompt, 0, len(r.prompts))
	for _, p := range r.prompts {
		prompts = append(prompts, p)
	}
	sort.Slice(prompts, func(i, j int) bool {
		return prompts[i].Name < prompts[j].Name
	})
	return prompts
}

func (r *fileRegistry) Get(name string) (Prompt, error) {
	if err := ValidateName(name); err != nil {
		return Prompt{}, err
	}
	p, ok := r.prompts[name]
	if !ok {
		return Prompt{}, fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	return p, nil
}

func (r *fileRegistry) Add(name, body string) (Prompt, error) {
	if err := ValidateName(name); err != nil {
		return Prompt{}, err
	}
	if _, exists := r.prompts[name]; exists {
		return Prompt{}, fmt.Errorf("%w: %s", ErrExists, name)
	}

	doc, err := ParseDocument(name, []byte(body))
	if err != nil {
		return Prompt{}, err
	}
	if strings.TrimSpace(doc.Body) == "" {
		return Prompt{}, fmt.Errorf("%w: %s", ErrEmpty, name)
	}

	path := r.userPath(name)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return Prompt{}, fmt.Errorf("create prompt dir %q: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(body), userPromptMode); err != nil {
		return Prompt{}, fmt.Errorf("write prompt %q: %w", path, err)
	}

	return r.loadUserFile(path)
}

func (r *fileRegistry) Remove(name string) error {
	if err := ValidateName(name); err != nil {
		return err
	}

	p, ok := r.prompts[name]
	if !ok {
		return fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	if p.Source != SourceUser {
		return fmt.Errorf("%w: %s", ErrBuiltIn, name)
	}

	if err := os.Remove(p.FilePath); err != nil {
		return fmt.Errorf("remove prompt %q: %w", p.FilePath, err)
	}

	if _, ok := r.builtInSources[name]; ok {
		source, err := loadBuiltInsFromDefault()
		if err != nil {
			return err
		}
		for _, builtIn := range source {
			if builtIn.Name == name {
				r.prompts[name] = builtIn
				return nil
			}
		}
	}
	delete(r.prompts, name)
	return nil
}

func (r *fileRegistry) PromoteToUser(name string) (Prompt, error) {
	if err := ValidateName(name); err != nil {
		return Prompt{}, err
	}

	p, ok := r.prompts[name]
	if !ok {
		return Prompt{}, fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	if p.Source == SourceUser {
		return p, nil
	}

	source, ok := r.builtInSources[name]
	if !ok {
		return Prompt{}, fmt.Errorf("%w: %s", ErrNotFound, name)
	}

	path := r.userPath(name)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return Prompt{}, fmt.Errorf("create prompt dir %q: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, source, userPromptMode); err != nil {
		return Prompt{}, fmt.Errorf("write prompt %q: %w", path, err)
	}

	return r.loadUserFile(path)
}

func (r *fileRegistry) loadUsers() error {
	entries, err := os.ReadDir(r.userDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read user prompts %q: %w", r.userDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		if _, err := r.loadUserFile(filepath.Join(r.userDir, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func (r *fileRegistry) loadUserFile(path string) (Prompt, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return Prompt{}, fmt.Errorf("resolve prompt path %q: %w", path, err)
	}

	name := strings.TrimSuffix(filepath.Base(absPath), ".md")
	if err := ValidateName(name); err != nil {
		return Prompt{}, fmt.Errorf("load user prompt %q: %w", absPath, err)
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		return Prompt{}, fmt.Errorf("read user prompt %q: %w", absPath, err)
	}

	doc, err := ParseDocument(name, content)
	if err != nil {
		return Prompt{}, fmt.Errorf("load user prompt %q: %w", absPath, err)
	}

	p := Prompt{
		Name:     name,
		Label:    doc.Label,
		Source:   SourceUser,
		FilePath: absPath,
		Body:     doc.Body,
	}
	r.prompts[name] = p
	return p, nil
}

func (r *fileRegistry) userPath(name string) string {
	return filepath.Join(r.userDir, name+".md")
}
