#!/usr/bin/env python3
from pathlib import Path
import json,re,sys
R=Path(__file__).resolve().parent; e=[]
def need(ok,msg):
    if not ok: e.append(msg)

ident=json.loads((R/'release_identity.json').read_text())
version=str(ident.get('version','')); channel=str(ident.get('channel','')); stable=str(ident.get('stable_baseline','')); previous=str(ident.get('previous_stable',''))
canon={
 'user.md':'# DE.PULSE — User documentation',
 'developer.md':'# DE.PULSE — Developer documentation',
 'limitations.md':'# DE.PULSE — Capabilities & Limitations',
}

# Permanent historical/product truth must survive new release families.
req={
 'user.md':['v17.5.1 STABLE','30 FULL / 0 PARTIAL / 0 MISSING','Trade Readiness','STALE','INCOMPLETE','Professional Trader','No Execution'],
 'developer.md':['v17.5.1 STABLE','Provider Router','release_identity.json','source_health_baseline.json','No v17 database','No Execution'],
 'limitations.md':['v17.5.1 STABLE','decision support, not a profit predictor','U.S.-listed only','stale/cached/history-only','physical native macOS/Windows','No Execution'],
}

# Documentation is intentionally two-tiered:
# - long-lived promoted/manual history stays in renderer/docs + renderer/qa/manifest.json;
# - TEST/RC patch candidates use a release-specific Documentation Impact + QA delta.
# This preserves full documentation quality while avoiding 38-80KB manual rewrites after
# every small patch. Stable promotion/material doc boundaries fold validated deltas back in.
impact_path=R/'release'/f'v{version}'/'DOCUMENTATION-IMPACT.md'
qa_delta_path=R/'release'/f'v{version}'/'QA-MANIFEST-DELTA.json'
patch_delta_mode=channel in {'TEST','RC'} and impact_path.exists() and qa_delta_path.exists()
impact=impact_path.read_text(errors='ignore') if impact_path.exists() else ''

for f,terms in req.items():
    s=(R/'renderer/docs'/f).read_text(errors='ignore')
    need(s.startswith(canon[f]+'\n'),f+' canonical first heading drift')
    need(len(re.findall(r'(?m)^#\s+',s))==1,f+' must contain exactly one H1')
    if patch_delta_mode:
        need(f'v{version} {channel}' in '\n'.join(impact.splitlines()[:35]),f+' current patch candidate missing from Documentation Impact')
        section={'user.md':'## User documentation impact','developer.md':'## Developer documentation impact','limitations.md':'## Capabilities & limitations impact'}[f]
        need(section in impact,f+' current patch impact section missing')
    else:
        need(f'v{version} {channel}' in '\n'.join(s.splitlines()[:35]),f+' current candidate section missing')
    for t in terms: need(t in s,f'{f} missing {t}')

if patch_delta_mode:
    for term in [
        f'v{version} {channel}',
        f'Immediate Stable predecessor:** {previous} STABLE',
        'Major v18 provenance anchor:** v17.5.1 STABLE',
        'No Execution',
        'macOS Apple Silicon',
        'Windows x64',
        'G0-G10',
        'G11-G15',
    ]:
        need(term in impact,'Documentation Impact missing '+term)

# Long-lived inherited QA assets remain available as regression/history evidence.
qa_required=['v16.11.0-master-scope.json','v16.11.0-acceptance-evidence.json','v16.11.0-traceability.md','v16.11.0-verification-plan.md','v16.11.0-g0-baseline-audit.md','v16.11.0-g1-major-closure-contract.md','v16.11.0-g2-source-architecture-audit.md','v16.11.0-adaptive-test-selection.json','v16.11.0-performance-baseline.md','v16.11.0-senior-engineer-review.md','v16.11.0-trader-investor-review.md','v16.11.0-adaptive-retrospective.md','v16.11.0-data-utility-source-hygiene-audit.md','v16.11.0-major-closure-scope-audit.md','v16.11.0.txt','v16.9.0-original-roadmap-status.json','original-professional-roadmap-acceptance.json','v17.5.1.txt','v18.0.0.txt','v18.0.0-identity-session-foundation.md']
for q in qa_required: need((R/'renderer/qa'/q).exists(),'missing '+q)

readme=(R/'README.md').read_text(errors='ignore')
# README remains the promoted/current Stable landing page while a patch is TEST/RC.
if channel=='STABLE':
    readme_terms=[f'DE.PULSE v{version}',f'Current Stable baseline:** {previous}']
else:
    prev=previous.lstrip('v')
    readme_terms=[f'DE.PULSE v{prev} STABLE','Major v18 provenance anchor:** v17.5.1']
for term in readme_terms+['No Execution Boundary','Adaptive Build Process v2','Data Utility','v18','v19','v20']:
    need(term in readme,'README missing '+term)

try:
    man=json.loads((R/'renderer/qa/manifest.json').read_text()); releases=man.get('releases',[]); first=releases[0]
    if channel=='STABLE':
        need(str(first.get('version'))==version and str(first.get('status'))==channel,'QA manifest first entry does not match canonical current release identity')
    elif patch_delta_mode:
        delta=json.loads(qa_delta_path.read_text())
        need(str(delta.get('version'))==version and str(delta.get('status'))==channel,'QA patch delta does not match canonical current release identity')
        need(str(delta.get('previousStable','')).lstrip('v')==previous.lstrip('v'),'QA patch delta previous Stable mismatch')
        need(str(delta.get('majorProvenanceAnchor','')).lstrip('v')==stable.lstrip('v'),'QA patch delta provenance anchor mismatch')
    else:
        need(False,'TEST/RC candidate requires current QA manifest entry or release-specific QA-MANIFEST-DELTA.json')

    stable_entries=[x for x in releases if x.get('status')=='STABLE']
    expected_latest_stable=version if channel=='STABLE' else previous
    need(bool(stable_entries) and str(stable_entries[0].get('version')).lstrip('v')==str(expected_latest_stable).lstrip('v'),'latest authoritative Stable entry does not match release identity baseline')

    reg=json.loads((R/'release_learning_registry.json').read_text())
    for rid in [f'RL-{i:03d}' for i in range(20,32)]:
        need(any(x.get('id')==rid and x.get('status')=='ACTIVE' for x in reg.get('entries',[])),rid+' not ACTIVE')
    status=json.loads((R/'renderer/qa/v16.9.0-original-roadmap-status.json').read_text())
    counts={k:sum(1 for x in status['items'] if x['status']==k) for k in ['FULL','PARTIAL','MISSING']}
    need(counts=={'FULL':30,'PARTIAL':0,'MISSING':0},'inherited original roadmap status not 30 FULL / 0 PARTIAL / 0 MISSING')
    sel=json.loads((R/'renderer/qa/v16.11.0-adaptive-test-selection.json').read_text()); need(len(sel.get('mandatory',[]))>=25,'Major Closure Adaptive Test Selection incomplete')
    scope=json.loads((R/'renderer/qa/v16.11.0-master-scope.json').read_text()); need(len(scope.get('scope_lock',[]))==10,'v16.11 Major Closure scope not 10 clauses')
    acc=json.loads((R/'renderer/qa/v16.11.0-acceptance-evidence.json').read_text()); need(len(acc.get('clauses',[]))==10 and acc.get('status')=='10/10 PASS','v16.11 acceptance evidence not 10/10 PASS')
    ci=json.loads((R/'ci_pipeline_plan.json').read_text()); need(version==str(ci.get('version')),'canonical release/CI identity drift')

    if version.startswith('18.'):
        need(channel in {'TEST','RC','STABLE'},'v18 release channel is not recognized')
        need(stable=='v17.5.1','v18 must preserve v17.5.1 as the major-family provenance anchor')
        need(previous=='v17.5.1' or (previous.startswith('v18.') and previous!=f'v{version}'),'v18 immediate Stable predecessor must be v17.5.1 or an earlier certified v18 Stable')
        for doc in canon:
            txt=(R/'renderer/docs'/doc).read_text(errors='ignore'); head='\n'.join(txt.splitlines()[:90])
            need('v18.0.0 TEST' in head,doc+' missing v18.0 TEST section')
            need('v17.5.1 STABLE' in head,doc+' missing authoritative v17.5.1 history')
        need('PersonalMarketTerminal-v18-TEST' in readme,'README missing isolated v18 TEST runtime truth')
    elif version.startswith('17.'):
        need(stable=='v16.11.0' and ident.get('previous_stable')=='v16.11.0','v17 incoming Stable baseline drift')
    elif version=='16.11.0':
        need(channel=='STABLE' and ident.get('previous_stable')=='v16.10.0','v16.11 Stable predecessor identity drift')
    else:
        need(False,'documentation governance does not recognize current major identity')
except Exception as ex:
    e.append('manifest/learning/status parse failed: '+str(ex))

if e:
    print('G16 Documentation Governance: FAIL')
    for x in e: print(' -',x)
    sys.exit(1)
mode='patch-delta' if patch_delta_mode else 'promoted/manual'
print(f'G16 Documentation Governance: PASS · v{version} {channel} current · immediate Stable {previous} · v18 provenance anchor {stable} preserved · one-H1 hierarchy · inherited roadmap/learning evidence retained · documentation mode {mode}')
