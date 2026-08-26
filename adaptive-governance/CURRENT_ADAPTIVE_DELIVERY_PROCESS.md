# CURRENT Adaptive Delivery Process

**Canonical rebaseline:** `governance/V19_V20_REBASELINE.md`  
**Certified Stable:** `v18.10.0` — immutable  
**Active version:** `v19.0.0`  
**Active issue/PR/branch:** #148 / PR #149 / `adapt-hosted-trust-foundation-001`

v18.10.0 delivery remains authoritative: exact-head Fast, Qualified, canonical Release G11-G16, actual macOS Apple Silicon + Windows x64 native evidence, provenance/SBOM and no-rebuild publication.

## Version/build delivery rule

Public planning and delivery are version-oriented. Requirement IDs, backlog acceptance bullets, commits and CI evidence remain granular inside the version but do not become separate release packets.

For each version:
- use one coherent development branch/PR unless governance explicitly requires otherwise;
- exact-head Fast validates coherent changed candidates;
- Qualified is impact/risk selected and mandatory at G10;
- Release G11-G16 is used only for a release candidate, not every implementation checkpoint;
- batch related small changes when they share owners/evidence;
- split genuinely heavy/high-risk features into a real patch version rather than creating hidden pseudo-releases;
- cancel/supersede obsolete candidate runs when supported, but never reuse evidence from the wrong head;
- keep the three canonical workflow families: Fast, Qualified, Release.

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

## CI efficiency without quality loss

CI savings come from coherent candidate batching, dependency-aware Impact Planner routing, avoiding requirement-sized PRs/releases, reusing frozen v18 evidence only when unchanged, and reserving expensive native/browser/release lanes for the risk surfaces that require them. Do not combine unrelated product scope simply to save Actions minutes and do not remove assurance to meet a budget.

## Cross-platform delivery

Shared product capability is not complete until every REQUIRED client is equivalent in meaning. Mac/Windows/Web may differ in platform mechanics, not domain truth, auth/role semantics, provider rights, freshness/provenance, intelligence conclusions or durable state.

## Adaptive/AI delivery

Adaptive influence must progress through governed evidence states and remain bounded/reversible. SHADOW results cannot silently become production truth. AI/agent/MCP changes must preserve canonical market-data owners, rights filtering, provenance, audit and a deterministic non-AI product path.

## Exactly one next action

Continue `v19.0.0` on PR #149. Fix the source-health/product-reachability gap for hosted identity/session helpers, then obtain fresh exact-head Fast. Later versions remain blocked until current-version exit criteria are satisfied.
