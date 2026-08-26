# Adaptive CI Rebaseline — Post-v18.10

**Status:** governing addendum to `ADAPTIVE_CI_OPERATING_CONTRACT.md` and `governance/CI-EFFICIENCY-CONTRACT.md` for v19+  
**Audit:** #145  
**Stable baseline:** v18.10.0

## Purpose

Preserve the exact-source quality achieved by v18 T1–T10 without allowing historical assurance machinery or candidate-SHA churn to scale faster than product value.

## Frozen historical conservation

v18 T1–T10 are immutable closure evidence. For v19+ routine PRs, CI should converge toward one fail-closed **v18 conservation** responsibility that verifies:
- the frozen 180-responsibility ledger remains intact;
- durable regression owners still exist or have proven retired-test equivalence;
- permanent product/architecture boundaries remain intact;
- canonical owners have not been silently replaced;
- immutable Stable/release provenance remains unchanged.

Deeper T1–T10 structural logic may remain separately executable and must run when its governing artifacts/owners/contracts are materially changed. This optimization may not remove or skip affected Go/renderer/browser/native/security/persistence tests.

## Candidate-SHA rule

A changed candidate always requires fresh exact-head Fast and, when applicable, Qualified. Efficiency is achieved by batching coherent edits and minimizing incomplete metadata/correction commits while a PR is open—not by reusing stale validation.

## Requirements are not CI events

The v19 `HOST-001..HOST-072` ledger is row-level traceability. It must not be interpreted as 72 required branches, PRs, Qualified runs or Stable releases. Group dependency-coherent requirements into the smallest sensible implementation/release bands while retaining exact row evidence.

## Main-push naming and documentation

The workflow named `DE.PULSE | CI Fast` currently also contains `main`-push continuity and branch-hygiene jobs while the actual PR Fast validation job is skipped. Documentation and future workflow presentation should distinguish **PR Fast** from **Main Continuity / Branch Hygiene** so operators do not misread a post-merge sentinel as another full validation run.

The continuity sentinel is useful and should be retained. Update documentation rather than deleting a safety control solely to match stale prose.

## Historical gate growth guard

Before adding a permanent CI step, ask:
1. Is this new assurance or a version-specific duplicate of an existing invariant?
2. Can it extend a canonical conservation/health/security/utility gate?
3. Can Impact Planner select it only when relevant?
4. Does it materially improve fault detection compared with its ongoing runner/maintenance cost?
5. What older check becomes redundant or subordinate?

A new major version must not automatically add another always-on T1–T10-style family to every future PR.

## Main protection

GitHub currently reports `main` as unprotected for this private user-owned repository. Internal exact-head status/provenance controls remain mandatory. When repository plan/capabilities permit, require canonical PR/status checks and prohibit accidental destructive direct changes.

## Permanent rule

**Reduce candidate churn and duplicated historical orchestration, never evidence for a changed risk surface. CI conservation grows by stronger shared invariants, not by permanently stacking every prior release's gate family.**
