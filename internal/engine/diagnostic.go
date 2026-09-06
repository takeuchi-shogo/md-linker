package engine

// Severity は診断の重大度。設定ファイルと JSON 出力の値と一致させる。
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityOff     Severity = "off"
)

// Diagnostic は 1 件の検査結果。JSON タグは要件書の schema に合わせる。
type Diagnostic struct {
	RuleID     string   `json:"rule"`
	Severity   Severity `json:"severity"`
	Path       string   `json:"path"`
	Line       int      `json:"line"`
	Column     int      `json:"column"`
	Message    string   `json:"message"`
	Suggestion string   `json:"suggestion,omitempty"`
}
