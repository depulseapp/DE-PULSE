# DE.PULSE — Adaptive Roadmap

## Permanent Build Recoverability Contract

The Adaptive Roadmap permanently includes **GitHub-backed build recoverability** from v18.2 onward.

Every roadmap release must be resumable after interruption from the last trustworthy GitHub/CI evidence without relying on conversation history. Each release therefore owns an active development/release branch, a durable build checkpoint, fingerprint-bound gate evidence, and immutable release artifacts as they become available.

Roadmap sequencing remains unchanged by an interruption. Recovery determines the last trustworthy PASS, resumes the current approved release at the next required step, completes G0–G16, then proceeds to the next approved roadmap release.

A source change invalidates affected downstream evidence. A conversation/tool interruption by itself does not invalidate unchanged-source evidence and must not force unnecessary reruns.

At G16, the current checkpoint is archived and the next approved roadmap branch/checkpoint is seeded from the exact promoted Stable commit/tag.

The canonical rules are defined in `adaptive-governance/BUILD_RESUME_PROTOCOL.md` and apply to all future roadmap families unless replaced by a stricter approved contract.
