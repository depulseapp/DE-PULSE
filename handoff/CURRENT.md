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
- Latest HOST-013/014 Azure checkpoint: `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/host013-014-azure-checkpoint-2026-08-28.json`.
- GitHub objects and executable evidence outrank this file and chat memory. Always fetch live branch/PR/issue/Actions before writing.

## Permanent readiness and zero-miss rules

Development Production Ready means technically robust, secure, persistent, cross-platform, tested and full provider/data-capable. Commercial/Public Ready additionally requires provider licensing/rights, public-user legal/compliance and the commercial activation audit.

Lifecycle:

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

**VERIFIED FOR v19.0 DEVELOPMENT.**

- Application checkpoint `5127b1599c052bb3c709b46bd7900cc46629d0ee`; Fast #1257 / run `33051229439` PASS.
- Managed recovery exact head `f357a1852640bf88aae69936d6593b63d9fd155d`.
- Fast #1303 / run `33138630810` PASS.
- Ordinary Qualified #231 / run `33139659529` PASS on the same head; manual recovery job skipped by design.
- Manual Qualified #232 / run `33154312937` PASS.
- Retained artifact `9678997463`, digest `sha256:15cddb9d278a9c469659c864cdb216d8d04cf348cbf2c3c00855bd29071e85f0`.
- Real Neon PITR restored a point before deletion, replayed canonical privacy enforcement, passed restart verification, resurrected zero deleted profile/workspace/device/session state, preserved the tombstone and measured RPO 7.753s / RTO 13.926746s within declared targets.

This evidence closes HOST-010..012 and is reusable for HOST-016. It does **not** close HOST-015..016 tenant-owned/scoped PostgreSQL isolation, migration, HA/failover and broader recovery obligations.

### HOST-013..014 — environment/service trust

**IMPLEMENTED_UNVERIFIED — CURRENT ACTIVE BAND.**

Canonical application policy remains `internal/hostedenv/desired_state_v1.json` and `internal/hostedenv/contract.go`; no second runtime authority was created.

Repository implementation includes:
- portable Kubernetes/Istio trust projection under `governance/hosted-infrastructure/`;
- canonical renderer `tools/hosted/render_kubernetes_trust.py`;
- fail-closed hosted infrastructure and external-egress conservation gates in the existing G2/source-health path;
- unique environment namespaces/service identities, default-deny networking, explicit managed ingress, canonical external-host allowlist, Istio STRICT mTLS, AuthorizationPolicy and `REGISTRY_ONLY` outbound behavior;
- Azure AKS Terraform substrate in Canada Central with private cluster, OIDC issuer, Microsoft Entra Workload Identity, local accounts disabled, Kubernetes RBAC, Azure Policy/network policy, Key Vault CSI hook and scoped BYO-VNet Network Contributor authority;
- dedicated workload managed identity plus federated Kubernetes service-account credential and exported non-secret client/principal/subject binding evidence;
- Microsoft-managed Istio with governed `asm-X-Y` revision selection;
- AzureRM remote Terraform state using OIDC + Microsoft Entra data-plane authentication;
- `tools/ci/host013_azure_operator.py` and secret-free live evidence collector;
- private-cluster inspection via `az aks command invoke`;
- mandatory post-verification `terraform plan -detailed-exitcode` zero-drift proof;
- manual-only Azure operator job inside the existing `.github/workflows/ci-qualified.yml` with `id-token: write`, pinned `azure/login` and `hashicorp/setup-terraform`, exact-head/dev-only enforcement and 30-day secret-free evidence retention;
- no fourth workflow and no long-lived Azure client-secret path.

Repository implementation checkpoint before governance reconciliation:
- `f2c003a444744bfd193ef7cef421b6b044734af8`.
- Exact-head Fast #1336 / run `33179213802`: **PASS**.
- The six commits after the prior Azure checkpoint touched only the planned integration surface: `ci-qualified.yml`, Azure workload identity/federated binding, dependency lock, reproducibility gate, workflow policy and Azure outputs.

The subsequent checkpoint/handoff reconciliation commits are governance-only and require their own exact-head Fast before being cited as final current-head evidence.

**Remaining HOST-013/014 verification requirement:** run the existing manual `CI Qualified` Azure operator against a real isolated non-production Azure subscription and prove:
- actual workload identity token acquisition/service-account binding;
- ingress TLS/certificate lifecycle;
- strict internal mTLS positive/adverse traffic;
- managed-ingress authorization and unauthorized ingress denial;
- canonical allowed egress success and unregistered-egress denial;
- cross-environment isolation/denial;
- secret-free live evidence;
- zero post-apply Terraform drift.

Repository IaC/CI cannot substitute for real managed infrastructure evidence.

### HOST-015..016 — PostgreSQL tenancy/recovery

**OPEN / PARTIALLY IMPLEMENTED.** Reuse HOST-012 Neon PITR evidence, but tenant-owned/scoped PostgreSQL persistence, live migration/isolation/adverse testing and broader HA/failover/rollback/recovery proof remain open.

### HOST-017..020

**OPEN** for managed secrets/KMS and supply-chain/deploy provenance.

### HOST-021..022

**OPEN** for measured provider/Data Health scorecards and canonical point-in-time/revision/no-lookahead truth.

### HOST-023

**OPEN** for aggregate Development Production Ready qualification after applicable HOST-001..022 technical closure.

Do not begin v19.1/#153 while #148 remains technically incomplete.

## Exactly one next action

Remain in **HOST-013..014 — Azure-backed Environment / Service Trust**. Obtain/configure the non-secret Azure OIDC/state identifiers for an isolated non-production subscription/environment, run the existing manual `CI Qualified` Azure operator on the exact candidate head, and retain live adverse/drift evidence. Do not mark HOST-013..014 VERIFIED or begin HOST-015+ until this managed-environment evidence passes, unless governance explicitly records an allowed Development-tier external blocker.

## Later dependency bands

1. HOST-015..016: tenant-owned/scoped PostgreSQL + live migration/isolation/adverse/HA/failover/recovery evidence.
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

## Resume rule

1. Read this file first.
2. Read the readiness-tier and strict zero-miss authorities.
3. Read `governance/current-state.json`, the canonical roadmap/build plan/build process/delivery process, `governance/V19_V20_REBASELINE.md`, active work-slice metadata and `closure.json`.
4. Fetch live `main`, active branch, PR #149, issue #148/latest comments, issue #164 and current Actions before writing.
5. Inspect commits since the latest verified checkpoint; do not duplicate work.
6. Keep Stable v18.10.0 immutable.
7. Do not create another branch/PR, merge/release, start v19.1.0, or weaken G0-G16/source/data-truth/security/platform/CI gates.
8. At each meaningful dependency-band transition, update durable GitHub state so any AI/account can resume independently.
