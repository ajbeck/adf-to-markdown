//go:build goexperiment.jsonv2

package adfmarkdown

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type emitter struct {
	buf bytes.Buffer
	cfg config
}

func newEmitter(cfg config) *emitter {
	return &emitter{cfg: cfg}
}

func (e *emitter) writeString(s string) {
	_, _ = e.buf.WriteString(s)
}

func (e *emitter) writeInlineNodes(nodes []inlineNode) {
	for _, n := range nodes {
		_ = n.renderInline(e)
	}
}

func (e *emitter) renderDocument(d document) ([]byte, error) {
	for i, n := range d.Content {
		if i > 0 {
			e.writeString("\n\n")
		}
		if err := n.renderBlock(e); err != nil {
			return nil, err
		}
	}
	return e.buf.Bytes(), nil
}

func (e *emitter) writeList(ordered bool, items []listItemNode, start int) error {
	for i, item := range items {
		prefix := "- "
		if ordered {
			prefix = strconv.Itoa(start+i) + ". "
		}
		if len(item.Content) == 0 {
			e.writeString(prefix)
			continue
		}

		// Render first block on the same line as marker; following blocks are
		// separated by an indented blank line.
		first := true
		for _, b := range item.Content {
			if !first {
				e.writeString("\n\n  ")
			} else {
				e.writeString(prefix)
			}
			inner := newEmitter(e.cfg)
			if err := b.renderBlock(inner); err != nil {
				return err
			}
			e.writeString(indentMultiline(inner.buf.String(), "  ", first))
			first = false
		}
		if i < len(items)-1 {
			e.writeString("\n")
		}
	}
	return nil
}

func (e *emitter) writeBlockquote(content []blockNode) error {
	inner := newEmitter(e.cfg)
	for i, b := range content {
		if i > 0 {
			inner.writeString("\n\n")
		}
		if err := b.renderBlock(inner); err != nil {
			return err
		}
	}
	e.writeString(prefixLines(inner.buf.String(), "> "))
	return nil
}

func (e *emitter) writeTable(rows []tableRowNode) error {
	if len(rows) == 0 {
		return nil
	}
	colCount := 0
	for _, r := range rows {
		if len(r.Cells) > colCount {
			colCount = len(r.Cells)
		}
	}
	if colCount == 0 {
		return nil
	}

	rowTexts := make([][]string, len(rows))
	for i, r := range rows {
		rowTexts[i] = make([]string, colCount)
		for j := 0; j < colCount; j++ {
			if j < len(r.Cells) {
				txt, err := e.renderTableCellText(r.Cells[j])
				if err != nil {
					return err
				}
				rowTexts[i][j] = txt
			}
		}
	}

	header := rowTexts[0]
	e.writeTableLine(header)
	sep := make([]string, colCount)
	for i := 0; i < colCount; i++ {
		sep[i] = "---"
	}
	e.writeString("\n")
	e.writeTableLine(sep)
	for i := 1; i < len(rowTexts); i++ {
		e.writeString("\n")
		e.writeTableLine(rowTexts[i])
	}
	return nil
}

func (e *emitter) writeTableLine(cells []string) {
	e.writeString("|")
	for _, c := range cells {
		e.writeString(" ")
		e.writeString(escapeTableCell(c))
		e.writeString(" |")
	}
}

func (e *emitter) renderTableCellText(c tableCellNode) (string, error) {
	if len(c.Content) == 0 {
		return "", nil
	}
	inner := newEmitter(e.cfg)
	for i, b := range c.Content {
		if i > 0 {
			inner.writeString("\n\n")
		}
		if err := b.renderBlock(inner); err != nil {
			return "", err
		}
	}
	s := inner.buf.String()
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\n\n", "<br><br>")
	s = strings.ReplaceAll(s, "\n", "<br>")
	return s, nil
}

func indentMultiline(s, indent string, firstLineAlreadyHasPrefix bool) string {
	var out bytes.Buffer
	lineStart := true
	for i := 0; i < len(s); i++ {
		if lineStart {
			if !firstLineAlreadyHasPrefix || i != 0 {
				_, _ = out.WriteString(indent)
			}
			lineStart = false
		}
		_ = out.WriteByte(s[i])
		if s[i] == '\n' {
			lineStart = true
		}
	}
	return out.String()
}

func prefixLines(s, prefix string) string {
	var out bytes.Buffer
	lineStart := true
	for i := 0; i < len(s); i++ {
		if lineStart {
			_, _ = out.WriteString(prefix)
			lineStart = false
		}
		_ = out.WriteByte(s[i])
		if s[i] == '\n' {
			lineStart = true
		}
	}
	if len(s) == 0 {
		_, _ = out.WriteString(prefix)
	}
	return out.String()
}

func applyMarks(text string, marks []mark) string {
	out := text
	for _, m := range marks {
		switch m.Type {
		case "strong":
			out = "**" + out + "**"
		case "em":
			out = "*" + out + "*"
		case "strike":
			out = "~~" + out + "~~"
		case "code":
			out = wrapCodeSpan(out)
		case "link":
			href, _ := m.Attrs["href"].(string)
			if href != "" {
				title, _ := m.Attrs["title"].(string)
				if title != "" {
					escapedTitle := strings.ReplaceAll(title, `"`, `\"`)
					out = "[" + out + "](" + href + ` "` + escapedTitle + `")`
				} else {
					out = "[" + out + "](" + href + ")"
				}
			}
		case "underline":
			out = "<u>" + out + "</u>"
		case "subsup":
			st, _ := m.Attrs["type"].(string)
			if st == "sub" {
				out = "<sub>" + out + "</sub>"
			} else if st == "sup" {
				out = "<sup>" + out + "</sup>"
			}
		case "textColor":
			color, _ := m.Attrs["color"].(string)
			if color != "" {
				out = `<span style="color:` + color + `">` + out + `</span>`
			}
		case "backgroundColor":
			color, _ := m.Attrs["color"].(string)
			if color != "" {
				out = `<span style="background-color:` + color + `">` + out + `</span>`
			}
		case "annotation":
			id, _ := m.Attrs["id"].(string)
			typ, _ := m.Attrs["annotationType"].(string)
			if id != "" || typ != "" {
				out = `<span data-annotation-id="` + id + `" data-annotation-type="` + typ + `">` + out + `</span>`
			}
		}
	}
	return out
}

func escapeTableCell(s string) string {
	return strings.ReplaceAll(s, "|", "\\|")
}

// escapeDelimiters escapes each character in delims that appears in s with a
// preceding backslash. This is used for round-trip custom syntax tokens where
// user-controlled text may contain delimiter characters.
func escapeDelimiters(s string, delims string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if strings.IndexByte(delims, s[i]) >= 0 {
			b.WriteByte('\\')
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func chooseCodeFence(text string, style CodeFenceStyle) string {
	fenceChar := '`'
	if style == CodeFenceTildes {
		fenceChar = '~'
	}
	n := maxRun(text, byte(fenceChar)) + 1
	if n < 3 {
		n = 3
	}
	return strings.Repeat(string(fenceChar), n)
}

func maxRun(s string, ch byte) int {
	maxN := 0
	cur := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ch {
			cur++
			if cur > maxN {
				maxN = cur
			}
		} else {
			cur = 0
		}
	}
	return maxN
}

func wrapCodeSpan(s string) string {
	n := maxRun(s, '`') + 1
	if n < 1 {
		n = 1
	}
	fence := strings.Repeat("`", n)
	if strings.HasPrefix(s, "`") || strings.HasSuffix(s, "`") {
		return fence + " " + s + " " + fence
	}
	return fence + s + fence
}

func mustInt(v any) (int, error) {
	switch x := v.(type) {
	case float64:
		return int(x), nil
	case int:
		return x, nil
	default:
		return 0, fmt.Errorf("not an integer: %T", v)
	}
}
