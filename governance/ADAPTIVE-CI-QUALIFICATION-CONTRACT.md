# DE.PULSE — Adaptive CI Cost & Qualification Contract

**Status:** APPROVED / PERMANENT / GOVERNING  
**Effective:** 2026-08-17  
**Applies to:** Adaptive Roadmap, Adaptive Build Plan, Adaptive Build Process, Adaptive Delivery Process, G0–G16 execution, GitHub Actions, native packaging/certification, AIPLC checkpoints and release evidence  
**Purpose:** Reduce redundant CI/CD execution and cost without reducing scope, test depth, behavioral proof, platform coverage, certification rigor or release quality.

This contract refines execution timing only. It does **not** add a G17+, remove a gate, weaken any gate, or permit evidence reuse across a materially different source fingerprint.

---

## 1. Permanent Quality Principle

DE.PULSE optimizes **when required tests execute**, never **which required tests exist**.

> **CI efficiency may eliminate redundant execution, never required evidence. No cost optimization may weaken scope, behavioral test depth, security, data/rights validation, performance/capacity/stability testing, professional acceptance, supported-platform certification, provenance, or actual packaged-runtime proof.**

Cost is an engineering constraint; it is never a justification for false PASS, skipped evidence, reduced platform coverage or weaker release assurance.

---

## 2. Permanent Execution Model

Use this lifecycle for ordinary feature/fix development:

**DEVELOP → ACCUMULATE COHERENT FIXES → DEVELOPMENT CHECKPOINT → G0–G5 FAST QUALIFICATION → BATCH FIX/RECHECKPOINT IF NEEDED → G6–G10 QUALIFICATION → IMMUTABLE RC → G11–G15 FULL/NATIVE CERTIFICATION → G16**

The default is **not**:

`small fix → full workflow → small fix → full workflow → copy/CSS/test change → full workflow`.

---

## 3. Four Adaptive Layers

### Adaptive Roadmap
Roadmap planning must include CI/runtime cost and qualification efficiency as part of engineering sustainability. Major/minor release planning should prefer a small canonical workflow set over accumulating release-specific orchestration indefinitely.

Roadmap quality protection:
- no required gate disappears because it is expensive;
- native macOS Apple Silicon and Windows x64 certification remains mandatory where already required;
- expensive evidence moves to the appropriate checkpoint/RC stage rather than being deleted;
- historical workflows may be retained as provenance but must not remain active against unrelated branches.

### Adaptive Build Plan
Each build plan must identify:
- coherent development batches/workstreams;
- the intended Development Checkpoint(s);
- FAST qualification scope;
- G6–G10 qualification scope;
- immutable RC boundary;
- G11–G15 release/native certification scope;
- which lanes may run safely in parallel;
- which evidence can be reused only when exact fingerprint/environment/test-definition equivalence is proven.

Development fixes should include their unit/behavior/regression tests during implementation, but Actions qualification should normally wait for the checkpoint.

### Adaptive Build Process
The build process uses three principal execution modes plus G16:

1. **DEV** — normal implementation commits. No expensive Actions by default.
2. **QUALIFY** — deliberate checkpoint. Run G0–G10 with cheap/fail-fast lanes before medium/expensive lanes.
3. **RELEASE** — immutable RC. Run G11–G15 full certification/native packaging/runtime/provenance.
4. **G16** — retrospective, handoff, cleanup and durable learning.

Rules:
- ordinary source/copy/CSS/test/documentation commits must not automatically trigger medium/full/native qualification;
- a qualification checkpoint must bind an exact source SHA/fingerprint and checkpoint ID;
- G0–G5 run first and fail fast;
- G6/G7/G8/G9 run only after FAST is green and should run in bounded parallel when independent;
- G10 binds all qualification evidence to the exact checkpoint source;
- if qualification fails, repair failures as a coherent batch and create a new checkpoint instead of rerunning the entire pipeline after every tiny repair;
- rerun only failed jobs when GitHub supports it and source/fingerprint/test definitions are unchanged;
- stale/superseded runs should be cancelled;
- CI-generated evidence/materialization must not recursively push to a watched branch and trigger itself;
- use path/checkpoint filters so unrelated commits do not launch qualification;
- use dependency caches and reusable workflows when safe, but cache use may never obscure a clean reproducibility problem;
- never reuse PASS evidence after a materially different source/tooling/test/environment fingerprint.

### Adaptive Delivery Process
Delivery begins only after a green G10 checkpoint produces an immutable RC candidate.

G11–G15 remain exhaustive:
- G11 freezes exact RC identity;
- G12 runs full certification/regression/security/performance/data/UI/professional acceptance as applicable;
- G13 packages required native artifacts/provenance;
- G14 audits the actual packaged runtime on required platforms;
- G15 proves promotion, rollback, reproducibility and provenance truth.

A source/tooling change after RC invalidates the affected certification evidence and requires the appropriate requalification. There is no “small enough to skip certification” exception when the change can affect the certified artifact.

---

## 4. Development Checkpoint Contract

A Development Checkpoint is a deliberate qualification request, not an ordinary commit.

Minimum checkpoint data:
- release/version;
- checkpoint ID;
- exact source SHA/fingerprint;
- qualification requested = true;
- current scope/freeze identity;
- optional note describing the coherent batch.

The checkpoint file/record is allowed to be a tiny trigger commit after the source batch. Its only purpose is to bind qualification to the exact source that preceded it.

A checkpoint must fail closed when:
- the source SHA is missing/invalid;
- it is not an ancestor of the checkpoint commit;
- release identity does not match;
- frozen scope cannot be verified;
- the requested evidence would be bound to a different fingerprint.

---

## 5. Cheap-First Qualification

Preferred order:

### FAST / Low Cost
- exact identity/baseline/scope;
- syntax/static checks;
- source/architecture/data utility governance;
- `go vet`;
- focused affected tests;
- full deterministic Go development-exit suite;
- focused renderer/Chromium behavior tests required for the changed slice.

### MEDIUM / Expensive only after FAST green
- PostgreSQL/DB integration;
- cross-module HTTP workflows;
- security/adversarial/RBAC/rights tests;
- performance/capacity/stability;
- race detector/randomized-order suites;
- broad responsive/UI/professional acceptance;
- failure injection/soak where required.

### RELEASE only after G10 green and RC frozen
- full certification;
- native macOS Apple Silicon package/runtime audit;
- native Windows x64 package/runtime audit;
- provenance/signing/manifest/rollback/reproducibility;
- promotion assurance.

---

## 6. Parallelism Contract

Parallelize only independent work whose evidence consumes the same immutable checkpoint/RC source.

Examples after FAST green:
- G6 integration;
- G7 data/security/adaptive intelligence;
- G8 performance/capacity/stability;
- G9 cross-module/UI/UX.

G10 joins them and fails unless every required lane passes.

Do not parallelize when one lane creates or mutates inputs required by another. Qualification workflows should prefer read-only/non-mutating execution.

---

## 7. Recursive Trigger / Workflow Hygiene

Permanent protections:
- qualification workflows should be read-only against source branches wherever feasible;
- generated evidence belongs in artifacts or deliberately committed later outside the watched qualification trigger;
- a workflow that changes source/identity/scope cannot certify the pre-change SHA as if nothing changed;
- path filters/checkpoint triggers prevent ordinary development commits from launching full qualification;
- `concurrency.cancel-in-progress` is a safety net, not the primary cost-control mechanism;
- obsolete historical workflow files may stay for provenance only when their branch/path triggers cannot affect active development;
- converge toward one logical Build Coordinator and reusable lane definitions rather than one accumulating orchestration system per minor release.

---

## 8. Evidence Reuse / Rerun Rules

Evidence may be reused or a failed job rerun without repeating successful jobs only when all relevant items are equivalent:
- source SHA/fingerprint;
- test implementation/version;
- workflow/gate definition;
- toolchain/environment;
- input/config/provider fixture contract;
- required platform/runtime assumptions.

If any materially relevant element changed, rerun the affected evidence.

No old PASS may be copied forward merely to save Actions minutes.

---

## 9. AIPLC Integration

AIPLC should align with meaningful development/qualification checkpoints.

Do not generate large redundant AIPLC reports for every tiny mechanical commit. A coherent batch/checkpoint should capture:
- what changed;
- root causes addressed;
- reusable prevention;
- tests added/changed;
- cost/efficiency learning;
- whether qualification revealed new product/architecture/test gaps.

Mechanically identical reruns with no new evidence may use **NO NEW LEARNING / EVIDENCE EQUIVALENT** when supported by the Adaptive Operating Contract.

---

## 10. v18.5.1 Immediate Adoption

For `v18.5.1-development`:
- ordinary pushes no longer trigger the recovery qualification chain;
- `release/v18.5.1/QUALIFICATION_CHECKPOINT.json` is the deliberate checkpoint trigger;
- checkpoint qualification is bound to an exact source SHA;
- qualification does not push generated source/evidence back into the watched branch;
- G0–G5 run before G6–G10;
- G6–G9 run only after FAST success;
- G10 binds the qualification result to the checkpoint identity/source SHA;
- G11–G15 remain a later immutable-RC release phase and are not run after every development repair.

This change is process hardening only. It does not expand the frozen v18.5.1 product scope or close any implementation defect by documentation alone.

---

## 11. Success Measures

Track over time:
- Actions minutes/cost per delivered release;
- number of workflow invocations per development checkpoint;
- cancelled/superseded runner minutes;
- duplicate tests avoided without loss of evidence;
- time from development-complete to G10;
- failed-lane-only rerun rate;
- release escape defects;
- defects caught by FAST vs G6–G10 vs G11–G15;
- native certification failures;
- false PASS / stale-evidence incidents (target: zero).

Efficiency is successful only when cost/runtime falls **without deterioration in release escape rate, evidence quality, supported-platform truth or user-visible reliability**.

---

## 12. Permanent Acceptance Standard

The operating model is compliant only when:
- normal development can proceed without repeatedly spending Actions minutes on the entire qualification chain;
- deliberate checkpoint qualification is exact-SHA/fingerprint bound;
- cheap failures prevent expensive downstream execution;
- required independent lanes parallelize safely;
- recursive/self-triggering qualification is prevented;
- immutable RC certification remains exhaustive;
- actual macOS Apple Silicon and Windows x64 artifact/runtime evidence is still required before applicable Stable promotion;
- no cost optimization lowers the DE.PULSE 10/10 quality standard.
