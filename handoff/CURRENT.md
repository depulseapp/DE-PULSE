# DE.PULSE Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Certified Stable candidate:** `e55d8d25b15cec2ffb0f5411bc358bc40b359cf9`  
**Certified Stable qualified source:** `d7276c3421dd2b4529ac2a987466be3cffa05678`  
**Certified Stable source fingerprint:** `0062f46dea5690d0b3fcd8a9ed3b1f71ebe1522c7dee2cb218e9d36b9e0076ff`  
**Certified Stable build ID:** `v18.9.1-stable-20260821`  
**Current closure branch:** `adapt-v18-final-closure-10-10-001`  
**Mandatory next release:** `v18.10.0` — **10/10 Future-Proof Final v18 Closure**  
**Parent:** #113 / `ADAPT-V18-FINAL-CLOSURE-10-10-001`  
**Completed T1:** #114 / PR #127 / merge `b434886751fde1804a4906c9cac41dcce4584834`  
**Completed T2:** #115 / PR #128 / merge `bd26f91964235a8fdb390184684f4fdd216eb1e6`  
**Completed T3:** #116 / PR #129 / merge `cf0cc9aa877188775d5c36cb805e6c8935d776bb`  
**Completed T4:** #117 / PR #130 / merge `fadd9247f63779195138f8d79cbaa0f304b16e61`  
**Active child:** #118 / T5 — persistence / restart / migration / install / upgrade / recovery  
**Recorded companion:** #119 / T6 — NOT STARTED  
**Frozen T1 responsibility count:** 180  
**T1 frozen assurance manifest:** `governance/programs/ADAPT-V18-FINAL-CLOSURE-10-10-001/feature-assurance-ledger-freeze.json`  
**T2 assurance:** `governance/programs/ADAPT-V18-FINAL-CLOSURE-10-10-001/T2_UNIT_CONTRACT_ASSURANCE.json`  
**T3 assurance:** `governance/programs/ADAPT-V18-FINAL-CLOSURE-10-10-001/T3_FUNCTIONAL_ASSURANCE.json`  
**T4 assurance:** `governance/programs/ADAPT-V18-FINAL-CLOSURE-10-10-001/T4_EDGE_FAILURE_DATA_TRUTH_ASSURANCE.json`  
**T5 assurance:** `governance/programs/ADAPT-V18-FINAL-CLOSURE-10-10-001/T5_PERSISTENCE_LIFECYCLE_ASSURANCE.json`  
**Parent closure ledger:** `governance/work-slices/ADAPT-V18-FINAL-CLOSURE-10-10-001/closure.json`  
**Future hosted umbrella:** #66 / `ADAPT-HOSTED-SYNC-001` — BLOCKED.

## Current authority

T1-T4 are durably closed. T4 final head `cf8bab0650e6e1d0d2b099ebaf2ec18757737368` passed exact-head Fast #1034/run `32869228199` and identical-head Qualified #209/run `32869561799`; PR #130 merged with expected-head guard as `fadd9247f63779195138f8d79cbaa0f304b16e61`. The T4 gate proved 180/180 adverse/failure/data-truth coverage with zero uncovered responsibilities.

T5/#118 is now the sole active closure action. It must certify durable state and lifecycle behavior across SQLite integrity, canonical persistence, migrations, fresh/warm restart, settings/desks/watchlists/research/config retention, backup/restore where supported, upgrade/rollback/recovery, installer/profile preservation and interrupted-write safety. T6/#119 remains NOT_STARTED.

## Exactly one next action

**Continue T5/#118 only:** reconstruct the frozen 180-responsibility inventory, resolve T5 applicability for every row, require deterministic persistence/restart/recovery evidence for every applicable row, and record the exact uncovered set from executable evidence. Do not start T6 until T5 is durably closed or GitHub machine state explicitly advances it.

## T5 closure discipline

- Preserve the immutable T1 inventory; do not create a smaller hand-picked persistence list.
- Every row must have a T5 expectation and an explicit applicability result.
- Applicable rows must use existing canonical persistence/cache/identity/session/workspace owners; no parallel persistence or recovery subsystem may be introduced.
- Prove SQLite integrity/migrations, fresh and warm restart, state retention, upgrade/rollback/recovery, interrupted-write/atomicity behavior and profile preservation where applicable.
- Read-only surfaces may consume restored canonical state but may not become independent persistence owners.
- Actual packaged native install/upgrade proof remains a T9 responsibility; T5 may require release/profile preservation contracts but must not claim T9 complete.
- If a genuine implementation miss appears, create a named corrective under #113; do not hide it in assurance metadata.
- T5 remains unverified until every applicable row is covered, all declared gaps are closed, exact-head Fast and identical-head Qualified pass, and the expected-head merge is durable.

## Retained architecture

Smart Provider Router v2 remains sole routing/admission authority. Direct SEC/EDGAR remains filing/Form 4 authority. Canonical Data Health/freshness/degradation/cache/persistence/subscription/telemetry/reconciliation/lifecycle and identity/session owners remain authoritative. Provider usefulness remains advisory only. U.S. equities processing, GLD/SLV/USO actionable exceptions, governed provider lifecycle and No Execution remain permanent. No parallel subsystem may be created merely to satisfy closure testing.

## Resume rule

1. Fetch live `main`, `adapt-v18-final-closure-10-10-001`, current closure PR if any, #113 and active child #118 first.
2. Read this file, `governance/current-state.json`, the v18.10 plan, frozen T1 manifest/reconciliation, T2-T4 assurance and T5 assurance.
3. Confirm T1-T4 durable evidence before relying on the frozen 180-responsibility inventory.
4. Resume **T5/#118 only** from current exact branch head. T6/#119 remains NOT_STARTED.
5. GitHub objects and executable evidence outrank this prose and all chat history.
6. Use G0-G16 only: Fast exact-head PASS -> identical-head Qualified PASS -> expected-head merge. Final v18.10.0 publication remains prohibited until T1-T10 are VERIFIED.
7. A new ChatGPT account, Codex or Claude must be able to resume from GitHub alone.
