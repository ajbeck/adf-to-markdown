# Markdown Extensions

adf-to-markdown uses custom markdown syntax to represent ADF nodes that have no native markdown equivalent. This page documents the syntax produced by each extension and how to work with it.

For the full round-trip specification shared with goldmark-adf, see [roundtrip-extensions.md](roundtrip-extensions.md).

## Standard Markdown

These ADF nodes map directly to CommonMark or GFM syntax and require no special handling:

| ADF Node | Markdown |
|---|---|
| `paragraph` | Paragraph text |
| `heading` | `# Heading` (levels 1-6) |
| `blockquote` | `> Quoted text` |
| `codeBlock` | `` ```language `` fenced code |
| `bulletList` / `listItem` | `- Item` |
| `orderedList` / `listItem` | `1. Item` |
| `rule` | `---` |
| `hardBreak` | Two trailing spaces or `\` |
| `text` + `strong` | `**bold**` |
| `text` + `em` | `*italic*` |
| `text` + `code` | `` `code` `` |
| `text` + `link` | `[text](url)` |
| `text` + `strike` | `~~strikethrough~~` (GFM) |
| `table` | GFM table syntax |
| `taskList` / `taskItem` | `- [x]` / `- [ ]` (GFM) |
| `mediaSingle` / `media` | `![alt](url)` |

## GitHub Convention Extensions

These use syntax that GitHub renders but that isn't part of the GFM spec.

### Panels

ADF `panel` nodes are rendered as GitHub alert syntax:

```markdown
> [!NOTE]
> This is an informational panel.

> [!WARNING]
> This is a warning panel.
```

Panel type mapping:

| ADF `panelType` | Alert keyword |
|---|---|
| `info` | `NOTE` |
| `note` | `NOTE` |
| `warning` | `WARNING` |
| `error` | `CAUTION` |
| `success` | `TIP` |
| `custom` | `NOTE` |

### Expand / Collapsible Sections

ADF `expand` and `nestedExpand` nodes are rendered as HTML `<details>` blocks:

```markdown
<details>
<summary>Click to expand</summary>

Content inside the collapsible section.

</details>
```

If the expand node has no title, the summary defaults to `"Details"`.

### Emoji

ADF `emoji` nodes are rendered using the `shortName` attribute:

```markdown
:smile:
```

The `id` and `text` attributes from ADF are not preserved.

## Custom Inline Extensions

These are syntax conventions defined specifically for ADF round-tripping. They are human-readable in any markdown viewer but only render as rich ADF nodes when parsed by goldmark-adf.

### Status

ADF `status` nodes are rendered with the text and color:

```markdown
[status:In Progress|yellow]
```

Format: `[status:TEXT|COLOR]`

Valid colors: `neutral`, `purple`, `blue`, `red`, `yellow`, `green`.

Special characters `|` and `]` in the text field are backslash-escaped:

```markdown
[status:Needs Review \| Approval|blue]
```

### Mentions

ADF `mention` nodes are rendered with display name and account ID:

```markdown
@[Jane Smith](abc-123)
```

Format: `@[DISPLAY_NAME](ACCOUNT_ID)`

If the mention has no display name, the ID is used for both:

```markdown
@[abc-123](abc-123)
```

Special characters `]` in the name and `)` in the ID are backslash-escaped.

### Dates

ADF `date` nodes are rendered with the timestamp:

```markdown
[date:1582152559]
```

Format: `[date:TIMESTAMP]`

The timestamp is the epoch milliseconds value from the ADF `timestamp` attribute.

### Placeholders

ADF `placeholder` nodes are rendered in double braces:

```markdown
{{Enter your name}}
```

Format: `{{TEXT}}`

The `}` character is escaped as `\}` when followed by another `}` in the text.

### Cards

ADF card nodes (`inlineCard`, `blockCard`, `embedCard`) are rendered with the URL:

```markdown
[card:https://atlassian.com/project]
[embed:https://youtube.com/watch?v=abc]
```

Inline and block cards use the same `[card:url]` syntax. The distinction is positional: a `[card:url]` that is the sole content of a paragraph becomes a `blockCard`; within other inline content it becomes an `inlineCard`.

Embed cards use `[embed:url]` to distinguish them from regular cards.

### Decision Lists

ADF `decisionList` / `decisionItem` nodes are rendered as list items with decision markers:

```markdown
- [!] Use json/v2 for performance
- [?] Pending design review
```

| Marker | ADF `state` |
|---|---|
| `[!]` | `DECIDED` |
| `[?]` | Any other state |

## Custom Block Extensions

### Layout Sections

ADF `layoutSection` / `layoutColumn` nodes are rendered with markers:

```markdown
[layout-section]

[layout-column 1]

Left column content.

[layout-column 2]

Right column content.
```

### Extension Nodes

ADF `extension`, `bodiedExtension`, `inlineExtension`, `syncBlock`, and `bodiedSyncBlock` nodes are rendered with type and key:

```markdown
[extension:extension:com.vendor.macro]
```

Bodied extensions include the body content after the marker. The `:` and `]` characters in the key are backslash-escaped.

Custom extension handlers (`WithExtensionBlockHandler`, `WithExtensionInlineHandler`) take precedence over this default syntax when provided.

## HTML Inline Marks

ADF marks that have no markdown equivalent are rendered as inline HTML:

| ADF Mark | HTML Output |
|---|---|
| `underline` | `<u>text</u>` |
| `subsup` (sub) | `<sub>text</sub>` |
| `subsup` (sup) | `<sup>text</sup>` |
| `textColor` | `<span style="color:#hex">text</span>` |

## Escaping

Custom inline tokens use backslash escaping for delimiter characters within user-controlled text fields. This follows the same convention as standard markdown.

| Context | Characters escaped |
|---|---|
| `[status:...\|...]` text field | `\|` and `]` |
| `@[...](...)` name field | `]` |
| `@[...](...)` id field | `)` |
| `{{...}}` text field | `}` (when followed by `}`) |
| `[extension:...:...]` key field | `:` and `]` |
