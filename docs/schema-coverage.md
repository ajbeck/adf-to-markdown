# ADF Schema Coverage (Nodes)

For ADF node details not covered by official Atlassian documentation, see [adf-nodes.md](adf-nodes.md). For the markdown syntax produced by each node, see [extensions.md](extensions.md).

Status labels:

- `implemented`: explicit decode + render support
- `partial`: decoded to fallback/generic behavior; not rich semantics
- `schema-variant`: covered via base node handling, not separate node model
- `missing`: not yet implemented

Markdown type labels:

- `standard`: native CommonMark syntax
- `GFM`: GitHub Flavored Markdown extension (tables, strikethrough, task lists)
- `GitHub convention`: syntax that GitHub renders but isn't part of the GFM spec (alerts, `<details>`)
- `custom extension`: syntax defined specifically for ADF round-tripping
- `HTML inline`: inline HTML tags for marks without markdown equivalents

## Nodes

| Schema Node | Status | Markdown Type | Markdown Syntax | ADF Docs |
|---|---|---|---|---|
| `blockCard_node` | implemented | custom extension | `[card:url]` | - |
| `blockTaskItem_node` | implemented | GFM | `- [x]` / `- [ ]` | - |
| `blockquote_node` | implemented | standard | `> text` | [blockquote](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/blockquote/) |
| `bodiedExtension_node` | partial | custom extension | `[extension:bodiedExtension:key]` | - |
| `bodiedSyncBlock_node` | partial | custom extension | `[extension:bodiedSyncBlock:key]` | - |
| `bulletList_node` | implemented | standard | `- item` | [bulletList](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/bulletList/) |
| `caption_node` | partial | standard | `![alt](url "caption")` (image title) | - |
| `codeBlock_node` | implemented | standard | `` ```lang `` | [codeBlock](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/codeBlock/) |
| `date_node` | implemented | custom extension | `[date:timestamp]` | [date](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/date/) |
| `decisionItem_node` | implemented | custom extension | `- [!] text` / `- [?] text` | - |
| `decisionList_node` | implemented | custom extension | (list of decision items) | - |
| `doc_node` | implemented | standard | (document root) | [doc](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/doc/) |
| `embedCard_node` | implemented | custom extension | `[embed:url]` | - |
| `emoji_node` | implemented | GitHub convention | `:shortcode:` | [emoji](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/emoji/) |
| `expand_node` | implemented | GitHub convention | `<details><summary>` | [expand](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/expand/) |
| `extension_node` | partial | custom extension | `[extension:extension:key]` | - |
| `hardBreak_node` | implemented | standard | two spaces + newline | [hardBreak](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/hardBreak/) |
| `heading_node` | implemented | standard | `# heading` | [heading](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/heading/) |
| `inlineCard_node` | implemented | custom extension | `[card:url]` | [inlineCard](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/inlineCard/) |
| `inlineExtension_node` | partial | custom extension | `[extension:inlineExtension:key]` | - |
| `layoutColumn_node` | implemented | custom extension | `[layout-column N]` | - |
| `layoutSection_node` | implemented | custom extension | `[layout-section]` | - |
| `listItem_node` | implemented | standard | `- item` or `1. item` | [listItem](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/listItem/) |
| `mediaGroup_node` | implemented | standard | consecutive `![alt](url)` | [mediaGroup](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/mediaGroup/) |
| `mediaInline_node` | implemented | standard | `![alt](url)` inline | - |
| `mediaSingle_node` | implemented | standard | `![alt](url)` | [mediaSingle](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/mediaSingle/) |
| `media_node` | implemented | standard | `![alt](url)` | [media](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/media/) |
| `mention_node` | implemented | custom extension | `@[name](id)` | [mention](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/mention/) |
| `nestedExpand_node` | implemented | GitHub convention | `<details><summary>` | [nestedExpand](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/nestedExpand/) |
| `orderedList_node` | implemented | standard | `1. item` | [orderedList](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/orderedList/) |
| `panel_node` | implemented | GitHub convention | `> [!TYPE]` | [panel](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/panel/) |
| `paragraph_node` | implemented | standard | (paragraph text) | [paragraph](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/paragraph/) |
| `placeholder_node` | implemented | custom extension | `{{text}}` | - |
| `rule_node` | implemented | standard | `---` | [rule](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/rule/) |
| `status_node` | implemented | custom extension | `[status:text\|color]` | [status](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/status/) |
| `syncBlock_node` | partial | custom extension | `[extension:syncBlock:key]` | - |
| `table_cell_node` | implemented | GFM | `\| cell \|` | [tableCell](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/table_cell/) |
| `table_header_node` | implemented | GFM | `\| header \|` | [tableHeader](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/table_header/) |
| `table_node` | implemented | GFM | `\| h1 \| h2 \|` | [table](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/table/) |
| `table_row_node` | implemented | GFM | `\| c1 \| c2 \|` | [tableRow](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/table_row/) |
| `taskItem_node` | implemented | GFM | `- [x]` / `- [ ]` | - |
| `taskList_node` | implemented | GFM | (list of task items) | - |
| `text_node` | implemented | standard | (inline text) | [text](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/text/) |

## Marks

| ADF Mark | Markdown Type | Markdown Syntax |
|---|---|---|
| `strong` | standard | `**text**` |
| `em` | standard | `*text*` |
| `code` | standard | `` `text` `` |
| `link` | standard | `[text](url)` |
| `strike` | GFM | `~~text~~` |
| `underline` | HTML inline | `<u>text</u>` |
| `subsup` (sub) | HTML inline | `<sub>text</sub>` |
| `subsup` (sup) | HTML inline | `<sup>text</sup>` |
| `textColor` | HTML inline | `<span style="color:#hex">text</span>` |
| `backgroundColor` | - | Not preserved (loss accepted) |

## Schema Variants

The following are schema variants of the above and are handled through base-node support:

- `codeBlock_root_only_node`
- `expand_root_only_node`
- `extension_with_marks_node`
- `formatted_text_inline_node`
- `heading_with_alignment_node`
- `heading_with_indentation_node`
- `heading_with_no_marks_node`
- `inlineExtension_with_marks_node`
- `layoutSection_full_node`
- `mediaSingle_caption_node`
- `mediaSingle_full_node`
- `nestedExpand_with_no_marks_node`
- `paragraph_with_alignment_node`
- `paragraph_with_indentation_node`
- `paragraph_with_no_marks_node`
- `text_with_no_marks_node`
- `bodiedExtension_with_marks_node`
