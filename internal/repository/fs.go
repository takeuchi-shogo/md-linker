package repository

// DirEntry はディレクトリ内の 1 エントリ。MDL007 の実名比較に使う。
type DirEntry interface {
	Name() string
	IsDir() bool
}

// FileSystem は Rule / Scanner が使う読み取り専用のファイル抽象。
// 実装は OS 実ファイルまたはテスト用 MemFS。
type FileSystem interface {
	Exists(path string) bool
	IsFile(path string) bool
	IsDir(path string) bool
	ReadFile(path string) ([]byte, error)
	ReadDir(path string) ([]DirEntry, error)
}
