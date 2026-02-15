//go:build goexperiment.jsonv2

package adfmarkdown

import "fmt"

type blockNode interface {
	renderBlock(*emitter) error
}

type inlineNode interface {
	renderInline(*emitter) error
}

type document struct {
	Version int
	Content []blockNode
}

type paragraphNode struct {
	Content []inlineNode
}

type headingNode struct {
	Level   int
	Content []inlineNode
}

type blockquoteNode struct {
	Content []blockNode
}

type codeBlockNode struct {
	Language string
	Text     string
}

type ruleNode struct{}

type panelNode struct {
	PanelType string
	Content   []blockNode
}

type expandNode struct {
	Title   string
	Content []blockNode
}

type mediaSingleNode struct {
	Content []inlineNode
	Caption []inlineNode
}

type mediaGroupNode struct {
	Content []inlineNode
}

type taskListNode struct {
	Items []taskItemNode
}

type taskItemNode struct {
	State   string
	Content []inlineNode
}

type blockTaskItemNode struct {
	State   string
	Content []blockNode
}

type decisionListNode struct {
	Items []decisionItemNode
}

type decisionItemNode struct {
	State   string
	Content []inlineNode
}

type cardBlockNode struct {
	URL string
}

type extensionBlockNode struct {
	NodeType string
	Key      string
	Path     string
	Content  []blockNode
}

type layoutSectionNode struct {
	Columns []layoutColumnNode
}

type layoutColumnNode struct {
	Content []blockNode
}

type tableNode struct {
	Rows []tableRowNode
}

type tableRowNode struct {
	Cells []tableCellNode
}

type tableCellNode struct {
	IsHeader bool
	Content  []blockNode
}

type bulletListNode struct {
	Items []listItemNode
}

type orderedListNode struct {
	Start int
	Items []listItemNode
}

type listItemNode struct {
	Content []blockNode
}

type unsupportedBlockNode struct {
	Type string
}

type rawMarkdownBlockNode struct {
	Markdown string
}

type textNode struct {
	Text  string
	Marks []mark
}

type hardBreakNode struct{}

type unsupportedInlineNode struct {
	Type string
}

type rawMarkdownInlineNode struct {
	Markdown string
}

type mark struct {
	Type  string         `json:"type"`
	Attrs map[string]any `json:"attrs"`
}

func (n paragraphNode) renderBlock(e *emitter) error {
	e.writeInlineNodes(n.Content)
	return nil
}

func (n headingNode) renderBlock(e *emitter) error {
	if n.Level < 1 || n.Level > 6 {
		return fmt.Errorf("invalid heading level: %d", n.Level)
	}
	for range n.Level {
		e.writeString("#")
	}
	e.writeString(" ")
	e.writeInlineNodes(n.Content)
	return nil
}

func (n blockquoteNode) renderBlock(e *emitter) error {
	return e.writeBlockquote(n.Content)
}

func (n codeBlockNode) renderBlock(e *emitter) error {
	fence := chooseCodeFence(n.Text, e.cfg.CodeFenceStyle)
	e.writeString(fence)
	if n.Language != "" {
		e.writeString(n.Language)
	}
	e.writeString("\n")
	e.writeString(n.Text)
	if n.Text != "" && n.Text[len(n.Text)-1] != '\n' {
		e.writeString("\n")
	}
	e.writeString(fence)
	return nil
}

func (n ruleNode) renderBlock(e *emitter) error {
	e.writeString("---")
	return nil
}

func (n panelNode) renderBlock(e *emitter) error {
	label := n.PanelType
	if label == "" {
		label = "info"
	}
	e.writeString("> [!")
	e.writeString(toUpperASCII(label))
	e.writeString("]\n")
	return e.writeBlockquote(n.Content)
}

func (n expandNode) renderBlock(e *emitter) error {
	e.writeString("<details>\n")
	e.writeString("<summary>")
	if n.Title != "" {
		e.writeString(n.Title)
	} else {
		e.writeString("Details")
	}
	e.writeString("</summary>\n\n")
	for i, b := range n.Content {
		if i > 0 {
			e.writeString("\n\n")
		}
		if err := b.renderBlock(e); err != nil {
			return err
		}
	}
	e.writeString("\n\n</details>")
	return nil
}

func (n mediaSingleNode) renderBlock(e *emitter) error {
	e.writeInlineNodes(n.Content)
	if len(n.Caption) > 0 {
		e.writeString("\n")
		e.writeString("_")
		e.writeInlineNodes(n.Caption)
		e.writeString("_")
	}
	return nil
}

func (n mediaGroupNode) renderBlock(e *emitter) error {
	for i, in := range n.Content {
		if i > 0 {
			e.writeString("\n")
		}
		if err := in.renderInline(e); err != nil {
			return err
		}
	}
	return nil
}

func (n taskListNode) renderBlock(e *emitter) error {
	for i, it := range n.Items {
		if i > 0 {
			e.writeString("\n")
		}
		if err := it.renderBlock(e); err != nil {
			return err
		}
	}
	return nil
}

func (n taskItemNode) renderBlock(e *emitter) error {
	if n.State == "DONE" {
		e.writeString("- [x] ")
	} else {
		e.writeString("- [ ] ")
	}
	e.writeInlineNodes(n.Content)
	return nil
}

func (n blockTaskItemNode) renderBlock(e *emitter) error {
	prefix := "- [ ] "
	if n.State == "DONE" {
		prefix = "- [x] "
	}
	if len(n.Content) == 0 {
		e.writeString(prefix)
		return nil
	}
	inner := newEmitter(e.cfg)
	for i, b := range n.Content {
		if i > 0 {
			inner.writeString("\n\n")
		}
		if err := b.renderBlock(inner); err != nil {
			return err
		}
	}
	e.writeString(prefix)
	e.writeString(indentMultiline(inner.buf.String(), "  ", true))
	return nil
}

func (n decisionListNode) renderBlock(e *emitter) error {
	for i, it := range n.Items {
		if i > 0 {
			e.writeString("\n")
		}
		if err := it.renderBlock(e); err != nil {
			return err
		}
	}
	return nil
}

func (n decisionItemNode) renderBlock(e *emitter) error {
	e.writeString("- [decision")
	if n.State != "" {
		e.writeString(":")
		e.writeString(n.State)
	}
	e.writeString("] ")
	e.writeInlineNodes(n.Content)
	return nil
}

func (n cardBlockNode) renderBlock(e *emitter) error {
	if n.URL == "" {
		e.writeString("[card]")
		return nil
	}
	e.writeString("[")
	e.writeString(n.URL)
	e.writeString("](")
	e.writeString(n.URL)
	e.writeString(")")
	return nil
}

func (n tableNode) renderBlock(e *emitter) error {
	return e.writeTable(n.Rows)
}

func (n bulletListNode) renderBlock(e *emitter) error {
	return e.writeList(false, n.Items, 1)
}

func (n orderedListNode) renderBlock(e *emitter) error {
	start := n.Start
	if start < 1 {
		start = 1
	}
	return e.writeList(true, n.Items, start)
}

func (n unsupportedBlockNode) renderBlock(e *emitter) error {
	e.writeString("[unsupported node: ")
	e.writeString(n.Type)
	e.writeString("]")
	return nil
}

func (n rawMarkdownBlockNode) renderBlock(e *emitter) error {
	e.writeString(n.Markdown)
	return nil
}

func (n listItemNode) renderBlock(e *emitter) error {
	for i, b := range n.Content {
		if i > 0 {
			e.writeString("\n\n")
		}
		if err := b.renderBlock(e); err != nil {
			return err
		}
	}
	return nil
}

func (n textNode) renderInline(e *emitter) error {
	e.writeString(applyMarks(n.Text, n.Marks))
	return nil
}

func (n hardBreakNode) renderInline(e *emitter) error {
	if e.cfg.HardBreakStyle == HardBreakBackslash {
		e.writeString("\\\n")
	} else {
		e.writeString("  \n")
	}
	return nil
}

type emojiInlineNode struct {
	Text string
}

func (n emojiInlineNode) renderInline(e *emitter) error {
	e.writeString(n.Text)
	return nil
}

type mentionInlineNode struct {
	Text string
	ID   string
}

func (n mentionInlineNode) renderInline(e *emitter) error {
	if n.Text != "" {
		e.writeString(n.Text)
		return nil
	}
	if n.ID != "" {
		e.writeString("@")
		e.writeString(n.ID)
		return nil
	}
	e.writeString("@mention")
	return nil
}

type dateInlineNode struct {
	Timestamp string
}

func (n dateInlineNode) renderInline(e *emitter) error {
	e.writeString("[date:")
	e.writeString(n.Timestamp)
	e.writeString("]")
	return nil
}

type statusInlineNode struct {
	Text string
}

func (n statusInlineNode) renderInline(e *emitter) error {
	e.writeString("[")
	e.writeString(n.Text)
	e.writeString("]")
	return nil
}

type inlineCardNode struct {
	URL string
}

func (n inlineCardNode) renderInline(e *emitter) error {
	if n.URL == "" {
		e.writeString("[inline-card]")
		return nil
	}
	e.writeString("[")
	e.writeString(n.URL)
	e.writeString("](")
	e.writeString(n.URL)
	e.writeString(")")
	return nil
}

type mediaInlineNode struct {
	URL string
	Alt string
}

type placeholderInlineNode struct {
	Text string
}

type extensionInlineNode struct {
	NodeType string
	Key      string
	Path     string
}

func (n mediaInlineNode) renderInline(e *emitter) error {
	if n.URL != "" {
		e.writeString("![")
		e.writeString(n.Alt)
		e.writeString("](")
		e.writeString(n.URL)
		e.writeString(")")
		return nil
	}
	e.writeString("[media]")
	return nil
}

func (n placeholderInlineNode) renderInline(e *emitter) error {
	e.writeString("{{")
	e.writeString(n.Text)
	e.writeString("}}")
	return nil
}

func (n extensionInlineNode) renderInline(e *emitter) error {
	if e.cfg.ExtensionInline != nil {
		md, ok, err := e.cfg.ExtensionInline(n.NodeType, n.Key, n.Path)
		if err != nil {
			return err
		}
		if ok {
			e.writeString(md)
			return nil
		}
	}
	e.writeString("[extension:")
	e.writeString(n.NodeType)
	if n.Key != "" {
		e.writeString(":")
		e.writeString(n.Key)
	}
	e.writeString("]")
	return nil
}

func (n unsupportedInlineNode) renderInline(e *emitter) error {
	e.writeString("[unsupported inline: ")
	e.writeString(n.Type)
	e.writeString("]")
	return nil
}

func (n rawMarkdownInlineNode) renderInline(e *emitter) error {
	e.writeString(n.Markdown)
	return nil
}

func (n extensionBlockNode) renderBlock(e *emitter) error {
	if e.cfg.ExtensionBlock != nil {
		md, ok, err := e.cfg.ExtensionBlock(n.NodeType, n.Key, n.Path)
		if err != nil {
			return err
		}
		if ok {
			e.writeString(md)
			return nil
		}
	}
	e.writeString("[extension:")
	e.writeString(n.NodeType)
	if n.Key != "" {
		e.writeString(":")
		e.writeString(n.Key)
	}
	e.writeString("]")
	if len(n.Content) == 0 {
		return nil
	}
	for _, b := range n.Content {
		e.writeString("\n\n")
		if err := b.renderBlock(e); err != nil {
			return err
		}
	}
	return nil
}

func (n layoutSectionNode) renderBlock(e *emitter) error {
	e.writeString("[layout-section]")
	for i, c := range n.Columns {
		e.writeString("\n\n")
		e.writeString("[layout-column ")
		e.writeString(fmt.Sprintf("%d", i+1))
		e.writeString("]")
		for _, b := range c.Content {
			e.writeString("\n\n")
			if err := b.renderBlock(e); err != nil {
				return err
			}
		}
	}
	return nil
}

func toUpperASCII(s string) string {
	b := []byte(s)
	for i := range b {
		if b[i] >= 'a' && b[i] <= 'z' {
			b[i] -= 'a' - 'A'
		}
	}
	return string(b)
}
