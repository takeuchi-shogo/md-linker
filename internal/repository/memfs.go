package repository

import (
	"fmt"
	"path"
	"strings"
)

// memEntry は MemFS 上のファイルまたはディレクトリ。
type memEntry struct {
	name  string
	isDir bool
	data  []byte
}

func (e memEntry) Name() string {
	return e.name
}

func (e memEntry) IsDir() bool {
	return e.isDir
}

// MemFS はテスト用のメモリ FileSystem。パスは / 区切りで正規化する。
type MemFS struct {
	files map[string][]byte
	dirs  map[string]struct{}
}

// NewMemFS は空のメモリファイルシステムを返す。
func NewMemFS() *MemFS {
	return &MemFS{
		files: map[string][]byte{},
		dirs:  map[string]struct{}{".": {}},
	}
}

// WriteFile はファイルを書き込み、親ディレクトリも作成する。
func (m *MemFS) WriteFile(p string, data []byte) {
	p = normalizeRepoPath(p)
	m.files[p] = append([]byte(nil), data...)
	m.addDirAncestors(p)
}

func (m *MemFS) addDirAncestors(filePath string) {
	dir := path.Dir(filePath)
	for dir != "." && dir != "/" && dir != "" {
		m.dirs[dir] = struct{}{}
		dir = path.Dir(dir)
	}
}

// Exists はファイルまたはディレクトリが存在するとき true。
// 実名は保持したまま、照合だけ大小文字を無視する（MDL007 用）。
func (m *MemFS) Exists(p string) bool {
	if _, ok := m.lookupFile(p); ok {
		return true
	}
	_, ok := m.lookupDir(p)
	return ok
}

// IsFile は通常ファイルが存在するとき true。大小文字は無視する。
func (m *MemFS) IsFile(p string) bool {
	_, ok := m.lookupFile(p)
	return ok
}

// IsDir はディレクトリが存在するとき true。大小文字は無視する。
func (m *MemFS) IsDir(p string) bool {
	_, ok := m.lookupDir(p)
	return ok
}

// ReadFile はファイル内容を返す。存在しない場合はエラー。
func (m *MemFS) ReadFile(p string) ([]byte, error) {
	data, ok := m.lookupFile(p)
	if !ok {
		return nil, fmt.Errorf("repository.MemFS.ReadFile: path=%q: file not found", p)
	}
	return append([]byte(nil), data...), nil
}

// ReadDir はディレクトリ直下のエントリを返す。引数の大小文字は無視し、実名を返す。
func (m *MemFS) ReadDir(p string) ([]DirEntry, error) {
	realDir, ok := m.lookupDir(p)
	if !ok {
		return nil, fmt.Errorf("repository.MemFS.ReadDir: path=%q: directory not found", p)
	}
	p = realDir

	seen := map[string]DirEntry{}
	prefix := p + "/"
	if p == "." {
		prefix = ""
	}

	for filePath := range m.files {
		name, ok := directChild(prefix, filePath)
		if !ok {
			continue
		}
		seen[name] = memEntry{name: name, isDir: false}
	}
	for dirPath := range m.dirs {
		if dirPath == "." || dirPath == p {
			continue
		}
		name, ok := directChild(prefix, dirPath)
		if !ok {
			continue
		}
		seen[name] = memEntry{name: name, isDir: true}
	}

	entries := make([]DirEntry, 0, len(seen))
	for _, e := range seen {
		entries = append(entries, e)
	}
	return entries, nil
}

func directChild(prefix, full string) (string, bool) {
	if prefix != "" && !strings.HasPrefix(full, prefix) {
		return "", false
	}
	rest := full
	if prefix != "" {
		rest = strings.TrimPrefix(full, prefix)
	}
	if rest == "" || strings.Contains(rest, "/") {
		return "", false
	}
	return rest, true
}

func (m *MemFS) lookupFile(p string) ([]byte, bool) {
	p = normalizeRepoPath(p)
	if data, ok := m.files[p]; ok {
		return data, true
	}
	for real, data := range m.files {
		if strings.EqualFold(real, p) {
			return data, true
		}
	}
	return nil, false
}

func (m *MemFS) lookupDir(p string) (string, bool) {
	p = normalizeRepoPath(p)
	if _, ok := m.dirs[p]; ok {
		return p, true
	}
	for real := range m.dirs {
		if strings.EqualFold(real, p) {
			return real, true
		}
	}
	return "", false
}

func normalizeRepoPath(p string) string {
	p = strings.ReplaceAll(p, "\\", "/")
	p = path.Clean(p)
	p = strings.TrimPrefix(p, "/")
	if p == "" {
		return "."
	}
	return p
}
