# DE.PULSE v18.5.1 — Hardened Final Closure Build Plan

Status: **OPEN / CURRENT RELEASE BLOCKER**  
Authority: `governance/ROADMAP.md` + `release/v18.5.1/V17-V18-IMPLEMENTATION-RECONCILIATION.json`  
Process: **G0–G16 only; no G17+**

## 1. Purpose

v18.5.1 is the final v18 closure build. It is not a cosmetic patch, a repository-only refactor, or another pass over old certification labels. Its purpose is to prove that every applicable v17/v18 commitment is implemented, behaves correctly in the current product, survives realistic failure/load conditions, and works in the actual macOS Apple Silicon and Windows x64 packages.

The earlier work is retained where it is still valid. Architecture, tests, provider controls, persistence, security, packaging and governance evidence are reusable inputs. What is rejected is inherited closure: an earlier PASS does not prove the current source, behavior or package.

## 2. Lessons converted into controls

| Learned failure | Permanent correction |
|---|---|
| A requirement was documented but disappeared between release slices. | The requirement ledger is the parent authority. Every slice contains ledger IDs and must conserve the parent count. |
| A remediation plan existed, so a gate passed even though implementation never entered scope. | Plan presence is not implementation evidence. Target assignment, source owner, behavior, regression and package proof are independently required. |
| Old gates were tied to old version strings. | Reconciliation is release-independent; identity assertions are separated from behavioral contracts. |
| Static markers/test names were treated as functional closure. | Static evidence is traceability only. User workflows require current executable/browser and package evidence. |
| Conversation promises were not durable. | Every accepted promise or defect receives an immutable ledger ID before G1 closes. |
| Repeat reports were treated as duplicates. | Reproduction reopens the original row and invalidates affected evidence. |
| Major Closure tested metadata but did not reconstruct every earlier promise. | G10 and G15 run full requirement conservation and evidence-freshness checks. |
| Packaging proof was inherited. | G13/G14 use the exact immutable candidate and directly exercise affected workflows on both required platforms. |

## 3. Authority and scope hierarchy

The authority order is:

1. canonical product roadmap;
2. v17/v18 reconciliation ledger;
3. this Final Closure Build Plan;
4. release/slice scope documents;
5. implementation tasks and tests;
6. evidence and certification reports.

A lower layer may refine a higher layer but cannot remove, rename, defer or close a parent requirement without updating the parent ledger and recording the approved reason.

### Requirement conservation

For every planning or build transition:

```text
input requirement IDs
= implemented
+ intentionally superseded
+ not applicable
+ explicitly approved future placement
+ still open
```

No IDs may be missing from the right side. Any count mismatch blocks G1, G3, G10 and G15.

Every applicable requirement moves through:

`PLACED → IMPLEMENTED → REGRESSION_PROVEN → BEHAVIOR_PROVEN → ARTIFACT_PROVEN → FRESH_PASS`

A requirement cannot skip directly from planned or source-present to `FRESH_PASS`.

## 4. Frozen closure inventory

The authoritative ledger currently contains 282 rows:

- 48 inherited approved requirements;
- 20 frozen v17 major items;
- 8 v18 major workstreams;
- 169 v18.0–v18.5 release entries;
- 13 orphaned functionality/utility remediations;
- 7 recovered conversational commitments;
- 7 confirmed implementation-miss groups;
- 6 escaped-defect groups;
- 4 approved future-only mature workstreams.

The four mature v20 workstreams are not current implementation misses. Their v18/v19 foundations remain applicable where already approved.

No additional row may be inserted only into a slice, issue or conversation. It must first enter the ledger with origin, owner, placement, affected behavior and required evidence.

## 5. Build waves

### Wave 0 — Reconstruct and freeze (G0–G3)

- Freeze the exact v18.5.0 Stable source/package baseline and v18.5.1 branch.
- Validate all 282 IDs, origins and counts.
- Add code owner, build lane, dependency graph, regression ID and evidence matrix for every applicable row.
- Record source/dependency/UI/state/persistence/security/package blast radius.
- Reject any slice whose input/output ID count is not conserved.
- Do not create RC or Stable packaging.

Exit: zero missing, duplicate, unowned or unplaced applicable IDs.

### Wave 1 — Escaped user defects (G4–G6/G9)

Fix and prove:

- stale v17/current-version copy;
- Day/Swing/Long final-membership removal;
- current-desk membership highlighting;
- scroll/focus/selection preservation through live/SSE and save/action rerenders;
- Research target/freshness/action layout and recovery behavior.

Required matrices include all desks, first/middle/final membership, reload/restart, desktop/tablet/narrow widths, keyboard focus and repeated live updates.

Exit: requirement-linked automated regression plus direct browser proof for each defect.

### Wave 2 — Confirmed implementation misses (G4–G8)

Implement and prove:

- TradeInsight bounded SHADOW/SECONDARY capabilities through Smart Router v2;
- shared Scanner/Radar snapshot broker/cache and in-flight coalescing;
- one Session Intelligence Coordinator for Pre-Market and Market Open;
- Market Activity seed surface demotion;
- News/Earnings/Filings legacy route retirement/redirection;
- role-aware documentation and Documentation Impact Manifest;
- External Dependency & Provider Readiness checkpoint plus User Action Required register.

TradeInsight remains fail-closed for entitlement and data rights. It cannot replace Finnhub/Alpaca, direct SEC truth or deterministic Day/Swing/Long logic.

Exit: canonical ownership, tests, failure/degraded behavior, performance/provider-budget proof and security/rights evidence.

### Wave 3 — Remaining v18.3 remediation closure (G4–G9)

Revalidate and either implement or evidence-close:

- Dashboard consolidation;
- Market Intelligence consolidation;
- Day Trade Desk consolidation;
- Swing Desk consolidation;
- Long-Term Desk consolidation;
- Catalyst Watch canonical ownership;
- Maintenance consolidation.

Exit: one canonical intelligence owner, one deep-evidence home and concise contextual reuse elsewhere. No registry-only PASS.

### Wave 4 — Full v17/v18 implementation proof (G5–G10)

For every applicable remaining row:

- identify current source owner;
- run focused unit/integration tests;
- exercise current HTTP/state/persistence behavior;
- exercise browser UI where applicable;
- validate failure, degraded, stale, restart and recovery semantics;
- attach exact evidence to the ledger.

Old evidence may be reused only when source and all dependencies are unchanged and its fingerprint remains valid. Otherwise it reruns.

Exit: every applicable row has a current final disposition and evidence; no `REVALIDATION_REQUIRED`, `REOPENED` or `NOT_IMPLEMENTED_CONFIRMED`.

### Wave 5 — Cross-cutting hardening (G7–G10)

Run the combined hardening campaign:

- OWNER/ADMIN/USER/DEMO authorization and role-composed UI;
- Day/Swing/Long/Discovery/Research/Market Intelligence/Maintenance/Settings;
- responsive desktop/tablet/narrow layouts and keyboard navigation;
- persistence migration, warm start, reload/restart and state isolation;
- provider outage, rate limit, disagreement, entitlement and stale data;
- PostgreSQL/SQLite pressure, recovery and readiness;
- queue saturation, backpressure, load shedding and recovery hysteresis;
- Scanner/Radar/Prep/Catalyst/background-work concurrency;
- performance/freshness SLOs under realistic supported load;
- deterministic formula and No Execution invariance;
- documentation audience, freshness and impact coverage;
- source/repository hygiene and dead/duplicate-route checks.

Exit: no broad unexplained `DATA DEGRADED`, no self-inflicted freshness delay, no cross-role data exposure, no unbounded work, no protected-formula drift.

### Wave 6 — Zero-gap pre-freeze (G10)

G10 is blocked unless:

- ledger count is exactly conserved;
- all applicable rows have final statuses;
- no missing owners, tests or evidence bindings;
- no open/reopened/not-implemented rows;
- all deferrals/supersessions have explicit approved placement/replacement;
- affected evidence fingerprints match the candidate source;
- the full regression, UI/UX, security, performance and repository suites pass.

Exit: one reviewed candidate is eligible for RC.

### Wave 7 — Immutable RC and actual packages (G11–G15)

- Freeze one immutable RC.
- Run full/randomized/race/static/security certification on that exact source.
- Build macOS Apple Silicon and Windows x64 packages once from the frozen RC.
- Directly exercise every affected workflow in both packages.
- Verify provenance, hashes, source fingerprint, migration/startup and rollback.
- Publish only already-certified artifacts; do not rebuild after certification.
- Any defect reopens the row, invalidates affected evidence and returns to the earliest impacted gate.

Exit: G15 GO only with zero unexplained gaps.

### Wave 8 — Learning and handoff (G16)

- Record root cause and prevention for each escaped/repeated item.
- Record what evidence was invalidated and why.
- Update regression ownership, documentation and User Action Required records.
- Preserve exact source/tag/package hashes and release artifacts.
- Hand off only explicit future work; no current v17/v18 item may disappear into a generic “later” note.
- Confirm the next release can reconstruct state without conversation memory.

## 6. Test and evidence model

Every applicable ledger row must bind:

- immutable requirement/defect ID;
- origin and approved wording;
- code owner and implementation commit;
- affected modules/surfaces/dependencies;
- automated regression ID and result;
- browser/API/state/persistence evidence as applicable;
- macOS Apple Silicon actual-package result;
- Windows x64 actual-package result;
- evidence fingerprint and invalidation dependencies;
- closure approver and date;
- recurrence/root-cause/prevention record where applicable.

A screenshot alone, source marker alone, unit test alone or historical package alone cannot close a user workflow.

## 7. Efficient execution without weakening closure

Use the CI/CD model inside G0–G16:

1. run inventory, syntax, scope and cheap static checks first;
2. run independent backend, renderer, security, persistence and documentation lanes in parallel;
3. serialize shared-state migrations and protected-formula checks;
4. run focused affected tests before expensive full suites;
5. merge evidence at G10;
6. build native packages once from the immutable RC;
7. rerun only invalidated evidence until RC; after RC, any source change creates a new candidate.

Parallel execution may reduce elapsed time. It may never split requirement authority or create multiple competing closure ledgers.

## 8. Promotion blockers

Stable promotion is forbidden with any:

- missing or duplicate applicable ID;
- `REVALIDATION_REQUIRED`, `REOPENED` or `NOT_IMPLEMENTED_CONFIRMED` status;
- unexplained deferral or future placement;
- missing owner/regression/evidence fingerprint;
- unresolved recurrence root cause;
- unsupported role, desk, viewport, restart or failure path;
- missing macOS or Windows actual-artifact proof;
- TradeInsight entitlement/rights ambiguity;
- protected formula, U.S. Equities or No Execution boundary drift;
- self-inflicted material degradation under supported load.

## 9. Definition of done

v18.5.1 Final Closure is complete only when:

1. all 282 rows are conserved and explicitly dispositioned;
2. every applicable v17/v18 item is `FRESH_PASS` or has a narrowly approved superseded/not-applicable disposition with proof;
3. every confirmed miss and reopened defect is fixed and regression-proven;
4. all 13 orphaned remediation records are closed;
5. zero open, unowned, untested or evidence-stale applicable rows remain;
6. full cross-role, cross-desk, cross-viewport, restart, failure and load matrices pass;
7. actual macOS Apple Silicon and Windows x64 packages pass direct runtime audit;
8. G15 releases the exact certified artifacts without rebuild;
9. G16 records every escape, invalidation and prevention control;
10. the next runner can reconstruct the entire release from repository evidence alone.

## 10. User involvement

Routine implementation, testing, certification and documentation proceed autonomously. User action is required only for unavoidable credentials, paid-service commitments, legal/licensing/data-rights decisions, irreversible external actions or genuinely new material product decisions.
