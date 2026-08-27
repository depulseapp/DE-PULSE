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
- Auth/session evidence issue: #164; v19.0 core is verified, issue remains open only for its roadmap-assigned v19.3 UX/client-parity residual unless a new core defect is discovered.
- Parent hosted program: #66 / `ADAPT-HOSTED-SYNC-001`.
- Draft PR: #149.
- Branch: `adapt-hosted-trust-foundation-001`.
- Baseline `main`: `7c8d0c6614ff4e8c14fc1fabb6aeadcf28a9e92c`.
- Permanent readiness authority: `governance/PRODUCTION-READINESS-TIERS.md` + `governance/production-readiness-tiers.json`.
- Canonical closure ledger: `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/closure.json`.
- Canonical work-slice state: `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/work-slice.json`.
- GitHub objects and executable evidence outrank this file and all chat memory. Always fetch live `main`, active branch, PR #149, issue #148/comments and Actions before writing.

## Permanent production-readiness tiers

These definitions are authoritative and supersede conflicting older wording:

> **Development Production Ready = technically robust, secure, persistent, cross-platform, tested, full provider/data capability.**
>
> **Commercial/Public Ready = Development Production Ready + provider licensing/rights + public-user legal/compliance + commercial activation audit.**

Consequences:

- v19.0.0 is currently being closed against **Development Production Ready**, not Commercial/Public Ready.
- Development keeps full configured/operational provider and data capability needed for technical validation.
- Technical provider-rights metadata/provenance/evaluator/audit/public fail-closed machinery remains a Development obligation.
- Actual provider-specific licence/contract/formal public/commercial reuse approval is a **Commercial/Public activation-only gate**, not a Development technical blocker.
- Public-user legal/compliance approval and the commercial activation audit are Commercial/Public-only gates.
- Technical privacy/security/deletion/retention/tenant isolation/MFA/session/device/infrastructure/database/secrets/data-truth obligations remain Development blockers where applicable.
- `Commercial/Public Ready => Development Production Ready`; the reverse is never implied.
- `publicProductionAuthorized` remains false until Development Production Ready plus all Commercial/Public activation gates pass and explicit activation is recorded.
- Missing Commercial/Public-only paperwork must be tracked separately; it must not make an otherwise satisfied Development technical row OPEN/BLOCKED_EXTERNAL/NOT_APPLICABLE for the wrong reason.

## Permanent zero-miss lifecycle

`PLANNED -> IMPLEMENTED -> PRODUCTION_INTEGRATED -> VERIFIED -> RELEASE_QUALIFIED`

OPEN/PARTIAL/IMPLEMENTED_UNVERIFIED/CONTRACT_ONLY/FRAMEWORK_ONLY are incomplete states only. They never authorize dependency-band advancement, merge, release, Development Production Ready, Commercial/Public Ready, 10/10, or handoff-as-done claims.

Before VERIFIED, prove the applicable chain:

`requirement -> canonical owner -> production integration -> consumer reachability -> positive behavior -> adverse/fail-closed behavior -> persistence/restart/lifecycle -> security/rights/privacy -> observability -> executable regression -> real external/infrastructure evidence where required -> exact-head CI -> closure ledger`

Zero-miss rigor applies independently inside the readiness tier being claimed. A newly discovered implementation miss automatically reopens the requirement and dependent completion claims.

## Provider-rights development/public-production rule

- Development and pre-public hosted validation use all configured/operationally eligible providers at full available capacity. Unfinished provider licensing must not suppress Smart Router routes, live subscriptions, cache, persistence or serving during development.
- Strict provider-rights evaluation remains active as audit/governance truth. Missing/unreviewed/expired/downgraded rights remain visible and never become fictional approval.
- Hard fail-closed rights enforcement activates only when the hosted runtime explicitly sets `DEPULSE_PROVIDER_RIGHTS_ENFORCEMENT_MODE=PUBLIC_PRODUCTION` after Commercial/Public activation gates are satisfied and explicit public activation is authorized.
- Actual provider-specific licence/contract/formal reuse approval remains mandatory before public/commercial use. Development disposition is not a waiver and never grants legal rights from credentials, successful calls or public terms.
- `isHostedRuntime()` remains independent because identity/session/persistence/infrastructure development still needs hosted behavior before public-user activation.

## Current v19 re-audit truth

The v19.0 core auth/session implementation checkpoint `de110ecfcbc6e385203eec90c631628871fd6e83` passed exact-head Fast #1254 / run `33049452646`. The canonical HOST-004..007 closure-ledger transition at `c964d303a4d2379972026844b3444fbef0ed9382` then passed exact-head Fast #1255 / run `33049748420`. These are dependency-band checkpoints, not final v19.0 qualification.

- `HOST-001..003` provider-rights development/pre-public control plane: **VERIFIED FOR DEVELOPMENT**. Actual provider-specific public/commercial approvals remain a separate Commercial/Public activation gate and grant no rights today.
- `HOST-004..007` tenant/account/device/session/reauth band: **VERIFIED FOR v19.0 DEVELOPMENT**.
  - **HOST-004 VERIFIED:** tenant/account isolation is canonical; privileged admin visibility/mutations are actor-tenant scoped; user creation binds actor tenant; cross-tenant role/status/password/session mutations are denied; critical-owner invariant is per tenant.
  - **HOST-005 VERIFIED:** canonical SUPER_OWNER/OWNER/ADMIN/USER/DEMO capability truth is preserved. ADMIN is capability-scoped, `roleHasHostedCapability` feeds `authorizeHostedIdentity`, and the production `/api/auth/device/*` path consumes that authorization. Positive/negative role-capability regressions are executable.
  - **HOST-006 VERIFIED:** durable STALE device retirement, LOST/REVOKED handling, bound-session revocation, cross-tenant denial and persistent privileged/tenant-scoped security audit are implemented/regression-protected.
  - **HOST-007 VERIFIED for the v19.0 core:** canonical identity/session owners implement a production-wired Ed25519 public-key MFA ceremony with durable credentials, session-bound one-time challenges, domain-separated signing payload, persisted challenge hash, signature verification, replay/cross-session/expiry rejection, credential revocation and restart persistence. Applicable #164 login/bootstrap/session discovery/renewal failure cleanup/rotation/reauth/logout/device and session revocation/role-downgrade/product-entitlement propagation evidence is executable.
  - **#164 remains OPEN only for later v19.3 UX/client parity:** safe app-context restoration and platform-specific Mac/Windows/Web login/reauth/deep-link presentation. That residual does not reopen HOST-007 unless it discovers a core security defect.
- `HOST-008..009` product entitlement/quota: **VERIFIED**.
- `HOST-010..012` privacy lifecycle: **EARLIEST OPEN BAND / PARTIALLY IMPLEMENTED**. Existing production-wired evidence includes versioned data inventory/classification, secret-free account export, destructive account deletion, sanitized tombstones and archive anti-resurrection. Remaining Development work includes account deactivation lifecycle, durable privacy-request audit/evidence, cross-store deletion/retention proof, and real backup/PITR/operator recovery deletion behavior where dependent on HOST-016.
- `HOST-013..014` environment/service trust: **OPEN / CONTRACT IMPLEMENTED, REAL INFRASTRUCTURE NOT PROVEN**.
- `HOST-015..016` PostgreSQL tenancy/recovery: **OPEN / PARTIALLY IMPLEMENTED**; tenant-owned DB isolation and executable HA/failover/backup/PITR/restore proof remain technical gaps.
- `HOST-017..020`: **OPEN** for managed secrets/KMS and supply-chain/deploy provenance technical evidence.
- `HOST-021..022`: **OPEN** for measured provider/Data Health scorecards and canonical point-in-time/revision/no-lookahead truth; these are technical/full-data-capability gaps, not provider-licensing paperwork.
- `HOST-023`: **OPEN** for aggregate Development Production Ready qualification after applicable HOST-001..022 technical closure.
- Final identical-head Fast + impact-selected Qualified: **OPEN**.

No Commercial/Public-only approval may masquerade as a technical v19.0.0 blocker. Conversely, no readiness-tier separation may waive the genuine technical gaps above.

Do not begin v19.1.0 / Hosted Provider Gateway while #148 remains technically incomplete.

## Approved future provider-registry / Market Data rebaseline

A new provider-onboarding architecture decision is now **durable approved future scope**. It does not change the current v19.0 next action and must not be implemented early.

Read these files before v19.1 G1/#153 provider work:
1. `adaptive-governance/ADAPTIVE_PROVIDER_REGISTRY_CONTRACT.md`
2. `governance/V19_V20_PROVIDER_REGISTRY_REBASELINE.md`
3. `governance/programs/V19-V20-REBASELINE/adaptive-provider-registry.json`
4. `handoff/PROVIDER-REGISTRY-REBASELINE.md`
5. issue #153 current body/comments.

Approved target:

`Provider Adapter -> Adaptive Provider Registry -> capability/entitlement probes -> rights/authority + Data Health -> Smart Provider Router v2 -> canonical state -> all useful consumers`

Permanent rules:
- provider-specific implementation stops at one standards-compliant adapter;
- adapters self-register with the Registry;
- consumers request capabilities rather than provider names;
- technical capability/configuration/entitlement/freshness/history/quota/health discovery is generic where observable;
- Smart Provider Router v2 remains sole general routing/admission authority;
- runtime technical eligibility, fallback, demotion, cooldown, recovery and plan/subscription reprobe may be automatic;
- lifecycle/authority promotion (`SHADOW -> VALIDATED -> APPROVED -> PRODUCTION`) remains explicitly governed and never automatic;
- direct authorities and provider public/commercial rights cannot be replaced/inferred automatically;
- every provider capability must cross-integrate through canonical state to every useful applicable consumer, including Research, Discovery/Radar, Desks, Prep, Market Intelligence/Regime contribution, alerts, history/options, Data Health/Maintenance and future Outcome Learning;
- Settings provider cards/API-key UX should become metadata-driven and reuse canonical secrets/test/clear/redaction owners.

Market Data (`marketdata.app`) is the first concrete adopter under v19.1/#153:
- `MARKETDATA_TOKEN` environment fallback;
- Bearer auth;
- canonical Data Providers Settings token UX using the TradeInsight-derived preserve/replace/clear/test/redaction pattern;
- current HTTP 200/203 success semantics;
- current trial/delayed capability represented truthfully;
- effective entitlement reprobe so a future paid/live subscription can expand technical eligibility without a DE.PULSE source release solely because the provider plan changed;
- SHADOW-first initial lifecycle;
- no automatic lifecycle/authority promotion;
- Router scoring/eligibility remains capability/authority/freshness/health/latency/headroom/cost/utility/rights aware.

Version placement:
- v19.1: Registry runtime + generic adapter contract + Market Data + Router/cross-integration.
- v19.3: role-aware Mac/Windows/Web provider Settings/Admin presentation using the same contract.
- v19.6.1: provider reliability/coverage/economics/readiness scorecards including Market Data.
- v20.5: bounded adaptive provider utility/cost/reliability priors; Router/lifecycle/rights authority remains intact.

## Exactly one next action

Continue **HOST-010..012 — Privacy Lifecycle** through the existing canonical identity/account/workspace/persistence/archive/security-audit owners. Implement durable self-deactivation with session revocation and last-critical-owner protection; record minimal durable privacy-request audit evidence for export/deactivation/deletion; preserve existing export/delete/tombstone/anti-resurrection behavior; add positive/adverse/restart regressions; and keep real hosted backup/PITR/operator-recovery proof explicitly dependent on HOST-016 where current source cannot truthfully prove it. Do not start HOST-013+ until this band satisfies the Zero-Miss chain and exact-head CI/required real evidence disposition.

## Later dependency bands

After HOST-010..012 is VERIFIED:

1. HOST-013..014: reproducible IaC/service identity/network/TLS/mTLS enforcement and real drift evidence.
2. HOST-015..016: tenant-owned/scoped PostgreSQL persistence + tagged real DB isolation/migration/adverse/recovery/failover/PITR/restore proof; this band also closes any HOST-010..012 recovery evidence explicitly dependent on real backup/PITR/operator recovery.
3. HOST-017..020: managed secrets/KMS and supply-chain/deploy provenance.
4. HOST-021..022: measured provider scorecards/Data Health and canonical point-in-time/no-lookahead truth.
5. HOST-023 + final identical-head Development Production Ready qualification.
6. Commercial/Public activation remains a later separate gate: provider licensing/rights + public-user legal/compliance + commercial activation audit.
7. Only after v19.0 Development closure may v19.1 begin, including #153 Adaptive Provider Registry / Market Data work.

## Permanent architecture/product boundaries

- U.S. Equities Processing only; GLD/SLV/USO remain actionable tradable exceptions.
- No Execution.
- Smart Provider Router v2 is the sole general routing/admission owner.
- Adaptive Provider Registry is a provider-registration/capability projection only; never a second Router.
- Direct SEC/EDGAR remains the governed filing/Form 4 authority.
- Extend canonical Data Health/freshness/cache/persistence/subscription/telemetry/reconciliation/lifecycle/identity/session owners; never create parallel owners.
- No automatic provider lifecycle/authority promotion.
- Point-in-time/no-lookahead truth precedes adaptive learning.
- Development provider-rights mode remains audit-only; public-production enforcement requires explicit Commercial/Public activation.
- Preserve existing product look-and-feel unless a justified truthful integration/defect correction requires change.

## Resume rule

1. Read this file first.
2. Read `governance/PRODUCTION-READINESS-TIERS.md` and its machine companion before interpreting any readiness/blocker status.
3. Then read the strict zero-miss contract/decision, `governance/current-state.json`, `governance/ROADMAP.md`, `governance/V19_V20_REBASELINE.md`, the active work-slice `work-slice.json`, `g1-scope.json`, `closure.json`, issue #148/latest comments, issue #164, PR #149 and current Actions.
4. The provider-registry addendum is approved future scope; read its four files above before any v19.1/#153 provider implementation. Do not use it to bypass v19.0 sequencing.
5. Where older wording conflicts with the readiness-tier contract, the readiness-tier contract wins and the conflicting artifact must be corrected at the next meaningful governance transition.
6. Fetch live `main` and active branch heads before making any change; another session may have advanced them.
7. Inspect commits since the latest verified checkpoint so nothing already implemented is duplicated.
8. Keep Stable v18.10.0 immutable.
9. Do not create another branch/PR, merge/release, start v19.1.0, or weaken G0-G16/source/data-truth/security/platform/CI gates.
10. Update durable GitHub state at meaningful dependency-band transitions so any AI/account can resume independently.
