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
      bar.setAttribute('aria-label', 'Market state and data coverage');
      topbar.insertAdjacentElement('afterend', bar);
    }

    const session = document.getElementById('market-session-context');
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

    if (session && session.parentElement !== bar) bar.appendChild(session);
    if (status && status.parentElement !== summary) summary.insertBefore(status, summary.firstChild);
    if (summary.parentElement !== bar) bar.appendChild(summary);
    if (toggle && toggle.parentElement !== bar) bar.appendChild(toggle);
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
