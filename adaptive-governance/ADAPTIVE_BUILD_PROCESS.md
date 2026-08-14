# DE.PULSE — Adaptive Build Process

## Permanent Resume Phase

A **Resume Reconciliation** phase is mandatory whenever work starts after interruption, handoff, new conversation, runner interruption, or uncertain prior execution state. It belongs inside the existing G0–G16 process and is not a new gate.

Resume Reconciliation performs:

**Detect active release → verify incoming Stable → read checkpoint → read actual branch SHA/release identity → inspect CI/artifacts → compare fingerprint/immutable RC identity → determine last trustworthy PASS → resume next required step.**

### Evidence reuse

Reuse a PASS only when its evidence is bound to the unchanged relevant source fingerprint or immutable RC/artifact identity. Do not infer PASS from chat text, a checkpoint label alone, or an older build with similar code.

### Invalidation

When source changes, determine the earliest affected G0–G16 gate and mark that gate and dependent downstream evidence for requalification. Do not invalidate unrelated earlier evidence, and do not retain later PASS evidence across a changed fingerprint.

### Persistence cadence

Commit meaningful implementation work before leaving its phase. Update checkpoint metadata after meaningful commits and gate transitions. Checkpoint-only commits live under the fingerprint-excluded `.depulse-certification/` area so recording recovery state does not mutate the product candidate fingerprint.

### Interrupted CI/native work

Inspect current GitHub run/job state before launching replacements. Resume/retry only incomplete or failed lanes where exact-source reuse is valid. A stopped conversation does not justify rerunning a green G12 or native platform audit.

## Permanent Adaptive CI execution

The release uses one logical Build Coordinator across G0–G16. Jobs may be parallelized only when dependencies and resource limits allow it, while each responsibility retains one canonical owner.

Before expensive qualification, run the preventative-learning preflight against known release-learning patterns. CI failures are classified as `PRODUCT_FAIL`, `GATE_TEST_FAIL`, `CI_HARNESS_FAIL`, `INFRA_FAIL`, `EXPECTED_NOOP`, or `SUPERSEDED` before candidate health is changed.

Mutating workflows are idempotent: a clean/no-change commit attempt is successful no-op behavior, not a failed build. Checkpoint/evidence-only changes must not wake unrelated full certification/native workflows.

The Build State Ledger is reconciled from actual branch HEAD, release identity, fingerprint, CI and artifact evidence. It is never trusted merely because a JSON checkpoint says PASS.

### Safety

Checkpoint/resume and adaptive CI never weaken G0–G16, native runtime requirements, Stable transformation/certification, source provenance, or the permanent No Execution Boundary.

Canonical protocols:
- `adaptive-governance/BUILD_RESUME_PROTOCOL.md`
- `adaptive-governance/ADAPTIVE_CI_OPERATING_CONTRACT.md`
