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
- Work slice: `ADAPT-HOSTED-TRUST-FOUNDATION-001`.
- Parent/closure issue: #148.
- Parent hosted program: #66 / `ADAPT-HOSTED-SYNC-001`.
- Draft PR: #149.
- Branch: `adapt-hosted-trust-foundation-001`.
- Baseline `main`: `7c8d0c6614ff4e8c14fc1fabb6aeadcf28a9e92c`.
- Last executable provider-rights implementation checkpoint before this governance transition: `13451ab96722ceea7729720ee89fe2da78abc642`, Fast #1218 / run `33031354009` PASS.
- Canonical closure ledger: `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/closure.json`.
- Canonical work-slice state: `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/work-slice.json`.
- GitHub objects and executable evidence outrank this file and all chat memory. Always fetch live `main`, the active branch, PR #149, issue #148/comments and Actions before writing.

## Permanent zero-miss lifecycle

`PLANNED -> IMPLEMENTED -> PRODUCTION_INTEGRATED -> VERIFIED -> RELEASE_QUALIFIED`

OPEN/PARTIAL/IMPLEMENTED_UNVERIFIED/CONTRACT_ONLY/FRAMEWORK_ONLY are incomplete states only. They never authorize dependency-band advancement, merge, release, production-ready/10/10 claims, or handoff-as-done.

Before VERIFIED, prove the applicable chain:

`requirement -> canonical owner -> production integration -> consumer reachability -> positive behavior -> adverse/fail-closed behavior -> persistence/restart/lifecycle -> security/rights/privacy -> observability -> executable regression -> real external/infrastructure evidence where required -> exact-head CI -> closure ledger`

A newly discovered implementation miss automatically reopens the requirement and dependent completion claims.

## Provider-rights development/public-production rule

User-approved G1 timing disposition on 2026-08-26:

- **Development and pre-public hosted validation use all configured/operationally eligible providers at full available capacity.** Unfinished provider licensing must not suppress Smart Router routes, live subscriptions, cache, persistence or serving during development.
- Strict provider-rights evaluation remains active as audit/governance truth. Missing/unreviewed/expired/downgraded rights remain visible and never become fictional approval.
- Hard fail-closed rights enforcement activates only when the hosted runtime explicitly sets `DEPULSE_PROVIDER_RIGHTS_ENFORCEMENT_MODE=PUBLIC_PRODUCTION` after the user explicitly decides to open DE.PULSE to users/public production.
- Actual provider-specific licence/contract/formal reuse approval is a mandatory deferred `PUBLIC_USER_PROVIDER_RIGHTS_ACTIVATION` gate. This development disposition is not a waiver and never grants legal rights from credentials, successful calls or public terms.
- `isHostedRuntime()` remains independent because identity/session/persistence/infrastructure development still needs hosted behavior before public-user activation.

## HOST-001..003 — verified development/pre-public rights foundation

The provider-rights control plane is production-integrated through existing canonical owners only:

- `provider_data_rights.go` remains the canonical rights model/evaluator/provenance owner.
- All 12 canonical provider identities remain covered by the public-terms review/audit inventory.
- Smart Provider Router v2 remains the sole general routing/admission authority; rights never alter numerical scoring/order.
- `runtime_environment.go` owns the explicit enforcement activation boundary.
- In development/pre-public hosted mode, strict blockers become `AUDIT_ONLY` runtime decisions and configured providers remain available.
- Finnhub/Alpaca subscription desired sets, MarketCache, canonical quote/evidence persistence and runtime/API/SSE serving retain development capacity when public enforcement is inactive.
- In explicit `PUBLIC_PRODUCTION`, the same strict evaluator still fails closed across router admission, live fanout, cache/replay, persistence and serving on missing/tampered/unreviewed/expired/downgraded rights.
- Rights-specific persistence source resolution remains registration-aware, preserves `SEC` versus `SEC EDGAR`, and recognizes BLS/EIA/other canonical providers.
- Fast #1218 / run `33031354009` passed on implementation SHA `13451ab96722ceea7729720ee89fe2da78abc642`, including workflow/source/migration/closure/conservation/portability, T1-T10, security/rights, persistence/lifecycle, gofmt, vet, full `go test ./...` and exact-head recording.
- The governance-transition head containing this handoff must also have fresh exact-head Fast before its dependency-pointer advancement is relied upon.

## Current v19 truth

- `HOST-001..003` provider-rights development/pre-public foundation: **VERIFIED**, subject to exact-head validation of the governance-transition candidate; public-user provider approvals remain a deferred activation gate, not a development blocker.
- `HOST-004..007` tenant/account/device/session/reauth: **OPEN / PARTIALLY IMPLEMENTED — earliest dependency band**.
- `HOST-008..009` product entitlement/quota: **VERIFIED**.
- `HOST-010..012` privacy lifecycle: **OPEN / PARTIALLY IMPLEMENTED**.
- `HOST-013..014` environment/service trust: **OPEN / CONTRACT IMPLEMENTED, REAL INFRASTRUCTURE NOT PROVEN**.
- `HOST-015..016` PostgreSQL tenancy/recovery: **OPEN / PARTIALLY IMPLEMENTED**.
- `HOST-017..023`: **OPEN**.
- Final identical-head Fast + impact-selected Qualified: **OPEN**.

Do not begin v19.1.0 / Hosted Provider Gateway while #148 remains incomplete.

## HOST-004..007 known zero-miss gaps

Existing canonical identity/session implementation already includes tenant IDs/status, role capability matrix, session tenant binding, device registration/binding, LOST/REVOKED denial, session rotation/revocation, recent-password reauthentication and cross-tenant negative authorization.

Known reopened gaps from the strict audit:

- stale devices are rejected by age but are not transitioned to a durable auditable `STALE/RETIRED` lifecycle state;
- device security events do not yet have an explicit privileged audit trail;
- `recordHostedMFAVerification` records externally established MFA proof but does not itself prove a real MFA/passkey ceremony;
- issue #164 remains the end-to-end login/session/silent-renewal/role-entitlement-refresh/logout/revocation/cross-platform evidence owner.

Reuse `identity_model.go`, `http_auth.go`, `identity_reauth.go`, canonical session/workspace persistence and existing security/observability owners. Do not create a second identity/session system.

## Exactly one next action

Continue `HOST-004..HOST-007` on the existing branch/PR: first fetch live GitHub state and exact-head CI, then close auditable stale-device retirement/device-security audit plus actual MFA/passkey-class proof and the applicable #164 end-to-end session lifecycle gaps through the existing canonical identity/session owners; do not start HOST-010+ or v19.1.0 until this dependency band is VERIFIED.

## Later dependency bands

After HOST-004..007 is VERIFIED:

1. Q1 authorization checkpoint covers HOST-001..009 as governed by the work-slice.
2. HOST-010..012: deactivation/privacy-request audit + real backup/PITR/operator recovery deletion proof.
3. HOST-013..014: reproducible IaC/service identity/network/TLS/mTLS enforcement and real drift evidence.
4. HOST-015..016: tenant-owned/scoped PostgreSQL persistence + tagged real DB isolation/migration/adverse/recovery/failover/PITR/restore proof.
5. HOST-017..020: managed secrets/KMS and supply-chain/deploy provenance.
6. HOST-021..022: measured provider scorecards and canonical point-in-time/no-lookahead truth.
7. HOST-023 + final identical-head qualification.

## Permanent architecture/product boundaries

- U.S. Equities Processing only; GLD/SLV/USO remain actionable tradable exceptions.
- No Execution.
- Smart Provider Router v2 is the sole general routing/admission owner.
- Direct SEC/EDGAR remains the governed filing/Form 4 authority.
- Extend canonical Data Health/freshness/cache/persistence/subscription/telemetry/reconciliation/lifecycle/identity/session owners; never create parallel owners.
- No automatic provider lifecycle promotion.
- Point-in-time/no-lookahead truth precedes adaptive learning.
- Provider-rights development mode is audit-only; public-production enforcement requires explicit activation.
- Preserve existing product look-and-feel unless a justified truthful integration/defect correction requires change.

## Resume rule

1. Read this file first.
2. Then read the strict zero-miss contract/decision, `governance/current-state.json`, `governance/ROADMAP.md`, `governance/V19_V20_REBASELINE.md`, the active work-slice `work-slice.json`, `g1-scope.json`, `closure.json`, issue #148/latest comments, issue #164, PR #149 and current Actions.
3. Fetch live `main` and active branch heads before making any change; another session may have advanced them.
4. Inspect commits since the latest verified checkpoint so nothing already implemented is duplicated.
5. Keep Stable v18.10.0 immutable.
6. Do not create another branch/PR, merge/release, start v19.1.0, or weaken G0-G16/source/data-truth/security/platform/CI gates.
7. Update durable GitHub state at meaningful dependency-band transitions so any AI/account can resume independently.
