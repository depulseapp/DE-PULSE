# DEC-2026-08-28-001 — Canonical Adaptive Documentation Rebaseline

**Status:** APPROVED / IMPLEMENTED IN GOVERNANCE  
**Date:** 2026-08-28  
**Affects:** Adaptive operating contract, roadmap, build plan, build process, delivery process, current-state projections and resume guidance

## Decision

Adopt one canonical narrative authority for each adaptive concern:

- permanent rules: `governance/ADAPTIVE-OPERATING-CONTRACT.md`;
- product placement: `governance/ROADMAP.md`;
- build planning: `adaptive-governance/ADAPTIVE_BUILD_PLAN.md`;
- build execution: `adaptive-governance/ADAPTIVE_BUILD_PROCESS.md`;
- delivery/release: `adaptive-governance/ADAPTIVE_DELIVERY_PROCESS.md`.

`adaptive-governance/CURRENT_ADAPTIVE_*.md` and `ADAPTIVE_ROADMAP.md` remain thin compatibility/status projections only. Exact current status and the one next action belong in `governance/current-state.json`, the active closure ledger and `handoff/CURRENT.md`.

`adaptive-governance/README.md` is the disposition index for canonical, specialized, projection and historical adaptive files.

## Why

The 2026-08-27 product audit was comprehensively conserved in new audit documents, machine registers and `CURRENT_*` overlays, but the older 400–675-line v18-era canonical build files remained authoritative. The two narrative sets conflicted and the overlays drifted from later HOST dependency progress. A future contributor could follow a stale v18 plan or stale next action even though the audit content existed in GitHub.

The correction makes audit-aligned files authoritative, removes obsolete release-specific narrative from permanent contracts and leaves volatile state in machine/handoff owners.

## Conserved decisions

This rebaseline preserves:

- all certified v18 responsibilities and immutable Stable evidence;
- all ten Executive audit findings and the full audit-risk/coverage/5×5 registers;
- one `SymbolIntelligenceSnapshot`, one shared Opportunity Lifecycle and Watchlist as selected-universe projection;
- deterministic vs registered learned vs LLM synthesis separation;
- Smart Provider Router v2 and direct-authority/Data Health rules;
- U.S. equities/approved ETFs, No Execution and hidden Data Engine boundaries;
- Go modular-monolith/Postgres-outbox/versioned-client direction;
- Development Production Ready versus Commercial/Public Ready separation;
- G0–G16, AIPLC, exact-head CI and artifact provenance.

It does not claim planned gaps are implemented, close the active v19.0 work slice, authorize v19.1, merge/release, or activate Commercial/Public use.

## Supersession

- Supersedes independent roadmap/build/process/delivery authority in `CURRENT_ADAPTIVE_*` files.
- Supersedes obsolete v18.2–v18.10 and old v19/v20 placement embedded in permanent base documents; immutable historical release evidence is retained.
- Supersedes the version-placement portion of the old permanent TDTI section while conserving its product invariants and approved scope.
- Does not supersede specialized permanent contracts unless this decision explicitly identifies a conflicting version-specific placement; conflicts must be reconciled through a later decision.

## Alternatives considered

1. **Keep canonical and CURRENT narratives in parallel.** Rejected because they already drifted and require every contributor to guess precedence.
2. **Delete old and historical files.** Rejected because certified traceability and referenced paths must remain available.
3. **Make handoff the roadmap.** Rejected because handoff is volatile execution state, not permanent intent.
4. **Generate every narrative from JSON.** Deferred; machine generation is useful for status projections, but durable architecture rationale remains clearer as one reviewed narrative.

## Migration and cost

Migration is documentation/gate-only:

- canonical files are rewritten once;
- `CURRENT_*` paths stay valid as projections;
- the Data Health gate checks canonical contracts instead of duplicated projections;
- resume instructions use fully qualified canonical paths;
- historical documents remain in place and are classified, not deleted.

The main cost is a one-time review and exact-head Fast qualification. Product runtime behavior and release numbering do not change.

## Enforcement

Future material changes update the single canonical file and Decision Log. Projection files may contain only machine identity/status links. A gate or review should reject a projection that introduces independent scope, rules or a next action.
