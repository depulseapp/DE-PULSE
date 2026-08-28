# DE.PULSE Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

## Stable authority

- **Certified Stable:** `v18.10.0` — immutable.
- Stable candidate SHA: `584e9e0ce91ec08e08cfd52c7cf60392ab74dd12`.
- Stable source fingerprint: `0adbd70aeb9a016b0e4ded93538cfb75d616494980c11d7d781cffa31b1e6037`.
- Stable build ID: `v18.10.0-stable-20260825`.
- Do not rebuild, republish, overwrite or reinterpret v18.10.0 from v19 work.

## Active source of truth

- Active version: `v19.0.0` — Hosted Trust & Identity Foundation.
- Current readiness target: **DEVELOPMENT_PRODUCTION_READY**.
- Work slice: `ADAPT-HOSTED-TRUST-FOUNDATION-001`.
- Parent/closure issue: #148.
- Draft PR: #149.
- Branch: `adapt-hosted-trust-foundation-001`.
- Baseline `main`: `7c8d0c6614ff4e8c14fc1fabb6aeadcf28a9e92c`.
- Auth/session issue #164 remains open only for roadmap-assigned v19.3 client/UX parity unless a new core security defect appears.
- Parent hosted program: #66 / `ADAPT-HOSTED-SYNC-001`.
- Permanent readiness authority: `governance/PRODUCTION-READINESS-TIERS.md` + `governance/production-readiness-tiers.json`.
- Canonical closure ledger: `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/closure.json`.
- GitHub objects and executable evidence outrank this file and chat memory. Always fetch live branch/PR/issue/Actions before writing.

## Permanent readiness and zero-miss rules

Development Production Ready means technically robust, secure, persistent, cross-platform, tested and full provider/data-capable. Commercial/Public Ready additionally requires provider licensing/rights, public-user legal/compliance and the commercial activation audit.

Lifecycle is:

`PLANNED -> IMPLEMENTED -> PRODUCTION_INTEGRATED -> VERIFIED -> RELEASE_QUALIFIED`

Before VERIFIED, prove the applicable chain:

`requirement -> canonical owner -> production integration -> consumer reachability -> positive/adverse behavior -> persistence/restart/lifecycle -> security/rights/privacy -> observability -> executable regression -> real external/infrastructure evidence where required -> exact-head CI -> closure ledger`

OPEN/PARTIAL/IMPLEMENTED_UNVERIFIED/CONTRACT_ONLY/FRAMEWORK_ONLY remain incomplete and never authorize merge/release/readiness claims.

## Current v19.0 dependency-band truth

### HOST-001..003 — provider-rights control plane

**VERIFIED FOR DEVELOPMENT.** Technical provenance/audit/public fail-closed machinery is implemented. Actual provider-specific public/commercial approvals remain separate Commercial/Public activation gates and do not grant legal rights today.

### HOST-004..007 — tenant/account/device/session/MFA

**VERIFIED FOR v19.0 DEVELOPMENT.** Tenant isolation, capability-scoped RBAC, device/session lifecycle and the production-wired Ed25519 MFA-class challenge/signature ceremony are executable. #164 remains open only for later v19.3 client/UX parity.

### HOST-008..009 — product entitlement/quota

**VERIFIED.** Product plan/status/capability/quota truth is independently composed with RBAC/provider rights and fails closed before protected projection.

### HOST-010..012 — privacy lifecycle and managed recovery

**VERIFIED FOR v19.0 DEVELOPMENT.** This supersedes the stale earlier handoff wording that kept the band open.

Application checkpoint:
- `5127b1599c052bb3c709b46bd7900cc46629d0ee`
- Fast #1257 / run `33051229439` PASS.
- Secret-free export, reversible self-deactivation, session revocation, governed reactivation, destructive deletion, sanitized tombstones, durable privacy audit events, restart persistence and application archive anti-resurrection are executable.

Managed recovery closure checkpoint:
- recovery repair head `f357a1852640bf88aae69936d6593b63d9fd155d`
- Fast #1303 / run `33138630810` PASS.
- ordinary Qualified #231 / run `33139659529` PASS on the same head; manual HOST-012 operator job skipped there by design.
- manual Qualified #232 / run `33154312937` PASS on exact head `f357a1852640bf88aae69936d6593b63d9fd155d`.
- retained artifact `9678997463`, digest `sha256:15cddb9d278a9c469659c864cdb216d8d04cf348cbf2c3c00855bd29071e85f0`.
- Real Neon PITR restored a point before deletion, canonical `enforceArchiveAccountDeletionPrivacy` replayed deletion truth, restart verification passed, deleted users/workspace/device/session state did not resurrect, tombstone remained, recovery was isolated, and measured RPO/RTO were within declared targets.

This evidence closes HOST-010..012. It is also relevant evidence for HOST-016, but **does not close HOST-015..016 as a whole** because tenant-owned/scoped PostgreSQL isolation, live migration/adverse testing and broader HA/failover/recovery obligations remain open.

### HOST-013..014 — environment/service trust

**IMPLEMENTED_UNVERIFIED — CURRENT ACTIVE BAND.**

Canonical application policy remains `internal/hostedenv/desired_state_v1.json` and `internal/hostedenv/contract.go`; no second runtime authority was created.

Repository implementation now includes:
- portable Kubernetes/Istio trust projection under `governance/hosted-infrastructure/`;
- canonical renderer `tools/hosted/render_kubernetes_trust.py`;
- fail-closed contract gate `tools/ci/hosted_infrastructure_contract_gate.py` bound into the existing G2/source-health path;
- namespace/environment isolation, dedicated service accounts, default-deny networking, explicit managed ingress, governed egress host registration, Istio STRICT mTLS, authorization policy and `REGISTRY_ONLY` outbound mesh behavior;
- desired-state SHA-256 binding and broad-egress rejection;
- Azure selected as the real hosted substrate while keeping the workload layer cloud-portable;
- Terraform AKS adapter under `governance/hosted-infrastructure/azure/` with private AKS, OIDC issuer, Microsoft Entra Workload Identity, local accounts disabled, Kubernetes RBAC, Azure Policy, Azure network policy, Key Vault CSI hook, user-assigned control-plane identity and Network Contributor prerequisite for the BYO VNet.

Repository evidence checkpoints:
- `bc3875ed0ff52d3fc18a2cc3938ed1834c8a4690` — portable HOST-013/014 gate integration; Fast #1305 / run `33155263302` PASS.
- `f306b47093d1adb2f11b1f0c62e0dbb8d3a5b268` — machine-readable hosted-trust checkpoint; Fast #1306 / run `33155375900` PASS.
- Azure substrate commits advanced through `ea5c83d8aaed0d8ecb9edfea127178bd13eb7cff`; Fast #1310 / run `33155849814` was triggered for that exact head and must be checked live before citing it as passed.

**Remaining verification requirement:** deploy a real non-production Azure AKS environment and execute adverse/drift evidence proving the declared isolation, workload identity, ingress/egress denial, TLS/mTLS behavior and reproducibility against managed infrastructure. Repository IaC/CI cannot substitute for this live evidence.

Azure/GitHub authentication policy for this verification is OIDC workload identity federation only. Do not store a long-lived Azure client secret. Expected GitHub configuration values are `AZURE_CLIENT_ID`, `AZURE_TENANT_ID`, and `AZURE_SUBSCRIPTION_ID`, preferably environment-scoped with protection rules.

### HOST-015..016 — PostgreSQL tenancy/recovery

**OPEN / PARTIALLY IMPLEMENTED.** The successful HOST-012 Neon PITR drill is real managed recovery evidence and must be reused, but this band still requires tenant-owned/scoped PostgreSQL persistence, live migration/isolation/adverse testing and broader HA/failover/rollback/recovery proof.

### HOST-017..020

**OPEN** for managed secrets/KMS and supply-chain/deploy provenance.

### HOST-021..022

**OPEN** for measured provider/Data Health scorecards and canonical point-in-time/revision/no-lookahead truth.

### HOST-023

**OPEN** for aggregate Development Production Ready qualification after applicable HOST-001..022 technical closure.

Do not begin v19.1/#153 while #148 remains technically incomplete.

## Exactly one next action

Remain in **HOST-013..014 — Azure-backed Environment / Service Trust**. Finish exact-head repository qualification for the current Azure baseline, then provision an isolated non-production AKS environment using OIDC-federated GitHub/Azure identity and run real infrastructure adverse/drift tests. Do not mark HOST-013..014 VERIFIED until live managed-environment evidence exists.

## Later dependency bands

1. HOST-015..016: tenant-owned/scoped PostgreSQL + live migration/isolation/adverse/HA/failover/recovery evidence, reusing the successful Neon PITR evidence rather than duplicating it.
2. HOST-017..020: managed secrets/KMS and supply-chain/deploy provenance.
3. HOST-021..022: measured provider scorecards/Data Health and canonical point-in-time/no-lookahead truth.
4. HOST-023 + final identical-head Development Production Ready qualification.
5. Commercial/Public activation remains a separate later gate.
6. Only after v19.0 Development closure may v19.1 begin, including #153 Adaptive Provider Registry / Market Data work.

## Permanent architecture/product boundaries

- U.S. Equities Processing only; GLD/SLV/USO remain actionable tradable exceptions.
- No Execution.
- Smart Provider Router v2 is the sole general routing/admission owner.
- Adaptive Provider Registry may register/project capabilities but never becomes a second Router.
- Direct SEC/EDGAR remains the governed filing/Form 4 authority.
- Extend canonical Data Health/freshness/cache/persistence/subscription/telemetry/reconciliation/lifecycle/identity/session owners; never create parallel owners.
- No automatic provider lifecycle/authority promotion.
- Point-in-time/no-lookahead truth precedes adaptive learning.
- Development provider-rights mode remains audit-only; public-production enforcement requires explicit Commercial/Public activation.
- Keep v18.10.0 Stable immutable.

## Approved future provider-registry / Market Data scope

Future v19.1/#153 remains governed by:
1. `adaptive-governance/ADAPTIVE_PROVIDER_REGISTRY_CONTRACT.md`
2. `governance/V19_V20_PROVIDER_REGISTRY_REBASELINE.md`
3. `governance/programs/V19-V20-REBASELINE/adaptive-provider-registry.json`
4. `handoff/PROVIDER-REGISTRY-REBASELINE.md`
5. issue #153.

Target remains:

`Provider Adapter -> Adaptive Provider Registry -> capability/entitlement probes -> rights/authority + Data Health -> Smart Provider Router v2 -> canonical state -> all useful consumers`

Do not implement this early.

## Resume rule

1. Read this file first.
2. Read the readiness-tier and strict zero-miss authorities.
3. Read `governance/current-state.json`, `governance/ROADMAP.md`, `governance/V19_V20_REBASELINE.md`, active work-slice metadata and `closure.json`.
4. Fetch live `main`, active branch, PR #149, issue #148/latest comments, issue #164 and current Actions before writing.
5. Inspect commits since the latest verified checkpoint; do not duplicate work.
6. Keep Stable v18.10.0 immutable.
7. Do not create another branch/PR, merge/release, start v19.1.0, or weaken G0-G16/source/data-truth/security/platform/CI gates.
8. At each meaningful dependency-band transition, update durable GitHub state so any AI/account can resume independently.
