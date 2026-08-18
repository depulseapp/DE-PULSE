# G1 Immutable Scope — v18.6.0

**Incoming immutable Stable:** `v18.5.2-stable` at `d30e54db4908ca57c52ae298cc4ada3416fab46b`  
**Development branch:** `v18.6-development`  
**Selection rule:** smallest coherent evidence-selected slice from the conserved v17/v18 recovery ledger; no inherited PASS and no silent scope loss.

## Assigned parent requirement IDs

1. **`IMPL-18-UTILITY-001` + `UTILITY-V18.3-06` + `UTILITY-V18.3-07` — Shared Scanner/Radar snapshot acquisition.**
   - One bounded canonical broad-snapshot broker/cache.
   - Freshness-aware reuse and in-flight coalescing.
   - Preserve distinct Scanner vs Opportunity Radar ranking responsibilities.
   - Record avoided duplicate provider calls and prove bounded provider/runtime pressure.

2. **`IMPL-18-UTILITY-002` + `UTILITY-V18.3-08` + `UTILITY-V18.3-09` — Session Intelligence Coordinator.**
   - One canonical coordinator for Pre-Market Prep and Market Open Prep.
   - Preserve distinct checkpoint timing/semantics while sharing scheduler/router/cache/canonical state.
   - No parallel broad acquisition or duplicate market-analysis pipeline.
   - Cover late start, catch-up, restart and market-calendar behavior.

3. **`IMPL-18-UTILITY-003` + `UTILITY-V18.3-11` — Market Activity surface demotion.**
   - Retain Market Activity Seeds as canonical Discovery/Radar input and optional drill-down/explanation.
   - Remove it as an equal/prominent normal-user decision surface.

4. **`IMPL-18-UTILITY-004` + `UTILITY-V18.3-12` — Legacy evidence-route retirement.**
   - Ticker News/Earnings/Filings deep evidence resolves to canonical Research subviews.
   - Market-wide material-event evidence resolves to Market/Event Intelligence.
   - Preserve safe direct/deep links; no broken navigation.

5. **`IMPL-18-DOC-001` + `CONVO-V17-005` — Role-aware documentation and Documentation Impact Manifest.**
   - Server-authoritative audience/capability policy.
   - OWNER/SUPER_OWNER, delegated ADMIN, USER and DEMO composition and direct-path enforcement.
   - USER/DEMO receive no developer/implementation machinery.
   - Every material product/process change carries documentation-impact disposition.

6. **`IMPL-17-DEPS-001` + `CONVO-V17-004` — External Dependency & Provider Readiness.**
   - Canonical dependency/readiness registry with owner, capability, status, blocker, user action, rights/entitlement and evidence fields.
   - Durable User Action Required register.
   - Gate binding without creating G17+.
   - Cover provider, database, package/runtime and credential/config dependencies.

7. **`AUDIT-18-AI-001` — AI bounded-context, cache identity, strict schema and continuous eval hardening.**
   - Semantic/materiality-aware compaction with a hard byte/token bound.
   - Complete inference cache identity including provider/model/prompt/safety/schema fingerprint and TTL/invalidation.
   - Structured-output capability where supported; safe abstention/failure semantics otherwise.
   - Golden, citation, contradiction, missing-evidence and injection/adversarial eval lane with bounded cost/latency telemetry.

8. **`AUDIT-18-AI-RIGHTS-001` — Rights-aware AI egress.**
   - Canonical provider × dataset AI-use decision source.
   - Fail closed when AI/commercial/redistribution rights are unknown or denied.
   - Working API entitlement alone never implies AI-use rights.
   - Approved/denied fixtures, redacted diagnostics and package behavior proof.

## Mandatory affected-area sentinels — not separate scope expansion

- `CONVO-V18-003` Provider → Market Mode method remains enforced: every affected provider/dataset capability has an explicit disposition; `NOT_IMPLEMENTED`/SHADOW/VALIDATED evidence has no production Market Mode influence; deterministic/statistical code owns numeric Market Mode truth.
- `CONVO-V17-003` runtime overload/degradation remains a risk sentinel: the shared acquisition/coordinator changes must measure duplicate work removed, queue/provider pressure and freshness impact and must not broaden `DATA DEGRADED` blast radius.
- Shared Symbol Intelligence remains canonical: equivalent lawful work scales with unique canonical demand rather than users × symbols.
- Existing Event Intelligence/Catalyst behavior is preserved; v18.6 does not create a duplicate event engine.

## Explicitly not assigned to v18.6.0 G1

The following remain conserved and named; they are **not dropped**:

- `IMPL-18-TRADEINSIGHT-001` / `CONVO-V18-001` — TradeInsight SHADOW/SECONDARY implementation: v18.7 candidate.
- `UTILITY-V18.3-01..05`, `UTILITY-V18.3-10`, `UTILITY-V18.3-13` — Dashboard, Market Intelligence, Day/Swing/Long, Catalyst and Maintenance consolidation/revalidation: v18.7 candidate.
- `AUDIT-18-PROVIDER-001`, `AUDIT-18-ARCH-001`, `AUDIT-18-TRADER-001` — v18.7 candidate/foundation as governed.
- `AUDIT-18-SECURITY-001` — broader credential/HTTP/link/supply-chain hardening retains its existing v18.7 completion lane; v18.6 must not weaken security and may implement dependency-compatible prerequisites only without claiming parent closure.
- `AUDIT-18-PROVENANCE-001` and final `AUDIT-18-QA-001` closure — final v18.x zero-gap release work.
- Mature ASBI/TDTI/AODR/Adaptive-13F intelligence remains governed future scope; v18.6 may preserve/strengthen evidence foundations only.

## Protected invariants

1. U.S. Equities Processing Boundary unchanged.
2. Permanent No Execution Boundary unchanged: no orders, paper trading, portfolio/P&L or journal execution features.
3. Deterministic Day/Swing/Long formulas unchanged.
4. Canonical Smart Provider Router remains the only provider-routing owner.
5. Provider count alone never changes a Market Mode.
6. GLD/SLV/USO remain explicit high-priority tradable live exceptions.
7. Desktop SQLite / hosted PostgreSQL architecture truth remains unchanged unless an assigned dependency requires compatible non-breaking plumbing.
8. No secrets/API keys committed; rights/entitlement remain fail-closed.
9. Adaptive/AI behavior cannot silently self-promote or self-modify protected production decision logic.
10. G0–G16 remains the only top-level gate model.

## v18.6.0 closure evidence

Every assigned parent ID must reach current implementation + fresh regression + current behavior proof before G10. Affected user workflows require browser/API/state evidence and the final candidate requires actual macOS Apple Silicon and Windows x64 package proof at G13/G14 before G15.

At minimum prove:

- duplicate Scanner/Radar acquisition is removed without ranking drift;
- Prep checkpoints share canonical ownership without timing/semantic loss;
- route/surface consolidation causes no broken navigation or lost material evidence;
- documentation is role-safe in UI and direct paths;
- dependency/User Action Required state is actionable and truthful;
- AI context is bounded, cache-safe, schema-safe and adversarially evaluated;
- AI egress cannot use provider/dataset evidence when rights are unknown/denied;
- provider/runtime load, freshness and degradation behavior improve or remain bounded;
- deterministic equivalence and No Execution invariants pass.

No item is `FRESH_PASS` from documentation, source markers or historical tests alone.
