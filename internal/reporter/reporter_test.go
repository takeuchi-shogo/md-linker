package reporter_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
	"github.com/takeuchi-shogo/md-linker/internal/reporter"
)

func sampleDiag() engine.Diagnostic {
	return engine.Diagnostic{
		RuleID:   "MDL001",
		Severity: engine.SeverityError,
		Path:     "README.md",
		Line:     42,
		Column:   5,
		Message:  `"./docs/setup.md" does not exist`,
	}
}

func TestExitCode_WhenNoDiagnosticsAndNoError_Returns0(t *testing.T) {
	if got := reporter.ExitCode(nil, nil); got != 0 {
		t.Fatalf("ExitCode() = %d, want 0", got)
	}
}

func TestExitCode_WhenWarningOnly_Returns0(t *testing.T) {
	diags := []engine.Diagnostic{{Severity: engine.SeverityWarning}}
	if got := reporter.ExitCode(diags, nil); got != 0 {
		t.Fatalf("ExitCode() = %d, want 0", got)
	}
}

func TestExitCode_WhenErrorDiagnostic_Returns1(t *testing.T) {
	diags := []engine.Diagnostic{{Severity: engine.SeverityError}}
	if got := reporter.ExitCode(diags, nil); got != 1 {
		t.Fatalf("ExitCode() = %d, want 1", got)
	}
}

func TestExitCode_WhenRunError_Returns2(t *testing.T) {
	if got := reporter.ExitCode(nil, errors.New("boom")); got != 2 {
		t.Fatalf("ExitCode() = %d, want 2", got)
	}
}

func TestWrite_WhenTextFormat_UsesHumanReadableLine(t *testing.T) {
	var buf bytes.Buffer
	if err := reporter.Write(&buf, "text", []engine.Diagnostic{sampleDiag()}); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	want := `README.md:42:5 error MDL001: "./docs/setup.md" does not exist`
	if strings.TrimSpace(buf.String()) != want {
		t.Fatalf("Write() = %q, want %q", buf.String(), want)
	}
}

func TestWrite_WhenJSONFormat_MatchesSpecSchema(t *testing.T) {
	var buf bytes.Buffer
	if err := reporter.Write(&buf, "json", []engine.Diagnostic{sampleDiag()}); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	var payload struct {
		Diagnostics []engine.Diagnostic `json:"diagnostics"`
	}
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v body=%s", err, buf.String())
	}
	if len(payload.Diagnostics) != 1 || payload.Diagnostics[0].RuleID != "MDL001" {
		t.Fatalf("payload = %+v, want 1 MDL001", payload)
	}
}

func TestWrite_WhenGitHubFormat_EmitsAnnotation(t *testing.T) {
	var buf bytes.Buffer
	if err := reporter.Write(&buf, "github", []engine.Diagnostic{sampleDiag()}); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	got := strings.TrimSpace(buf.String())
	want := `::error file=README.md,line=42,col=5::MDL001: "./docs/setup.md" does not exist`
	if got != want {
		t.Fatalf("Write() = %q, want %q", got, want)
	}
}
