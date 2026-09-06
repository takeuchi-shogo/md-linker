package engine

import "sort"

// Engine は有効な Rule を各 Document に適用する。
type Engine struct {
	rules []Rule
}

// New は指定 Rule を持つ Engine を返す。
func New(rules []Rule) *Engine {
	return &Engine{rules: rules}
}

// Run は全 Document に Rule を適用し、設定 Severity で上書きして安定ソートする。
func (e *Engine) Run(ctx *Context, docs []*Document, severities map[string]Severity) []Diagnostic {
	var out []Diagnostic
	for _, doc := range docs {
		for _, rule := range e.rules {
			severity, ok := severities[rule.ID()]
			if !ok || severity == SeverityOff {
				continue
			}
			for _, diag := range rule.Check(ctx, doc) {
				diag.Severity = severity
				if diag.RuleID == "" {
					diag.RuleID = rule.ID()
				}
				out = append(out, diag)
			}
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		if a.Column != b.Column {
			return a.Column < b.Column
		}
		return a.RuleID < b.RuleID
	})
	return out
}

// Register は全 Document の path とアンカーを Context に載せる。
func (c *Context) Register(docs []*Document) {
	if c.Documents == nil {
		c.Documents = map[string]*Document{}
	}
	if c.Anchors == nil {
		c.Anchors = map[string]map[string]struct{}{}
	}
	for _, doc := range docs {
		c.Documents[doc.Path] = doc
		anchors := map[string]struct{}{}
		for _, heading := range doc.Headings {
			anchors[heading.Anchor] = struct{}{}
		}
		c.Anchors[doc.Path] = anchors
	}
}
