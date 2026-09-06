package engine

// Rule は 1 つの Lint 規則。実装は FileSystem を Context 経由でのみ使う。
type Rule interface {
	ID() string
	Check(ctx *Context, doc *Document) []Diagnostic
}
