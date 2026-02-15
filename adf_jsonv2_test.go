//go:build goexperiment.jsonv2

package adfmarkdown

import (
	"bytes"
	"errors"
	"testing"
)

func TestUnmarshalADF(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		opts     []Option
	}{
		{
			name: "paragraph with text",
			input: `{
				"version": 1,
				"type": "doc",
				"content": [
					{"type": "paragraph", "content": [{"type": "text", "text": "Hello world"}]}
				]
			}`,
			expected: "Hello world",
		},
		{
			name: "heading level 2",
			input: `{
				"version": 1,
				"type": "doc",
				"content": [
					{"type": "heading", "attrs": {"level": 2}, "content": [{"type": "text", "text": "Section"}]}
				]
			}`,
			expected: "## Section",
		},
		{
			name: "bullet list",
			input: `{
				"version": 1,
				"type": "doc",
				"content": [
					{"type": "bulletList", "content": [
						{"type": "listItem", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "One"}]}]},
						{"type": "listItem", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "Two"}]}]}
					]}
				]
			}`,
			expected: "- One\n- Two",
		},
		{
			name: "ordered list starting at 3",
			input: `{
				"version": 1,
				"type": "doc",
				"content": [
					{"type": "orderedList", "attrs": {"order": 3}, "content": [
						{"type": "listItem", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "One"}]}]},
						{"type": "listItem", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "Two"}]}]}
					]}
				]
			}`,
			expected: "3. One\n4. Two",
		},
		{
			name: "multiple blocks with spacing",
			input: `{
				"version": 1,
				"type": "doc",
				"content": [
					{"type": "heading", "attrs": {"level": 1}, "content": [{"type": "text", "text": "Title"}]},
					{"type": "paragraph", "content": [{"type": "text", "text": "Body"}]}
				]
			}`,
			expected: "# Title\n\nBody",
		},
		{
			name: "blockquote",
			input: `{
				"version": 1,
				"type": "doc",
				"content": [
					{"type": "blockquote", "content": [
						{"type": "paragraph", "content": [{"type": "text", "text": "Quoted"}]}
					]}
				]
			}`,
			expected: "> Quoted",
		},
		{
			name: "rule",
			input: `{
				"version": 1,
				"type": "doc",
				"content": [
					{"type": "paragraph", "content": [{"type": "text", "text": "Above"}]},
					{"type": "rule"},
					{"type": "paragraph", "content": [{"type": "text", "text": "Below"}]}
				]
			}`,
			expected: "Above\n\n---\n\nBelow",
		},
		{
			name: "code block with language",
			input: `{
				"version": 1,
				"type": "doc",
				"content": [
					{"type": "codeBlock", "attrs": {"language": "go"}, "content": [
						{"type": "text", "text": "fmt.Println(\"ok\")\n"}
					]}
				]
			}`,
			expected: "```go\nfmt.Println(\"ok\")\n```",
		},
		{
			name: "hard break default two spaces",
			input: `{
				"version": 1,
				"type": "doc",
				"content": [
					{"type": "paragraph", "content": [
						{"type": "text", "text": "a"},
						{"type": "hardBreak"},
						{"type": "text", "text": "b"}
					]}
				]
			}`,
			expected: "a  \nb",
		},
		{
			name: "hard break backslash style",
			input: `{
				"version": 1,
				"type": "doc",
				"content": [
					{"type": "paragraph", "content": [
						{"type": "text", "text": "a"},
						{"type": "hardBreak"},
						{"type": "text", "text": "b"}
					]}
				]
			}`,
			opts:     []Option{WithHardBreakStyle(HardBreakBackslash)},
			expected: "a\\\nb",
		},
		{
			name: "code fence tildes option",
			input: `{
				"version": 1,
				"type": "doc",
				"content": [
					{"type": "codeBlock", "content": [{"type": "text", "text": "x"}]}
				]
			}`,
			opts:     []Option{WithCodeFenceStyle(CodeFenceTildes)},
			expected: "~~~\nx\n~~~",
		},
		{
			name: "link mark with title",
			input: `{
				"version": 1,
				"type": "doc",
				"content": [
					{"type":"paragraph","content":[
						{"type":"text","text":"site","marks":[{"type":"link","attrs":{"href":"https://example.com","title":"Example"}}]}
					]}
				]
			}`,
			expected: `[site](https://example.com "Example")`,
		},
		{
			name:     "code span with backtick content",
			input:    "{\"version\":1,\"type\":\"doc\",\"content\":[{\"type\":\"paragraph\",\"content\":[{\"type\":\"text\",\"text\":\"a`b\",\"marks\":[{\"type\":\"code\"}]}]}]}",
			expected: "``a`b``",
		},
		{
			name:     "code block fence collision backticks",
			input:    "{\"version\":1,\"type\":\"doc\",\"content\":[{\"type\":\"codeBlock\",\"attrs\":{\"language\":\"md\"},\"content\":[{\"type\":\"text\",\"text\":\"```\\ncode\\n```\\n\"}]}]}",
			expected: "````md\n```\ncode\n```\n````",
		},
		{
			name:     "code block fence collision tildes",
			input:    "{\"version\":1,\"type\":\"doc\",\"content\":[{\"type\":\"codeBlock\",\"content\":[{\"type\":\"text\",\"text\":\"~~~\\nblock\\n~~~\\n\"}]}]}",
			opts:     []Option{WithCodeFenceStyle(CodeFenceTildes)},
			expected: "~~~~\n~~~\nblock\n~~~\n~~~~",
		},
		{
			name: "nested bullet list",
			input: `{
				"version": 1,
				"type": "doc",
				"content": [
					{"type":"bulletList","content":[
						{"type":"listItem","content":[
							{"type":"paragraph","content":[{"type":"text","text":"Parent"}]},
							{"type":"bulletList","content":[
								{"type":"listItem","content":[{"type":"paragraph","content":[{"type":"text","text":"Child"}]}]}
							]}
						]}
					]}
				]
			}`,
			expected: "- Parent\n\n    - Child",
		},
		{
			name: "blockquote inside list item",
			input: `{
				"version": 1,
				"type": "doc",
				"content": [
					{"type":"orderedList","content":[
						{"type":"listItem","content":[
							{"type":"blockquote","content":[{"type":"paragraph","content":[{"type":"text","text":"Quoted"}]}]}
						]}
					]}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "1. > Quoted",
		},
		{
			name: "nested blockquote",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"blockquote","content":[
						{"type":"blockquote","content":[
							{"type":"paragraph","content":[{"type":"text","text":"deep"}]}
						]}
					]}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "> > deep",
		},
		{
			name: "list item hard break indentation",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"bulletList","content":[
						{"type":"listItem","content":[
							{"type":"paragraph","content":[
								{"type":"text","text":"first"},
								{"type":"hardBreak"},
								{"type":"text","text":"second"}
							]}
						]}
					]}
				]
			}`,
			expected: "- first  \n  second",
		},
		{
			name: "panel node",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"panel","attrs":{"panelType":"warning"},"content":[
						{"type":"paragraph","content":[{"type":"text","text":"Careful"}]}
					]}
				]
			}`,
			expected: "> [!WARNING]\n> Careful",
		},
		{
			name: "expand node",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"expand","attrs":{"title":"More"},"content":[
						{"type":"paragraph","content":[{"type":"text","text":"Details text"}]}
					]}
				]
			}`,
			expected: "<details>\n<summary>More</summary>\n\nDetails text\n\n</details>",
		},
		{
			name: "media single external",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"mediaSingle","attrs":{"layout":"center"},"content":[
						{"type":"media","attrs":{"type":"external","url":"https://example.com/a.png","alt":"A"}}
					]}
				]
			}`,
			expected: "![A](https://example.com/a.png)",
		},
		{
			name: "media single with caption",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"mediaSingle","attrs":{"layout":"center"},"content":[
						{"type":"media","attrs":{"type":"external","url":"https://example.com/a.png","alt":"A"}},
						{"type":"caption","content":[{"type":"text","text":"Figure 1"}]}
					]}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "![A](https://example.com/a.png)\n_Figure 1_",
		},
		{
			name: "media single file type",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"mediaSingle","attrs":{"layout":"center"},"content":[
						{"type":"media","attrs":{"type":"file","id":"abc123","collection":"col1","alt":"Asset"}}
					]}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "![Asset](atlassian-media://col1/abc123)",
		},
		{
			name: "media single link type",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"mediaSingle","attrs":{"layout":"center"},"content":[
						{"type":"media","attrs":{"type":"link","id":"l-7","collection":"links"}}
					]}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "![media](atlassian-media-link://links/l-7)",
		},
		{
			name: "media group",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"mediaGroup","content":[
						{"type":"media","attrs":{"type":"external","url":"https://example.com/a.png","alt":"A"}},
						{"type":"media","attrs":{"type":"external","url":"https://example.com/b.png","alt":"B"}}
					]}
				]
			}`,
			expected: "![A](https://example.com/a.png)\n![B](https://example.com/b.png)",
		},
		{
			name: "table basic",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"table","attrs":{"isNumberColumnEnabled":false,"layout":"center","displayMode":"default"},"content":[
						{"type":"tableRow","content":[
							{"type":"tableHeader","attrs":{},"content":[{"type":"paragraph","content":[{"type":"text","text":"Name"}]}]},
							{"type":"tableHeader","attrs":{},"content":[{"type":"paragraph","content":[{"type":"text","text":"Age"}]}]}
						]},
						{"type":"tableRow","content":[
							{"type":"tableCell","attrs":{},"content":[{"type":"paragraph","content":[{"type":"text","text":"Alice"}]}]},
							{"type":"tableCell","attrs":{},"content":[{"type":"paragraph","content":[{"type":"text","text":"30"}]}]}
						]}
					]}
				]
			}`,
			expected: "| Name | Age |\n| --- | --- |\n| Alice | 30 |",
		},
		{
			name: "table without header row",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"table","attrs":{"isNumberColumnEnabled":false,"layout":"center","displayMode":"default"},"content":[
						{"type":"tableRow","content":[
							{"type":"tableCell","attrs":{},"content":[{"type":"paragraph","content":[{"type":"text","text":"A"}]}]},
							{"type":"tableCell","attrs":{},"content":[{"type":"paragraph","content":[{"type":"text","text":"B"}]}]}
						]},
						{"type":"tableRow","content":[
							{"type":"tableCell","attrs":{},"content":[{"type":"paragraph","content":[{"type":"text","text":"1"}]}]},
							{"type":"tableCell","attrs":{},"content":[{"type":"paragraph","content":[{"type":"text","text":"2"}]}]}
						]}
					]}
				]
			}`,
			expected: "| A | B |\n| --- | --- |\n| 1 | 2 |",
		},
		{
			name: "task list",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"taskList","content":[
						{"type":"taskItem","attrs":{"state":"TODO"},"content":[{"type":"text","text":"Ship parser"}]},
						{"type":"taskItem","attrs":{"state":"DONE"},"content":[{"type":"text","text":"Write tests"}]}
					]}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "- [ ] Ship parser\n- [x] Write tests",
		},
		{
			name: "block task item",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"blockTaskItem","attrs":{"state":"DONE","localId":"bt-1"},"content":[
						{"type":"paragraph","content":[{"type":"text","text":"Standalone task"}]}
					]}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "- [x] Standalone task",
		},
		{
			name: "decision list",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"decisionList","content":[
						{"type":"decisionItem","attrs":{"state":"DECIDED"},"content":[{"type":"text","text":"Use json/v2"}]}
					]}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "- [decision:DECIDED] Use json/v2",
		},
		{
			name: "block card",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"blockCard","attrs":{"url":"https://example.com/card"}}
				]
			}`,
			expected: "[https://example.com/card](https://example.com/card)",
		},
		{
			name: "embed card",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"embedCard","attrs":{"url":"https://example.com/embed"}}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "[https://example.com/embed](https://example.com/embed)",
		},
		{
			name: "extension default fallback",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"extension","attrs":{"extensionKey":"x.y.z"}}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "[extension:extension:x.y.z]",
		},
		{
			name: "layout section with columns",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"layoutSection","content":[
						{"type":"layoutColumn","content":[
							{"type":"paragraph","content":[{"type":"text","text":"Left"}]}
						]},
						{"type":"layoutColumn","content":[
							{"type":"paragraph","content":[{"type":"text","text":"Right"}]}
						]}
					]}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "[layout-section]\n\n[layout-column 1]\n\nLeft\n\n[layout-column 2]\n\nRight",
		},
		{
			name: "card url from attrs data",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"blockCard","attrs":{"data":{"link":{"href":"https://example.com/from-data"}}}}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "[https://example.com/from-data](https://example.com/from-data)",
		},
		{
			name: "card url from attrs data array",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"blockCard","attrs":{"data":{"items":[{"href":"https://example.com/from-array"}]}}}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "[https://example.com/from-array](https://example.com/from-array)",
		},
		{
			name: "bodied extension default fallback",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"bodiedExtension","attrs":{"extensionType":"macro"},"content":[
						{"type":"paragraph","content":[{"type":"text","text":"inside"}]}
					]}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "[extension:bodiedExtension:macro]\n\ninside",
		},
		{
			name: "inline extension default fallback",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"paragraph","content":[
						{"type":"text","text":"A "},
						{"type":"inlineExtension","attrs":{"extensionKey":"chip"}},
						{"type":"text","text":" B"}
					]}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "A [extension:inlineExtension:chip] B",
		},
		{
			name: "placeholder inline",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"paragraph","content":[{"type":"placeholder","attrs":{"text":"name"}}]}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "{{name}}",
		},
		{
			name: "sync block fallback",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"syncBlock","attrs":{"localId":"sb-1"}}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "[extension:syncBlock:sb-1]",
		},
		{
			name: "annotation mark renders span",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"paragraph","content":[
						{"type":"text","text":"hello","marks":[{"type":"annotation","attrs":{"id":"a1","annotationType":"inlineComment"}}]}
					]}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "<span data-annotation-id=\"a1\" data-annotation-type=\"inlineComment\">hello</span>",
		},
		{
			name: "underline mark",
			input: `{
				"version":1,"type":"doc","content":[
					{"type":"paragraph","content":[
						{"type":"text","text":"hello","marks":[{"type":"underline"}]}
					]}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "<u>hello</u>",
		},
		{
			name: "subsup mark",
			input: `{
				"version":1,"type":"doc","content":[
					{"type":"paragraph","content":[
						{"type":"text","text":"2","marks":[{"type":"subsup","attrs":{"type":"sup"}}]}
					]}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "<sup>2</sup>",
		},
		{
			name: "text color mark",
			input: `{
				"version":1,"type":"doc","content":[
					{"type":"paragraph","content":[
						{"type":"text","text":"red","marks":[{"type":"textColor","attrs":{"color":"#ff0000"}}]}
					]}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "<span style=\"color:#ff0000\">red</span>",
		},
		{
			name: "background color mark",
			input: `{
				"version":1,"type":"doc","content":[
					{"type":"paragraph","content":[
						{"type":"text","text":"hl","marks":[{"type":"backgroundColor","attrs":{"color":"#ffff00"}}]}
					]}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "<span style=\"background-color:#ffff00\">hl</span>",
		},
		{
			name: "inline nodes emoji mention date status card",
			input: `{
				"version":1,
				"type":"doc",
				"content":[
					{"type":"paragraph","content":[
						{"type":"emoji","attrs":{"text":"😀"}},
						{"type":"text","text":" "},
						{"type":"mention","attrs":{"id":"abc","text":"@Brad"}},
						{"type":"text","text":" "},
						{"type":"date","attrs":{"timestamp":"1582152559"}},
						{"type":"text","text":" "},
						{"type":"status","attrs":{"text":"In Progress","color":"yellow"}},
						{"type":"text","text":" "},
						{"type":"inlineCard","attrs":{"url":"https://atlassian.com"}}
					]}
				]
			}`,
			opts:     []Option{WithBuiltInSchemaValidation(false)},
			expected: "😀 @Brad [date:1582152559] [In Progress] [https://atlassian.com](https://atlassian.com)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalADF([]byte(tt.input), tt.opts...)
			if err != nil {
				t.Fatalf("UnmarshalADF failed: %v", err)
			}
			if string(got) != tt.expected {
				t.Fatalf("unexpected markdown\nwant:\n%s\n\ngot:\n%s", tt.expected, got)
			}
		})
	}
}

func TestStrictAndErrorPaths(t *testing.T) {
	t.Run("invalid version strict", func(t *testing.T) {
		_, err := UnmarshalADF([]byte(`{"version":2,"type":"doc","content":[]}`), WithBuiltInSchemaValidation(false))
		if err == nil {
			t.Fatal("expected error")
		}
		var de *Error
		if !errors.As(err, &de) {
			t.Fatalf("expected *Error, got %T", err)
		}
		if de.Path != "/" || de.Kind != ErrKindInvalidRoot {
			t.Fatalf("unexpected error metadata: %+v", de)
		}
	})

	t.Run("invalid version lenient", func(t *testing.T) {
		got, err := UnmarshalADF(
			[]byte(`{"version":2,"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"ok"}]}]}`),
			WithStrictSchema(false),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(got) != "ok" {
			t.Fatalf("unexpected output: %q", got)
		}
	})

	t.Run("empty text strict", func(t *testing.T) {
		_, err := UnmarshalADF([]byte(`{
			"version":1,"type":"doc","content":[
				{"type":"paragraph","content":[{"type":"text","text":""}]}
			]}`), WithBuiltInSchemaValidation(false))
		if err == nil {
			t.Fatal("expected error")
		}
		var de *Error
		if !errors.As(err, &de) {
			t.Fatalf("expected *Error, got %T", err)
		}
		if de.Path != "/content/0/content/0" || de.Kind != ErrKindInvalidText {
			t.Fatalf("unexpected error metadata: %+v", de)
		}
	})

	t.Run("bad link mark path", func(t *testing.T) {
		_, err := UnmarshalADF([]byte(`{
			"version":1,"type":"doc","content":[
				{"type":"paragraph","content":[{"type":"text","text":"x","marks":[{"type":"link","attrs":{}}]}]}
			]}`), WithBuiltInSchemaValidation(false))
		if err == nil {
			t.Fatal("expected error")
		}
		var de *Error
		if !errors.As(err, &de) {
			t.Fatalf("expected *Error, got %T", err)
		}
		if de.Path != "/content/0/content/0/marks/0" || de.Kind != ErrKindInvalidMark {
			t.Fatalf("unexpected error metadata: %+v", de)
		}
	})

	t.Run("invalid code mark combination", func(t *testing.T) {
		_, err := UnmarshalADF([]byte(`{
			"version":1,"type":"doc","content":[
				{"type":"paragraph","content":[{"type":"text","text":"x","marks":[{"type":"code"},{"type":"strong"}]}]}
			]}`), WithBuiltInSchemaValidation(false))
		if err == nil {
			t.Fatal("expected error")
		}
		var de *Error
		if !errors.As(err, &de) {
			t.Fatalf("expected *Error, got %T", err)
		}
		if de.Path != "/content/0/content/0/marks" || de.Kind != ErrKindInvalidMarkCombo {
			t.Fatalf("unexpected error metadata: %+v", de)
		}
	})

	t.Run("invalid subsup mark attrs", func(t *testing.T) {
		_, err := UnmarshalADF([]byte(`{
			"version":1,"type":"doc","content":[
				{"type":"paragraph","content":[{"type":"text","text":"x","marks":[{"type":"subsup","attrs":{"type":"bad"}}]}]}
			]}`), WithBuiltInSchemaValidation(false))
		if err == nil {
			t.Fatal("expected error")
		}
		var de *Error
		if !errors.As(err, &de) {
			t.Fatalf("expected *Error, got %T", err)
		}
		if de.Kind != ErrKindInvalidMark {
			t.Fatalf("unexpected error metadata: %+v", de)
		}
	})

	t.Run("schema validator hook", func(t *testing.T) {
		called := false
		_, err := UnmarshalADF(
			[]byte(`{"version":1,"type":"doc","content":[]}`),
			WithSchemaValidator(func([]byte) error {
				called = true
				return errors.New("schema failed")
			}),
		)
		if !called {
			t.Fatal("expected schema validator to be called")
		}
		if err == nil || err.Error() != "schema failed" {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("schema validator runs when strict false", func(t *testing.T) {
		called := false
		_, err := UnmarshalADF(
			[]byte(`{"version":2,"type":"doc","content":[]}`),
			WithStrictSchema(false),
			WithSchemaValidator(func([]byte) error {
				called = true
				return errors.New("validator invoked")
			}),
		)
		if !called {
			t.Fatal("expected schema validator to be called")
		}
		if err == nil || err.Error() != "validator invoked" {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("unsupported node handlers", func(t *testing.T) {
		got, err := UnmarshalADF(
			[]byte(`{
				"version":1,"type":"doc","content":[
					{"type":"mysteryBlock"},
					{"type":"paragraph","content":[{"type":"mysteryInline"}]}
				]
			}`),
			WithBuiltInSchemaValidation(false),
			WithUnsupportedBlockHandler(func(nodeType, path string) (string, bool, error) {
				if nodeType == "mysteryBlock" && path == "/content/0" {
					return "<block-hook>", true, nil
				}
				return "", false, nil
			}),
			WithUnsupportedInlineHandler(func(nodeType, path string) (string, bool, error) {
				if nodeType == "mysteryInline" && path == "/content/1/content/0" {
					return "<inline-hook>", true, nil
				}
				return "", false, nil
			}),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(got) != "<block-hook>\n\n<inline-hook>" {
			t.Fatalf("unexpected output: %q", got)
		}
	})

	t.Run("extension handlers", func(t *testing.T) {
		got, err := UnmarshalADF(
			[]byte(`{
				"version":1,"type":"doc","content":[
					{"type":"extension","attrs":{"extensionKey":"k1"}},
					{"type":"paragraph","content":[{"type":"inlineExtension","attrs":{"extensionKey":"k2"}}]}
				]
			}`),
			WithBuiltInSchemaValidation(false),
			WithExtensionBlockHandler(func(nodeType, key, path string) (string, bool, error) {
				if nodeType == "extension" && key == "k1" && path == "/content/0" {
					return "<ext-block>", true, nil
				}
				return "", false, nil
			}),
			WithExtensionInlineHandler(func(nodeType, key, path string) (string, bool, error) {
				if nodeType == "inlineExtension" && key == "k2" && path == "/content/1/content/0" {
					return "<ext-inline>", true, nil
				}
				return "", false, nil
			}),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(got) != "<ext-block>\n\n<ext-inline>" {
			t.Fatalf("unexpected output: %q", got)
		}
	})

	t.Run("built-in schema validator disabled", func(t *testing.T) {
		_, err := UnmarshalADF(
			[]byte(`{"version":1,"type":"doc","content":[{"type":"unknown-node"}]}`),
			WithBuiltInSchemaValidation(false),
		)
		if err == nil {
			t.Fatal("expected decode-layer unsupported node error")
		}
		var de *Error
		if !errors.As(err, &de) {
			t.Fatalf("expected *Error, got %T", err)
		}
		if de.Kind != ErrKindUnsupportedNode {
			t.Fatalf("unexpected error kind: %+v", de)
		}
	})
}

func TestValidateADFSchema(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		err := ValidateADFSchema([]byte(`{
			"version": 1,
			"type": "doc",
			"content": [{"type":"paragraph","content":[{"type":"text","text":"ok"}]}]
		}`))
		if err != nil {
			t.Fatalf("expected valid schema, got %v", err)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		err := ValidateADFSchema([]byte(`{"version":1,"type":"doc","content":[{"type":"text","text":"bad"}]}`))
		if err == nil {
			t.Fatal("expected schema validation error")
		}
		if got := err.Error(); got == "" || !contains(got, "ADF schema validation failed") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}

func TestUnmarshalADFTo(t *testing.T) {
	input := []byte(`{
		"version": 1,
		"type": "doc",
		"content": [
			{"type": "paragraph", "content": [{"type": "text", "text": "Hello"}]}
		]
	}`)
	var buf bytes.Buffer
	if err := UnmarshalADFTo(&buf, input); err != nil {
		t.Fatalf("UnmarshalADFTo failed: %v", err)
	}
	if buf.String() != "Hello" {
		t.Fatalf("unexpected output: %q", buf.String())
	}
}

func TestUnsupportedNodeOption(t *testing.T) {
	input := []byte(`{
		"version": 1,
		"type": "doc",
		"content": [
			{"type": "totallyUnknown", "content": []}
		]
	}`)

	if _, err := UnmarshalADF(input, WithBuiltInSchemaValidation(false)); err == nil {
		t.Fatal("expected error for unsupported node")
	}

	got, err := UnmarshalADF(input, WithAllowUnsupportedNodes(true), WithBuiltInSchemaValidation(false))
	if err != nil {
		t.Fatalf("expected success with unsupported-node fallback: %v", err)
	}
	if string(got) != "[unsupported node: totallyUnknown]" {
		t.Fatalf("unexpected fallback output: %q", got)
	}
}
