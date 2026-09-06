package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/takeuchi-shogo/md-linker/internal/app"
	"github.com/takeuchi-shogo/md-linker/internal/config"
	"github.com/takeuchi-shogo/md-linker/internal/repository"
)

// Version は --version で表示する SemVer。
const Version = "0.1.0"

const usage = `mdlinker — Repository-aware Markdown linter

Usage:
  mdlinker                 Show this help
  mdlinker --help          Show this help
  mdlinker --version       Show version
  mdlinker check [path]    Lint Markdown under path (default: .)

Options for check:
  --format text|json|github   Output format (default: text)
  --config <file>             Config file (default: .mdlinker.yaml)

Examples:
  mdlinker check .
  mdlinker check README.md
  mdlinker check docs/
  mdlinker check . --format json
`

// Run は CLI 引数を解釈して終了コードを返す。args は os.Args[1:] 相当。
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || isHelp(args) {
		fmt.Fprint(stdout, usage)
		return 0
	}
	if isVersion(args) {
		fmt.Fprintln(stdout, Version)
		return 0
	}
	if args[0] != "check" {
		fmt.Fprintf(stderr, "mdlinker: unknown command %q\n%s", args[0], usage)
		return 2
	}
	return runCheck(args[1:], stdout, stderr)
}

func runCheck(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	fs.SetOutput(stderr)
	format := fs.String("format", "text", "output format: text, json, github")
	configPath := fs.String("config", "", "path to .mdlinker.yaml")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	target := "."
	if fs.NArg() > 0 {
		target = fs.Arg(0)
	}

	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(stderr, "mdlinker: getwd: %v\n", err)
		return 2
	}

	cfg, err := loadConfig(root, *configPath)
	if err != nil {
		fmt.Fprintf(stderr, "mdlinker: %v\n", err)
		return 2
	}

	code, err := app.Run(repository.NewOSFS(root), root, target, cfg, *format, stdout)
	if err != nil {
		fmt.Fprintf(stderr, "mdlinker: %v\n", err)
		return 2
	}
	return code
}

func loadConfig(root, explicit string) (config.Config, error) {
	if explicit != "" {
		return config.Load(explicit)
	}
	return config.LoadOrDefault(root)
}

func isHelp(args []string) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}

func isVersion(args []string) bool {
	for _, arg := range args {
		if arg == "--version" || arg == "-v" {
			return true
		}
	}
	return false
}
