'use strict';
/* DE.PULSE v18.6.1 patch contract.
   The watchlist extension is a separate classic script and must not depend on
   an undeclared desk-key variable. The patch also binds build identity and QA
   history to v18.6.1 while the large frozen renderer remains unchanged. */
var DESKS = Object.freeze(['day', 'swing', 'long']);
var DEPULSE_PATCH_VERSION = '18.6.1';
var DEPULSE_PATCH_BUILD_ID = 'v18.6.1-stable-20260819';
var DEPULSE_PATCH_QA_ENTRY = Object.freeze({
  version: DEPULSE_PATCH_VERSION,
  date: '2026-08-19',
  status: 'STABLE',
  summary: 'Focused watchlist removal, membership-toggle and header-alert hardening with inherited v18.6.0 certification coverage and v18.6.1 edge/browser proof.',
  file: 'v18.6.1.txt',
  buildId: DEPULSE_PATCH_BUILD_ID,
  checkpoint: 'release/v18.6.1/patch_contract.json'
});

buildIdentityIntegrityMarkup = function buildIdentityIntegrityMarkupV1861(){
  const v=String(state?.version||''), b=String(state?.buildId||'');
  const ok=v===DEPULSE_PATCH_VERSION && b===DEPULSE_PATCH_BUILD_ID;
  return `<div class="build-integrity ${ok?'ok':'bad'}"><b>${ok?'BUILD IDENTITY VERIFIED':'BUILD IDENTITY MISMATCH'}</b><span>${esc(v||'unknown')} · ${esc(b||'unknown')}</span>${ok?'':'<small>Running backend does not match the packaged v18.6.1 Stable renderer. Restart De-Pulse.app.</small>'}</div>`;
};

if (typeof qaHistoryMarkup === 'function') {
  const inheritedQAHistoryMarkup = qaHistoryMarkup;
  qaHistoryMarkup = function qaHistoryMarkupV1861(){
    if (Array.isArray(qaManifest) && !qaManifest.some((x)=>String(x?.version||'')===DEPULSE_PATCH_VERSION)) {
      qaManifest = [DEPULSE_PATCH_QA_ENTRY, ...qaManifest];
    }
    return inheritedQAHistoryMarkup();
  };
}
