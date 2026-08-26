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

Continue `v19.0.0` on PR #149. Exact-head Fast #1165 on `f36417fda84e063d5a9cafcc31c464b051f5b3af` proved the original hosted-helper orphan/source-health defect is resolved, but Fast remained red because the CURRENT Adaptive Data Health delivery contract had drifted. Restore the conserved Roadmap / Build Plan / Build Process / Delivery Process projection, then obtain fresh exact-head Fast and continue current-version closure only from that evidence. Do not advance to `v19.1.0` yet.
