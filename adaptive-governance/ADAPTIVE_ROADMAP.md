# DE.PULSE — Active Release Operational Roadmap Overlay

**Status:** OPERATIONAL OVERLAY / NOT PRODUCT-SEQUENCE AUTHORITY  
**Canonical product roadmap:** `governance/ROADMAP.md`  
**Canonical v19/v20 rebaseline:** `governance/V19_V20_REBASELINE.md`  
**Permanent operating contract:** `governance/ADAPTIVE-OPERATING-CONTRACT.md`  
**Current projection:** `adaptive-governance/CURRENT_ADAPTIVE_ROADMAP.md`

This overlay carries active-release execution obligations only. It does not redefine product sequencing. If it conflicts with canonical governance, canonical governance wins.

Permanent top-level release model: **G0-G16 only. No G17+.**

## 1. Current product state

- Certified immutable Stable: `v18.10.0`.
- v18 T1-T10 and 180 shipped-responsibility assurance remain frozen historical baseline evidence.
- Post-v18 source-overlap audit passed.
- Active version: `v19.0.0` Hosted Trust & Identity Foundation.
- Active issue/PR/branch: #148 / PR #149 / `adapt-hosted-trust-foundation-001`.
- Current PR head at rebaseline audit was not qualified: Fast #1141 failed source-health because hosted identity/session helpers were not yet production-reachable.

## 2. Recoverability / resume

Every active version must be resumable from GitHub/CI evidence without conversation history. Durable state must identify:
- incoming Stable;
- active version, branch, issue and PR;
- frozen requirement/backlog maps;
- exact source head and qualification status;
- current blocker/next action;
- closure/evidence owners;
- release artifacts/provenance when created.

Conversation interruption does not invalidate unchanged-source evidence. Source/tooling changes invalidate affected/dependent evidence only.

## 3. Version-first adaptive execution

Current future planning uses coherent **versions/builds**, not requirement packets.

`diff -> impact map -> source overlap -> canonical-owner decision -> implementation/evidence -> Delta AIPLC -> exact-head Fast -> risk-selected Qualified -> G10 whole-version reconciliation -> immutable RC -> G11-G16`

Requirements/backlog bullets remain granular evidence rows. They do not become branches, PRs, public versions or release events.

Heavy work may receive a real patch version when that materially improves correctness/risk isolation. Smaller related work should be combined when it shares architecture and evidence.

## 4. Zero-miss adaptive audit

For every changed responsibility, bind:
- requirement/backlog/source provenance;
- canonical owner and duplicate/consolidation disposition;
- upstream evidence and downstream consumers;
- trader/investor utility;
- positive/negative/failure evidence;
- freshness/degradation/recovery behavior;
- persistence/restart/migration;
- role/RBAC/product entitlement/provider rights;
- load/backpressure/resource behavior;
- required Mac/Windows/Web evidence;
- durable regression ownership.

For intelligence work also bind #170:
- cross-integration `REQUIRED / CONDITIONAL / NOT_USEFUL`;
- Market Regime `YES / CONDITIONAL / NO`;
- Outcome Learning `YES / CONDITIONAL / NO`;
- point-in-time/no-lookahead lineage;
- intelligence maturity and deterministic fallback.

For visible work also bind #171:
- `KEEP / IMPROVE / CONSOLIDATE / COLLAPSE / MOVE / REMOVE`;
- decision value and cognitive-load rationale;
- role-aware composition;
- no deletion of useful canonical evidence solely to simplify a page.

## 5. Smart/intelligent maturity

Permanent target:

`canonical evidence -> deterministic intelligence -> cross-feature synthesis -> point-in-time outcomes -> bounded adaptive learning -> optional AI/agent synthesis/orchestration`

Do not equate intelligence with more UI or more LLM usage. Deterministic logic is preferred where it is safer and sufficient. Learning must be bounded, auditable, shadow-first and reversible. AI/agents cannot independently fetch provider truth, invent missing evidence, bypass Router/rights/authority or become required for core product behavior.

## 6. UI protections

- Preserve Day/Swing/Long Desk look-and-feel/workflow.
- Preserve Dashboard Market Regime and Desk Control materially.
- Preserve Data Engine look-and-feel except genuine defects.
- Preserve current AI Copilot engine/header visual treatment unless separately justified.
- Reduce crowding through synthesis/progressive disclosure/correct placement rather than deleting analytical depth.

## 7. Current version sequence

Canonical sequence and exact backlog/HOST mapping are owned by `governance/ROADMAP.md` and `governance/V19_V20_REBASELINE.md`. Machine maps:
- `governance/programs/V19-V20-REBASELINE/backlog-version-matrix.json`;
- `governance/programs/V19-V20-REBASELINE/host-requirement-version-map.json`.

This overlay must not reproduce stale historical future-version reservations.

## 8. Exactly one next action

Continue `v19.0.0` on existing PR #149. Production-wire or correctly consolidate the hosted identity/session helpers through existing canonical auth/HTTP owners, then obtain fresh exact-head Fast. Do not start `v19.1.0` while v19.0.0 remains unqualified.
