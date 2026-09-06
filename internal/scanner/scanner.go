package scanner

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/takeuchi-shogo/md-linker/internal/config"
	"github.com/takeuchi-shogo/md-linker/internal/repository"
)

// Scan は target 配下の Markdown パスを Repository 相対で列挙する。
func Scan(fs repository.FileSystem, root, target string, cfg config.Config) ([]string, error) {
	relTarget, err := repoRel(root, target)
	if err != nil {
		return nil, err
	}

	if fs.IsFile(relTarget) {
		if !isMarkdown(relTarget) {
			return nil, fmt.Errorf("scanner.Scan: target=%q: not a markdown file", target)
		}
		return []string{relTarget}, nil
	}

	if !fs.IsDir(relTarget) && !fs.Exists(relTarget) {
		return nil, fmt.Errorf("scanner.Scan: target=%q: path not found", target)
	}

	var found []string
	err = walk(fs, relTarget, func(p string) error {
		if !isMarkdown(p) {
			return nil
		}
		if !matchAny(cfg.Include, p) {
			return nil
		}
		if matchAny(cfg.Exclude, p) {
			return nil
		}
		found = append(found, p)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(found)
	return found, nil
}

func walk(fs repository.FileSystem, dir string, visit func(string) error) error {
	entries, err := fs.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("scanner.walk: dir=%q: %w", dir, err)
	}
	for _, entry := range entries {
		child := entry.Name()
		if dir != "." {
			child = path.Join(dir, entry.Name())
		}
		if entry.IsDir() {
			if err := walk(fs, child, visit); err != nil {
				return err
			}
			continue
		}
		if err := visit(child); err != nil {
			return err
		}
	}
	return nil
}

func matchAny(patterns []string, p string) bool {
	for _, pattern := range patterns {
		ok, err := doublestar.Match(pattern, p)
		if err == nil && ok {
			return true
		}
	}
	return false
}

func isMarkdown(p string) bool {
	ext := strings.ToLower(path.Ext(p))
	return ext == ".md" || ext == ".mdx"
}

func repoRel(root, target string) (string, error) {
	if target == "" || target == "." {
		return ".", nil
	}
	normalized := strings.ReplaceAll(target, "\\", "/")
	root = strings.ReplaceAll(root, "\\", "/")
	root = strings.TrimRight(root, "/")

	if path.IsAbs(normalized) && root != "" && root != "." {
		if normalized == root {
			return ".", nil
		}
		if strings.HasPrefix(normalized, root+"/") {
			return strings.TrimPrefix(normalized, root+"/"), nil
		}
	}
	return path.Clean(strings.TrimPrefix(normalized, "./")), nil
}
