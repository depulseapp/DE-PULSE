# DE.PULSE — Current Handoff

**GitHub source of truth:** `depulseapp/DE-PULSE`  
**Certified Stable:** `v18.8.1-stable`  
**Stable candidate:** `410679ba0d6459f66a44db15a0a55f30741a7c53`  
**Certified fingerprint:** `bfefa3605ab29b4678275936a3e60e45133d0b592b91298551731f6d629a9d92`  
**Certified Stable build ID:** `v18.8.1-stable-20260820`  
**Current branch:** `v18.8.2-development`  
**Current PR:** #59 — single PR, Draft during RC promotion  
**Current issue:** #57 / `ADAPT-FRESHNESS-001 REOPENED`

## v18.8.2 scope

Bounded Market Intelligence reliability repair only. The G0 diagnosis proved the existing live/snapshot allocator already owns SPY/QQQ and the canonical 15-symbol breadth universe, while VIX already uses the canonical special-index path. The escaped defect was canonical quote freshness/recovery accountability for Market Intelligence breadth plus unavailable-vs-zero presentation truth.

No second router, freshness engine, subscription manager or data engine is permitted. Smart Provider Router v2 remains sole routing authority; canonical freshness/recovery + routed refresh remain sole recovery owner; deterministic Day/Swing/Long and No Execution are preserved.

## Implemented product repair

- `engine_core.go`: existing `broadBreadthUniverse` participates in canonical quote freshness scope with deterministic dedupe.
- `renderer/market-intelligence-truth.js`: `DATA DEGRADED` / `UNAVAILABLE` renders Market Tradeability score as `UNAVAILABLE`, while a genuinely evaluated numeric `0/100` remains valid.
- `v18_8_2_market_intelligence_freshness_test.go`: deterministic coverage for allocator ownership, SPY/QQQ protection, individual SPY/QQQ/VIX loss, breadth missing/recovery, stale evidence and 0/15 unavailable truth.
- `tests/renderer/surface_consolidation_test.js`: v18.8.2 renderer truth assertions are integrated into the existing Fast renderer lane; no new workflow/job.

## Evidence already earned before RC identity promotion

Exact product head `5f2d229a9d63780e539705aa6c94cb62b36bf51d`:
- Fast #432 / run `32433235205`: **PASS**.
- Qualified #149 / run `32433851064`: **PASS** across CI/harness, backend full suite, race, randomized order, renderer, Chrome, WebKit and exact-head evidence summary.

These runs prove the bounded implementation but become superseded for merge/release identity once this RC metadata commit changes the branch head. They must not be reused as exact-head G11 evidence.

## Release-candidate promotion in this commit

The same PR #59 is intentionally Draft while the release-capable head is created. This commit aligns:
- `release_identity.json` → v18.8.2 / build `v18.8.2-stable-20260820`, previous Stable `v18.8.1`;
- `VERSION.txt` and `app_bootstrap.go`;
- release-coupled renderer cache/title identity and last-loaded `renderer/release-identity-v18.8.2.js`;
- `release/v18.8.2/release_contract.json`;
- `release/v18.8.2/run_full_certification.sh` for exact-source G12;
- CURRENT Adaptive Build/Process/Delivery overlays and this handoff.

The canonical Release workflow remains unchanged. Because `release_identity.json` now changes in PR #59, a future exact-head merge can trigger the existing single G11–G16 Release path; merging before fresh RC-head Fast + Qualified is prohibited.

## Exactly one next action

Inspect the single automatic CI Fast run on the new v18.8.2 RC head. If and only if Fast passes, mark the **same PR #59** Ready for Review once to trigger one fresh exact-head Qualified run. Do not manually duplicate Fast/Qualified, do not create another branch/PR, and do not merge until the new RC head has both `DE.PULSE/fast-head` and `DE.PULSE/qualified-head` success.

After fresh Qualified passes, the next action is exact-head merge authorization; the merge then enters the existing one-run G11–G16 path. Stable may be claimed only after G12, macOS Apple Silicon + Windows x64 G13/G14 actual packaged-runtime audits, G15 assurance and same-run no-rebuild publication pass.

After v18.8.2 Stable and continuity reconciliation, resume v18.9.0 `ADAPT-TRADEINSIGHT-001` SHADOW integration through Smart Provider Router v2.

## Resume rule

Any ChatGPT account, Claude or other assistant must read `AGENTS.md` / `CLAUDE.md`, `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`, this file, CURRENT Adaptive overlays, `release_identity.json`, `release/v18.8.2/release_contract.json`, both `.depulse-certification/resume/` checkpoints, the v18.8.1 Stable evidence manifest, issue #57, PR #59 and live workflow/branch state before changing source. Never resume from model memory alone.
