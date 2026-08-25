# DE.PULSE Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Certified Stable candidate:** `e55d8d25b15cec2ffb0f5411bc358bc40b359cf9`  
**Certified Stable qualified source:** `d7276c3421dd2b4529ac2a987466be3cffa05678`  
**Certified Stable source fingerprint:** `0062f46dea5690d0b3fcd8a9ed3b1f71ebe1522c7dee2cb218e9d36b9e0076ff`  
**Certified Stable build ID:** `v18.9.1-stable-20260821`  
**Live program baseline:** `5bda4fa96612423e79087f1738728670fa002834`  
**Mandatory next release:** `v18.10.0` — **10/10 Future-Proof Final v18 Closure**  
**Parent:** #113 / `ADAPT-V18-FINAL-CLOSURE-10-10-001`  
**Program branch:** `adapt-v18-final-closure-10-10-001`  
**Canonical program plan:** `governance/programs/ADAPT-V18-FINAL-CLOSURE-10-10-001/V18_10_FINAL_CLOSURE_PLAN.md`  
**Machine feature ledger:** `governance/programs/ADAPT-V18-FINAL-CLOSURE-10-10-001/feature-assurance-ledger.json`  
**Parent closure ledger:** `governance/work-slices/ADAPT-V18-FINAL-CLOSURE-10-10-001/closure.json`  
**Future hosted umbrella:** #66 / `ADAPT-HOSTED-SYNC-001` — BLOCKED until v18.10.0 final closure and post-closure audit.

## Current authority

The previous v18 provider-intelligence program (#65/#107) is complete, but before v19 the owner requires one final exhaustive v18 release. Therefore `v18.10.0` is now the mandatory next product release. It is not a feature expansion; it is a full-product audit/certification allowed to fix any genuine implementation/test/data-truth/UI/UX/security/persistence/performance/platform gap it discovers.

The release may be called **10/10 Future-Proof Final v18 Closure** only after all ten assurance tracks are VERIFIED and one immutable candidate completes G0–G16 publication. Documentation alone cannot close a gap.

## Ten assurance tracks

1. #114 / T1 — complete feature / requirement / owner / test traceability.
2. #115 / T2 — exhaustive unit / contract / static / property assurance.
3. #116 / T3 — full functional / integration / end-to-end feature matrix.
4. #117 / T4 — edge / adversarial / failure / data-truth matrix.
5. #118 / T5 — persistence / restart / migration / install / upgrade / recovery.
6. #119 / T6 — security / roles / secrets / rights / negative authorization.
7. #120 / T7 — UI / UX / information architecture / content / accessibility. Every visible item receives KEEP/MOVE/MERGE/REMOVE/RENAME/REDESIGN disposition; technical functionality alone is insufficient.
8. #121 / T8 — performance / load / soak / concurrency / resource safety.
9. #122 / T9 — macOS Apple Silicon + Windows x64 packaged runtime/release/provenance certification plus Chrome/WebKit renderer truth.
10. #123 / T10 — durable regression ownership, executable v19 requirement conservation, GitHub-only portability, zero-gap final certification and v18.10.0 publication.

## 10/10 zero-miss rule

Every shipped v18 feature must ultimately have: requirement provenance; canonical owner; positive functional evidence; applicable unit/contract evidence; edge/negative/failure evidence; persistence/restart evidence where applicable; role/security/provider-right evidence where applicable; UI/UX/IA/content evidence if visible; required platform evidence; durable regression ownership.

The T1 machine ledger starts `DISCOVERY_PENDING`. Do not populate it from memory. Audit current `main`, renderer/navigation, GitHub issues/comments/PRs, Stable/release evidence and canonical contracts. Any `UNOWNED`, `UNTESTED`, `UNKNOWN`, duplicate owner, closed-issue-without-executable-evidence or unexplained carry-forward is a blocker. A genuine implementation miss becomes an explicit corrective under #113 and is fixed/re-qualified before closure.

## Retained architecture

Smart Provider Router v2 is sole routing/admission authority. Direct SEC/EDGAR remains filing/Form 4 authority. Canonical Data Health/freshness/degradation/cache/persistence/subscription/telemetry/reconciliation/lifecycle and identity/session owners remain authoritative. Provider usefulness remains `ADVISORY_ONLY`. U.S. equities processing, GLD/SLV/USO actionable exceptions and No Execution remain permanent. No parallel subsystem may be created just to satisfy closure testing.

## Future-proof requirement

Before first v19 product G1, the 72-row #66 requirement-conservation ledger must be mechanically enforced in existing CI. After v18.10.0, future changes must fail closed if a conserved requirement or durable regression owner disappears.

## Exactly one next action

After the #113 program-registration governance PR is Fast/Qualified and merged, fetch fresh `main` and execute **#114 / T1 only**: exhaustive shipped-feature/requirement/owner/test discovery and machine-ledger population. Do not start T2–T10 or v19 product work before T1 inventory is frozen and its discovered corrective gaps are explicitly dispositioned.

## Resume rule

1. Fetch live `main` and active/open PRs/issues first.
2. Read this file, `governance/current-state.json`, `governance/ROADMAP.md`, the v18.10 plan/feature ledger, #113 and child issues #114–#123.
3. GitHub objects/executable evidence outrank roadmap prose and chat memory.
4. If program registration is not yet merged, finish that exact-head governance delivery first; otherwise resume the current T1–T10 child named by machine state.
5. Never infer a feature is covered because an old issue is closed; require executable evidence.
6. Never hide a newly discovered implementation miss inside a closure-only track.
7. Preserve Smart Provider Router v2, Data Health/freshness/cache/persistence/subscription/telemetry/reconciliation/lifecycle, identity/session, direct SEC/EDGAR, U.S. equities, GLD/SLV/USO and No Execution.
8. Use G0–G16 only: canonical Fast exact-head PASS -> identical-head Qualified PASS -> expected-head merge; final public v18.10.0 publication only after T1–T10 VERIFIED.
9. A new ChatGPT account, Codex or Claude must resume from GitHub source-of-truth. No old chat handoff upload is required.
