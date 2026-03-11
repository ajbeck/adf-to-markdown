//go:build goexperiment.jsonv2

package adfmarkdown

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
)

type decoder struct {
	cfg config
}

type rootEnvelope struct {
	Version int              `json:"version"`
	Type    string           `json:"type"`
	Content []jsontext.Value `json:"content"`
}

type nodeEnvelope struct {
	Type    string           `json:"type"`
	Attrs   map[string]any   `json:"attrs"`
	Content []jsontext.Value `json:"content"`
	Text    string           `json:"text"`
	Marks   []mark           `json:"marks"`
}

func newDecoder(cfg config) *decoder {
	return &decoder{cfg: cfg}
}

func (d *decoder) decodeDocument(data []byte) (document, error) {
	var root rootEnvelope
	if err := json.Unmarshal(data, &root); err != nil {
		return document{}, err
	}
	if root.Type != "doc" {
		return document{}, newDecodeError("/", ErrKindInvalidRoot, "expected type doc")
	}
	if d.cfg.StrictSchema && root.Version != 1 {
		return document{}, newDecodeError("/", ErrKindInvalidRoot, "expected version 1")
	}
	doc := document{
		Version: root.Version,
		Content: make([]blockNode, 0, len(root.Content)),
	}
	for i, raw := range root.Content {
		node, err := d.decodeBlock(raw, fmt.Sprintf("/content/%d", i))
		if err != nil {
			return document{}, err
		}
		doc.Content = append(doc.Content, node)
	}
	return doc, nil
}

func (d *decoder) decodeEnvelope(raw jsontext.Value, path string) (nodeEnvelope, error) {
	var env nodeEnvelope
	if err := json.Unmarshal([]byte(raw), &env); err != nil {
		return nodeEnvelope{}, newDecodeError(path, ErrKindInvalidJSON, err.Error())
	}
	if env.Type == "" {
		return nodeEnvelope{}, newDecodeError(path, ErrKindMissingType, "node type is required")
	}
	return env, nil
}

func (d *decoder) decodeBlock(raw jsontext.Value, path string) (blockNode, error) {
	env, err := d.decodeEnvelope(raw, path)
	if err != nil {
		return nil, err
	}
	switch env.Type {
	case "paragraph":
		content, err := d.decodeInlineContent(env.Content, path+"/content")
		if err != nil {
			return nil, err
		}
		return paragraphNode{Content: content}, nil
	case "heading":
		levelAny, ok := env.Attrs["level"]
		if !ok {
			return nil, newDecodeError(path, ErrKindInvalidAttr, "heading requires attrs.level")
		}
		level, err := mustInt(levelAny)
		if err != nil {
			return nil, newDecodeError(path, ErrKindInvalidAttr, "heading attrs.level must be integer")
		}
		content, err := d.decodeInlineContent(env.Content, path+"/content")
		if err != nil {
			return nil, err
		}
		return headingNode{Level: level, Content: content}, nil
	case "blockquote":
		content := make([]blockNode, 0, len(env.Content))
		for i, c := range env.Content {
			n, err := d.decodeBlock(c, fmt.Sprintf("%s/content/%d", path, i))
			if err != nil {
				return nil, err
			}
			content = append(content, n)
		}
		return blockquoteNode{Content: content}, nil
	case "codeBlock":
		lang := ""
		if v, ok := env.Attrs["language"]; ok {
			s, ok := v.(string)
			if !ok {
				return nil, newDecodeError(path, ErrKindInvalidAttr, "codeBlock attrs.language must be string")
			}
			lang = s
		}
		text, err := d.decodeCodeBlockText(env.Content, path+"/content")
		if err != nil {
			return nil, err
		}
		return codeBlockNode{Language: lang, Text: text}, nil
	case "rule":
		return ruleNode{}, nil
	case "panel":
		pt := ""
		if v, ok := env.Attrs["panelType"]; ok {
			if s, ok := v.(string); ok {
				pt = s
			}
		}
		content := make([]blockNode, 0, len(env.Content))
		for i, c := range env.Content {
			n, err := d.decodeBlock(c, fmt.Sprintf("%s/content/%d", path, i))
			if err != nil {
				return nil, err
			}
			content = append(content, n)
		}
		return panelNode{PanelType: pt, Content: content}, nil
	case "expand", "nestedExpand":
		title := ""
		if v, ok := env.Attrs["title"]; ok {
			if s, ok := v.(string); ok {
				title = s
			}
		}
		content := make([]blockNode, 0, len(env.Content))
		for i, c := range env.Content {
			n, err := d.decodeBlock(c, fmt.Sprintf("%s/content/%d", path, i))
			if err != nil {
				return nil, err
			}
			content = append(content, n)
		}
		return expandNode{Title: title, Content: content}, nil
	case "mediaSingle":
		content := make([]inlineNode, 0, 1)
		caption := make([]inlineNode, 0)
		for i, c := range env.Content {
			p := fmt.Sprintf("%s/content/%d", path, i)
			itemEnv, err := d.decodeEnvelope(c, p)
			if err != nil {
				return nil, err
			}
			switch itemEnv.Type {
			case "media":
				in, err := d.decodeMediaInlineFromBlock(c, p)
				if err != nil {
					return nil, err
				}
				content = append(content, in)
			case "caption":
				in, err := d.decodeInlineContent(itemEnv.Content, p+"/content")
				if err != nil {
					return nil, err
				}
				caption = append(caption, in...)
			default:
				if d.cfg.StrictSchema {
					return nil, newDecodeError(p, ErrKindInvalidStructure, "mediaSingle content must be media or caption")
				}
			}
		}
		return mediaSingleNode{Content: content, Caption: caption}, nil
	case "mediaGroup":
		content := make([]inlineNode, 0, len(env.Content))
		for i, c := range env.Content {
			in, err := d.decodeMediaInlineFromBlock(c, fmt.Sprintf("%s/content/%d", path, i))
			if err != nil {
				return nil, err
			}
			content = append(content, in)
		}
		return mediaGroupNode{Content: content}, nil
	case "taskList":
		items, err := d.decodeTaskItems(env.Content, path+"/content")
		if err != nil {
			return nil, err
		}
		return taskListNode{Items: items}, nil
	case "blockTaskItem":
		state, _ := env.Attrs["state"].(string)
		blocks := make([]blockNode, 0, len(env.Content))
		for i, raw := range env.Content {
			b, err := d.decodeBlock(raw, fmt.Sprintf("%s/content/%d", path, i))
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, b)
		}
		return blockTaskItemNode{
			State:   state,
			Content: blocks,
		}, nil
	case "decisionList":
		items, err := d.decodeDecisionItems(env.Content, path+"/content")
		if err != nil {
			return nil, err
		}
		return decisionListNode{Items: items}, nil
	case "blockCard", "embedCard":
		url := extractCardURL(env.Attrs)
		return cardBlockNode{URL: url, NodeType: env.Type}, nil
	case "extension", "bodiedExtension", "extensionFrame", "multiBodiedExtension", "bodiedSyncBlock", "syncBlock":
		key := extractExtensionKey(env.Attrs)
		content := make([]blockNode, 0, len(env.Content))
		for i, c := range env.Content {
			n, err := d.decodeBlock(c, fmt.Sprintf("%s/content/%d", path, i))
			if err != nil {
				return nil, err
			}
			content = append(content, n)
		}
		return extensionBlockNode{
			NodeType: env.Type,
			Key:      key,
			Path:     path,
			Content:  content,
		}, nil
	case "layoutSection":
		cols, err := d.decodeLayoutColumns(env.Content, path+"/content")
		if err != nil {
			return nil, err
		}
		return layoutSectionNode{Columns: cols}, nil
	case "table":
		rows, err := d.decodeTableRows(env.Content, path+"/content")
		if err != nil {
			return nil, err
		}
		return tableNode{Rows: rows}, nil
	case "bulletList":
		items, err := d.decodeListItems(env.Content, path+"/content")
		if err != nil {
			return nil, err
		}
		return bulletListNode{Items: items}, nil
	case "orderedList":
		start := 1
		if order, ok := env.Attrs["order"]; ok {
			v, err := mustInt(order)
			if err != nil {
				return nil, newDecodeError(path, ErrKindInvalidAttr, "orderedList attrs.order must be integer")
			}
			start = v
		}
		items, err := d.decodeListItems(env.Content, path+"/content")
		if err != nil {
			return nil, err
		}
		return orderedListNode{Start: start, Items: items}, nil
	case "listItem":
		blocks := make([]blockNode, 0, len(env.Content))
		for i, itemRaw := range env.Content {
			block, err := d.decodeBlock(itemRaw, fmt.Sprintf("%s/content/%d", path, i))
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, block)
		}
		return listItemNode{Content: blocks}, nil
	default:
		if d.cfg.UnsupportedBlock != nil {
			md, ok, err := d.cfg.UnsupportedBlock(env.Type, path)
			if err != nil {
				return nil, err
			}
			if ok {
				return rawMarkdownBlockNode{Markdown: md}, nil
			}
		}
		if d.cfg.AllowUnsupportedNodes {
			return unsupportedBlockNode{Type: env.Type}, nil
		}
		return nil, newDecodeError(path, ErrKindUnsupportedNode, env.Type)
	}
}

func (d *decoder) decodeLayoutColumns(raws []jsontext.Value, path string) ([]layoutColumnNode, error) {
	cols := make([]layoutColumnNode, 0, len(raws))
	for i, raw := range raws {
		p := fmt.Sprintf("%s/%d", path, i)
		env, err := d.decodeEnvelope(raw, p)
		if err != nil {
			return nil, err
		}
		if env.Type != "layoutColumn" {
			return nil, newDecodeError(p, ErrKindInvalidStructure, "layoutSection content must be layoutColumn nodes")
		}
		content := make([]blockNode, 0, len(env.Content))
		for j, c := range env.Content {
			n, err := d.decodeBlock(c, fmt.Sprintf("%s/content/%d", p, j))
			if err != nil {
				return nil, err
			}
			content = append(content, n)
		}
		cols = append(cols, layoutColumnNode{Content: content})
	}
	return cols, nil
}

func (d *decoder) decodeTaskItems(raws []jsontext.Value, path string) ([]taskItemNode, error) {
	items := make([]taskItemNode, 0, len(raws))
	for i, raw := range raws {
		p := fmt.Sprintf("%s/%d", path, i)
		env, err := d.decodeEnvelope(raw, p)
		if err != nil {
			return nil, err
		}
		if env.Type != "taskItem" {
			return nil, newDecodeError(p, ErrKindInvalidStructure, "taskList content must be taskItem nodes")
		}
		state, _ := env.Attrs["state"].(string)
		content, err := d.decodeInlineContent(env.Content, p+"/content")
		if err != nil {
			return nil, err
		}
		items = append(items, taskItemNode{
			State:   state,
			Content: content,
		})
	}
	return items, nil
}

func (d *decoder) decodeDecisionItems(raws []jsontext.Value, path string) ([]decisionItemNode, error) {
	items := make([]decisionItemNode, 0, len(raws))
	for i, raw := range raws {
		p := fmt.Sprintf("%s/%d", path, i)
		env, err := d.decodeEnvelope(raw, p)
		if err != nil {
			return nil, err
		}
		if env.Type != "decisionItem" {
			return nil, newDecodeError(p, ErrKindInvalidStructure, "decisionList content must be decisionItem nodes")
		}
		state, _ := env.Attrs["state"].(string)
		content, err := d.decodeInlineContent(env.Content, p+"/content")
		if err != nil {
			return nil, err
		}
		items = append(items, decisionItemNode{
			State:   state,
			Content: content,
		})
	}
	return items, nil
}

func extractCardURL(attrs map[string]any) string {
	if u, _ := attrs["url"].(string); u != "" {
		return u
	}
	if data, ok := attrs["data"].(map[string]any); ok {
		if u := findURLLike(data); u != "" {
			return u
		}
	}
	return ""
}

func findURLLike(m map[string]any) string {
	return findURLLikeAny(m)
}

func findURLLikeAny(v any) string {
	switch x := v.(type) {
	case map[string]any:
		for _, k := range []string{"url", "href", "@id"} {
			if u, _ := x[k].(string); u != "" {
				return u
			}
		}
		for _, child := range x {
			if u := findURLLikeAny(child); u != "" {
				return u
			}
		}
	case []any:
		for _, child := range x {
			if u := findURLLikeAny(child); u != "" {
				return u
			}
		}
	}
	return ""
}

func extractExtensionKey(attrs map[string]any) string {
	if k, _ := attrs["extensionKey"].(string); k != "" {
		return k
	}
	if k, _ := attrs["extensionType"].(string); k != "" {
		return k
	}
	if k, _ := attrs["localId"].(string); k != "" {
		return k
	}
	return ""
}

func (d *decoder) decodeTableRows(raws []jsontext.Value, path string) ([]tableRowNode, error) {
	rows := make([]tableRowNode, 0, len(raws))
	for i, raw := range raws {
		p := fmt.Sprintf("%s/%d", path, i)
		env, err := d.decodeEnvelope(raw, p)
		if err != nil {
			return nil, err
		}
		if env.Type != "tableRow" {
			return nil, newDecodeError(p, ErrKindInvalidStructure, "table content must be tableRow nodes")
		}
		row := tableRowNode{Cells: make([]tableCellNode, 0, len(env.Content))}
		for j, cellRaw := range env.Content {
			cellPath := fmt.Sprintf("%s/content/%d", p, j)
			cellEnv, err := d.decodeEnvelope(cellRaw, cellPath)
			if err != nil {
				return nil, err
			}
			isHeader := false
			switch cellEnv.Type {
			case "tableCell":
			case "tableHeader":
				isHeader = true
			default:
				return nil, newDecodeError(cellPath, ErrKindInvalidStructure, "tableRow content must be tableCell or tableHeader")
			}
			blocks := make([]blockNode, 0, len(cellEnv.Content))
			for k, braw := range cellEnv.Content {
				b, err := d.decodeBlock(braw, fmt.Sprintf("%s/content/%d", cellPath, k))
				if err != nil {
					return nil, err
				}
				blocks = append(blocks, b)
			}
			row.Cells = append(row.Cells, tableCellNode{
				IsHeader: isHeader,
				Content:  blocks,
			})
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func (d *decoder) decodeInline(raw jsontext.Value, path string) (inlineNode, error) {
	env, err := d.decodeEnvelope(raw, path)
	if err != nil {
		return nil, err
	}
	switch env.Type {
	case "text":
		if d.cfg.StrictSchema && env.Text == "" {
			return nil, newDecodeError(path, ErrKindInvalidText, "text must be non-empty")
		}
		if err := d.validateMarks(env.Marks, path+"/marks"); err != nil {
			return nil, err
		}
		return textNode{
			Text:  env.Text,
			Marks: env.Marks,
		}, nil
	case "hardBreak":
		return hardBreakNode{}, nil
	case "emoji":
		shortName, _ := env.Attrs["shortName"].(string)
		if shortName == "" {
			shortName = ":emoji:"
		}
		return emojiInlineNode{ShortName: shortName}, nil
	case "mention":
		text, _ := env.Attrs["text"].(string)
		id, _ := env.Attrs["id"].(string)
		return mentionInlineNode{Text: text, ID: id}, nil
	case "date":
		ts, _ := env.Attrs["timestamp"].(string)
		return dateInlineNode{Timestamp: ts}, nil
	case "status":
		text, _ := env.Attrs["text"].(string)
		color, _ := env.Attrs["color"].(string)
		return statusInlineNode{Text: text, Color: color}, nil
	case "inlineCard":
		url, _ := env.Attrs["url"].(string)
		return inlineCardNode{URL: url}, nil
	case "inlineExtension":
		return extensionInlineNode{
			NodeType: env.Type,
			Key:      extractExtensionKey(env.Attrs),
			Path:     path,
		}, nil
	case "placeholder":
		txt, _ := env.Attrs["text"].(string)
		return placeholderInlineNode{Text: txt}, nil
	case "mediaInline":
		u, a, err := d.mediaAttrsToURLAlt(env.Attrs, path)
		if err != nil {
			return nil, err
		}
		return mediaInlineNode{URL: u, Alt: a}, nil
	default:
		if d.cfg.UnsupportedInline != nil {
			md, ok, err := d.cfg.UnsupportedInline(env.Type, path)
			if err != nil {
				return nil, err
			}
			if ok {
				return rawMarkdownInlineNode{Markdown: md}, nil
			}
		}
		if d.cfg.AllowUnsupportedNodes {
			return unsupportedInlineNode{Type: env.Type}, nil
		}
		return nil, newDecodeError(path, ErrKindUnsupportedInline, env.Type)
	}
}

func (d *decoder) decodeMediaInlineFromBlock(raw jsontext.Value, path string) (inlineNode, error) {
	env, err := d.decodeEnvelope(raw, path)
	if err != nil {
		return nil, err
	}
	if env.Type != "media" {
		if d.cfg.AllowUnsupportedNodes {
			return unsupportedInlineNode{Type: env.Type}, nil
		}
		return nil, newDecodeError(path, ErrKindInvalidStructure, "media container content must be media nodes")
	}
	u, a, err := d.mediaAttrsToURLAlt(env.Attrs, path+"/attrs")
	if err != nil {
		return nil, err
	}
	return mediaInlineNode{URL: u, Alt: a}, nil
}

func (d *decoder) mediaAttrsToURLAlt(attrs map[string]any, path string) (url, alt string, err error) {
	alt, _ = attrs["alt"].(string)
	mt, _ := attrs["type"].(string)
	switch mt {
	case "external":
		url, _ = attrs["url"].(string)
		if url == "" && d.cfg.StrictSchema {
			return "", "", newDecodeError(path, ErrKindInvalidAttr, "external media requires attrs.url")
		}
		return url, alt, nil
	case "file", "link":
		id, _ := attrs["id"].(string)
		collection, _ := attrs["collection"].(string)
		if id == "" && d.cfg.StrictSchema {
			return "", "", newDecodeError(path, ErrKindInvalidAttr, "file/link media requires attrs.id")
		}
		if alt == "" {
			alt = "media"
		}
		scheme := "atlassian-media://"
		if mt == "link" {
			scheme = "atlassian-media-link://"
		}
		if collection != "" {
			return scheme + collection + "/" + id, alt, nil
		}
		return scheme + id, alt, nil
	default:
		if d.cfg.StrictSchema {
			return "", "", newDecodeError(path, ErrKindInvalidAttr, "media attrs.type must be external, file, or link")
		}
		return "", alt, nil
	}
}

func (d *decoder) decodeInlineContent(raws []jsontext.Value, path string) ([]inlineNode, error) {
	out := make([]inlineNode, 0, len(raws))
	for i, raw := range raws {
		n, err := d.decodeInline(raw, fmt.Sprintf("%s/%d", path, i))
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, nil
}

func (d *decoder) decodeListItems(raws []jsontext.Value, path string) ([]listItemNode, error) {
	items := make([]listItemNode, 0, len(raws))
	for i, raw := range raws {
		n, err := d.decodeBlock(raw, fmt.Sprintf("%s/%d", path, i))
		if err != nil {
			return nil, err
		}
		item, ok := n.(listItemNode)
		if !ok {
			return nil, newDecodeError(fmt.Sprintf("%s/%d", path, i), ErrKindInvalidStructure, "expected listItem")
		}
		items = append(items, item)
	}
	return items, nil
}

func (d *decoder) validateMarks(marks []mark, path string) error {
	hasCode := false
	hasLink := false
	otherCount := 0
	for i, m := range marks {
		switch m.Type {
		case "strong", "em", "strike", "code", "link", "annotation", "underline", "subsup", "textColor", "backgroundColor":
			switch m.Type {
			case "code":
				hasCode = true
			case "link":
				hasLink = true
			default:
				otherCount++
			}
			if m.Type == "link" {
				if _, ok := m.Attrs["href"].(string); !ok {
					return newDecodeError(fmt.Sprintf("%s/%d", path, i), ErrKindInvalidMark, "link mark requires attrs.href")
				}
			}
			if m.Type == "subsup" && d.cfg.StrictSchema {
				st, _ := m.Attrs["type"].(string)
				if st != "sub" && st != "sup" {
					return newDecodeError(fmt.Sprintf("%s/%d", path, i), ErrKindInvalidMark, "subsup mark requires attrs.type=sub|sup")
				}
			}
		default:
			if d.cfg.AllowUnsupportedNodes {
				continue
			}
			return newDecodeError(fmt.Sprintf("%s/%d", path, i), ErrKindUnsupportedMark, m.Type)
		}
	}
	if d.cfg.StrictSchema && hasCode && (otherCount > 0 || (!hasLink && len(marks) > 1)) {
		return newDecodeError(path, ErrKindInvalidMarkCombo, "code mark can only be combined with link")
	}
	return nil
}

func (d *decoder) decodeCodeBlockText(raws []jsontext.Value, path string) (string, error) {
	out := ""
	for i, raw := range raws {
		env, err := d.decodeEnvelope(raw, fmt.Sprintf("%s/%d", path, i))
		if err != nil {
			return "", err
		}
		if env.Type != "text" {
			if d.cfg.AllowUnsupportedNodes {
				continue
			}
			return "", newDecodeError(fmt.Sprintf("%s/%d", path, i), ErrKindInvalidStructure, "codeBlock content must be text nodes")
		}
		if d.cfg.StrictSchema && len(env.Marks) > 0 {
			return "", newDecodeError(fmt.Sprintf("%s/%d", path, i), ErrKindInvalidStructure, "codeBlock text must not have marks")
		}
		out += env.Text
	}
	return out, nil
}
