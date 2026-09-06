package repository

import (
	"fmt"
	"os"
	"path/filepath"
)

// OSFS は Repository root を基準にした読み取り専用 OS アダプタ。
type OSFS struct {
	root string
}

// NewOSFS は root 配下を見る FileSystem を返す。
func NewOSFS(root string) *OSFS {
	return &OSFS{root: root}
}

func (o *OSFS) resolve(p string) string {
	p = filepath.FromSlash(p)
	if filepath.IsAbs(p) {
		return p
	}
	if o.root == "" || o.root == "." {
		return p
	}
	return filepath.Join(o.root, p)
}

func (o *OSFS) Exists(p string) bool {
	_, err := os.Stat(o.resolve(p))
	return err == nil
}

func (o *OSFS) IsFile(p string) bool {
	info, err := os.Stat(o.resolve(p))
	return err == nil && !info.IsDir()
}

func (o *OSFS) IsDir(p string) bool {
	info, err := os.Stat(o.resolve(p))
	return err == nil && info.IsDir()
}

func (o *OSFS) ReadFile(p string) ([]byte, error) {
	full := o.resolve(p)
	data, err := os.ReadFile(full)
	if err != nil {
		return nil, fmt.Errorf("repository.OSFS.ReadFile: path=%q: %w", full, err)
	}
	return data, nil
}

func (o *OSFS) ReadDir(p string) ([]DirEntry, error) {
	full := o.resolve(p)
	entries, err := os.ReadDir(full)
	if err != nil {
		return nil, fmt.Errorf("repository.OSFS.ReadDir: path=%q: %w", full, err)
	}
	out := make([]DirEntry, 0, len(entries))
	for _, entry := range entries {
		out = append(out, osDirEntry{name: entry.Name(), isDir: entry.IsDir()})
	}
	return out, nil
}

type osDirEntry struct {
	name  string
	isDir bool
}

func (e osDirEntry) Name() string { return e.name }
func (e osDirEntry) IsDir() bool  { return e.isDir }
