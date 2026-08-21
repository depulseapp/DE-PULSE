# DE.PULSE — Current Adaptive Build Process

**Certified Stable:** `v18.8.2-stable`  
**Current activity:** v18.9.0 G0 exact-baseline / TradeInsight capability discovery before G1 scope freeze.  
**Active development branch:** none.

Process remains:

**GitHub source of truth → reconcile actual Stable → perform G0 exact-baseline → freeze verified G1 scope → one version-development branch → one Draft PR → exact-head Fast → same PR Ready → Qualified → exact-head merge → one G11–G16 Release when release-capable → post-Stable continuity reconciliation before the next product line.**

## v18.8.2 process closure

G0–G16 are PASS. Final Fast #437 / `32435845178`, Qualified #151 / `32435920048`, and Release #31 / `32436189650` are the release authority. macOS Apple Silicon and Windows x64 packaged-runtime audits, G15 evidence binding, no-rebuild publication and G16 all passed.

Release #30 exposed one release-process defect after G11: `version_consistency_test.py` coupled a semantic portability invariant to one obsolete README presentation heading. The durable learning is:
- release/readiness harnesses must validate semantic invariants rather than cosmetic presentation literals unless the literal itself is contractual;
- a proven post-merge harness-only failure must be classified before recovery;
- unchanged-SHA reruns remain prohibited;
- the only permitted recovery is the repository's bounded same-version-development-line exception: recreate the same branch name from the failed merged candidate, change only the proven harness/governance defect, earn fresh Fast + full Qualified, merge, and re-enter the same canonical Release workflow;
- never create retry/certification/promotion branches or duplicate workflows.

Repair order remains REUSE → CONSOLIDATE → REFACTOR → DELETE/REPLACE → ADD. Failure handling remains classify-first: `PRODUCT_FAIL`, `GATE_TEST_FAIL`, `CI_HARNESS_FAIL`, `INFRA_FAIL`, `EXPECTED_NOOP`, `SUPERSEDED`. No new top-level gates; G0–G16 remains the only release model.

## v18.9.0 entry rule

G0 must enumerate and disposition the configured TradeInsight beta capability surface before source implementation. G1 then freezes the bounded useful SHADOW scope and explicit exclusions. Any executable provider use must pass through Smart Provider Router v2; no provider-specific parallel router, data engine, freshness owner or subscription manager may be introduced.
