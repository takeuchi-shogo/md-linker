package main

import (
	"os"

	"github.com/takeuchi-shogo/md-linker/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
