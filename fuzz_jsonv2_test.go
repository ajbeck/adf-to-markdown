//go:build goexperiment.jsonv2

package adfmarkdown

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func FuzzUnmarshalADF(f *testing.F) {
	f.Add([]byte(`{"version":1,"type":"doc","content":[]}`))
	f.Add([]byte(`{"version":1,"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"hello"}]}]}`))

	entries, err := os.ReadDir("testdata/fixtures")
	if err == nil {
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") || strings.HasSuffix(e.Name(), ".opts.json") {
				continue
			}
			b, readErr := os.ReadFile(filepath.Join("testdata/fixtures", e.Name()))
			if readErr == nil {
				f.Add(b)
			}
		}
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 1<<20 {
			t.Skip()
		}
		_, _ = UnmarshalADF(
			data,
			WithBuiltInSchemaValidation(false),
			WithAllowUnsupportedNodes(true),
			WithStrictSchema(false),
		)
	})
}
