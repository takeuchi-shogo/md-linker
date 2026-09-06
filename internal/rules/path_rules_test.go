package rules_test

import (
	"testing"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
	"github.com/takeuchi-shogo/md-linker/internal/repository"
	"github.com/takeuchi-shogo/md-linker/internal/rules"
)

func TestMDL001_WhenRelativeLinkMissing_ReturnsBrokenLink(t *testing.T) {
	fs := repository.NewMemFS()
	fs.WriteFile("README.md", []byte("x"))
	ctx := engine.NewContext(".", fs)
	doc := &engine.Document{
		Path:  "README.md",
		Links: []engine.Link{{Dest: "./docs/setup.md", Line: 2, Column: 5}},
	}

	diags := rules.MDL001{}.Check(ctx, doc)
	if len(diags) != 1 || diags[0].RuleID != "MDL001" {
		t.Fatalf("Check() = %+v, want 1 MDL001 diagnostic", diags)
	}
}

func TestMDL001_WhenHTTPLink_IgnoresExternalScheme(t *testing.T) {
	ctx := engine.NewContext(".", repository.NewMemFS())
	doc := &engine.Document{
		Path:  "README.md",
		Links: []engine.Link{{Dest: "https://example.com"}},
	}

	if diags := (rules.MDL001{}).Check(ctx, doc); len(diags) != 0 {
		t.Fatalf("Check() = %+v, want empty for http(s)", diags)
	}
}

func TestMDL001_WhenTargetExists_ReturnsNoDiagnostic(t *testing.T) {
	fs := repository.NewMemFS()
	fs.WriteFile("docs/setup.md", []byte("# Setup"))
	ctx := engine.NewContext(".", fs)
	doc := &engine.Document{
		Path:  "README.md",
		Links: []engine.Link{{Dest: "./docs/setup.md", Line: 1, Column: 1}},
	}

	if diags := (rules.MDL001{}).Check(ctx, doc); len(diags) != 0 {
		t.Fatalf("Check() = %+v, want empty when file exists", diags)
	}
}

func TestMDL004_WhenImageMissing_ReturnsMissingImage(t *testing.T) {
	ctx := engine.NewContext(".", repository.NewMemFS())
	doc := &engine.Document{
		Path:   "README.md",
		Images: []engine.Image{{Dest: "./images/logo.png", Line: 3, Column: 1}},
	}

	diags := rules.MDL004{}.Check(ctx, doc)
	if len(diags) != 1 || diags[0].RuleID != "MDL004" {
		t.Fatalf("Check() = %+v, want 1 MDL004 diagnostic", diags)
	}
}

func TestMDL007_WhenCaseDiffers_ReturnsCaseMismatch(t *testing.T) {
	fs := repository.NewMemFS()
	fs.WriteFile("Docs/Setup.md", []byte("x"))
	ctx := engine.NewContext(".", fs)
	doc := &engine.Document{
		Path:  "README.md",
		Links: []engine.Link{{Dest: "./docs/setup.md", Line: 1, Column: 1}},
	}

	diags := rules.MDL007{}.Check(ctx, doc)
	if len(diags) != 1 || diags[0].RuleID != "MDL007" {
		t.Fatalf("Check() = %+v, want 1 MDL007 diagnostic", diags)
	}
}

func TestMDL001_WhenOnlyCaseDiffers_DoesNotReportBrokenLink(t *testing.T) {
	fs := repository.NewMemFS()
	fs.WriteFile("Docs/Setup.md", []byte("x"))
	ctx := engine.NewContext(".", fs)
	doc := &engine.Document{
		Path:  "README.md",
		Links: []engine.Link{{Dest: "./docs/setup.md", Line: 1, Column: 1}},
	}

	if diags := (rules.MDL001{}).Check(ctx, doc); len(diags) != 0 {
		t.Fatalf("Check() = %+v, want empty; case mismatch is MDL007", diags)
	}
}
