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
- Canonical closure ledger: `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/closure.json`.
- Canonical work-slice state: `governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/work-slice.json`.
- GitHub objects and executable evidence outrank this file and all chat memory. Always fetch live `main`, the active branch, PR #149, issue #148/comments and Actions before writing.

## Permanent zero-miss lifecycle

`PLANNED -> IMPLEMENTED -> PRODUCTION_INTEGRATED -> VERIFIED -> RELEASE_QUALIFIED`

OPEN/PARTIAL/IMPLEMENTED_UNVERIFIED/CONTRACT_ONLY/FRAMEWORK_ONLY are incomplete states only. They never authorize dependency-band advancement, merge, release, production-ready/10/10 claims, or handoff-as-done.

Before VERIFIED, prove the applicable chain:

`requirement -> canonical owner -> production integration -> consumer reachability -> positive behavior -> adverse/fail-closed behavior -> persistence/restart/lifecycle -> security/rights/privacy -> observability -> executable regression -> real external/infrastructure evidence where required -> exact-head CI -> closure ledger`

A newly discovered implementation miss automatically reopens the requirement and dependent completion claims.

## Current v19 truth

- `HOST-001..003` provider rights: **OPEN / PRODUCTION-INTEGRATED / NOT VERIFIED** — earliest dependency band. Do not advance.
- `HOST-004..007` tenant/account/device/session/reauth: **OPEN / PARTIALLY IMPLEMENTED**.
- `HOST-008..009` product entitlement/quota: **VERIFIED**.
- `HOST-010..012` privacy lifecycle: **OPEN / PARTIALLY IMPLEMENTED**.
- `HOST-013..014` environment/service trust: **OPEN / CONTRACT IMPLEMENTED, REAL INFRASTRUCTURE NOT PROVEN**.
- `HOST-015..016` PostgreSQL tenancy/recovery: **OPEN / PARTIALLY IMPLEMENTED**.
- `HOST-017..023`: **OPEN**.
- Final identical-head Fast + impact-selected Qualified: **OPEN**.

Do not begin v19.1.0 / Hosted Provider Gateway while #148 remains incomplete.

## HOST-001..003 — implemented application-side boundary

The strict audit found missing production consumption and later found provider-inventory/source-identity gaps. Those application-side gaps are now remediated through existing canonical owners only:

- `provider_data_rights.go` remains the canonical rights owner.
- Rights metadata binds explicit provider identity and a SHA-256-pinned reviewed evidence bundle; missing/malformed/mismatched/unreviewed evidence fails closed.
- A working key, successful request or public web terms never grant executable rights.
- Smart Provider Router v2 remains the sole general routing/admission owner. Hosted rights constrain eligibility/admission but do not alter Smart Router scoring/order.
- Live subscription desired sets re-evaluate rights; downgrade/expiry removes future fanout.
- Market cache replay/save and canonical persistence re-check hosted rights.
- HTTP/runtime serving checks hosted rights before per-user privacy projection.
- Provider-unbound legacy/derived payloads fail closed pending HOST-022 provenance.
- Public-terms review coverage is derived from all canonical `providerRegistrations()` identities: Alpaca, Finnhub, TradeInsight, Twelve Data, Marketaux, FRED, SEC, SEC EDGAR, yfinance, CBOE, BLS and EIA.
- Source-to-provider persistence resolution recognizes all registered provider identities, keeps `SEC` separate from `SEC EDGAR`, and prevents BLS/EIA/other registered external evidence from being treated as internal semantic evidence.
- Tests proving these paths are consolidated into existing root owners; permanent root/G16 policy was not weakened.

## HOST-001..003 — remaining external blocker

Real provider-specific approved rights evidence is still not bound. Therefore the band remains OPEN.

Private/commercial providers require the applicable agreement/licence/order form or equivalent reviewed approval evidence for the exact DE.PULSE use. Government/public providers may be closable from official reuse/API terms only after a documented owner/legal review maps every required hosted-use dimension and operational obligation.

Current public-terms triage is intentionally non-executable and non-approval:
`governance/work-slices/ADAPT-HOSTED-TRUST-FOUNDATION-001/provider-public-terms-review.json`

Notable current distinctions:
- BLS public data is public-domain/secondary-use friendly, but citation/access-date/non-vouching, truthful representation and API-limit obligations must be mapped and implemented.
- EIA permits API-backed search/display/analysis and public-domain data reuse/distribution with attribution, subject to protected third-party material and trademark restrictions.
- SEC/EDGAR public filing/data access is broadly available but fair-access controls and third-party/filer-supplied content exceptions prevent treating it as an automatic blanket redistribution/AI grant.
- FRED remains materially restrictive, including series-owner rights, caching/archiving and AI/ML constraints.
- Private/commercial feeds remain blocked until account-specific rights evidence exists.

Do not create fictional approvals or convert the public-terms review into the executable provider-rights bundle.

## Exactly one next action

Stay on HOST-001..003.

1. Fetch live GitHub state and reconcile any newer commits first.
2. Obtain/review real provider-specific rights evidence for all applicable providers and use dimensions.
3. Map commercial/multi-user, proxy, cache/retention, redistribution/display, derived/AI use, environment, limits/attribution and effective/expiry/review requirements.
4. Bind only actually approved records through the existing SHA-256-pinned provider-rights bundle.
5. Prove real hosted positive and fail-closed behavior.
6. Obtain exact-head Fast and the governed impact-selected Qualified checkpoint on the resulting candidate.
7. Only then consider HOST-001..003 VERIFIED and move to HOST-004..007.

## Later dependency bands

After HOST-001..003 is VERIFIED:

1. HOST-004..007: auditable stale-device retirement/security audit + actual MFA/passkey-class proof + applicable #164 end-to-end session lifecycle evidence.
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
- Preserve existing product look-and-feel unless a justified truthful integration/defect correction requires change.

## Resume rule

1. Read this file first.
2. Then read the strict zero-miss contract/decision, `governance/current-state.json`, `governance/ROADMAP.md`, `governance/V19_V20_REBASELINE.md`, the active work-slice `work-slice.json`, `g1-scope.json`, `closure.json`, issue #148/latest comments, issue #164, PR #149 and current Actions.
3. Fetch live `main` and active branch heads before making any change; another session may have advanced them.
4. Inspect commits since the latest verified checkpoint so nothing already implemented is duplicated.
5. Keep Stable v18.10.0 immutable.
6. Do not create another branch/PR, merge/release, start v19.1.0, or weaken G0-G16/source/data-truth/security/platform/CI gates.
7. Update durable GitHub state at meaningful dependency-band transitions so any AI/account can resume independently.
