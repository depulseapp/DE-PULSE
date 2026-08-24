'use strict';
/* DE.PULSE v18.8.1 release identity overlay.
   Adaptive Trust & Data Truth Closure changes release/process/data-truth and
   focused Research/Symbol/Readiness/Freshness contracts, not the large legacy
   renderer. This last-loaded overlay owns only current release/build integrity
   presentation and the v18.8.1 QA history entry. */
var DEPULSE_RELEASE_VERSION = '18.8.1';
var DEPULSE_RELEASE_BUILD_ID = 'v18.8.1-stable-20260820';
var DEPULSE_RELEASE_QA_ENTRY = Object.freeze({
  version: DEPULSE_RELEASE_VERSION,
  date: '2026-08-20',
  status: 'STABLE',
  summary: 'Adaptive Trust & Data Truth Closure: zero-miss reconciliation, evidence-time/freshness truth, Research and Symbol/Desk correctness, prep/readiness semantics, cost-aware CI telemetry and release-neutral Market Header ownership.',
  file: 'v18.8.1.txt',
  buildId: DEPULSE_RELEASE_BUILD_ID,
  checkpoint: 'release/v18.8.1/release_contract.json'
});

buildIdentityIntegrityMarkup = function buildIdentityIntegrityMarkupV1881(){
  const v=String(state?.version||''), b=String(state?.buildId||'');
  const ok=v===DEPULSE_RELEASE_VERSION && b===DEPULSE_RELEASE_BUILD_ID;
  return `<div class="build-integrity ${ok?'ok':'bad'}"><b>${ok?'BUILD IDENTITY VERIFIED':'BUILD IDENTITY MISMATCH'}</b><span>${esc(v||'unknown')} · ${esc(b||'unknown')}</span>${ok?'':'<small>Running backend does not match the packaged v18.8.1 renderer identity. Restart De-Pulse.app.</small>'}</div>`;
};

if (typeof qaHistoryMarkup === 'function') {
  const inheritedQAHistoryMarkupV1881 = qaHistoryMarkup;
  qaHistoryMarkup = function qaHistoryMarkupV1881(){
    if (Array.isArray(qaManifest) && !qaManifest.some((x)=>String(x?.version||'')===DEPULSE_RELEASE_VERSION)) {
      qaManifest = [DEPULSE_RELEASE_QA_ENTRY, ...qaManifest];
    }
    return inheritedQAHistoryMarkupV1881();
  };
}
