# DE.PULSE — Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

**GitHub source of truth:** `depulseapp/DE-PULSE`  
**Certified Stable:** `v18.9.0-stable`  
**Certified Stable candidate:** `9ea81cddae4875ae15d3719ca028519a36c597b6`  
**Certified Stable qualified source:** `9e86b5e731f7a585cc77c1521f3639fc7a208efc`  
**Certified Stable fingerprint:** `a8719090c341c874dbd1279cc31ad98e84075d5701c46a800bf951340780ecb9`  
**Certified Stable build ID:** `v18.9.0-stable-20260821`  
**Release PR:** #62 — merged  
**Completed release scope:** #61 / `ADAPT-TRADEINSIGHT-001` — closed completed  
**Active development branch:** none  
**Master corrective program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Immediate blocker/next patch:** #64 / `ADAPT-RUNTIME-CRASH-001` -> `v18.9.1`.

## v18.9.0 — COMPLETE / IMMUTABLE STABLE

Fast #481 / `32525637987`, Qualified #153 / `32525738828`, merged candidate `9ea81cddae4875ae15d3719ca028519a36c597b6`, and Release #32 / `32526121817` are the release authority. G11-G16, actual macOS Apple Silicon + Windows x64 packaged-runtime audits, G15 assurance, same-run no-rebuild publication and G16 evidence passed. Stable tag: `v18.9.0-stable`. Durable manifest: `release/v18.9.0/stable-evidence-manifest.json`.

## Post-Stable audit direction — issue #65

The architectural foundation remains valid: Smart Provider Router v2 is the sole executable router; canonical freshness/cache/telemetry/persistence/state owners remain authoritative; direct SEC/EDGAR remains authoritative; no TradeInsight-specific parallel system is allowed.

The audit identified corrective/product gaps: TradeInsight Settings/API-key UX, first-success vs coverage-aware fulfillment, canonical company-name identity/presentation, behavior-oriented Market Data Modes, Form 4 enrichment, ticker/company search, movers/ranking evidence, full useful-capability admission and stronger provider-efficiency/adaptive telemetry.

The user explicitly prefers **many small complete builds over heavy builds** to reduce implementation misses. This is now a permanent process rule: one primary responsibility per patch, explicit non-goals, implementation-miss audit before moving on, and no known miss left only in chat.

## Ordered v18.9.x patch train

1. `v18.9.1` — #64 runtime crash corrective ONLY.
2. `v18.9.2` — TradeInsight Settings/API-key UX ONLY.
3. `v18.9.3` — coverage-aware Smart Provider Router v2 core ONLY.
4. `v18.9.4` — canonical company identity + all-desk presentation ONLY.
5. `v18.9.5` — Market Data Modes + capability diagnostics ONLY.
6. `v18.9.6` — TradeInsight SEC Form 4 enrichment ONLY.
7. `v18.9.7` — TradeInsight ticker/company search ONLY.
8. `v18.9.8` — TradeInsight movers/ranking evidence ONLY.
9. `v18.9.9` — remaining useful TradeInsight capability sweep ONLY.
10. `v18.9.10` — provider efficiency + Adaptive Intelligence telemetry ONLY.
11. `v18.9.11` — whole v18.9.x professional closure audit ONLY; no new feature scope.

A later patch may be split further at G0/G1 if it is still too broad. Never merge patches merely to reduce version count.

## Permanent adaptive provider contract

`consumer requirement -> current cache/state -> exact missing coverage -> eligible-provider ranking -> targeted acquisition -> canonical merge/provenance -> coverage re-evaluation -> next provider only if still needed -> synthesized consumer state`

Provider success does not equal consumer completeness. No fixed provider chain is authoritative. Static ordering is at most a prior/tiebreaker. Validation lifecycle (`SHADOW/VALIDATED/...`) is separate from serving role (`PRIMARY/FALLBACK/BACKFILL/ENRICH/CORROBORATE/...`).

## Other continuity truth

Issue #57 / v18.8.1 Market Intelligence escape is now closed completed because `release/v18.8.2/stable-evidence-manifest.json` already records `RESOLVED_IN_V18.8.2_STABLE`. Its regression remains mandatory in final v18.9.x closure.

The real v18.9.0 macOS Apple Silicon `EXC_CRASH (SIGABRT)` remains unresolved under #64. Do not guess root cause or delete `PersonalMarketTerminal` state/API keys as a first troubleshooting step.

## Exactly one next action

Diagnose #64 using the complete macOS crash evidence or deterministic reproduction and freeze the narrow `v18.9.1` G1. Do not create a `v18.9.2` branch until `v18.9.1` is truthfully closed or the crash is proven external/non-product.

## Resume rule

Any ChatGPT account, Codex session, Claude or human maintainer must read `AGENTS.md`, `CLAUDE.md`, `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`, this file, all four `adaptive-governance/CURRENT_ADAPTIVE_*` files, `release_identity.json`, `release/v18.9.0/stable-evidence-manifest.json`, both `.depulse-certification/resume/` checkpoints, issue #65, issue #64 and live GitHub state. GitHub objects and executable evidence outrank chat memory.
