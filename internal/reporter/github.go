package reporter

import (
	"fmt"
	"io"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
)

func writeGitHub(w io.Writer, diags []engine.Diagnostic) error {
	for _, diag := range diags {
		level := "warning"
		if diag.Severity == engine.SeverityError {
			level = "error"
		}
		_, err := fmt.Fprintf(w, "::%s file=%s,line=%d,col=%d::%s: %s\n",
			level, diag.Path, diag.Line, diag.Column, diag.RuleID, diag.Message)
		if err != nil {
			return fmt.Errorf("reporter.writeGitHub: %w", err)
		}
	}
	return nil
}
