# DE.PULSE — Functionality Utility, Reuse, Correlation & Surface Checkpoint

Status: **Permanent blocking checkpoint inside canonical G0–G16**  
Effective: **v18.2 and all later releases**

This checkpoint creates no new G-gate. It is executed during G1/G2/G3, re-verified at G9/G10, and closed at G16.

## Purpose

DE.PULSE must not accumulate features, jobs, datasets, calculations, cards, tabs or background work simply because they can be built. Every new or materially changed capability must prove that it adds distinct decision-support value, reuses canonical evidence where possible, correlates with the rest of the intelligence system, and is presented in the smallest useful surface.

The governing principle is:

**Need → reuse → correlate → consolidate → place → measure.**

The default is **not** to create a new tab, fetch path, scheduler, data store, metric or parallel engine.

## Mandatory scope

The checkpoint applies to every proposed or materially changed:

- application tab, sub-tab, page, card, panel, table or shell surface;
- scanner, detector, ranking engine, preparation/checkpoint job, watcher, scheduler or background task;
- provider/API request, subscription, cache, dataset or persistence field;
- metric, score, derived feature, model, AI context item or learning signal;
- alert, notification, decision-queue input or readiness input;
- administrative/security operation or role-sensitive payload;
- research/evidence workflow and drill-down path.

## Blocking questions

Every item must answer all applicable questions before G3 exit:

1. **Purpose** — What user/system problem does this solve? What decision or workflow consumes it?
2. **Canonical owner** — Which existing owner should own it? If an owner already exists, extend/reuse that owner instead of creating a parallel implementation.
3. **Existing evidence reuse** — Which already-fetched, already-computed or already-persisted evidence can satisfy some or all of the need?
4. **Duplicate acquisition** — Can the new requirement share a canonical fetch, subscription, snapshot, cache, calculation or in-flight request instead of calling the provider again?
5. **Duplicate computation** — Is the same or substantially similar metric/intelligence already calculated elsewhere? Prefer calculate-once/reuse.
6. **Functional overlap** — Does an existing tab, job, detector, watcher or workflow already solve the same problem? If partially, consolidate or extend before adding another owner.
7. **Correlation/integration** — Which existing evidence, modules, events, outcomes or user context must this correlate with so it becomes intelligence rather than an isolated datum?
8. **Materiality/freshness** — When is it useful, stale, irrelevant or suppressible? What happens when evidence is absent, contradictory or degraded?
9. **Retention** — Does it need persistence? If yes, why, for how long, at what granularity, and for which future consumer/learning purpose?
10. **Performance/provider cost** — What provider calls, CPU, memory, storage, fan-out, scheduler pressure or UI rerender load does it add? Can work be coalesced, bounded, incremental, event-driven or moved off the critical path?
11. **Surface necessity** — Does it need any visible UI at all? If visible, can it be a compact context line, existing section, drill-down or Maintenance diagnostic instead of a new tab/card?
12. **New-tab justification** — A new tab is allowed only when existing surfaces cannot preserve clear workflow, hierarchy, security/capability separation or maintainability. Organizational convenience is not sufficient.
13. **Role/security composition** — Who may see/use it, which payload fields must be withheld server-side, and how does the remaining layout reflow for roles without it?
14. **Outcome/learning** — For intelligent capabilities, what success/failure/outcome evidence will show whether it is useful enough to retain or needs consolidation/removal?
15. **Deletion/retirement** — What existing obsolete, duplicated or superseded implementation can be removed as part of this change?

## Required disposition

Each audited item receives one or more explicit dispositions:

- `REUSE` — existing canonical implementation/data fully satisfies the need;
- `EXTEND_EXISTING` — add the capability to the existing canonical owner;
- `CONSOLIDATE` — merge overlapping implementations/surfaces/jobs;
- `DRILLDOWN_ONLY` — keep detailed evidence, but only behind the canonical deep-evidence home;
- `INTERNAL_ONLY` — useful operational/supporting evidence, not a normal user surface;
- `NEW_SURFACE_JUSTIFIED` — distinct visible surface is materially necessary and documented;
- `REMOVE_OR_RETIRE` — obsolete/redundant owner, route, surface or work should be removed;
- `DEFER` — insufficient utility/consumer/readiness to build now.

A new implementation cannot proceed merely with `NEW`. It must resolve to one of the dispositions above.

## UI surface contract

The default visual rule is:

**one canonical intelligence owner → one deep-evidence home → concise contextual reuse elsewhere.**

Do not duplicate full evidence packages across Dashboard, Desks, Discovery, Research and Market Intelligence. A surface may reuse a compact conclusion, material risk, freshness state or navigation link without becoming another deep-evidence owner.

A proposed new tab must prove at least one of:

- a genuinely distinct high-frequency workflow;
- a materially different information hierarchy that would damage an existing page if merged;
- a security/capability boundary that is cleaner and safer when separated;
- a deep operational/research function that cannot be represented as progressive disclosure or drill-down;
- maintainability benefit that also improves user clarity, not merely developer organization.

Otherwise the item must be integrated into an existing surface or remain internal.

## Data and correlation contract

Raw availability never justifies collection or display. New data must declare:

- canonical owner;
- active consumer(s);
- reuse relationship to existing data;
- correlation targets;
- freshness/materiality semantics;
- retention/rights/sensitivity requirements;
- failure/degraded behavior;
- provider/performance budget;
- whether the datum influences deterministic logic, contextual intelligence only, diagnostics only, or adaptive learning only.

Independent corroboration must remain distinct from duplicated/syndicated evidence; repeated copies cannot inflate confidence.

## Provider → Market Mode assessment

Every new or materially changed provider capability must receive a durable Market Mode disposition before G3 and be revalidated at G7/G10:

- `INTEGRATED` — canonical, freshness-qualified evidence may enter named Market Mode scopes only after its lifecycle permits production influence;
- `CONTEXTUAL_ONLY` — may explain risk, confidence, contradictions or material changes, but cannot independently control a mode;
- `NOT_RELEVANT` — useful elsewhere, with no Market Mode consumer;
- `INTENTIONALLY_HIDDEN` — withheld from Market Modes because of rights, safety, quality or product-boundary constraints.

Provider count never changes a mode. Every input must flow through the canonical Smart Provider Router and canonical state; no provider-specific mode engine is allowed. The assessment must name dataset/capability, router role, canonical owner, consumers, Market Mode scopes, freshness/independence/rights policy, degraded behavior, lifecycle, production influence and acceptance evidence.

The smart/AI/LLM-style behavior belongs in evidence selection, usefulness learning, contradiction detection, synthesis and explanation. Deterministic/statistical code continues to own price truth, numeric Market Mode calculation, calibration and promotion measurement. LLM output cannot directly set a Market Mode or silently reweight protected Day/Swing/Long logic.

All provider influence follows **SHADOW → VALIDATED → APPROVED → PRODUCTION**. `NOT_IMPLEMENTED`, `SHADOW` and `VALIDATED` rows have zero production Market Mode influence. Promotion requires sample sufficiency, cutoff-safe replay/live comparison, rollback and explicit approval.

Maintain `provider_market_mode_integration_registry.json`. The existing `functionality_utility_checkpoint_gate.py` validates current Router/capability providers against this registry, so adding a provider without an assessment fails closed. TradeInsight is the first explicit future-provider case: its historical/corporate-action evidence may target integrated Swing/Long/Sector/Industry use after validation; congressional/Form 4/metadata remain contextual; Top Movers is not relevant to modes; rights-gated AI/MCP evidence is intentionally hidden.

## Background-job / checkpoint contract

A new scheduled/background function must not become an independent mini-engine by default. It must first determine whether it is actually:

- a temporal checkpoint over existing canonical state;
- an event-triggered evaluation;
- a shared refresh request;
- a derived intelligence calculation;
- a maintenance/integrity operation.

When multiple jobs require the same evidence, they must request it through the canonical scheduler/router/cache and coalesce in-flight work where practical.

## G0–G16 enforcement

- **G1 — Immutable Scope:** enumerate every new/changed functionality, dataset, job and surface; record initial overlap/necessity disposition.
- **G2 — Architecture / Data Utility:** prove canonical ownership, reuse, correlation, freshness, retention and no parallel acquisition/computation.
- **G3 — Design / Dependency Readiness:** prove UI placement, new-tab necessity, role composition, dependencies and removal/consolidation plan before coding.
- **G4 — Development Exit:** implementation follows the approved owner/disposition; no convenience duplication introduced during coding.
- **G7 — Data / Security / Adaptive Intelligence:** verify data utility, evidence independence, sensitivity, point-in-time truth and adaptive-governance behavior.
- **G8 — Performance / Capacity / Stability:** verify added work is bounded and provider/runtime efficient; duplicate/in-flight work is coalesced where practical.
- **G9 — Cross-Module / UI / UX:** audit every tab/surface and major function for repeated information, duplicated workflows, hierarchy drift and unnecessary UI.
- **G10 — Pre-Freeze Qualification:** rerun the checkpoint against the actual candidate; unresolved duplication/orphan data/unjustified surface is blocking.
- **G16 — Adaptive Retrospective / Handoff:** record what was reused/consolidated/removed, identify new overlap discovered after implementation, and add regressions/governance learning.

## Required evidence

Every release must maintain `functionality_utility_registry.json`. The registry is reviewed by `functionality_utility_checkpoint_gate.py` and must cover all primary navigation tabs plus material engines/jobs introduced or changed by the release.

A release cannot pass G10 when:

- a primary navigation tab is absent from the registry;
- an item lacks owner/purpose/consumer/reuse/correlation/placement disposition;
- a new tab lacks documented separation justification;
- a retained overlapping function lacks a reason it must remain separate;
- a duplicated acquisition/computation path is knowingly left without an approved bounded exception;
- newly introduced data has no active/strategic consumer or required governance metadata;
- the candidate leaves obsolete duplicate surfaces/jobs without an explicit retirement decision.

## Permanent rule

**DE.PULSE should add the minimum new machinery necessary to produce materially better intelligence. Reuse and correlate first; consolidate before extending; create a new tab or engine only when separation is genuinely valuable.**
