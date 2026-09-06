package reporter

import (
	"fmt"
	"io"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
)

const (
	FormatText   = "text"
	FormatJSON   = "json"
	FormatGitHub = "github"
)

// Write は診断を指定形式で w へ書く。
func Write(w io.Writer, format string, diags []engine.Diagnostic) error {
	switch format {
	case "", FormatText:
		return writeText(w, diags)
	case FormatJSON:
		return writeJSON(w, diags)
	case FormatGitHub:
		return writeGitHub(w, diags)
	default:
		return fmt.Errorf("reporter.Write: unknown format %q", format)
	}
}

// ExitCode は診断と実行エラーからプロセス終了コードを決める。
func ExitCode(diags []engine.Diagnostic, runErr error) int {
	if runErr != nil {
		return 2
	}
	for _, diag := range diags {
		if diag.Severity == engine.SeverityError {
			return 1
		}
	}
	return 0
}
