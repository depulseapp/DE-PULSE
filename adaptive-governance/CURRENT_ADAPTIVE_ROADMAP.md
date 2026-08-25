# CURRENT Adaptive Roadmap

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Retained completed process authority:** #73 / `ADAPT-ROOT-CONVERGENCE-001`  
**Completed provider/data-health program:** #79 with final #84 / PR #91 / merge `733d90ca125a4fe5abd38a2ea40de0623703dfd4`  
**Completed canonical identity:** #92 / PR #93 / merge `57d530e58bfb0b38cc108980cd5cd4a041014db8`  
**Active product work:** #95 / `ADAPT-PROVIDER-ONBOARDING-001` / `adapt-provider-onboarding-001`  
**Separate sibling residual:** #94 / `ADAPT-PROVIDER-TELEMETRY-001`  
**Parent program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`

The current v18.9.x direction is to make the completed provider/data-health architecture easier to extend **without creating a second provider architecture**. #95 turns provider onboarding itself into a bounded adaptive contract while preserving the already-qualified Smart Provider Router v2, lifecycle/readiness, Data Health, freshness, persistence, telemetry, reconciliation and degradation owners.

## #95 adaptive destination
A supported provider/capability should require one real adapter/normalizer plus one canonical provider registration. That registration declares only the metadata/evidence needed by existing owners: canonical datasets/capabilities, route priority, configuration requirements, quota/cost/cadence metadata, capability diagnostics, canonical owner/consumer purpose, adapter/schema/timestamp/freshness/failure/rights contracts, evidence/approval references and invalidation rules.

The registration is **not a new authority**. Smart Provider Router v2 remains sole general routing/admission authority; lifecycle/readiness remains the promotion authority; Data Health/freshness/degradation remain their canonical owners; provider data-rights evidence remains independently governed. A `PRODUCTION` label cannot make an incomplete provider contract routable.

## Entitlement adaptation
When credentials or plan-affecting configuration change, only stale entitlement/configuration suppression for the affected provider is reopened before the next canonical Router v2 decision. Real provider evidence remains authoritative; API-key presence never proves a paid plan. If the provider changes entitlement server-side while the same key remains, the existing privileged **Provider Capabilities → Recheck** action may request a bounded fresh entitlement observation instead of waiting for a stale cooldown.

Upgrade/downgrade observations are capability-scoped. A 402/403/plan-limited response may suppress the affected capability and permit existing fallback, but must not poison unrelated capabilities or become a generic provider outage. Successful evidence restores eligibility through the existing capability state owner; production authority still requires governed lifecycle/rights evidence.

## Portability and continuity
#95 is durable only when executable source, the #95 work-slice/closure ledger, `governance/current-state.json`, `handoff/CURRENT.md`, these four CURRENT Adaptive projections, issue #95 and exact-head CI/merge objects agree. Documentation alone cannot close the slice. A new ChatGPT account, another assistant, or another implementation session must be able to fetch the live branch and continue from those repository objects without chat memory.

## Retained boundaries
The completed #79/#84 Data Health foundation remains an invariant. Smart Provider Router v2 remains sole general routing authority; direct SEC/EDGAR remains Form 4 authority; U.S. equities processing, GLD/SLV/USO actionable exceptions and No Execution remain permanent. TradeInsight SEC/search/movers capabilities whose REST contracts are unverified remain fail-closed/gated. #94 stays a separate telemetry/usefulness residual and must not be silently bundled into #95. #66 remains future-blocked/not started.

After #95 closes through exact-head Fast -> impact-selected Qualified -> expected-head merge, #65 resumes fresh semantic-overlap selection of the next genuine residual. v19/v20 inherit this canonical onboarding contract rather than rebuilding a new provider registry/data plane.
