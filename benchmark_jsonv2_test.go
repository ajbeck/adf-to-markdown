//go:build goexperiment.jsonv2

package adfmarkdown

import (
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkUnmarshalADF(b *testing.B) {
	cases := []struct {
		name string
		file string
	}{
		{name: "small_paragraph", file: "paragraph_basic.json"},
		{name: "medium_comment", file: "jira_comment_mixed.json"},
		{name: "large_epic", file: "jira_epic_description.json"},
	}

	for _, tc := range cases {
		data, err := os.ReadFile(filepath.Join("testdata/fixtures", tc.file))
		if err != nil {
			b.Fatalf("read fixture %s: %v", tc.file, err)
		}
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, err := UnmarshalADF(data, WithBuiltInSchemaValidation(false))
				if err != nil {
					b.Fatalf("UnmarshalADF failed: %v", err)
				}
			}
		})
	}
}
