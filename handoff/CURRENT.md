# DE.PULSE Current Handoff

**SUPERSEDES ALL PRIOR CHAT HANDOFFS**

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Certified Stable candidate:** `e55d8d25b15cec2ffb0f5411bc358bc40b359cf9`  
**Certified Stable source fingerprint:** `0062f46dea5690d0b3fcd8a9ed3b1f71ebe1522c7dee2cb218e9d36b9e0076ff`  
**Certified Stable build ID:** `v18.9.1-stable-20260821`  
**Immutable predecessor resume checkpoint release:** `v18.9.0` / `v18.9.0-stable`  
**Retained completed process authority:** #73 / `ADAPT-ROOT-CONVERGENCE-001` / `adapt-root-convergence-001` / closure `governance/work-slices/ADAPT-ROOT-CONVERGENCE-001/closure.json`  
**Completed provider/data-health program:** #79 with #80, #81, #82, #83, #78 and final #84 / PR #91 / merge `733d90ca125a4fe5abd38a2ea40de0623703dfd4`  
**Active product slice:** #92 / `ADAPT-COMPANY-INSTRUMENT-IDENTITY-001` / `adapt-company-instrument-identity-001` / closure `governance/work-slices/ADAPT-COMPANY-INSTRUMENT-IDENTITY-001/closure.json`  
**Parent program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`

## Current authority
Actual GitHub objects and executable evidence outrank this handoff. #73 / `ADAPT-ROOT-CONVERGENCE-001` / `adapt-root-convergence-001` remains the completed retained repository-architecture process authority represented by canonical `activeWorkSlice`; it is not reopened. #79/#84 are complete and must not be restarted. The separate `productCapabilityGate` now reserves #92 as the active canonical-company/instrument-identity packet under #65.

## #92 active contract
Identity must reuse the existing Smart Provider Router v2 Alpaca `US Asset Universe` acquisition and canonical symbol/persistence owners. The existing `/v2/assets` response already carries useful slow-changing identity fields; #92 captures them from that exact response with **zero additional provider calls**.

`SymbolRegistryRecord` remains a complete active trading-registry snapshot and must not be reused for partial identity writes because that path resets `active/selected` state outside the supplied snapshot. #92 therefore uses a dedicated instrument-identity persistence capability behind the existing PersistenceManager/backend owners. Native macOS SQLite and Windows system SQLite share schema version 5 identity storage; PostgreSQL implements the same logical contract. The unsupported file fallback reuses its existing persistence container rather than creating a second runtime store.

Persisted company/instrument identity is slow-changing evidence. It is warm-reused before provider refresh where available and must never be represented as live quote/current-market truth. Existing U.S.-equity eligibility remains authoritative; GLD/SLV/USO rules are unchanged. TradeInsight `symbol-search` remains hard-gated/non-executable and is not used by #92.

## Permanent boundaries
- U.S. equities processing only.
- No Execution/order routing.
- Smart Provider Router v2 is the sole general routing/admission authority.
- Direct SEC/EDGAR remains authoritative for Form 4.
- GLD, SLV and USO remain governed actionable exceptions.
- Reuse canonical freshness/cache/persistence/telemetry/state/validation owners.
- No secrets in diagnostics, tests, governance or handoff artifacts.
- G0–G16 and canonical Fast/Qualified/Release workflows only.
- No automatic lifecycle promotion and no parallel provider-specific lifecycle/health/router/cache system.
- No parallel company-profile database/router/cache/service.

## Exactly one next action
Complete #92 on `adapt-company-instrument-identity-001`: finish executable identity/persistence regressions and governance traceability, run canonical Fast on the exact candidate, reconcile any real gate finding without weakening assurance, then run impact-selected Qualified on that identical head. Merge only with an expected-head guard against unchanged live `main`, record immutable run/job/merge evidence on #92 and #65, then close #92 as completed. Do not create a Stable/public SemVer release merely for this adaptive child.

## Resume rule
1. Fetch live `main` and live `adapt-company-instrument-identity-001` first; another session may have advanced either.
2. Read `AGENTS.md`, `governance/AI-ASSISTANT-PORTABILITY-CONTRACT.md`, `governance/CI-EFFICIENCY-CONTRACT.md`, `governance/README.md`, this file, `governance/current-state.json`, issue #65 latest comments, issue #92 and its comments, and the #92 work-slice/closure ledger.
3. Inspect commits since `733d90ca125a4fe5abd38a2ea40de0623703dfd4` before changing code so #92 work is never duplicated.
4. Continue actual #92 implementation/qualification from the exact branch head; do not restart #73, #79, or #84.
5. Preserve Router v2, direct SEC/EDGAR, canonical symbol/persistence/Data Health owners, No Execution and U.S.-equities boundaries.
6. Use only canonical Fast, Qualified and Release workflows. Never weaken a gate or substitute documentation for executable/native evidence.
7. #92 cannot close until exact-head Fast + impact-selected Qualified + expected-head merge evidence all pass.
