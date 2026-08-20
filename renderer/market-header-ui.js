(() => {
  'use strict';

  const OWNER='renderer/market-header-ui.js';
  const OWNER_VERSION=1;

  function ensureSecondaryMarketStatus() {
    const topbar = document.querySelector('.topbar');
    if (!topbar) return null;
    let bar = document.getElementById('market-status-bar');
    if (!bar) {
      bar = document.createElement('div');
      bar.id = 'market-status-bar';
      bar.className = 'market-status-bar';
      bar.setAttribute('aria-label', 'Market Pulse Ribbon: session, coverage, clocks and data control');
      topbar.insertAdjacentElement('afterend', bar);
    }

    let content = bar.querySelector('.market-status-content');
    if (!content) {
      content = document.createElement('div');
      content.className = 'market-status-content';
      bar.appendChild(content);
    }

    const session = document.getElementById('market-session-context');
    const clocks = document.querySelector('.market-clocks');
    const status = document.getElementById('runtime-status');
    const toggle = document.getElementById('runtime-toggle');
    let summary = document.getElementById('market-data-summary');
    if (!summary) {
      summary = document.createElement('div');
      summary.id = 'market-data-summary';
      summary.className = 'market-data-summary';
      summary.setAttribute('aria-label', 'Market data coverage');
      const detail = document.createElement('span');
      detail.id = 'market-data-detail';
      detail.className = 'market-data-detail';
      detail.textContent = 'Market data state unavailable';
      summary.appendChild(detail);
    }

    if (status && status.parentElement !== summary) summary.insertBefore(status, summary.firstChild);

    // Keep the ribbon compact and deterministic: market state, coverage,
    // complete clocks, then the data control. appendChild also repairs order
    // when this capability is hot-reloaded over an older ribbon.
    if (session) content.appendChild(session);
    content.appendChild(summary);
    if (clocks) content.appendChild(clocks);
    if (toggle) content.appendChild(toggle);
    return bar;
  }

  const baseUpdateChrome = updateChrome;

  function updateSecondaryMarketStatus() {
    ensureSecondaryMarketStatus();
    const detail = document.getElementById('market-data-detail');
    if (!detail) return;
    const health = headerDataHealth();
    detail.textContent = health.detail || 'Market data state unavailable';
    detail.title = health.detail || health.label || 'Market data state unavailable';
  }

  ensureSecondaryMarketStatus();

  updateChrome = function updateChromeMarketHeader(changedSymbol = '') {
    baseUpdateChrome(changedSymbol);
    updateSecondaryMarketStatus();
  };

  const registry=globalThis.__DE_PULSE_RENDERER_OWNERS__||(globalThis.__DE_PULSE_RENDERER_OWNERS__={});
  registry.marketHeader={
    owner:OWNER,
    version:OWNER_VERSION,
    state:'ACTIVE_OWNER_WITH_COMPAT_ALIAS',
    responsibilities:['market-pulse-ribbon','session-context','market-data-health','market-clocks','data-runtime-control'],
    dependencies:['shared-chrome-update','header-data-health','existing-header-dom'],
    compatibilityAliases:['__v1851HeaderContracts'],
    legacyCompatibilityFile:'renderer/header-v18.5.1.js'
  };

  const api=Object.freeze({
    owner:OWNER,
    version:OWNER_VERSION,
    ensureSecondaryMarketStatus,
    updateSecondaryMarketStatus,
    registry:()=>registry.marketHeader
  });
  globalThis.__DE_PULSE_MARKET_HEADER_UI__=api;
  globalThis.__v1851HeaderContracts=api;
})();
