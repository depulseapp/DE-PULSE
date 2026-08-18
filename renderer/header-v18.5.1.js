(() => {
  'use strict';

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
    // when this extension is hot-reloaded over an older ribbon.
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

  updateChrome = function updateChromeV1851(changedSymbol = '') {
    baseUpdateChrome(changedSymbol);
    updateSecondaryMarketStatus();
  };

  window.__v1851HeaderContracts = Object.freeze({
    ensureSecondaryMarketStatus,
    updateSecondaryMarketStatus
  });
})();
