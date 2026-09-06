# mdlinker

Repository-aware Markdown linter for Go. It treats Markdown as part of the repository and checks that local links, anchors, images, and high-confidence file references actually exist.

v0.1 is deterministic, local-only, and read-only. It does not call the network, run commands, or use an LLM.

## Install

```bash
go install github.com/takeuchi-shogo/md-linker/cmd/mdlinker@latest
```

From this repository:

```bash
go run ./cmd/mdlinker check .
```

## Usage

```bash
mdlinker                 # help
mdlinker --version
mdlinker check           # lint the current directory
mdlinker check .
mdlinker check README.md
mdlinker check docs/
mdlinker check . --format json
mdlinker check . --format github
mdlinker check . --config .mdlinker.yaml
```

## Output

text (default):

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

## Exit codes

| Code | Meaning |
|---|---|
| 0 | Success, or warnings only |
| 1 | One or more error diagnostics |
| 2 | Config / I/O / parser failure |

## Configuration

`.mdlinker.yaml` in the repository root is optional. Defaults:

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

Severity values: `error`, `warning`, `off`. CLI `--format` and `--config` override the file.

## Rules

| ID | Name | Default | What it checks |
|---|---|---|---|
| MDL001 | broken-link | error | Relative local link target is missing |
| MDL002 | broken-anchor | error | `#fragment` is missing in the target Markdown |
| MDL003 | missing-file-reference | warning | High-confidence path in inline code is missing |
| MDL004 | missing-image | error | Image target is missing |
| MDL006 | absolute-local-path | warning | Machine-local absolute path (home directory or Windows drive) |
| MDL007 | case-mismatch | error | Path exists but letter case does not match |

`MDL005` is unused on purpose.

## Out of scope (v0.1)

- External HTTP(S) link checks
- `--fix`
- Git `--changed` / `--affected` (planned for v0.2)
- LLM / semantic checks
- LSP / IDE plugin
