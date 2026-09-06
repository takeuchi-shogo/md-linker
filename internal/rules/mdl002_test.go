package rules_test

import (
	"testing"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
	"github.com/takeuchi-shogo/md-linker/internal/repository"
	"github.com/takeuchi-shogo/md-linker/internal/rules"
)

func TestMDL002_WhenSameFileAnchorMissing_ReturnsBrokenAnchor(t *testing.T) {
	doc := &engine.Document{
		Path:     "README.md",
		Headings: []engine.Heading{{Text: "Setup", Anchor: "setup"}},
		Links:    []engine.Link{{Dest: "", Fragment: "missing", Line: 3, Column: 1}},
	}
	ctx := engine.NewContext(".", repository.NewMemFS())
	ctx.Register([]*engine.Document{doc})

	diags := (rules.MDL002{}).Check(ctx, doc)
	if len(diags) != 1 || diags[0].RuleID != "MDL002" {
		t.Fatalf("Check() = %+v, want 1 MDL002 diagnostic", diags)
	}
}

func TestMDL002_WhenSameFileAnchorExists_ReturnsNoDiagnostic(t *testing.T) {
	doc := &engine.Document{
		Path:     "README.md",
		Headings: []engine.Heading{{Text: "Setup", Anchor: "setup"}},
		Links:    []engine.Link{{Dest: "", Fragment: "setup", Line: 3, Column: 1}},
	}
	ctx := engine.NewContext(".", repository.NewMemFS())
	ctx.Register([]*engine.Document{doc})

	if diags := (rules.MDL002{}).Check(ctx, doc); len(diags) != 0 {
		t.Fatalf("Check() = %+v, want empty", diags)
	}
}

func TestMDL002_WhenOtherFileMissing_DoesNotReportAnchor(t *testing.T) {
	doc := &engine.Document{
		Path:  "README.md",
		Links: []engine.Link{{Dest: "./gone.md", Fragment: "setup", Line: 1, Column: 1}},
	}
	ctx := engine.NewContext(".", repository.NewMemFS())
	ctx.Register([]*engine.Document{doc})

	if diags := (rules.MDL002{}).Check(ctx, doc); len(diags) != 0 {
		t.Fatalf("Check() = %+v, want empty; missing file is MDL001", diags)
	}
}

func TestMDL002_WhenOtherFileAnchorMissing_ReturnsBrokenAnchor(t *testing.T) {
	readme := &engine.Document{
		Path:  "README.md",
		Links: []engine.Link{{Dest: "./guide.md", Fragment: "missing", Line: 4, Column: 2}},
	}
	guide := &engine.Document{
		Path:     "guide.md",
		Headings: []engine.Heading{{Text: "Setup", Anchor: "setup"}},
	}
	fs := repository.NewMemFS()
	fs.WriteFile("guide.md", []byte("# Setup"))
	ctx := engine.NewContext(".", fs)
	ctx.Register([]*engine.Document{readme, guide})

	diags := (rules.MDL002{}).Check(ctx, readme)
	if len(diags) != 1 || diags[0].RuleID != "MDL002" {
		t.Fatalf("Check() = %+v, want 1 MDL002 diagnostic", diags)
	}
}
