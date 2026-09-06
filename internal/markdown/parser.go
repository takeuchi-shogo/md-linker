package markdown

import (
	"fmt"
	"strings"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Parse は Markdown を Rule 向け Document Model へ変換する。
func Parse(path string, src []byte) (*engine.Document, error) {
	reader := text.NewReader(src)
	root := goldmark.New().Parser().Parse(reader)

	doc := &engine.Document{
		Path:   path,
		Source: append([]byte(nil), src...),
	}

	err := ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n.Kind() {
		case ast.KindFencedCodeBlock, ast.KindCodeBlock:
			return ast.WalkSkipChildren, nil
		case ast.KindHeading:
			heading := n.(*ast.Heading)
			line, col := positionOf(src, n)
			doc.Headings = append(doc.Headings, engine.Heading{
				Text:   nodeText(n, src),
				Level:  heading.Level,
				Line:   line,
				Column: col,
			})
		case ast.KindLink:
			appendLink(doc, src, n, string(n.(*ast.Link).Destination))
		case ast.KindAutoLink:
			appendLink(doc, src, n, string(n.(*ast.AutoLink).URL(src)))
		case ast.KindImage:
			raw := string(n.(*ast.Image).Destination)
			dest, fragment := splitDest(raw)
			line, col := positionOf(src, n)
			doc.Images = append(doc.Images, engine.Image{
				Raw:      raw,
				Dest:     dest,
				Fragment: fragment,
				Line:     line,
				Column:   col,
			})
		case ast.KindCodeSpan:
			line, col := positionOf(src, n)
			doc.Code = append(doc.Code, engine.CodeSpan{
				Text:   nodeText(n, src),
				Line:   line,
				Column: col,
			})
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, fmt.Errorf("markdown.Parse: path=%q: %w", path, err)
	}

	texts := make([]string, 0, len(doc.Headings))
	for _, h := range doc.Headings {
		texts = append(texts, h.Text)
	}
	anchors := UniqueAnchors(texts)
	for i := range doc.Headings {
		doc.Headings[i].Anchor = anchors[i]
	}
	return doc, nil
}

func appendLink(doc *engine.Document, src []byte, n ast.Node, raw string) {
	dest, fragment := splitDest(raw)
	line, col := positionOf(src, n)
	doc.Links = append(doc.Links, engine.Link{
		Raw:      raw,
		Dest:     dest,
		Fragment: fragment,
		Line:     line,
		Column:   col,
	})
}

func splitDest(raw string) (dest, fragment string) {
	dest, fragment, _ = strings.Cut(raw, "#")
	return dest, fragment
}

func nodeText(n ast.Node, src []byte) string {
	var b strings.Builder
	_ = ast.Walk(n, func(child ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if child == n {
			return ast.WalkContinue, nil
		}
		if t, ok := child.(*ast.Text); ok {
			b.Write(t.Segment.Value(src))
		}
		return ast.WalkContinue, nil
	})
	return b.String()
}

func positionOf(src []byte, n ast.Node) (line, column int) {
	offset := nodeOffset(n)
	return offsetToLineCol(src, offset)
}

func nodeOffset(n ast.Node) int {
	if n.Type() == ast.TypeBlock && n.Lines() != nil && n.Lines().Len() > 0 {
		return n.Lines().At(0).Start
	}
	if t, ok := n.(*ast.Text); ok {
		return t.Segment.Start
	}
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if off := nodeOffset(c); off > 0 {
			return off
		}
	}
	return 0
}

func offsetToLineCol(src []byte, offset int) (line, column int) {
	line = 1
	column = 1
	if offset < 0 {
		offset = 0
	}
	if offset > len(src) {
		offset = len(src)
	}
	for i := 0; i < offset; i++ {
		if src[i] == '\n' {
			line++
			column = 1
			continue
		}
		column++
	}
	return line, column
}
