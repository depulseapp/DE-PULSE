# CURRENT Adaptive Build Plan

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Completed provider/data-health program:** #79 with final #84 / PR #91 / merge `733d90ca125a4fe5abd38a2ea40de0623703dfd4`  
**Active product work:** #92 / `ADAPT-COMPANY-INSTRUMENT-IDENTITY-001` / `adapt-company-instrument-identity-001`  
**Parent program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`

#92 required build outputs are:
1. Extend the existing Alpaca asset-universe decoder only with useful identity fields already returned by `/v2/assets`; do not add a provider call.
2. Continue using `canonicalUSUniverseAssetEligible` and existing U.S.-equity/tradability/exchange/symbol boundaries before identity admission.
3. Store normalized `InstrumentIdentityRecord` values in the Engine canonical symbol/universe owner and expose only bounded canonical lookup to actual consumers.
4. Reuse `PersistenceManager` and existing persistence backends. Native SQLite/Windows and PostgreSQL use dedicated instrument-identity records with source and observation time; no second database/cache is permitted.
5. Never use partial `SymbolRegistryRecord` persistence for identity. Regression proof must show identity writes cannot deactivate or clear `active/selected` trading-registry state.
6. Warm-reuse persisted slow-changing identity on restart before provider refresh when valid; never present its observation time as current market/quote truth.
7. Keep TradeInsight `symbol-search` gated/non-executable; the same-response Alpaca path is canonical for #92.
8. Regression-protect same-request/zero-extra-fetch capture, normalization/filtering, persistence round-trip, restart reuse, monotonic stale-overwrite rejection and registry isolation.
9. Preserve Smart Provider Router v2, direct SEC authority, GLD/SLV/USO actionable exceptions, U.S. equities and No Execution.
10. Maintain Roadmap, Build Plan, Build Process, Delivery Process, current-state, work-slice and handoff alignment.

Delivery remains one branch/one PR, canonical Fast on the exact final candidate followed by impact-selected Qualified on the identical head, then expected-head merge only if live `main` is unchanged. No Stable/public SemVer release is created solely for #92.
