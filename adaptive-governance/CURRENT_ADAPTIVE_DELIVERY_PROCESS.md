# DE.PULSE — Current Adaptive Delivery Process

**Certified Stable:** `v18.8.1-stable`  
**Stable candidate:** `410679ba0d6459f66a44db15a0a55f30741a7c53`  
**Next delivery line:** `v18.8.2 — Market Intelligence Reliability`  
**PR:** #59 — single PR, Draft during RC promotion.

v18.8.1 delivery remains complete and immutable. v18.8.2 does not rebuild or redefine the v18.8.1 package.

## v18.8.2 proof already earned

Bounded product head `5f2d229a9d63780e539705aa6c94cb62b36bf51d` passed Fast #432 / `32433235205` and Qualified #149 / `32433851064` across backend/full/race/randomized, renderer, Chrome, WebKit and CI/harness lanes.

The current RC promotion changes canonical package/release identity, so those statuses become historical implementation evidence rather than merge-eligible exact-head evidence. Fresh Fast + Qualified are mandatory on the new RC head.

## Release-capable delivery state

The same PR now carries v18.8.2 `release_identity.json`, version/build identity, renderer release overlay, `release/v18.8.2/release_contract.json` and exact-source G12 script. The canonical Release workflow remains unchanged and can trigger once on PR merge because the PR changes `release_identity.json`.

Normal delivery remains one development branch → one Draft PR → Fast → same PR Ready → Qualified → exact-head merge → one canonical G11–G16 run → exact certified artifacts published without rebuild → repository continuity reconciliation.

G11 must require exact-head Fast + Qualified statuses and source-head/merged-candidate fingerprint equivalence. G12 performs exact-source full certification. G13/G14 must build and runtime-audit macOS Apple Silicon and Windows x64. G15 binds both native evidence graphs. Publication uploads the same certified artifacts without rebuild. G16 records the durable workflow handoff.

No retry/certification/promotion branches, duplicate release workflows or manual duplicate CI runs. G0–G16 only.

## Exactly one next action

Inspect the automatic Fast run on the new RC head. If PASS, mark the same PR #59 Ready exactly once to trigger fresh Qualified. Do not merge before both exact-head statuses are green.

After v18.8.2 Stable and continuity reconciliation: v18.9.0 TradeInsight SHADOW → v18.9.1 provider intelligence → v18.10 zero-gap closure.
