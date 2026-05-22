package prompt

import (
	"embed"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type registry struct {
	userDir string
	prompts map[string]Prompt
}

func NewRegistry(builtin embed.FS, userDir string) (Registry, error) {
	r := &registry{userDir: userDir, prompts: map[string]Prompt{}}
	if err := r.loadBuiltins(builtin); err != nil {
		return nil, err
	}
	if err := r.loadUsers(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *registry) List() []Prompt {
	out := make([]Prompt, 0, len(r.prompts))
	for _, p := range r.prompts {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func (r *registry) Get(name string) (Prompt, error) {
	if p, ok := r.prompts[name]; ok {
		return p, nil
	}
	return Prompt{}, ErrNotFound
}

func (r *registry) Add(name, body string) (Prompt, error) {
	if err := ValidateName(name); err != nil {
		return Prompt{}, err
	}
	if strings.TrimSpace(body) == "" {
		return Prompt{}, ErrEmpty
	}
	if _, ok := r.prompts[name]; ok {
		return Prompt{}, ErrExists
	}
	path := filepath.Join(r.userDir, name+".md")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		return Prompt{}, err
	}
	label, parsedBody, err := ParseFrontmatter(name, body)
	if err != nil {
		return Prompt{}, err
	}
	p := Prompt{Name: name, Label: label, Source: SourceUser, FilePath: path, Body: parsedBody}
	r.prompts[name] = p
	return p, nil
}

func (r *registry) Remove(name string) error {
	p, ok := r.prompts[name]
	if !ok {
		return ErrNotFound
	}
	if p.Source != SourceUser {
		return ErrBuiltIn
	}
	if err := os.Remove(p.FilePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNotFound
		}
		return err
	}
	delete(r.prompts, name)
	return nil
}

func (r *registry) PromoteToUser(name string) (Prompt, error) {
	p, ok := r.prompts[name]
	if !ok {
		return Prompt{}, ErrNotFound
	}
	if p.Source == SourceUser {
		return p, nil
	}
	path := filepath.Join(r.userDir, name+".md")
	if err := os.WriteFile(path, []byte(p.Body), 0o600); err != nil {
		return Prompt{}, err
	}
	p.Source = SourceUser
	p.FilePath = path
	r.prompts[name] = p
	return p, nil
}

func (r *registry) loadBuiltins(builtin embed.FS) error {
	entries, err := fs.ReadDir(builtin, "prompts")
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".md" {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".md")
		bytes, err := builtin.ReadFile(filepath.ToSlash(filepath.Join("prompts", e.Name())))
		if err != nil {
			return err
		}
		label, body, err := ParseFrontmatter(name, string(bytes))
		if err != nil {
			return err
		}
		r.prompts[name] = Prompt{Name: name, Label: label, Source: SourceBuiltIn, Body: body}
	}
	return nil
}

func (r *registry) loadUsers() error {
	if err := os.MkdirAll(r.userDir, 0o700); err != nil {
		return err
	}
	entries, err := os.ReadDir(r.userDir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".md" {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".md")
		content, err := os.ReadFile(filepath.Join(r.userDir, e.Name()))
		if err != nil {
			return err
		}
		label, body, err := ParseFrontmatter(name, string(content))
		if err != nil {
			return err
		}
		r.prompts[name] = Prompt{Name: name, Label: label, Source: SourceUser, FilePath: filepath.Join(r.userDir, e.Name()), Body: body}
	}
	return nil
}
