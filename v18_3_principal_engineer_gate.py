#!/usr/bin/env python3
from pathlib import Path
import sys
R=Path(__file__).resolve().parent; e=[]
def need(ok,msg):
    if not ok:e.append(msg)
repo=(R/'persistence_repository.go').read_text(); select=(R/'persistence_backend_select.go').read_text(); pg=(R/'persistence_backend_postgres.go').read_text(); workspace=(R/'user_workspace.go').read_text(); main=(R/'main.go').read_text(); api=(R/'http_api.go').read_text(); scope=(R/'release/v18.3.0/G1-IMMUTABLE-SCOPE.md').read_text()
need(repo.count('type PersistenceBackend interface')==1,'PersistenceBackend ownership duplicated')
need('newPersistenceManagerWithBackend' in repo,'backend injection/test seam missing')
need('newLocalPersistenceBackend' in select and 'newPostgresPersistenceBackend' in select,'central backend selection missing')
need('newUnavailablePersistenceBackend' in select,'fail-closed backend owner missing')
need('db.SetMaxOpenConns' in pg and 'db.SetMaxIdleConns' in pg,'bounded PostgreSQL pool missing')
need('db.Conn(ctx)' in pg and 'pg_advisory_lock' in pg and 'pg_advisory_unlock' in pg,'same-session migration serialization missing')
need('LevelSerializable' in pg,'transactional migration isolation missing')
need('DataState = "persisted"' in pg and 'FeedType = "persisted"' in pg,'warm-start truth relabel missing')
need('union, never a per-user provider pipeline' in workspace,'shared processing owner drift')
need('isHostedRuntime()' in main and 'openAppWindow' in main,'desktop/hosted outer runtime seam missing')
need('"/api/health"' in api and '"/api/ready"' in api,'liveness/readiness separation missing')
need('No Execution Boundary' in scope and 'Day/Swing/Long formulas' in scope,'protected product boundaries missing')
if e:
    print('v18.3 Principal Engineer Gate: FAIL'); [print(' -',x) for x in e]; sys.exit(2)
print('v18.3 Principal Engineer Gate: PASS · one persistence owner · fail-closed hosted selection · bounded pool · same-session migration lock · shared market core preserved')
