package rules_test

import (
	"testing"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
	"github.com/takeuchi-shogo/md-linker/internal/repository"
	"github.com/takeuchi-shogo/md-linker/internal/rules"
)

func TestMDL003_WhenHighConfidencePathMissing_ReturnsMissingFileReference(t *testing.T) {
	ctx := engine.NewContext(".", repository.NewMemFS())
	doc := &engine.Document{
		Path: "README.md",
		Code: []engine.CodeSpan{
			{Text: "internal/auth/service.go", Line: 1, Column: 1},
			{Text: "fmt.Println", Line: 2, Column: 1},
		},
	}

	diags := (rules.MDL003{}).Check(ctx, doc)
	if len(diags) != 1 || diags[0].RuleID != "MDL003" {
		t.Fatalf("Check() = %+v, want 1 MDL003 diagnostic", diags)
	}
}

func TestMDL003_WhenRepoRootPathFromNestedDoc_ResolvesFromRoot(t *testing.T) {
	fs := repository.NewMemFS()
	fs.WriteFile("internal/auth/service.go", []byte("package auth"))
	ctx := engine.NewContext(".", fs)
	doc := &engine.Document{
		Path: "docs/guide.md",
		Code: []engine.CodeSpan{{Text: "internal/auth/service.go", Line: 1, Column: 1}},
	}

	if diags := (rules.MDL003{}).Check(ctx, doc); len(diags) != 0 {
		t.Fatalf("Check() = %+v, want empty; repo-root paths resolve from root", diags)
	}
}

func TestMDL003_WhenRelativePathExists_ReturnsNoDiagnostic(t *testing.T) {
	fs := repository.NewMemFS()
	fs.WriteFile("internal/auth/service.go", []byte("package auth"))
	ctx := engine.NewContext(".", fs)
	doc := &engine.Document{
		Path: "README.md",
		Code: []engine.CodeSpan{{Text: "./internal/auth/service.go", Line: 1, Column: 1}},
	}

	if diags := (rules.MDL003{}).Check(ctx, doc); len(diags) != 0 {
		t.Fatalf("Check() = %+v, want empty", diags)
	}
}

func TestMDL006_WhenUnixHomePath_ReturnsAbsoluteLocalPath(t *testing.T) {
	ctx := engine.NewContext(".", repository.NewMemFS())
	doc := &engine.Document{
		Path: "README.md",
		Code: []engine.CodeSpan{{Text: "/Users/demo/project/main.go", Line: 2, Column: 3}},
	}

	diags := (rules.MDL006{}).Check(ctx, doc)
	if len(diags) != 1 || diags[0].RuleID != "MDL006" {
		t.Fatalf("Check() = %+v, want 1 MDL006 diagnostic", diags)
	}
}

func TestMDL006_WhenWindowsUserPathInLink_ReturnsAbsoluteLocalPath(t *testing.T) {
	ctx := engine.NewContext(".", repository.NewMemFS())
	doc := &engine.Document{
		Path:  "README.md",
		Links: []engine.Link{{Dest: `C:\Users\demo\file.md`, Line: 1, Column: 1}},
	}

	diags := (rules.MDL006{}).Check(ctx, doc)
	if len(diags) != 1 || diags[0].RuleID != "MDL006" {
		t.Fatalf("Check() = %+v, want 1 MDL006 diagnostic", diags)
	}
}
