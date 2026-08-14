# DE.PULSE — Adaptive CI Operating Contract

Status: **Permanent governing contract**  
Applies to: Adaptive Build Plan, Adaptive Build Process, Adaptive Delivery Process, Build Resume Protocol  
Effective: **v18.2 and all later releases**

## Purpose

DE.PULSE CI/CD behaves as an adaptive engineering system:

**observe → classify → diagnose → generalize → prevent → impact-map → execute → measure → learn.**

A failure is not fully resolved when it merely turns green; useful root-cause learning must reduce recurrence where practical.

**G0–G16 is the permanent top-level release model. No G17+.** New assurance responsibilities belong in checkpoints, sub-stages, shards, bounded parallel lanes, or strengthened ownership within an existing gate.

---

## 1. One logical Build Coordinator

Every release has one authoritative orchestration/dependency owner.

Conceptual dependency flow:

**G0–G4 → G5–G10 → G11 → G12 → G13/G14 macOS + Windows → G15 → G16**

Separate workflow files may exist only as subordinate/reusable lanes. They must not independently own conflicting release mutations or authoritative gate conclusions.

Use path/event filters, concurrency/supersession and idempotent mutations to prevent self-trigger loops and duplicate expensive work.

---

## 2. Mandatory failure classification

Every failed/terminated lane is classified before candidate health changes:
- `PRODUCT_FAIL` — reproducible product/contract defect;
- `GATE_TEST_FAIL` — invalid/stale/brittle gate/test logic preventing trustworthy evidence;
- `CI_HARNESS_FAIL` — workflow/script/path/tooling failure unrelated to product correctness;
- `INFRA_FAIL` — runner/browser/network/toolchain/capacity/external infrastructure failure;
- `EXPECTED_NOOP` — successful idempotent no-change behavior;
- `SUPERSEDED` — newer candidate replaced the run.

Only reproducible `PRODUCT_FAIL` proves the candidate itself is defective. Other classes may still block trustworthy evidence, but they are not mislabeled as product defects.

---

## 3. Preventative-learning preflight

Before expensive qualification, check known recurrence patterns including:
- generated debris/temp/build outputs entering source;
- platform-specific shell assumptions where portable tooling is appropriate;
- unguarded no-op commits/mutations;
- workflow self-trigger loops;
- checkpoint-only commits waking full certification/native work;
- duplicate ownership of expensive tests/gates;
- stale fingerprint/RC evidence;
- insufficient history/provenance;
- hard-coded release literals where canonical identity should be derived;
- tests mutating frozen source;
- package evidence reused after package/source identity changes;
- missing functionality-utility registry coverage;
- new capability/data/provider path without owner/consumer/reuse/correlation/materiality/rights disposition;
- unjustified new tab/surface;
- duplicate acquisition/computation/in-flight work;
- obsolete owners/surfaces/jobs retained after consolidation.

Recurring failure patterns should become preventative controls/tests where useful.

---

## 4. Build State Ledger / Checkpoint

The ledger is a **derived/reconcilable representation**, never an optimistic source of truth.

It reconciles against:
- actual branch HEAD;
- incoming Stable/release identity;
- candidate/source fingerprint;
- CI runs/jobs;
- artifacts;
- immutable RC;
- native package/runtime evidence;
- promotion/release objects.

Required truth includes:
- release/channel/incoming Stable;
- active branch and observed current HEAD;
- candidate source commit/fingerprint and state;
- G0–G16 status using canonical states;
- evidence reference for every PASS;
- reuse eligibility;
- blocker/failure classification;
- last trustworthy PASS;
- earliest required resume responsibility;
- exactly one next action;
- macOS/Windows package/hash/G14 state;
- G15/G16/user-delivery state.

If ledger and GitHub disagree, **correct the ledger**.

---

## 5. Metadata isolation / idempotency

Checkpoint/evidence metadata remains fingerprint-excluded where defined.

Metadata-only commits must not trigger full product qualification/native packaging unless an evidence-relevant dependency changed.

Mutating steps must be idempotent. No-change is `EXPECTED_NOOP`, not failure.

---

## 6. Delta-first evidence reuse / partial resume

Before rerunning expensive work:

**diff → owner → dependency blast radius → affected responsibilities → reusable evidence → smallest rerun set.**

Reuse exact-source/artifact evidence only when relevant fingerprint, dependency/input contract, test definition and environment assumptions remain equivalent.

Examples:
- failed Windows G14 does not invalidate unchanged G12 or successful macOS G14;
- harness/infra repair reruns only affected lane plus any necessary repair-validation;
- changed surface/engine/provider/data ownership revalidates affected functionality-utility/G9/G10 evidence;
- checkpoint-only metadata does not invalidate immutable RC evidence.

---

## 7. Adaptive decomposition / bounded parallelism

Heavy CI responsibilities may be split when that improves fault isolation, resumability, evidence reuse or speed.

Every lane has:
- canonical owner;
- dependencies/source/artifact inputs;
- PASS/FAIL/BLOCKED criteria;
- evidence output;
- downstream consumers;
- invalidation/reuse rules.

Parallelism is bounded by total runner minutes, CPU/memory, browser load, artifacts, provider/API/model calls and DB pressure.

Share setup, fixtures, cached/canonical data and evidence packages so sharding does not multiply work.

After failure, rerun the smallest affected/dependent set.

Canonical details: `adaptive-governance/ADAPTIVE_WORK_DECOMPOSITION_CONTRACT.md`.

---

## 8. AIPLC in CI

Meaningful CI/build events feed the Delta AIPLC.

CI should provide concise evidence for:
- what changed/was exercised;
- failure classification;
- root cause;
- fix/disposition;
- prevention added;
- affected coverage/inheritance state;
- metrics/calls/reruns avoided where useful;
- next action.

Mechanically identical reruns may record `NO NEW LEARNING / EVIDENCE EQUIVALENT`.

G10 performs full coverage reconciliation; G16 aggregates deep release learning.

---

## 9. Provider/data/model efficiency

Prefer deterministic fixtures/replay/canonical cached evidence when live behavior is not under test.

Live-provider tests are bounded to cases requiring actual routing/fallback/freshness/entitlement proof.

Equivalent lanes reuse acquisition rather than refetching.

AI/LLM regression prefers canonical evidence packages/fingerprints; true-model calls are bounded to behavior that cannot be certified deterministically. Materially identical reusable synthesis should not be rerun per user/test without need.

---

## 10. Native and user-delivery invariant

Completed native delivery requires:
1. required G13 packages from certified RC;
2. macOS Apple Silicon G14 PASS;
3. Windows x64 G14 PASS;
4. G15 PASS;
5. permanent release assets;
6. hashes/provenance;
7. user-facing handoff surfacing required runnable assets/status;
8. `userDelivery = DELIVERED` at closure.

GitHub packaging without user-facing asset handoff is `READY`, not `DELIVERED`.

Preserve an unchanged platform PASS when exact RC/package identity is unchanged.

---

## 11. Governance-to-Implementation CI closure

CI proves governance is executable, not merely documented.

Requirements follow:

**Governed → Implemented → Enforced → Evidenced → Delivered → Learned.**

CI cannot turn `Governed` into PASS without implementation/evidence.

Before affected evidence is trusted:
- overlapping authoritative workflows are consolidated/subordinated;
- Build State Ledger is reconciled from actual GitHub truth;
- canonical naming applies to active workflow/job/artifact/checkpoint/gate/capability identities;
- current-release blockers/process-hardening obligations have owners and closure evidence.

Canonical details: `adaptive-governance/GOVERNANCE_IMPLEMENTATION_CLOSURE_CONTRACT.md`.

---

## 12. G16 adaptive CI closure

G16 reviews:
- product vs test/harness/infra/no-op/superseded incidents;
- root cause/prevention status;
- duplicated/obsolete workflows/checks;
- evidence reuse/reruns avoided;
- functionality-utility alignment;
- provider/model/resource efficiency;
- branch/workflow hygiene;
- native/user-delivery closure.

Redundant CI machinery is removed/consolidated when it adds ceremony without assurance.

---

## 13. Canonical concise status

Normal release reporting:

**Code | G0–G12 | macOS G13/G14 | Windows G13/G14 | G15 | G16 | User Delivery**

---

## Permanent rule

**DE.PULSE CI learns from failures, protects exact-source evidence, reruns only the smallest affected set, uses bounded parallelism, avoids duplicate provider/model/runtime work, and remains inside permanent G0–G16. More jobs or more gates are never treated as inherently better.**
