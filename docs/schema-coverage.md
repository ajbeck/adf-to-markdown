# ADF Schema Coverage (Nodes)

Status labels:

- `implemented`: explicit decode + render support
- `partial`: decoded to fallback/generic behavior; not rich semantics
- `schema-variant`: covered via base node handling, not separate node model
- `missing`: not yet implemented

## Nodes

| Schema Node | Status | Notes | ADF Docs |
|---|---|---|---|
| `blockCard_node` | implemented | URL fallback render | - |
| `blockTaskItem_node` | implemented | Rendered as markdown checklist item | - |
| `blockquote_node` | implemented | | [blockquote](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/blockquote/) |
| `bodiedExtension_node` | partial | Extension fallback + handler hooks | - |
| `bodiedSyncBlock_node` | partial | Extension fallback + handler hooks | - |
| `bulletList_node` | implemented | | [bulletList](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/bulletList/) |
| `caption_node` | partial | Supported as `mediaSingle` child path | - |
| `codeBlock_node` | implemented | | [codeBlock](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/codeBlock/) |
| `date_node` | implemented | Inline token output | [date](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/date/) |
| `decisionItem_node` | implemented | Decision list item fallback syntax | - |
| `decisionList_node` | implemented | | - |
| `doc_node` | implemented | | [doc](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/doc/) |
| `embedCard_node` | implemented | URL fallback render | - |
| `emoji_node` | implemented | | [emoji](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/emoji/) |
| `expand_node` | implemented | | [expand](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/expand/) |
| `extension_node` | partial | Extension fallback + handler hooks | - |
| `hardBreak_node` | implemented | | [hardBreak](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/hardBreak/) |
| `heading_node` | implemented | | [heading](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/heading/) |
| `inlineCard_node` | implemented | | [inlineCard](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/inlineCard/) |
| `inlineExtension_node` | partial | Extension fallback + handler hooks | - |
| `layoutColumn_node` | implemented | Via `layoutSection` content | - |
| `layoutSection_node` | implemented | | - |
| `listItem_node` | implemented | | [listItem](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/listItem/) |
| `mediaGroup_node` | implemented | | [mediaGroup](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/mediaGroup/) |
| `mediaInline_node` | implemented | | - |
| `mediaSingle_node` | implemented | Includes caption handling | [mediaSingle](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/mediaSingle/) |
| `media_node` | implemented | External/file/link mapped to markdown/image forms | [media](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/media/) |
| `mention_node` | implemented | | [mention](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/mention/) |
| `nestedExpand_node` | implemented | Mapped with `expand` behavior | [nestedExpand](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/nestedExpand/) |
| `orderedList_node` | implemented | | [orderedList](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/orderedList/) |
| `panel_node` | implemented | | [panel](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/panel/) |
| `paragraph_node` | implemented | | [paragraph](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/paragraph/) |
| `placeholder_node` | implemented | | - |
| `rule_node` | implemented | | [rule](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/rule/) |
| `status_node` | implemented | | [status](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/status/) |
| `syncBlock_node` | partial | Extension fallback + handler hooks | - |
| `table_cell_node` | implemented | | [tableCell](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/table_cell/) |
| `table_header_node` | implemented | | [tableHeader](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/table_header/) |
| `table_node` | implemented | | [table](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/table/) |
| `table_row_node` | implemented | | [tableRow](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/table_row/) |
| `taskItem_node` | implemented | | - |
| `taskList_node` | implemented | | - |
| `text_node` | implemented | | [text](https://developer.atlassian.com/cloud/jira/platform/apis/document/nodes/text/) |

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

