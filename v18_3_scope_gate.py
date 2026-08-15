#!/usr/bin/env python3
from pathlib import Path
import json, subprocess, sys
R=Path(__file__).resolve().parent; errors=[]
def need(ok,msg):
    if not ok: errors.append(msg)
identity=json.loads((R/'release_identity.json').read_text())
scope=json.loads((R/'v18_3_scope.json').read_text())
channel=identity.get('channel'); need(identity.get('version')=='18.3.0' and channel in {'TEST','STABLE'},'v18.3 release identity missing')
need(identity.get('previous_stable')=='v18.2.0' and identity.get('patch_predecessor')=='v18.2.0','v18.2.0 predecessor drift')
expected_build='v18.3.0-test-postgresql-hosted-shared-state-20260814' if channel=='TEST' else 'v18.3.0-stable-postgresql-hosted-shared-state-20260815'
expected_runtime='PersonalMarketTerminal-v18.3.0-TEST' if channel=='TEST' else 'PersonalMarketTerminal'
expected_bundle='De-Pulse-v18.3.0-TEST.app' if channel=='TEST' else 'De-Pulse.app'
need(identity.get('build_id')==expected_build,'v18.3 build identity drift')
need(identity.get('runtime_config')==expected_runtime,'v18.3 runtime identity drift')
need(identity.get('application_bundle')==expected_bundle,'v18.3 application bundle drift')
need(scope.get('incomingStableCommit')=='78c42d739bff64b0e4c8676cb341fee65c4a3e67','G0 Stable commit drift')
need(scope.get('incomingStableFingerprint')=='d867dc2b7c58be879d1fa15a3e27a7f67dbf03509d52f7925963c1c74758c7ed','G0 Stable fingerprint drift')
need(len(scope.get('clauses',[]))==16,'scope clause count drift')
repo=(R/'persistence_repository.go').read_text(); archive=(R/'persistence_archive.go').read_text(); select=(R/'persistence_backend_select.go').read_text(); pg=(R/'persistence_backend_postgres.go').read_text(); stub=(R/'persistence_backend_postgres_stub.go').read_text(); workspace=(R/'user_workspace.go').read_text(); main=(R/'main.go').read_text(); api=(R/'http_api.go').read_text(); health=(R/'http_health.go').read_text()
for token in ['type PersistenceBackend interface','newPersistenceManagerWithBackend','PersistencePoolDiagnostics','PersistenceDatabaseDiagnostics','ProbeReady','HealthState','RetryBackoffMs']:
    need(token in repo, token+' missing from canonical persistence owner')
need('persistenceArchiveSchemaVersion' in archive and 'ExportArchiveFile' in archive and 'RestoreArchiveFile' in archive,'versioned backup/restore/migration archive missing')
need('sha256.Sum256' in archive and '0600' in archive,'archive integrity/private-file contract missing')
need('persistenceRestoreModeEmpty' in archive and 'persistenceRestoreModeReplace' in archive,'safe restore modes missing')
need('newLocalPersistenceBackend(configDir)' in select,'desktop local backend delegation missing')
need('case "postgres", "postgresql"' in select and 'newPostgresPersistenceBackend' in select,'explicit PostgreSQL selection missing')
need('//go:build postgres' in pg and 'github.com/jackc/pgx/v5/stdlib' in pg,'hosted PostgreSQL driver/build boundary missing')
need('//go:build !postgres' in stub and 'newUnavailablePersistenceBackend' in stub,'fail-closed non-postgres build behavior missing')
for token in ['pg_advisory_lock','schema_migrations','user_workspaces','identity_state','quote_history','PoolDiagnostics','DatabaseDiagnostics','ExportPersistenceArchive','RestorePersistenceArchive','HealthCheck']:
    need(token in pg, 'PostgreSQL parity/readiness token missing: '+token)
need('union, never a per-user provider pipeline' in workspace,'shared market processing contract drift')
need('isHostedRuntime()' in main and 'hostedListenAddress()' in main,'hosted runtime outer seam missing')
need('registerHealthRoutes' in api and '"/api/ready"' in health and 'StatusServiceUnavailable' in health,'hosted readiness contract missing')
for f in ('release/v18.3.0/G0-EXACT-BASELINE.md','release/v18.3.0/G1-IMMUTABLE-SCOPE.md','release/v18.3.0/G2-ARCHITECTURE-DATA-UTILITY.md','release/v18.3.0/G3-DESIGN-DEPENDENCY-READINESS.md','v18_3_g0_g3_contract.json'):
    need((R/f).exists(), f+' missing')
# When a real Git checkout/tag is available, protected deterministic owners must stay untouched.
protected=['scanner.go','preparation_types_liquidity.go','validation_learning.go','signal_validation.go']
if (R/'.git').exists():
    try:
        changed=subprocess.check_output(['git','diff','--name-only','v18.2.0-stable...HEAD','--',*protected],cwd=R,text=True).splitlines()
        need(not changed,'protected deterministic/formula owner drift: '+', '.join(changed))
    except Exception as exc:
        errors.append('could not verify protected-source diff: '+str(exc))
if errors:
    print('v18.3 Scope Gate: FAIL')
    for err in errors: print(' -',err)
    sys.exit(2)
print(f'v18.3 Scope Gate: PASS · {channel} · 16/16 clauses · PostgreSQL/hosted shared-state boundaries preserved')
