package scanner_test

import (
	"reflect"
	"testing"

	"github.com/takeuchi-shogo/md-linker/internal/config"
	"github.com/takeuchi-shogo/md-linker/internal/repository"
	"github.com/takeuchi-shogo/md-linker/internal/scanner"
)

func TestScan_WhenVendorPresent_ExcludesVendorMarkdown(t *testing.T) {
	fs := repository.NewMemFS()
	fs.WriteFile("README.md", []byte("x"))
	fs.WriteFile("vendor/x.md", []byte("x"))
	fs.WriteFile("docs/a.mdx", []byte("x"))

	got, err := scanner.Scan(fs, ".", ".", config.Default())
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	want := []string{"README.md", "docs/a.mdx"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Scan() = %v, want %v", got, want)
	}
}

func TestScan_WhenTargetIsSingleMarkdown_ReturnsThatFile(t *testing.T) {
	fs := repository.NewMemFS()
	fs.WriteFile("README.md", []byte("x"))
	fs.WriteFile("docs/a.md", []byte("x"))

	got, err := scanner.Scan(fs, ".", "README.md", config.Default())
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	want := []string{"README.md"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Scan() = %v, want %v", got, want)
	}
}

func TestScan_WhenTargetIsNotMarkdown_ReturnsError(t *testing.T) {
	fs := repository.NewMemFS()
	fs.WriteFile("main.go", []byte("package main"))

	_, err := scanner.Scan(fs, ".", "main.go", config.Default())
	if err == nil {
		t.Fatal("Scan() error = nil, want error for non-markdown file")
	}
}
