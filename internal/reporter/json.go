package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/takeuchi-shogo/md-linker/internal/engine"
)

type jsonReport struct {
	Diagnostics []engine.Diagnostic `json:"diagnostics"`
}

func writeJSON(w io.Writer, diags []engine.Diagnostic) error {
	if diags == nil {
		diags = []engine.Diagnostic{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(jsonReport{Diagnostics: diags}); err != nil {
		return fmt.Errorf("reporter.writeJSON: %w", err)
	}
	return nil
}
