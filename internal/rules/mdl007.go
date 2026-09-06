package rules

import (
	"fmt"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
)

// MDL007 は参照パスと実ファイル名の大小文字不一致を検出する。
type MDL007 struct{}

func (MDL007) ID() string { return "MDL007" }

func (MDL007) Check(ctx *engine.Context, doc *engine.Document) []engine.Diagnostic {
	var diags []engine.Diagnostic
	for _, link := range doc.Links {
		resolved, ok := resolveLocal(ctx, doc.Path, link.Dest)
		if !ok {
			continue
		}
		if !hasCaseMismatch(ctx.FS, resolved.RepoPath) {
			continue
		}
		diags = append(diags, engine.Diagnostic{
			RuleID:  "MDL007",
			Path:    doc.Path,
			Line:    link.Line,
			Column:  link.Column,
			Message: fmt.Sprintf("%q differs in case from the file on disk", link.Dest),
		})
	}
	for _, image := range doc.Images {
		resolved, ok := resolveLocal(ctx, doc.Path, image.Dest)
		if !ok {
			continue
		}
		if !hasCaseMismatch(ctx.FS, resolved.RepoPath) {
			continue
		}
		diags = append(diags, engine.Diagnostic{
			RuleID:  "MDL007",
			Path:    doc.Path,
			Line:    image.Line,
			Column:  image.Column,
			Message: fmt.Sprintf("%q differs in case from the file on disk", image.Dest),
		})
	}
	return diags
}
