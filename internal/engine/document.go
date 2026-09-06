package engine

// Document は Markdown 1 ファイルから Rule が参照する共通モデル。
// goldmark など parser 固有 AST はここへ変換し、Rule には漏らさない。
type Document struct {
	Path     string
	Source   []byte
	Links    []Link
	Images   []Image
	Headings []Heading
	Code     []CodeSpan
}

// Link は Markdown リンク 1 件。Dest は fragment を除いた参照先。
type Link struct {
	Raw      string
	Dest     string
	Fragment string
	Line     int
	Column   int
}

// Image は Markdown 画像参照 1 件。
type Image struct {
	Raw      string
	Dest     string
	Fragment string
	Line     int
	Column   int
}

// Heading は見出し 1 件。Anchor は GitHub 互換の一意な fragment。
type Heading struct {
	Text   string
	Level  int
	Anchor string
	Line   int
	Column int
}

// CodeSpan は inline code 1 件。fenced code block は v0.1 では含めない。
type CodeSpan struct {
	Text   string
	Line   int
	Column int
}
