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

For every API-readback fixture, add a matching `<name>.meta.md` file with its
product/API representation, observation date, and a statement that it has been
structurally minimized and sanitized. Do not commit raw API responses.
