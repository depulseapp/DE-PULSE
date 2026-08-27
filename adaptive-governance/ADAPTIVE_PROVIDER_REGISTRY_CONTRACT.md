# DE.PULSE — Adaptive Provider Registry & Smart Provider Onboarding Contract

**Status:** PERMANENT / GOVERNING ADDITIVE CONTRACT  
**Applies to:** Adaptive Roadmap, Adaptive Build Plan, Adaptive Build Process, Adaptive Delivery Process  
**Canonical routing authority:** Smart Provider Router v2  
**First concrete adoption under this contract:** Market Data (`marketdata.app`)  
**Primary planned implementation version:** `v19.1.0 — Canonical Data Runtime & Global Symbol Processing` / issue #153

This contract extends the existing Adaptive Data Health / Provider Production contract and Smart Provider Router v2. It does **not** introduce another router, health owner, freshness owner, cache, subscription manager, lifecycle owner, canonical state store, Market Regime engine or provider-specific data engine.

## 1. North Star

DE.PULSE must treat provider change as a normal adaptive operating event rather than an application redesign.

Target architecture:

`Provider Adapter -> Adaptive Provider Registry -> capability/entitlement probes -> rights/authority + Data Health -> Smart Provider Router v2 -> canonical state -> all authorized consumers`

Permanent rule:

> Once a standards-compliant provider adapter is added, the provider must self-register with the Adaptive Provider Registry. Smart Provider Router v2 must automatically discover its currently usable capabilities and consider it wherever those capabilities are requested. Features must not require provider-specific routing wiring.

A completely unknown third-party API cannot be understood magically. A DE.PULSE adapter is still required. The future-proof requirement is that **all work after adapter registration is generic**.

## 2. Canonical ownership

### Adaptive Provider Registry owns
- registered provider identities and adapter manifests;
- declared capability surface and authority class;
- authentication/configuration requirements without storing secret values itself;
- runtime capability and entitlement observations;
- plan/feature/freshness/history/credit/rate-limit observations where discoverable;
- provider lifecycle state reference;
- provider availability/eligibility projection consumed by Smart Provider Router v2;
- adapter/version/schema identity;
- capability probe results and timestamps;
- provider-specific transport quirks only when they are normalization concerns.

### Existing canonical owners remain unchanged
- **Smart Provider Router v2:** sole general routing/admission/selection authority.
- **Data Health/freshness:** sole downstream health/freshness truth owner.
- **Provider lifecycle:** `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`; no automatic authority promotion.
- **Provider rights:** canonical rights/governance/public-production enforcement owner.
- **Secrets:** existing canonical secret/managed-secret owner; Registry never becomes a credential store.
- **Cache/persistence/canonical state:** existing owners.
- **Dynamic Multi-Feed Subscription Manager:** sole upstream live-subscription allocator/owner.
- **Direct source authorities:** remain explicit; direct SEC/EDGAR remains authoritative where already governed, including Form 4.

## 3. Standard Provider Adapter Contract

Every new routable provider adapter must expose enough normalized metadata/behavior for generic registration and routing. Applicable fields include:

- stable provider ID and adapter/schema version;
- authentication type and secret-reference name;
- supported markets/instrument classes;
- capabilities/datasets/endpoints;
- authority class: routable / corroborative / fallback / direct-authority / non-routable;
- real-time/delayed/historical semantics;
- historical depth/completeness;
- entitlement/plan-effective capabilities when discoverable;
- quota/credits/rate-limit headroom when discoverable;
- provider health/status and latency observations;
- cache/batch/stream support;
- provider observation/source-time semantics;
- cost class/economics when known;
- rights/licensing metadata reference;
- capability-specific probe functions;
- normalized fetch functions for supported canonical capabilities;
- deterministic unsupported/withheld/unentitled states;
- redacted diagnostics only; never secret values.

Consumers request **capabilities**, never provider names. Example:

`OPTIONS_CHAIN_REALTIME -> Registry eligible candidates -> Smart Provider Router v2 -> selected provider -> canonical state -> consumers`

## 4. Automatic vs governed behavior

### Automatic behavior is required for
- adapter discovery/registration;
- configuration-presence detection;
- capability probing;
- effective entitlement/freshness/history/quota discovery when technically observable;
- provider health/latency/headroom observations;
- candidate eligibility recalculation;
- runtime fallback/demotion/cooldown/suppression/recovery through existing Router/Data Health owners;
- capability re-evaluation after subscription/entitlement changes;
- dependent consumer re-evaluation after canonical capability recovery;
- diagnostics/scorecard refresh;
- safe degradation when credentials expire, quota is exhausted or a capability disappears.

### Automatic behavior is prohibited for
- `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION` authority promotion;
- replacing an explicit direct-authority source;
- inferring legal/public/commercial rights from a successful API call, API key or paid plan;
- changing protected deterministic Day/Swing/Long logic;
- allowing AI/LLM/agent output to override Router/rights/authority;
- silently expanding the U.S. Equities Processing product boundary;
- introducing execution capability.

A plan upgrade may automatically make a provider **technically eligible** for new capabilities. It does not automatically grant a higher governed lifecycle/authority state.

## 5. Subscription/entitlement future-proofing

Provider adapters must prefer **effective capability discovery** over hard-coded marketing plan names.

Example desired behavior:

`trial/delayed entitlement -> provider eligible only for delayed/history/shadow workloads`

then later, without a DE.PULSE code change:

`same account/token becomes paid/live -> probe discovers live capability -> Registry updates eligibility -> Smart Provider Router v2 may consider provider for live capability if lifecycle/rights/health policy permits`

Rules:
- do not encode `if plan == Trader` as the routing truth when effective capability can be observed;
- do not require an application release solely because a provider account moves between trial/free/paid tiers;
- plan labels may be displayed as informational evidence when reliably known, but effective capability is the technical routing contract;
- capability loss/downgrade must remove eligibility without corrupting unrelated provider/capability health;
- hysteresis/backoff prevents provider flapping.

## 6. Generic Provider Settings / Credential UX

The TradeInsight Settings/API-key work established the reusable pattern. Every token/API-key provider should use the same generic owner rather than introducing provider-specific settings infrastructure.

Required behavior:
- provider card appears inside the existing Data Providers Settings surface from provider-registration metadata where practical;
- secret entry uses the canonical settings/secrets path;
- renderer/client never receives the stored secret back;
- public state exposes only redacted configuration state such as `has<Key>`/configured status;
- blank save preserves an existing secret;
- non-empty value replaces the secret;
- explicit clear/remove uses the canonical clear-secret owner;
- environment-variable fallback is supported for CI/headless/developer workflows when defined by the adapter;
- stored configured credential takes precedence according to the canonical secret-resolution contract;
- provider `Test` uses the shared provider-test path and returns redacted/auth/capability truth;
- test/diagnostic output must never leak Authorization headers/tokens;
- Settings shows concise configuration/eligibility truth, not raw Router internals;
- detailed runtime provider ranking, active/fallback/recovery, quota, latency and capability evidence belongs in Maintenance/Data Health;
- hosted Commercial/Public zero-key mode ultimately uses server-managed secrets/KMS and must not expose platform provider credentials to ordinary clients.

Provider-card implementation should become metadata-driven enough that adding the next provider does not require redesigning Settings.

## 7. Mandatory cross-integration

A provider is not considered integrated merely because its adapter responds successfully.

Every admitted capability must receive an explicit cross-integration disposition (`REQUIRED / CONDITIONAL / NOT_USEFUL`) against applicable canonical consumers, including:
- Dashboard/ticker/canonical quote consumers;
- Day / Swing / Long Desks;
- Market Tradeability / Readiness;
- Market Intelligence / Market Regime contribution;
- Research;
- Discovery / Opportunity Radar;
- Watchlist / Global Symbols;
- Pre-Market / Market Open Prep;
- Earnings/material catalysts;
- Rapid Move / shock / volume-volatility alerting;
- fundamentals/history/intraday bars;
- options structure where supported;
- SEC/ownership/congressional evidence where supported;
- Data Health / Maintenance;
- canonical persistence/cache/reconciliation;
- future point-in-time Outcome Learning / provider usefulness;
- adaptive processing priority.

Rules:
- no UI surface independently fetches the same provider capability;
- cross-integration consumes canonical routed state;
- specialized provider evidence contributes to Market Regime only through the governed aggregation contract; one provider/ticker observation cannot directly flip global regime;
- unavailable optional evidence cannot falsely degrade unrelated consumers;
- recovered/new evidence must trigger applicable dependent re-evaluation.

## 8. Market Data (`marketdata.app`) — first concrete adoption

Market Data is the first provider intentionally added using this generic contract rather than as another one-off integration.

### Credential contract
- canonical environment fallback: `MARKETDATA_TOKEN`;
- authenticate with `Authorization: Bearer <token>`;
- never place the token in logged URLs;
- add a Market Data card to the existing Data Providers Settings UI using the generic credential UX above;
- persisted secret is redacted from renderer/public state;
- support save-preserve/replace/clear/test semantics identical in principle to TradeInsight;
- never commit a real token to GitHub, fixtures, logs or telemetry.

### Current vendor transport/plan semantics to implement and test
Vendor documentation audited 2026-08-26 states:
- HTTP `200` and `203` are both successful API responses; `203` indicates cache-tier fulfillment and must not be treated as failure;
- official SDKs use `MARKETDATA_TOKEN`;
- free trials of paid plans remain delayed; paid Trader and above can provide real-time stock/options capabilities subject to endpoint-level freshness;
- current Trader Trial advertises 100,000 credits/day, delayed stock/options data, one-year non-AAPL historical depth and restricted premium endpoints;
- if the trial ends without subscription, the account converts to Free Forever under current vendor policy;
- paid entitlement changes should therefore be discovered at runtime rather than requiring a new DE.PULSE release.

These are **provider observations, not permanent assumptions**. Adapter probes and vendor-contract tests must tolerate plan/API evolution and fail truthfully when semantics change.

### Initial Market Data capability posture
Until real configured-token probes and evidence complete:
- enter `SHADOW` / validation-first;
- delayed stock/option evidence may support corroboration, historical/backfill, validation and non-live research where freshness is sufficient;
- delayed data is never represented as live;
- real-time capabilities remain ineligible until the effective entitlement actually proves them;
- after a paid upgrade, live capabilities may become technically eligible automatically, but lifecycle promotion remains governed;
- Market Data must not automatically displace healthy authoritative/current providers merely because it is paid;
- credit/headroom, latency, freshness, health, quality, rights and consumer utility feed normal Router scoring/eligibility;
- current provider single-IP/session restrictions must be handled as operational constraints and surfaced truthfully when encountered, not hidden as generic degradation.

## 9. AI/LLM-style adaptive behavior

“AI/LLM-style” means the provider system learns and adapts from evidence while retaining deterministic safety boundaries.

The Registry/Router/Data Health/provider scorecard should accumulate applicable evidence such as:
- capability availability over time;
- successful/failed requests;
- freshness compliance;
- p50/p95 latency;
- quota/credit pressure;
- cache effectiveness;
- fallback success/recovery time;
- cross-provider disagreement;
- schema/normalization integrity;
- consumer usefulness;
- provider calls avoided;
- cost/economic value where known;
- false degradation/miss contribution;
- outcome usefulness when point-in-time Outcome Learning is available.

This evidence may automatically affect runtime eligibility/selection inside existing Router policy and may produce **promotion recommendations**. It cannot silently grant governed authority or legal rights.

Future adaptive operations may use bounded learned utility/cost priors in SHADOW/champion-challenger form, with deterministic fallback, audit, drift detection and rollback. Smart Provider Router v2 remains the decision owner.

## 10. Adaptive Roadmap propagation

Permanent roadmap direction:
- provider independence is a core platform property;
- new providers enter through one Adapter + Registry contract;
- v19.1 establishes generic registration/capability/entitlement discovery and Market Data adoption;
- v19.3 completes applicable cross-platform/role-aware provider Settings/admin presentation using the same contract, not a second implementation;
- v19.6.1 closes provider reliability/economics/readiness evidence across the provider estate;
- v20.5 may add bounded adaptive provider utility/cost learning while preserving Router/lifecycle/rights authority.

No provider-specific public version is required merely because a new adapter is added if the change fits an existing coherent build.

## 11. Adaptive Build Plan propagation

For every new provider:
1. source-overlap audit and existing-owner reuse;
2. provider manifest + canonical adapter;
3. secret/config metadata and Settings integration;
4. deterministic endpoint/schema/transport fixtures;
5. real configured-key capability/entitlement probe where appropriate;
6. provider-capability matrix + fetch-path classification;
7. Data Health/freshness/authority/rights classification;
8. Registry self-registration;
9. Router-only invocation proof;
10. cross-integration matrix;
11. SHADOW usefulness/quality evidence;
12. adverse/recovery/load/rate-limit tests;
13. Mac/Windows/Web applicability;
14. governed lifecycle promotion decision only after evidence;
15. AIPLC provider-utility learning and next action.

## 12. Adaptive Build Process propagation

G2/G3/G7/G10 must reject a new provider integration if any applicable item is absent:
- canonical owner / no parallel subsystem;
- adapter/manifest registration;
- secret redaction/security;
- entitlement/freshness capability truth;
- Router-only general routing;
- direct-authority exception classification;
- Data Health/freshness/recovery integration;
- rate-limit/quota/backpressure behavior;
- capability-specific fallback;
- no duplicate acquisition/subscription;
- cross-integration disposition;
- lifecycle state and non-auto-promotion proof;
- rights/public-production behavior;
- deterministic fixtures + bounded live smoke when necessary;
- provider plan upgrade/downgrade scenario;
- provider outage/recovery scenario;
- actual consumer reachability;
- required platform evidence.

## 13. Adaptive Delivery Process propagation

A release containing a new provider cannot claim delivery merely because the provider test button succeeds.

Required delivery evidence for applicable capabilities:
- provider adapter and registry appear in the actual runtime;
- credentials remain secret/redacted in actual packaged/web runtime;
- effective capabilities/entitlements are truthful;
- Router selects/skips provider for the correct reasons;
- delayed vs live truth is preserved;
- lifecycle promotion has not occurred without approval;
- adverse auth/403/429/5xx/timeout/stale/malformed/quota/outage/recovery behavior is proven;
- canonical consumers receive routed data without page-specific fetches;
- Settings/Maintenance/Data Health agree about configured/eligible/serving/fallback states;
- Mac Apple Silicon + Windows x64 + Web parity is proven where the shared capability is required;
- Commercial/Public provider rights remain fail closed until separately activated;
- no Stable artifact or previous evidence is reinterpreted.

## 14. Acceptance

This contract is complete only when the runtime eventually proves:

`add adapter once -> self-register -> probe capabilities/entitlement -> enter governed lifecycle -> Smart Provider Router v2 automatically considers eligible capabilities -> canonical state fans out to all useful consumers -> health/quality/economics evidence accumulates -> runtime adapts safely -> subscription changes require no application rewrite`

and all of the following remain true:
- no automatic provider lifecycle/authority promotion;
- no second Router/Data Health/cache/subscription/canonical-state system;
- no provider-specific consumer routing;
- no secret leakage;
- no false live/freshness claims;
- no direct-authority replacement;
- U.S. Equities Processing boundary preserved;
- GLD/SLV/USO actionable exceptions preserved;
- No Execution preserved;
- point-in-time/no-lookahead truth precedes learned provider influence.
