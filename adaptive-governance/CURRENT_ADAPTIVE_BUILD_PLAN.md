# DE.PULSE — Current Adaptive Build Plan

**Immutable Stable:** `v18.9.0-stable`  
**Stable candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Stable fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Active development branch:** none  
**Active PR:** none  
**Immediate blocker:** issue #64 / `ADAPT-RUNTIME-CRASH-001`.

## v18.9.0 closure — G0–G16 PASS

Issue #61 delivered TradeInsight integration through existing canonical owners only. Exact source head `9e86b5e731f7a585cc77c1521f3639fc7a208efc` passed Fast #481 and full Qualified #153. Merged candidate `9ea81cddae4875ae15d3719ca028519a36c597b6` passed Release #32 / `32526121817` through G11–G16, including native macOS Apple Silicon + Windows x64 package/runtime audits, G15 assurance and no-rebuild Stable publication.

Delivered executable scope includes validated Congressional SHADOW evidence, daily adjusted OHLCV fallback/backfill, bounded admission-controlled multi-symbol history and shared provider telemetry. Smart Provider Router v2 remains sole executable routing authority. Direct SEC/EDGAR remains authoritative. TradeInsight Form 4 enrichment, top movers and ticker/company search remain contract-gated until exact production REST contracts are proven.

## Corrective release entry

The user reported a v18.9.0 macOS Apple Silicon `EXC_CRASH (SIGABRT)` / `abort() called` after Stable publication. Issue #63 is superseded by issue #64, which is the only immediate corrective product packet.

G0 must capture the complete `.ips`/symbolized crash evidence or establish deterministic reproduction on the published Stable artifact. G1 then freezes only the proven affected lifecycle surface. Do not guess between startup, renderer/WebView, persistence/bootstrap, runtime lifecycle, shutdown or install path without evidence. Preserve `PersonalMarketTerminal` state/API keys unless corruption is proven and a safe migration/backup path exists.

Normal execution remains one bounded corrective branch → one Draft PR → exact-head Fast → same PR Ready → Qualified → exact-head merge → one canonical G11–G16 release. No duplicate release workflow, no weakened gate, no unrelated feature bundling unless root cause requires it.

## Protected boundaries

Smart Provider Router v2 sole routing authority; canonical freshness/recovery sole freshness owner; existing multi-feed allocator sole subscription owner; BroadSnapshotBroker canonical reuse owner; deterministic Day/Swing/Long; U.S. Equities Processing; GLD/SLV/USO actionable exceptions; No Execution.

## Exactly one next action

Perform issue #64 G0 crash diagnosis from concrete macOS evidence/reproduction and produce the bounded G1 corrective scope before product-source changes.
