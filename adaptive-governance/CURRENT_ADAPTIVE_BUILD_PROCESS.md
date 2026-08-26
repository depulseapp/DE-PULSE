# CURRENT Adaptive Build Process

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.10.0` — immutable  
**Retained process-control authority (historical COMPLETE):** #107 / `ADAPT-PROVIDER-PROFESSIONAL-CLOSURE-001` / `adapt-provider-professional-closure-001`  
**Active v19 G1:** #148 / `ADAPT-HOSTED-TRUST-FOUNDATION-001` / `adapt-hosted-trust-foundation-001`  
**Closure ledger:** `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/closure.json`.

The permanent execution loop remains source-driven and exact-head: LOOKUP -> COMPARE -> CLASSIFY -> DECIDE -> UPDATE -> Fast -> Qualified -> G11–G16 when a release is actually produced. #148 uses one development branch and one long-lived Draft PR for the coherent `HOST-001..HOST-023` v19.0.x trust-foundation band.

The inherited Data Health process remains in force: **#81/#82/#83/#78/#84** continue to govern Router adoption, common health/recovery, lifecycle, TradeInsight admission and zero-gap closure. Unclassified provider/fetch paths fail closed, and **canonical freshness** remains authoritative.

## Packet / CI operating rule

Use the smallest dependency-correct engineering packet that can be reasoned about and tested completely. Small packets are encouraged when they prevent implementation misses, clarify ownership, isolate failures, or make adverse testing tractable. They are **not** automatically separate releases, branches, PRs, Qualified runs or public versions.

For each packet: implement behavior and evidence together; run focused local/unit/static checks while editing where available; batch a coherent candidate; then use exact-head Fast. Mark only `IMPLEMENTED_UNVERIFIED` after Fast when final band qualification is still outstanding. Qualified is risk-boundary based: after `HOST-001..009`, after `HOST-010..020`, and final G10/`HOST-001..023`, unless a genuinely new risk surface justifies another checkpoint.

## Active process rules

1. Implement dependency-first: provider rights -> tenant identity/session -> product entitlement -> privacy/environment/trust -> PostgreSQL -> managed secrets/supply chain -> scorecards/point-in-time -> band closure.
2. Bind evidence while coding. Development is not complete until canonical owner, consumer, positive/adverse evidence and persistence/security/UI applicability are recorded.
3. Classify assurance findings as product behavior, test/evidence, ownership binding, or N/A before adding machinery.
4. Keep Smart Provider Router v2, Data Health/freshness, cache/persistence, multi-feed subscription, telemetry/reconciliation/lifecycle, identity/session, Research/Discovery, decision/outcome lineage and direct SEC/EDGAR as canonical owners.
5. Hosted threat modeling includes tenant escape, cache/coalescing leakage, long-lived stream revocation, secret leakage, mixed-client downgrade, noisy-neighbor pressure and backup/restore isolation where applicable.
6. Historical/adaptive evaluations reconstruct what was genuinely knowable at the decision timestamp; later revisions cannot leak backward.
7. PostgreSQL connection success is not hosted readiness. Tenant isolation, migration/recovery, privacy and adverse evidence are required.
8. One coherent candidate per coherent correction. Changed candidates get fresh Fast; do not manufacture requirement-sized branches or CI events.
9. Frozen v18 T1–T10 remain durable baseline evidence; deeper historical logic is impact-triggered while changed product regressions remain fully enforced.
10. Main-push continuity/branch hygiene remains operationally distinct from PR Fast.
11. Public version numbers identify coherent product/band states; internal packet IDs identify engineering progress. Do not make version numbering drive architecture decomposition.

## Exactly one next action

Finish production reachability for `HOST-004..007` through existing identity/auth HTTP owners, validate the coherent candidate with Fast, and continue the same v19.0.x band without opening another requirement-sized PR/release.
