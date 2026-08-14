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

### Safety

Checkpoint/resume never weakens G0–G16, native runtime requirements, Stable transformation/certification, source provenance, or the permanent No Execution Boundary.

Canonical protocol: `adaptive-governance/BUILD_RESUME_PROTOCOL.md`.
