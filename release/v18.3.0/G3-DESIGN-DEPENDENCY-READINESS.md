# G3 Design / Dependency Readiness — v18.3.0

- Hosted PostgreSQL uses `github.com/jackc/pgx/v5` through `database/sql` under the `postgres` build tag; desktop builds do not import or initialize the PostgreSQL driver.
- Selected pgx baseline: **v5.10.0**, MIT licensed. Its module requires Go 1.25; the hosted CI lane therefore uses the repository CI toolchain Go 1.26.6 while the default desktop source remains compatible with the repository `go 1.23` contract.
- PostgreSQL integration qualification uses an ephemeral PostgreSQL 17 service and an isolated test database. No paid database/service or production credentials are required for development qualification.
- Pool defaults are bounded (16 max open / 8 max idle) with explicit hard caps and configurable connection lifetime/idle limits.
- Migrations execute transactionally under a PostgreSQL advisory lock so concurrent hosted instances cannot race schema initialization.
- PostgreSQL selection/environment keys are runtime-only and must never be included in public diagnostics or logs beyond backend/readiness state.
- Focused tests cover default desktop persistence, fail-closed selection, pool bounds, hosted runtime/config resolution and readiness semantics.
- Real PostgreSQL tests cover repository parity, warm-start truth, identity/workspace persistence, concurrent workspace writes and serialized migration startup.
- Before G10, add and prove SQLite→PostgreSQL export/import, backup/restore, realistic contention/load, DB outage/recovery, bounded retry/backpressure, and hosted runtime restart scenarios.
- Required later release qualification remains G0–G16 only, including fresh G12 full certification and independent G13/G14 macOS Apple Silicon + Windows x64 runtime/package audits.
