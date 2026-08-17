# DE.PULSE v18.5.1 TEST — Documentation Impact

**Current candidate:** v18.5.1 TEST  
**Immediate Stable predecessor:** v18.5.0 STABLE  
**Major v18 provenance anchor:** v17.5.1 STABLE  
**Build ID:** `v18.5.1-test-recovery-20260817`  
**Runtime profile:** `PersonalMarketTerminal-v18.5.1-TEST`  
**Status:** DEVELOPMENT COMPLETE / AWAITING QUALIFICATION

This patch does not replace the long-lived User, Developer, or Capabilities & Limitations manuals. Those manuals retain the complete historical release narrative and are refreshed in full at meaningful Stable/material architecture boundaries. This release-specific impact record is the authoritative current-candidate delta while v18.5.1 remains TEST/RC.

## User documentation impact

- Live market updates use incremental DOM reconciliation so same-page updates preserve row identity, hover/focus/selection and scroll rather than destructively remounting the active surface.
- A desk-row **×** means remove the tracked symbol from all Day/Swing/Long memberships. This is intentionally different from a desk-membership pill, where final-membership removal remains protected.
- Undo restores the exact prior desk-membership combination and prior selection context.
- The active desk has a visible non-color-only `CURRENT` state plus `aria-current`; the selected ticker remains a separate state.
- Research Target and header hierarchy are hardened for desktop/tablet/narrow widths. Both complete ET/PT clocks remain visible; market/data/start-stop status is a secondary layer and Market Instruments remains tertiary.
- First-run wording is version-neutral and truthful: compatible prior Stable profile/data can be preserved without claiming an obsolete v17-specific migration.
- v18.5.1 TEST uses an isolated `PersonalMarketTerminal-v18.5.1-TEST` profile cloned from v18.5.0 Stable on first use without writing into Stable.

## Developer documentation impact

- Canonical release identity is `18.5.1` / `TEST`, build `v18.5.1-test-recovery-20260817`.
- `stable_baseline` remains the v18 major-family provenance anchor `v17.5.1`; `previous_stable` is the immediate certified predecessor `v18.5.0`.
- `renderer/live-dom-reconcile.js` owns incremental live-update reconciliation; full structural navigation may still use the full renderer path.
- `renderer/watchlist-v18.5.1.js` reuses the existing global master-symbol remove/restore API instead of creating a second backend removal owner.
- `renderer/header-v18.5.1.js` and `renderer/ui-v18.5.1.css` refine information hierarchy without changing provider/scoring/market logic.
- TEST profile migration is owned by `v18_test_profile.go`; migration excludes transient runtime files and preserves Stable isolation.
- Qualification is checkpoint-based: ordinary development commits are quiet; exact-SHA FAST G0-G5 runs before G6-G10; G11-G15 are reserved for an immutable RC.

## Capabilities & limitations impact

- DE.PULSE remains decision support, not a profit predictor and not an execution product. The permanent **No Execution** boundary is unchanged.
- Actionable processing remains **U.S.-listed only** under the existing eligibility contracts.
- Stale/cached/history-only evidence remains subject to existing freshness/degradation semantics; this patch does not loosen truthfulness or readiness rules.
- The watchlist UI change does not create a portfolio, position, P&L or order-management concept; it only changes tracked-symbol membership semantics.
- macOS Apple Silicon and Windows x64 actual packaged-runtime certification remain required before Stable promotion. No physical native macOS/Windows PASS is claimed by this development record.
- G0-G10 qualification and G11-G15 native/release certification remain pending until executable evidence passes on the exact frozen source.

## Historical continuity retained

The bundled manuals continue to preserve `v18.0.0 TEST`, `v17.5.1 STABLE`, the 30 FULL / 0 PARTIAL / 0 MISSING inherited roadmap evidence, Trade Readiness/freshness truth, Provider Router ownership, physical native macOS/Windows requirements and the permanent No Execution boundary.

## Documentation delivery rule

For TEST/RC patch candidates, this release-specific delta prevents 38–80 KB historical manuals from being redundantly rewritten after every small patch while still making the current candidate fully auditable. At Stable promotion or another material documentation boundary, the validated delta is folded into the long-lived manuals/README/QA release history as appropriate.
