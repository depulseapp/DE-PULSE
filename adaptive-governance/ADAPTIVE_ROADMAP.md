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

## Permanent Role-Aware UI Composition & Capability Governance

From v18.2 onward every roadmap release must preserve a deliberate, role-aware product experience across the complete application shell and every current or newly introduced tab.

The governing principle is **role/capability changes composition, not information hierarchy**. DE.PULSE must not design one OWNER page and merely hide cards for everyone else. Each role receives an intentional composition built from the same canonical information hierarchy, authorized capabilities and reusable responsive primitives.

Permanent roadmap direction:

- **SUPER_OWNER / OWNER** retain full authorized governance and implementation/operational visibility.
- **ADMIN is capability-based, not automatically full-power.** Selected admins may receive deep/full operational rights when explicitly delegated; other admins receive only their assigned administrative capabilities.
- **USER / DEMO receive no implementation machinery** such as provider plumbing, queue/cache/database/scheduler internals, subscription capacity, routing/circuit/fallback internals, secrets or policy/model implementation controls.
- role visibility is never security; matching backend/API authorization and server-side sensitive-data suppression are mandatory.
- removing unauthorized content must trigger natural layout recomposition: no blank grid cells, orphan headings, dead space, placeholder geometry, hierarchy inversion, clipping or awkward page endings.
- privileged diagnostics/administration remain contextually placed; they must not move above higher-value market intelligence merely because the viewer is privileged.
- a new tab is allowed only when it materially improves utility, clarity, security or workflow separation. It must not be created merely to relocate content. A dedicated **Administration** tab is permitted/expected when it cleanly serves delegated administrative capabilities; its navigation and contents remain capability-scoped.
- every new tab or major surface automatically enters the role/capability audit contract.
- protected deterministic Day/Swing/Long logic and the No Execution Boundary are role-invariant.

This roadmap contract is enforced inside existing G0–G16, especially G1/G2/G3/G9/G10/G12/G16. It creates no new gate.

Canonical role-aware rules: `adaptive-governance/ROLE_AWARE_UI_COMPOSITION_CONTRACT.md`.

## Permanent Functionality Utility, Reuse, Correlation & Minimal-Surface Direction

From v18.2 onward every roadmap release also inherits the **Functionality Utility, Reuse, Correlation & Surface Checkpoint**.

DE.PULSE must not grow by appending independent features, tabs, cards, background jobs, datasets, metrics or provider paths when existing canonical owners/evidence can be reused, extended or consolidated.

Permanent roadmap direction:

- every proposed functionality must have a clear purpose and active consumer/workflow;
- existing canonical data, computations, caches, subscriptions, evidence and persistence must be reused before new acquisition/storage is approved;
- overlapping engines/jobs are consolidated or assigned explicit non-overlapping responsibilities;
- new data must correlate with relevant existing evidence rather than remain an isolated display value;
- provider requests and expensive computations are fetch-once/calculate-once/reuse where practical, including in-flight coalescing for simultaneous demand;
- background preparation/watch/integrity work is modeled as temporal checkpoints, event evaluations or maintenance tasks over canonical state unless a genuinely separate engine is required;
- deep evidence has one canonical home; other surfaces reuse concise conclusions, material risks, freshness and navigation rather than duplicating full evidence packages;
- supporting/operational data remains internal or drill-down when it does not justify normal-user prominence;
- a new tab is the exception, not the default, and requires documented workflow/security/clarity separation value;
- every G9 audit reviews **all tabs and major functionality**, not just the features changed by the release, for accumulated repetition, hierarchy drift and obsolete surfaces/jobs;
- G16 records reuse, consolidation, removals and newly discovered duplication so the product becomes simpler as intelligence capability grows.

This direction is enforced within G1/G2/G3/G7/G8/G9/G10/G16 and does not create a new G-gate.

Canonical rules: `adaptive-governance/FUNCTIONALITY_UTILITY_INTEGRATION_CHECKPOINT.md`.

## Permanent Adaptive Work Decomposition & Process Evolution Direction

From v18.2 onward, the engineering roadmap follows the same AI/LLM-style adaptive reasoning pattern as the product:

**Understand → decompose → map dependencies → reuse → execute → checkpoint → evaluate evidence → adapt next work → integrate → certify → learn.**

Large roadmap items, implementation phases, qualification suites, all-tab audits, performance exercises and delivery work should be split into smaller independently verifiable units whenever that materially improves reliability, speed, recovery, fault isolation, evidence reuse, ownership or learning.

Permanent direction:

- do not force a heavy responsibility into one monolithic task when meaningful checkpoints or independent lanes are safer or more efficient;
- use dependency-aware parallelism for genuinely independent work while preserving one logical Build Coordinator and one canonical owner per responsibility;
- preserve exact-source PASS evidence across unrelated failures and rerun only affected/dependent work;
- use intermediate evidence to adapt sequencing and remediation rather than blindly continuing a fixed script;
- measure the total cost of decomposition so more jobs do not accidentally create more CPU/memory/provider/API/browser/database load than the original task;
- remove redundant checkpoints/jobs when G16 evidence shows they add ceremony without assurance value.

### Gate-model evolution

G0–G16 remains the default canonical gate map. However, the release-gate model itself may evolve when accumulated evidence proves a materially distinct assurance responsibility cannot be represented cleanly as a checkpoint, sub-stage, parallel lane or strengthened existing gate.

Any proposed new G-gate must pass the Gate Utility Test in `adaptive-governance/ADAPTIVE_WORK_DECOMPOSITION_CONTRACT.md` and must update the complete governance/automation model coherently. Historical Stable releases keep the gate map under which they were certified; the process is never retroactively rewritten.

Preferred evolution order:

**reuse existing evidence → checkpoints → sharded/parallel lanes → strengthen/reassign an existing gate → add a new canonical gate only when materially justified.**

Canonical rules: `adaptive-governance/ADAPTIVE_WORK_DECOMPOSITION_CONTRACT.md`.

## v18.2 Audit Carry-Forward — Mandatory v18.3 Entry Workstream

The all-tab/functionality audit performed during v18.2 identified source-changing consolidation/removal work after the immutable v18.2.0 Admin/Presence/Sessions scope had already been frozen. To preserve G1 integrity, those unrelated source changes are not smuggled into the v18.2.0 candidate.

They are mandatory inputs to **v18.3.0 G1–G3 before hosted-web/PostgreSQL implementation proceeds**. The authoritative item-by-item plan and acceptance criteria live in `functionality_utility_remediation.json`.

The v18.3 entry workstream includes:

- shared/coalesced broad snapshot acquisition for Discovery Scanner and Opportunity Radar while retaining distinct intelligence responsibilities;
- Session Intelligence Coordinator ownership for Pre-Market Prep and Market Open Prep as temporal checkpoints over canonical state;
- Event Intelligence lifecycle ownership for Earnings & Material Catalyst Reaction Watch;
- Dashboard/context and Market Intelligence presentation consolidation;
- Day/Swing/Long deep-evidence de-duplication with Research as the canonical ticker deep-evidence home;
- Market Activity retained as a reusable discovery input but removed from unnecessary normal-user prominence;
- legacy standalone News/Earnings/Filings evidence routes retired or redirected to canonical homes;
- Maintenance preparation/event diagnostics consolidated without mixing operational truth into Settings;
- regression coverage proving no material intelligence, protected deterministic logic or authorized operational visibility is lost.

v18.3 cannot treat this workstream as optional backlog. G1/G2/G3 must either implement the recorded remediation or produce a stronger audited replacement disposition before hosted implementation advances.

## Permanent Governance-to-Implementation Closure Direction

From v18.2 onward, DE.PULSE distinguishes **governance adopted** from **implementation closed**. A roadmap item is not complete merely because the rule exists in governance documents.

Permanent lifecycle:

**Governed → Implemented → Enforced → Evidenced → Delivered → Learned.**

Every new permanent rule must be classified as `CURRENT_RELEASE_BLOCKER`, `CURRENT_RELEASE_PROCESS_HARDENING`, `NEXT_RELEASE_MANDATORY_ENTRY`, or `FUTURE_STRATEGIC`. Ambiguous documented-but-unimplemented state is forbidden.

### v18.2 mandatory closure

Before v18.2 Stable promotion, the roadmap requires closure of:

- explicit capability-based ADMIN authorization rather than blanket `ADMIN` power;
- the justified capability-scoped **Administration** tab/composition rather than permanent embedding in Settings;
- full role × tab × viewport audit across SUPER_OWNER/OWNER, full-capability ADMIN, limited ADMIN, USER and DEMO;
- an authoritative auto/deterministically reconciled Build State Ledger v2 rather than a stale manual checkpoint;
- one actual Build Coordinator/orchestration owner with overlapping/self-triggering release workflows consolidated or subordinated;
- repo-wide canonical naming migration for active gates/workflows/jobs/artifacts/capabilities, with obsolete ambiguous names removed or explicitly deprecated.

These items cannot be treated as optional backlog merely because their governing contracts already exist. Process-hardening items must close before the affected evidence is trusted; product/security blockers must close before freeze/promotion.

The v18.3 functionality-consolidation workstream above remains mandatory `NEXT_RELEASE_MANDATORY_ENTRY`, preserving v18.2 immutable product scope while preventing the audit findings from being lost.

Canonical closure rules: `adaptive-governance/GOVERNANCE_IMPLEMENTATION_CLOSURE_CONTRACT.md`.

## Permanent Shared Symbol Intelligence Processing Direction

From v18.2 onward DE.PULSE also adopts the **Shared Symbol Intelligence Processing Contract**. The architectural North Star is:

**Global/shared market intelligence first; authorized personal composition second.**

The Global Symbol Registry owns instrument identity/shared processing membership. User workspaces contribute demand, membership, preferences and workflow context but must not create independent provider, quote, history, event, indicator or reusable intelligence pipelines for symbols that are already active in the shared processing union.

Permanent roadmap direction:

- build one shared demand union across legitimate user/system consumers;
- process each equivalent canonical symbol/dataset key once where lawful and technically equivalent;
- use shared provider subscriptions/acquisition, canonical state, deterministic calculations, correlation, evidence packages and reusable AI/adaptive synthesis;
- collapse simultaneous equivalent work with freshness-aware cache and in-flight coalescing;
- propagate material changes to affected downstream consumers instead of blindly recomputing all state;
- dynamically allocate attention by materiality, session, event risk, decision relevance and capacity rather than treating all symbols equally;
- personalize authorization, ranking, notifications and UI composition at the consumer layer without multiplying the market-data pipeline;
- partition sharing whenever provider entitlement, data-rights, tenant/security domain or private user/LLM context differs;
- measure duplicate acquisition/calculation, cache/coalescing reuse, fan-out, marginal overlapping-user cost, latency, CPU/memory/storage, provider pressure, fairness and leakage incidents;
- treat uncontrolled per-user duplicate market pipelines or repeated identical AI synthesis as architecture/performance defects.

### Release progression

- **v18.2:** governance/invariant protection — multi-user identity/workspace work must not introduce per-user market/provider/intelligence engines.
- **v18.3:** mandatory shared-execution implementation — shared Scanner/Radar acquisition, Session Intelligence Coordinator, Event Intelligence ownership, hosted shared demand union and canonical shared state compatible with persistence/recovery.
- **v18.4:** security/data-rights isolation closure for shared caches, synthesis, provider entitlements and private context.
- **v18.5:** full 10/10 multi-user efficiency/capacity/security closure under realistic overlapping-demand load; processing cost must scale primarily with unique canonical demand rather than `users × symbols`.

This direction is enforced through existing G1/G2/G3/G4/G6/G7/G8/G9/G10/G12/G16 and does not create a new gate by itself.

Canonical rules: `adaptive-governance/SHARED_SYMBOL_INTELLIGENCE_PROCESSING_CONTRACT.md`.
