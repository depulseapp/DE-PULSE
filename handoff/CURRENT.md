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
- GitHub/source/executable evidence outranks prose and chat memory. Always fetch live `main`, branch, PR #149, issue #148/comments and Actions before writing because another session may advance them.

## 2026-08-26 zero-miss implementation audit + HOST-001..003 remediation checkpoint

A fresh requirement-by-requirement audit invalidated the earlier projection that HOST-001..016 were all VERIFIED. Fast #1184 is green on `1f1f766f0dd89350f29d513735d4ef1f34db9086`, but Fast/unit/policy evidence cannot substitute for missing production integration, real infrastructure or real external/legal evidence.

HOST-001..003 production-integration remediation has now closed the application-side consumption gaps found by that audit without creating a second router or weakening the permanent root policy:

- `provider_data_rights.go` is the canonical provider-rights owner. Rights are provider-bound, fail closed, and reviewed evidence is loaded only through a SHA-256-pinned provider-rights bundle. A credential/successful call never grants rights.
- Smart Provider Router v2 remains the sole routing/admission authority; hosted rights affect eligibility/admission only and do not mutate numerical route scoring.
- Finnhub/Alpaca live subscription desired sets re-evaluate hosted rights, so downgrade/expiry removes future live fanout instead of relying on reconnect behavior.
- Canonical market-cache serialization/replay and canonical persistence save/load re-check current provider rights; mixed legacy data without enforceable provider provenance fails closed in hosted mode.
- Runtime/bootstrap/SSE/direct symbol/news/earnings/filing serving now applies hosted-rights enforcement before the existing per-user privacy projection. Provider-unbound legacy/derived payloads remain withheld pending HOST-022 provenance rather than receiving inferred permission.
- Positive/adverse lifecycle regression coverage is consolidated into existing root owners; temporary duplicate root files were removed rather than relaxing `governance/root-layout-policy.json`.
- Fast #1191 on pre-cleanup head `acf0b1f4045271ffa84a4e9c444a0b8d67cd895a` failed at G16/root ownership before Go because those temporary duplicate root files still existed. That failure is superseded by the structural cleanup; it is not a product PASS or Go failure.

**HOST-001..003 remains OPEN / NOT VERIFIED.** Real provider-specific approved legal/data-rights evidence/provenance is still not present/bound for actual providers, and fresh exact-head executable CI after the cleanup is still required. Test rights bundles are fixtures only and are not legal evidence.

Current truthful band state:

- `HOST-001..003` provider rights: **OPEN / PRODUCTION-INTEGRATED, NOT VERIFIED**. Application-side router/live/cache/persistence/serving consumption is wired fail closed; real provider-specific approved rights evidence and fresh exact-head CI remain blocking.
- `HOST-004..007` tenant/account/device/session/reauth: **OPEN / PARTIALLY IMPLEMENTED**. Tenant/capability/session/device/revocation/rotation/reauth paths exist with cross-tenant denial. Still missing auditable stale-device retirement and an actual MFA/passkey-class ceremony/integration; #164 end-to-end auth/session lifecycle audit remains open.
- `HOST-008..009` product entitlement/quota: **VERIFIED**. Canonical identity persistence owns plan/status/metering; product entitlement is separate from RBAC/provider rights; authenticated HTTP enforces CSRF -> entitlement -> anti-abuse -> product meter -> protected projection; transitions/quota/reload/adverse tests exist.
- `HOST-010..012` privacy lifecycle: **OPEN / PARTIALLY IMPLEMENTED**. Account export/delete/tombstone/anti-resurrection restore are real and production-wired. Required deactivation/privacy-request audit and actual backup/PITR/operator recovery enforcement are not complete.
- `HOST-013..014` environment/service trust: **OPEN / CONTRACT IMPLEMENTED, INFRASTRUCTURE NOT PROVEN**. Desired-state/runtime drift validation exists. No repository IaC/deployment enforcement currently proves environment isolation, workload/service identity, network policy or mTLS in provisioned infrastructure.
- `HOST-015..016` PostgreSQL tenancy/recovery: **OPEN / PARTIALLY IMPLEMENTED**. Environment search_path/TLS/pool/recovery-policy validation exists, but environment schemas are not tenant schemas. PostgreSQL still has global `identity_state(id=1)` and `user_workspaces(user_id PK)` without tenant-owned keys/query scoping. HA/backup/PITR/restore/RPO-RTO/rollback are self-declared runtime readiness values, not live infrastructure proof. Tagged Postgres integration tests are not run by Fast.
- `HOST-017..023`: **OPEN**. Do not advance into these while earlier dependency bands remain open.
- final identical-head Fast + impact-selected Qualified: **OPEN**.

## Required remediation order

1. Finish `HOST-001..003`: obtain/bind reviewable real provider-specific approved rights evidence/provenance, run fresh exact-head executable Fast on the cleaned canonical-owner implementation, fix any gate/test failure, and only then determine whether the band can advance to VERIFIED.
2. Reclose `HOST-004..007`: complete stale-device retirement/audit, real MFA/passkey-class proof integration and applicable #164 session lifecycle evidence.
3. Reclose `HOST-010..012`: add deactivation/privacy-request audit and prove cross-store deletion/retention against real hosted recovery behavior.
4. Reclose `HOST-013..014`: implement/prove reproducible environment/IaC, service identity, network/TLS/mTLS enforcement rather than self-declaration alone.
5. Reclose `HOST-015..016`: implement tenant-owned/scoped PostgreSQL persistence and real Postgres integration/recovery/failover/backup/PITR/restore evidence.
6. Only then continue `HOST-017..018` managed secrets/KMS, followed by HOST-019..023 in frozen dependency order.
7. Final exact-head Fast + impact-selected Qualified on the identical final candidate, then governed expected-head merge/release only if every closure row is truthful.

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

## Exactly one next action

1. Fetch live `main`, branch `adapt-hosted-trust-foundation-001`, PR #149, issue #148/comments and Actions state.
2. Reconcile any commits after this checkpoint before writing.
3. Stay at the **earliest reopened dependency band: HOST-001..003 provider rights**.
4. Require real provider evidence/provenance plus fresh exact-head executable Fast over the production-integrated canonical boundaries before marking that band VERIFIED.
5. Do not merge, release, start v19.1.0 or activate Hosted Provider Gateway while any HOST-001..023 row/final qualification remains open.

## Resume rule

1. Read this file first, then `governance/V19_V20_REBASELINE.md`, `governance/ROADMAP.md`, `governance/current-state.json`, all v19/v20 machine companions, current adaptive Roadmap/Build Plan/Build Process/Delivery Process/Gap Closure/CI Convergence, issue #148, parent #66, PR #149 and current CI evidence.
2. Inspect commits since the latest verified/audit checkpoint so nothing already implemented is duplicated.
3. GitHub/source/executable evidence outranks prose and chat memory.
4. Preserve v18.10.0 Stable immutability and permanent architecture boundaries.
5. Do not create another branch/PR or requirement-sized public version for active v19.0.0.
6. Do not weaken G0-G16, source/data-truth/security/platform/CI gates.
7. Update durable GitHub state at every meaningful dependency-band transition so another assistant/account/model can resume independently.
