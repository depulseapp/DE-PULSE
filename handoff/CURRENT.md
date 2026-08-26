# DE.PULSE Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

## Stable authority

- **Certified Stable:** `v18.10.0` — immutable.
- Stable candidate SHA: `584e9e0ce91ec08e08cfd52c7cf60392ab74dd12`.
- Stable qualified source SHA: `ec39319c86dee5e5976751abc42bc96a402a6d46`.
- Stable source fingerprint: `0adbd70aeb9a016b0e4ded93538cfb75d616494980c11d7d781cffa31b1e6037`.
- Stable build ID: `v18.10.0-stable-20260825`; platform build `181000`; release run `32917159547`.
- Do not rebuild, republish, overwrite or reinterpret v18.10.0 from v19 work.

## Active source of truth

- **Active version:** `v19.0.0` — Hosted Trust & Identity Foundation.
- Work slice: `ADAPT-HOSTED-TRUST-FOUNDATION-001`.
- Parent issue: #148; parent hosted program: #66 / `ADAPT-HOSTED-SYNC-001`.
- Draft PR: #149.
- Branch: `adapt-hosted-trust-foundation-001`.
- Canonical closure ledger: `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/closure.json`.
- Canonical planning: `governance/ROADMAP.md`, `governance/V19_V20_REBASELINE.md`, and all five `governance/programs/V19-V20-REBASELINE/*` machine companions.
- GitHub objects and executable evidence outrank this handoff and all chat memory. Always fetch live `main`, live branch head, PR #149, issue #148/comments and Actions state before writing because another session may advance them.

## Current v19.0.0 verified dependency bands

The active v19 work retains one coherent branch/PR; do not create requirement-sized branches/releases.

- `HOST-001..003` provider legal/data rights: **VERIFIED** through canonical rights/router owners.
- `HOST-004..007` tenant/account/device/session/reauth: **VERIFIED** through canonical identity/auth/session owners with cross-tenant and revoked-state denial.
- `HOST-008..009` product entitlement/quota: **VERIFIED**; product capability remains separate from RBAC/provider rights and metering is durable/fail-closed.
- `HOST-010..012` privacy lifecycle: **VERIFIED**; secret-free portability/export, deletion/tombstones and anti-resurrection restore semantics are executable.
- `HOST-013..014` environment/service trust: **VERIFIED**; dev/test/stage/prod desired state, unique isolation/service identities, default-deny network posture, HTTPS/TLS and internal mTLS fail closed on drift.
- `HOST-015..016` PostgreSQL tenancy/recovery: **VERIFIED** on implementation head `33c74e399cc33baea3acf8761b65ab5be7d42dc3` by Fast #1183 / run `33011966479`.

PostgreSQL details:
- existing PostgreSQL persistence remains the sole hosted persistence stack;
- `persistence_backend_select.go` production-wires the hosted database policy before backend construction;
- `internal/hostedpersistence/postgres_policy_v1.json` binds unique `depulse_dev/test/stage/prod` `search_path` schemas, `sslmode=verify-full`, bounded pools, explicit application-authorization RLS disposition, expand-contract migrations, HA/encrypted backup/PITR, RPO/RTO, restore-drill freshness and tested rollback;
- policy tests fail closed on schema/TLS/capacity/recovery/rollback drift and do not leak database credentials;
- successful connection/ping alone is not hosted readiness.

Exact-head evidence:
- Fast #1182 / run `33009612801` PASS on `fcfaae2f5f234eb5ef790e1b6679c809b0d30cb5`, validating the privacy/environment/service-trust candidate and all inherited earlier v19 bands.
- Fast #1183 / run `33011966479` PASS on `33c74e399cc33baea3acf8761b65ab5be7d42dc3`, including workflow policy, source health, migration/closure/conservation/resume, v18 T1-T10 assurance, gofmt, go vet, full Go unit/package suite and exact-head status.

## Remaining v19.0.0 dependency order

1. `HOST-017..018` — managed secrets/KMS lifecycle.
2. `HOST-019..020` — supply-chain/deploy provenance.
3. `HOST-021` — provider scorecards/usefulness without changing Smart Provider Router v2 authority.
4. `HOST-022` — point-in-time/revision-aware/no-lookahead truth.
5. `HOST-023` — zero-gap closure.
6. Final exact-head Fast + impact-selected Qualified on the identical final candidate, then governed expected-head merge/release only if every closure row is truthful.

Do not begin `v19.1.0` or Hosted Provider Gateway work before #148 closes.

## Permanent architecture/product boundaries

- U.S. Equities Processing only; GLD/SLV/USO remain actionable exceptions.
- **No Execution**.
- Smart Provider Router v2 remains the sole general routing/admission owner.
- Direct SEC/EDGAR remains governed filing/Form 4 authority.
- Extend canonical Data Health/freshness/cache/persistence/subscription/telemetry/reconciliation/lifecycle/identity/session owners; never create parallel owners.
- No automatic provider lifecycle promotion.
- Point-in-time/no-lookahead truth must precede adaptive learning.
- Preserve Day/Swing/Long Desk look-and-feel, Dashboard Market Regime/Desk Control, Data Engine look-and-feel, and existing AI Copilot visual treatment except justified truthful integration/defect corrections.
- Future mirrored SHORT setup correction remains conserved for v19.3.0; do not pull it into v19.0.0 opportunistically.

## Adaptive/intelligent north star

`canonical evidence -> deterministic intelligence -> cross-feature synthesis -> point-in-time outcome accumulation -> bounded adaptive learning -> optional AI/agent explanation/orchestration`

Every applicable feature/data point must have purpose/consumer/materiality/freshness/retention, reuse canonical data, and explicitly integrate Market Regime/Outcome Learning only where materially justified. One symbol-level signal never directly flips the global regime; missing/stale evidence is not neutral evidence.

## Exactly one next action

1. Fetch live `main`, branch `adapt-hosted-trust-foundation-001`, PR #149, issue #148/comments and current Actions state.
2. Confirm the latest head includes the HOST-015/016 closure projection and has executable exact-head Fast; if another session advanced it, reconcile that newer evidence first.
3. Continue **HOST-017..018 managed secrets/KMS** through the existing local-secret/provider configuration owners. Do not create a second secret subsystem, expose raw provider secrets to clients/sync/logs/traces/backups, or weaken source/root/security gates.
4. Require server-side managed references, environment/least-privilege isolation, health-before-cutover, rotation/cutover/rollback/emergency revoke and privileged audit evidence; obtain fresh exact-head Fast before advancing to HOST-019..020.
5. Do not merge, release, or start v19.1.0 while any HOST-001..023 row or final Fast/Qualified closure remains open.

## Resume rule

1. Read this file first, then `governance/V19_V20_REBASELINE.md`, `governance/ROADMAP.md`, `governance/current-state.json`, all five v19/v20 rebaseline machine companions, current adaptive Roadmap/Build Plan/Build Process/Delivery Process/Gap Closure/CI Convergence, issue #148, parent #66, PR #149 and current CI evidence.
2. Inspect commits since the latest verified checkpoint so nothing already implemented is duplicated.
3. GitHub/source/executable evidence outranks prose and chat memory.
4. Preserve v18.10.0 Stable immutability and permanent architecture boundaries.
5. Do not create another branch/PR or requirement-sized public version for the active v19.0.0 work.
6. Do not weaken G0-G16, source/data-truth/security/platform/CI gates.
7. Update durable GitHub state at every meaningful dependency-band transition so another assistant/account/model can resume independently.
