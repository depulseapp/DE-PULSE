# DE.PULSE — v18.6.1 Current Reconciliation Baseline

**Baseline:** `v18.6.1-stable`  
**Created:** 2026-08-19  
**Purpose:** establish an honest current-state reconciliation checkpoint after the v18.6.1 Stable release without fabricating or renumbering historical requirement rows.

## Important scope statement

Earlier recovery work referenced a conserved row-by-row requirement ledger (including a historical 296-row snapshot). That exact source ledger was not conclusively located during the post-v18.6.1 repository audit. Therefore this file **does not invent 296 replacement rows** and does not overwrite any historical reconciliation evidence.

When the original immutable ledger is located, fresh reconciliation must preserve every original requirement ID/history and classify each row against current Stable source/evidence.

## Canonical fresh-status vocabulary

Each conserved requirement row must resolve to exactly one current status:

- `FRESH_PASS` — current Stable source/evidence freshly proves the requirement.
- `REOPENED` — previously accepted requirement has a current reproducible gap or regression.
- `NOT_IMPLEMENTED` — applicable requirement remains unimplemented.
- `INTENTIONALLY_SUPERSEDED` — a governed newer contract explicitly replaces the old implementation requirement while preserving intent/traceability.
- `NOT_APPLICABLE` — requirement is no longer applicable under permanent product boundaries, with rationale.
- `ROADMAP_FUTURE_SCOPE` — valid requirement intentionally remains future work and is placed in current roadmap scope.

No row may be silently dropped because a feature appears old, a file is old, or a previous model assumed completion.

## Confirmed current Stable facts

The following are current baseline facts, not substitutes for row-level reconciliation:

1. `v18.6.1-stable` is the immutable Stable baseline and completed G0–G16.
2. U.S. Equities Processing and No Execution remain permanent boundaries.
3. Smart Provider Router v2 remains the sole provider-routing owner.
4. GLD/SLV/USO remain explicit tradable exceptions.
5. Provider → Market Mode treatment remains explicit and provider count alone cannot change Market Mode.
6. v18.6.1 fixed the watchlist `DESKS` global-removal failure and retains exact-membership Undo behavior.
7. Watchlist membership UI uses toggle/pressed semantics and does not use the deprecated `CURRENT` state label.
8. The header notification is centered in the available middle header lane.
9. The canonical CI model is Fast → Qualified → Release, with exact-head Fast/Qualified evidence and same-run no-rebuild Stable publication.
10. Historical release/governance evidence and actively loaded compatibility assets are intentionally preserved unless a later cleanup proof establishes they are unreferenced.

## Current reconciliation workstreams

### Engineering process hardening

Status: `ROADMAP_FUTURE_SCOPE` with Packet A currently in development.

- Impact Planner v2.
- Release Rehearsal.
- current operational roadmap/build/process/delivery overlays.
- reproducible CI dependencies.
- targeted WebKit.
- durable Stable evidence/CI telemetry.

### User-trust defects

Status: requires fresh G0–G3 intake against current source/issues before individual rows can be marked `FRESH_PASS`, `REOPENED` or future scope.

Priority categories include freshness truth, degraded-state truth, focus/refresh behavior, state/membership integrity and visual stability.

### ADR-GDI / runtime reliability

Status: `ROADMAP_FUTURE_SCOPE` for remaining hardening unless a conserved ledger row/source proof shows a completed item. Continue capability health, freshness SLO, provider circuits, backpressure, coalescing, bounded retry, blast-radius and recovery-state work.

### Shared intelligence utility consolidation

Status: requires fresh architecture/source reconciliation. Scanner/Opportunity Radar, session preparation, catalyst lifecycle and Research should converge on shared canonical owners where duplicate internal pipelines remain.

### Renderer maintainability

Status: `ROADMAP_FUTURE_SCOPE`. The active versioned compatibility stack is not classified as junk merely by age. Modularization is incremental and requires proof before duplicate-owner deletion.

### TradeInsight

Status: controlled `SHADOW` / `SECONDARY` roadmap work through Smart Provider Router only, subject to capability/entitlement/rights/freshness/consumer/Market-Mode review.

## Fresh reconciliation procedure

When resuming the conserved row ledger:

1. Locate the actual immutable ledger file/snapshot and record its source SHA.
2. Preserve IDs and historical status; never regenerate IDs from prose.
3. For each row, identify current source owner, current consumer and current release evidence.
4. Apply exactly one fresh-status value from this file.
5. Record source/evidence citations or explicit absence of evidence.
6. Reopen a requirement if current behavior/source contradicts the former pass state.
7. Place valid incomplete work in the current roadmap/build plan.
8. Block v18 Major Closure on any unexplained applicable row.

## File-age rule

GitHub modification age is not reconciliation evidence. An unchanged five-day-old file can be correct and current; a recently touched file can still be obsolete. Reconciliation uses ownership, runtime references, tests, consumers and release evidence—not timestamps.

## Closure criterion

This baseline may be superseded only by a newer reconciliation artifact that either:

- imports the original conserved row ledger with preserved IDs/history and completes fresh current-state classification; or
- explicitly proves a canonical replacement ledger and maps every prior ID without loss.

Until then, this file is the authoritative statement of what is and is not currently known; it must not be misrepresented as a completed 296-row reconciliation.
