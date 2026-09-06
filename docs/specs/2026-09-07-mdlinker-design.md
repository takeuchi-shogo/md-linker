# mdlinker v0.1 Design

## Goal

Markdown を Repository の構成要素として静的検査する Go CLI `mdlinker` を作る。
リンク切れ、存在しないアンカー / 画像 / ファイル参照、端末固有の絶対パス、path の case mismatch を決定論的に検出する。

出典: `mdguard_requirements_v0.1_updated.docx`（製品名だけ mdlinker に置換）

## Naming

| 項目 | 値 |
|---|---|
| リポジトリ | `https://github.com/takeuchi-shogo/md-linker` |
| module path | `github.com/takeuchi-shogo/md-linker` |
| バイナリ / コマンド | `mdlinker` |
| 設定ファイル | `.mdlinker.yaml` |
| Rule ID | `MDL001` … `MDL007`（`MDL005` は欠番。要件書 MDG と番号対応） |

## CLI

```
mdlinker                 → ヘルプ。exit 0
mdlinker --help          → ヘルプ。exit 0
mdlinker --version       → バージョン表示。exit 0
mdlinker check           → カレントディレクトリを検査
mdlinker check .         → 同上
mdlinker check README.md
mdlinker check docs/
mdlinker check . --format json
mdlinker check . --format github
mdlinker check . --config path/to/.mdlinker.yaml
```

- `--format`: `text` / `json` / `github`。既定は `text`
- `--config`: 省略時は対象 Repository root の `.mdlinker.yaml`。無くても既定設定で動く
- CLI オプションが設定ファイルを上書きする

## Requirements

- `.md` / `.mdx` を探索する
- 既定除外: `.git/**`, `node_modules/**`, `vendor/**`
- 読み取り専用。ファイルを変更しない
- ネットワークアクセスなし
- AI / LLM なし
- `--fix` なし
- 同一入力・同一設定なら同一診断

## Rules

| ID | 名称 | 既定 Severity | 内容 |
|---|---|---|---|
| MDL001 | broken-link | error | 相対ローカルリンク先が存在しない |
| MDL002 | broken-anchor | error | `#fragment` がリンク先 Markdown に無い |
| MDL003 | missing-file-reference | warning | 高確度な Repository パス参照が存在しない |
| MDL004 | missing-image | error | 画像参照先が存在しない |
| MDL006 | absolute-local-path | warning | `/Users/...` や `C:\Users\...` 等 |
| MDL007 | case-mismatch | error | 参照パスと実ファイル名の大小文字が違う |

MDL001:

- `http://` / `https://` / `mailto:` は対象外
- `#anchor` のみは MDL002
- 相対パスは参照元 Markdown のディレクトリ基準
- URL デコードする
- Repository root 外へ出るパスは v0.1 では検査しない

MDL002:

- GitHub 互換アンカー
- 同名見出しは `-1`, `-2` suffix
- 同一ファイル / 別ファイルの両方

MDL003:

- inline code を主対象。fenced code block は v0.1 対象外
- `./` `../` または拡張子付きパスのみ

## Exit codes

| code | 意味 |
|---|---|
| 0 | 正常完了。warning のみでも 0 |
| 1 | error Severity が 1 件以上 |
| 2 | 設定不正、読込失敗、parser 失敗などの実行エラー |

## Architecture

```
CLI (cmd/mdlinker)
  → app.Check
    → config + scanner + parser + engine + reporter
```

Rule は `os` を直接呼ばず `FileSystem` を使う。
Rule は goldmark AST を見ず `Document` だけを見る。
Path 解決と GitHub 互換アンカー生成は共通処理に集約する。

```
cmd/mdlinker/main.go
internal/app/check.go
internal/config/
internal/engine/          # Document, Diagnostic, Rule, Engine, Context
internal/markdown/        # parser, anchor
internal/repository/      # FileSystem, PathResolver
internal/rules/           # MDL001-007
internal/scanner/
internal/reporter/        # text, json, github
testdata/fixtures/
```

## Core types

```go
type Severity string

const (
    SeverityError   Severity = "error"
    SeverityWarning Severity = "warning"
    SeverityOff     Severity = "off"
)

type Document struct {
    Path     string
    Source   []byte
    Links    []Link
    Images   []Image
    Headings []Heading
    Code     []CodeSpan
}

type Diagnostic struct {
    RuleID     string
    Severity   Severity
    Path       string
    Line       int
    Column     int
    Message    string
    Suggestion string
}

type Rule interface {
    ID() string
    Check(ctx *Context, doc *Document) []Diagnostic
}

type FileSystem interface {
    Exists(path string) bool
    IsFile(path string) bool
    IsDir(path string) bool
    ReadFile(path string) ([]byte, error)
    ReadDir(path string) ([]DirEntry, error)
}
```

`ReadDir` は要件書の最小 interface への追加。Scanner と MDL007 のディレクトリエントリ比較に必要。

## Config

設定が無い場合の既定:

```yaml
include:
  - "**/*.md"
  - "**/*.mdx"
exclude:
  - ".git/**"
  - "node_modules/**"
  - "vendor/**"
rules:
  MDL001: error
  MDL002: error
  MDL003: warning
  MDL004: error
  MDL006: warning
  MDL007: error
```

未知の Rule ID または不正な Severity は設定エラー（exit 2）。

## Output

text:

```
README.md:42:5 error MDL001: "./docs/setup.md" does not exist
```

json:

```json
{
  "diagnostics": [
    {
      "rule": "MDL001",
      "severity": "error",
      "path": "README.md",
      "line": 42,
      "column": 5,
      "message": "./docs/setup.md does not exist"
    }
  ]
}
```

github: GitHub Actions annotation 形式。

診断は path, line, column, rule, severity で安定ソートする。

## Tech choices

- Go 1.22+
- Markdown parser: `github.com/yuin/goldmark`
- YAML: `gopkg.in/yaml.v3`
- CLI: 標準 `flag` + サブコマンド分岐。Cobra は使わない
- 公開 library API は切らない。Engine は `internal/` に置く

## Acceptance Criteria

AC-001 存在しない相対リンク → MDL001  
AC-002 存在しないアンカー → MDL002  
AC-003 存在しない画像 → MDL004  
AC-004 高確度ファイル参照の欠落 → MDL003  
AC-005 ローカル絶対パス → MDL006  
AC-006 case mismatch → MDL007  
AC-007 text / json 出力  
AC-008 `.mdlinker.yaml` で severity / include / exclude を変更できる  
AC-009 exit 0 / 1 / 2  
AC-010 検査前後で Repository ファイルが変わらない  
AC-011 ネットワークなしで全機能が動く  
AC-012 README、設定例、主要 Rule 説明がある  

## Out of Scope (v0.1)

- 外部 HTTP 疎通
- LLM / 意味的整合
- 文章校正
- コマンド実行
- `--fix`
- `--changed` / `--target` / `--affected`（v0.2）
- LSP / IDE
- 参照グラフ永続化
- 公開 Go API

## Constraints

- 読み取り専用
- ネットワーク不要
- panic をユーザーへ出さない
- Rule ID の意味を再利用しない
- OS の case sensitivity だけに頼らず、ディレクトリエントリ比較で MDL007 を判定する
