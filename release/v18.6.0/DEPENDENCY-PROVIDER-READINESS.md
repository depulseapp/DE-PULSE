# External Dependency & Provider Readiness — v18.6.0

**Requirement:** `IMPL-17-DEPS-001` / `CONVO-V17-004`  
**Gate model:** existing `G0–G16` only  
**Canonical registry:** `dependency_readiness_registry.json`  
**Durable user-action register:** `user_action_required_registry.json`

## Contract

v18.6 makes external dependencies explicit instead of treating a configured endpoint or successful request as sufficient readiness evidence. Every registered dependency carries an owner, capability, readiness status, blocker, user action, rights/entitlement disposition, and source evidence.

Operational provider readiness remains owned by the existing Provider Capability Registry and Smart Provider Router. This slice does **not** create a second provider router, does not alter deterministic Day/Swing/Long formulas, and does not let provider count influence Market Mode.

Provider licensing/data rights remain separate from technical entitlement. `provider_data_rights.go` continues to fail closed: a working credential, public endpoint, successful response, or plan entitlement never implies commercial-use, redistribution, or AI/LLM-use permission.

## Coverage

The canonical registry covers:

- market-data and AI providers;
- desktop SQLite and conditional hosted PostgreSQL;
- Go, Node, Python, Chromium/browser proof, and `govulncheck` release tooling;
- live-data, AI-provider, and hosted-runtime credential/config dependencies.

Desktop persistence remains SQLite. Hosted shared-state persistence remains PostgreSQL and requires a `postgres`-tagged artifact plus `DEPULSE_PERSISTENCE_BACKEND=postgres` and `DEPULSE_DATABASE_URL`. The non-PostgreSQL build fails closed rather than silently substituting another hosted database.

## Durable User Action Required policy

`user_action_required_registry.json` preserves unresolved or conditional actions across conversations/releases. It currently records:

1. provider-specific commercial/redistribution/AI-use rights evidence binding before commercialization or external AI egress;
2. hosted PostgreSQL build/configuration requirements when a hosted target is produced;
3. live Finnhub/Alpaca credential and entitlement setup when live primary equity capability is enabled;
4. selected AI-provider credential setup, explicitly separated from source-data AI-use rights.

No secret value is stored in either registry. Credentials continue through the existing Settings/secret-storage path or deployment secret store.

## Gate binding

`dependency_readiness_gate.py` is bound into CI Fast and CI Qualified. It verifies:

- all four required dependency categories are represented;
- required fields and evidence exist;
- dependency and user-action identifiers are unique and linked;
- market-provider rights remain review-required/fail-closed;
- user actions bind only existing `G0–G16` gates;
- the Documentation Impact Manifest contains the dependency-readiness disposition;
- CI Go pins match the release-environment preferred version and remain inside an approved minor line.

The v18.6 readiness slice therefore adds evidence to the established gates; it creates no `G17+`.

## Readiness interpretation

`RUNTIME_EVALUATED` means source/build integration exists but operational readiness depends on current secret, entitlement, health, quota, or selected-mode state. `CONDITIONAL` means the dependency is required only for the named deployment/capability. `READY_NO_CREDENTIAL` means no credential setup is required, not that the upstream source can never fail. `CI_ENFORCED` and `RELEASE_ENFORCED` mean the dependency is proven by the corresponding pipeline stage.

This distinction prevents static source-control evidence from pretending to know live credentials, paid-plan entitlement, provider uptime, or deployment secrets.
