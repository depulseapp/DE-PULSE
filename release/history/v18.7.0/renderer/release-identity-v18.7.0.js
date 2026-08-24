'use strict';
/* DE.PULSE v18.7.0 release identity overlay.
   Runtime Reliability & Data Truth does not require rewriting the large legacy
   renderer merely to change build identity. This last-loaded overlay owns only
   release/build integrity presentation and the v18.7 QA history entry. */
var DEPULSE_RELEASE_VERSION = '18.7.0';
var DEPULSE_RELEASE_BUILD_ID = 'v18.7.0-stable-20260819';
var DEPULSE_RELEASE_QA_ENTRY = Object.freeze({
  version: DEPULSE_RELEASE_VERSION,
  date: '2026-08-19',
  status: 'STABLE',
  summary: 'Runtime Reliability & Data Truth: fail-closed degradation semantics, recovery hysteresis, bounded-load truth, active-market reliability evidence, and reproducible G11-G16 release dependencies.',
  file: 'v18.7.0.txt',
  buildId: DEPULSE_RELEASE_BUILD_ID,
  checkpoint: 'release/v18.7.0/release_contract.json'
});

buildIdentityIntegrityMarkup = function buildIdentityIntegrityMarkupV1870(){
  const v=String(state?.version||''), b=String(state?.buildId||'');
  const ok=v===DEPULSE_RELEASE_VERSION && b===DEPULSE_RELEASE_BUILD_ID;
  return `<div class="build-integrity ${ok?'ok':'bad'}"><b>${ok?'BUILD IDENTITY VERIFIED':'BUILD IDENTITY MISMATCH'}</b><span>${esc(v||'unknown')} · ${esc(b||'unknown')}</span>${ok?'':'<small>Running backend does not match the packaged v18.7.0 Stable renderer identity. Restart De-Pulse.app.</small>'}</div>`;
};

if (typeof qaHistoryMarkup === 'function') {
  const inheritedQAHistoryMarkupV1870 = qaHistoryMarkup;
  qaHistoryMarkup = function qaHistoryMarkupV1870(){
    if (Array.isArray(qaManifest) && !qaManifest.some((x)=>String(x?.version||'')===DEPULSE_RELEASE_VERSION)) {
      qaManifest = [DEPULSE_RELEASE_QA_ENTRY, ...qaManifest];
    }
    return inheritedQAHistoryMarkupV1870();
  };
}
