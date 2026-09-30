package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

const initNoticePrefix = "<!-- kp: "

// editInitAnswers preloads the prompt template into a temporary Markdown file,
// opens the user's editor, and parses the saved sections. An empty or
// untouched file cancels. Missing or empty sections reopen the same file with
// a notice at the top, so nothing the user wrote is lost.
func (a *app) editInitAnswers() (initAnswers, error) {
	editor, err := a.resolveEditor()
	if err != nil {
		return initAnswers{}, err
	}

	file, err := os.CreateTemp("", "kp-task-*.md")
	if err != nil {
		return initAnswers{}, NewExitError(ExitSystem, err)
	}
	path := file.Name()
	_ = file.Close()
	defer os.Remove(path)

	draft, cursorLine := renderInitDraft()
	for {
		if err := os.WriteFile(path, []byte(draft), 0o600); err != nil {
			return initAnswers{}, NewExitError(ExitSystem, err)
		}
		if err := a.runEditorAt(editor, path, cursorLine); err != nil {
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				return initAnswers{}, NewExitError(ExitCancel, errPickerCanceled)
			}
			return initAnswers{}, err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return initAnswers{}, NewExitError(ExitSystem, err)
		}
		answers, problems, empty := parseInitDraft(string(content))
		if empty {
			return initAnswers{}, NewExitError(ExitCancel, errPickerCanceled)
		}
		if len(problems) == 0 {
			return answers, nil
		}
		draft = initNoticePrefix + strings.Join(problems, " ") +
			" Fix this and save, or delete everything to cancel. -->\n" +
			stripInitNotice(string(content))
		cursorLine = 1
	}
}

// renderInitDraft returns the editor template and the line where the first
// answer goes. Guidance lives in HTML comments, which parsing removes.
func renderInitDraft() (string, int) {
	var b strings.Builder
	b.WriteString("<!-- kp task\n")
	b.WriteString("Write under each heading. Blank lines and indentation are kept.\n")
	b.WriteString("Save and quit to print and copy the prompt. Comments like this one\n")
	b.WriteString("are removed. Delete everything to cancel.\n")
	b.WriteString("-->\n\n")

	cursorLine := 0
	for _, field := range initFields {
		b.WriteString(field.heading() + "\n\n")
		b.WriteString("<!-- " + field.question + " -->\n")
		if cursorLine == 0 {
			cursorLine = strings.Count(b.String(), "\n") + 1
		}
		b.WriteString("\n\n")
	}
	b.WriteString("<!-- kp adds this closing line: " + initInvestigationSentence + " -->\n")
	return b.String(), cursorLine
}

// parseInitDraft splits the saved file into sections by heading. It reports
// empty when every section is blank, and problems for missing headings,
// empty sections, or text outside any section.
func parseInitDraft(content string) (answers initAnswers, problems []string, empty bool) {
	bodies := make([][]string, len(initFields))
	seen := make([]bool, len(initFields))
	current := -1
	stray := false
	for _, line := range strings.Split(stripInitComments(content), "\n") {
		line = strings.TrimRight(line, " \t\r")
		if i := initHeadingIndex(line); i >= 0 {
			current, seen[i] = i, true
			continue
		}
		if current < 0 {
			stray = stray || strings.TrimSpace(line) != ""
			continue
		}
		bodies[current] = append(bodies[current], line)
	}

	empty = true
	for i, field := range initFields {
		body := strings.Join(trimBlankLines(bodies[i]), "\n")
		field.set(&answers, body)
		switch {
		case !seen[i]:
			problems = append(problems, fmt.Sprintf("The %q heading is missing.", field.title))
		case body == "":
			problems = append(problems, fmt.Sprintf("%s is empty.", field.title))
		default:
			empty = false
		}
	}
	if stray {
		problems = append(problems, "Move the text above the first heading into a section.")
		empty = false
	}
	return answers, problems, empty
}

// initHeadingIndex matches a section heading with or without its emoji or a
// Markdown "#" prefix, ignoring case.
func initHeadingIndex(line string) int {
	title := strings.TrimLeftFunc(line, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	for i, field := range initFields {
		if strings.EqualFold(strings.TrimSpace(title), field.title) {
			return i
		}
	}
	return -1
}

// stripInitComments removes <!-- ... --> spans, dropping lines that held only
// comments so guidance does not leave extra blank lines behind.
func stripInitComments(content string) string {
	var out []string
	inComment := false
	for _, line := range strings.Split(content, "\n") {
		var kept strings.Builder
		rest, touched := line, inComment
		for rest != "" {
			if inComment {
				end := strings.Index(rest, "-->")
				if end < 0 {
					rest = ""
					break
				}
				rest, inComment = rest[end+3:], false
				continue
			}
			start := strings.Index(rest, "<!--")
			if start < 0 {
				kept.WriteString(rest)
				break
			}
			kept.WriteString(rest[:start])
			rest, inComment, touched = rest[start+4:], true, true
		}
		if touched && strings.TrimSpace(kept.String()) == "" {
			continue
		}
		out = append(out, kept.String())
	}
	return strings.Join(out, "\n")
}

func stripInitNotice(content string) string {
	for strings.HasPrefix(content, initNoticePrefix) {
		_, rest, ok := strings.Cut(content, "-->\n")
		if !ok {
			break
		}
		content = rest
	}
	return content
}

func trimBlankLines(lines []string) []string {
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// runEditorAt opens path, starting vi-family editors on line.
func (a *app) runEditorAt(editor editorCommand, path string, line int) error {
	switch filepath.Base(editor.name) {
	case "vi", "vim", "nvim":
		if line > 0 {
			editor.args = append(append([]string{}, editor.args...), "+"+strconv.Itoa(line))
		}
	}
	return a.runEditor(editor, path)
}
