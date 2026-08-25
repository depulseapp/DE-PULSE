'use strict';
/* DE.PULSE v18.10.0 final v18 closure release identity overlay.
   10/10 future-proof closure; shared provider-routing, persistence,
   data-truth and No Execution semantics remain governed and unchanged. */
var DEPULSE_RELEASE_VERSION = '18.10.0';
var DEPULSE_RELEASE_BUILD_ID = 'v18.10.0-stable-20260825';
var DEPULSE_RELEASE_QA_ENTRY = Object.freeze({
  version: DEPULSE_RELEASE_VERSION,
  date: '2026-08-25',
  status: 'CANDIDATE',
  summary: '10/10 Future-Proof Final v18 Closure candidate: exhaustive assurance, packaged macOS/Windows lifecycle, provenance and no-rebuild publication readiness.',
  file: 'v18.10.0.txt',
  buildId: DEPULSE_RELEASE_BUILD_ID,
  checkpoint: 'release/v18.10.0/release_contract.json'
});

buildIdentityIntegrityMarkup = function buildIdentityIntegrityMarkupV18100(){
  const v=String(state?.version||''), b=String(state?.buildId||'');
  const ok=v===DEPULSE_RELEASE_VERSION && b===DEPULSE_RELEASE_BUILD_ID;
  return `<div class="build-integrity ${ok?'ok':'bad'}"><b>${ok?'BUILD IDENTITY VERIFIED':'BUILD IDENTITY MISMATCH'}</b><span>${esc(v||'unknown')} · ${esc(b||'unknown')}</span>${ok?'':'<small>Running backend does not match the packaged v18.10.0 renderer identity. Restart DE.PULSE.</small>'}</div>`;
};

if (typeof qaHistoryMarkup === 'function') {
  const inheritedQAHistoryMarkupV18100 = qaHistoryMarkup;
  qaHistoryMarkup = function qaHistoryMarkupV18100(){
    if (Array.isArray(qaManifest) && !qaManifest.some((x)=>String(x?.version||'')===DEPULSE_RELEASE_VERSION)) {
      qaManifest = [DEPULSE_RELEASE_QA_ENTRY, ...qaManifest];
    }
    return inheritedQAHistoryMarkupV18100();
  };
}
