//go:build goexperiment.jsonv2

package adfmarkdown

import (
	"encoding/json/v2"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fixtureOptions struct {
	StrictSchema       *bool `json:"strict_schema"`
	BuiltInSchemaCheck *bool `json:"builtin_schema_check"`
	AllowUnsupported   *bool `json:"allow_unsupported"`
}

func TestConformanceFixtures(t *testing.T) {
	entries, err := os.ReadDir("testdata/fixtures")
	if err != nil {
		t.Fatalf("read fixtures: %v", err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") || strings.HasSuffix(e.Name(), ".opts.json") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".json")
		t.Run(name, func(t *testing.T) {
			inPath := filepath.Join("testdata/fixtures", name+".json")
			outPath := filepath.Join("testdata/fixtures", name+".md")
			optsPath := filepath.Join("testdata/fixtures", name+".opts.json")

			in, err := os.ReadFile(inPath)
			if err != nil {
				t.Fatalf("read input: %v", err)
			}
			want, err := os.ReadFile(outPath)
			if err != nil {
				t.Fatalf("read expected output: %v", err)
			}

			opts := []Option{}
			if optBytes, err := os.ReadFile(optsPath); err == nil {
				var fx fixtureOptions
				if err := json.Unmarshal(optBytes, &fx); err != nil {
					t.Fatalf("parse opts: %v", err)
				}
				if fx.StrictSchema != nil {
					opts = append(opts, WithStrictSchema(*fx.StrictSchema))
				}
				if fx.BuiltInSchemaCheck != nil {
					opts = append(opts, WithBuiltInSchemaValidation(*fx.BuiltInSchemaCheck))
				}
				if fx.AllowUnsupported != nil {
					opts = append(opts, WithAllowUnsupportedNodes(*fx.AllowUnsupported))
				}
			}

			got, err := UnmarshalADF(in, opts...)
			if err != nil {
				t.Fatalf("UnmarshalADF failed: %v", err)
			}
			if string(got) != strings.TrimRight(string(want), "\n") {
				t.Fatalf("unexpected markdown\nwant:\n%s\n\ngot:\n%s", want, got)
			}
		})
	}
}
