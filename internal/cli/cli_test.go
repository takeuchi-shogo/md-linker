package cli_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/takeuchi-shogo/md-linker/internal/cli"
)

func TestRun_WhenNoArgs_ShowsHelpAndExit0(t *testing.T) {
	var out bytes.Buffer
	code := cli.Run([]string{}, &out, io.Discard)
	if code != 0 {
		t.Fatalf("Run() code = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "mdlinker") {
		t.Fatalf("help = %q, want it to contain mdlinker", out.String())
	}
}

func TestRun_WhenUnknownCommand_ReturnsExit2(t *testing.T) {
	code := cli.Run([]string{"lint"}, io.Discard, io.Discard)
	if code != 2 {
		t.Fatalf("Run() code = %d, want 2", code)
	}
}

func TestRun_WhenVersion_PrintsVersionAndExit0(t *testing.T) {
	var out bytes.Buffer
	code := cli.Run([]string{"--version"}, &out, io.Discard)
	if code != 0 {
		t.Fatalf("Run() code = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "0.1.0") {
		t.Fatalf("version = %q, want 0.1.0", out.String())
	}
}
