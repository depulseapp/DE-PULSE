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

## Permanent Role-Aware Build Planning Requirements

Every release that adds, removes, moves or materially changes a tab, card, shell control, administrative operation or role-sensitive dataset must define its role/capability composition before implementation exits G3.

The build plan must contain or generate a role/capability matrix covering:

- every current application tab and global shell/navigation surface;
- SUPER_OWNER/OWNER composition;
- selected/full-capability ADMIN composition where delegated;
- limited/delegated ADMIN composition;
- USER composition;
- DEMO composition;
- visible controls and server-authorized actions for each capability;
- data/payload fields that must be withheld server-side, not merely hidden;
- intended information hierarchy and placement of each surviving section;
- responsive reflow behavior when privileged sections are absent;
- direct-navigation/API denial expectations for unauthorized capabilities;
- supported desktop/tablet/narrow-browser role × viewport test coverage.

Planning rules:

- ADMIN capability grants are explicit; `ADMIN` alone must not imply full runtime/provider/security/maintenance authority.
- USER/DEMO plans must exclude implementation machinery rather than merely cosmetically hiding it.
- removal of privileged content must include a planned recomposition/reflow outcome; inherited OWNER geometry is not accepted.
- OWNER/admin controls are placed according to utility and hierarchy, not automatically at the top.
- a proposed new tab must pass a utility/separation test: create it only when it materially improves workflow, security boundary, clarity or maintainability. Do not add tabs for organizational convenience alone.
- when justified, a capability-scoped Administration tab is preferred to mixing delegated identity/session administration into unrelated Settings or Maintenance surfaces.
- frontend visibility and backend authorization changes are planned together as one capability boundary.

G9 is the main UI-composition enforcement point, with security/data authorization additionally verified in G7/G12 and final role-audit closure recorded at G16.

Canonical role-aware rules: `adaptive-governance/ROLE_AWARE_UI_COMPOSITION_CONTRACT.md`.

## Permanent Functionality Utility & Integration Planning Requirements

Every release build plan must include a **Functionality Utility, Reuse, Correlation & Surface Checkpoint** before implementation exits G3.

The plan must inventory every new or materially changed functionality, including tabs/sub-tabs/cards, engines, scanners, detectors, preparation jobs, watchers, schedulers, APIs/provider calls, datasets, derived metrics/models, alerts, persistence fields and administrative operations.

For each item, planning must record:

- purpose and active consumer/workflow;
- canonical owner;
- existing implementation/data that will be reused;
- whether provider acquisition, computation, cache, persistence or in-flight work can be shared/coalesced;
- functional overlap with existing engines/jobs/surfaces and the chosen consolidation/retirement decision;
- required correlations with existing canonical evidence and outcome/learning state;
- freshness/materiality, retention, rights/sensitivity and degraded behavior;
- provider/runtime/storage/UI performance impact;
- intended UI disposition: existing surface, compact contextual reuse, drill-down/internal-only, justified new surface, remove/retire, or defer;
- explicit new-tab justification whenever a new primary/conditional tab is proposed;
- role/capability visibility and backend authorization requirements;
- obsolete implementation/surface/job that can be removed as part of the change.

Planning default: **reuse and correlate before adding; consolidate before creating another owner; create no extra tab unless separation materially improves the product.**

The release must maintain `functionality_utility_registry.json` and keep it aligned with current primary navigation and material background/intelligence capabilities. The plan must include `functionality_utility_checkpoint_gate.py` in pre-freeze qualification.

Canonical rules: `adaptive-governance/FUNCTIONALITY_UTILITY_INTEGRATION_CHECKPOINT.md`.
