# DE.PULSE — Canonical Naming & Identity Contract

Status: **Permanent governing contract**  
Applies to: product UI, source, APIs, data models, Adaptive Roadmap, Adaptive Build Plan, Adaptive Build Process, Adaptive Delivery Process, CI/CD, checkpoints, artifacts, release evidence and handoffs  
Effective: **v18.2 and all later releases**

## Purpose

DE.PULSE must use one stable vocabulary for the same concept across product, engineering and delivery. Ambiguous, duplicate, historical or convenience names create implementation drift, incorrect ownership, confusing UI, brittle tests and unreliable release evidence.

Permanent rule:

**One concept → one canonical name → one canonical machine identifier → one canonical owner.**

Aliases may exist only for backward compatibility or migration and must be explicitly marked deprecated/legacy. They must not become competing active names.

## 1. Product branding

Canonical product name: **DE.PULSE**.

Use `DE.PULSE` in user-visible titles, documentation and release handoffs. Repository/package/file-system identifiers may use `De-Pulse` or `DE-PULSE` only where platform/tooling conventions require it; do not invent additional product spellings.

## 2. Version, channel and release identity

Canonical semantic version: `vMAJOR.MINOR.PATCH`.

Canonical channels:
- `TEST`
- `RC`
- `STABLE`

Canonical branch pattern:
- active development: `v<major>.<minor>-development`
- release-specific development when patch distinction is required: `v<major>.<minor>.<patch>-development`

Canonical Stable tag:
- `v<major>.<minor>.<patch>-stable`

Canonical build ID pattern:
- `v<version>-<channel-lower>-<short-release-purpose>-<YYYYMMDD>`

Canonical native asset pattern:
- `De-Pulse-v<version>-<CHANNEL>-macOS-Apple-Silicon.zip`
- `De-Pulse-v<version>-<CHANNEL>-Windows-x64.zip`

Names must be derived from canonical release identity where practical rather than duplicated as hard-coded literals across scripts/workflows.

## 3. Canonical release-gate names

Under the current gate map, use these exact gate names in governance, workflow/job labels, evidence, checkpoints and handoffs:

- **G0 — Exact Baseline**
- **G1 — Immutable Scope**
- **G2 — Architecture & Data Utility**
- **G3 — Design & Dependency Readiness**
- **G4 — Development Exit**
- **G5 — FAST Qualification**
- **G6 — Integration & MEDIUM Qualification**
- **G7 — Data, Security & Adaptive Intelligence**
- **G8 — Performance, Capacity & Stability**
- **G9 — Cross-Module UI/UX**
- **G10 — Pre-Freeze Qualification**
- **G11 — Immutable Release Candidate**
- **G12 — Full Certification**
- **G13 — Native Packaging & Provenance**
- **G14 — Actual Artifact Runtime Audit**
- **G15 — Release Assurance & Promotion**
- **G16 — Adaptive Retrospective & Handoff**

Do not shorten a gate to a vague label such as `security`, `UI`, `package`, `final`, `release`, or `cert` in authoritative evidence. Short labels may be used only as secondary UI text when the canonical `G#` identity remains visible/derivable.

Historical Stable releases retain the names/evidence under which they were certified. A future gate-map change follows the Gate Utility Test and updates this contract atomically with the other governing contracts.

## 4. Checkpoint and work-package naming

A checkpoint is not a gate. Name checkpoints so their parent gate/responsibility is obvious.

Preferred checkpoint identifier:

`CP-<GATE>-<AREA>-<PURPOSE>`

Examples:
- `CP-G2-FUNCTIONALITY-UTILITY`
- `CP-G3-ROLE-COMPOSITION-READINESS`
- `CP-G8-RUNTIME-LOAD-PROFILE`
- `CP-G9-ROLE-UI-COMPOSITION`
- `CP-G14-MACOS-RUNTIME`

Heavy decomposed implementation/audit units use:

`WP-<VERSION>-<AREA>-<NN>`

Each work package/checkpoint records a human title, machine ID, canonical owner, dependencies, evidence and invalidation rules.

Avoid names such as `phase2`, `final2`, `new-check`, `misc`, `temp`, `fix`, `latest`, `actual`, or `final-final` as durable release identities.

## 5. CI workflow and job naming

Canonical workflow display pattern:

`DE.PULSE <version> | <gate/range> | <purpose>`

Examples:
- `DE.PULSE v18.2 | G0–G4 | Foundation Qualification`
- `DE.PULSE v18.2 | G5–G10 | Pre-Freeze Qualification`
- `DE.PULSE v18.2 | G13–G14 | macOS Apple Silicon`
- `DE.PULSE v18.2 | G13–G14 | Windows x64`
- `DE.PULSE v18.2 | G15–G16 | Promotion & Closure`

Canonical job display pattern:

`<GATE> · <LANE> · <PURPOSE>`

Examples:
- `G9 · USER Desktop · Role Composition`
- `G9 · Limited ADMIN Narrow · Role Composition`
- `G8 · Runtime · Load & Stability`
- `G14 · Windows x64 · Actual Artifact Runtime`

Workflow file names should be lowercase kebab-case and describe the durable responsibility, not the temporary troubleshooting history.

Deprecated/one-shot workflows must be deleted or moved out of the active release path at G16.

## 6. Failure, state and evidence vocabulary

Use the canonical failure classes exactly:
- `PRODUCT_FAIL`
- `GATE_TEST_FAIL`
- `CI_HARNESS_FAIL`
- `INFRA_FAIL`
- `EXPECTED_NOOP`
- `SUPERSEDED`

Canonical gate/checkpoint states:
- `PASS`
- `FAIL`
- `PENDING`
- `BLOCKED`
- `INVALIDATED`
- `SUPERSEDED`

Canonical user-delivery states:
- `NOT_READY`
- `READY`
- `DELIVERED`

Do not invent synonyms such as `OK`, `DONE`, `GOOD`, `ERROR`, `SKIPPED`, or `COMPLETE` in authoritative ledgers unless mapped explicitly to a canonical state.

## 7. Role and capability naming

Canonical roles:
- `SUPER_OWNER`
- `OWNER`
- `ADMIN`
- `USER`
- `DEMO`

`ADMIN` is a role family, not a synonym for full operational authority.

Administrative rights use explicit capability names. Capability machine identifiers should follow:

`<DOMAIN>.<RESOURCE>.<ACTION>`

Examples:
- `identity.users.read`
- `identity.users.manage`
- `identity.sessions.revoke`
- `runtime.engine.control`
- `providers.configuration.manage`
- `maintenance.diagnostics.read`

UI labels may be friendlier, but authorization, audit evidence and tests use the canonical capability identifier.

Avoid vague capabilities such as `full_admin`, `super_access`, `power_user`, or `manage_all` unless they are documented bundles composed from explicit canonical capabilities.

## 8. Product surface naming

Every primary/conditional tab, sub-view and deep-evidence home has one canonical user-visible name and one stable machine `surfaceId`.

Current primary/conditional surface names:
- Dashboard — `dashboard`
- Market Intelligence — `market-intelligence`
- Day Trade Desk — `day`
- Swing Desk — `swing`
- Long-Term Desk — `long`
- Discovery — `discovery`
- Research — `research`
- AI Copilot — `ai`
- Administration — `administration` (capability-scoped)
- Maintenance — `maintenance`
- Settings — `settings`
- Documentation — `documentation`

A new surface must be registered in the functionality utility registry before G3 exit. Renames require a migration/redirect plan and cannot silently create a second destination for the same workflow.

## 9. Engine, coordinator, watcher and job naming

Names must describe responsibility, not implementation mechanism.

Use the following semantic categories consistently:
- **Engine** — owns a distinct continuous/intelligence responsibility.
- **Coordinator** — orchestrates existing canonical owners/checkpoints without becoming a duplicate data/analysis engine.
- **Checkpoint** — time/state-specific reconciliation or qualification boundary.
- **Watcher** — selective event-lifecycle monitoring with explicit trigger/exit conditions.
- **Scanner** — intentional candidate-search/ranking workflow.
- **Detector** — low-latency event recognition from canonical live state.
- **Router** — canonical provider/request serving decision owner.
- **Store/Repository** — canonical state/evidence persistence owner.

Canonical examples:
- `Discovery Scanner`
- `Opportunity Radar`
- `Rapid Move / Market Shock Detector`
- `Session Intelligence Coordinator`
- `Pre-Market Checkpoint`
- `Market Open Checkpoint`
- `Earnings & Material Catalyst Reaction Watch`
- `Event Intelligence`
- `Smart Provider Router`

Do not use `scanner`, `watcher`, `engine`, and `prep` interchangeably for different responsibilities.

## 10. Dataset and evidence naming

Every dataset/evidence family has:
- canonical dataset name;
- machine key/identifier where needed;
- canonical owner;
- purpose/consumer;
- retention/sensitivity;
- lineage/provenance identity.

Names must distinguish raw/provider evidence, canonical normalized evidence, derived intelligence and outcome/learning evidence. Avoid generic names such as `data`, `info`, `payload`, `result`, or `history` in durable governance when the actual semantic class can be named.

## 11. File, registry and artifact naming

Governance files use uppercase snake-style descriptive names where they represent permanent contracts, for example:
- `ADAPTIVE_BUILD_PROCESS.md`
- `ROLE_AWARE_UI_COMPOSITION_CONTRACT.md`
- `NAMING_AND_IDENTITY_CONTRACT.md`

Machine registries/checkpoint files use lowercase snake_case or established JSON path conventions consistently, for example:
- `data_utility_registry.json`
- `functionality_utility_registry.json`
- `.depulse-certification/resume/build-checkpoint.json`

Evidence artifacts should include release version, gate/checkpoint identity, platform/lane when relevant, and fingerprint/hash reference in their manifest rather than relying on ambiguous filenames.

## 12. Naming audit checkpoint

Naming consistency is a cross-gate checkpoint:

- **G1:** scope introduces canonical names and deprecates replaced names.
- **G2/G3:** owners, data, surfaces, capabilities and work packages use canonical identifiers before implementation.
- **G4:** source/API/UI implementation must not introduce competing names.
- **G9:** user-visible navigation/headings and role-specific compositions use the same canonical product vocabulary.
- **G10:** pre-freeze naming audit checks release identity, gates, workflows, registries, primary surfaces and durable evidence for drift.
- **G12:** certification fails on material identity ambiguity that could affect authorization, evidence reuse, routing or delivery.
- **G16:** remove obsolete aliases/workflow names/files where compatibility no longer requires them and record approved remaining aliases.

## 13. Rename rule

A rename is an architecture/release change, not cosmetic search-and-replace, when the name is used by routes, APIs, persistence keys, capabilities, CI evidence or external artifacts.

Before a material rename:
1. identify canonical old identity;
2. choose canonical new identity;
3. enumerate consumers;
4. define compatibility/redirect/migration behavior;
5. update tests/docs/governance/CI together;
6. remove the old identity when safe;
7. record the rename in release evidence.

## Permanent rule

**Names are part of system architecture. DE.PULSE uses precise, durable, non-overlapping names so a user, developer, CI runner, checkpoint and future AI agent all mean the same thing when they refer to a gate, role, capability, surface, engine, dataset, artifact or release.**
