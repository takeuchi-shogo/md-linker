package engine_test

import (
	"testing"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
	"github.com/takeuchi-shogo/md-linker/internal/repository"
)

type stubRule struct {
	id string
}

func (s stubRule) ID() string { return s.id }

func (s stubRule) Check(_ *engine.Context, doc *engine.Document) []engine.Diagnostic {
	return []engine.Diagnostic{{
		RuleID:  s.id,
		Path:    doc.Path,
		Line:    1,
		Column:  1,
		Message: "x",
	}}
}

func TestEngineRun_WhenRuleOff_SkipsDiagnostics(t *testing.T) {
	eng := engine.New([]engine.Rule{stubRule{id: "MDL001"}})
	ctx := engine.NewContext(".", repository.NewMemFS())
	docs := []*engine.Document{{Path: "README.md"}}

	got := eng.Run(ctx, docs, map[string]engine.Severity{
		"MDL001": engine.SeverityOff,
	})
	if len(got) != 0 {
		t.Fatalf("Run() = %+v, want empty when rule is off", got)
	}
}

func TestEngineRun_WhenRuleError_OverridesSeverityAndKeepsDiagnostic(t *testing.T) {
	eng := engine.New([]engine.Rule{stubRule{id: "MDL001"}})
	ctx := engine.NewContext(".", repository.NewMemFS())
	docs := []*engine.Document{{Path: "README.md"}}

	got := eng.Run(ctx, docs, map[string]engine.Severity{
		"MDL001": engine.SeverityError,
	})
	if len(got) != 1 {
		t.Fatalf("Run() len = %d, want 1", len(got))
	}
	if got[0].Severity != engine.SeverityError {
		t.Fatalf("Severity = %q, want %q", got[0].Severity, engine.SeverityError)
	}
}

func TestEngineRun_WhenMultipleDocuments_SortsByPathLineColumnRule(t *testing.T) {
	eng := engine.New([]engine.Rule{stubRule{id: "MDL001"}})
	ctx := engine.NewContext(".", repository.NewMemFS())
	docs := []*engine.Document{{Path: "b.md"}, {Path: "a.md"}}

	got := eng.Run(ctx, docs, map[string]engine.Severity{
		"MDL001": engine.SeverityWarning,
	})
	if len(got) != 2 {
		t.Fatalf("Run() len = %d, want 2", len(got))
	}
	if got[0].Path != "a.md" || got[1].Path != "b.md" {
		t.Fatalf("sorted paths = %q, %q, want a.md then b.md", got[0].Path, got[1].Path)
	}
}

func TestContextRegister_WhenHeadingsPresent_IndexesAnchorsByPath(t *testing.T) {
	ctx := engine.NewContext(".", repository.NewMemFS())
	doc := &engine.Document{
		Path:     "README.md",
		Headings: []engine.Heading{{Text: "Setup", Anchor: "setup"}},
	}
	ctx.Register([]*engine.Document{doc})

	if ctx.Documents["README.md"] != doc {
		t.Fatal("Documents[README.md] was not registered")
	}
	if _, ok := ctx.Anchors["README.md"]["setup"]; !ok {
		t.Fatal("Anchors[README.md][setup] missing")
	}
}
