package rules

import (
	"fmt"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
)

// MDL004 は画像参照先の欠落を検出する。
type MDL004 struct{}

func (MDL004) ID() string { return "MDL004" }

func (MDL004) Check(ctx *engine.Context, doc *engine.Document) []engine.Diagnostic {
	var diags []engine.Diagnostic
	for _, image := range doc.Images {
		resolved, ok := resolveLocal(ctx, doc.Path, image.Dest)
		if !ok {
			continue
		}
		if ctx.FS.Exists(resolved.RepoPath) {
			continue
		}
		diags = append(diags, engine.Diagnostic{
			RuleID:  "MDL004",
			Path:    doc.Path,
			Line:    image.Line,
			Column:  image.Column,
			Message: fmt.Sprintf("%q does not exist", image.Dest),
		})
	}
	return diags
}
