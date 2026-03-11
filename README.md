# adf-to-markdown

Go library for converting Atlassian Document Format (ADF) JSON into Markdown.

## Contributing

See `CONTRIBUTING.md` for local development, testing, versioning, and release workflow details.

## Requirements

- Go `1.25+`
- `GOEXPERIMENT=jsonv2` enabled (this library uses `encoding/json/v2`)

## Install

```bash
go get github.com/ajbeck/adf-to-markdown
```

## Generate (Before Build/Test)

If your project uses code generation, run:

```bash
GOEXPERIMENT=jsonv2 go generate ./...
```

Then run tests:

```bash
GOEXPERIMENT=jsonv2 go test ./...
```

Note: this repository is safe to run with `go generate ./...` even when no generators are configured.

## Basic Usage

```go
package main

import (
	"fmt"
	"log"

	adfmarkdown "github.com/ajbeck/adf-to-markdown"
)

func main() {
	input := []byte(`{
		"version": 1,
		"type": "doc",
		"content": [
			{
				"type": "heading",
				"attrs": {"level": 2},
				"content": [{"type": "text", "text": "Overview"}]
			},
			{
				"type": "paragraph",
				"content": [{"type": "text", "text": "Hello from ADF"}]
			}
		]
	}`)

	md, err := adfmarkdown.UnmarshalADF(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(md))
}
```

Run with:

```bash
GOEXPERIMENT=jsonv2 go run .
```

## Writing to an `io.Writer`

```go
var buf bytes.Buffer
err := adfmarkdown.UnmarshalADFTo(&buf, adfJSON)
```

## Using a Typed Package in a Consuming App

In your app, keep API payloads typed (for example in `internal/jiratypes`) and pass the typed
ADF field to this library.

```go
package jiratypes

type Document struct {
	Version int    `json:"version"`
	Type    string `json:"type"`
	Content []Node `json:"content"`
}

type Node struct {
	Type    string         `json:"type"`
	Attrs   map[string]any `json:"attrs,omitempty"`
	Content []Node         `json:"content,omitempty"`
	Text    string         `json:"text,omitempty"`
}

type Issue struct {
	Fields struct {
		Description Document `json:"description"`
	} `json:"fields"`
}
```

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	adfmarkdown "github.com/ajbeck/adf-to-markdown"
	"yourmodule/internal/jiratypes"
)

func renderIssueDescription(issueJSON []byte) (string, error) {
	var issue jiratypes.Issue
	if err := json.Unmarshal(issueJSON, &issue); err != nil {
		return "", err
	}

	adfBytes, err := json.Marshal(issue.Fields.Description)
	if err != nil {
		return "", err
	}

	md, err := adfmarkdown.UnmarshalADF(adfBytes)
	if err != nil {
		return "", err
	}
	return string(md), nil
}

func main() {
	md, err := renderIssueDescription([]byte(`{"fields":{"description":{"version":1,"type":"doc","content":[]}}}`))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(md)
}
```

## Round-Trip with goldmark-adf

This library is designed to work with [goldmark-adf](https://github.com/ajbeck/goldmark-adf) for lossless ADF round-tripping:

```
ADF JSON --> adf-to-markdown --> Markdown --> goldmark-adf --> ADF JSON
```

**adf-to-markdown** converts ADF JSON to Markdown, using custom syntax extensions for ADF-specific nodes that have no native Markdown equivalent. **goldmark-adf** parses that Markdown — including the custom extensions — back into ADF JSON.

### Custom Extension Syntax

| ADF Node | Markdown Output |
|---|---|
| `status` | `[status:text\|color]` |
| `mention` | `@[name](id)` |
| `date` | `[date:timestamp]` |
| `placeholder` | `{{text}}` |
| `inlineCard` / `blockCard` | `[card:url]` |
| `embedCard` | `[embed:url]` |
| `emoji` | `:shortcode:` |
| `panel` | `> [!NOTE]` (GitHub alert syntax) |
| `decisionList` / `decisionItem` | `- [!] text` / `- [?] text` |
| `taskList` / `taskItem` | `- [x]` / `- [ ]` |
| `expand` | `<details><summary>` |

Delimiter characters inside user text are backslash-escaped to prevent ambiguity. See [docs/roundtrip-extensions.md](docs/roundtrip-extensions.md) for the full specification and [docs/library-integration.md](docs/library-integration.md) for integration details.

## Useful Options

- `WithStrictSchema(bool)`
- `WithBuiltInSchemaValidation(bool)`
- `WithAllowUnsupportedNodes(bool)`
- `WithHardBreakStyle(...)`
- `WithCodeFenceStyle(...)`
- `WithSchemaValidator(func([]byte) error)`
- `WithUnsupportedBlockHandler(...)`
- `WithUnsupportedInlineHandler(...)`
- `WithExtensionBlockHandler(...)`
- `WithExtensionInlineHandler(...)`

## Error Handling

The library returns typed errors:

- `*adfmarkdown.Error` with fields:
- `Path` (ADF path, e.g. `/content/0/content/1`)
- `Kind` (`adfmarkdown.ErrorKind`)
- `Detail`
- `Cause` (wrapped error; available via `errors.Unwrap` / `errors.As`)

Common `ErrorKind` values include:

- `ErrKindInvalidRoot`
- `ErrKindInvalidJSON`
- `ErrKindMissingType`
- `ErrKindUnsupportedNode`
- `ErrKindUnsupportedInline`
- `ErrKindUnsupportedMark`
- `ErrKindInvalidAttr`
- `ErrKindInvalidMark`
- `ErrKindInvalidMarkCombo`
- `ErrKindInvalidStructure`
- `ErrKindInvalidText`

## v1 Compatibility Contract

- Build/runtime requirements: Go `1.25+` and `GOEXPERIMENT=jsonv2`.
- Stable entry points (v1): `UnmarshalADF`, `UnmarshalADFTo`.
- Stable options (v1): `WithStrictSchema`, `WithBuiltInSchemaValidation`, `WithAllowUnsupportedNodes`, `WithHardBreakStyle`, `WithCodeFenceStyle`, `WithSchemaValidator`, `WithUnsupportedBlockHandler`, `WithUnsupportedInlineHandler`, `WithExtensionBlockHandler`, `WithExtensionInlineHandler`.
- Stable error surface (v1): returned typed error `*adfmarkdown.Error` and the `ErrorKind` constants listed above.

## Tests

```bash
GOEXPERIMENT=jsonv2 go test ./...
```

Or via `make`:

```bash
make test
make test-nojsonv2
```

## Fuzzing

```bash
GOEXPERIMENT=jsonv2 go test -fuzz=FuzzUnmarshalADF -run=^$ ./...
```

Or via `make`:

```bash
make fuzz
```

## Benchmarks

```bash
GOEXPERIMENT=jsonv2 go test -run=^$ -bench=BenchmarkUnmarshalADF -benchmem ./...
```

Or via `make`:

```bash
make bench
```
