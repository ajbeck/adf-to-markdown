# Confluence persisted-API compatibility

Source: Confluence Cloud v2 page API `atlas_doc_format` read-back, observed
2026-08-19. This fixture is minimized and sanitized from the structural API
shapes only; it contains no customer page content or identifiers.

It covers an annotated inline card, a `mediaInline` node with no `type` and an
API-added `__fileName`, and a task item with a null `localId`.
