# CURRENT Adaptive Roadmap

**Canonical machine state:** `governance/current-state.json`  
**Certified Stable:** `v18.9.1-stable`  
**Completed provider/data-health program:** #79 with #80, #81, #82, #83, #78 and #84 / PR #91 / merge `733d90ca125a4fe5abd38a2ea40de0623703dfd4`  
**Active product work:** #92 / `ADAPT-COMPANY-INSTRUMENT-IDENTITY-001` / `adapt-company-instrument-identity-001`  
**Parent program:** #65 / `ADAPT-PROVIDER-INTELLIGENCE-010`

The current v18.9.x direction advances from completed provider production/data-health closure into canonical slow-changing company/instrument identity. #92 must strengthen the existing symbol/universe and persistence owners rather than create a company-profile subsystem.

The canonical identity path is persistence-first and reuse-first: existing in-memory identity -> persisted canonical identity -> validate the slow-changing record -> acquire only missing/expired evidence. When acquisition is needed, identity is captured from the **same Smart Provider Router v2 Alpaca `/v2/assets` U.S. Asset Universe response** already used for canonical universe discovery. No additional company-profile fetch is introduced.

Identity includes only useful canonical fields already present in that response: normalized symbol, display/company name, exchange, asset class, stable provider asset ID when supplied, source and observation time. The observation time is slow-changing identity provenance, not live quote/market evidence time. Existing universe eligibility remains authoritative.

`SymbolRegistryRecord` remains the complete trading-registry snapshot owner. Partial identity must never be written through it because full registry persistence intentionally resets `active/selected` state outside the supplied snapshot. #92 therefore uses dedicated instrument-identity storage behind the existing `PersistenceManager`/backend abstraction, with native macOS SQLite, Windows system SQLite and PostgreSQL logical parity.

## Retained Adaptive Data Health contract
The completed #79/#84 provider program remains an active invariant, not superseded prose. Provider/capability coverage continues to use truthful states including **PARTIAL COVERAGE** and **DATA DEGRADED** whenever current evidence is incomplete, stale, unavailable, suppressed, or fallback-constrained. Healthy-looking UI must never be synthesized from missing evidence. Smart Provider Router v2, canonical freshness, cache/persistence, telemetry, lifecycle/readiness, fault-recovery and zero-gap closure evidence remain the production data-health foundation beneath #92 and every later #65 child.

TradeInsight `symbol-search` remains hard-gated/non-executable. Smart Provider Router v2 remains the sole general routing authority. U.S. equities processing, direct SEC/EDGAR Form 4 authority, GLD/SLV/USO actionable exceptions, canonical freshness/cache/persistence/telemetry/state ownership and No Execution remain permanent boundaries.

After #92 closes through exact-head Fast -> impact-selected Qualified -> expected-head merge, #65 continues in dependency order; no Stable/public SemVer release is required merely for this capability child.
