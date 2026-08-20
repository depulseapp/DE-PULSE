'use strict';
/* DE.PULSE v18.8.0 release identity overlay.
   Shared Intelligence Consolidation changes Scanner/Radar universe ownership,
   not the large legacy renderer. This last-loaded overlay owns only current
   release/build integrity presentation and the v18.8 QA history entry. */
var DEPULSE_RELEASE_VERSION = '18.8.0';
var DEPULSE_RELEASE_BUILD_ID = 'v18.8.0-stable-20260819';
var DEPULSE_RELEASE_QA_ENTRY = Object.freeze({
  version: DEPULSE_RELEASE_VERSION,
  date: '2026-08-19',
  status: 'STABLE',
  summary: 'Shared Intelligence Consolidation: one canonical U.S.-equity universe owner for Scanner and Opportunity Radar with coalesced refresh, truthful stale evidence, bounded retry and preserved ranking/promotion behavior.',
  file: 'v18.8.0.txt',
  buildId: DEPULSE_RELEASE_BUILD_ID,
  checkpoint: 'release/v18.8.0/release_contract.json'
});

buildIdentityIntegrityMarkup = function buildIdentityIntegrityMarkupV1880(){
  const v=String(state?.version||''), b=String(state?.buildId||'');
  const ok=v===DEPULSE_RELEASE_VERSION && b===DEPULSE_RELEASE_BUILD_ID;
  return `<div class="build-integrity ${ok?'ok':'bad'}"><b>${ok?'BUILD IDENTITY VERIFIED':'BUILD IDENTITY MISMATCH'}</b><span>${esc(v||'unknown')} · ${esc(b||'unknown')}</span>${ok?'':'<small>Running backend does not match the packaged v18.8.0 Stable renderer identity. Restart De-Pulse.app.</small>'}</div>`;
};

if (typeof qaHistoryMarkup === 'function') {
  const inheritedQAHistoryMarkupV1880 = qaHistoryMarkup;
  qaHistoryMarkup = function qaHistoryMarkupV1880(){
    if (Array.isArray(qaManifest) && !qaManifest.some((x)=>String(x?.version||'')===DEPULSE_RELEASE_VERSION)) {
      qaManifest = [DEPULSE_RELEASE_QA_ENTRY, ...qaManifest];
    }
    return inheritedQAHistoryMarkupV1880();
  };
}
