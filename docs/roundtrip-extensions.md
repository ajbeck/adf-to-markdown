# ADF Round-Trip Markdown Extensions

This document specifies the markdown syntax conventions used to represent Atlassian Document Format (ADF) nodes that have no native markdown equivalent. These conventions are shared between two libraries:

- **adf-to-markdown** — converts ADF JSON to markdown using these conventions
- **goldmark-adf** — parses markdown (including these conventions) back into ADF JSON

The goal is **lossless round-tripping**: an ADF document converted to markdown and back should produce semantically equivalent ADF.

## Escaping

Several extensions use delimiters within user-controlled text. To prevent ambiguity, both libraries use **backslash escaping** — the same convention markdown itself uses.

When **writing** markdown (adf-to-markdown), delimiter characters inside user text are escaped with `\`. When **parsing** markdown (goldmark-adf), `\` followed by a delimiter is consumed as the literal character.

Each extension section documents which characters require escaping.

### Escape rules

| Context | Characters to escape |
|---|---|
| Inside `[status:...\|...]` text field | `\|` `]` |
| Inside `@[...](...)` name field | `]` |
| Inside `@[...](...)` id field | `)` |
| Inside `{{...}}` text field | `}` (when followed by `}`) |
| Inside `[extension:...:...]` key field | `:` `]` |

Characters that don't conflict with the surrounding delimiter never need escaping. For example, `|` inside a mention name is fine because mention syntax doesn't use `|` as a delimiter.

## Tier 1: GFM-Native Constructs

These use standard GitHub Flavored Markdown syntax. No custom extensions are needed in the markdown itself — only correct ADF node mapping.

### Task Lists

ADF nodes: `taskList`, `taskItem`, `blockTaskItem`

```markdown
- [ ] Incomplete task
- [x] Completed task
```

| Checkbox | ADF `state` attribute |
|---|---|
| `[ ]` | `TODO` |
| `[x]` | `DONE` |

**adf-to-markdown**: Already emits this syntax. No change needed.

**goldmark-adf**: Currently emits `[x]`/`[ ]` as text prefixes inside `bulletList`/`listItem`. Must be changed to emit `taskList`/`taskItem` ADF nodes with `state` and `localId` attributes.

### Tables

ADF nodes: `table`, `tableRow`, `tableHeader`, `tableCell`

Standard GFM table syntax. Both libraries already handle this correctly.

### External Media

ADF nodes: `mediaSingle`, `media` (type: external), `caption`

```markdown
![alt text](https://example.com/image.png)
```

With caption (image title becomes caption):

```markdown
![alt text](https://example.com/image.png "Caption text")
```

**goldmark-adf**: Requires `WithExternalMedia(true)` option.

## Tier 2: GitHub Convention Constructs

These use syntax that GitHub renders but that isn't part of the GFM spec. Both libraries need parser/renderer support.

### Panels (GitHub Alerts)

ADF node: `panel`

```markdown
> [!NOTE]
> This is an info panel.

> [!WARNING]
> Be careful here.
```

**Type mapping:**

| ADF `panelType` | Alert keyword |
|---|---|
| `info` | `NOTE` |
| `note` | `NOTE` |
| `warning` | `WARNING` |
| `error` | `CAUTION` |
| `success` | `TIP` |
| `custom` | `NOTE` |

**Reverse mapping** (goldmark-adf parsing):

| Alert keyword | ADF `panelType` |
|---|---|
| `NOTE` | `info` |
| `TIP` | `success` |
| `IMPORTANT` | `info` |
| `WARNING` | `warning` |
| `CAUTION` | `error` |

**adf-to-markdown**: Already emits this syntax. No change needed.

**goldmark-adf**: New parser extension to detect `[!KEYWORD]` as the first line of a blockquote and emit a `panel` node instead.

### Expand / Collapsible Sections

ADF nodes: `expand`, `nestedExpand`

```markdown
<details>
<summary>Click to expand</summary>

Content inside the collapsible section.

</details>
```

If no title is provided, the summary defaults to `"Details"`.

**adf-to-markdown**: Already emits this syntax. No change needed.

**goldmark-adf**: New parser extension to parse `<details>` HTML blocks (currently skipped) into `expand` ADF nodes.

### Emoji

ADF node: `emoji`

```markdown
:shortcode:
```

The `shortName` attribute from ADF is used directly (it already includes the colons in ADF, e.g. `":smile:"`). The `id` and `text` attributes are not preserved in the markdown representation.

ADF schema: `shortName` (required), `id` (optional), `text` (optional)

**adf-to-markdown**: Currently emits the `text` attr falling back to `shortName`. Must be changed to always emit `shortName` to ensure the shortcode syntax is consistent.

**goldmark-adf**: New parser extension to detect `:shortcode:` patterns and emit `emoji` nodes with `shortName`.

## Tier 3: Custom Inline Syntax

These are conventions defined specifically for ADF round-tripping. They don't render specially in standard markdown viewers but are unambiguous and human-readable.

### Status

ADF node: `status`

```markdown
[status:In Progress|yellow]
```

ADF schema attributes:
- `text` (required, string) — the display text
- `color` (required, enum: `neutral`, `purple`, `blue`, `red`, `yellow`, `green`)
- `localId` (optional) — not preserved
- `style` (optional) — not preserved

Format: `[status:TEXT|COLOR]`

**Escaping**: `|` and `]` in TEXT must be escaped as `\|` and `\]`.

**adf-to-markdown**: Must capture `color` from attrs and emit the new syntax.

**goldmark-adf**: New inline parser to detect `[status:...|...]` and emit `status` nodes.

### Mentions

ADF node: `mention`

```markdown
@[Brad](abc123)
```

ADF schema attributes:
- `id` (required) — the account identifier
- `text` (optional) — display name
- `accessLevel` (optional) — not preserved
- `userType` (optional, enum: `DEFAULT`, `SPECIAL`, `APP`) — not preserved
- `localId` (optional) — not preserved

Format: `@[TEXT](ID)`

If `text` is empty, the `id` is used as both display and identifier: `@[abc123](abc123)`.

**Escaping**: `]` in TEXT must be escaped as `\]`. `)` in ID must be escaped as `\)`.

**adf-to-markdown**: Must change from current `@Brad` output.

**goldmark-adf**: New inline parser to detect `@[...](...)` and emit `mention` nodes.

### Dates

ADF node: `date`

```markdown
[date:1582152559]
```

ADF schema attributes:
- `timestamp` (required, string) — epoch milliseconds as string

Format: `[date:TIMESTAMP]`

No escaping needed — timestamps are numeric strings.

**adf-to-markdown**: Already emits this syntax. No change needed.

**goldmark-adf**: New inline parser to detect `[date:DIGITS]` and emit `date` nodes.

### Placeholders

ADF node: `placeholder`

```markdown
{{name}}
```

ADF schema attributes:
- `text` (required) — the placeholder label

Format: `{{TEXT}}`

**Escaping**: `}` must be escaped as `\}` when followed by another `}` in TEXT.

**adf-to-markdown**: Already emits this syntax. No change needed.

**goldmark-adf**: New inline parser to detect `{{...}}` and emit `placeholder` nodes.

### Inline Cards

ADF node: `inlineCard`

```markdown
[card:https://atlassian.com]
```

ADF schema attributes:
- `url` (required) — the card URL

Format: `[card:URL]`

**Escaping**: `]` in URL must be escaped as `\]` (rare in practice).

**adf-to-markdown**: Must change from current `[url](url)` link output.

**goldmark-adf**: New inline parser to detect `[card:...]` and emit `inlineCard` nodes.

### Block Cards

ADF node: `blockCard`

```markdown
[card:https://example.com/card]
```

Same syntax as inline cards. Distinguished by position: a `[card:...]` token that is the sole content of a paragraph is a block card. Within other inline content, it is an inline card.

**adf-to-markdown**: Must change from current `[url](url)` link output.

**goldmark-adf**: Detect block vs inline by paragraph context.

### Embed Cards

ADF node: `embedCard`

```markdown
[embed:https://youtube.com/watch?v=abc]
```

ADF schema attributes:
- `url` (required) — the embed URL

Format: `[embed:URL]`

**adf-to-markdown**: Must change from current `[url](url)` link output.

**goldmark-adf**: New inline parser to detect `[embed:...]` and emit `embedCard` nodes.

### Decision Lists

ADF nodes: `decisionList`, `decisionItem`

```markdown
- [!] Use json/v2 for performance
- [?] Pending design review
```

| Checkbox | ADF `state` attribute |
|---|---|
| `[!]` | `DECIDED` |
| `[?]` | any other state |

ADF schema attributes on `decisionItem`:
- `localId` (required) — not preserved
- `state` (required) — `DECIDED` or other values

Format: `- [!] text` or `- [?] text`

**adf-to-markdown**: Must change from current `- [decision:STATE] text` output.

**goldmark-adf**: New parser extension to detect `- [!]`/`- [?]` list items and emit `decisionList`/`decisionItem` nodes.

## Tier 4: Complex Structures

### Layout Sections

ADF nodes: `layoutSection`, `layoutColumn`

```markdown
[layout-section]

[layout-column 1]

Left column content.

[layout-column 2]

Right column content.
```

**adf-to-markdown**: Already emits this syntax. No change needed.

**goldmark-adf**: New parser extension to detect layout markers and emit `layoutSection`/`layoutColumn` nodes.

### Media (Non-External)

ADF nodes: `media` (type: file or link), `mediaInline`, `mediaGroup`

File media uses a custom URL scheme to preserve collection and ID:

```markdown
![alt](atlassian-media://collection-id/media-id)
```

Link media:

```markdown
![alt](atlassian-media-link://collection-id/media-id)
```

Media groups are rendered as consecutive images:

```markdown
![A](https://example.com/a.png)
![B](https://example.com/b.png)
```

Inline media uses the same `![alt](url)` syntax but appears inline within paragraph text.

**goldmark-adf**: Recognize `atlassian-media://` and `atlassian-media-link://` URL schemes in image syntax. Group consecutive block-level images into `mediaGroup`. Distinguish inline vs block images.

### Extensions

ADF nodes: `extension`, `bodiedExtension`, `inlineExtension`, `syncBlock`, `bodiedSyncBlock`

```markdown
[extension:extension:com.vendor.macro]
```

With body content:

```markdown
[extension:bodiedExtension:com.vendor.macro]

Body content here.
```

Format: `[extension:TYPE:KEY]`

**Escaping**: `:` and `]` in KEY must be escaped as `\:` and `\]`.

**adf-to-markdown**: Already emits this syntax. No change needed.

**goldmark-adf**: New parser to detect `[extension:...]` markers and emit the appropriate extension ADF nodes.

### Marks Without Markdown Equivalents

| ADF Mark | Markdown | Round-trips? |
|---|---|---|
| `underline` | `<u>text</u>` | Requires HTML inline parsing |
| `subsup` (sub) | `<sub>text</sub>` | Requires HTML inline parsing |
| `subsup` (sup) | `<sup>text</sup>` | Requires HTML inline parsing |
| `textColor` | Not preserved | Loss accepted |
| `backgroundColor` | Not preserved | Loss accepted |

**goldmark-adf**: New parser extension to handle inline HTML tags `<u>`, `<sub>`, `<sup>` and apply the corresponding ADF marks.

## Implementation Phases

### Phase 1: GFM-native (lowest effort)

| Library | Change |
|---|---|
| goldmark-adf | Fix task list rendering to emit `taskList`/`taskItem` ADF nodes |

### Phase 2: GitHub conventions

| Library | Change |
|---|---|
| adf-to-markdown | Normalize emoji output to `:shortName:` |
| goldmark-adf | Panel alert parser extension |
| goldmark-adf | `<details>` expand parser extension |
| goldmark-adf | `:shortcode:` emoji parser extension |

### Phase 3: Custom inline syntax

| Library | Change |
|---|---|
| adf-to-markdown | Add escape helpers |
| adf-to-markdown | Update `status` renderer (add color, new syntax) |
| adf-to-markdown | Update `mention` renderer (new syntax) |
| adf-to-markdown | Update `inlineCard`/`blockCard`/`embedCard` renderers |
| adf-to-markdown | Update `decisionItem` renderer |
| goldmark-adf | Add unescape helpers |
| goldmark-adf | Inline token parsers: status, mention, date, placeholder, cards |
| goldmark-adf | Decision list parser extension |

### Phase 4: Complex structures

| Library | Change |
|---|---|
| goldmark-adf | Layout section parser |
| goldmark-adf | Non-external media handling (URL schemes, grouping) |
| goldmark-adf | Extension marker parser |
| goldmark-adf | Inline HTML mark parsing (`<u>`, `<sub>`, `<sup>`) |

### Phase 5: Validation

| Library | Change |
|---|---|
| adf-to-markdown | Round-trip test harness with comprehensive ADF fixture |
| adf-to-markdown | Documentation: how both libraries work together |
