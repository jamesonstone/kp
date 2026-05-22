package prompt

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Source int

const (
	SourceBuiltIn Source = iota
	SourceUser
)

type Prompt struct {
	Name     string
	Label    string
	Source   Source
	FilePath string
	Body     string
}

type Registry interface {
	List() []Prompt
	Get(name string) (Prompt, error)
	Add(name, body string) (Prompt, error)
	Remove(name string) error
	PromoteToUser(name string) (Prompt, error)
}

var (
	ErrNotFound = errors.New("prompt not found")
	ErrEmpty    = errors.New("prompt body is empty")
	ErrBuiltIn  = errors.New("cannot modify built-in prompt; promote first")
	ErrExists   = errors.New("prompt already exists")
)

var validName = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

func ValidateName(name string) error {
	if !validName.MatchString(name) {
		return fmt.Errorf("invalid prompt name %q", name)
	}
	return nil
}

func DefaultLabel(name string) string {
	parts := strings.Split(name, "-")
	for i := range parts {
		if parts[i] == "" {
			continue
		}
		parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
	}
	return strings.Join(parts, " ")
}
