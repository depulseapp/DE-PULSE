# DE.PULSE — Adaptive Roadmap

## Permanent Build Recoverability Contract

The Adaptive Roadmap permanently includes **GitHub-backed build recoverability** from v18.2 onward.

Every roadmap release must be resumable after interruption from the last trustworthy GitHub/CI evidence without relying on conversation history. Each release therefore owns an active development/release branch, a durable build checkpoint, fingerprint-bound gate evidence, and immutable release artifacts as they become available.

Roadmap sequencing remains unchanged by an interruption. Recovery determines the last trustworthy PASS, resumes the current approved release at the next required step, completes G0–G16, then proceeds to the next approved roadmap release.

A source change invalidates affected downstream evidence. A conversation/tool interruption by itself does not invalidate unchanged-source evidence and must not force unnecessary reruns.

At G16, the current checkpoint is archived and the next approved roadmap branch/checkpoint is seeded from the exact promoted Stable commit/tag.

## Permanent Adaptive CI learning direction

From v18.2 onward the build system itself follows the DE.PULSE adaptive principle: **observe → classify → diagnose → generalize → prevent → execute → measure → learn**.

Future roadmap releases inherit:

- one logical Build Coordinator inside canonical G0–G16;
- mandatory CI failure classification so harness/infra/no-op/superseded events are not mislabeled as product defects;
- a preventative-learning preflight driven by accumulated release lessons;
- a GitHub-reconciled Build State Ledger / Checkpoint v2;
- idempotent/path-isolated CI to reduce duplicate and self-triggered workflow noise;
- exact-source partial resume rather than unnecessary full reruns;
- mandatory Mac/Windows native user-delivery closure for completed native releases;
- G16 consolidation of obsolete/duplicate workflow machinery while preserving assurance.

The canonical rules are defined in:
- `adaptive-governance/BUILD_RESUME_PROTOCOL.md`
- `adaptive-governance/ADAPTIVE_CI_OPERATING_CONTRACT.md`

They apply to all future roadmap families unless replaced by a stricter approved contract.
