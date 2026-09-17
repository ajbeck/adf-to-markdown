//go:build goexperiment.jsonv2

package adfmarkdown

import "testing"

func TestPersistedAPICompatibility(t *testing.T) {
	paragraph := func(inline string) string {
		return `{"type":"paragraph","content":[` + inline + `]}`
	}
	tests := []struct {
		name, block, want string
	}{
		{"annotated inline card", paragraph(`{"type":"inlineCard","attrs":{"url":"https://example.com"},"marks":[{"type":"annotation","attrs":{"id":"a","annotationType":"inlineComment"}}]}`), "[card:https://example.com]"},
		{"date with url", paragraph(`{"type":"date","attrs":{"timestamp":"0","url":"https://example.com/date"}}`), "[date:0]"},
		{"emoji with url", paragraph(`{"type":"emoji","attrs":{"shortName":":smile:","url":"https://example.com/emoji"}}`), ":smile:"},
		{"inline media with url", paragraph(`{"type":"mediaInline","attrs":{"id":"a","collection":"b","url":"https://example.com/media"}}`), "![media](atlassian-media://b/a)"},
		{"heading without attrs", `{"type":"heading","content":[{"type":"text","text":"Title"}]}`, "# Title"},
		{"heading without level", `{"type":"heading","attrs":{"localId":"h"},"content":[{"type":"text","text":"Title"}]}`, "# Title"},
		{"heading with empty attrs", `{"type":"heading","attrs":{},"content":[{"type":"text","text":"Title"}]}`, "# Title"},
		{"panel without attrs", `{"type":"panel","content":[{"type":"paragraph","content":[{"type":"text","text":"Body"}]}]}`, "> [!INFO]\n> Body"},
		{"panel without panelType", `{"type":"panel","attrs":{"localId":"p"},"content":[{"type":"paragraph","content":[{"type":"text","text":"Body"}]}]}`, "> [!INFO]\n> Body"},
		{"panel with empty attrs", `{"type":"panel","attrs":{},"content":[{"type":"paragraph","content":[{"type":"text","text":"Body"}]}]}`, "> [!INFO]\n> Body"},
		{"code block with marked text", `{"type":"codeBlock","attrs":{"language":"go"},"content":[{"type":"text","text":"a < b\n","marks":[{"type":"strong"},{"type":"link","attrs":{"href":"https://example.com"}}]},{"type":"text","text":"next","marks":[{"type":"code"}]}]}`, "```go\na < b\nnext\n```"},
		{"root code block with breakout and marked text", `{"type":"codeBlock","marks":[{"type":"breakout","attrs":{"mode":"wide"}}],"content":[{"type":"text","text":"code","marks":[{"type":"em"}]}]}`, "```\ncode\n```"},
		{"nested code block with marked text", `{"type":"panel","content":[{"type":"codeBlock","content":[{"type":"text","text":"code","marks":[{"type":"strong"}]}]}]}`, "> [!INFO]\n> ```\n> code\n> ```"},
		{"status with hex color", paragraph(`{"type":"status","attrs":{"text":"Ready","color":"#aB12Ef"}}`), "[status:Ready|#aB12Ef]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := []byte(`{"version":1,"type":"doc","content":[` + tt.block + `]}`)
			if err := ValidateADFSchema(data); err != nil {
				t.Fatalf("schema rejected compatible ADF: %v", err)
			}
			for _, builtIn := range []bool{true, false} {
				got, err := UnmarshalADF(data, WithBuiltInSchemaValidation(builtIn))
				if err != nil {
					t.Fatalf("built-in validation %t: %v", builtIn, err)
				}
				if string(got) != tt.want {
					t.Fatalf("built-in validation %t: got %q, want %q", builtIn, got, tt.want)
				}
			}
		})
	}
}

func TestPersistedAPISchemaRejectsMalformedNodes(t *testing.T) {
	paragraph := func(inline string) string {
		return `{"type":"paragraph","content":[` + inline + `]}`
	}
	tests := []struct{ name, block string }{
		{"unknown inline card mark", paragraph(`{"type":"inlineCard","attrs":{"url":"https://example.com"},"marks":[{"type":"unknown"}]}`)},
		{"annotation missing id", paragraph(`{"type":"inlineCard","attrs":{"url":"https://example.com"},"marks":[{"type":"annotation","attrs":{"annotationType":"inlineComment"}}]}`)},
		{"date url must be string", paragraph(`{"type":"date","attrs":{"timestamp":"0","url":42}}`)},
		{"date still requires timestamp", paragraph(`{"type":"date","attrs":{"url":"https://example.com"}}`)},
		{"emoji url must be string", paragraph(`{"type":"emoji","attrs":{"shortName":":smile:","url":42}}`)},
		{"emoji still requires shortName", paragraph(`{"type":"emoji","attrs":{"url":"https://example.com"}}`)},
		{"media url must be string", paragraph(`{"type":"mediaInline","attrs":{"id":"a","collection":"b","url":42}}`)},
		{"media still requires id", paragraph(`{"type":"mediaInline","attrs":{"collection":"b","url":"https://example.com"}}`)},
		{"media still requires collection", paragraph(`{"type":"mediaInline","attrs":{"id":"a","url":"https://example.com"}}`)},
		{"heading level out of range", `{"type":"heading","attrs":{"level":7}}`},
		{"heading level is null", `{"type":"heading","attrs":{"level":null}}`},
		{"heading level wrong type", `{"type":"heading","attrs":{"level":"1"}}`},
		{"heading unknown attribute", `{"type":"heading","attrs":{"unexpected":true}}`},
		{"panelType invalid", `{"type":"panel","attrs":{"panelType":"unknown"},"content":[{"type":"paragraph"}]}`},
		{"panelType is null", `{"type":"panel","attrs":{"panelType":null},"content":[{"type":"paragraph"}]}`},
		{"panel still requires content", `{"type":"panel"}`},
		{"code block with nontext content", `{"type":"codeBlock","content":[{"type":"hardBreak"}]}`},
		{"code block with empty text", `{"type":"codeBlock","content":[{"type":"text","text":"","marks":[{"type":"strong"}]}]}`},
		{"code block with unknown mark", `{"type":"codeBlock","content":[{"type":"text","text":"code","marks":[{"type":"unknown"}]}]}`},
		{"code block with malformed link", `{"type":"codeBlock","content":[{"type":"text","text":"code","marks":[{"type":"link","attrs":{}}]}]}`},
		{"status short hex", paragraph(`{"type":"status","attrs":{"text":"Ready","color":"#abc"}}`)},
		{"status invalid hex", paragraph(`{"type":"status","attrs":{"text":"Ready","color":"#12345g"}}`)},
		{"status unknown name", paragraph(`{"type":"status","attrs":{"text":"Ready","color":"orange"}}`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := []byte(`{"version":1,"type":"doc","content":[` + tt.block + `]}`)
			if err := ValidateADFSchema(data); err == nil {
				t.Fatal("schema accepted malformed ADF")
			}
			if _, err := UnmarshalADF(data); err == nil {
				t.Fatal("default conversion accepted malformed ADF")
			}
		})
	}
}
