package rules

import (
	"fmt"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
)

// MDL001 は相対ローカルリンク先の欠落を検出する。
type MDL001 struct{}

func (MDL001) ID() string { return "MDL001" }

func (MDL001) Check(ctx *engine.Context, doc *engine.Document) []engine.Diagnostic {
	var diags []engine.Diagnostic
	for _, link := range doc.Links {
		resolved, ok := resolveLocal(ctx, doc.Path, link.Dest)
		if !ok {
			continue
		}
		if ctx.FS.Exists(resolved.RepoPath) {
			continue
		}
		diags = append(diags, engine.Diagnostic{
			RuleID:  "MDL001",
			Path:    doc.Path,
			Line:    link.Line,
			Column:  link.Column,
			Message: fmt.Sprintf("%q does not exist", link.Dest),
		})
	}
	return diags
}
