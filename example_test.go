//go:build goexperiment.jsonv2

package adfmarkdown_test

import (
	"fmt"
	"log"

	adfmarkdown "github.com/ajbeck/adf-to-markdown"
)

func ExampleUnmarshalADF() {
	input := []byte(`{"version":1,"type":"doc","content":[{"type":"heading","attrs":{"level":2},"content":[{"type":"text","text":"Overview"}]},{"type":"paragraph","content":[{"type":"text","text":"Hello from ADF"}]}]}`)

	md, err := adfmarkdown.UnmarshalADF(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(md))
	// Output:
	// ## Overview
	//
	// Hello from ADF
}
