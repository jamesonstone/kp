package scaffold

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Action string

const (
	ActionCreate      Action = "create"
	ActionUpdate      Action = "update"
	ActionSkip        Action = "skip"
	ActionWouldCreate Action = "would-create"
	ActionWouldUpdate Action = "would-update"
	ActionWouldSkip   Action = "would-skip"
)

type Options struct {
	Dir    string
	DryRun bool
	Force  bool
}

type Result struct {
	Path   string
	Action Action
}

type artifact struct {
	path     string
	template string
	mode     fsMode
}

type fsMode os.FileMode

const (
	fileMode fsMode = 0o644
	envMode  fsMode = 0o600
)

func Run(opts Options) ([]Result, error) {
	root, err := targetRoot(opts.Dir)
	if err != nil {
		return nil, err
	}

	artifacts, err := scaffoldArtifacts()
	if err != nil {
		return nil, err
	}

	results := make([]Result, 0, len(artifacts)+1)
	gitignoreResult, err := applyGitignore(root, opts)
	if err != nil {
		return nil, err
	}
	results = append(results, gitignoreResult)

	for _, artifact := range artifacts {
		result, err := applyArtifact(root, artifact, opts)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func targetRoot(dir string) (string, error) {
	if strings.TrimSpace(dir) == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("get working directory: %w", err)
		}
		dir = cwd
	}

	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("resolve scaffold target %q: %w", dir, err)
	}
	return abs, nil
}

func scaffoldArtifacts() ([]artifact, error) {
	specs := []struct {
		path     string
		template string
		mode     fsMode
	}{
		{path: ".env", template: "", mode: envMode},
		{path: ".envrc", template: "envrc", mode: fileMode},
		{path: ".coderabbit.yaml", template: "coderabbit.yaml", mode: fileMode},
		{path: ".github/pull_request_template.md", template: "github/pull_request_template.md", mode: fileMode},
		{path: "AGENTS.md", template: "AGENTS.md", mode: fileMode},
		{path: "CLAUDE.md", template: "CLAUDE.md", mode: fileMode},
		{path: ".github/copilot-instructions.md", template: "github/copilot-instructions.md", mode: fileMode},
		{path: "docs/agents/README.md", template: "docs/agents/README.md", mode: fileMode},
		{path: "docs/agents/WORKFLOWS.md", template: "docs/agents/WORKFLOWS.md", mode: fileMode},
		{path: "docs/agents/RLM.md", template: "docs/agents/RLM.md", mode: fileMode},
		{path: "docs/agents/TOOLING.md", template: "docs/agents/TOOLING.md", mode: fileMode},
		{path: "docs/agents/GUARDRAILS.md", template: "docs/agents/GUARDRAILS.md", mode: fileMode},
		{path: "docs/references/README.md", template: "docs/references/README.md", mode: fileMode},
		{path: "docs/references/testing.md", template: "docs/references/testing.md", mode: fileMode},
		{path: "docs/references/tooling.md", template: "docs/references/tooling.md", mode: fileMode},
		{path: "docs/references/external-systems.md", template: "docs/references/external-systems.md", mode: fileMode},
	}

	artifacts := make([]artifact, 0, len(specs))
	for _, spec := range specs {
		content := ""
		if spec.template != "" {
			templateContent, err := readTemplate(spec.template)
			if err != nil {
				return nil, err
			}
			content = templateContent
		}
		artifacts = append(artifacts, artifact{
			path:     spec.path,
			template: content,
			mode:     spec.mode,
		})
	}
	return artifacts, nil
}

func applyArtifact(root string, artifact artifact, opts Options) (Result, error) {
	path := filepath.Join(root, filepath.FromSlash(artifact.path))
	exists, err := fileExists(path)
	if err != nil {
		return Result{}, err
	}

	if opts.DryRun {
		if exists && opts.Force {
			return Result{Path: artifact.path, Action: ActionWouldUpdate}, nil
		}
		if exists {
			return Result{Path: artifact.path, Action: ActionWouldSkip}, nil
		}
		return Result{Path: artifact.path, Action: ActionWouldCreate}, nil
	}

	if exists && !opts.Force {
		return Result{Path: artifact.path, Action: ActionSkip}, nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return Result{}, fmt.Errorf("create scaffold directory %q: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(artifact.template), os.FileMode(artifact.mode)); err != nil {
		return Result{}, fmt.Errorf("write scaffold file %q: %w", path, err)
	}
	if exists {
		return Result{Path: artifact.path, Action: ActionUpdate}, nil
	}
	return Result{Path: artifact.path, Action: ActionCreate}, nil
}

func applyGitignore(root string, opts Options) (Result, error) {
	const relativePath = ".gitignore"
	path := filepath.Join(root, relativePath)
	content, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return Result{}, fmt.Errorf("read scaffold file %q: %w", path, err)
		}
		if opts.DryRun {
			return Result{Path: relativePath, Action: ActionWouldCreate}, nil
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return Result{}, fmt.Errorf("create scaffold directory %q: %w", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, []byte(gitignoreTemplate), 0o644); err != nil {
			return Result{}, fmt.Errorf("write scaffold file %q: %w", path, err)
		}
		return Result{Path: relativePath, Action: ActionCreate}, nil
	}

	missing := missingGitignorePatterns(string(content))
	if len(missing) == 0 {
		if opts.DryRun {
			return Result{Path: relativePath, Action: ActionWouldSkip}, nil
		}
		return Result{Path: relativePath, Action: ActionSkip}, nil
	}
	if opts.DryRun {
		return Result{Path: relativePath, Action: ActionWouldUpdate}, nil
	}

	updated := appendGitignorePatterns(string(content), missing)
	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		return Result{}, fmt.Errorf("write scaffold file %q: %w", path, err)
	}
	return Result{Path: relativePath, Action: ActionUpdate}, nil
}

func missingGitignorePatterns(content string) []string {
	existing := make(map[string]struct{})
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		existing[line] = struct{}{}
	}

	missing := make([]string, 0, len(gitignorePatterns))
	for _, pattern := range gitignorePatterns {
		if _, ok := existing[pattern]; !ok {
			missing = append(missing, pattern)
		}
	}
	return missing
}

func appendGitignorePatterns(content string, patterns []string) string {
	var builder strings.Builder
	builder.WriteString(strings.TrimRight(content, "\n"))
	if strings.TrimSpace(content) != "" {
		builder.WriteString("\n\n")
	}
	builder.WriteString("# Kit local generated environment, cache, and scratch artifacts\n")
	for _, pattern := range patterns {
		builder.WriteString(pattern)
		builder.WriteString("\n")
	}
	return builder.String()
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, fmt.Errorf("stat scaffold file %q: %w", path, err)
}
