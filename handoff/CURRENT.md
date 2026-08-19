# DE.PULSE — Current Authoritative Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS.** GitHub is authoritative; chat history and temporary workspaces are advisory only.

**Release:** `v18.6.1`  
**Active branch:** `v18.6.1-development`  
**Last Certified Stable:** `v18.6.0-stable`  
**Stable promotion commit:** `2abc4a4a3fbbe623aff57948ec875f45e7ef0a1c`  
**Last certified source:** `d375852d846f8c9f0045ac929da1830b85ad629e`  
**Last certified source fingerprint:** `e8c009c16eedb448ed5b9731d8dd24026a7ea0b5a2b5c82e26490a2941b7b4c8`  
**Last canonical certification run:** `32225064225`  
**Repository:** `depulseapp/DE-PULSE`  
**Status:** v18.6.1 implementation complete on the single development branch; exact-head CI qualification and Release are pending.  
**Last updated:** 2026-08-19 America/Vancouver

## Resume rule

Read `AGENTS.md` or `CLAUDE.md`, then `governance/CI-EFFICIENCY-CONTRACT.md`, this handoff, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, the current GitHub PR/check state, and the immutable Stable tag. Never resume from model memory alone. Resume from the last trustworthy PASS recorded in GitHub evidence.

## v18.6.1 committed scope

This is the focused CI-efficiency trial requested after v18.6.0 Stable. Use only `v18.6.1-development` and one PR to `main`.

1. **Repository hygiene:** audit branches, PRs, issues and repository files. At pre-PR intake only `main` and `v18.6.1-development` are needed; no stale trigger/retry/certification branches remain. Historical release evidence, governance and still-loaded compatibility assets are not junk and must not be deleted merely because they are old.
2. **Watchlist membership UI:** Day/Swing/Long membership is represented only by the pressed-state toggle (`aria-pressed`). The deprecated `CURRENT` / current-desk marker remains absent.
3. **Global symbol removal:** repair the production WebKit/Safari failure `Global Remove Failed · Can't find variable: DESKS`. `renderer/watchlist-desk-contract-v18.6.1.js` explicitly owns the canonical `day/swing/long` desk-key binding and loads before the active watchlist extension. Desk-row removal and Master Market Store removal continue through the one canonical global-remove path with exact-membership Undo.
4. **Extensive watchlist regression:** `release/v18.6.1/browser_watchlist_global_remove_test.py` covers all seven legal desk-membership combinations, Undo, Master Market Store removal, rapid double activation, backend-failure behavior, toggle semantics and absence of `CURRENT`. Fast static policy now fails if the production DESKS contract/order or G12 proof disappears.
5. **Header alert acceptance:** the notification bar fills and centers within the middle topbar lane between product identity and runtime/account controls.

No deterministic Day/Swing/Long formula, Smart Provider Router ownership, provider rule, Adaptive Intelligence contract, U.S. Equities Processing Boundary or No Execution Boundary is changed.

## Patch certification model

`release/v18.6.1/patch_contract.json` explicitly inherits the complete v18.6.0 certification/CI matrix. `release/v18.6.1/run_full_certification.sh` executes the full inherited v18.6.0 G12 matrix against the exact v18.6.1 source, then executes the new v18.6.1 browser edge proof. This is an efficiency optimization in metadata/orchestration only; no G0–G16 quality gate is removed.

## CI/process contract

The permanent streamlined lifecycle is:

`v18.6.1-development` → one Draft PR → exact-head **CI Fast** → mark the same PR Ready → exact-head **CI Qualified / G10** → merge the same PR → main-push branch hygiene only → one **Release G11–G16** because `release_identity.json` changed → exact same-run certified macOS Apple Silicon + Windows x64 artifacts → `v18.6.1-stable`.

Never manufacture CI events with trigger/retry/fallback/certification/promotion branches or PRs. Same-source infrastructure failures rerun only failed work; source fixes remain on this branch and PR.

## Last trustworthy PASS

The last fully certified product PASS is immutable `v18.6.0-stable`. v18.6.1 has implementation evidence committed but has not yet earned exact-head Fast/Qualified or G11-G16 PASS. The checkpoints therefore correctly mark source fingerprint and downstream gates pending rather than claiming Stable early.

## Provider-neutral continuation instruction

> Connect to `depulseapp/DE-PULSE`, read `AGENTS.md` or `CLAUDE.md`, `governance/CI-EFFICIENCY-CONTRACT.md`, `handoff/CURRENT.md`, `release_identity.json`, both `.depulse-certification/resume/` checkpoints, current PR/check state and latest Stable release. Preserve G0–G16, U.S. Equities Processing, No Execution, deterministic formulas and Smart Provider Router ownership. Continue v18.6.1 only on `v18.6.1-development` and its single PR; never create CI trigger/retry/certification/promotion branches.

## Exactly one next action

Open the single Draft PR `v18.6.1-development → main` and let CI Fast validate the exact source head. If Fast passes, mark that same PR Ready exactly once to start CI Qualified.
