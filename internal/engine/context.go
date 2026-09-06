package engine

import "github.com/takeuchi-shogo/md-linker/internal/repository"

// Context は 1 回の check 実行で Rule が共有する検査コンテキスト。
type Context struct {
	Root      string
	FS        repository.FileSystem
	Documents map[string]*Document
	Anchors   map[string]map[string]struct{}
}

// NewContext は空の Document / Anchor 登録済み Context を返す。
func NewContext(root string, fs repository.FileSystem) *Context {
	return &Context{
		Root:      root,
		FS:        fs,
		Documents: map[string]*Document{},
		Anchors:   map[string]map[string]struct{}{},
	}
}
