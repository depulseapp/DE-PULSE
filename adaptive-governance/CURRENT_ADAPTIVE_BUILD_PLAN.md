# CURRENT Adaptive Build Plan

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Retained completed process authority:** #73 / `ADAPT-ROOT-CONVERGENCE-001`  
**Completed provider/data-health program:** #79 with final #84 / PR #91 / merge `733d90ca125a4fe5abd38a2ea40de0623703dfd4`  
**Completed canonical identity:** #92 / PR #93 / merge `57d530e58bfb0b38cc108980cd5cd4a041014db8`  
**Active product work:** #95 / `ADAPT-PROVIDER-ONBOARDING-001` / `adapt-provider-onboarding-001`  
**Separate sibling residual:** #94 / `ADAPT-PROVIDER-TELEMETRY-001`  
**Parent program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`

#95 required build outputs are:
1. `provider_registration.go` is the single provider-onboarding descriptor consumed by existing owners for route membership/priority, configuration, quota/cost/expected-delay metadata and capability diagnostics.
2. Every production-routable dataset declaration carries canonical owner/consumer purpose plus adapter/schema/timestamp/freshness/failure/rights/evidence/approval/invalidation contracts. Incomplete declarations fail closed even if lifecycle text says `PRODUCTION`.
3. Existing `routeChains()`, configured-provider checks and provider quota/cost/delay helpers project from the canonical registration rather than requiring new provider-specific switches for each integration.
4. `provider_capabilities.go` projects existing truthful provider capability diagnostics from the same registration without losing special semantics such as TradeInsight SHADOW, Alpaca plan-limited activity, Finnhub premium evidence or Twelve Data direct-FX detection.
5. `provider_entitlement_refresh.go` owns only configuration-triggered entitlement invalidation. One-way fingerprints remain process-local and never enter logs, diagnostics, governance or persistence.
6. Smart Provider Router v2 invokes targeted configuration revalidation before ranking, not during snapshot construction, so snapshot locks are not mutated/re-entered.
7. Existing **Provider Capabilities → Recheck** may force bounded fresh entitlement observation for same-key server-side plan changes; it may not erase genuine outage/rate-limit/NOT_SUPPORTED/healthy evidence.
8. Regression proof covers synthetic-provider adoption, route parity, fail-closed incomplete contracts, free-plan exclusion, config-change upgrade, same-key manual recheck, successful re-eligibility, downgrade/402/403 fallback and unrelated-capability isolation.
9. Regression preserves current behavior for Finnhub, Alpaca, Twelve Data, TradeInsight, Marketaux, FRED, BLS, EIA, SEC/EDGAR, yfinance and CBOE, plus direct SEC authority, U.S. equities, GLD/SLV/USO and No Execution.
10. Preserve existing Data Engine capability-row behavior/order; registration refactoring must not introduce unrelated UI rearrangement.
11. Keep #94 separate. #95 does not introduce semantic usefulness scoring or alter production provider order from new #94 evidence.
12. Maintain `governance/current-state.json`, #95 work-slice/closure ledger, all four CURRENT Adaptive projections, `handoff/CURRENT.md` and issue #95 as a cross-session executable resume set.

## Retained Adaptive Data Health build owners
#95 inherits rather than replaces the completed #80 Data Health artifacts: `governance/data-health/provider-capability-matrix.json`, `governance/data-health/data-health-slo.json`, and `governance/data-health/provider-fetch-paths.json`. Those remain executable classification/SLO/fetch-path inputs to Smart Provider Router v2 and source-health recurrence. The live **Adaptive Roadmap**, **Build Plan**, **Build Process**, and **Delivery Process** must retain this foundation while identifying #95 as the active extension.

The completed #79 program remains executable input: canonical provider capability, Data Health SLO, fetch-path classification, Router v2, lifecycle/readiness, freshness, persistence, telemetry and degradation owners cannot be replaced by #95. Provider onboarding describes/adopts capabilities into those owners; it does not become a new health/router/lifecycle system.

## Executable owners already introduced on #95 branch
- `provider_registration.go`
- `provider_entitlement_refresh.go`
- `provider_router.go`
- `smart_router_v2.go`
- `provider_capabilities.go`
- `data_engine_handlers.go`
- `provider_registration_regression_test.go`
- `provider_onboarding_adaptive_regression_test.go`
- `provider_capability_projection_regression_test.go`

These source changes are implementation evidence but are **not yet qualification evidence**. The #95 closure ledger remains fail-closed until focused regressions and canonical exact-head CI succeed.

Delivery remains one existing branch/one PR: reconcile source compatibility first, run focused tests, run canonical Fast on the exact final candidate, then impact-selected Qualified on the identical head, then expected-head merge only if live `main` is unchanged. No Stable/public SemVer release is created solely for #95.
