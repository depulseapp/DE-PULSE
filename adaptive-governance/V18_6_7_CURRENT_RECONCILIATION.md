# DE.PULSE — v18.6.7 Current Reconciliation

**Current engineering slice:** `v18.6.7-development`  
**Immutable Stable baseline:** `v18.6.1-stable`  
**Conserved authority ledger:** `release/v18.5.1/V17-V18-IMPLEMENTATION-RECONCILIATION.json`  
**Ledger blob at v18.6.7 intake:** `2a32b3f93203d61b1aca55172530652d736bbf55`  
**Declared conserved rows:** `296`  
**Purpose:** supersede the post-v18.6.1 operational uncertainty about whether the historical row ledger could be located, while preserving the historical ledger and every original requirement ID/status as evidence rather than pretending the v18.5.1 snapshot is current truth.

## 1. What changed since the v18.6.1 baseline reconciliation

`adaptive-governance/V18_6_1_CURRENT_RECONCILIATION.md` correctly refused to invent a replacement 296-row ledger when the exact source had not yet been located during that audit. The source has now been located and verified in the repository at:

`release/v18.5.1/V17-V18-IMPLEMENTATION-RECONCILIATION.json`

That discovery resolves the **location uncertainty only**. It does not make the old v18.5.1 row statuses current.

The historical ledger's own policy is explicit: no historical PASS/status is current evidence. Each approved item must map to current source/behavior, owner, regression/runtime/browser evidence where applicable, native evidence where applicable, and an explicit current disposition.

## 2. Conservation rule

The 296 figure means **296 tracked authority records**, not 296 open defects and not 296 unfinished features.

The existing `v18_5_1_v17_v18_reconciliation_gate.py` remains the authoritative structural/conservation parser for this ledger. In inventory mode it verifies, among other things:

- ledger schema/release identity;
- all 48 inherited approved-scope items;
- the frozen v17 major scope and closure matrix;
- v18 workstreams/delivery slices;
- all listed v18 release-scope entries;
- the 13 functionality-utility carry-forward rows;
- required known blocker IDs;
- the declared tracked-row count against the reconstructed tracked rows;
- non-empty immutable IDs and duplicate-ID absence;
- status-vocabulary integrity.

For v18.6.7, Fast executes this gate in **inventory mode**. Inventory mode proves conservation; it does not grant current `FRESH_PASS` to historical rows.

## 3. Current fresh-status vocabulary

For current reconciliation, use the operational vocabulary already established after v18.6.1:

- `FRESH_PASS` — current source/evidence freshly proves the requirement.
- `REOPENED` — previously accepted requirement has a current reproducible gap or regression.
- `NOT_IMPLEMENTED` — applicable requirement remains unimplemented.
- `INTENTIONALLY_SUPERSEDED` — a governed newer contract replaces the old implementation requirement while preserving intent/traceability.
- `NOT_APPLICABLE` — requirement is no longer applicable under permanent product boundaries, with rationale.
- `ROADMAP_FUTURE_SCOPE` — valid requirement intentionally remains future work and is placed in current roadmap scope.

The historical ledger may use older status labels such as `REVALIDATION_REQUIRED`, `PLACED_NEXT_V18` or `ROADMAP_PLACED_FUTURE`. Those values remain preserved inside the historical artifact. A fresh v18.6.7 disposition is a separate current assessment, not an in-place rewrite of historical evidence.

## 4. Fresh current facts already proven after the historical snapshot

The following are current facts with post-v18.5.1 evidence. They are **not automatically assigned to old row IDs by title similarity**; row-level mapping must preserve exact IDs and evidence.

1. `v18.6.1-stable` completed the governed Stable path and remains the immutable application Stable baseline.
2. v18.6.1 fixed the watchlist global-remove `DESKS` failure, preserves exact-membership Undo, removes the deprecated `CURRENT` membership label, and centers the header notification in the available middle lane.
3. Phase 0 Packet A / PR #46 completed Impact Planner v2, Release Rehearsal and current operating overlays.
4. Packet B / PR #47 completed Fast/Qualified reproducibility hardening and immutable Action pins for those two workflows.
5. Packet C / PR #48 established Chrome + WebKit as co-primary browser engines with targeted, cost-aware execution.
6. Packet D / PR #49 established durable Stable evidence indexing plus CI runtime/queue/platform/browser-setup telemetry and amplification warnings.
7. Packet E / PR #50 established the first capability-oriented renderer owner for Documentation with explicit legacy fallback and direct Chrome + WebKit proof.
8. Phase 0 A–E is therefore complete on current `main`; older documents describing one of those packets as active are operationally superseded by the current roadmap/build-plan/handoff overlays.

## 5. Current v18.6.7 reconciliation work

### A. Legacy test/gate hygiene

The version-stacked root history is now treated as a traceability/ownership problem, not as junk by age.

- all targeted root executable tests/gates are dynamically inventoried;
- Go `*_test.go` remains active through `go test ./...` unless deliberately migrated;
- explicit current CI/certification consumers make a versioned executable `ACTIVE_REQUIRED`;
- unreferenced executable evidence defaults to `UNREFERENCED_USEFUL`, never automatically `SAFE_TO_REMOVE`;
- six low-risk current v18.6 test files have a first capability-oriented rename/move wave with byte-for-byte test-body preservation;
- organized renderer tests are routed as `RENDERER_UI`, so Chrome + WebKit protection cannot be silently lost because files moved under `tests/renderer/` or `tests/browser/`.

### B. Row-level fresh reconciliation

Do not rewrite all 296 historical rows to a guessed modern status in one mechanical pass. The current process is:

1. conserve every row/ID through the historical ledger gate;
2. prioritize rows related to still-reproducible user-trust, runtime/degradation, provider, shared-intelligence and renderer risks;
3. map exact ID → current source owner → current consumer → current regression/runtime/browser/native evidence as applicable;
4. assign one current fresh status only when evidence supports it;
5. carry unresolved applicable rows forward explicitly rather than hiding them behind historical PASS or release age.

## 6. G0–G3 direction for the next product slice

The provisional next product target remains **v18.7.0 Runtime Reliability & Data Truth**, unless this fresh reconciliation identifies a still-reproducible higher-severity user-trust defect that should be bundled ahead of or into that slice.

Candidate v18.7.0 themes:

- exact `DATA DEGRADED` semantics and blast-radius isolation;
- capability/provider health aggregation;
- freshness SLO truth;
- duplicate-work suppression / single-flight / request coalescing;
- bounded retries and provider circuits;
- queue/backpressure/load shedding and runtime overload behavior;
- disagreement handling and recovery hysteresis;
- `UNKNOWN` / `ABSTAIN` behavior when evidence is insufficient;
- realistic active-market load/latency evidence.

The scope becomes immutable only at G1 after current-source verification shows which controls already exist and which gaps remain.

## 7. v18 closure rule

The located ledger removes the excuse for unexplained requirement loss. v18 Major Closure is blocked until every applicable conserved row has an explicit current disposition and the closure evidence required by its risk is fresh.

A historical row may remain preserved forever as history; what may not remain is an unexplained applicable requirement with no current disposition.

## 8. Authority order

For current operational decisions use, in order:

1. immutable `v18.6.1-stable` release evidence;
2. current `main` source and exact-head CI evidence;
3. this v18.6.7 reconciliation overlay;
4. `CURRENT_ADAPTIVE_ROADMAP.md`, `CURRENT_ADAPTIVE_BUILD_PLAN.md`, current process/delivery overlays and `handoff/CURRENT.md`;
5. the conserved v18.5.1 ledger for immutable IDs/history;
6. older reconciliation/audit files as historical evidence.

Do not overwrite the historical ledger merely to make its old status labels look current.
