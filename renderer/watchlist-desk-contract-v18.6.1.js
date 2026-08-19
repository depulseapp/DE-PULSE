'use strict';
/* DE.PULSE v18.6.1 patch contract.
   The watchlist extension is a separate classic script and must not depend on
   an undeclared desk-key variable. The patch also binds build-identity UI to
   the v18.6.1 backend while the large frozen renderer remains unchanged. */
var DESKS = Object.freeze(['day', 'swing', 'long']);
var DEPULSE_PATCH_VERSION = '18.6.1';
var DEPULSE_PATCH_BUILD_ID = 'v18.6.1-stable-20260819';

buildIdentityIntegrityMarkup = function buildIdentityIntegrityMarkupV1861(){
  const v=String(state?.version||''), b=String(state?.buildId||'');
  const ok=v===DEPULSE_PATCH_VERSION && b===DEPULSE_PATCH_BUILD_ID;
  return `<div class="build-integrity ${ok?'ok':'bad'}"><b>${ok?'BUILD IDENTITY VERIFIED':'BUILD IDENTITY MISMATCH'}</b><span>${esc(v||'unknown')} · ${esc(b||'unknown')}</span>${ok?'':'<small>Running backend does not match the packaged v18.6.1 Stable renderer. Restart De-Pulse.app.</small>'}</div>`;
};
