package main

import (
	"fmt"
	"path/filepath"
)

func writeDecisionRecord(root string, endpoints []endpoint) {
	writeText(filepath.Join(root, "docs", "integrations", "public-api-documentation.md"), fmt.Sprintf(`# Public API documentation decision

## Result

The public API contract is colocated with the version that owns it:

- `+"`contracts/gateway/v1/public-api/`"+` documents all %d Gateway v1 operations.
- `+"`contracts/realtime-gateway/v1/public-api/`"+` documents RealtimeGateway HTTP and WebSocket output.

Gateway REST/GraphQL-over-HTTP uses OpenAPI 3.1. RealtimeGateway HTTP also uses OpenAPI 3.1, while its multiplex WebSocket messages use AsyncAPI 3.0 because OpenAPI cannot fully describe bidirectional message channels and binary frames. Postman 2.1 collections and environments are committed as importable derived artifacts.

## Content boundary

These directories contain only API contracts, operational rules, complete request examples, and test-tool imports. Product usage teaching and UI workflows remain separate frontend concerns.

The human-readable directories are named `+"`rules/`"+` because they describe constraints: authentication, cookies, CSRF, CORS, rate limits, retry behavior, version compatibility, WebSocket admission, frames, backpressure, presence, and stable errors.

Each contract also owns `+"`versions/dev-log.md`"+` and `+"`versions/comparison.md`"+`. The current comparison records the single published v1 baseline and defines the fields that must be compared when another version is introduced, without inventing unreleased behavior.

## Authentication decision

During Beta and early v1, a consumer registers or logs in with an account and password. Gateway stores access and refresh JWTs in HttpOnly cookies. The client keeps a cookie jar and sends the CSRF token returned by register/login on protected mutations. API keys are intentionally deferred to a later Beta-to-v1 milestone.

This mechanism is suitable for a user's own client and trusted server-side integrations. Third-party browser apps must not collect Notezy credentials and also require an explicitly allowed Origin.

## Maintenance

Routes and Go DTO contracts remain authoritative. Run:

`+"```bash"+`
make -C contracts public-api-gen
`+"```"+`

The generator refreshes both OpenAPI documents, AsyncAPI, Postman files, examples, indexes, rules, and this decision record. CI should later run the generator and reject a dirty diff so route changes cannot silently bypass public documentation.

No generated example contains a real credential. Exported developer environments containing passwords, cookies, CSRF values, OAuth codes, or realtime tickets must never be committed.`, len(endpoints)))
}
