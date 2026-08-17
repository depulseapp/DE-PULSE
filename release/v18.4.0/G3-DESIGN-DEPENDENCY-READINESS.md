# DE.PULSE v18.4.0 — G3 Design / Dependency Readiness

Status: PASS

Delivery is delta-first. Slice 1 hardens the web/auth perimeter: explicitly trusted hosted proxy headers, effective secure-request detection, same-origin enforcement for unsafe browser requests, browser security headers/cookies, bounded login-abuse throttling, and recent-authentication groundwork. Forwarded headers are fail-closed unless hosted proxy trust is explicitly enabled.

Provider entitlement/data-rights and broader licensing/redistribution/AI-use suitability are separate bounded slices because they have different owners and evidence. G10 reconciles all eight clauses before freeze; v18.5 closure is not pulled forward.
