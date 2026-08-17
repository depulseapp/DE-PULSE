# DE.PULSE v18.5.0 — G2 Architecture / Data Utility

Status: PASS — closure contract frozen; implementation/evidence remains subject to later gates.

v18.5 preserves singular canonical ownership. Identity/session truth remains in the canonical identity service; provider selection and operational entitlement remain in Smart Provider Router v2; commercial/data-rights readiness remains explicit and structurally separate from executable routing; symbol intelligence remains shared/canonical; persistence remains SQLite for desktop and PostgreSQL for hosted shared state. No duplicate router, scanner, persistence stack, quote owner, recommendation silo or per-user market-wide computation may be introduced for closure.

Data Utility / Correlation remains mandatory: fetch/store/compute/display only with a named purpose and consumer; prefer fetch-once/calculate-once, material-change propagation and point-in-time lineage; suppress stale/redundant/low-value clutter; reuse canonical evidence across Research, Opportunity Radar, Decision Queue, readiness, preparation and validation consumers.

ADR-GDI architecture ownership is capability/dependency aware. A provider, database, queue, background job or dataset failure must map to the smallest truthful blast radius. Broad `DATA DEGRADED` is not a substitute for capability health. Work queues and provider/DB pressure must be bounded; decision-critical freshness receives priority; lower-value work sheds or defers before live/current evidence becomes misleading.

G2 release blocker: architecture that can create unbounded fan-out, duplicate provider/calculation work, broad unexplained degradation, or stale/misstated decision-critical evidence under realistic supported load must be corrected or bounded with explicit truthful operating limits before promotion.
