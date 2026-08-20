# DE.PULSE — Current Adaptive Build Process

**Certified Stable:** `v18.8.0-stable`  
**Current activity:** `v18.8.1-development` governance hardening and G0 entry preparation.

Process remains:

**GitHub source of truth → reconcile actual Stable → freeze only verified G1 scope → one version-development branch → one Draft PR → exact-head Fast → same PR Ready → Qualified → exact-head merge → one G11–G16 Release when release-capable.**

## v18.8.1 process hardening

- Run `ADAPT-REL-001` first so v18.8.0 Stable tag/Release, checkpoints, Stable evidence manifest, handoff and CURRENT overlays agree before normal product work.
- Add `ADAPT-CI-001` Release State Coherence so one cheap preflight reports every release-state mismatch instead of discovering VERSION/checkpoint/manifest/handoff drift sequentially.
- Add `ADAPT-CI-002` at G11 so a conflicting target Stable tag/version/build/predecessor fails before G12/native certification. Keep the publication-time tag guard too.
- Reorder Fast per `ADAPT-CI-003`: Python/impact/coherence/governance/identity/provenance before expensive Go/Node/browser setup.
- Manual CI defaults follow `ADAPT-CI-004`: smallest safe/adaptive lane by default; `full` requires explicit intent.
- For market evidence, `ADAPT-DATA-002` separates provider evidence time from retrieval time. Missing provider time never becomes `time.Now()` freshness.
- For Scanner/Radar, `ADAPT-DATA-001` makes universe eligibility explicit; acquisition filters such as `has_options` cannot silently redefine the advertised broad U.S.-equity universe.
- `ADAPT-ARCH-001` keeps one neutral shared-universe owner with cancellation/panic-safe refresh cleanup and no duplicate router/broker/freshness engine.
- Renderer/test cleanup uses strangler migration: stable capability owners first, equivalence proof, then retire superseded release-number owners. No mass delete/rename.

## Zero-miss carry-forward execution

`ADAPT-RECON-001` is an explicit process responsibility, not a generic v18.10 note. At v18.8.1 G0–G3, load the conserved v17/v18 ledger plus post-ledger approved Adaptive requirements and map each applicable item to current owner, current behavior, regression/evidence and one current disposition. Prioritize `ADAPT-UX-RESEARCH-001`, `ADAPT-SYMBOL-001`, `ADAPT-READINESS-001`, `ADAPT-FRESHNESS-001` and `ADAPT-RESEARCH-002` for fresh current-source proof.

If current implementation is correct, close by fresh evidence rather than rebuilding it. If a gap reproduces, keep the original requirement identity, invalidate affected evidence and implement the fix in its assigned coherent slice. G10 cannot hide an unexplained item behind a broad “reconciliation” label, and v18.10 is final verification rather than the first discovery point.

## TradeInsight full-capability execution

At v18.9.0, `ADAPT-TRADEINSIGHT-001` executes through Smart Provider Router v2 only:

1. enumerate every capability actually exposed to the configured TradeInsight beta account/API;
2. classify each as `USE`, `CORROBORATE`, `FALLBACK`, `STORE_FOR_HISTORY`, `FUTURE`, or `NOT_USEFUL`;
3. for every useful capability bind canonical owner/consumer, source authority, entitlement/rights, provider evidence time/freshness, rate/cost budget, cache/retention, disagreement and degradation behavior;
4. give every capability an explicit provider → Market Mode disposition and preserve deterministic Market Mode ownership;
5. run bounded SHADOW fixtures/live smoke and collect availability, latency, freshness, completeness, disagreement, unique evidence, usefulness, calls avoided, errors/rate pressure and cost;
6. promote only through `SHADOW → VALIDATED → APPROVED → PRODUCTION`, with rollback and no silent adaptive production change;
7. never equate “full capability” with “call every endpoint”: request only purpose-driven evidence, reuse canonical state and isolate TradeInsight beta failures from unrelated consumers;
8. prohibit a direct TradeInsight-to-UI silo or duplicate router/freshness/reconciliation owner.

v18.9.1 consumes the resulting capability evidence for smarter multi-provider routing. v18.10 proves no materially useful TradeInsight/provider capability is orphaned or unexplained. v19/v20 mature quality/history/adaptive usefulness without changing deterministic market truth.

Failure handling remains classify-first: `PRODUCT_FAIL`, `GATE_TEST_FAIL`, `CI_HARNESS_FAIL`, `INFRA_FAIL`, `EXPECTED_NOOP`, `SUPERSEDED`. No unchanged reruns or CI event manufacturing.

Historical v18.5.1 reconciliation remains provenance, not current-release identity. Actual GitHub objects and CURRENT overlays identify the current Stable/release line.

No new top-level gates. G0–G16 remains the only release model.