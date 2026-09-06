package config_test

import (
	"testing"

	"github.com/takeuchi-shogo/md-linker/internal/config"
	"github.com/takeuchi-shogo/md-linker/internal/engine"
)

func TestDefault_WhenCalled_SetsWarningForMDL003(t *testing.T) {
	got := config.Default()

	if got.Rules["MDL003"] != engine.SeverityWarning {
		t.Fatalf("Default().Rules[MDL003] = %q, want %q", got.Rules["MDL003"], engine.SeverityWarning)
	}
}

func TestLoadBytes_WhenUnknownRule_ReturnsError(t *testing.T) {
	_, err := config.LoadBytes([]byte("rules:\n  MDL999: error\n"))
	if err == nil {
		t.Fatal("LoadBytes() error = nil, want unknown rule error")
	}
}

func TestLoadBytes_WhenInvalidSeverity_ReturnsError(t *testing.T) {
	_, err := config.LoadBytes([]byte("rules:\n  MDL001: fatal\n"))
	if err == nil {
		t.Fatal("LoadBytes() error = nil, want invalid severity error")
	}
}

func TestLoadBytes_WhenValidOverride_ReplacesDefaultSeverity(t *testing.T) {
	got, err := config.LoadBytes([]byte("rules:\n  MDL003: error\n"))
	if err != nil {
		t.Fatalf("LoadBytes() error = %v", err)
	}
	if got.Rules["MDL003"] != engine.SeverityError {
		t.Fatalf("Rules[MDL003] = %q, want %q", got.Rules["MDL003"], engine.SeverityError)
	}
}
