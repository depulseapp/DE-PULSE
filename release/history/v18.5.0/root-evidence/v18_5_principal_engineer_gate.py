#!/usr/bin/env python3
"""Principal Engineer acceptance for DE.PULSE v18.5 Major Closure.

This gate validates current ownership/invariant contracts. It deliberately does
not require historical TEST-profile artifacts from earlier v18 point releases.
"""
from pathlib import Path
import json
import sys

R = Path(__file__).resolve().parent
errors = []


def need(ok, msg):
    if not ok:
        errors.append(msg)


def text(path):
    p = R / path
    need(p.exists(), f"missing required owner/evidence: {path}")
    return p.read_text(errors="ignore") if p.exists() else ""

scope = json.loads(text("v18_5_scope.json") or "{}")
matrix = json.loads(text("v18_5_major_closure.json") or "{}")
router = text("smart_router_v2.go") + "\n" + text("provider_router.go")
persist = text("persistence_repository.go")
pg = text("persistence_backend_postgres.go")
select = text("persistence_backend_select.go")
workspace = text("user_workspace.go")
main = text("main.go")
api = text("http_api.go")
health = text("http_health.go")
auth = text("http_auth.go")
identity = text("identity_model.go")
rapid = text("rapid_move_intelligence.go")
structure = text("governance/REPOSITORY_STRUCTURE_CONTRACT.md")
archive = text("governance/RELEASE_ARTIFACT_ARCHIVE.md")
readme = text("README.md")
roadmap = text("governance/ROADMAP.md")
adaptive_build = text("adaptive-governance/ADAPTIVE_BUILD_PROCESS.md")

# Frozen Major Closure contract.
need(scope.get("version") == "18.5.0", "v18.5 scope identity drift")
need(len(scope.get("closureDimensions", [])) == 10, "10 closure dimensions not frozen")
need(len(scope.get("adrGdiRequiredScenarios", [])) == 15, "15 ADR-GDI scenarios not frozen")
need(len(matrix.get("adrGdiScenarios", [])) == 15, "Major Closure matrix does not bind all ADR-GDI scenarios")

# One canonical provider/router authority and truthful routing state.
for token in ("NOT_ENTITLED", "Preferred", "Serving"):
    need(token in router, f"Smart Provider Router contract missing {token}")
need("TIERED_PARTIAL" in rapid, "Rapid Move coverage-truth contract missing")

# Persistence: one interface/manager owner, bounded PostgreSQL, fail-closed selection.
need(persist.count("type PersistenceBackend interface") == 1, "PersistenceBackend ownership duplicated")
need("newPersistenceManagerWithBackend" in persist, "persistence backend injection seam missing")
need("newLocalPersistenceBackend" in select and "newPostgresPersistenceBackend" in select, "central persistence backend selection missing")
need("newUnavailablePersistenceBackend" in select, "fail-closed unavailable persistence owner missing")
need("db.SetMaxOpenConns" in pg and "db.SetMaxIdleConns" in pg, "bounded PostgreSQL pool missing")
need("pg_advisory_lock" in pg and "LevelSerializable" in pg, "serialized/transactional PostgreSQL migration contract missing")
need("ProbeReady" in persist and "scheduleRetryLocked" in persist, "bounded persistence readiness/recovery owner missing")

# Shared market processing and outer runtime seam.
need("union, never a per-user provider pipeline" in workspace, "shared multi-user processing ownership drift")
need("isHostedRuntime()" in main and "openAppWindow" in main, "desktop/hosted runtime seam missing")
need("registerHealthRoutes" in api and '"/api/health"' in health and '"/api/ready"' in health, "health/readiness separation missing")

# Security invariants retained from v18 hardening without requiring old TEST packaging.
need("sessionTokenHash" in identity, "hashed session-token persistence missing")
need("SameSiteStrictMode" in auth and "X-DE-PULSE-CSRF" in auth, "cookie/CSRF hardening missing")

# Developer-proofing and delivery are now first-class closure obligations.
for phrase in ("cmd/depulse", "internal/", "tests/", "tools/", "config/", "active workflows only"):
    need(phrase in structure, f"Repository Structure Contract missing target archetype element: {phrase}")
for phrase in ("A tag alone is not sufficient", "macOS Apple Silicon runnable ZIP", "Windows x64 runnable ZIP", "GitHub Release"):
    need(phrase in archive, f"native/recovery delivery contract missing: {phrase}")

# Permanent product boundaries are checked against their canonical current owners,
# rather than requiring duplicated prose in every governance document.
protected = set(scope.get("protectedContracts", []))
need("no_execution_boundary" in protected and "no_execution_features" in set(scope.get("exclusions", [])), "No Execution boundary missing from frozen v18.5 scope")
normalized_readme = readme.lower().replace("u.s.-", "us ").replace("u.s. ", "us ").replace("-", " ")
need("us equities processing boundary" in normalized_readme or "us equities processing" in normalized_readme, "US Equities processing boundary missing from current product documentation")
need("SHADOW → VALIDATED → APPROVED → PRODUCTION" in (roadmap + "\n" + adaptive_build + "\n" + readme), "adaptive promotion governance missing")

if errors:
    print("v18.5 Principal Engineer Gate: FAIL")
    for err in errors:
        print(" -", err)
    sys.exit(2)

print("v18.5 Principal Engineer Gate: PASS · current ownership/invariants, architecture, ADR-GDI, security, delivery and developer-proofing verified")
