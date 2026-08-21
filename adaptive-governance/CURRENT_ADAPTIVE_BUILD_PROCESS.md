# DE.PULSE — Current Adaptive Build Process

**Certified Stable:** `v18.9.0-stable`  
**Current program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Immediate activity:** v18.9.1 corrective diagnosis for #64 / `ADAPT-RUNTIME-CRASH-001`  
**Active development branch:** none.

## Permanent process

**GitHub source of truth -> exact G0 baseline -> one bounded G1 scope -> G2 canonical-owner audit -> G3 dependency/contract readiness -> one version-development branch -> coherent code+tests -> one Draft PR -> exact-head Fast -> same PR Ready -> Qualified -> exact-head merge -> one canonical G11-G16 Release when release-capable -> post-release implementation-miss audit -> continuity reconciliation -> next patch.**

## Small-patch discipline

The v18.9.x line deliberately prefers multiple small patches over heavy builds.

1. One primary responsibility per patch.
2. Explicit non-goals are mandatory at G1.
3. Do not combine crash correction, router architecture, provider admission, company-identity UX and Market Mode redesign in one patch.
4. Reuse existing owners first. Repair order remains `REUSE -> CONSOLIDATE -> REFACTOR -> DELETE/REPLACE -> ADD`.
5. Every patch must have deterministic acceptance tests tied directly to its G1 scope.
6. Before closure, re-audit implementation against the frozen scope and search for misses, bypasses, duplicate ownership and misleading UI truth.
7. Any newly discovered out-of-scope miss must receive a durable issue/target patch before closure; chat-only carry-forward is prohibited.
8. The next patch cannot begin until current patch evidence, issue state, handoff and checkpoints agree.
9. CI is intentionally economical: batch coherent source/test changes before PR, avoid duplicate/manual runs, classify failures before rerun, and never create retry/certification branch families.
10. G0-G16 remains the only gate model; no G17+.

## v18.9.x process sequence

- `v18.9.1`: #64 runtime crash only.
- `v18.9.2`: TradeInsight Settings/API-key UX only.
- `v18.9.3`: coverage-aware Smart Provider Router v2 core only.
- `v18.9.4`: canonical company identity/all-desk presentation only.
- `v18.9.5`: Market Data Modes/capability diagnostics only.
- `v18.9.6`: TradeInsight Form 4 enrichment only.
- `v18.9.7`: TradeInsight ticker/company search only.
- `v18.9.8`: TradeInsight movers/ranking evidence only.
- `v18.9.9`: remaining useful TradeInsight capability admission only.
- `v18.9.10`: provider efficiency/Adaptive Intelligence telemetry only.
- `v18.9.11`: professional closure audit only; no new feature scope.

## Adaptive provider process contract

Provider work must begin from a consumer requirement and existing canonical evidence, not from a provider-first fetch loop. The sole routing owner evaluates missing coverage and eligible providers, acquires only what is needed, merges with provenance, re-evaluates remaining gaps and stops only when the bounded requirement is met or eligible budget is exhausted.

A provider response marked successful does not imply consumer completeness. Static provider ordering is only a prior/tiebreaker. TradeInsight is never allowed to create its own router/cache/scanner/Market Mode/SEC truth/symbol/persistence system.

Provider validation lifecycle and runtime serving role are distinct concepts. SHADOW/VALIDATED/APPROVED describe evidence maturity; PRIMARY/FALLBACK/BACKFILL/ENRICH/CORROBORATE describe serving purpose. Promotion/demotion requires telemetry/evidence and must not silently alter deterministic Day/Swing/Long truth.

## Failure handling

Classify before action: `PRODUCT_FAIL`, `GATE_TEST_FAIL`, `CI_HARNESS_FAIL`, `INFRA_FAIL`, `EXPECTED_NOOP`, `SUPERSEDED`. Never weaken a gate to make a patch pass. A real post-certification escape becomes the next learning/corrective loop without rewriting prior release evidence.

## Permanent owners/boundaries

Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; direct SEC/EDGAR authoritative; existing canonical cache/persistence/telemetry/symbol/state owners; deterministic Day/Swing/Long truth; U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution.

## Exactly one next action

Run #64 / v18.9.1 G0 from complete macOS crash evidence or deterministic reproduction and freeze the narrow G1 before any product-source change.
