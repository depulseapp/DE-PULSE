'use strict';
(function(root){
  function normalizedTradeabilityState(state){
    return String(state || 'DATA DEGRADED').trim().toUpperCase();
  }

  function tradeabilityScoreLabel(state, score){
    const normalized = normalizedTradeabilityState(state);
    if(normalized === 'DATA DEGRADED' || normalized === 'UNAVAILABLE') return 'UNAVAILABLE';
    const value = Number(score);
    return Number.isFinite(value) ? `${Math.round(value)}/100` : 'UNAVAILABLE';
  }

  function reconcileDashboardTradeability(doc){
    for(const row of doc.querySelectorAll('.market-intelligence-summary .compact-row')){
      const label = row.querySelector('b');
      const value = row.querySelector('span');
      if(!label || !value || label.textContent.trim() !== 'Tradeability') continue;
      const state = String(value.textContent || '').split('·')[0].trim() || 'DATA DEGRADED';
      const score = tradeabilityScoreLabel(state, null);
      if(score !== 'UNAVAILABLE') continue;
      const desired = `${normalizedTradeabilityState(state)} · UNAVAILABLE`;
      if(value.textContent !== desired) value.textContent = desired;
    }
  }

  function reconcileMarketIntelligenceCard(doc){
    for(const card of doc.querySelectorAll('.market-intelligence-page article.card')){
      const label = card.querySelector('small');
      if(!label || label.textContent.trim() !== 'Market Tradeability') continue;
      const heading = card.querySelector('h3');
      const score = card.querySelector('b');
      if(!heading || !score) continue;
      const desired = tradeabilityScoreLabel(heading.textContent, score.textContent);
      if(desired === 'UNAVAILABLE' && score.textContent !== desired) score.textContent = desired;
    }
  }

  function reconcileMarketIntelligenceTruth(doc){
    if(!doc || typeof doc.querySelectorAll !== 'function') return;
    reconcileDashboardTradeability(doc);
    reconcileMarketIntelligenceCard(doc);
  }

  const api = {
    normalizedTradeabilityState,
    tradeabilityScoreLabel,
    reconcileMarketIntelligenceTruth,
  };

  if(typeof module !== 'undefined' && module.exports) module.exports = api;
  if(root) root.DEPULSEMarketIntelligenceTruth = api;

  if(typeof document !== 'undefined' && typeof MutationObserver !== 'undefined'){
    let queued = false;
    const schedule = () => {
      if(queued) return;
      queued = true;
      const run = () => {
        queued = false;
        reconcileMarketIntelligenceTruth(document);
      };
      if(typeof requestAnimationFrame === 'function') requestAnimationFrame(run);
      else setTimeout(run, 0);
    };
    new MutationObserver(schedule).observe(document.documentElement, {childList:true, subtree:true, characterData:true});
    schedule();
  }
})(typeof window !== 'undefined' ? window : null);
