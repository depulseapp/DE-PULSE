#!/usr/bin/env python3
from pathlib import Path
import sys,re
ROOT=Path(__file__).resolve().parent
errors=[]
def need(cond,msg):
    if not cond: errors.append(msg)
idx=(ROOT/'renderer/index.html').read_text()
js=(ROOT/'renderer/renderer.js').read_text()
css=(ROOT/'renderer/styles.css').read_text()
docs='\n'.join((ROOT/'renderer/docs'/x).read_text() for x in ['user.md','developer.md','limitations.md'])
active='\n'.join([idx,js])
# Canonical visible terminology required by the current product.
for label in ['Pre-Market Prep','Market Open Prep','Earnings & Material Catalyst Watch','Provider Capability Registry','Live Market Coverage','Global / FX Refresh','Data Engine Detail']:
    need(label in active or label in docs,f'canonical label missing: {label}')
# Known inconsistent variants that have caused UI drift must not return in active UI source.
for bad in ['Pre-market Prep','Market open Prep','Long term Desk','Websocket','FinnHub','Alpacca','Twelvedata']:
    need(bad not in active,f'inconsistent active UI terminology: {bad}')
# Side panel intentionally has no generic manual-actions console; Maintenance owns generic operations.
need('id="data-engine-manual-actions"' not in idx,'generic sidebar Manual Actions host returned')
need('maintenance-card-action' in js and 'data-engine-action' in js,'Maintenance contextual action controls missing')
# Action copy follows action-oriented Title Case / concise control language.
for label in ["'Run Now'","'Evaluate'","'Reconnect'","'Recheck'","'Refresh'","'Refresh Due'","'Retry'"]:
    need(label in js,f'expected action label missing: {label}')
# Defined preparation/action states remain canonical uppercase tokens.
for state in ['READY','RUNNING','COMPLETE','DEGRADED','FAILED']:
    need(state in js,f'defined state missing: {state}')
# Current header/status geometry/copy rules are explicitly represented.
need('fitHeaderNotification' in js,'content-aware header notification fitting missing')
need('short-message' in js and 'long-message' in js,'short/long header notification behavior missing')
need('.sidebar-prep-state' in css,'Data Engine status typography rule missing')
# Documentation/traceability must not describe the current rich sidebar as compact.
trace=(ROOT/'renderer/qa/v14-traceability.md').read_text()
need('Rich anchored sidebar Data Engine' in trace,'traceability does not reflect rich side Data Engine')
need('Compact anchored sidebar Data Engine' not in trace,'obsolete compact Data Engine requirement remains')
if errors:
    print('content & copy consistency audit: FAIL')
    print('\n'.join('- '+e for e in errors)); sys.exit(1)
print('content & copy consistency audit: PASS')
