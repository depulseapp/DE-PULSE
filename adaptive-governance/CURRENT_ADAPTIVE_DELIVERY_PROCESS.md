# CURRENT Adaptive Delivery Process

**Canonical rebaseline:** `governance/V19_V20_REBASELINE.md`  
**Provider-registry additive rebaseline:** `governance/V19_V20_PROVIDER_REGISTRY_REBASELINE.md`  
**Provider-registry permanent contract:** `adaptive-governance/ADAPTIVE_PROVIDER_REGISTRY_CONTRACT.md`  
**Machine current-state authority:** `governance/current-state.json`  
**Provider-registry machine map:** `governance/programs/V19-V20-REBASELINE/adaptive-provider-registry.json`  
**Certified Stable:** `v18.10.0` — immutable  
**Active version:** `v19.0.0`  
**Active work slice:** `ADAPT-HOSTED-TRUST-FOUNDATION-001`  
**Active issue/PR/branch:** #148 / PR #149 / `adapt-hosted-trust-foundation-001`  
**Canonical closure ledger:** `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/closure.json`

v18.10.0 delivery remains authoritative: exact-head Fast, Qualified, canonical Release G11-G16, actual macOS Apple Silicon + Windows x64 native evidence, provenance/SBOM and no-rebuild publication.

## Version/build delivery rule

Public planning and delivery are version-oriented. Requirement IDs, backlog acceptance bullets, commits and CI evidence remain granular inside the version but do not become separate release packets.

For each version:
- use one coherent development branch/PR unless governance explicitly requires otherwise;
- a **canonical Fast exact-head PASS** is required for the coherent candidate before it can be treated as Fast-qualified evidence;
- a **Qualified exact-head PASS** is required at G10 and at other material risk boundaries selected by Impact Planner/current governance;
- Release G11-G16 is used only for a release candidate, not every implementation checkpoint;
- batch related small changes when they share owners/evidence;
- split genuinely heavy/high-risk features into a real patch version rather than creating hidden pseudo-releases;
- cancel/supersede obsolete candidate runs when supported, but never reuse evidence from the wrong head;
- keep the three canonical workflow families: Fast, Qualified, Release.

## Conserved Data Health delivery boundary

Provider/Data Health delivery evidence remains governed by the completed #80/#81/#82/#83/#78/#84 program and its machine-readable matrix/SLO/fetch-path contracts. Delivery must prove that routable capabilities remain under Smart Provider Router v2, while explicit source-authority boundaries remain intact.

Direct **SEC/EDGAR** remains authoritative where the canonical provider contract classifies it as direct authority, including Form 4. A secondary provider, router preference or AI/adaptive layer cannot silently replace that authority.

Delivery evidence for affected provider/data surfaces must show truthful freshness, minimally scoped `PARTIAL COVERAGE` / `DATA DEGRADED`, eligible fallback/cache recovery, hysteresis/no-flapping, optional-provider isolation, workload/backpressure protection and capability-level lifecycle/readiness. Missing required evidence must remain visible rather than being normalized into false health.

**No Execution** remains a permanent product boundary: no order routing, broker execution, paper execution, P&L, portfolio or execution-adjacent feature can be introduced or implied by provider/data-health, AI/agent or hosted-platform delivery work.

## Adaptive Provider Registry delivery boundary

A release containing a new provider cannot be considered delivery-complete because the provider authenticates or a Settings Test button succeeds.

For every new Registry provider/capability, actual release evidence must prove applicable:
- the adapter/manifest is present in the actual runtime and self-registers through the canonical Adaptive Provider Registry;
- Smart Provider Router v2 remains the sole general routing/admission/selection owner;
- the Registry does not duplicate Router, Data Health, cache, persistence, subscription, telemetry, lifecycle or canonical state;
- effective capabilities/entitlements/freshness/history/quota are represented truthfully where observable;
- Settings secret state is redacted and the actual packaged/web runtime cannot recover/display the stored token;
- blank save preserve, non-empty replace, explicit clear and provider-test semantics work through canonical owners;
- Router selects/skips/demotes/falls back/recovers the provider for correct capability-specific reasons;
- delayed data is never represented as live;
- provider lifecycle promotion has not occurred without explicit governed approval;
- direct-authority providers remain protected;
- rights/public-production activation remains separate from credentials/technical eligibility;
- all required cross-integration consumers receive canonical routed state rather than page-specific provider fetches;
- Settings, Maintenance and Data Health agree on configured, eligible, serving, fallback/recovery, freshness and entitlement truth;
- optional provider failure does not contaminate unrelated capability/application health;
- required Mac Apple Silicon, Windows x64 and Web parity is proven for shared provider Settings/admin/runtime behavior where applicable.

## Market Data delivery requirements

Market Data (`marketdata.app`) is the first concrete adopter of the Registry contract and therefore must prove the reusable onboarding path rather than only MarketData-specific happy-path behavior.

Applicable evidence includes:
- `MARKETDATA_TOKEN` environment fallback and canonical stored-secret behavior;
- Bearer header authentication with no token in logged URLs;
- HTTP 200 and 203 success handling under current vendor semantics;
- missing/invalid token, 401/403, 429/quota exhaustion, 5xx, timeout, malformed/partial/schema-drift responses;
- current delayed trial/free freshness represented truthfully;
- current trial limits treated as observed vendor semantics rather than permanent hard-coded assumptions;
- trial -> Free downgrade;
- trial/free -> paid/live effective-capability expansion;
- no application source release required solely because the external subscription changes when effective capability can be reprobed;
- paid/live technical eligibility does not imply primary-provider status or automatic lifecycle promotion;
- quota/credit headroom, latency, freshness, health, cost/utility and rights feed the normal Router policy;
- provider single-IP/account operational restrictions are surfaced/classified truthfully when encountered;
- cross-integration to all required Research, Discovery/Radar, Desk, Prep, Market Intelligence/Regime contribution, alert, history/options, Data Health/Maintenance and future Outcome Learning consumers through canonical state;
- no secret leakage in logs, telemetry, fixtures, artifacts, browser state or GitHub evidence.

## Delivery acceptance beyond code correctness

A version cannot advance merely because its feature works on a happy path. Applicable closure requires:
- every mapped backlog/HOST/source-discovered requirement owned;
- positive + failure/adverse evidence;
- truthful Data Health/freshness/recovery behavior;
- cross-feature re-evaluation after canonical state changes;
- role/RBAC/product-entitlement/provider-right negative evidence;
- persistence/restart/migration/recovery evidence;
- load/resource/backpressure evidence;
- required Mac/Windows/Web evidence;
- #170 cross-integration/Market Regime disposition for intelligence;
- #171 UI/data-density/intelligence-maturity disposition for visible surfaces;
- point-in-time/no-lookahead lineage for historical/outcome/adaptive evidence;
- deterministic fallback for AI/adaptive features;
- current roadmap/build/process/delivery/handoff convergence.

For provider work, closure additionally requires adapter/Registry ownership, capability/entitlement truth, secure Settings/secret behavior, Router eligibility/selection evidence, Data Health integration, cross-integration, lifecycle disposition, direct-authority/rights protection and subscription-change behavior.

## CI efficiency without quality loss

CI savings come from coherent candidate batching, dependency-aware Impact Planner routing, avoiding requirement-sized PRs/releases, reusing frozen v18 evidence only when unchanged, and reserving expensive native/browser/release lanes for the risk surfaces that require them. Do not combine unrelated product scope simply to save Actions minutes and do not remove assurance to meet a budget.

For provider tests, use deterministic fixtures/historical replay/canonical cached evidence for normal CI when live vendor behavior is not the subject. Use bounded explicit live smoke only when authentication, entitlement, real provider transport or effective capability must be proven; never leak secrets or consume live credits unnecessarily.

## Cross-platform delivery

Shared product capability is not complete until every REQUIRED client is equivalent in meaning. Mac/Windows/Web may differ in platform mechanics, not domain truth, auth/role semantics, provider rights, freshness/provenance, intelligence conclusions or durable state.

Provider Settings should consume the same generic metadata/secret contract across required platforms. Do not build a MarketData-only cross-platform UX stack after the native/runtime adapter exists.

## Adaptive/AI delivery

Adaptive influence must progress through governed evidence states and remain bounded/reversible. SHADOW results cannot silently become production truth. AI/agent/MCP changes must preserve canonical market-data owners, rights filtering, provenance, audit and a deterministic non-AI product path.

Provider intelligence follows this same boundary. Availability, latency, freshness, quota, cache value, fallback/recovery quality, disagreement, cost and consumer usefulness may influence Router selection within governed policy and may create promotion recommendations. Adaptive systems cannot silently promote provider lifecycle/authority, infer rights, or become a parallel router.

## Exactly one next action

Continue `v19.0.0` on PR #149 from current live GitHub/executable evidence. The Adaptive Provider Registry / Market Data rebaseline is approved future v19.1 scope and must not bypass unfinished v19.0 Development Production Ready work. Governance-only rebaseline commits require fresh exact-head Fast before dependency-band advancement or release qualification.
