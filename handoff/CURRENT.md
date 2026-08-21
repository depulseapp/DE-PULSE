# DE.PULSE — Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

**GitHub source of truth:** `depulseapp/DE-PULSE`  
**Certified Stable:** `v18.9.0-stable`  
**Certified Stable candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Certified Stable qualified source:** `9e86b5e731f7a585cc77c1521f3639fc7a208efc`  
**Certified Stable fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Certified Stable build ID:** `v18.9.0-stable-20260821`  
**Release PR:** #62 — merged  
**Completed release scope:** issue #61 / `ADAPT-TRADEINSIGHT-001` — closed completed  
**Active development branch:** none  
**Open corrective blocker:** issue #64 / `ADAPT-RUNTIME-CRASH-001`

## v18.9.0 — COMPLETE / STABLE

v18.9.0 was qualified and published through the single canonical G0–G16 process:
- Fast #481 / run `32525637987`: PASS on exact source head `9e86b5e731f7a585cc77c1521f3639fc7a208efc`;
- Qualified #153 / run `32525738828`: PASS across backend/full/race/randomized, renderer, Chrome, WebKit, CI/provenance and exact-head evidence;
- merged certified candidate `9ea81cddae4875ae15d3719ca028519a36c597b6`;
- Release #32 / run `32526121817`: G11–G16 PASS;
- macOS Apple Silicon and Windows x64 actual packaged-runtime audits: PASS;
- G15 release assurance and exact same-run no-rebuild Stable publication: PASS;
- G16 durable handoff artifact: PASS;
- Stable tag: `v18.9.0-stable`;
- durable manifest: `release/v18.9.0/stable-evidence-manifest.json`.

The immutable Stable tag/candidate/fingerprint remain the release authority even if later fingerprint-excluded continuity metadata advances `main`.

## v18.9.0 delivered scope

- Smart Provider Router v2 remains the sole executable routing authority.
- TradeInsight daily adjusted OHLCV is admitted only through the canonical Historical Bars owner as fallback/backfill; no intraday capability is claimed.
- Bounded multi-symbol history is SHADOW-admitted as sequential client-side fan-out over the verified per-symbol history endpoint, with canonical normalization/deduplication, VIX exclusion and a 50-symbol ceiling.
- Congressional Trading Intelligence is validated and SHADOW-only inside canonical Research alternative evidence after direct SEC refresh.
- Optional TradeInsight failure does not downgrade healthy canonical readiness.
- Shared provider telemetry, freshness, cache, persistence, corporate-action and canonical state owners are reused.
- No second router, scanner, scheduler, Market Mode engine, SEC truth owner, symbol authority or persistence subsystem was introduced.
- Direct SEC/EDGAR remains authoritative for Form 4.
- TradeInsight Form 4 enrichment, top movers and ticker/company search remain deliberately contract-gated until exact executable production REST contracts are independently verified; this is an explicit disposition, not forgotten v18.9.0 work.
- Deterministic Day/Swing/Long truth, U.S. Equities Processing, GLD/SLV/USO actionable exceptions and No Execution remain preserved.

## Post-Stable runtime escape

After v18.9.0 publication, the user reported a real macOS Apple Silicon crash on bundle version `18900`: `EXC_CRASH (SIGABRT)` / signal 6 / `abort() called`.

Version-scoped issue #63 is closed as superseded, **not fixed**. Corrective ownership is issue #64 / `ADAPT-RUNTIME-CRASH-001`. The screenshot proves the v18.9.0 crash class and native ARM64 identity but not a symbolized root cause, so future work must obtain the full `.ips`/backtrace or reproduce before changing product code. Do not delete `PersonalMarketTerminal` state/API keys as a first troubleshooting step.

## Exactly one next action

Diagnose issue #64 from concrete crash evidence/reproduction and produce the bounded corrective-release G0/G1 proposal before optional provider/Market-Mode expansion or unrelated feature work.

## Resume rule

Any ChatGPT account, Codex session, Claude or human maintainer must read `AGENTS.md`, `CLAUDE.md`, `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`, this file, `release_identity.json`, `release/v18.9.0/release_contract.json`, `release/v18.9.0/stable-evidence-manifest.json`, both `.depulse-certification/resume/` checkpoints, issue #64 and live GitHub state. GitHub objects and executable evidence outrank chat memory. No upload of an old chat handoff is required.
