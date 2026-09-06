package reporter

import (
	"fmt"
	"io"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
)

func writeText(w io.Writer, diags []engine.Diagnostic) error {
	for _, diag := range diags {
		_, err := fmt.Fprintf(w, "%s:%d:%d %s %s: %s\n",
			diag.Path, diag.Line, diag.Column, diag.Severity, diag.RuleID, diag.Message)
		if err != nil {
			return fmt.Errorf("reporter.writeText: %w", err)
		}
	}
	return nil
}
