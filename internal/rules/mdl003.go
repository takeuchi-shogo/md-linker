package rules

import (
	"fmt"
	"path"
	"strings"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
)

var highConfidenceExt = map[string]struct{}{
	".go": {}, ".md": {}, ".mdx": {}, ".sh": {}, ".yml": {}, ".yaml": {},
	".json": {}, ".txt": {}, ".toml": {}, ".rs": {}, ".py": {}, ".js": {},
	".ts": {}, ".tsx": {}, ".css": {}, ".html": {},
}

// MDL003 は inline code 内の高確度ファイル参照の欠落を検出する。
type MDL003 struct{}

func (MDL003) ID() string { return "MDL003" }

func (MDL003) Check(ctx *engine.Context, doc *engine.Document) []engine.Diagnostic {
	var diags []engine.Diagnostic
	for _, span := range doc.Code {
		if !isHighConfidencePath(span.Text) {
			continue
		}
		if existsAsRepoPath(ctx, doc.Path, span.Text) {
			continue
		}
		diags = append(diags, engine.Diagnostic{
			RuleID:  "MDL003",
			Path:    doc.Path,
			Line:    span.Line,
			Column:  span.Column,
			Message: fmt.Sprintf("%q does not exist", span.Text),
		})
	}
	return diags
}

func existsAsRepoPath(ctx *engine.Context, docPath, raw string) bool {
	if strings.HasPrefix(raw, "./") || strings.HasPrefix(raw, "../") {
		resolved, ok := resolveLocal(ctx, docPath, raw)
		return ok && ctx.FS.Exists(resolved.RepoPath)
	}
	resolved, ok := resolveLocal(ctx, ".", raw)
	return ok && ctx.FS.Exists(resolved.RepoPath)
}

func isHighConfidencePath(s string) bool {
	s = strings.TrimSpace(s)
	if strings.ContainsAny(s, " \t") {
		return false
	}
	if strings.HasPrefix(s, "./") || strings.HasPrefix(s, "../") {
		return true
	}
	if !strings.Contains(s, "/") {
		return false
	}
	_, ok := highConfidenceExt[strings.ToLower(path.Ext(s))]
	return ok
}
