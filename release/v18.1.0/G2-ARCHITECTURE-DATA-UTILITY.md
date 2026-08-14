# v18.1.0 G2 — Architecture / Data Utility

Build order: **REUSE → CONSOLIDATE → REFACTOR → DELETE/REPLACE → ADD**.

`UserWorkspace` is the only new durable personal-state owner. Existing IdentityService, PersistenceBackend, Hub, Tracked Symbols helpers, Provider Router v2 and canonical Engine stores are reused. Provider/evidence/scoring work stays global and processes one deduplicated union of workspace symbols. Per-user runtime/event filtering is a presentation/privacy boundary, not a second computation pipeline.

`data_utility_registry.json` records User Workspaces as active durable state with explicit consumers/retention. No raw market data is copied per user.
