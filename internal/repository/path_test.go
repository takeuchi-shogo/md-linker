package repository_test

import (
	"testing"

	"github.com/takeuchi-shogo/md-linker/internal/repository"
)

func TestResolve_WhenRelativeWithFragment_NormalizesAgainstDocDir(t *testing.T) {
	got, err := repository.Resolve("/repo", "docs/guide.md", "./setup.md#install")
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if got.RepoPath != "docs/setup.md" {
		t.Fatalf("RepoPath = %q, want %q", got.RepoPath, "docs/setup.md")
	}
	if got.Fragment != "install" {
		t.Fatalf("Fragment = %q, want %q", got.Fragment, "install")
	}
	if got.Outside {
		t.Fatal("Outside = true, want false")
	}
}

func TestResolve_WhenPathEscapesRoot_MarksOutside(t *testing.T) {
	got, err := repository.Resolve("/repo", "README.md", "../../etc/passwd")
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if !got.Outside {
		t.Fatal("Outside = false, want true")
	}
}

func TestResolve_WhenURLEncoded_DecodesPath(t *testing.T) {
	got, err := repository.Resolve("/repo", "README.md", "./docs/set%20up.md")
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if got.RepoPath != "docs/set up.md" {
		t.Fatalf("RepoPath = %q, want %q", got.RepoPath, "docs/set up.md")
	}
}
