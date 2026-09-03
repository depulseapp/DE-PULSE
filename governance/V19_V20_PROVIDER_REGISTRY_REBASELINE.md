# DE.PULSE v19/v20 Adaptive Provider Registry Rebaseline Addendum

**Status:** AUTHORITATIVE ADDITIVE REBASELINE  
**Baseline:** `v18.10.0` remains immutable  
**Active version remains:** `v19.0.0 — Hosted Trust & Identity Foundation`  
**Active implementation branch/PR remain:** `adapt-hosted-trust-foundation-001` / PR #149  
**Primary future implementation owner:** `v19.1.0 — Canonical Data Runtime & Global Symbol Processing` / issue #153  
**Permanent contract:** `adaptive-governance/ADAPTIVE_PROVIDER_REGISTRY_CONTRACT.md`

This addendum rebaselines provider onboarding architecture after the TradeInsight implementation lessons and the decision to add Market Data (`marketdata.app`). It does not reopen v18.10.0, change the active v19.0 dependency order, start v19.1 early, create a new Router, or create a provider-specific public version.

## Audit conclusion

Existing DE.PULSE architecture is directionally correct: Smart Provider Router v2 is already the sole general routing/admission owner and canonical Data Health/lifecycle/cache/persistence/subscription owners already exist. The gap is **provider onboarding standardization**.

Current integrations still require too much provider-specific admission/configuration knowledge. The future-proof architecture must make the provider adapter the only provider-specific implementation boundary. Once an adapter is registered, capability discovery, entitlement changes, eligibility, Router participation, Data Health, diagnostics and consumer reuse must be generic.

Therefore:

`provider-specific API -> one adapter -> Adaptive Provider Registry -> Smart Provider Router v2 -> canonical state -> all useful consumers`

replaces any residual pattern of:

`provider-specific API -> provider-specific feature wiring / provider-specific Settings / provider-specific consumer routing`.

## Rebaseline decision — no version renumbering required

The existing v19/v20 sequence remains valid. This decision strengthens already-planned responsibilities rather than adding a new release.

### v19.0.0 — unchanged active scope
Finish HOST-001..023 Development Production Ready. Do not implement the new Market Data/provider-registry runtime while HOST-007 and later v19.0 dependency bands remain open.

The v19.0 provider-rights/secrets/scorecard foundations are prerequisites used later by the Registry, but this addendum does not expand v19.0 G1 with Market Data product behavior.

### v19.1.0 — strengthened scope
Existing #153 Smart Provider Router v2 adoption/cross-integration is expanded to include:
- Adaptive Provider Registry canonical runtime owner;
- standard provider adapter/manifest contract;
- self-registration of provider adapters;
- automatic capability/configuration/entitlement/freshness/history/quota probing where technically observable;
- generic eligible-provider projection consumed by Smart Provider Router v2;
- no provider-name routing in product consumers;
- generic Provider Settings metadata contract;
- Market Data as the first new provider implemented through this pattern;
- Market Data API-key/token configuration through existing canonical Settings/secrets owners;
- Market Data provider test/redaction/clear/preserve/replace flow modeled after the successful TradeInsight Settings work;
- Market Data effective-capability discovery so trial/free/paid changes do not require application code changes solely for entitlement changes;
- Market Data SHADOW-first quality/usefulness validation;
- capability-specific cross-integration to applicable Dashboard/Desks/Tradeability/Market Intelligence/Research/Discovery/Global Symbols/Prep/alerts/history/options/Data Health/Maintenance consumers;
- provider-capability matrix/fetch-path/rights/authority updates;
- recovery/downgrade/reprobe behavior through existing Router/Data Health owners;
- cross-provider diagnostics showing configured vs eligible vs serving provider truth;
- no automatic lifecycle/authority promotion.

### v19.2.0 — unchanged architectural role
Hosted Gateway/shared serving consumes the same Registry/Router output and provider-rights/secret owners. It must not create another hosted provider registry or provider router.

### v19.3.0 — strengthened cross-platform/role obligation
Provider configuration/admin presentation across required Mac/Windows/Web surfaces must consume the generic provider Settings metadata contract. Do not implement separate MarketData-specific cross-platform UX stacks.

Ordinary hosted commercial users remain zero-key. Owner/developer provider-secret administration is capability/role controlled and server-managed in hosted mode.

### v19.6.1 — strengthened provider economics/readiness closure
Market Data and every registered provider participate in the same provider reliability/coverage/economics/readiness scorecard, including call/credit pressure, latency, freshness, failures, fallback value, cache value and useful unique coverage where measurable.

### v20.5.0 — clarified adaptive operations role
Adaptive Operations may learn bounded provider utility/cost/reliability priors from point-in-time evidence and may recommend or influence routing within governed Router policy. It cannot become a parallel router, grant lifecycle authority, infer rights or autonomously promote SHADOW providers to production authority.

## Market Data concrete adoption

Market Data is intentionally the first provider used to prove the new onboarding contract.

Current official vendor docs audited 2026-08-26 establish implementation inputs that must be treated as changeable vendor semantics, not permanent hard-coded truth:
- Bearer token authentication;
- official SDK environment variable `MARKETDATA_TOKEN`;
- HTTP 203 may represent successful cache-tier fulfillment and must be accepted alongside 200;
- current free trials provide delayed data even when based on a paid plan;
- current Trader Trial advertises 100,000 daily credits and converts to Free Forever if the user does not subscribe;
- paid Trader currently unlocks real-time stock/options capabilities subject to endpoint-level freshness.

Runtime rule: **effective capability wins over hard-coded plan name**. If a configured account later becomes paid, a fresh capability/entitlement probe updates technical eligibility automatically. Smart Provider Router v2 can then consider the new capability if lifecycle, rights, health, freshness, quota and policy permit. No DE.PULSE code release should be required solely because the account tier changed.

## TradeInsight learning conserved

The following successful TradeInsight integration lessons become generic requirements for Market Data and future token/API-key providers:
- existing Data Providers Settings surface;
- canonical secrets owner, no parallel credential store;
- stored secret never returned to renderer/public state;
- blank save preserves existing secret;
- non-empty save replaces;
- explicit clear through canonical clear-secret owner;
- environment fallback for headless/CI/developer operation;
- shared provider-test path;
- Router configuration recognizes the canonical configured secret;
- provider lifecycle stays capability-scoped;
- no provider-specific Router/cache/freshness/persistence/telemetry system;
- live smoke is bounded/explicit and secrets never enter CI logs;
- Settings reports concise configuration truth while Maintenance/Data Health owns detailed runtime truth.

## Smart/AI/LLM-style provider behavior

Provider adaptation is evidence-driven, not vendor-hard-coded.

Automatically observable evidence should include where applicable:
- capabilities and entitlement changes;
- current freshness class;
- history depth/completeness;
- quota/credits and headroom;
- success/failure/auth/rate-limit patterns;
- latency distributions;
- availability/circuit/recovery behavior;
- cache effectiveness;
- schema integrity;
- disagreement/corroboration quality;
- calls avoided;
- consumer usefulness;
- cost/economics;
- outcome usefulness once point-in-time learning exists.

Use this evidence to recalculate runtime eligibility, selection and fallback within Smart Provider Router v2 and to produce lifecycle-promotion recommendations. Do not silently self-promote provider authority.

## Cross-integration is mandatory

Every provider capability receives `REQUIRED / CONDITIONAL / NOT_USEFUL` dispositions for applicable canonical consumers. Market Data cannot be declared integrated after only a successful API test.

At minimum audit:
- Dashboard/ticker;
- Day/Swing/Long Desks;
- Tradeability/Readiness;
- Market Intelligence / Market Regime contribution;
- Research;
- Discovery / Opportunity Radar;
- Global Symbols/watchlist processing;
- Prep;
- catalysts/rapid-move/volume-volatility alerts;
- fundamentals/history/intraday bars;
- options capabilities;
- Data Health/Maintenance;
- persistence/cache/reconciliation;
- future Outcome Learning/provider usefulness.

Consumers use canonical routed state. No page-specific provider fetch is permitted without an explicit direct-authority justification.

## Acceptance for v19.1 entry/closure

v19.1 G1-G3 must treat this addendum and `ADAPTIVE_PROVIDER_REGISTRY_CONTRACT.md` as mandatory inputs. #153 is the primary requirement owner.

v19.1 cannot close the provider-registry portion until evidence proves:
1. provider adapters self-register through one generic Registry contract;
2. Registry feeds Smart Provider Router v2 rather than replacing it;
3. Market Data is implemented using the generic contract;
4. secure Settings/API-key UX follows the canonical TradeInsight-derived secret pattern;
5. trial/delayed capability is never represented as live;
6. effective entitlement changes can alter technical eligibility without code changes;
7. automatic runtime downgrade/fallback/recovery works;
8. lifecycle/authority promotion remains governed;
9. cross-integration uses canonical state across useful consumers;
10. Data Health/Maintenance/Settings agree on configured, eligible, serving, delayed/live and degraded/recovery truth;
11. no secret leakage, duplicate acquisition, second provider health system or second Router exists;
12. U.S. Equities Processing, GLD/SLV/USO actionable exceptions, direct SEC authority, point-in-time/no-lookahead and No Execution remain intact.

## Resume / continuity rule

Any next assistant/account/AI resuming DE.PULSE must treat this addendum as durable approved future scope. It does not supersede the current `handoff/CURRENT.md` next action: finish v19.0 first. When v19.1 begins, read this addendum and `adaptive-governance/ADAPTIVE_PROVIDER_REGISTRY_CONTRACT.md` before freezing G1 for #153/provider work.
