(() => {
  'use strict';

  const baseUpdateChrome = updateChrome;

  function updateSecondaryMarketStatus() {
    const detail = document.getElementById('market-data-detail');
    if (!detail) return;
    const health = headerDataHealth();
    detail.textContent = health.detail || 'Market data state unavailable';
    detail.title = health.detail || health.label || 'Market data state unavailable';
  }

  updateChrome = function updateChromeV1851(changedSymbol = '') {
    baseUpdateChrome(changedSymbol);
    updateSecondaryMarketStatus();
  };

  window.__v1851HeaderContracts = Object.freeze({ updateSecondaryMarketStatus });
})();
