# DE.PULSE — Adaptive Provider Registry / Market Data Handoff Addendum

**Status:** DURABLE APPROVED FUTURE SCOPE  
**Recorded while active version remains:** `v19.0.0 — Hosted Trust & Identity Foundation`  
**Do not start this implementation before v19.0 exit.**

## Decision

The user approved a future-proof evolution of DE.PULSE provider onboarding:

`Provider Adapter -> Adaptive Provider Registry -> capability/entitlement probes -> rights/authority + Data Health -> Smart Provider Router v2 -> canonical state -> all useful consumers`

A provider-specific adapter remains necessary for a new vendor API. After adapter registration, provider discovery, effective capability/entitlement observation, Router eligibility, Data Health, diagnostics and consumer reuse should be generic.

Smart Provider Router v2 remains the sole general routing/admission authority. The Registry is not Router v3 and is not a parallel Router.

## Permanent authoritative files

Read before any v19.1 #153/provider implementation:
1. `adaptive-governance/ADAPTIVE_PROVIDER_REGISTRY_CONTRACT.md`
2. `governance/V19_V20_PROVIDER_REGISTRY_REBASELINE.md`
3. `governance/programs/V19-V20-REBASELINE/adaptive-provider-registry.json`
4. issue #153 — expanded title/body is the primary v19.1 implementation owner.
5. existing #79/#80/#81/#82/#83/#84 provider/Data Health architecture and evidence remain inherited.

## Market Data first adopter

Market Data (`marketdata.app`) is the first new provider required to use the generic provider-registry architecture.

Required implementation includes:
- provider adapter + manifest self-registration;
- `MARKETDATA_TOKEN` environment fallback;
- Bearer-header auth;
- existing Data Providers Settings API-key/token UI using the generic secret pattern learned from TradeInsight #76;
- secret never returned to renderer/public state;
- blank save preserves existing key;
- non-empty save replaces;
- explicit clear through canonical clear-secret owner;
- shared provider-test path;
- no token in URL/logs/telemetry;
- HTTP 200 and 203 treated as successful Market Data responses where current vendor semantics apply;
- current trial/delayed entitlement represented truthfully;
- effective capability probing so future paid/live entitlement can become technically eligible without a DE.PULSE code release solely because the subscription changed;
- SHADOW-first initial lifecycle;
- no automatic lifecycle/authority promotion;
- normal Router scoring by capability/authority/freshness/health/latency/headroom/cost/utility/rights;
- cross-integration to all useful canonical consumers, never provider-specific page fetching.

## Cross-integration rule

Every Market Data and future-provider capability must receive `REQUIRED / CONDITIONAL / NOT_USEFUL` dispositions for applicable:
- Dashboard/ticker;
- Day/Swing/Long Desks;
- Tradeability/Readiness;
- Market Intelligence/Market Regime contribution;
- Research;
- Discovery/Opportunity Radar;
- Global Symbols/watchlist;
- Prep;
- earnings/catalyst/rapid-move/volume-volatility alerts;
- fundamentals/history/intraday;
- options where supported;
- Data Health/Maintenance;
- cache/persistence/reconciliation;
- future Outcome Learning/provider usefulness;
- adaptive processing priority.

Consumers use canonical routed state. Provider-specific UI/page fetches are prohibited unless an explicit direct-authority exception applies.

## Automatic vs governed

Automatic:
- provider self-registration;
- configured-state discovery;
- capability/entitlement/freshness/history/quota probes where observable;
- health/latency/headroom updates;
- technical eligibility recalculation;
- fallback/demotion/cooldown/recovery;
- plan/subscription change reprobe;
- dependent consumer re-evaluation.

Governed only:
- `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION` promotion;
- direct-authority replacement;
- provider public/commercial rights activation;
- protected deterministic logic changes;
- product-boundary changes.

## Version placement

- **v19.0.0:** unchanged active scope; finish HOST-001..023 Development Production Ready first.
- **v19.1.0:** Adaptive Provider Registry runtime + generic adapter contract + Market Data adapter + generic Settings/token integration + entitlement reprobe + Router/cross-integration + existing-provider adoption audit under #153.
- **v19.3.0:** required role-aware Mac/Windows/Web provider Settings/Admin presentation consumes the same metadata contract; no MarketData-specific parallel UX.
- **v19.6.1:** whole-provider reliability/coverage/economics/readiness scorecards include Market Data.
- **v20.5.0:** bounded adaptive provider usefulness/cost priors may mature, but Router/lifecycle/rights authority remains intact.

No new public version was inserted solely for Market Data, and the existing v19/v20 dependency chain is unchanged.

## AI/LLM-style intent

The provider system should behave intelligently by accumulating evidence and adapting safely:
- capability availability;
- freshness;
- latency;
- quota/credits;
- errors/rate limits;
- cache value;
- fallback/recovery quality;
- schema integrity;
- cross-provider disagreement;
- consumer usefulness;
- calls avoided;
- cost/economics;
- later point-in-time outcome usefulness.

This evidence may affect runtime eligibility/selection within Smart Provider Router v2 and create promotion recommendations. It must not silently self-promote provider authority, infer rights, or make AI/LLM/agents a market-truth owner.

## Current sequencing protection

This addendum is future scope only while v19.0 remains active. Resume `handoff/CURRENT.md` first and finish its exact next action. Do not use this provider work to bypass HOST-007 or other v19.0 Development Production Ready blockers.
