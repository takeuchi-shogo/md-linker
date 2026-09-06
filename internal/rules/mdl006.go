package rules

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
)

// MDL006 は端末固有の絶対パスを検出する。
type MDL006 struct{}

func (MDL006) ID() string { return "MDL006" }

func (MDL006) Check(ctx *engine.Context, doc *engine.Document) []engine.Diagnostic {
	var diags []engine.Diagnostic
	for _, span := range doc.Code {
		if !isAbsoluteLocalPath(span.Text) {
			continue
		}
		diags = append(diags, engine.Diagnostic{
			RuleID:  "MDL006",
			Path:    doc.Path,
			Line:    span.Line,
			Column:  span.Column,
			Message: fmt.Sprintf("%q looks like a machine-local absolute path", span.Text),
		})
	}
	for _, link := range doc.Links {
		if !isAbsoluteLocalPath(link.Dest) {
			continue
		}
		diags = append(diags, engine.Diagnostic{
			RuleID:  "MDL006",
			Path:    doc.Path,
			Line:    link.Line,
			Column:  link.Column,
			Message: fmt.Sprintf("%q looks like a machine-local absolute path", link.Dest),
		})
	}
	for _, image := range doc.Images {
		if !isAbsoluteLocalPath(image.Dest) {
			continue
		}
		diags = append(diags, engine.Diagnostic{
			RuleID:  "MDL006",
			Path:    doc.Path,
			Line:    image.Line,
			Column:  image.Column,
			Message: fmt.Sprintf("%q looks like a machine-local absolute path", image.Dest),
		})
	}
	return diags
}

func isAbsoluteLocalPath(s string) bool {
	s = strings.TrimSpace(s)
	switch {
	case strings.HasPrefix(s, "/Users/"),
		strings.HasPrefix(s, "/home/"),
		strings.HasPrefix(s, "/tmp/"):
		return true
	}
	return isWindowsDrivePath(s)
}

func isWindowsDrivePath(s string) bool {
	if len(s) < 3 {
		return false
	}
	if !unicode.IsLetter(rune(s[0])) || s[1] != ':' {
		return false
	}
	return s[2] == '\\' || s[2] == '/'
}
