# DE.PULSE — Adaptive Build Plan

## Permanent Checkpoint & Resume Requirements

Every release build plan must define its durable recovery state before implementation proceeds beyond G3.

Required plan items:

- exact incoming Stable tag/commit/release;
- one active development/release branch;
- canonical release identity and expected runtime profile;
- build checkpoint and release-evidence checkpoint locations;
- fingerprint-keyed qualification strategy;
- checkpoint updates after meaningful code commits and at G3/G4/G10/G11/G12/G14/G15/G16;
- next-step/blocker recording whenever a gate stops;
- G16 archive/cleanup and next-release checkpoint seed.

The build plan must distinguish **implementation complete** from **evidence reusable**. A feature may be built while a later source change requires requalification; the checkpoint must preserve that truth rather than labeling later gates PASS.

Expensive tests and native audits may resume from exact-source evidence when source/RC/artifact identity is unchanged. Any affected evidence must be invalidated after source or packaging changes.

No meaningful local-only work is accepted as a durable build-plan milestone. It must be committed to the active GitHub branch before it can be treated as resumable.

## Permanent Adaptive CI planning requirements

Every release plan must also define:

- one logical Build Coordinator and the G0–G16 dependency graph;
- which qualification lanes may safely run in parallel;
- one canonical owner for each expensive test/gate responsibility;
- preventative-learning preflight before expensive qualification;
- mandatory CI failure classification;
- Build State Ledger / Checkpoint v2 reconciliation against actual GitHub evidence;
- metadata-only path isolation so checkpoint commits do not cause unnecessary recertification;
- idempotent mutation behavior so no-change operations cannot create false red builds;
- separate macOS Apple Silicon and Windows x64 native evidence;
- explicit `User Delivery` completion for native TEST/RC/Stable deliveries as prescribed;
- G16 CI-learning review and workflow/check consolidation.

The detailed reconciliation and invalidation rules are mandatory from `adaptive-governance/BUILD_RESUME_PROTOCOL.md`.
The CI learning/orchestration/delivery rules are mandatory from `adaptive-governance/ADAPTIVE_CI_OPERATING_CONTRACT.md`.
