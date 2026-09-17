# ADF Schema Sources

`adf-full-schema-57.5.0.json` is the unmodified `dist/json-schema/v1/full.json`
file from [`@atlaskit/adf-schema@57.5.0`](https://unpkg.com/@atlaskit/adf-schema@57.5.0/dist/json-schema/v1/full.json)
(Apache-2.0). Keep the upstream file unchanged so that future updates can be
compared independently of this project's compatibility policy.

## Upstream changes

A structural JSON comparison of 57.0.0 and 57.5.0 finds one change:
`status.attrs.color` accepts six-digit hex colors (`#RRGGBB`) in addition to
`neutral`, `purple`, `blue`, `red`, `yellow`, and `green`. The existing renderer
already preserves the color in `[status:text|color]`.

Compared with the previously embedded 56.1.3 schema, this is also the only
validation change. That older committed file already used a draft-07 declaration;
the upstream 57.5.0 file still declares draft-04. None of the compatibility cases
in [issue #7](https://github.com/ajbeck/adf-to-markdown/issues/7) are fixed upstream.

## Persisted-API compatibility

`adf-persisted-api.patch.json` is a narrowly scoped RFC 6902 overlay for
Confluence/Jira API compatibility. It includes the existing API-readback fixtures
and the additional shapes reported in issue #7. Issue-derived regression cases
are synthetic reproductions of the report, not newly captured API responses.
These allowances are a library compatibility policy, not claims that the upstream
JSON schema accepts these shapes.

| Shape | Validation and Markdown behavior | Evidence |
|---|---|---|
| `inlineCard.marks` | Accept annotation marks with their required attributes; omit annotation metadata from `[card:url]`. Other mark types remain invalid. | Existing Confluence readback fixture; #7 |
| `mediaInline.attrs.__fileName` | Accept a string; omit service metadata. | Existing Confluence readback fixture |
| `taskItem.attrs.localId: null` | Accept null or a string; omit the session identifier. | Existing Confluence readback fixture |
| `date`, `emoji`, `mediaInline` with `attrs.url` | Accept a string without relaxing the other required attributes. Omit the extra URL; media continues to use its `id` and `collection`. | #7 |
| `heading` without `attrs.level` | Allow missing/empty `attrs`; render level 1. Explicit invalid levels remain invalid. | #7; upstream runtime default |
| `panel` without `attrs.panelType` | Allow missing/empty `attrs`; render an INFO alert. Explicit invalid types remain invalid. | #7; upstream runtime default |
| Marked text inside `codeBlock` | Accept the upstream formatted-text and code-text variants in root and nested code blocks. Preserve literal text and omit marks inside the fence. Nontext content and malformed marks remain invalid. | #7 |

The heading and panel defaults are defined in
`@atlaskit/adf-schema@57.5.0/dist/esm/next-schema/generated/nodeTypes.js`.
Existing inline-media support also follows the runtime default of `file` when
`type` is omitted and accepts `image` media, both permitted by the upstream JSON
schema. Paragraph `fontSize: small` is likewise upstream-supported and renders
as `<small>…</small>`.

## Generation

```sh
GOEXPERIMENT=jsonv2 go generate ./...
GOEXPERIMENT=jsonv2 go test ./...
```

The generator updates the upstream draft-04 declaration to draft-07, which is
the draft supported by `jsonschema-go`, then applies the overlay to produce
`adf-persisted-api-schema-57.5.0.json`. Both source and generated schema are
committed; consumers do not need to run generation. The generator supports the
`add` and `replace` operations used by this overlay, not every RFC 6902 operation.

`schema_compat_jsonv2_test.go` covers issue #7, the new status colors, conversion
with and without the built-in validator, and malformed nearby shapes that must
still be rejected. Existing Confluence fixtures cover the readback corrections
and paragraph font size.
