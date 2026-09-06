package markdown_test

import (
	"testing"

	"github.com/takeuchi-shogo/md-linker/internal/markdown"
)

func TestParse_WhenLinkHeadingAndInlineCode_ExtractsDocumentModel(t *testing.T) {
	src := []byte("# Setup\n\nSee [guide](./guide.md#install) and `internal/app.go`.\n")

	doc, err := markdown.Parse("README.md", src)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if doc.Path != "README.md" {
		t.Fatalf("Path = %q, want README.md", doc.Path)
	}
	if len(doc.Links) != 1 {
		t.Fatalf("len(Links) = %d, want 1", len(doc.Links))
	}
	if doc.Links[0].Dest != "./guide.md" || doc.Links[0].Fragment != "install" {
		t.Fatalf("Link = %+v, want dest=./guide.md fragment=install", doc.Links[0])
	}
	if doc.Links[0].Line < 1 {
		t.Fatalf("Link.Line = %d, want >= 1", doc.Links[0].Line)
	}
	if len(doc.Headings) != 1 || doc.Headings[0].Text != "Setup" {
		t.Fatalf("Headings = %+v, want Setup", doc.Headings)
	}
	if len(doc.Code) != 1 || doc.Code[0].Text != "internal/app.go" {
		t.Fatalf("Code = %+v, want internal/app.go", doc.Code)
	}
}

func TestParse_WhenFencedCodeContainsPath_DoesNotAddCodeSpan(t *testing.T) {
	src := []byte("```\ninternal/secret.go\n```\n")

	doc, err := markdown.Parse("README.md", src)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(doc.Code) != 0 {
		t.Fatalf("Code = %+v, want empty for fenced block", doc.Code)
	}
}

func TestParse_WhenImagePresent_ExtractsImageDest(t *testing.T) {
	src := []byte("![alt](./images/logo.png)\n")

	doc, err := markdown.Parse("README.md", src)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(doc.Images) != 1 || doc.Images[0].Dest != "./images/logo.png" {
		t.Fatalf("Images = %+v, want ./images/logo.png", doc.Images)
	}
}
