# v18.5 G4/G5 Qualification Request

Status: REQUALIFICATION REQUESTED AFTER DOCUMENTATION / ARCHIVE SYNC

Current documented source head: `9ca706ce3576667630c51a66836ba2a8df034689`.

This marker requests a fresh G0–G5 qualification after v18.5 TEST documentation, QA manifest, v18 provenance anchor, release-artifact archive policy, diagrams, and explicit recovery paths were synchronized. The change is governance/documentation/release-metadata only; product runtime logic is unchanged.

The qualification must remain fail-closed. A PASS may be recorded only after the live v18.5 workflow completes G0–G3 scope/architecture/documentation checks, source/static/developer-proof checks, full and focused inherited degradation/backpressure/readiness tests, performance/stability gates, focused race coverage, and CGO-disabled fallback checks. Historical `go test -run` selector families must be proven non-empty before their results are trusted.
