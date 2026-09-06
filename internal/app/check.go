package app

import (
	"fmt"
	"io"

	"github.com/takeuchi-shogo/md-linker/internal/config"
	"github.com/takeuchi-shogo/md-linker/internal/engine"
	"github.com/takeuchi-shogo/md-linker/internal/markdown"
	"github.com/takeuchi-shogo/md-linker/internal/reporter"
	"github.com/takeuchi-shogo/md-linker/internal/repository"
	"github.com/takeuchi-shogo/md-linker/internal/rules"
	"github.com/takeuchi-shogo/md-linker/internal/scanner"
)

// Run は対象 Markdown を検査し、診断を書いて終了コードを返す。
func Run(fs repository.FileSystem, root, target string, cfg config.Config, format string, w io.Writer) (int, error) {
	files, err := scanner.Scan(fs, root, target, cfg)
	if err != nil {
		return 2, fmt.Errorf("app.Run: scan root=%q target=%q: %w", root, target, err)
	}

	docs := make([]*engine.Document, 0, len(files))
	for _, filePath := range files {
		data, readErr := fs.ReadFile(filePath)
		if readErr != nil {
			return 2, fmt.Errorf("app.Run: ReadFile path=%q: %w", filePath, readErr)
		}
		doc, parseErr := markdown.Parse(filePath, data)
		if parseErr != nil {
			return 2, fmt.Errorf("app.Run: Parse path=%q: %w", filePath, parseErr)
		}
		docs = append(docs, doc)
	}

	ctx := engine.NewContext(root, fs)
	ctx.Register(docs)
	diags := engine.New(rules.All()).Run(ctx, docs, cfg.Rules)
	if writeErr := reporter.Write(w, format, diags); writeErr != nil {
		return 2, fmt.Errorf("app.Run: Write format=%q: %w", format, writeErr)
	}
	return reporter.ExitCode(diags, nil), nil
}
