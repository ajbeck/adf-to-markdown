// Command generate-adf-schema normalizes Atlassian's generated draft-04
// schema and applies the project's persisted-API compatibility overlay.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const draft07Schema = "http://json-schema.org/draft-07/schema#"

type operation struct {
	Op    string          `json:"op"`
	Path  string          `json:"path"`
	Value json.RawMessage `json:"value"`
}

func main() {
	source := flag.String("source", "", "upstream Atlassian full.json schema")
	patch := flag.String("patch", "", "RFC 6902 compatibility patch")
	output := flag.String("output", "", "generated persisted-API schema")
	flag.Parse()
	if *source == "" || *patch == "" || *output == "" {
		fail(errors.New("source, patch, and output are required"))
	}

	document, err := readJSON(*source)
	if err != nil {
		fail(err)
	}
	root, ok := document.(map[string]any)
	if !ok {
		fail(fmt.Errorf("%s: schema root must be an object", *source))
	}
	root["$schema"] = draft07Schema

	patchDocument, err := os.ReadFile(*patch)
	if err != nil {
		fail(fmt.Errorf("read patch: %w", err))
	}
	var operations []operation
	if err := json.Unmarshal(patchDocument, &operations); err != nil {
		fail(fmt.Errorf("parse patch: %w", err))
	}
	for _, op := range operations {
		if err := apply(document, op); err != nil {
			fail(err)
		}
	}

	encoded, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		fail(fmt.Errorf("encode generated schema: %w", err))
	}
	encoded = append(encoded, '\n')
	if err := os.WriteFile(*output, encoded, 0o644); err != nil {
		fail(fmt.Errorf("write generated schema: %w", err))
	}
}

func readJSON(path string) (any, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var document any
	if err := json.Unmarshal(b, &document); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return document, nil
}

func apply(document any, op operation) error {
	if op.Op != "add" && op.Op != "replace" {
		return fmt.Errorf("unsupported RFC 6902 operation %q", op.Op)
	}
	if len(op.Value) == 0 {
		return fmt.Errorf("%s %q: value is required", op.Op, op.Path)
	}
	parts, err := pointerParts(op.Path)
	if err != nil {
		return err
	}
	if len(parts) == 0 {
		return fmt.Errorf("%s %q: replacing the root is not supported", op.Op, op.Path)
	}

	parent, err := resolve(document, parts[:len(parts)-1])
	if err != nil {
		return fmt.Errorf("%s %q: %w", op.Op, op.Path, err)
	}
	var value any
	if err := json.Unmarshal(op.Value, &value); err != nil {
		return fmt.Errorf("%s %q: parse value: %w", op.Op, op.Path, err)
	}
	last := parts[len(parts)-1]
	switch p := parent.(type) {
	case map[string]any:
		if op.Op == "replace" {
			if _, ok := p[last]; !ok {
				return fmt.Errorf("target does not exist")
			}
		}
		p[last] = value
		return nil
	case []any:
		if last == "-" && op.Op == "add" {
			return errors.New("array append is not supported by this generator")
		}
		index, err := strconv.Atoi(last)
		if err != nil || index < 0 || index >= len(p) {
			return fmt.Errorf("invalid array index %q", last)
		}
		if op.Op == "add" {
			return errors.New("array insertion is not supported by this generator")
		}
		p[index] = value
		return nil
	default:
		return fmt.Errorf("target parent is %T, not object or array", parent)
	}
}

func pointerParts(path string) ([]string, error) {
	if path == "" {
		return nil, nil
	}
	if !strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("invalid JSON pointer %q", path)
	}
	parts := strings.Split(path[1:], "/")
	for i, part := range parts {
		parts[i] = strings.ReplaceAll(strings.ReplaceAll(part, "~1", "/"), "~0", "~")
	}
	return parts, nil
}

func resolve(document any, parts []string) (any, error) {
	current := document
	for _, part := range parts {
		switch node := current.(type) {
		case map[string]any:
			var ok bool
			current, ok = node[part]
			if !ok {
				return nil, fmt.Errorf("pointer segment %q does not exist", part)
			}
		case []any:
			index, err := strconv.Atoi(part)
			if err != nil || index < 0 || index >= len(node) {
				return nil, fmt.Errorf("invalid array index %q", part)
			}
			current = node[index]
		default:
			return nil, fmt.Errorf("pointer segment %q has non-container parent %T", part, current)
		}
	}
	return current, nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "generate-adf-schema:", err)
	os.Exit(1)
}
