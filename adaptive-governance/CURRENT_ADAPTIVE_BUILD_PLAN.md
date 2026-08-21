# DE.PULSE — Current Adaptive Build Plan

**Immutable Stable:** `v18.8.1-stable`  
**Stable candidate:** `410679ba0d6459f66a44db15a0a55f30741a7c53`  
**Current engineering branch:** `v18.8.2-development`  
**PR:** #59  
**Open defect:** #57 / `ADAPT-FRESHNESS-001 REOPENED`  
**Current line:** `v18.8.2 — Market Intelligence Reliability / RC exact-head requalification`.

## G0–G4 — PASS

G0 proved the existing canonical allocator already owns SPY/QQQ and all 15 breadth symbols; VIX remains canonical special-index demand. G1 froze issue #57 only. G2/G3 preserved the existing owner chain. G4 repaired the narrower quote-freshness accountability gap and unavailable-vs-zero presentation truth without adding routing, freshness or subscription owners.

## G5–G10 — implementation head PASS

Product head `5f2d229a9d63780e539705aa6c94cb62b36bf51d` passed Fast #432 / `32433235205` and Qualified #149 / `32433851064`, including full Go, race, randomized order, renderer, Chrome, WebKit, CI/harness and exact-head status evidence.

That evidence is not reused after release-candidate identity changes the branch head. Exact-head release provenance requires a new Fast + Qualified pair on the RC head.

## RC identity promotion — current commit

The same branch/PR now becomes release-capable by aligning:
- v18.8.2 canonical `release_identity.json`, `VERSION.txt` and app bootstrap;
- release-coupled renderer cache/title identity plus `release-identity-v18.8.2.js`;
- `release/v18.8.2/release_contract.json` and exact-source `run_full_certification.sh`;
- durable handoff/build/delivery overlays.

Stable baseline and previous Stable are both v18.8.1. Build ID is `v18.8.2-stable-20260820`.

## Fresh exact-head qualification required

PR #59 remains Draft through this source-changing RC promotion. The new head must earn:
1. one automatic Fast run;
2. if Fast passes, one Ready-for-Review transition on the same PR;
3. one fresh Qualified run, expected to cover backend/race/randomized, renderer, Chrome and WebKit because the RC changes backend/renderer/release tooling surfaces;
4. no source change after Qualified unless a classified real defect requires it.

Only a head with both `DE.PULSE/fast-head` and `DE.PULSE/qualified-head` success may merge. Release G11 then checks source-head → merged-candidate fingerprint equivalence and runs the existing G12–G16 no-rebuild path.

## Protected boundaries

Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription allocator; BroadSnapshotBroker canonical reuse owner; deterministic Day/Swing/Long; U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution.

## Exactly one next action

Inspect the automatic Fast result on the new RC head. If PASS, mark the same PR #59 Ready exactly once for fresh Qualified. No retry/certification branch, duplicate workflow or manual duplicate run.

After v18.8.2 Stable: v18.9.0 TradeInsight SHADOW → v18.9.1 provider/Market-Mode hardening → v18.10 zero-gap closure.
