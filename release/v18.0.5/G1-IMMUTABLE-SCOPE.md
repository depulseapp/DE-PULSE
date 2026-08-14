# DE.PULSE v18.0.5 — G1 Immutable Scope

Status: FROZEN
Baseline: v18.0.4 Stable
Baseline commit: 9b725526d2749e14ec3fcd41fd64ca75eed0770d
Release: v18.0.5 — UI/UX + Symbol Management Hardening

## Governing contracts

- G0–G16 Adaptive Operating Model only.
- REUSE → CONSOLIDATE → REFACTOR → DELETE/REPLACE → ADD.
- Preserve the No Execution Boundary and all permanent adaptive-intelligence, data-utility/correlation, point-in-time truth, performance/scalability, security, multi-platform and UI/UX contracts.
- Do not alter the immutable v18.0.4 Stable release.
- Do not mix v18.1 multi-user scope into this patch.

## Approved implementation scope

### A. Tracked Symbols redesign
- Replace implementation-facing symbol terminology with user-facing Tracked Symbols where appropriate.
- Explain that tracked symbols feed Day / Swing / Long-Term views.
- Clean ticker input and Add Symbol action.
- Compact symbol chips and visible symbol count.
- Quiet danger-style Remove All only when symbols exist.
- Remove user-facing implementation terms such as Master Symbol Store and Add to All Desks.

### B. Canonical symbol ownership + Remove All defect
- One canonical symbol owner; desk lists must not independently mutate canonical state.
- Valid ticker add.
- Deterministic duplicate handling.
- Invalid ticker handling.
- Propagation to all relevant desks.
- Single remove propagation.
- Remove All clears canonical state.
- Day / Swing / Long update immediately.
- Persistence survives reload and app restart.
- No stale rehydration from caches or secondary lists.
- Rapid add/remove actions remain deterministic.
- Empty-state Remove All is safe.

### C. Opportunity Radar cleanup
- Headings and hierarchy.
- Column consistency and numeric alignment.
- Row spacing and status pills.
- Research / Stage actions.
- Score, multiple and percent context.
- Responsive behavior, wrapping and clipping.

### D. Compact Research Target only
- Reduce card height and spacing for Research Target only.
- Compact ticker selection/add and Research Freshness.
- Smaller aligned actions.
- CURRENT status pill.
- Concise Opened from Dashboard context row.
- Narrow-width reflow.
- Do not apply compact-everything globally.

### E. Hide internal machinery from USER / DEMO
- Hide/remove exposed Data Engine and raw provider plumbing.
- Hide subscription/capacity, queue/cache/scheduler/database and circuit/fallback internals from normal-user surfaces.
- Translate technical state into conclusions, confidence/freshness and material degraded-data warnings.
- Preserve deep diagnostics in Maintenance / privileged drill-down where appropriate.

### F. Responsive/layout hardening
- MacBook widths.
- Windows desktop widths.
- Larger screens.
- iPad/narrow browser widths where applicable.
- No overlap, clipped controls, unreadable compression or inconsistent hierarchy.

## Repository hygiene carried from G0

- `De-Pulse-v18.0.4-TEST-Source.zip` is obsolete release debris still present in the baseline repository root.
- Preserve v18.0.4 provenance first; remove the obsolete TEST package through the v18.0.5 release path, not by rewriting historical Stable provenance.

## Explicitly out of scope

- v18.1 per-user/My Market Symbols architecture.
- Hosted/PostgreSQL architecture.
- New provider/router or Rapid Move production scope; those remain the next authorized foundation stage after v18.0.5.
- Changes to protected deterministic Day/Swing/Long formulas.
- Any trading execution, paper trading, portfolio/P&L, journal or autonomous trading capability.
