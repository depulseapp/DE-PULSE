# DE.PULSE Production Readiness Tier Contract

**Status:** AUTHORITATIVE PERMANENT GOVERNANCE  
**Applies to:** all active and future DE.PULSE versions, work slices, closure ledgers, handoffs, CI/release claims and public/commercial activation decisions.  
**Precedence:** this contract supersedes any older repository wording that conflates provider licensing/rights or public-user commercial/legal approval with technical Development Production Ready closure.

## 1. Canonical readiness definitions

> **Development Production Ready = technically robust, secure, persistent, cross-platform, tested, full provider/data capability.**
>
> **Commercial/Public Ready = Development Production Ready + provider licensing/rights + public-user legal/compliance + commercial activation audit.**

These are two distinct readiness claims. They must never be collapsed into one gate.

### 1.1 Evidence-based 5/5 maturity target

The long-term technical/product target is **5 / 5 in every Development maturity domain**, not a rounded average and not a documentation score. The audit baseline domains are:

- Local product usefulness;
- Deterministic intelligence;
- Provider architecture;
- Canonical state;
- Adaptive intelligence;
- Persistence;
- Security foundation;
- Testing / release assurance;
- Web + multi-platform readiness.

A domain may be called `5 / 5` only when current-source, production-reachable behavior and required executable evidence show no known material gap for the intended Development target. Strong scores in other domains cannot compensate for a weaker domain. Missing real infrastructure/provider/platform evidence cannot be replaced by prose, mocks or a self-assigned score.

**Capability conservation is mandatory while improving the score.** Rebaseline/refactor/consolidation may not silently drop, weaken, bypass or reduce an existing certified/approved capability to make architecture cleaner or scoring easier. A replacement must prove functional equivalence or improvement for all applicable consumers, roles, platforms, persistence/recovery behavior and regressions before the prior owner can be retired. Intentional capability removal requires an explicit user-approved product decision and durable traceability.

DE.PULSE remains in **Development** until the user explicitly declares the product ready to enter Commercial/Public readiness work. Until that explicit declaration, **Commercial distribution is `NOT_ACTIVATED` rather than a Development maturity failure**. Commercial/Public distribution itself has the same target of `5 / 5` once activated, but provider licensing/rights, public-user legal/compliance and the commercial activation audit remain separate activation obligations and may not be used to cripple Development provider/data capability.

The 5/5 target is governed by the Strict Zero-Miss Implementation & Closure Contract: no maturity claim may conceal an implementation miss, unconsumed helper, duplicate/parallel authority, untested adverse path, missing persistence/restart behavior, unresolved cross-platform gap, or known material defect.

## 2. Development Production Ready

A Development Production Ready claim requires all applicable technical obligations to be closed through the normal zero-miss lifecycle and executable evidence. At minimum, where applicable, this includes:

- technically robust production-reachable implementation through canonical owners;
- security, tenant/account isolation, authentication/session/device controls, authorization and fail-closed behavior;
- persistent canonical state, restart/lifecycle/recovery behavior and truthful degradation;
- required Mac Apple Silicon + Windows x64 parity and Web parity for shared hosted capabilities unless explicitly justified N/A;
- full applicable automated testing, adverse-path evidence, exact-head CI and impact-selected qualification;
- full configured/operational provider and data capability needed by the product, without suppressing technical capability merely because later commercial rights approval is unfinished;
- provider-rights metadata, provenance, evaluator, audit visibility, policy separation and an explicit public-production fail-closed activation mechanism;
- privacy/data lifecycle, deletion/retention, secrets, infrastructure, database tenancy/recovery, observability, supply-chain, data-truth and no-lookahead controls where they are part of the technical product.

Development Production Ready does **not** mean provider-specific commercial/public rights have been granted. Working credentials, successful provider calls, public terms, development use or a technical rights-control implementation never imply legal/commercial approval.

## 3. Commercial/Public Ready

Commercial/Public Ready may be claimed only after Development Production Ready is already satisfied and all applicable activation-only obligations are complete:

1. provider-specific licensing/rights approvals for the intended public/commercial use, including applicable redistribution, display, proxy/cache/retention, derived/AI and multi-user rights;
2. public-user legal/compliance approvals and required public-facing legal/compliance artifacts;
3. a commercial activation audit proving the intended deployment, providers, rights, user exposure, controls and evidence match the approved commercial/public operating model.

`Commercial/Public Ready => Development Production Ready`.

The reverse is never implied. Development Production Ready alone never authorizes opening DE.PULSE to public/commercial users.

## 4. Provider/data capability rule

During development and pre-public validation, DE.PULSE must retain the full configured and operationally eligible provider/data capability required for technical validation. Unfinished provider licensing/rights approval must not be used to cripple Smart Provider Router v2, live subscriptions, cache, persistence, serving or other technical data paths merely to obtain a readiness label.

The technical provider-rights control plane remains a Development obligation: rights metadata/provenance/evaluation must be truthful and auditable, and the public-production fail-closed path must be implemented and tested.

Actual provider-specific legal/commercial approval is a Commercial/Public activation obligation. It is tracked separately and cannot be cited as a missing Development implementation artifact.

For the current hosted foundation, hard provider-rights enforcement remains gated by the explicit public-production activation mode. The existence of audit-only development behavior grants no provider rights.

## 5. Technical vs activation-only classification

The following remain **Development technical obligations** when applicable and must not be deferred merely because they have privacy/security/rights implications:

- tenant/account isolation and RBAC/product-entitlement separation;
- authentication, MFA/passkey-class proof, session/device lifecycle and revocation;
- privacy data inventory, minimization, retention, deletion/deactivation, tombstones and recovery behavior;
- cross-store deletion and backup/PITR behavior;
- secrets/KMS, service identity, network/TLS/mTLS and IaC enforcement;
- PostgreSQL tenant ownership/isolation, migrations, HA/failover/backup/PITR/restore;
- supply-chain/provenance/SBOM/signing/attestation where applicable;
- provider scorecards, Data Health, freshness, point-in-time/revision/no-lookahead truth;
- cross-platform behavior and exact-head executable qualification.

The following are **Commercial/Public activation-only obligations** unless a specific requirement explicitly makes an additional technical implementation necessary:

- final provider-specific licence/contract/formal public/commercial reuse approval;
- public-user legal/compliance approval and public-facing legal/compliance acceptance;
- final commercial activation audit and authorization to enable public production.

## 6. Closure and blocker semantics

Zero-miss closure rigor applies independently within the readiness tier being claimed.

Technical lifecycle remains:

`PLANNED -> IMPLEMENTED -> PRODUCTION_INTEGRATED -> VERIFIED -> RELEASE_QUALIFIED`

A Development requirement may be OPEN only for a Development-applicable technical/evidence gap. Missing Commercial/Public-only paperwork must be recorded as a separate activation gate and must not make an otherwise satisfied Development technical requirement OPEN, `BLOCKED_EXTERNAL`, or `NOT_APPLICABLE` for the wrong reason.

External blockers must name the readiness tier they block.

A requirement cannot be marked VERIFIED from types, helpers, policy prose, unit tests or self-declared runtime configuration alone when production reachability, adverse behavior, persistence, security, external infrastructure/data evidence or exact-head qualification is applicable.

## 7. Public-production authorization

Public/commercial activation remains fail closed.

`publicProductionAuthorized = false` until:

- Development Production Ready is complete for the intended public candidate;
- all intended provider licensing/rights approvals are current and applicable;
- public-user legal/compliance gates pass;
- the commercial activation audit passes;
- the explicit activation decision is recorded.

No development build, CI pass, provider credential or Development Production Ready declaration may silently flip this state.

## 8. Current v19.0.0 interpretation

The active v19.0.0 Hosted Trust & Identity Foundation is being closed against **Development Production Ready**. HOST-001..HOST-023 therefore close on their technical responsibilities and executable evidence.

HOST-001..003 may be Development-verified when the provider-rights control plane, provenance, audit-only development behavior and explicit public-production fail-closed path are technically proven. Actual provider licensing/rights approval remains a separately tracked Commercial/Public activation gate.

All genuine technical HOST gaps remain blockers. In particular, reclassifying commercial provider approval does not waive tenant isolation, MFA/session/device lifecycle, privacy/recovery, IaC/service trust, PostgreSQL tenancy/recovery, managed secrets, supply chain, provider/Data Health, point-in-time/no-lookahead or final qualification requirements.

## 9. Governance conservation rule

Every roadmap, work slice, closure ledger, handoff, issue, PR and release decision must preserve these two readiness tiers and the evidence-based 5/5 maturity target. If older wording conflicts, this contract wins and the conflicting artifact must be corrected at the next meaningful governance transition.

No future AI, account, implementation session or release process may reinterpret Commercial/Public-only obligations as Development technical blockers, use readiness separation to waive genuine technical closure evidence, or improve a maturity score by silently reducing product capability.