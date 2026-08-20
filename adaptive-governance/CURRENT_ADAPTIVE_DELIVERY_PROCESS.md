# DE.PULSE — Current Adaptive Delivery Process

**Certified Stable:** `v18.8.0-stable`  
**Stable candidate:** `3a32d57dd4c74c6f812cc942a9d8049a7b517718`  
**Fast:** #419 / `32336519003`  
**Qualified:** #145 / `32336619446`  
**Release:** #28 / `32336898662`  
**Next delivery line:** `v18.8.1 — Renderer Modularization II + 10/10 Audit Hardening`.

v18.8.0 delivery is complete: exact-head qualification, G11–G16, macOS Apple Silicon + Windows x64 actual packaged-runtime proof, G15 assurance and same-artifact no-rebuild Stable publication all passed.

## Immediate post-Stable continuity closure

`ADAPT-REL-001` is the first v18.8.1 entry responsibility. Reconcile the actual v18.8.0 Stable tag/Release into durable build/release checkpoints, Stable evidence manifest and `handoff/CURRENT.md`. This does **not** justify rebuilding or recertifying unchanged v18.8.0 native artifacts.

## v18.8.1 delivery hardening

- Release State Coherence must be green before promotion and must aggregate release metadata/checkpoint/handoff mismatches.
- G11 must reject a conflicting target Stable tag/version/build/predecessor before G12/native spend; publication repeats the immutable collision guard.
- Actual packages must preserve provider evidence time separately from retrieval time; missing provider time cannot be presented as fresh `now`.
- Scanner/Radar package behavior must match an explicit discovery-universe eligibility contract.
- Renderer/test modularization must prove browser/native behavior equivalence for each migrated capability; no mass cleanup without evidence lineage.
- G15/G16 report CI work avoided by cheap-first checks/evidence reuse and confirm no mandatory evidence was removed.
- G16 is incomplete until actual Stable release identity, checkpoints, evidence manifest, handoff and CURRENT overlays converge and a new assistant/account can identify exactly one next action from GitHub alone.

## Zero-miss delivery proof

For `ADAPT-RECON-001`, `ADAPT-UX-RESEARCH-001`, `ADAPT-SYMBOL-001`, `ADAPT-READINESS-001`, `ADAPT-FRESHNESS-001` and `ADAPT-RESEARCH-002`, delivery evidence must distinguish three outcomes: freshly proven unchanged behavior, newly fixed behavior, or explicit approved future placement. Historical PASS/source markers alone are insufficient. Any affected user-facing behavior must carry appropriate browser/runtime evidence and actual macOS Apple Silicon + Windows x64 package proof before the responsible release claims closure.

v18.10 is the final zero-gap verification point, not a holding area for known gaps. No applicable conserved requirement or post-ledger Adaptive requirement may reach v18 closure as missing, unowned, unexplained or supported only by generic reconciliation wording.

## TradeInsight full-capability delivery proof

For v18.9.0 `ADAPT-TRADEINSIGHT-001`, delivery must include a capability inventory for the actual configured TradeInsight beta surface and a disposition for every discovered capability: `USE`, `CORROBORATE`, `FALLBACK`, `STORE_FOR_HISTORY`, `FUTURE`, or `NOT_USEFUL`.

Every capability used by DE.PULSE must ship with explicit Smart Provider Router v2 ownership, consumer/purpose, authority, rights/entitlement, evidence-time/freshness, budget/rate-limit, cache/retention, disagreement/degradation behavior, Market Mode disposition and SHADOW status. The guaranteed minimum roles remain Congressional Trading Intelligence, SEC Form 4 enrichment secondary to direct SEC, and historical OHLCV fallback/backfill; additional useful beta capabilities are included when the audit proves purpose and lawful access.

Package/runtime proof must demonstrate that full capability does not mean indiscriminate API traffic: canonical evidence is reused, calls are purpose-driven and bounded, TradeInsight outages/rate limits do not broadly degrade unrelated consumers, and there is no direct provider-to-UI silo. Any production influence requires `SHADOW → VALIDATED → APPROVED → PRODUCTION` evidence and rollback.

v18.9.1 delivery must prove the wider provider pool can use measured capability health/freshness/latency/disagreement/headroom/rights/cost/usefulness without changing deterministic market truth. v18.10 must show zero orphan materially useful provider capability or duplicate routing owner. v19/v20 carry mature provider/data infrastructure and adaptive usefulness/outcome learning.

Normal delivery remains: one development branch → one Draft PR → automatic Fast → same PR Ready → Qualified → exact-head merge → one canonical `release.yml` G11–G16 run when release-capable → exact certified artifacts published without rebuild.

No retry/certification/promotion branches, duplicate release workflows or manual duplicate CI runs. G0–G16 only.