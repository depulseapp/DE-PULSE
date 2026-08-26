# CURRENT Adaptive Build Process

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.10.0` — immutable  
**Post-v18 audit:** #145 / `adapt-post-v18-overlap-rebaseline-001` — PASS candidate pending exact-head qualification/merge.

The permanent execution loop remains source-driven and exact-head: LOOKUP -> COMPARE -> CLASSIFY -> DECIDE -> UPDATE -> Fast -> Qualified -> G11–G16 when a release is actually being produced. v18.10.0 remains immutable and is not reopened by this audit.

## Rebaselined process rules

1. **Start from live source, not roadmap labels.** Before G1 implementation, compare the requirement to current canonical owners and classify it `REUSE`, `EXTEND_EXISTING`, `CONSOLIDATE`, `JUSTIFIED_RESIDUAL`, or `NOT_APPLICABLE`.
2. **Bind evidence while coding.** A requirement is not development-complete until its canonical owner, consumer, positive/adverse evidence and persistence/security/UI applicability are recorded. Do not postpone evidence ownership until major closure.
3. **Classify assurance gaps correctly.** Distinguish product behavior gaps from test/evidence gaps, ownership-binding gaps and legitimate N/A. Do not create product machinery merely to satisfy a bookkeeping gap.
4. **One coherent candidate per coherent correction.** Prepare related changes before PR creation/advancement where practical. A changed candidate still requires fresh exact-head Fast; efficiency comes from reducing needless candidate SHAs, not suppressing validation.
5. **Requirements are not releases.** `HOST-001..072` remain traceability rows. Implement dependency-coherent bands/slices; preserve row-level closure without forcing one PR/release per row.
6. **Reuse canonical runtime owners.** Smart Provider Router v2, Data Health/freshness, cache/persistence, multi-feed subscription, telemetry/reconciliation/lifecycle, identity/session, Research/Discovery, decision/outcome lineage and direct SEC/EDGAR remain authoritative.
7. **Hosted threat modeling begins before code exit.** Every hosted slice asks how tenant escape, cache/coalescing leakage, long-lived stream revocation, secret leakage, mixed-client downgrade, noisy-neighbor pressure and backup/restore isolation fail closed.
8. **Point-in-time before learning.** Historical/adaptive evaluations must reconstruct what was actually knowable at the decision timestamp; later revisions cannot leak backward.
9. **Measure usefulness, not just availability.** Where a slice affects intelligence, name the future outcome/calibration signal and reuse canonical decision/outcome history rather than creating a parallel analytics store.
10. **Trader-grade adverse state.** Explicitly classify exchange halt/LULD/volatility pause/resume separately from provider outage, stale quote and ordinary closed-session state.
11. **Historical CI is conserved, not multiplied.** Frozen v18 T1–T10 remain durable baseline evidence. Future CI should use a fail-closed conservation concept and invoke deeper historical logic only when materially affected, while keeping changed product/regression tests fully enforced.
12. **Post-merge continuity is not PR Fast.** Main-push continuity/branch hygiene should be named/documented as such so operators do not mistake it for another full PR validation run.

## Exactly one next action

Qualify and merge #145. Only then may a v19 G1 band be reserved from the rebaselined process.
