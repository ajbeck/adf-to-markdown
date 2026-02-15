//go:build goexperiment.jsonv2

// Package adfmarkdown converts Atlassian Document Format (ADF) JSON to Markdown.
//
// # Build Requirements
//
// This package requires Go 1.25+ with jsonv2 enabled:
//
//	GOEXPERIMENT=jsonv2 go build ./...
//	GOEXPERIMENT=jsonv2 go test ./...
//
// # Basic Usage
//
//	md, err := adfmarkdown.UnmarshalADF(adfJSON)
//
// # Options
//
// Use functional options to control schema validation and rendering:
//
//	adfmarkdown.UnmarshalADF(adfJSON,
//		adfmarkdown.WithStrictSchema(true),
//		adfmarkdown.WithCodeFenceStyle(adfmarkdown.CodeFenceBackticks),
//	)
//
// # Errors
//
// Decode failures return *Error with a stable ErrorKind and ADF path metadata.
package adfmarkdown
