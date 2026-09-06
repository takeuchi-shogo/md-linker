package markdown

import (
	"fmt"
	"strings"
	"unicode"
)

// GitHubAnchor は GitHub 互換の見出し fragment を生成する。
func GitHubAnchor(heading string) string {
	heading = strings.ToLower(heading)
	var b strings.Builder
	for _, r := range heading {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			b.WriteRune(r)
		case r == ' ' || r == '-':
			b.WriteRune('-')
		case r == '_':
			b.WriteRune('_')
		}
	}
	return b.String()
}

// UniqueAnchors は同名見出しに -1, -2 を付けて一意化する。
func UniqueAnchors(headings []string) []string {
	out := make([]string, 0, len(headings))
	seen := map[string]int{}
	for _, heading := range headings {
		base := GitHubAnchor(heading)
		n := seen[base]
		if n == 0 {
			out = append(out, base)
		} else {
			out = append(out, fmt.Sprintf("%s-%d", base, n))
		}
		seen[base] = n + 1
	}
	return out
}
