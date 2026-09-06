package repository_test

import (
	"testing"

	"github.com/takeuchi-shogo/md-linker/internal/repository"
)

func TestMemFS_ExistsAndRead_WhenFileWritten_ReturnsContents(t *testing.T) {
	fs := repository.NewMemFS()
	fs.WriteFile("docs/setup.md", []byte("# Setup"))

	if !fs.Exists("docs/setup.md") {
		t.Fatal("Exists(docs/setup.md) = false, want true")
	}

	got, err := fs.ReadFile("docs/setup.md")
	if err != nil {
		t.Fatalf("ReadFile(docs/setup.md) error = %v", err)
	}
	if string(got) != "# Setup" {
		t.Fatalf("ReadFile(docs/setup.md) = %q, want %q", got, "# Setup")
	}
}

func TestMemFS_IsFileAndIsDir_WhenNestedFileWritten_DistinguishesPathTypes(t *testing.T) {
	fs := repository.NewMemFS()
	fs.WriteFile("docs/setup.md", []byte("# Setup"))

	if !fs.IsFile("docs/setup.md") {
		t.Fatal("IsFile(docs/setup.md) = false, want true")
	}
	if !fs.IsDir("docs") {
		t.Fatal("IsDir(docs) = false, want true")
	}
	if fs.IsDir("docs/setup.md") {
		t.Fatal("IsDir(docs/setup.md) = true, want false")
	}
}

func TestMemFS_ReadDir_WhenDirectoryHasFiles_ReturnsEntryNames(t *testing.T) {
	fs := repository.NewMemFS()
	fs.WriteFile("docs/setup.md", []byte("# Setup"))
	fs.WriteFile("docs/guide.md", []byte("# Guide"))

	entries, err := fs.ReadDir("docs")
	if err != nil {
		t.Fatalf("ReadDir(docs) error = %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("ReadDir(docs) len = %d, want 2", len(entries))
	}

	names := map[string]bool{}
	for _, e := range entries {
		names[e.Name()] = true
	}
	if !names["setup.md"] || !names["guide.md"] {
		t.Fatalf("ReadDir(docs) names = %v, want setup.md and guide.md", names)
	}
}
