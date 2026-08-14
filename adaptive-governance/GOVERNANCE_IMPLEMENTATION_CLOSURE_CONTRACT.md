# DE.PULSE — Governance-to-Implementation Closure Contract

Status: **Permanent governing contract**  
Effective: **v18.2 and all later releases**

## Purpose

A DE.PULSE rule is not complete merely because it exists in the Adaptive Roadmap, Adaptive Build Plan, Adaptive Build Process, Adaptive Delivery Process, registry, or documentation. Governance must be connected to actual implementation, automated enforcement, release evidence, and delivery behavior.

Permanent rule:

**Governed → Implemented → Enforced → Evidenced → Delivered → Learned.**

A release may not claim a governance requirement is closed while the source, CI/CD, runtime authorization, UI composition, checkpoint/ledger state, naming, native artifacts, or release evidence still contradict it.

## 1. Mandatory closure classification

Every permanent governance addition must be classified as one of:

- `CURRENT_RELEASE_BLOCKER` — must be implemented and evidenced before the current release can freeze/promote.
- `CURRENT_RELEASE_PROCESS_HARDENING` — build/release machinery required to make the current release evidence trustworthy; must close before the affected gate is claimed PASS.
- `NEXT_RELEASE_MANDATORY_ENTRY` — source-changing work discovered after immutable scope freeze that must enter the next release at G1–G3 and cannot become optional backlog.
- `FUTURE_STRATEGIC` — intentionally deferred architecture/product direction with an explicit target/review point.

No governance item may remain in an ambiguous “documented but not implemented” state without one of these dispositions.

## 2. Current v18.2 implementation-closure obligations

The following are mandatory v18.2 closure items because they either belong directly to Admin / Presence / Sessions scope or are required for trustworthy governance/release evidence.

### 2.1 Capability-Based Administration — CURRENT_RELEASE_BLOCKER

- `ADMIN` role alone must not imply full administrative, provider, runtime, Maintenance, security, or implementation-diagnostics authority.
- Administrative authority is determined by explicit delegated capabilities under the canonical authorization owner.
- OWNER / SUPER_OWNER retain full authorized governance according to ownership policy.
- Selected/full-capability ADMIN may receive deep operational rights only through explicit capability delegation.
- Limited ADMIN receives only assigned capabilities.
- UI visibility, direct navigation, API authorization, payload redaction, and audit evidence must use the same capability truth.

### 2.2 Dedicated Administration Composition — CURRENT_RELEASE_BLOCKER

- Identity, user, role/status, presence, session, password-reset, and delegated-capability administration must use the canonical conditional **Administration** surface when its separation is justified.
- It must not remain permanently embedded in unrelated Settings merely because that was convenient during implementation.
- Navigation and page contents are capability-scoped.
- USER/DEMO must not receive the Administration route/surface or privileged payload.
- Removing the tab for unauthorized roles must leave a seamless navigation/page hierarchy with no gap or orphan state.

### 2.3 Complete Role × Tab × Viewport Audit — CURRENT_RELEASE_BLOCKER

Before v18.2 qualification can claim role-aware closure, every application tab plus global shell/navigation must be audited for:

- SUPER_OWNER / OWNER;
- selected/full-capability ADMIN;
- limited/delegated ADMIN;
- USER;
- DEMO;
- supported desktop/tablet/narrow viewport families.

Evidence must cover visibility, backend authorization, direct-route/API denial, sensitive-payload suppression, responsive reflow, information hierarchy, keyboard/focus order, and absence of implementation machinery for USER/DEMO.

### 2.4 Authoritative Build State Ledger v2 — CURRENT_RELEASE_PROCESS_HARDENING

The durable checkpoint must be automatically or deterministically reconciled from actual GitHub evidence rather than aging as an optimistic hand-maintained JSON note.

Required behavior:

- reconcile actual branch HEAD, canonical release identity, source fingerprint, workflow/job evidence, artifacts, immutable RC, native packages, G15/G16 state, and `User Delivery`;
- use only canonical evidence-state vocabulary;
- correct stale checkpoint values when GitHub truth differs;
- record last trustworthy PASS, earliest required resume gate, one next action, blocker/failure class, platform artifact status and evidence IDs;
- keep metadata fingerprint-isolated and avoid self-triggered full qualification.

A stale ledger cannot be cited as authoritative release state.

### 2.5 One Real Build Coordinator — CURRENT_RELEASE_PROCESS_HARDENING

“One logical Build Coordinator” must become actual orchestration behavior, not only documentation.

- overlapping workflows that independently formalize, mutate, qualify or publish the same candidate must be consolidated, made reusable/subordinate, or path/event scoped so they do not create workflow storms;
- each gate/checkpoint responsibility has one canonical orchestration owner;
- independent lanes may still run in parallel under dependency-aware coordination;
- bot/checkpoint/docs mutations must not recursively wake unrelated expensive release workflows;
- superseded runs must stop consuming authoritative status.

### 2.6 Canonical Naming Migration — CURRENT_RELEASE_PROCESS_HARDENING

The Naming & Identity Contract applies to actual active repository objects, not only new documentation.

Before release closure:

- active workflow/job/artifact names must use canonical gate/capability/release vocabulary;
- obsolete numbered or ambiguous names such as temporary “sync 2”, duplicate “formalize”, “final”, “new”, or phase aliases are renamed/consolidated/removed where they are no longer required for compatibility;
- one concept has one canonical human name and one stable machine identifier;
- aliases remain only where migration/backward compatibility requires them and are explicitly marked deprecated.

## 3. Mandatory v18.3 entry workstream

Source-changing product utility/remediation discovered after v18.2 immutable scope remains `NEXT_RELEASE_MANDATORY_ENTRY` for v18.3 and cannot be treated as optional backlog.

This includes the authoritative items in `functionality_utility_remediation.json`, including:

- shared/coalesced Discovery Scanner + Opportunity Radar acquisition;
- Session Intelligence Coordinator ownership for Pre-Market Prep and Market Open Prep;
- Event Intelligence ownership for Earnings & Material Catalyst Reaction Watch;
- Dashboard and Market Intelligence presentation consolidation;
- Day/Swing/Long deep-evidence de-duplication with Research as canonical ticker deep-evidence home;
- Market Activity demotion to supporting/drill-down data;
- retirement/redirect of legacy standalone News/Earnings/Filings routes;
- Maintenance Session & Event Intelligence diagnostic consolidation;
- regression proof that no material intelligence, protected deterministic logic, role visibility, or decision-support value is lost.

## 4. Gate integration

Under the current canonical gate map:

- **G1 — Immutable Scope:** classify every governance-to-implementation obligation and freeze which are current blockers/process hardening vs next-release mandatory entry.
- **G2 — Architecture & Data Utility:** define canonical owner and reuse/authorization/data boundaries.
- **G3 — Design & Dependency Readiness:** define implementation/evidence plan, naming, role matrix, orchestration dependencies and closure criteria.
- **G4 — Development Exit:** current-release blocker implementation must actually exist in source and backend/frontend boundaries must agree.
- **G7 — Data, Security & Adaptive Intelligence:** capability enforcement, redaction, audit evidence and governed data classes must pass.
- **G9 — Cross-Module UI/UX:** all-tab role/capability/viewport composition and hierarchy evidence must pass.
- **G10 — Pre-Freeze Qualification:** current-release blockers and required process hardening must have source-bound evidence; a documented-only requirement is not PASS.
- **G11/G12:** immutable RC/full certification must consume the reconciled ledger and canonical orchestration evidence.
- **G13/G14/G15:** native/runtime/promotion evidence must preserve the same governance truth.
- **G16 — Adaptive Retrospective & Handoff:** verify every governance addition has an implementation disposition, close completed items, carry mandatory next-release entries forward explicitly, clean obsolete workflow/naming machinery, and convert recurring gaps into preventative checks.

## 5. Delivery invariant

A release is not `DELIVERED` merely because governance documents are correct. Delivery closure requires that the certified runtime and released artifacts actually implement the applicable governance contracts.

If a governance requirement is intentionally deferred, the handoff must name its classification, target release, reason, and blocking/non-blocking status. Silent carry-forward is forbidden.

## Permanent rule

**DE.PULSE governance is executable governance. A rule must either be implemented and evidenced for the release where it applies, or explicitly classified and carried to a named future release. Documentation alone never converts an unfinished requirement into PASS.**
