# T1 — Shipped v18 Feature Discovery Baseline

**Issue:** #114  
**Parent:** #113 / `ADAPT-V18-FINAL-CLOSURE-10-10-001`  
**Baseline main:** `c89fd8e94d8afa2aec5bb3adc666dcb918b721cd`  
**Branch:** `adapt-v18-closure-t1-traceability`  
**Status:** DISCOVERY IN PROGRESS — no T1 or parent assurance gap is VERIFIED by this document.

## Sources inspected so far

1. Live repository tree at the exact baseline.
2. `renderer/index.html` primary navigation, shell, Data Engine/sidebar and loaded capability owners.
3. `governance/registries/functionality_utility_registry.json` — permanent utility/surface registry covering material tabs, engines, checkpoints and watchers.
4. `governance/data-health/provider-capability-matrix.json` — provider/capability/consumer/authority/freshness ownership.
5. `tests/integration/http_workflow_test.py` — authenticated runtime, API, watchlist/desk, Settings, Data Engine, Discovery, event, validation and runtime structure workflows.
6. `tests/renderer/` — renderer logic, Research correctness/IA, responsive UI, symbol/desk, documentation and surface-consolidation owners.
7. `tests/acceptance/` — professional expert/runtime and trader acceptance owners.
8. GitHub v18 issue/program history including #11, #12, #57, #61, #63/#64, #65, #70, #73, #76, #78–#84, #92, #94, #95, #102 and #107.
9. `adaptive-governance/FUNCTIONALITY_UTILITY_INTEGRATION_CHECKPOINT.md` and canonical Roadmap/handoff.
10. Current Stable/release identity and the v18.10 program registration evidence.

## Current material surface census

Primary navigation currently exposes Dashboard, Market Intelligence, Day Trade Desk, Swing Desk, Long-Term Desk, Discovery, Research, AI Copilot, Maintenance, Settings and Documentation. Administration is a capability-scoped conditional surface. The shell also owns market session/time, runtime state/start-stop, identity/sign-out, macro-event banner, ticker tape and the Data Engine sidebar.

The permanent utility registry additionally identifies Opportunity Radar, Rapid Move / Market Shock, Pre-Market Prep, Market Open Prep, Earnings & Material Catalyst Reaction Watch and Market Activity Seeds as material engines/checkpoints even when they are embedded rather than primary pages.

## Current non-surface capability census

The integration/runtime and provider registries prove additional shipped capability families that must not disappear under generic labels:
- secure owner bootstrap, session authentication and role/capability authorization;
- runtime lifecycle and demo provenance;
- canonical Day/Swing/Long watchlist membership, final-desk protection and Master Symbol remove/undo/remove-all;
- U.S. symbol universe/eligibility and canonical instrument identity;
- Smart Provider Router v2, provider registration/onboarding, capability/entitlement state, provider lifecycle/readiness, transport diagnostics and semantic usefulness telemetry;
- canonical freshness/Data Health, provider reconciliation, live allocation/subscription and broad snapshot reuse;
- live quotes, historical bars, VIX/index, global market context, macro/official data and calendars;
- News, Earnings, SEC/Ownership/Form 4, corporate actions and TradeInsight governed evidence;
- Options and alternative intelligence, including untrusted community evidence;
- Event Intelligence, smart notifications, catalyst reactions, both preparation checkpoints, Decision Queue/Trade Readiness and market modes/regime;
- signal validation, seasonality/calibration/concentration learning, evidence snapshots/provenance, research package and SHADOW/adaptive-data-policy controls;
- Settings/secrets, persistence/migrations/restart, release identity/provenance and operational load diagnostics.

## Discovery rules

- Existing issue closure is evidence provenance, not automatic feature verification.
- Existing test files are candidate regression owners, not proof of complete evidence.
- Every user-visible row will enter T7 with `PENDING_T7` and must receive KEEP/MOVE/MERGE/REMOVE/RENAME/REDESIGN.
- Every feature row starts `DISCOVERED_UNVERIFIED` until the relevant assurance tracks bind executable evidence.
- A duplicate canonical owner, missing owner, no meaningful test owner, or source/issue mismatch becomes an explicit blocker/corrective under #113.
- T1 does not implement fixes; it discovers and disposes them.

## Next T1 work

Populate `feature-assurance-ledger.json` with the current feature-level census, then run a second source/test/issue reconciliation pass for omissions, duplicates and unowned rows before freezing T1.
