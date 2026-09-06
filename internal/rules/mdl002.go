package rules

import (
	"fmt"
	"strings"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
	"github.com/takeuchi-shogo/md-linker/internal/repository"
)

// MDL002 は同一 / 別 Markdown の存在しないアンカーを検出する。
type MDL002 struct{}

func (MDL002) ID() string { return "MDL002" }

func (MDL002) Check(ctx *engine.Context, doc *engine.Document) []engine.Diagnostic {
	var diags []engine.Diagnostic
	for _, link := range doc.Links {
		if link.Fragment == "" {
			continue
		}
		if isExternalScheme(link.Dest) {
			continue
		}

		targetPath := doc.Path
		if link.Dest != "" {
			resolved, err := repository.Resolve(ctx.Root, doc.Path, link.Dest)
			if err != nil || resolved.Outside {
				continue
			}
			targetPath = resolved.RepoPath
			if _, known := ctx.Documents[targetPath]; !known {
				continue
			}
		}

		if hasAnchor(ctx.Anchors[targetPath], link.Fragment) {
			continue
		}
		diags = append(diags, engine.Diagnostic{
			RuleID:  "MDL002",
			Path:    doc.Path,
			Line:    link.Line,
			Column:  link.Column,
			Message: fmt.Sprintf("anchor %q does not exist in %s", link.Fragment, targetPath),
		})
	}
	return diags
}

func hasAnchor(anchors map[string]struct{}, fragment string) bool {
	if _, ok := anchors[fragment]; ok {
		return true
	}
	_, ok := anchors[strings.ToLower(fragment)]
	return ok
}
