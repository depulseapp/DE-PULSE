# DE.PULSE v18.x — Adaptive Recovery, Hardening & Final Closure Program

Status: **OPEN / ACTIVE v18.x RECOVERY PROGRAM**  
Authority: `governance/ROADMAP.md` + the current v18 reconciliation ledger  
Process: **G0–G16 only; no G17+**

## 1. Purpose

This program governs v18.5.1, v18.6, v18.7 and any further evidence-required v18.x slices. v18.5.1 is the audit/containment and urgent-recovery entry slice; it is not automatically the final closure build. The final v18 closure version is selected only when every applicable v17/v18 commitment is implemented, current behavior is proven, hardening thresholds are met, and the exact macOS Apple Silicon and Windows x64 packages are ready for certification.

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

## 4. Frozen program inventory

The authoritative ledger currently contains 282 rows:

- 48 inherited approved requirements;
- 20 frozen v17 major items;
- 8 v18 major workstreams;
- 169 v18.0–v18.5 release entries;
- 13 orphaned functionality/utility remediations;
- 7 recovered conversational commitments;
- 7 confirmed implementation-miss groups;
- 7 escaped-defect groups;
- 4 approved future-only mature workstreams.

The four mature v20 workstreams are not current implementation misses. Their v18/v19 foundations remain applicable where already approved.

No additional row may be inserted only into a slice, issue or conversation. It must first enter the ledger with origin, owner, placement, affected behavior and required evidence.

## 5. Adaptive release-train model

The release number follows the evidence; the evidence does not get forced into a preselected release number.

### v18.5.1 — Audit, containment and urgent recovery

Mandatory outcomes:

- establish and enforce the complete reconciliation ledger;
- prevent new silent slicing or inherited PASS;
- preserve every reported defect and missing implementation as a blocker or explicit next slice;
- fix only urgent/high-confidence items that can be safely completed without hiding larger dependencies;
- finish with an evidence-backed handoff to the next v18.x slice if open applicable work remains.

v18.5.1 may close as a recovery release while the v18 major remains open.

### v18.6+ — Evidence-selected implementation slices

Before each new minor release, run the adaptive loop:

`Observe → Reconcile → Prioritize → Slice → Build → Validate → Measure → Learn → Replan`

Prioritization considers:

- user trust and real-money impact;
- defect recurrence and implementation absence;
- architecture/dependency readiness;
- security, provider entitlement and data rights;
- performance/freshness risk;
- coupling and safest build order;
- evidence that can be reused versus must rerun;
- expected decision-support value and provider/runtime cost.

Likely clusters may include user-trust UI/state recovery, provider/utility architecture, TradeInsight, documentation/dependency readiness, consolidation and final hardening. These are planning candidates, not pre-frozen promises. G0–G3 selects the smallest coherent next slice from the open ledger.

### Provisional starting sequence

This is the recommended starting plan, subject to G0–G3 replanning after each slice:

| Candidate | Primary objective | Candidate contents | Why this grouping |
|---|---|---|---|
| v18.5.1 | Trust recovery + anti-miss control plane | ledger/gate; version copy; symbol removal/highlighting; hover/render stability; scroll/focus continuity; Research top-area; exact next-slice placement | Direct user pain, repeated defects and high-confidence boundaries should be corrected first. |
| v18.6 | Canonical utility/ownership recovery | shared Scanner/Radar acquisition; Session Intelligence Coordinator; Market Activity demotion; legacy route redirects; role-aware docs/impact manifest; dependency readiness register | These items share architecture, ownership, surface and governance dependencies. |
| v18.7 | Provider intelligence + remaining consolidation | TradeInsight SHADOW/SECONDARY; remaining Dashboard/Market Intelligence/Desk/Catalyst/Maintenance remediation; provider-rights/readiness validation | External/provider work and larger consolidation need isolated rights, performance and regression qualification. |
| v18.8+ | Major closure candidate or additional hardening | full 283-row evidence convergence; cross-role/failure/load/native hardening; any newly revealed repair | Designate closure only if zero-gap readiness is already true; otherwise create another v18.x slice. |

This table is not immutable scope. Each candidate becomes binding only when its own G1 closes. Findings may pull a dependency-compatible blocker earlier, split an unsafe slice or create another v18.x release.

### Final v18.x closure release

The final closure version may be v18.6, v18.7, v18.8 or later. It is designated only when:

- no confirmed implementation miss remains;
- no reopened defect remains;
- all 13 orphaned remediations are final;
- all applicable v17/v18 rows have current evidence;
- cross-cutting hardening is clean;
- the remaining work fits one immutable RC without unsafe compression.

If new evidence appears, the program creates another v18.x recovery/hardening slice rather than forcing closure.

## 6. Adaptive build waves

### Wave 0 — Reconstruct and freeze (G0–G3)

- Freeze the exact v18.5.0 Stable source/package baseline and v18.5.1 branch.
- Validate all 283 IDs, origins and counts.
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
- stable ticker/row hover with no table or screen blinking during repeated live/SSE updates;
- scroll/focus/selection preservation through live/SSE and save/action rerenders;
- Research target/freshness/action layout and recovery behavior.

Required matrices include all desks, first/middle/final membership, ticker cards and tables, pointer hover dwell through repeated live/SSE updates, reload/restart, desktop/tablet/narrow widths, reduced-motion behavior and keyboard focus.

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

## 7. Test and evidence model

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

## 8. Efficient execution without weakening closure

Use the CI/CD model inside G0–G16:

1. run inventory, syntax, scope and cheap static checks first;
2. run independent backend, renderer, security, persistence and documentation lanes in parallel;
3. serialize shared-state migrations and protected-formula checks;
4. run focused affected tests before expensive full suites;
5. merge evidence at G10;
6. build native packages once from the immutable RC;
7. rerun only invalidated evidence until RC; after RC, any source change creates a new candidate.

Parallel execution may reduce elapsed time. It may never split requirement authority or create multiple competing closure ledgers.

## 9. Promotion blockers

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

## 10. Definition of done

The v18 major Final Closure is complete only when:

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

## 11. User involvement

Routine implementation, testing, certification and documentation proceed autonomously. User action is required only for unavoidable credentials, paid-service commitments, legal/licensing/data-rights decisions, irreversible external actions or genuinely new material product decisions.
