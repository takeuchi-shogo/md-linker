package rules

import (
	"strings"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
	"github.com/takeuchi-shogo/md-linker/internal/repository"
)

func isExternalScheme(dest string) bool {
	lower := strings.ToLower(dest)
	return strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "mailto:")
}

func resolveLocal(ctx *engine.Context, docPath, dest string) (repository.ResolvedPath, bool) {
	if dest == "" || isExternalScheme(dest) {
		return repository.ResolvedPath{}, false
	}
	resolved, err := repository.Resolve(ctx.Root, docPath, dest)
	if err != nil || resolved.Outside {
		return repository.ResolvedPath{}, false
	}
	return resolved, true
}

func hasCaseMismatch(fs repository.FileSystem, repoPath string) bool {
	if repoPath == "" || repoPath == "." {
		return false
	}
	parts := strings.Split(repoPath, "/")
	current := "."
	for _, part := range parts {
		entries, err := fs.ReadDir(current)
		if err != nil {
			return false
		}
		matched := ""
		for _, entry := range entries {
			if strings.EqualFold(entry.Name(), part) {
				matched = entry.Name()
				break
			}
		}
		if matched == "" {
			return false
		}
		if matched != part {
			return true
		}
		if current == "." {
			current = matched
		} else {
			current = current + "/" + matched
		}
	}
	return false
}
