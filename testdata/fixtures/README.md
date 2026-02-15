# Conformance Fixtures

Each fixture is a pair:

- `<name>.json`: input ADF payload
- `<name>.md`: expected markdown output

Optional:

- `<name>.opts.json`: per-fixture decoder options

Example opts file:

```json
{
  "strict_schema": false,
  "builtin_schema_check": false,
  "allow_unsupported": true
}
```

Use these fixtures for sanitized real Jira/Confluence payloads.

