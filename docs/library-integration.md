# Library Integration: adf-to-markdown + goldmark-adf

This document describes how the two libraries work together to provide round-trip conversion between Atlassian Document Format (ADF) and Markdown.

## Overview

```
ADF JSON ──► adf-to-markdown ──► Markdown ──► goldmark-adf ──► ADF JSON
```

- **adf-to-markdown**: Converts ADF JSON documents to GitHub Flavored Markdown, using custom syntax extensions for ADF-specific nodes.
- **goldmark-adf**: A goldmark renderer (with parser extensions) that converts Markdown back to ADF JSON.

Together, these libraries enable lossless round-tripping for most ADF content.

## Extension Summary

Both libraries implement coordinated syntax extensions to represent ADF nodes that have no native markdown equivalent. The authoritative specification is in [roundtrip-extensions.md](roundtrip-extensions.md).

### What adf-to-markdown produces

| ADF Node | Markdown Output | Extension Type |
|---|---|---|
| `paragraph`, `heading`, `blockquote`, `codeBlock`, `bulletList`, `orderedList`, `listItem`, `rule`, `hardBreak` | Standard GFM | Native |
| `text` + `strong`/`em`/`strike`/`code`/`link` marks | Standard GFM inline | Native |
| `table`, `tableRow`, `tableHeader`, `tableCell` | GFM tables | Native |
| `taskList`, `taskItem`, `blockTaskItem` | `- [x]` / `- [ ]` | GFM |
| `mediaSingle`, `media` (external) | `![alt](url)` | GFM image |
| `panel` | `> [!TYPE]` (GitHub alerts) | GitHub convention |
| `expand`, `nestedExpand` | `<details><summary>` | GitHub convention |
| `emoji` | `:shortcode:` | GitHub convention |
| `status` | `[status:text\|color]` | Custom extension |
| `mention` | `@[name](id)` | Custom extension |
| `date` | `[date:timestamp]` | Custom extension |
| `placeholder` | `{{text}}` | Custom extension |
| `inlineCard`, `blockCard` | `[card:url]` | Custom extension |
| `embedCard` | `[embed:url]` | Custom extension |
| `decisionList`, `decisionItem` | `- [!]` / `- [?]` | Custom extension |
| `layoutSection`, `layoutColumn` | `[layout-section]` / `[layout-column N]` | Custom extension |
| `extension`, `bodiedExtension`, etc. | `[extension:type:key]` | Custom extension |
| `text` + `underline` mark | `<u>text</u>` | HTML inline |
| `text` + `subsup` mark | `<sub>text</sub>` / `<sup>text</sup>` | HTML inline |
| `text` + `textColor` mark | `<span style="color:#hex">text</span>` | HTML inline |

### What goldmark-adf must parse

goldmark-adf needs parser extensions to recognize the custom syntax and emit the corresponding ADF nodes. The extensions are organized by type:

**Block-level extensions** (goldmark parser extensions):
- GitHub alert syntax in blockquotes → `panel`
- `<details><summary>` HTML blocks → `expand`
- `- [!]` / `- [?]` list items → `decisionList` / `decisionItem`
- `[layout-section]` / `[layout-column N]` markers → `layoutSection` / `layoutColumn`
- `[extension:type:key]` markers → extension nodes

**Inline extensions** (goldmark inline parsers):
- `[status:text|color]` → `status`
- `@[name](id)` → `mention`
- `[date:timestamp]` → `date`
- `{{text}}` → `placeholder`
- `[card:url]` → `inlineCard` or `blockCard`
- `[embed:url]` → `embedCard`
- `:shortcode:` → `emoji`

**HTML inline parsing** (goldmark inline parsers):
- `<u>text</u>` → `underline` mark
- `<sub>text</sub>` → `subsup` mark (type: sub)
- `<sup>text</sup>` → `subsup` mark (type: sup)

**GFM behavior changes**:
- Task lists (`- [x]`/`- [ ]`) must emit `taskList`/`taskItem` ADF nodes instead of text prefixes

## Escaping Convention

Custom inline tokens use backslash escaping for delimiter characters within user-controlled text fields. Both libraries must implement matching escape/unescape logic:

- **adf-to-markdown**: escapes delimiters when writing tokens
- **goldmark-adf**: unescapes delimiters when parsing tokens

See the escaping section in [roundtrip-extensions.md](roundtrip-extensions.md) for the full specification.

## Known Losses

Some ADF attributes are intentionally not preserved in the markdown representation:

| Attribute | Node | Reason |
|---|---|---|
| `localId` | status, mention, decision items, task items | Session-specific identifiers, not meaningful across conversions |
| `accessLevel`, `userType` | mention | Authorization metadata, not content |
| `style` | status | Rarely used, color is sufficient |
| `textColor`, `backgroundColor` marks | text | Preserved as HTML spans but may not round-trip through all parsers |
| `layout` attributes | mediaSingle | Image layout hints (center, wide, etc.) |
| `isNumberColumnEnabled`, `layout` | table | Table display configuration |
| `colspan`, `rowspan`, `colwidth`, `background` | table cells | Complex table formatting |
| `media` (file/link type) attributes | media | Collection IDs preserved in custom URL scheme but may not resolve |

## Configuration

### adf-to-markdown

No special configuration is needed — the new syntax is the default output.

Custom extension handlers (`ExtensionBlock`, `ExtensionInline`) take precedence over the default `[extension:type:key]` syntax when provided.

### goldmark-adf

The following options are relevant for round-trip usage:

```go
md := adf.NewWithGFM(
    adf.WithExternalMedia(true), // Required for ![alt](url) → mediaSingle
)
```

Parser extensions for the custom syntax will need to be registered when available. See the goldmark-adf documentation for details.

## Implementation Status

### adf-to-markdown (this library)

All syntax changes are implemented:
- Status: `[status:text|color]` with backslash escaping
- Mention: `@[name](id)` with backslash escaping
- Emoji: normalized to `:shortName:`
- Cards: `[card:url]` and `[embed:url]`
- Decisions: `- [!]` / `- [?]`
- Escape helpers: `escapeDelimiters()` for safe token rendering

### goldmark-adf

The following extensions are planned (see [roundtrip-extensions.md](roundtrip-extensions.md) for the full implementation plan):

- [ ] Task list fix (emit `taskList`/`taskItem` ADF nodes)
- [ ] Panel alert parser extension
- [ ] `<details>` expand parser extension
- [ ] `:shortcode:` emoji parser extension
- [ ] Inline token parsers (status, mention, date, placeholder, cards)
- [ ] Decision list parser extension
- [ ] Layout section parser
- [ ] Extension marker parser
- [ ] Inline HTML mark parsing (`<u>`, `<sub>`, `<sup>`)
