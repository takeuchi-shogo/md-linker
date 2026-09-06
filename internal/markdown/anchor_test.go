package markdown_test

import (
	"reflect"
	"testing"

	"github.com/takeuchi-shogo/md-linker/internal/markdown"
)

func TestGitHubAnchor_WhenEnglishHeading_LowercasesAndHyphenates(t *testing.T) {
	cases := map[string]string{
		"Install": "install",
		"Set up":  "set-up",
		"インストール":  "インストール",
	}
	for in, want := range cases {
		if got := markdown.GitHubAnchor(in); got != want {
			t.Fatalf("GitHubAnchor(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestUniqueAnchors_WhenDuplicateHeadings_AppendsNumericSuffix(t *testing.T) {
	got := markdown.UniqueAnchors([]string{"Install", "Install"})
	want := []string{"install", "install-1"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("UniqueAnchors() = %v, want %v", got, want)
	}
}
