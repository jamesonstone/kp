package prompt

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Source int

const (
	SourceBuiltIn Source = iota
	SourceUser
)

func (s Source) String() string {
	switch s {
	case SourceBuiltIn:
		return "builtin"
	case SourceUser:
		return "user"
	default:
		return "unknown"
	}
}

type Prompt struct {
	Name     string
	Label    string
	Source   Source
	FilePath string
	Body     string
}

var (
	ErrNotFound     = errors.New("prompt not found")
	ErrEmpty        = errors.New("prompt body is empty")
	ErrBuiltIn      = errors.New("cannot modify built-in prompt")
	ErrExists       = errors.New("prompt already exists")
	ErrInvalidName  = errors.New("invalid prompt name")
	ErrReservedName = errors.New("reserved prompt name")
)

var namePattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

var reservedNames = map[string]struct{}{
	"edit":     {},
	"help":     {},
	"list":     {},
	"new":      {},
	"prompt":   {},
	"rm":       {},
	"scaffold": {},
	"version":  {},
}

func ValidateName(name string) error {
	if !namePattern.MatchString(name) {
		return fmt.Errorf("%w: %s", ErrInvalidName, name)
	}
	if IsReservedName(name) {
		return fmt.Errorf("%w: %s", ErrReservedName, name)
	}
	return nil
}

func IsReservedName(name string) bool {
	_, ok := reservedNames[name]
	return ok
}

func reservedNameList() []string {
	names := make([]string, 0, len(reservedNames))
	for name := range reservedNames {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func defaultLabel(name string) string {
	if name == "" {
		return ""
	}

	parts := strings.Split(name, "-")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}
