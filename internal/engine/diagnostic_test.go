package engine_test

import (
	"encoding/json"
	"testing"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
)

func TestDiagnostic_MarshalJSON_WhenFilled_UsesSpecFieldNames(t *testing.T) {
	diag := engine.Diagnostic{
		RuleID:   "MDL001",
		Severity: engine.SeverityError,
		Path:     "README.md",
		Line:     42,
		Column:   5,
		Message:  "./docs/setup.md does not exist",
	}

	got, err := json.Marshal(diag)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	want := `{"rule":"MDL001","severity":"error","path":"README.md","line":42,"column":5,"message":"./docs/setup.md does not exist"}`
	if string(got) != want {
		t.Fatalf("json.Marshal() = %s, want %s", got, want)
	}
}
