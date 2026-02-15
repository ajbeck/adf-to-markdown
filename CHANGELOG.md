# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog and this project follows Semantic Versioning.

## [1.0.0] - TBD

### Added
- `UnmarshalADF(data []byte, opts ...Option) ([]byte, error)` for ADF JSON to Markdown conversion.
- `UnmarshalADFTo(w io.Writer, data []byte, opts ...Option) error` for streaming write targets.
- Schema validation controls via `WithStrictSchema`, `WithBuiltInSchemaValidation`, and `WithSchemaValidator`.
- Rendering controls via `WithHardBreakStyle` and `WithCodeFenceStyle`.
- Unsupported/extension hooks via `WithUnsupportedBlockHandler`, `WithUnsupportedInlineHandler`, `WithExtensionBlockHandler`, and `WithExtensionInlineHandler`.
- Typed decode errors with `*Error`, `ErrorKind`, path metadata, and wrapped causes.
- Built-in embedded ADF schema validation using `jsonschema-go`.
- Conformance fixture tests, fuzz target, and benchmarks.

### Behavior
- Requires Go `1.25+` and `GOEXPERIMENT=jsonv2`.
- Strict mode enforces ADF structural checks and mark constraints.
- Custom schema validator runs when provided, independent of strict mode.

### Limitations
- API intentionally depends on experimental `encoding/json/v2`.
- Extension nodes use fallback rendering unless handled by consumer-provided hooks.
