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

## Permanent Adaptive Work Decomposition Planning

Every release plan must evaluate heavy work before execution and choose the smallest useful evidence boundaries rather than defaulting to monolithic tasks.

For each heavy implementation/qualification/delivery responsibility, the plan should decide:

- keep as one unit, or split into work packages/checkpoints/sub-stages/shards;
- dependency order and which packages may run safely in parallel;
- canonical owner for each package;
- exact inputs/fingerprint/artifact identity;
- PASS/FAIL/BLOCKED or completion criteria and evidence location;
- downstream consumers and invalidation rules;
- shared setup/data/artifacts that prevent duplicated work;
- bounded concurrency based on runner, browser, provider/API, DB, CPU and memory capacity;
- resume behavior after a partial failure or interruption;
- whether intermediate evidence can safely prevent rerunning unrelated work.

Planning follows the adaptive loop:

**Understand → decompose → map dependencies → reuse → execute → checkpoint → evaluate evidence → adapt next work → integrate → certify → learn.**

The planner may adapt sequencing after checkpoints when actual evidence shows a better next action, provided immutable scope, release contracts and required assurance are preserved.

### Gate-model planning

The current G0–G16 map is the default. If a recurring/material responsibility appears to require a new release gate, planning must first attempt checkpointing, sharding/parallel lanes, or strengthening/reassigning an existing gate.

A new gate may be proposed only after the Gate Utility Test proves distinct risk/responsibility, independent evidence, non-duplication, material value, canonical ownership, process-wide updates, migration clarity and planned G16 review. No release plan may introduce an isolated ad-hoc G-gate.

Canonical rules: `adaptive-governance/ADAPTIVE_WORK_DECOMPOSITION_CONTRACT.md`.

## Permanent Governance-to-Implementation Closure Planning

Every build plan must distinguish **governance adoption** from **implementation closure**. Each new permanent rule receives an explicit implementation disposition before G3 exit:

- `CURRENT_RELEASE_BLOCKER`;
- `CURRENT_RELEASE_PROCESS_HARDENING`;
- `NEXT_RELEASE_MANDATORY_ENTRY`;
- `FUTURE_STRATEGIC`.

For every applicable item, the plan records implementation owner, source/workflow targets, dependencies, evidence gate/checkpoint, naming, regression coverage, delivery impact, and exact completion criteria. A documented rule with no implementation/evidence plan is a planning defect.

### v18.2 required plan closure

The v18.2 build plan must explicitly schedule and evidence:

- capability-based ADMIN authorization shared by UI and backend;
- dedicated capability-scoped Administration navigation/composition;
- all-tab/global-shell role × viewport audit for SUPER_OWNER/OWNER, full-capability ADMIN, limited ADMIN, USER and DEMO;
- deterministic Build State Ledger v2 reconciliation from actual GitHub HEAD/fingerprint/CI/artifact truth using canonical evidence-state names;
- consolidation/subordination of overlapping release workflows under one actual Build Coordinator;
- canonical naming migration/cleanup for active workflows, jobs, artifacts, checkpoints and gate labels.

These are not satisfied by documentation alone. Product/security items block the applicable product gates; process-hardening items block trust in the affected release evidence until closed.

### v18.3 mandatory entry

All source-changing utility/consolidation items in `functionality_utility_remediation.json` are mandatory G1–G3 inputs for v18.3. They may be replaced only by a stronger audited disposition, never silently dropped.

Canonical closure rules: `adaptive-governance/GOVERNANCE_IMPLEMENTATION_CLOSURE_CONTRACT.md`.

## Permanent Shared Symbol Intelligence Planning Requirements

Every release that adds or materially changes symbol demand, provider acquisition, subscriptions, scanner/radar behavior, research, preparation/event processing, user workspaces, AI/adaptive synthesis, persistence or hosted multi-user behavior must include a **Shared Symbol Intelligence Processing Plan** before G3 exit.

The plan must define:

- Global Symbol Registry ownership and all demand contributors;
- the shared demand union and rules for entering/leaving it;
- canonical processing keys, including symbol/instrument, dataset/capability, session/time window, freshness/materiality requirement, provider/entitlement/data-rights domain and model/policy version where relevant;
- one canonical owner for acquisition, normalization, validation, canonical state, deterministic calculations, correlation, reusable intelligence and fan-out;
- which consumers reuse each canonical result;
- cache/freshness strategy and in-flight coalescing/single-flight ownership;
- material-change invalidation/propagation dependencies rather than blanket recomputation;
- dynamic attention/priority rules and provider/runtime budgets;
- fairness/backpressure so one user or large watchlist cannot starve higher-value shared work;
- memory-first live state and bounded durable persistence/warm-start behavior;
- AI/evidence fingerprinting and reusable synthesis boundaries;
- explicit non-shareable boundaries for private prompts/context, tenant/security isolation, provider entitlements and data rights;
- multi-user load scenarios with overlapping and non-overlapping symbol demand;
- efficiency scorecard baselines and release acceptance thresholds.

Planning default:

**shared canonical processing → material-change reuse → authorized personal composition.**

A proposal that introduces a per-user provider engine, duplicate symbol pipeline, separate equivalent scanner/prep acquisition path or repeated identical AI synthesis must justify why canonical sharing is impossible or unsafe. Otherwise it is rejected or consolidated.

### Required efficiency planning metrics

Where applicable, the release plan must capture unique active symbols/processing keys, total consumer demand, provider calls/subscriptions per unique key, duplicate acquisition/calculation rate, in-flight coalescing, cache reuse, shared-synthesis reuse, fan-out ratio, marginal cost of an overlapping user, CPU/memory/storage per active key, provider pressure, freshness/material-change latency, stale/degraded fan-out, fairness/starvation and authorization/data-rights leakage.

For equivalent overlapping demand, the target architecture is that cost scales primarily with **unique canonical demand**, not `users × symbols`.

### Release-specific planning

- **v18.2:** prove multi-user `UserWorkspace`/role work does not introduce per-user market/provider/intelligence ownership.
- **v18.3:** mandatory implementation plan for shared Scanner/Radar acquisition, Session Intelligence Coordinator, Event Intelligence lifecycle ownership, hosted demand union and shared canonical persistence/recovery behavior.
- **v18.4:** explicit rights/security/cache/synthesis partitioning plan.
- **v18.5:** realistic overlapping-user performance/capacity/security certification plan using the full 10/10 efficiency scorecard.

Canonical shared-processing rules: `adaptive-governance/SHARED_SYMBOL_INTELLIGENCE_PROCESSING_CONTRACT.md`.
