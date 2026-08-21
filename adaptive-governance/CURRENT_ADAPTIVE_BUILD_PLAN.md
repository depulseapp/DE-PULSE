# DE.PULSE — Current Adaptive Build Plan

**Immutable Stable:** `v18.9.0-stable`  
**Stable candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Stable fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Active development branch:** none  
**Active PR:** none  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Immediate next patch:** `v18.9.1` / #64 / `ADAPT-RUNTIME-CRASH-001`.

## Build philosophy — completeness through small patches

The v18.9.x line MUST NOT use heavy multi-domain builds. Each patch owns one primary responsibility and its directly necessary supporting work only. The purpose is to reduce blast radius, make review/test scope exact, catch implementation misses early and keep CI/release evidence interpretable.

Before any patch implementation:
- G0 re-fetches live GitHub and exact predecessor Stable/patch truth;
- G1 freezes one bounded scope plus explicit non-goals;
- G2 maps canonical owners and proves no parallel subsystem is needed;
- G3 freezes dependency/provider contracts and deterministic acceptance tests.

Before the next patch may begin:
- G4–G10 evidence for the current patch must be coherent;
- implementation-miss audit must be performed against G1 acceptance;
- all discovered out-of-scope findings must have a durable target patch/issue;
- open issues must be reconciled;
- handoff/checkpoints must identify exactly one next action.

## Ordered v18.9.x build packets

### v18.9.1 — Runtime Reliability
Scope owner #64. Crash diagnosis/fix only. No provider expansion, router redesign, Market Modes or company-name work. Exact root cause/reproduction, lifecycle regression, persisted-state/API-key continuity and packaged macOS Apple Silicon proof required.

### v18.9.2 — TradeInsight Settings/API Key
Settings integration only: existing Data Provider Settings + existing local secret owner, masked field, Save/Test/Clear, configured/connected/error/capability status, safe environment override. No provider-routing behavior change.

### v18.9.3 — Coverage-Aware Router Core
Smart Provider Router v2 evolution only. Add consumer requirement/coverage contracts, cache-first gap calculation, eligible-provider ranking against remaining need, targeted acquisition, canonical merge/provenance/conflict handling, coverage re-evaluation and bounded stop criteria. Separate validation lifecycle from serving role. No UI/provider-feature expansion beyond what tests require.

### v18.9.4 — Company Identity
Canonical identity owner and presentation only. All Day/Swing/Long state headings show symbol + company name when known, e.g. `APP - AppLovin : In Entry Zone`; reuse in Research/Discovery/Add Symbol; symbol-only fallback. No TradeInsight search admission yet.

### v18.9.5 — Market Data Modes
Behavior/quality-oriented Adaptive modes and capability diagnostics only. Remove hard-coded provider-brand semantics where misleading. Surface actual provider contribution/freshness/coverage in diagnostics without creating a new Market Mode owner.

### v18.9.6 — TradeInsight Form 4
Contract-validated Corporate Insider/Form 4 enrichment only. SHADOW-first, direct SEC/EDGAR authoritative, source-family de-duplication, existing SEC/Ownership model reused, optional provider failure non-degrading.

### v18.9.7 — TradeInsight Symbol Search
Contract-validated ticker/company search only. Plug into canonical symbol validation/company identity as fallback/corroboration. U.S.-equity boundary final; GLD/SLV/USO exceptions preserved.

### v18.9.8 — TradeInsight Movers
Contract-validated mover/ranking evidence only. Opportunity Radar consumes it as SHADOW candidate evidence; existing scanner/ranker remains canonical. No undocumented REST endpoint assumptions.

### v18.9.9 — Remaining TradeInsight Capability Admission
Full useful-capability inventory/disposition only. Revalidate Congress, daily adjusted/raw OHLCV, corporate actions and bounded history under coverage-aware routing. Every useful entitlement needs a named consumer/disposition/freshness/retention/rate policy. No Python/MCP production dependency and no inferred intraday support.

### v18.9.10 — Provider Efficiency / Adaptive Telemetry
Measure and harden only: coverage completion, residual gaps, calls avoided, cache reuse, provider usefulness, latency/errors/rate limits/backpressure, freshness failures, disagreement/corroboration, bounded fan-out, CPU/memory/goroutine stability and consumer materiality. Feed evidence into adaptive ranking without changing deterministic trading truth opaquely.

### v18.9.11 — Professional Closure Audit
No new feature scope. Audit whole v18.9.x for misses/duplication/bypasses/orphan useful capabilities; retest #57 and #64 regressions; deterministic Day/Swing/Long equivalence; all-desk identity; partial coverage/provider failure/recovery; actual macOS Apple Silicon + Windows x64 packages; Adaptive Intelligence Scorecard. Closure requires zero unexplained carry-forward and zero unowned useful capability.

## CI/release efficiency rule

For each patch: develop the coherent code+test batch before opening the PR; one PR only; exact-head Fast once for the coherent candidate, Qualified once when Ready, then one canonical G11–G16 release when release-capable. Do not spend CI budget on avoidable duplicate runs, retry branches or certification branches. Failure classification remains mandatory before rerunning.

## Protected boundaries

Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; direct SEC/EDGAR authoritative; existing cache/persistence/telemetry/symbol/state owners reused; deterministic Day/Swing/Long protected; U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution; G0–G16 only.

## Exactly one next action

Execute G0 for issue #64 / v18.9.1 using concrete macOS crash evidence or deterministic reproduction, then freeze its narrow G1. No v18.9.2 branch until v18.9.1 is truthfully closed.
