# DE.PULSE v18.5.1 — v17 + v18 Implementation Reconciliation Audit

**Status:** OPEN / STABLE BLOCKING  
**Baseline:** `v18.5.0-stable` at `0d37ca35f5fc3ad89cebed506cc5a4c2d6a7a680`  
**Machine ledger:** `release/v18.5.1/V17-V18-IMPLEMENTATION-RECONCILIATION.json`

## Why this audit exists

v17 and v18 used adaptive slicing, but v18 did not maintain one enforced major-scope ledger from original approval through final package certification. Local slice gates could pass while a major workstream remained unassigned.

TradeInsight proves the failure mode:

1. It was approved in `v18_major_scope.json`.
2. It was marked `COMMITTED` in `v18_delivery_slices.json`.
3. It was explicitly deferred from v18.0.
4. It was not added to v18.1, v18.2, v18.3, v18.4 or v18.5 immutable scope.
5. v18.5 reconstructed numbered minor-release scope but did not require every original major workstream to reach a current disposition.
6. No provider adapter, endpoint, configuration, router capability, rights record or test exists in current source.

The adaptive method therefore existed as documentation but lacked end-to-end enforcement.

## Root-cause analysis of the slicing failure

1. **No v18 equivalent of the v17 20-item major closure matrix.** v18 had a major workstream list and separate minor scopes, but no enforced one-to-one final disposition for every original workstream.
2. **Deferred work did not require reassignment.** TradeInsight was marked committed, deferred from v18.0 and then left without a target owner/release.
3. **Most v18 slice gates are release-identity pinned.** Earlier gates expect their original v18.0.x/v18.1/v18.2/v18.3 identity and cannot be rerun unchanged on v18.5.
4. **v18.5 checked closure metadata, not every earlier behavioral contract.** Its scope gate validates closure dimensions, scenario names, documents and boundaries; it does not enumerate all original v18 workstreams or directly rerun every earlier user workflow.
5. **The inherited 48-item gate is marker-oriented.** It confirms source ownership markers, while separate suites were assumed to prove behavior. The failing Research/desk workflows show that this separation did not guarantee requirement-level runtime acceptance.
6. **Conversation commitments were not all converted into durable requirement IDs.** GitHub currently has only the new v18.5.1 issues, so older reported items could rely on handoffs or conversation memory instead of one authoritative defect ledger.
7. **A duplicate report had no mandatory reopen rule.** “Already tracked” could suppress recurrence instead of invalidating the prior closure evidence.

The corrective model is version-independent, requirement-level and package-bound. A slice may optimize execution, but it cannot own the only record of whether a major commitment was eventually implemented.

## Audit universe

The ledger inventories:

| Layer | Current inventory |
|---|---:|
| Inherited approved product requirements that v17/v18 had to preserve | 48 |
| Frozen v17 major-scope items | 20 |
| v18 major workstreams | 8 |
| v18.0–v18.5 release items/clauses/contracts/scenarios | 169 |
| Recovered conversational commitments requiring explicit reconciliation | 7 |
| Confirmed implementation misses | 1 |
| Confirmed escaped defect groups | 5 |

The counts are inventory counts, not PASS counts.

## Current audit results

| Audit | Result | Meaning |
|---|---|---|
| Inherited 48-item production-source marker reconstruction | 48/48 markers present | Static ownership evidence exists; it does not prove working behavior. Research/desk items remain reopened by direct user evidence. |
| v17.0–v17.4 static scope-gate reconstruction | 0 missing contracts/files/test names | v17 implementation structure is still present; fresh test execution and package behavior remain required. |
| v18 major-workstream-to-release mapping | FAIL | TradeInsight was committed and deferred but never reassigned or implemented. |
| Fresh current behavioral campaign | NOT STARTED | Blocking. |
| Current macOS Apple Silicon package audit | NOT STARTED | Blocking. |
| Current Windows x64 package audit | NOT STARTED | Blocking. |

A static marker or test-name check can support traceability but can never close a user workflow. This distinction is now explicit because the earlier gates reported PASS while Research, symbol removal, desk-state styling and scroll continuity were still defective.

## Current confirmed blockers

1. `IMPL-18-TRADEINSIGHT-001` — TradeInsight SHADOW / SECONDARY integration is not implemented.
2. `COPY-18.5.1-001` — stale v17 profile-preservation copy.
3. `SYMBOL-18.5.1-001` — final-membership row removal failure across Day/Swing/Long.
4. `SYMBOL-18.5.1-002` — current desk membership is not distinctly represented.
5. `NAV-18.5.1-001` — same-page/live refresh can lose scroll/focus context.
6. `RESEARCH-v15.1.0-17-19-REOPENED` — Research ticker/freshness/top-area regression.

## High-risk earlier commitments requiring fresh proof

- Master Market Symbols control layout and semantics.
- Pre-Market/Market Open exception-list consolidation.
- Runtime overload / intermittent `DATA DEGRADED` / slow response.
- External Dependency & Provider Readiness checkpoint and User Action Required register.
- Role-aware User/Admin/Developer documentation, Documentation Impact Manifest and stale-doc enforcement.
- All 48 inherited approved product requirements, especially Research and desk/watchlist behavior.
- All 20 v17 persistence/runtime/reliability items.
- Every identity, Smart Router, multi-user, admin/presence, PostgreSQL, security/data-rights and closure clause from v18.

## Status policy

Historical release text, source-marker existence and an earlier PASS do not qualify as current proof.

Allowed current dispositions:

- `FRESH_PASS`
- `REOPENED`
- `NOT_IMPLEMENTED_CONFIRMED`
- `INTENTIONALLY_SUPERSEDED`
- `NOT_APPLICABLE`
- `ROADMAP_PLACED_FUTURE`
- `REVALIDATION_REQUIRED`
- `REVALIDATION_REQUIRED_HIGH`

v18.5.1 Stable is blocked by every `REOPENED`, `NOT_IMPLEMENTED_CONFIRMED`, `REVALIDATION_REQUIRED` and `REVALIDATION_REQUIRED_HIGH` row.

## Required evidence chain

Every applicable row must close with:

**origin → current code owner → current observed behavior → owner → fix/disposition → regression test → browser/runtime retest → macOS package proof → Windows package proof → closure approval**

A matching repeat report reopens the original row, increments recurrence and invalidates affected/dependent evidence. It is never discarded as “already tracked.”

## Audit execution order

1. Reconcile the 48 inherited requirements against current browser behavior and current code—not string markers alone.
2. Revalidate all 20 v17 items against current persistence/runtime behavior and supported packages.
3. Reconcile all v18 workstreams and 169 release entries; implement TradeInsight through the canonical Smart Router.
4. Fix the five escaped defect groups and add behavioral regressions.
5. Run cross-role, cross-desk, cross-viewport, restart/reload, failure/degradation and actual-package matrices.
6. Freeze an immutable RC only after the ledger has no unexplained or blocking status.
7. Record at G16 why prior gates missed each escape and which machine control prevents recurrence.

## Non-miss future scope

Mature ASBI, TDTI, AODR and adaptive 13F intelligence remain explicitly placed in v20 and are not incorrectly relabeled as current implementation misses. Their v18/v19 evidence-foundation obligations remain part of this reconciliation where applicable.
