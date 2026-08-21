# DE.PULSE — Current Adaptive Roadmap

**Certified Stable:** `v18.9.0-stable`  
**Certified candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Certified fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Build ID:** `v18.9.0-stable-20260821`  
**Immediate blocker:** issue #64 / `ADAPT-RUNTIME-CRASH-001`.

## v18.9.0 — COMPLETE / STABLE

Issue #61 / `ADAPT-TRADEINSIGHT-001` is closed completed. Exact source head `9e86b5e731f7a585cc77c1521f3639fc7a208efc` passed Fast #481 / `32525637987` and Qualified #153 / `32525738828`. Merged candidate `9ea81cddae4875ae15d3719ca028519a36c597b6` then passed canonical Release #32 / `32526121817` through G11–G16, including full G12 certification, actual macOS Apple Silicon and Windows x64 packaged-runtime audits, G15 assurance, no-rebuild publication and G16 evidence. Durable release evidence is `release/v18.9.0/stable-evidence-manifest.json`.

TradeInsight delivered bounded SHADOW/canonical-owner integration for Congressional Trading Intelligence, daily adjusted OHLCV fallback/backfill, selective admission-controlled multi-symbol history and shared provider telemetry. Smart Provider Router v2 remains sole routing authority. Direct SEC/EDGAR remains authoritative. Form 4 enrichment, top movers and ticker search remain explicitly contract-gated until exact production REST contracts are proven.

## Post-Stable corrective priority

A real user-observed macOS Apple Silicon v18.9.0 crash (`EXC_CRASH (SIGABRT)` / `abort() called`) escaped certification. Version-scoped issue #63 is closed as superseded; corrective ownership is issue #64 / `ADAPT-RUNTIME-CRASH-001`. The next engineering work must establish the exact root cause from a full `.ips`/symbolized backtrace or deterministic reproduction before changing runtime code. Do not erase user state/API-key continuity as a first troubleshooting action.

## Following sequence

1. Corrective runtime reliability for issue #64.
2. Provider Intelligence & Market-Mode hardening after the crash blocker is resolved or proven external/non-product.
3. v18 Major Closure Candidate: zero unexplained carry-forward, zero orphan useful provider capability, zero duplicate routing owner.
4. v19 Professional Data Infrastructure.
5. v20 Adaptive Intelligence maturation.

Permanent constraints: U.S. Equities Processing, No Execution, G0–G16 only, Smart Provider Router v2 sole routing owner, BroadSnapshotBroker canonical reuse owner, GLD/SLV/USO actionable tradable exceptions and deterministic Day/Swing/Long truth protected.

## Exactly one next action

Diagnose issue #64 from concrete macOS crash evidence/reproduction and freeze a bounded corrective-release G1 before any unrelated feature expansion.
