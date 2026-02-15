# ADF Schema Coverage (Nodes)

Status labels:

- `implemented`: explicit decode + render support
- `partial`: decoded to fallback/generic behavior; not rich semantics
- `schema-variant`: covered via base node handling, not separate node model
- `missing`: not yet implemented

## Nodes

| Schema Node | Status | Notes |
|---|---|---|
| `blockCard_node` | implemented | URL fallback render |
| `blockTaskItem_node` | implemented | Rendered as markdown checklist item |
| `blockquote_node` | implemented | |
| `bodiedExtension_node` | partial | Extension fallback + handler hooks |
| `bodiedSyncBlock_node` | partial | Extension fallback + handler hooks |
| `bulletList_node` | implemented | |
| `caption_node` | partial | Supported as `mediaSingle` child path |
| `codeBlock_node` | implemented | |
| `date_node` | implemented | Inline token output |
| `decisionItem_node` | implemented | Decision list item fallback syntax |
| `decisionList_node` | implemented | |
| `doc_node` | implemented | |
| `embedCard_node` | implemented | URL fallback render |
| `emoji_node` | implemented | |
| `expand_node` | implemented | |
| `extension_node` | partial | Extension fallback + handler hooks |
| `hardBreak_node` | implemented | |
| `heading_node` | implemented | |
| `inlineCard_node` | implemented | |
| `inlineExtension_node` | partial | Extension fallback + handler hooks |
| `layoutColumn_node` | implemented | Via `layoutSection` content |
| `layoutSection_node` | implemented | |
| `listItem_node` | implemented | |
| `mediaGroup_node` | implemented | |
| `mediaInline_node` | implemented | |
| `mediaSingle_node` | implemented | Includes caption handling |
| `media_node` | implemented | External/file/link mapped to markdown/image forms |
| `mention_node` | implemented | |
| `nestedExpand_node` | implemented | Mapped with `expand` behavior |
| `orderedList_node` | implemented | |
| `panel_node` | implemented | |
| `paragraph_node` | implemented | |
| `placeholder_node` | implemented | |
| `rule_node` | implemented | |
| `status_node` | implemented | |
| `syncBlock_node` | partial | Extension fallback + handler hooks |
| `table_cell_node` | implemented | |
| `table_header_node` | implemented | |
| `table_node` | implemented | |
| `table_row_node` | implemented | |
| `taskItem_node` | implemented | |
| `taskList_node` | implemented | |
| `text_node` | implemented | |

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

