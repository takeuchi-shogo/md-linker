package repository

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

// ResolvedPath は Markdown 参照を Repository 相対パスへ正規化した結果。
type ResolvedPath struct {
	RepoPath string
	Fragment string
	Outside  bool
}

// Resolve は raw 参照を root 相対パスへ解決する。
// fromDoc は Repository 相対の参照元 Markdown パス。
func Resolve(root, fromDoc, raw string) (ResolvedPath, error) {
	dest, fragment := splitFragment(raw)
	decoded, err := url.PathUnescape(dest)
	if err != nil {
		return ResolvedPath{}, fmt.Errorf("repository.Resolve: unescape dest=%q fromDoc=%q: %w", dest, fromDoc, err)
	}
	decoded = strings.ReplaceAll(decoded, "\\", "/")

	if decoded == "" {
		return ResolvedPath{
			RepoPath: normalizeRepoPath(fromDoc),
			Fragment: fragment,
		}, nil
	}

	if isAbsoluteRef(decoded) {
		return resolveAbsolute(root, decoded, fragment), nil
	}

	fromDir := path.Dir(normalizeRepoPath(fromDoc))
	joined := path.Clean(path.Join(fromDir, decoded))
	if isOutsideRepoRel(joined) {
		return ResolvedPath{Outside: true, Fragment: fragment}, nil
	}
	return ResolvedPath{
		RepoPath: strings.TrimPrefix(joined, "/"),
		Fragment: fragment,
	}, nil
}

func splitFragment(raw string) (dest, fragment string) {
	dest, fragment, _ = strings.Cut(raw, "#")
	return dest, fragment
}

func isAbsoluteRef(p string) bool {
	if strings.HasPrefix(p, "/") {
		return true
	}
	// Windows drive path: C:\ or C:/
	return len(p) >= 3 && p[1] == ':' && (p[2] == '/' || p[2] == '\\')
}

func resolveAbsolute(root, absPath, fragment string) ResolvedPath {
	rootClean := path.Clean(strings.ReplaceAll(root, "\\", "/"))
	absClean := path.Clean(absPath)
	if rootClean != "/" && (absClean == rootClean || strings.HasPrefix(absClean, rootClean+"/")) {
		rel := strings.TrimPrefix(absClean, rootClean)
		rel = strings.TrimPrefix(rel, "/")
		if rel == "" {
			rel = "."
		}
		return ResolvedPath{RepoPath: rel, Fragment: fragment}
	}
	return ResolvedPath{Outside: true, Fragment: fragment}
}

func isOutsideRepoRel(p string) bool {
	return p == ".." || strings.HasPrefix(p, "../")
}
