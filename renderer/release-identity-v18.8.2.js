'use strict';
/* DE.PULSE v18.8.2 release identity overlay.
   Market Intelligence Reliability repairs canonical freshness/recovery
   accountability and unavailable-vs-zero truth without changing the large
   legacy renderer or protected trading semantics. */
var DEPULSE_RELEASE_VERSION = '18.8.2';
var DEPULSE_RELEASE_BUILD_ID = 'v18.8.2-stable-20260820';
var DEPULSE_RELEASE_QA_ENTRY = Object.freeze({
  version: DEPULSE_RELEASE_VERSION,
  date: '2026-08-20',
  status: 'STABLE',
  summary: 'Market Intelligence Reliability: canonical breadth quote freshness/recovery accountability, truthful SPY/QQQ/VIX and breadth evidence state, and unavailable/unknown distinct from observed zero.',
  file: 'v18.8.2.txt',
  buildId: DEPULSE_RELEASE_BUILD_ID,
  checkpoint: 'release/v18.8.2/release_contract.json'
});

buildIdentityIntegrityMarkup = function buildIdentityIntegrityMarkupV1882(){
  const v=String(state?.version||''), b=String(state?.buildId||'');
  const ok=v===DEPULSE_RELEASE_VERSION && b===DEPULSE_RELEASE_BUILD_ID;
  return `<div class="build-integrity ${ok?'ok':'bad'}"><b>${ok?'BUILD IDENTITY VERIFIED':'BUILD IDENTITY MISMATCH'}</b><span>${esc(v||'unknown')} · ${esc(b||'unknown')}</span>${ok?'':'<small>Running backend does not match the packaged v18.8.2 renderer identity. Restart De-Pulse.app.</small>'}</div>`;
};

if (typeof qaHistoryMarkup === 'function') {
  const inheritedQAHistoryMarkupV1882 = qaHistoryMarkup;
  qaHistoryMarkup = function qaHistoryMarkupV1882(){
    if (Array.isArray(qaManifest) && !qaManifest.some((x)=>String(x?.version||'')===DEPULSE_RELEASE_VERSION)) {
      qaManifest = [DEPULSE_RELEASE_QA_ENTRY, ...qaManifest];
    }
    return inheritedQAHistoryMarkupV1882();
  };
}
