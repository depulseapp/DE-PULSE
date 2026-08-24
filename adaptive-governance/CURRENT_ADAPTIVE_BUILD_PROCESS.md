# CURRENT Adaptive Build Process

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Completed provider/data-health program:** #79 with final #84 / PR #91 / merge `733d90ca125a4fe5abd38a2ea40de0623703dfd4`  
**Active product work:** #92 / `ADAPT-COMPANY-INSTRUMENT-IDENTITY-001` / `adapt-company-instrument-identity-001`  
**Parent program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`

For #92, implementation is executable-first, reuse-first and canonical-owner-first:
- reuse Smart Provider Router v2 and the existing `US Asset Universe` acquisition; a company/instrument identity implementation that adds a second provider request or bypasses Router v2 fails architecture review;
- reuse canonical symbol eligibility and persistence owners before adding data structures;
- capture identity only for assets that pass the existing canonical U.S. eligibility boundary;
- treat identity as slow-changing provenance-bound evidence, not live/current market truth;
- load valid persisted identity before refetching equivalent data when possible;
- never route partial identity through `PersistenceBackend.UpsertSymbols` / `SymbolRegistryRecord`, because that API owns a complete registry snapshot and resets `active/selected`; use dedicated identity persistence behind the same backend instead;
- require logical persistence parity across supported native macOS SQLite and Windows system SQLite; maintain PostgreSQL compatibility for later hosted use without changing caller ownership;
- keep TradeInsight `symbol-search` hard-gated/non-executable; do not use beta reachability as permission to change lifecycle or authority;
- fail identity persistence independently: a storage error may degrade identity persistence diagnostics but must not rewrite a successful provider universe request as provider failure or corrupt registry state;
- regression-test zero-extra-fetch, filter/normalization, restart reuse, registry isolation and stale-overwrite rejection before CI;
- use canonical Fast exact-head PASS, then deliberate impact-selected Qualified exact-head PASS; fix real findings without weakening gates;
- do not add workflow families, parallel caches/databases/services or weaken G0-G16/source-health/architecture gates.

Documentation alone never verifies #92. Actual GitHub run/job/merge evidence remains the post-commit authority.
