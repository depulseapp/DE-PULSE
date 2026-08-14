# DE.PULSE — Adaptive Delivery Process

## Permanent Delivery Recovery Contract

Delivery state must be durable, independently verifiable, and resumable from GitHub evidence through G13–G16.

### G13 / G14

Packaging and actual artifact runtime audits are bound to the immutable candidate source/RC and artifact hashes. macOS Apple Silicon and Windows x64 results are tracked separately. If one platform is PASS and the other is interrupted, the PASS platform may be reused only when the exact source and packaging identity remain unchanged.

### G15

Release Assurance consumes immutable source/RC provenance plus required native artifact evidence. An interrupted conversation does not change G15 eligibility; GitHub artifacts/hashes and CI job evidence determine whether G15 may resume or must rerun.

### Promotion

Stable promotion is allowed only from certified evidence. The checkpoint records the intended promotion source, but the actual `main`, Stable tag, permanent release, source/package hashes and release evidence remain authoritative.

### G16

G16 must:

- verify `main` / immutable Stable tag / permanent release identity agreement;
- record final source and native artifact provenance;
- archive the completed release checkpoint/evidence state;
- close or clean obsolete development/promotion branches according to repository hygiene;
- seed the next approved development branch from the exact Stable commit/tag;
- create the next release's initial build checkpoint.

### Recovery after partial promotion

If publication is interrupted, inspect each durable GitHub object independently before retrying. Never assume that a failed workflow means nothing was published, and never create a duplicate tag/release without checking existing refs/releases first.

### User burden

Normal delivery recovery is automated from GitHub/CI and requires no routine manual work from the user. Only the approved external-blocker exceptions may interrupt autonomous delivery.

Canonical protocol: `adaptive-governance/BUILD_RESUME_PROTOCOL.md`.
