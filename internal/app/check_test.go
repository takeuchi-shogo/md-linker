package app_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/takeuchi-shogo/md-linker/internal/app"
	"github.com/takeuchi-shogo/md-linker/internal/config"
	"github.com/takeuchi-shogo/md-linker/internal/engine"
	"github.com/takeuchi-shogo/md-linker/internal/repository"
)

func TestRun_WhenBrokenRelativeLink_ReturnsJSONMDL001AndExit1(t *testing.T) {
	fs := repository.NewMemFS()
	fs.WriteFile("README.md", []byte("See [setup](./docs/setup.md).\n"))

	var buf bytes.Buffer
	code, err := app.Run(fs, ".", ".", config.Default(), "json", &buf)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if code != 1 {
		t.Fatalf("Run() code = %d, want 1", code)
	}

	var payload struct {
		Diagnostics []engine.Diagnostic `json:"diagnostics"`
	}
	if uerr := json.Unmarshal(buf.Bytes(), &payload); uerr != nil {
		t.Fatalf("json.Unmarshal() error = %v body=%s", uerr, buf.String())
	}
	if len(payload.Diagnostics) != 1 || payload.Diagnostics[0].RuleID != "MDL001" {
		t.Fatalf("diagnostics = %+v, want 1 MDL001", payload.Diagnostics)
	}
}

func TestRun_WhenCleanRepo_ReturnsExit0(t *testing.T) {
	fs := repository.NewMemFS()
	fs.WriteFile("README.md", []byte("# Hello\n"))

	var buf bytes.Buffer
	code, err := app.Run(fs, ".", ".", config.Default(), "text", &buf)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if code != 0 {
		t.Fatalf("Run() code = %d, want 0 body=%s", code, buf.String())
	}
	if strings.TrimSpace(buf.String()) != "" {
		t.Fatalf("Run() stdout = %q, want empty", buf.String())
	}
}
