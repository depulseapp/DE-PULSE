# DE.PULSE Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Certified Stable candidate:** `e55d8d25b15cec2ffb0f5411bc358bc40b359cf9`  
**Certified Stable source fingerprint:** `0062f46dea5690d0b3fcd8a9ed3b1f71ebe1522c7dee2cb218e9d36b9e0076ff`  
**Certified Stable build ID:** `v18.9.1-stable-20260821`  
**Immutable predecessor resume checkpoint release:** `v18.9.0` / `v18.9.0-stable`  
**Current live-main baseline for this slice:** `57d530e58bfb0b38cc108980cd5cd4a041014db8`  
**Completed provider/data-health program:** #79 with final #84 / PR #91 / merge `733d90ca125a4fe5abd38a2ea40de0623703dfd4`  
**Completed canonical identity slice:** #92 / PR #93 / merge `57d530e58bfb0b38cc108980cd5cd4a041014db8`  
**Active product slice:** #95 / `ADAPT-PROVIDER-ONBOARDING-001` / `adapt-provider-onboarding-001`  
**Work-slice:** `governance/work-slices/ADAPT-PROVIDER-ONBOARDING-001/work-slice.json`  
**Closure ledger:** `governance/work-slices/ADAPT-PROVIDER-ONBOARDING-001/closure.json`  
**Parent program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`  
**Separate sibling residual:** #94 / `ADAPT-PROVIDER-TELEMETRY-001` — do not fold into #95  
**Future work:** #66 remains blocked/not started.

## Current authority
GitHub objects and executable evidence outrank this handoff. #92 is complete and must not be reopened. #95 is the active implementation branch and is not documentation-only: executable source already exists on `adapt-provider-onboarding-001`, while its closure ledger intentionally remains OPEN until exact-head qualification and guarded merge evidence exist.

Closed issues #96–#100 are accidental connector artifacts explicitly titled `[ACCIDENTAL TOOL INVOCATION — IGNORE]`; they carry no DE.PULSE scope or authority and must never be treated as work items.

## #95 product contract
DE.PULSE provider onboarding must be adaptive without creating a second router/registry/health/lifecycle system. A supported provider capability joins the existing architecture through one canonical registration plus its real adapter/normalizer and governed evidence. Smart Provider Router v2 remains the sole general routing/admission authority.

A production-routable capability must fail closed unless its registration carries the required canonical owner, consumer purpose, adapter/schema/timestamp/freshness/failure/rights contracts, evidence/approval references, invalidation rule, lifecycle and priority. A `PRODUCTION` label alone is insufficient.

Provider credential/configuration changes must reopen only affected stale entitlement/configuration suppression before the next canonical Router v2 decision. Provider plan changes with the same API key are handled by the existing **Provider Capabilities → Recheck** operator action, which can force a bounded fresh entitlement observation. Real provider responses remain authoritative; key presence never proves a paid plan.

## Executable implementation already on branch
- `provider_registration.go` — canonical provider/capability onboarding contract; route/config/quota/cost/delay and capability-diagnostic metadata.
- `provider_entitlement_refresh.go` — process-local one-way configuration fingerprints and targeted entitlement revalidation; no secret/hash export or persistence.
- `provider_router.go` / `smart_router_v2.go` — Router v2 consumes registration-driven metadata and refreshes configuration entitlement state before ranking.
- `provider_capabilities.go` — capability diagnostics project from provider registration instead of a second manual provider list.
- `data_engine_handlers.go` — manual capability recheck reopens stale entitlement suppression before bounded provider probes.
- `provider_registration_regression_test.go` — adaptive registration/fail-closed/synthetic-provider proof.
- `provider_capability_projection_regression_test.go` — registration-driven diagnostics and same-key recheck proof.

Implementation commits through the first full code pass include `c9323d16...`, `092b153e...`, `7924fdde...`, `0403f403...`, `447fe3cf...`, `662deb81...`, `021c20da...`, `be550080...`, `8054cbbb...`, and `f083b7cf...`. Fetch the live branch before relying on any recorded SHA because governance/test fixes may have advanced it.

## Remaining #95 acceptance before merge
1. Preserve existing Data Engine provider-capability row behavior/order; no incidental UI rearrangement from registration refactor.
2. Run focused/exact-head tests and fix any compile/regression findings without weakening architecture or gates.
3. Prove configured free-plan exclusion, config-change upgrade, same-key manual recheck, successful evidence re-eligibility, downgrade/402/403 fallback, and unrelated-capability isolation.
4. Prove all existing providers retain current routing/diagnostic behavior: Finnhub, Alpaca, Twelve Data, TradeInsight, Marketaux, FRED, BLS, EIA, SEC/EDGAR, yfinance and CBOE.
5. Bind #95 into all four CURRENT Adaptive governance projections and keep #94 separate.
6. Obtain canonical Fast PASS on the exact final candidate, then impact-selected Qualified PASS on that identical head.
7. Merge only with expected-head guard against the then-current `main`. No Stable/public SemVer release solely for #95.

## Permanent boundaries
- U.S. equities processing only.
- No Execution/order routing.
- Smart Provider Router v2 is the sole general routing/admission authority.
- Direct SEC/EDGAR remains authoritative for Form 4.
- GLD, SLV and USO remain governed actionable exceptions.
- Reuse canonical freshness/cache/persistence/telemetry/state/validation/lifecycle owners.
- No secrets in diagnostics, tests, governance or handoff artifacts.
- No automatic provider lifecycle promotion; `SHADOW -> VALIDATED -> APPROVED -> PRODUCTION` remains evidence/governance controlled.
- G0–G16 and canonical Fast/Qualified/Release workflows only; no gate weakening or temporary workflow family.

## Exactly one next action
Continue #95 implementation/qualification on the existing `adapt-provider-onboarding-001` branch: first reconcile diagnostic-order compatibility, then run focused/exact-head tests and correct real findings, update all four CURRENT Adaptive projections, and proceed to canonical Fast followed by impact-selected Qualified. Do not create another branch, do not merge early, and do not start #66.

## Resume rule
1. Fetch live `main` and live `adapt-provider-onboarding-001` first because another process/session may have advanced either.
2. Read this file first, then `governance/current-state.json`, `AGENTS.md`, the AI-assistant portability and CI-efficiency contracts, issue #65 latest comments, issue #95, #95 comments, and the #95 work-slice/closure ledger. Read #94 only to preserve sibling-scope separation.
3. Inspect commits from live `main` baseline through the live #95 branch head before changing source so existing work is never duplicated.
4. Continue actual implementation from the exact branch state; do not restart planning and do not reopen #92/#79/#84.
5. Preserve Router v2, direct SEC/EDGAR, canonical Data Health/freshness/persistence/telemetry/lifecycle owners, U.S.-equities boundaries, GLD/SLV/USO exceptions and No Execution.
6. Treat registration metadata as admission description, not independent authority: lifecycle/readiness, data-rights evidence, canonical freshness/degradation and Router v2 remain their existing owners.
7. Never close #95 from documentation alone. Exact-head executable evidence and guarded merge evidence are mandatory.
