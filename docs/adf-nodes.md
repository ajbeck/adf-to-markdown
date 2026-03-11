# ADF Node Reference

This document covers ADF nodes that are defined in the [ADF JSON schema](https://unpkg.com/@atlaskit/adf-schema@52.2.2/dist/json-schema/v1/full.json) but lack entries in the [official Atlassian ADF documentation](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/). It complements the markdown extension documentation in [extensions.md](extensions.md).

For nodes that have official documentation, links are provided in [schema-coverage.md](schema-coverage.md).

## Inline Nodes

### status

A colored badge displaying a short text label. Used in Jira for workflow states, labels, and tags.

**ADF JSON:**
```json
{
  "type": "status",
  "attrs": {
    "text": "In Progress",
    "color": "yellow",
    "localId": "abc-123",
    "style": ""
  }
}
```

| Attribute | Required | Type | Description |
|---|---|---|---|
| `text` | yes | string (min 1 char) | Display text for the badge |
| `color` | yes | enum | Badge color: `neutral`, `purple`, `blue`, `red`, `yellow`, `green` |
| `localId` | no | string | Internal identifier (not preserved in markdown) |
| `style` | no | string | Visual style variant (not preserved in markdown) |

**Content:** None (leaf node).

**Markdown:** `[status:In Progress|yellow]` — see [extensions.md](extensions.md#status).

---

### placeholder

A placeholder element that prompts users to enter content. Typically rendered as a greyed-out hint in Atlassian editors.

**ADF JSON:**
```json
{
  "type": "placeholder",
  "attrs": {
    "text": "Enter your name"
  }
}
```

| Attribute | Required | Type | Description |
|---|---|---|---|
| `text` | yes | string | Placeholder label |
| `localId` | no | string | Internal identifier |

**Content:** None (leaf node).

**Markdown:** `{{Enter your name}}` — see [extensions.md](extensions.md#placeholders).

---

### mediaInline

An inline media element (image, file, or link) that appears within paragraph text rather than as a standalone block.

**ADF JSON:**
```json
{
  "type": "mediaInline",
  "attrs": {
    "id": "media-id-123",
    "collection": "collection-name",
    "type": "image",
    "alt": "Description"
  }
}
```

| Attribute | Required | Type | Description |
|---|---|---|---|
| `id` | yes | string (min 1 char) | Media file identifier |
| `collection` | yes | string | Media collection name |
| `type` | no | enum | `link`, `file`, or `image` |
| `alt` | no | string | Alt text |
| `localId` | no | string | Internal identifier |
| `occurrenceKey` | no | string | Unique occurrence key |
| `width` | no | number | Display width |
| `height` | no | number | Display height |
| `data` | no | object | Additional media data |

**Content:** None (inline leaf node).

**Marks:** `link`, `annotation`, `border`.

**Markdown:** `![alt](url)` inline within paragraph text.

---

## Block Nodes

### blockCard

A block-level smart link that renders as a rich preview card. Unlike `inlineCard` which appears inline, `blockCard` occupies its own block.

**ADF JSON:**
```json
{
  "type": "blockCard",
  "attrs": {
    "url": "https://atlassian.com/project",
    "localId": "card-123"
  }
}
```

| Attribute | Required | Type | Description |
|---|---|---|---|
| `url` | yes* | string | Card URL (*one of `url`, `data`, or `datasource` is required) |
| `data` | yes* | object | Resolved card data |
| `datasource` | yes* | object | Datasource with `id`, `parameters`, `views` |
| `localId` | no | string | Internal identifier |
| `width` | no | number | Display width |
| `layout` | no | enum | `wide`, `full-width`, `center`, `wrap-right`, `wrap-left`, `align-end`, `align-start` |

**Content:** None (leaf node).

**Markdown:** `[card:url]` as sole paragraph content — see [extensions.md](extensions.md#cards).

---

### embedCard

An embedded rich media card, typically used for video embeds, interactive content, or full-page previews.

**ADF JSON:**
```json
{
  "type": "embedCard",
  "attrs": {
    "url": "https://youtube.com/watch?v=abc",
    "layout": "center"
  }
}
```

| Attribute | Required | Type | Description |
|---|---|---|---|
| `url` | yes | string | Embed URL |
| `layout` | yes | enum | `wide`, `full-width`, `center`, `wrap-right`, `wrap-left`, `align-end`, `align-start` |
| `width` | no | number (0-100) | Width as percentage |
| `originalHeight` | no | number | Original content height |
| `originalWidth` | no | number | Original content width |
| `localId` | no | string | Internal identifier |

**Content:** None (leaf node).

**Markdown:** `[embed:url]` — see [extensions.md](extensions.md#cards).

---

### caption

A caption element that appears as a child of `mediaSingle` to provide a text description below an image.

**ADF JSON:**
```json
{
  "type": "caption",
  "content": [
    {"type": "text", "text": "Figure 1: Architecture diagram"}
  ]
}
```

| Attribute | Required | Type | Description |
|---|---|---|---|
| `localId` | no | string | Internal identifier |

**Content:** Inline nodes including text, hardBreak, mention, emoji, date, placeholder, inlineCard, status.

**Markdown:** Image title syntax: `![alt](url "Caption text")`.

---

### taskList

A container for task items, rendered as a checklist in Atlassian editors.

**ADF JSON:**
```json
{
  "type": "taskList",
  "attrs": {"localId": "list-123"},
  "content": [
    {
      "type": "taskItem",
      "attrs": {"localId": "item-1", "state": "DONE"},
      "content": [{"type": "text", "text": "Completed task"}]
    },
    {
      "type": "taskItem",
      "attrs": {"localId": "item-2", "state": "TODO"},
      "content": [{"type": "text", "text": "Pending task"}]
    }
  ]
}
```

| Attribute | Required | Type | Description |
|---|---|---|---|
| `localId` | yes | string | List identifier |

**Content:** First item must be `taskItem` or `blockTaskItem`. Subsequent items can also include nested `taskList` for hierarchical checklists.

**Markdown:** `- [x]` / `- [ ]` — standard GFM task list syntax.

---

### taskItem

A single task within a `taskList`. Contains inline content and a completion state.

**ADF JSON:**
```json
{
  "type": "taskItem",
  "attrs": {"localId": "item-1", "state": "TODO"},
  "content": [{"type": "text", "text": "Review the pull request"}]
}
```

| Attribute | Required | Type | Description |
|---|---|---|---|
| `localId` | yes | string | Item identifier |
| `state` | yes | enum | `TODO` or `DONE` |

**Content:** Inline nodes (text, mentions, emoji, links, etc.).

**Markdown:** `- [x] text` (DONE) or `- [ ] text` (TODO).

---

### blockTaskItem

A block-level variant of `taskItem` that contains paragraph content instead of inline content. Used in some Jira contexts.

**ADF JSON:**
```json
{
  "type": "blockTaskItem",
  "attrs": {"localId": "item-1", "state": "DONE"},
  "content": [
    {
      "type": "paragraph",
      "content": [{"type": "text", "text": "Completed task"}]
    }
  ]
}
```

| Attribute | Required | Type | Description |
|---|---|---|---|
| `localId` | yes | string | Item identifier |
| `state` | yes | enum | `TODO` or `DONE` |

**Content:** 1-2 items: `paragraph` (no marks variant) and optionally an `extension`.

**Markdown:** `- [x] text` / `- [ ] text` — same as `taskItem`.

---

### decisionList

A container for decision items, used in Jira and Confluence for recording team decisions.

**ADF JSON:**
```json
{
  "type": "decisionList",
  "attrs": {"localId": "dec-list-1"},
  "content": [
    {
      "type": "decisionItem",
      "attrs": {"localId": "dec-1", "state": "DECIDED"},
      "content": [{"type": "text", "text": "Use json/v2 for performance"}]
    }
  ]
}
```

| Attribute | Required | Type | Description |
|---|---|---|---|
| `localId` | yes | string | List identifier |

**Content:** Array of `decisionItem` nodes (min 1).

**Markdown:** List of `- [!]` / `- [?]` items — see [extensions.md](extensions.md#decision-lists).

---

### decisionItem

A single decision within a `decisionList`.

**ADF JSON:**
```json
{
  "type": "decisionItem",
  "attrs": {"localId": "dec-1", "state": "DECIDED"},
  "content": [{"type": "text", "text": "Use json/v2 for performance"}]
}
```

| Attribute | Required | Type | Description |
|---|---|---|---|
| `localId` | yes | string | Item identifier |
| `state` | yes | string | Decision state (typically `DECIDED` or other values) |

**Content:** Inline nodes (text, mentions, emoji, links, etc.).

**Markdown:** `- [!] text` (DECIDED) or `- [?] text` (other state).

---

### layoutSection

A container for multi-column layouts. Contains 2-3 `layoutColumn` children.

**ADF JSON:**
```json
{
  "type": "layoutSection",
  "content": [
    {
      "type": "layoutColumn",
      "attrs": {"width": 50},
      "content": [
        {"type": "paragraph", "content": [{"type": "text", "text": "Left column"}]}
      ]
    },
    {
      "type": "layoutColumn",
      "attrs": {"width": 50},
      "content": [
        {"type": "paragraph", "content": [{"type": "text", "text": "Right column"}]}
      ]
    }
  ]
}
```

| Attribute | Required | Type | Description |
|---|---|---|---|
| `localId` | no | string | Section identifier |

**Content:** 2-3 `layoutColumn` nodes (full variant).

**Marks:** `breakout` (with `mode`: `wide` or `full-width`).

**Markdown:** `[layout-section]` / `[layout-column N]` markers — see [extensions.md](extensions.md#layout-sections).

---

### layoutColumn

A single column within a `layoutSection`.

**ADF JSON:**
```json
{
  "type": "layoutColumn",
  "attrs": {"width": 50},
  "content": [
    {"type": "paragraph", "content": [{"type": "text", "text": "Column content"}]}
  ]
}
```

| Attribute | Required | Type | Description |
|---|---|---|---|
| `width` | yes | number (0-100) | Column width as a percentage |
| `localId` | no | string | Column identifier |

**Content:** Block content (paragraphs, lists, headings, etc.) with at least 1 item.

**Markdown:** Content under `[layout-column N]` marker.

---

### extension

A block-level Atlassian extension (macro, app, or plugin content). Used by Confluence macros and Jira apps.

**ADF JSON:**
```json
{
  "type": "extension",
  "attrs": {
    "extensionType": "com.atlassian.confluence.macro.core",
    "extensionKey": "toc",
    "parameters": {"maxLevel": "3"}
  }
}
```

| Attribute | Required | Type | Description |
|---|---|---|---|
| `extensionType` | yes | string (min 1 char) | Extension type identifier (typically a reverse domain name) |
| `extensionKey` | yes | string (min 1 char) | Specific extension key (e.g., macro name) |
| `parameters` | no | object | Extension-specific parameters |
| `text` | no | string | Fallback text representation |
| `layout` | no | enum | `wide`, `full-width`, `default` |
| `localId` | no | string | Instance identifier |

**Content:** None (leaf node).

**Marks:** `dataConsumer`, `fragment`.

**Markdown:** `[extension:extension:key]` — see [extensions.md](extensions.md#extension-nodes).

---

### bodiedExtension

A block-level extension that contains body content. Unlike `extension`, this wraps child block nodes that the extension renders around.

**ADF JSON:**
```json
{
  "type": "bodiedExtension",
  "attrs": {
    "extensionType": "com.atlassian.confluence.macro.core",
    "extensionKey": "panel",
    "parameters": {"title": "Note"}
  },
  "content": [
    {"type": "paragraph", "content": [{"type": "text", "text": "Body content"}]}
  ]
}
```

| Attribute | Required | Type | Description |
|---|---|---|---|
| `extensionType` | yes | string (min 1 char) | Extension type identifier |
| `extensionKey` | yes | string (min 1 char) | Specific extension key |
| `parameters` | no | object | Extension-specific parameters |
| `text` | no | string | Fallback text representation |
| `layout` | no | enum | `wide`, `full-width`, `default` |
| `localId` | no | string | Instance identifier |

**Content:** Block nodes (paragraphs, lists, headings, etc.) with at least 1 item.

**Marks:** `dataConsumer`, `fragment`.

**Markdown:** `[extension:bodiedExtension:key]` followed by body content.

---

### inlineExtension

An inline extension element that appears within paragraph text.

**ADF JSON:**
```json
{
  "type": "inlineExtension",
  "attrs": {
    "extensionType": "com.atlassian.jira",
    "extensionKey": "issue-key"
  }
}
```

| Attribute | Required | Type | Description |
|---|---|---|---|
| `extensionType` | yes | string (min 1 char) | Extension type identifier |
| `extensionKey` | yes | string (min 1 char) | Specific extension key |
| `parameters` | no | object | Extension-specific parameters |
| `text` | no | string | Fallback text representation |
| `localId` | no | string | Instance identifier |

**Content:** None (inline leaf node).

**Marks:** `dataConsumer`, `fragment`.

**Markdown:** `[extension:inlineExtension:key]` inline within paragraph text.

---

### syncBlock

A reference to a synchronized block of content (live pages, whiteboards). Points to shared content by resource ID.

**ADF JSON:**
```json
{
  "type": "syncBlock",
  "attrs": {
    "resourceId": "resource-123",
    "localId": "sync-1"
  }
}
```

| Attribute | Required | Type | Description |
|---|---|---|---|
| `resourceId` | yes | string | Reference to the synchronized resource |
| `localId` | yes | string | Local instance identifier |

**Content:** None (reference node).

**Marks:** `breakout` (with `mode`: `wide` or `full-width`).

**Markdown:** `[extension:syncBlock:key]`.

---

### bodiedSyncBlock

A synchronized block that contains body content. Combines sync reference semantics with inline body content.

**ADF JSON:**
```json
{
  "type": "bodiedSyncBlock",
  "attrs": {
    "resourceId": "resource-123",
    "localId": "bsync-1"
  },
  "content": [
    {"type": "paragraph", "content": [{"type": "text", "text": "Synced content"}]}
  ]
}
```

| Attribute | Required | Type | Description |
|---|---|---|---|
| `resourceId` | yes | string | Reference to the synchronized resource |
| `localId` | yes | string | Local instance identifier |

**Content:** Block nodes including paragraphs, lists, tables, panels, etc.

**Marks:** `breakout` (with `mode`: `wide` or `full-width`).

**Markdown:** `[extension:bodiedSyncBlock:key]` followed by body content.
