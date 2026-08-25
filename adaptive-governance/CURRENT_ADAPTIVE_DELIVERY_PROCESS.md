# CURRENT Adaptive Delivery Process

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Retained completed process authority:** #73 / `ADAPT-ROOT-CONVERGENCE-001`  
**Completed provider/data-health program:** #79 with final #84 / PR #91 / merge `733d90ca125a4fe5abd38a2ea40de0623703dfd4`  
**Completed canonical identity:** #92 / PR #93 / merge `57d530e58bfb0b38cc108980cd5cd4a041014db8`  
**Active product work:** #95 / `ADAPT-PROVIDER-ONBOARDING-001` / `adapt-provider-onboarding-001`  
**Separate sibling residual:** #94 / `ADAPT-PROVIDER-TELEMETRY-001`  
**Parent program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`

Delivery remains PR-first, exact-head governed and evidence-bound. #95 uses the existing `adapt-provider-onboarding-001` branch; do not create another branch for the same slice. No Stable/public SemVer release is performed merely to close #95.

Before qualification, the branch must prove through focused regression that:
- route/configuration/quota/cost/delay/capability diagnostic metadata are adopted from the canonical provider registration without introducing a second provider router/registry/health/lifecycle owner;
- incomplete production contracts fail closed even when a lifecycle label says `PRODUCTION`;
- existing provider route ordering/roles and Data Engine capability semantics/order are preserved unless an intentional separately-scoped change is approved;
- free-plan/`NOT_ENTITLED` evidence excludes only the affected capability;
- relevant credential/config change reopens only affected entitlement/configuration suppression and the next canonical Router v2 attempt can establish fresh evidence;
- same-key server-side plan upgrades can be rechecked through the existing privileged Provider Capabilities action without waiting for the stale entitlement cooldown;
- downgrade/402/403/plan-limited evidence safely re-suppresses the capability and uses existing fallback/degradation behavior without poisoning unrelated capabilities;
- current Finnhub, Alpaca, Twelve Data, TradeInsight, Marketaux, FRED, BLS, EIA, SEC/EDGAR, yfinance and CBOE behavior remains compatible;
- direct SEC/EDGAR Form 4 authority, U.S. equities, GLD/SLV/USO actionable exceptions and No Execution remain unchanged;
- no credential value or reusable credential-derived material appears in source-visible evidence, diagnostics, persistence, logs or artifacts.

## Canonical qualification order
1. Re-fetch live `main` and live #95 branch; inspect candidate diff and ensure no concurrent work was lost.
2. Run focused/local/static regression where available and fix genuine findings without weakening gates.
3. Commit final source/governance candidate.
4. Run canonical **CI Fast** on that exact head.
5. If Fast passes, run deliberate impact-selected **CI Qualified** on the identical head. Backend/race/randomized order and any planner-selected renderer/native/DB lanes must pass as selected.
6. Any source change after either gate creates a new candidate; rerun the required exact-head gates.
7. Merge only if live `main` still matches the expected merge base and the PR head matches the qualified candidate.
8. Record immutable run IDs/job conclusions/candidate SHA/merge SHA in #95 and its closure evidence after GitHub produces them. Never fabricate future CI evidence into source.

## Closure boundary
The #95 closure ledger remains fail-closed until every blocking gap is VERIFIED. Source code plus governance text is insufficient. If CI exposes a real defect, correct the architecture/source and requalify; do not waive or weaken assurance. #94 remains a separate residual and #66 remains future-blocked/not started.
