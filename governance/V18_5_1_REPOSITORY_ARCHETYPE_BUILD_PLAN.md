# DE.PULSE v18.5.1 — Audit/Recovery + 10/10 Repository Archetype Build Plan

**Release type:** v18 audit/containment entry slice, selected urgent recovery and repository/source-architecture hardening  
**Dependency:** starts only after v18.5.0 Stable promotion and native delivery  
**Qualification standard:** closure-grade G0–G16; not a lightweight cleanup  
**Product scope:** implement and prove already-approved v17/v18 commitments only; no unapproved expansion, execution scope or protected decision-formula change.

## Mission

Reconstruct every applicable v17/v18 commitment and escaped defect, implement the safest G1-selected urgent scope, explicitly place every remainder into the adaptive v18.x program, and leave the source in a predictable structure that supports later slices.

Requirement continuity is the primary acceptance boundary; repository restructuring is a controlled hardening lane and may not hide or invalidate current or next-slice work. This release implements `governance/REPOSITORY_STRUCTURE_CONTRACT.md`, follows `governance/V18_ADAPTIVE_RECOVERY_AND_CLOSURE_PROGRAM.md`, and is accepted only at **10/10** for its frozen slice.

## Non-negotiable rules

1. **Structure changes must preserve behavior; explicit requirement-linked defect fixes and approved missing implementations are separately qualified.**
2. Every moved/deleted file must have an inventoried disposition and dependency/path analysis.
3. Do not export or duplicate Go symbols merely to make package moves compile.
4. Preserve canonical shared state, fetch-once/calculate-once, Smart Provider Router v2, bounded work/backpressure, truthful ADR-GDI degradation, user isolation, US Equities Processing and No Execution boundaries.
5. A structural move that exposes or requires a semantic defect fix is treated as a separate change and requalified from the earliest invalidated gate.
6. v18.5.1 must complete the canonical G0–G16 process and native runtime delivery despite being a patch version.

## Mandatory escaped-defect and implementation closure scope

v18.5.1 is no longer accepted as structure-only equivalence work. It must preserve the repository-archetype mission, reconstruct every escape/miss below and freeze which items are completed now versus explicitly assigned to the next v18.x slice:

1. **COPY-18.5.1-001:** stale `v17` profile-preservation copy.
2. **SYMBOL-18.5.1-001:** inert/inconsistent row removal in Day, Swing and Long, including final desk membership.
3. **SYMBOL-18.5.1-002:** current desk membership is not visually distinguished from other memberships.
4. **NAV-18.5.1-001:** live/SSE refresh can unexpectedly jump a user from mid-section to the page top or lose interaction state.
5. **RESEARCH-v15.1.0-17-19-REOPENED:** approved Research ticker input, freshness placement and responsive top-area requirements require fresh closure, including truthful labels, action semantics, disabled states and evidence recovery.
6. **IMPL-18-TRADEINSIGHT-001:** the committed TradeInsight SHADOW / SECONDARY workstream is absent from current source and must be implemented through the canonical Smart Router with entitlement, rights, provenance, tests and package proof.
7. **IMPL-18-UTILITY-001:** Scanner and Opportunity Radar must use one bounded/coalesced snapshot acquisition owner rather than independent equivalent Alpaca fetch paths.
8. **IMPL-18-UTILITY-002:** Pre-Market and Market Open must be governed by one Session Intelligence Coordinator.
9. **IMPL-18-UTILITY-003:** Market Activity Seeds must cease to be a prominent normal-user Discovery surface.
10. **IMPL-18-UTILITY-004:** legacy standalone News/Earnings/Filings routes must retire or safely redirect to canonical evidence homes.
11. **IMPL-18-DOC-001:** role-aware documentation composition and the Documentation Impact Manifest must be implemented and enforced.
12. **IMPL-17-DEPS-001:** the External Dependency & Provider Readiness checkpoint and User Action Required register must be implemented over the existing capability/rights foundations.
13. **VERSION-18.5.1-002:** obsolete active release strings in setup, renderer, runtime and TEST-profile paths must be removed or explicitly classified and tested.
14. **UTILITY-v18.3-CARRY-FORWARD:** all 13 mandatory `functionality_utility_remediation.json` records must receive current implementation/evidence disposition; none may inherit closure from plan presence.

The release must reconstruct all 48 inherited requirements, all 20 frozen v17 items, every v18 major workstream and v18.0–v18.5 release entry, all 13 orphaned remediations, accepted conversational commitments and defects. TradeInsight plus six additional implementation-gap groups are confirmed missing; six escaped defect groups are reopened. Items not frozen into v18.5.1 require explicit next-v18.x placement, owner, reason and user impact.

These items are **CURRENT_RELEASE_BLOCKER**. A repository move may not hide, defer or claim equivalence for them. Each item requires:

**origin → current observation → owner → fix/disposition → regression test → browser retest → macOS package retest → Windows package retest → closure evidence**.

A repeated report reopens the original item, increments recurrence and triggers G16 root-cause/prevention work. Historical PASS evidence is invalid for an affected surface after relevant source, renderer, dependency or packaging changes.


## Target archetype

```text
DE-PULSE/
├── cmd/
│   └── depulse/                    # executable composition / entry point
├── internal/
│   ├── app/                        # application orchestration/state
│   ├── domain/                     # domain models/invariants
│   ├── intelligence/               # research/adaptive/decision support
│   ├── market/                     # market state/scanners/readiness/catalysts
│   ├── providers/                  # provider adapters/router/subscriptions
│   ├── persistence/                # repository/database/PostgreSQL/SQLite
│   ├── runtime/                    # lifecycle/jobs/load/backpressure
│   ├── security/                   # identity/auth/RBAC/session/security
│   └── transport/                  # HTTP/API/DTO boundaries
├── renderer/                       # desktop/web renderer assets
├── tests/
│   ├── integration/
│   ├── acceptance/
│   ├── ui/
│   ├── performance/
│   ├── degradation/
│   └── fixtures/
├── tools/
│   ├── gates/
│   ├── ci/
│   ├── release/
│   └── dev/
├── config/
│   ├── policies/
│   └── schemas/
├── docs/
│   ├── architecture/
│   ├── governance/
│   ├── operations/
│   └── releases/
├── assets/
│   └── reference/
├── release/<version>/
├── .depulse-certification/
├── .github/workflows/              # active workflows only
├── README.md
├── VERSION.txt
├── go.mod
└── go.sum
```

## Workstreams

### 1. Machine-readable inventory + dependency graph

Inventory every tracked file with:

- current path;
- category and owning subsystem;
- production/test/tool/config/doc/release classification;
- callers/importers/path references where determinable;
- disposition `KEEP / MOVE / DELETE / REVIEW`;
- target path;
- risk class;
- verification owner.

No broad move group begins until path-dependent callers have been enumerated.

### 2. Dead/duplicate repository hygiene

After dependency verification, remove from the active branch:

- completed one-shot/discovery workflows;
- version-specific workflows superseded by current canonical workflows;
- generated/transient files accidentally tracked;
- obsolete duplicate gates/scripts;
- dead reference outputs no longer used by docs/tests/runtime.

Git history and immutable Stable tags are the historical recovery mechanism. Uncertain items may temporarily enter `archive/to-review/<date>/` with a disposition manifest; this is not a permanent dumping ground.

### 3. Deterministic non-Go relocation first

Move lower-risk categories before changing Go package boundaries:

- qualification gates/helpers → `tools/gates/` or `tools/ci/`;
- identity/package/provenance helpers → `tools/release/`;
- acceptance/UI harnesses → `tests/acceptance/` / `tests/ui/`;
- performance/degradation harnesses → matching `tests/` subtree;
- policy/registry JSON → `config/policies/`;
- reference screenshots/media → `assets/reference/`;
- long-lived documentation/governance → `docs/` hierarchy where compatible.

Every caller/workflow path change lands atomically with the move that requires it.

### 4. Go package decomposition

After non-Go relocation is green, migrate the broad root `package main` ownership into `cmd/depulse` + focused `internal/...` packages.

Rules:

- boundaries follow ownership and invariants, not arbitrary file counts;
- dependency direction must be documented and cycle-free;
- interfaces stay narrow;
- provider, persistence, intelligence, runtime/backpressure, security and transport ownership remain explicit;
- UI/transport DTO concerns stay out of core market/intelligence packages where practical;
- no duplicate canonical state, router, scanner or provider ownership is introduced.

### 5. Developer experience and anti-regression

Add/update:

- concise root README developer map;
- architecture/module ownership diagram;
- dependency-direction diagram;
- build/test/certification/package commands;
- active-workflow manifest;
- repository-structure gate preventing renewed root sprawl;
- path-reference gate preventing broken script/workflow links;
- package ownership guidance / CODEOWNERS where useful.

### 6. Adaptive GitHub Actions archetype

Implement `CI-ADAPTIVE-18.5.1-001` as part of repository cleanup:

- inventory all workflows and delete/retire obsolete v18.5 development, pre-freeze, discovery and one-shot variants from the active branch after dependency verification;
- converge on durable `ci-fast.yml`, reusable `ci-qualified.yml` and `release.yml` ownership;
- remove CI source-edit/self-delete/push patterns and use parameterized dispatch for diagnostics;
- add action/workflow lint, PowerShell lint, UTF-8 and OS-aware native-harness tests;
- enable supported dependency caching, bounded concurrency and evidence-based retention;
- record planner decisions, failure classes, runtime/cost/cache telemetry and evidence fingerprints;
- keep read-only permissions by default and isolate publication writes;
- preserve final exact-candidate macOS Apple Silicon and Windows x64 proof regardless of budget.

The $5 Actions budget is a soft telemetry threshold, not a completion criterion or quality ceiling.

## 10/10 Repository Archetype acceptance rubric

**All ten dimensions must independently PASS. There is no averaging. A 9/10 result blocks v18.5.1 Stable.**

| # | Dimension | 10/10 requirement |
|---|---|---|
| 1 | Root clarity | Root contains only approved entry-point metadata/directories; no test/gate/config/version clutter. |
| 2 | Production ownership | Runtime code is located by subsystem ownership and a new engineer can identify the owner without filename archaeology. |
| 3 | Package/dependency direction | `cmd` composition and `internal` dependencies are documented, cycle-free and architecture-gated. |
| 4 | Test organization | Unit/integration/acceptance/UI/performance/degradation evidence has predictable homes and no broken discovery. |
| 5 | Tool/CI organization | Gates, CI orchestration, release tooling and dev utilities are separated; active workflows contain no historical accident. |
| 6 | Config/policy organization | Policies, registries and schemas are centralized, typed/validated where applicable, and not scattered through root. |
| 7 | Docs/governance navigation | Architecture, operations, governance and release docs have an understandable index and diagrams where useful. |
| 8 | Workflow hygiene | `.github/workflows` contains only current purposeful workflows; one-shot workflows retire automatically/explicitly after use. |
| 9 | Dead/duplicate artifact hygiene | No known dead/duplicate file remains without a documented retention reason; REVIEW items are bounded and manifested. |
| 10 | Build/test/release discoverability | A developer can find how to build, test, certify, package and recover a release within minutes; structure/path gates enforce this continuously. |

Each dimension must have machine-verifiable evidence where technically possible and explicit human-readable evidence otherwise. “Looks cleaner” is not evidence.

## Closure-grade test strategy

v18.5.1 receives the same seriousness as a Major Closure because repository/package movement can silently invalidate build, test, runtime and release assumptions.

### G0–G4 · baseline / architecture / source integrity

- exact immutable v18.5.0 Stable baseline;
- frozen no-feature structural scope;
- complete inventory/dependency graph;
- repository structure + path-reference gates;
- `gofmt` and `go vet ./...`;
- source-quality/developer-proofing gates;
- bounded-diff review proving structural intent.

### G5–G10 · behavioral equivalence / integration / performance / security / UI

- full `go test -count=1 ./...`;
- selector-existence guards before every focused `go test -run`;
- CGO-disabled fallback;
- focused and bounded/full race coverage appropriate to package moves;
- real ephemeral PostgreSQL integration with mandatory tests failing if skipped;
- SQLite desktop continuity;
- ADR-GDI provider failure/rate limit/stale evidence/source disagreement/blast radius/backpressure/load shedding/fanout/restart/recovery hysteresis/UNKNOWN-ABSTAIN truth;
- v16.11+ performance/capacity gates and supported operating envelopes;
- deterministic equivalence;
- randomized test order;
- Extreme-30;
- functional HTTP workflow;
- Professional Trader/Investor acceptance;
- renderer/browser/responsive/accessibility/role visibility regression;
- security/auth/session/adversarial authorization/provider-rights fail-closed qualification;
- approved-scope traceability and adaptive-governance gates;
- requirement-linked regression coverage for COPY-18.5.1-001, SYMBOL-18.5.1-001/002, NAV-18.5.1-001 and RESEARCH-v15.1.0-17-19-REOPENED;
- full inherited + v17 + v18 implementation-reconciliation ledger with no unexplained gaps;
- inventory gate PASS and release-mode reconciliation gate PASS from `v18_5_1_v17_v18_reconciliation_gate.py`.

### G11–G15 · immutable RC / full certification / native runtime / release assurance

- immutable exact-source RC and source fingerprint;
- fresh full certification after freeze;
- macOS Apple Silicon package build and **actual runtime audit**;
- Windows x64 package build and **actual runtime audit**;
- exact source ZIP;
- per-artifact SHA-256 manifest and provenance;
- fresh Stable identity-only promotion followed by exact Stable recertification;
- GitHub Release must visibly contain runnable macOS + Windows ZIPs, source ZIP, checksums and certification evidence before G15 can pass.

### G16 · retrospective / handoff / hygiene

- 10/10 rubric all PASS;
- no temporary branches/workflows left behind;
- repository map and diagrams match actual paths;
- final handoff records exact commit/tag/fingerprint/artifact SHAs and GitHub Release location;
- adaptive retrospective records any structural lesson that should prevent future repo sprawl.

## Native delivery is mandatory again for v18.5.1

v18.5.1 does not inherit v18.5 binaries. It must generate and certify its own:

1. macOS Apple Silicon runnable ZIP;
2. Windows x64 runnable ZIP;
3. exact source ZIP;
4. SHA-256 manifest;
5. provenance/certification evidence;
6. actual GitHub Release listing those assets.

A tag alone is not sufficient.

## Definition of done

v18.5.1 is complete only when:

1. all 10/10 repository-archetype dimensions independently PASS;
2. root contains only intentionally allowed entry-point files/directories;
3. active workflows are current, minimal and purpose-named;
4. tests/tools/config/docs are categorized and discoverable;
5. production package ownership and dependency direction are clear from the tree and diagrams;
6. no dead/duplicate file remains without explicit reason;
7. all canonical G0–G16 applicable gates and closure-grade regression suites pass;
8. no approved product functionality, decision semantics or protected boundary changed;
9. actual macOS Apple Silicon and Windows x64 artifacts pass runtime/provenance audit and are published in the GitHub Release;
10. a developer unfamiliar with DE.PULSE history can locate the relevant subsystem and build/test/release path without scanning the repository root;
11. every escaped-defect group is either completed in the frozen v18.5.1 slice with required proof or explicitly owned/placed in the next v18.x slice;
12. every confirmed implementation-gap group is either completed now with required qualification or explicitly owned/placed by dependency and risk;
13. all 13 orphaned remediation records have current-vs-next placement and required evidence defined;
14. the full inherited + v17 + v18 ledger contains no unexplained, silently dropped, unowned or falsely inherited item;
15. every recurrence records why prior closure missed it and the regression/prevention control that now blocks repetition;
16. `CI-ADAPTIVE-18.5.1-001` is implemented: stale workflows are retired, mandatory evidence cannot be budget-skipped, and planner/failure/cost telemetry is enforced.
