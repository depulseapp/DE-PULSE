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
- create the next release's initial build checkpoint;
- classify release CI incidents and preserve genuinely new generalized lessons;
- verify both native runnable assets were surfaced to the user for completed native deliveries.

### Recovery after partial promotion

If publication is interrupted, inspect each durable GitHub object independently before retrying. Never assume that a failed workflow means nothing was published, and never create a duplicate tag/release without checking existing refs/releases first.

## User Delivery completion invariant

A native Stable release is `READY` when permanent certified Mac/Windows assets exist, but it is `DELIVERED` only after the user-facing handoff surfaces both runnable assets with provenance/hash status.

Required concise release status:

**Code | G0–G12 | macOS G13/G14 | Windows G13/G14 | G15 | G16 | User Delivery**

A Stable build must not be presented as fully delivered while `User Delivery` is `NOT_READY` or `READY`.

### User burden

Normal delivery recovery is automated from GitHub/CI and requires no routine manual work from the user. Only the approved external-blocker exceptions may interrupt autonomous delivery.

Canonical protocols:
- `adaptive-governance/BUILD_RESUME_PROTOCOL.md`
- `adaptive-governance/ADAPTIVE_CI_OPERATING_CONTRACT.md`

## Permanent Role-Aware Delivery Invariant

A release is not delivery-complete merely because the role-sensitive code compiled or because privileged content was visually hidden. Role-aware composition and authorization are release-quality requirements.

Before G15/Stable promotion, delivery evidence must establish that:

- every current/new tab and global shell/navigation surface passed the required role/capability audit;
- OWNER/SUPER_OWNER receives the intended full authorized composition;
- selected/full-capability ADMIN receives only the explicitly delegated deep capabilities;
- limited ADMIN receives only assigned administrative capabilities;
- USER/DEMO receives no implementation machinery or privileged payloads;
- frontend visibility and backend authorization agree;
- unauthorized direct page/API access is denied/redirected;
- information hierarchy remains intentional for every role rather than being reordered accidentally by hidden content;
- conditional sections reflow seamlessly without blank cells, dead space, orphan headings, uneven leftovers, clipping or awkward whitespace;
- justified conditional tabs such as Administration appear only for authorized capability profiles and do not leave navigation gaps when absent;
- role composition does not alter protected deterministic market logic;
- role × viewport G9 evidence is bound to the certified source/RC and remains valid for the delivered native packages.

G16 must include role-aware delivery closure in the handoff: role/capability matrix result, any known limitations, regressions added, and new generalized learning. Recurring composition/security defects must be converted into reusable primitives/tests rather than rediscovered manually in later releases.

macOS and Windows deliverables must expose the same certified role/capability behavior. A platform package that diverges materially in authorization, navigation or composition cannot inherit the other platform's PASS.

Canonical role-aware rules: `adaptive-governance/ROLE_AWARE_UI_COMPOSITION_CONTRACT.md`.

## Permanent Functionality Utility Delivery Invariant

A release is not delivery-complete if it introduces working but unnecessary or unintegrated product machinery.

Before G15/Stable promotion, certified evidence must establish that the actual candidate passed the Functionality Utility, Reuse, Correlation & Surface Checkpoint and that:

- every primary navigation tab is represented in the functionality utility registry;
- every materially introduced/changed engine, job, watcher, scheduler, dataset-facing workflow and major surface has a canonical owner and active consumer;
- existing data/computation/acquisition was reused where practical and known duplicate provider work is either consolidated or explicitly justified/bounded;
- new data is correlated with relevant canonical evidence and has freshness/materiality/retention/governance behavior;
- repeated deep-evidence presentation is consolidated into a canonical home with concise contextual reuse elsewhere;
- supporting/operational data is not promoted into normal-user UI without material decision-support value;
- any new tab has explicit documented separation justification;
- obsolete or superseded routes/surfaces/jobs identified by the checkpoint are removed, redirected, scheduled for the current release, or explicitly carried as a known closure item rather than silently retained;
- G9 performed an all-tab/major-function repetition and hierarchy audit, not only a changed-file visual check;
- the delivered macOS and Windows packages preserve the same certified information hierarchy and utility dispositions.

G16 must record the final reuse/consolidation/removal result and any remaining approved exceptions. A recurring duplication pattern must become a preventative regression or governance rule.

Canonical checkpoint rules: `adaptive-governance/FUNCTIONALITY_UTILITY_INTEGRATION_CHECKPOINT.md`.
