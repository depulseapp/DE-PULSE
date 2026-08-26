# CURRENT Adaptive Build Plan

**Canonical rebaseline:** `governance/V19_V20_REBASELINE.md`  
**Backlog/version matrix:** `governance/programs/V19-V20-REBASELINE/backlog-version-matrix.json`  
**HOST/version map:** `governance/programs/V19-V20-REBASELINE/host-requirement-version-map.json`  
**Legacy future-commitment conservation:** `governance/programs/V19-V20-REBASELINE/legacy-future-commitment-conservation.json`  
**Cross-integration matrix:** `governance/programs/V19-V20-REBASELINE/cross-integration-matrix.json`  
**Whole-product surface map:** `governance/programs/V19-V20-REBASELINE/whole-product-surface-rebaseline.json`  
**Certified Stable:** `v18.10.0` — immutable  
**Active build:** `v19.0.0` — Hosted Trust & Identity Foundation  
**Active issue/PR/branch:** #148 / PR #149 / `adapt-hosted-trust-foundation-001`

## Build sizing rule

Plan and communicate work by **real version/build**, not requirement packets. Requirement rows, issue acceptance bullets and CI evidence rows remain granular traceability units inside the version.

- Combine small, related changes when they share canonical owners and acceptance evidence.
- Give feature-heavy/risk-heavy work its own version or patch version. Current deliberate heavy splits include `v19.4.1` Discovery, `v19.5.1` Options/GEX and `v20.3.1` AODR.
- Do not create a version for every requirement, card, page defect or CI checkpoint.
- Do not enlarge a version so much that ownership, adverse testing or rollback becomes unclear.
- Use commits and risk-based CI checkpoints for implementation progress; they are not product planning units.

## Required build matrix for every version

Before G3, every assigned requirement/issue/legacy-conservation row must resolve:
- source-overlap disposition: `INHERITED / EXTEND_EXISTING_OWNER / REPLACE_CONSOLIDATE / NEW_RESIDUAL / EXTERNAL_BLOCKED / N_A`;
- canonical owner and upstream evidence;
- actual consumers and user/trader decision purpose;
- data freshness/provenance/point-in-time semantics;
- positive + negative/failure evidence;
- persistence/restart/migration applicability;
- role/auth/product-entitlement/provider-right applicability;
- load/resource/backpressure applicability;
- Mac/Windows/Web requirement or justified N/A;
- durable regression owner.

For intelligence-bearing work also record:
- intelligence maturity: `DETERMINISTIC_ONLY / ADAPTIVE_CANDIDATE / LEARNING_ENABLED / AI_ASSISTED / NOT_USEFUL`;
- downstream cross-integration `REQUIRED / CONDITIONAL / NOT_USEFUL` for Market Regime, Tradeability, Discovery, Watchlist, Research, Desks, Prep, alerts, Outcome/Pattern, Data Health/Maintenance and processing priority;
- Market Regime contribution `YES / CONDITIONAL / NO`;
- Outcome Learning contribution `YES / CONDITIONAL / NO`;
- duplicate/isolation conflict and consolidation decision;
- re-evaluation behavior after canonical recovery/new evidence.

For visible work also record #171 disposition: `KEEP / IMPROVE / CONSOLIDATE / COLLAPSE / MOVE / REMOVE`.

## v19 build plan

- `v19.0.0` trust/identity: HOST-001..023, #164 core auth/session, security portions of #156.
- `v19.1.0` canonical runtime/global symbols: #150/#151/#153/#154/#155/#160/#167 core + conserved Router/Data Health/Shared Symbol responsibilities.
- `v19.2.0` gateway/shared serving/sync: HOST-024..039.
- `v19.3.0` cross-platform roles/product IA: HOST-040..047/053 + #152/#156/#159/#160 presentation/#167 Admin/#171/#164 UX + `LEGACY-TRADER-SETUP-SHORT-001`.
- `v19.4.0` Market Intelligence/Research: HOST-049 + #158/#161/#162/#171.
- `v19.4.1` Discovery: HOST-048 + #163/#171 + conserved halt/LULD/pause/resume behavior.
- `v19.5.0` price-volume/event intelligence: #168/#169.
- `v19.5.1` options/GEX structure: #157.
- `v19.6.0` point-in-time/outcome-ready foundation: HOST-057..064 + deterministic #165 + institutional/13F/two-sided thesis lineage.
- `v19.6.1` reliability/economics/readiness: HOST-050..056/065..071 + ADR-GDI + provider-gap + final #170/#171 reconciliation.
- `v19.7.0` Major Closure: HOST-072; no feature scope.

### v19.3.0 two-sided deterministic Desk setup acceptance

Current source is not a true SHORT implementation: bearish labels are possible, but plan geometry/action-state logic remains long-oriented. `v19.3.0` must correct the **contract**, not merely change text or colors.

Required implementation behavior:
- explicit side: `LONG / SHORT / NO_SETUP-WAIT`;
- separate directional evidence from 0–100 setup-quality score;
- a strong SHORT setup can score highly for quality without being encoded as a low numeric score;
- LONG: target/trim above entry and invalidation below;
- SHORT: cover/target below entry and invalidation above;
- R-multiple/action-state/entry-distance/sort/chart overlays/replay snapshots work correctly for both sides;
- Research and Discovery consume the same canonical setup side/geometry rather than reconstructing it;
- existing Day/Swing/Long Desk look-and-feel remains materially unchanged; only truthful side-aware labels/content change where necessary;
- No Execution remains permanent.

The institutional/TDTI two-sided thesis substrate remains separate in v19.6.0/v20.3.0.

## v20 build plan

- `v20.0.0` Outcome Learning & Adaptive Control Plane, including conserved experiment/calibration/guardrail responsibilities.
- `v20.1.0` Adaptive Chart Pattern & Similarity Intelligence (#166).
- `v20.2.0` Adaptive Market Synthesis, Market Regime & Discovery Learning, including conserved ASBI normalization/synthesis/contradiction/abstention/outcome responsibilities.
- `v20.3.0` Adaptive Institutional / Two-Sided Thesis.
- `v20.3.1` AODR Adaptive Opportunity Intelligence.
- `v20.4.0` Agent Orchestration & Controlled MCP/API.
- `v20.5.0` Adaptive Operations.
- `v20.6.0` Professional Adaptive Closure; no feature scope.

## Zero-miss rule

No version closes with an unassigned applicable certified-v18 responsibility, backlog row, HOST requirement, legacy future-roadmap/build-plan commitment, source-discovered responsibility, cross-integration, role/right case, UI disposition or regression owner. #170 and #171 are applied inside the affected versions rather than deferred as documentation-only work.

## Exactly one next action

Continue `v19.0.0` on PR #149. Fetch the live head/checks, then resolve the known hosted identity/session production-reachability/source-health gap by production-wiring or correctly consolidating those helpers through the existing canonical auth/HTTP owners. Run exact-head Fast on the coherent candidate. Do not begin the next version until current-version closure criteria are met.
