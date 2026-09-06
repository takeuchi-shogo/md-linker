package app_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/takeuchi-shogo/md-linker/internal/app"
	"github.com/takeuchi-shogo/md-linker/internal/config"
	"github.com/takeuchi-shogo/md-linker/internal/engine"
	"github.com/takeuchi-shogo/md-linker/internal/repository"
)

func TestFixtures_WhenBrokenLink_ReportsMDL001(t *testing.T) {
	assertFixtureRule(t, "broken-link", "MDL001")
}

func TestFixtures_WhenBrokenAnchor_ReportsMDL002(t *testing.T) {
	assertFixtureRule(t, "broken-anchor", "MDL002")
}

func TestFixtures_WhenMissingImage_ReportsMDL004(t *testing.T) {
	assertFixtureRule(t, "missing-image", "MDL004")
}

func TestFixtures_WhenMissingFileReference_ReportsMDL003(t *testing.T) {
	assertFixtureRule(t, "file-reference", "MDL003")
}

func TestFixtures_WhenAbsolutePath_ReportsMDL006(t *testing.T) {
	assertFixtureRule(t, "absolute-path", "MDL006")
}

func TestFixtures_WhenCaseMismatch_ReportsMDL007(t *testing.T) {
	assertFixtureRule(t, "case-mismatch", "MDL007")
}

func TestFixtures_WhenCheckRuns_DoesNotModifyRepositoryFiles(t *testing.T) {
	root := repoRoot(t)
	target := filepath.Join("testdata", "fixtures", "broken-link")
	readme := filepath.Join(root, target, "README.md")
	before := fileHash(t, readme)

	var buf bytes.Buffer
	_, err := app.Run(repository.NewOSFS(root), root, target, config.Default(), "text", &buf)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	after := fileHash(t, readme)
	if before != after {
		t.Fatal("README.md hash changed after check; tool must be read-only")
	}
}

func assertFixtureRule(t *testing.T, name, ruleID string) {
	t.Helper()
	root := repoRoot(t)
	target := filepath.Join("testdata", "fixtures", name)

	var buf bytes.Buffer
	_, err := app.Run(repository.NewOSFS(root), root, target, config.Default(), "json", &buf)
	if err != nil {
		t.Fatalf("Run(%s) error = %v", name, err)
	}

	var payload struct {
		Diagnostics []engine.Diagnostic `json:"diagnostics"`
	}
	if uerr := json.Unmarshal(buf.Bytes(), &payload); uerr != nil {
		t.Fatalf("json.Unmarshal(%s) error = %v body=%s", name, uerr, buf.String())
	}

	for _, diag := range payload.Diagnostics {
		if diag.RuleID == ruleID {
			return
		}
	}
	t.Fatalf("fixture %s: missing rule %s in %+v", name, ruleID, payload.Diagnostics)
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("go.mod not found")
	return ""
}

func fileHash(t *testing.T, path string) [32]byte {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		t.Fatal(err)
	}
	var sum [32]byte
	copy(sum[:], h.Sum(nil))
	return sum
}
