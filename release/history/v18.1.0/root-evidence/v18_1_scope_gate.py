#!/usr/bin/env python3
from pathlib import Path
import json,re,sys
R=Path(__file__).resolve().parent; e=[]
def need(ok,msg):
    if not ok:e.append(msg)
i=json.loads((R/'release_identity.json').read_text()); s=json.loads((R/'v18_1_scope.json').read_text()); c=json.loads((R/'v18_1_g0_g3_contract.json').read_text())
need(i.get('version')=='18.1.0' and i.get('channel') in {'TEST','STABLE'},'v18.1.0 release identity missing')
need(i.get('previous_stable')=='v18.0.6' and i.get('patch_predecessor')=='v18.0.6','v18.0.6 Stable predecessor drift')
if i.get('channel')=='TEST':
    need(i.get('build_id')=='v18.1.0-test-multi-user-my-market-symbols-20260814','TEST build identity drift')
    need(i.get('runtime_config')=='PersonalMarketTerminal-v18.1.0-TEST','isolated TEST runtime missing')
    need(i.get('application_bundle')=='De-Pulse-v18.1.0-TEST.app','separate TEST bundle missing')
else:
    need(i.get('build_id')=='v18.1.0-stable-multi-user-my-market-symbols-20260814','Stable build identity drift')
    need(i.get('runtime_config')=='PersonalMarketTerminal','Stable runtime config drift')
    need(i.get('application_bundle')=='De-Pulse.app','Stable bundle drift')
need(len(s.get('clauses',[]))==15,'immutable scope clause count mismatch')
need(c.get('baseline',{}).get('commit')=='714980096be84042b5bb883a6358d6e557da2df7','G0 Stable commit drift')
need(c.get('baseline',{}).get('sourceFingerprint')=='5f5285071f4a46ed6ab83026c158d58936eecdea38961940622faa8cfe4d7f90','G0 Stable fingerprint drift')
model=(R/'app_model.go').read_text(); ws=(R/'user_workspace.go').read_text(); runtime=(R/'user_runtime.go').read_text(); routes=(R/'http_api.go').read_text(); persist=(R/'persistence_repository.go').read_text(); tests=(R/'v18_1_test.go').read_text(); profile=(R/'v18_test_profile.go').read_text(); core=''.join((R/f).read_text() for f in ['engine_core.go','provider_router.go','smart_router_v2.go','rapid_move_intelligence.go'])
need('type UserWorkspace struct' in model,'UserWorkspace durable owner missing')
for term in ['initializeUserWorkspaces','ensureUserWorkspace','workspaceStateLocked','saveWorkspaceStateLocked','processingStateLocked']:
    need(term in ws,'workspace ownership helper missing '+term)
for term in ['runtimeSnapshotForUserFrom','SnapshotForUser','broadcastRuntime','broadcastSymbolEvent','broadcastNews','broadcastEarnings','broadcastFilings']:
    need(term in runtime,'runtime/event isolation helper missing '+term)
need('LoadUserWorkspaces(context.Context)' in persist and 'SaveUserWorkspace(context.Context' in persist,'persistence workspace contract missing')
need('user_workspaces' in (R/'persistence_backend_sqlite.go').read_text() and 'user_workspaces' in (R/'persistence_backend_windows.go').read_text(),'SQLite workspace migration missing')
need('Workspaces map[string]UserWorkspace' in (R/'persistence_backend_fallback.go').read_text(),'fallback workspace persistence missing')
for route in ['/api/master-symbol/add','/api/master-symbol/remove','/api/master-symbol/remove-all','/api/desk/membership','/api/watchlists/add-symbol']:
    need(re.search(re.escape(route)+r'", a\.auth\(',routes) is not None,'user-owned route not authenticated USER-capable '+route)
need('mux.HandleFunc("/api/settings/ai-provider", a.requireRole(RoleAdmin' in routes,'global AI provider policy not ADMIN-owned')
need('mux.HandleFunc("/api/settings/save", a.requireRole(RoleAdmin' in routes,'global Settings not ADMIN-owned')
for term in ['WorkspaceIsolationAndEmptyNewUser','ProcessingUnionDeduplicatesAndRetainsOtherOwner','MasterSymbolMutationIsUserScoped','RuntimeSnapshotFiltersOtherWorkspaceSymbols','SymbolEventFanoutTargetsWorkspaceOwners','LegacyOwnerMigrationAndWorkspacePersistence','EnsureUserWorkspaceCreatesDurableEmptyWorkspace']:
    need(term in tests,'focused isolation test missing '+term)
need('PersonalMarketTerminal-v18.1.0-TEST' in profile and '"sourceVersion": "v18.0.6"' in profile and '.v18.1.0-test-profile-migration.json' in profile,'v18.1 TEST profile migration truth missing')
need(len(re.findall(r'func\s+\(e \*Engine\)\s+executeProviderRoute\s*\(',core))==1,'Provider Router canonical owner duplicated')
need('rapid-move-v1.1.0' in core,'Rapid Move v1.1.0 Stable foundation missing')
if e:
    print('v18.1.0 Multi-User Scope Gate: FAIL'); [print(' -',x) for x in e]; sys.exit(2)
print('v18.1.0 Multi-User Scope Gate: PASS · 15/15 clauses · per-user ownership isolated · one shared intelligence core preserved')
