# DE.PULSE — Current Adaptive Build Process

**Certified Stable:** `v18.8.1-stable`  
**Current activity:** `v18.8.2-development` has completed G0–G3 diagnosis/design freeze and the bounded G4 implementation for GitHub issue #57; CI Fast is the next gate.

Process remains:

**GitHub source of truth → reconcile actual Stable → freeze verified G1 scope → one version-development branch → one Draft PR → exact-head Fast → same PR Ready → Qualified → exact-head merge → one G11–G16 Release when release-capable → post-Stable continuity reconciliation before the next product line.**

## Post-Stable continuity rule

A Release G16 workflow artifact is necessary evidence but does not by itself update repository checkpoints/handoff/overlays. After Stable publication, the repository must converge on the actual Stable tag/candidate/fingerprint, Stable evidence manifest and exactly one next action before a new product branch begins.

The cheap `tools/ci/post_stable_continuity_gate.py` owns this repository-level convergence check. `main` may not silently carry a later STABLE package identity than the durable Stable checkpoint.

## v18.8.2 reliability execution

Issue #57 reopens `ADAPT-FRESHNESS-001` for the affected Market Intelligence path.

- G0: **COMPLETE** — allocator/provider demand already exists; freshness/recovery accountability is the escaped gap.
- G1: **FROZEN** — issue #57 only; TradeInsight and unrelated product scope excluded.
- G2: **FROZEN** — existing allocation → Smart Provider Router v2 → canonical freshness/recovery remains the only owner chain.
- G3: **FROZEN** — SPY/QQQ protected by existing allocation, VIX remains canonical special path, breadth remains existing bounded 15-symbol market context; unknown/unavailable is not observed zero.
- G4: **IMPLEMENTED** — existing breadth demand added to canonical quote freshness scope, unavailable-vs-zero renderer truth added, focused Go regressions added and renderer assertions integrated into the existing Fast lane. G4 exit is pending Fast evidence.
- G5: **NEXT** — one Draft PR, one automatic exact-head Fast run, classify any failure before changing or rerunning anything.
- G6-G10: integration, data truth, provider failure/recovery, performance/backpressure, browser/UI and exact-head reconciliation.
- G11-G16: immutable RC, full certification, native packages, actual artifact runtime audit, assurance, publication and learning.

Repair order follows REUSE → CONSOLIDATE → REFACTOR → DELETE/REPLACE → ADD. No Market Intelligence-specific fetch loop was introduced. Smart Provider Router v2 and canonical freshness/recovery remain sole executable owners. No workflow was added or duplicated.

Failure handling remains classify-first: `PRODUCT_FAIL`, `GATE_TEST_FAIL`, `CI_HARNESS_FAIL`, `INFRA_FAIL`, `EXPECTED_NOOP`, `SUPERSEDED`. No unchanged reruns or CI event manufacturing.

After v18.8.2, resume v18.9.0 TradeInsight SHADOW through Smart Provider Router v2, then v18.9.1 and v18.10.

No new top-level gates. G0–G16 remains the only release model.
