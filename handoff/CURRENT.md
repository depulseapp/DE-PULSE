# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS.** GitHub is authoritative; chat history and temporary workspaces are advisory only.

**Release:** `v18.6.1`  
**Active branch:** `main`  
**Certified Stable:** `v18.6.1-stable`  
**Stable target / certified merged candidate:** `42e8432f7530ae39cbfd6ceb0b0bd5f6311dc5cc`  
**Qualified source head:** `5c3fae486f3e4b4a39a0b1d549916aea9e1295fd`  
**Certified source fingerprint:** `b01c14e1d54b736785eab6c03407801c527edd7769ff6f3d41fd4b20dabebd75`  
**Build ID:** `v18.6.1-stable-20260819`  
**Canonical certification/publication run:** `32279232665` (Release #24)  
**Repository:** `depulseapp/DE-PULSE`  
**Status:** STABLE PUBLISHED · G0–G16 PASS · no open release blockers  
**Last updated:** 2026-08-19 America/Vancouver

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, `governance/CI-EFFICIENCY-CONTRACT.md`, this handoff, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, current GitHub PR/issue state and the immutable Stable tag. Never resume from model memory alone.

## v18.6.1 delivered scope

1. **Repository hygiene:** stale release/retry/trigger branches removed; after closure only `main` remains. All PRs and issues were reviewed; issue #12 was closed after Stable proof. Disposable-file patterns such as `.DS_Store`, `.bak`, `.tmp`, `.swp`, `Thumbs.db`, `_invalid`, `__pycache__`, `.pyc` and `.orig` are absent. Historical release evidence, governance and still-loaded compatibility assets remain intentionally preserved.
2. **Watchlist membership UI:** Day/Swing/Long membership uses pressed-state toggles (`aria-pressed`) only; deprecated `CURRENT` / current-desk wording is absent.
3. **Global symbol removal:** fixed `Global Remove Failed · Can't find variable: DESKS` by loading `renderer/watchlist-desk-contract-v18.6.1.js` before the watchlist extension. Desk-row and Master Market Store removal share the canonical global-remove path with exact-membership Undo.
4. **Extensive watchlist regression:** `release/v18.6.1/browser_watchlist_global_remove_test.py` covers all seven legal non-empty membership combinations, Undo, Master Market Store removal, rapid/double activation, backend-failure preservation, toggle semantics and absence of `CURRENT`.
5. **Header alert:** centered in the available middle header lane between the DE.PULSE identity and runtime/account controls.
6. **Release-harness hardening learned during the trial:** inherited v18.6 UI hierarchy proof is patch-line compatible; publication uses a robust missing-tag lookup; changes to canonical `release.yml` can enter the governed release path without fake product-identity edits.

Permanent boundaries remain unchanged: U.S. Equities Processing, No Execution, deterministic Day/Swing/Long formulas, Smart Provider Router ownership, provider/Market-Mode rules, GLD/SLV/USO tradable exceptions and Adaptive Intelligence governance.

## Final evidence

- Final exact-head Fast: run `32279055139` · PASS.
- Final exact-head Qualified: run `32279113032` · PASS using the intended process-only CI-harness portability lane.
- Product candidate full Qualified: run `32276304863` · backend/race/randomized, renderer, deterministic 2403/2403 and browser/global-remove edge coverage PASS.
- Canonical Release #24: run `32279232665` · G11, G12, macOS Apple Silicon G13/G14, Windows x64 G13/G14, G15, same-run no-rebuild publication and G16 all PASS.
- Stable tag: `v18.6.1-stable` → certified candidate `42e8432f7530ae39cbfd6ceb0b0bd5f6311dc5cc`.
- Earlier Release #22 exposed and stopped on an inherited hard-coded patch-version test; Release #23 passed G11–G15 and artifact verification but stopped on a missing-tag publication lookup. Both defects were repaired through governed source changes before Release #24 passed. Neither failed attempt published Stable.

## CI-efficiency trial result

The event-amplification problem is fixed: no trigger/retry/certification/promotion branches were manufactured, exact-head evidence was enforced, main pushes perform hygiene rather than product Fast testing, and the final process-only publication correction skipped unnecessary product suites. The trial still needed two small recovery PRs because it surfaced previously latent release-harness defects. Those were real defects, not duplicate CI events.

Measured GitHub build clock: first v18.6.1 Draft PR created `2026-08-19T15:11:42Z`; final Stable Release #24 completed `2026-08-19T17:08:02Z` = **1h 56m 20s end-to-end**. Final canonical G11–G16 run itself took **6m 43s** (`17:01:19Z` → `17:08:02Z`).

## Exactly one next action

When the next product slice is requested, start intake from immutable `v18.6.1-stable` / current `main` using one version-development branch and one Draft PR under `governance/CI-EFFICIENCY-CONTRACT.md`. Do not create trigger/retry/certification/promotion branches.
