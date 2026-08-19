(() => {
  'use strict';

  const baseRenderDiscovery = renderDiscovery;
  const baseSetPage = setPage;
  const baseBindDynamic = bindDynamic;
  const legacyEvidencePages = new Set(['news', 'earnings', 'filings']);
  let pendingResearchEvidenceTab = '';

  function marketActivitySeedRows() {
    const activity = runtime?.marketActivity || {};
    const combined = [
      ...(activity.mostActive || []).slice(0, 5),
      ...(activity.gainers || []).slice(0, 3),
      ...(activity.losers || []).slice(0, 2)
    ];
    return combined.filter((row, index, all) => {
      const symbol = String(row?.symbol || '').toUpperCase();
      return symbol && all.findIndex(candidate => String(candidate?.symbol || '').toUpperCase() === symbol) === index;
    });
  }

  function marketActivitySupportingDrilldown() {
    if (runtime?.marketActivity?.status !== 'AVAILABLE') return '';
    const rows = marketActivitySeedRows();
    return `<details class="discovery-supporting-input" data-market-activity-drilldown>
      <summary>
        <span><small>Supporting discovery input</small><b>Market Activity Seeds</b></span>
        <em>${rows.length} current seed${rows.length === 1 ? '' : 's'} · optional detail</em>
      </summary>
      <div class="discovery-supporting-input-body">
        <p>Most-active and mover observations are reused as low-cost Scanner / Opportunity Radar seeds. They do not bypass horizon filters, deterministic desk scoring, or Research qualification.</p>
        <div class="market-activity-seeds supporting-seeds">${rows.map(row => `<span><b>${esc(String(row.symbol || '').toUpperCase())}</b><small>${num(row.changePercent) ? pctSigned(row.changePercent) : 'Active'}</small></span>`).join('') || '<div class="empty">No entitled activity seed data.</div>'}</div>
      </div>
    </details>`;
  }

  function demoteProminentMarketActivity(html) {
    const marker = '<section class="card discovery-market-activity full-width">';
    const start = html.indexOf(marker);
    if (start < 0) return html;
    const endMarker = '</section>';
    const end = html.indexOf(endMarker, start);
    if (end < 0) return html;
    return html.slice(0, start) + marketActivitySupportingDrilldown() + html.slice(end + endMarker.length);
  }

  renderDiscovery = function renderDiscoveryV186() {
    return demoteProminentMarketActivity(baseRenderDiscovery());
  };

  function evidenceResearchTab(kind) {
    return kind === 'filings' || kind === 'sec' ? 'filings' : kind === 'news' || kind === 'earnings' ? 'catalysts' : 'overview';
  }

  function openCanonicalTickerEvidence(symbol, kind, origin = page) {
    const target = String(symbol || state?.ui?.selectedTicker || researchSymbol || '').trim().toUpperCase();
    if (!target) {
      baseSetPage('market-intelligence');
      return;
    }
    pendingResearchEvidenceTab = evidenceResearchTab(kind);
    openResearch(target, origin);
    researchTab = pendingResearchEvidenceTab;
    render();
  }

  // Standalone News / Earnings / Filings were legacy market-wide evidence homes.
  // Safe direct/internal navigation remains supported by redirecting to the
  // canonical Market Intelligence / Event Intelligence surface.
  setPage = function setPageV186(target) {
    const next = String(target || 'dashboard').toLowerCase();
    if (legacyEvidencePages.has(next)) {
      pendingResearchEvidenceTab = '';
      return baseSetPage('market-intelligence');
    }
    return baseSetPage(target);
  };

  function bindCanonicalEvidenceRoutes() {
    $$('[data-jump]').forEach(button => {
      const target = String(button.dataset.jump || '').toLowerCase();
      if (!legacyEvidencePages.has(target)) return;
      button.onclick = event => {
        event.preventDefault();
        event.stopPropagation();
        const symbol = String(button.dataset.symbol || '').trim().toUpperCase();
        if (symbol) openCanonicalTickerEvidence(symbol, target, page);
        else baseSetPage('market-intelligence');
      };
    });

    $$('[data-readiness-drill]').forEach(button => {
      const [target, symbol] = String(button.dataset.readinessDrill || '').split(':');
      if (!['research', 'sec', 'earnings'].includes(target)) return;
      button.onclick = event => {
        event.preventDefault();
        event.stopPropagation();
        openCanonicalTickerEvidence(symbol, target, 'maintenance');
      };
    });
  }

  // Research hydration resets the visible tab to Overview while it reconciles
  // evidence. Preserve an explicitly routed canonical evidence destination and
  // apply it after the source refresh completes.
  if (typeof v151SelectResearchTarget === 'function') {
    const baseSelectResearchTarget = v151SelectResearchTarget;
    v151SelectResearchTarget = async function v151SelectResearchTargetV186(symbol) {
      const desiredTab = pendingResearchEvidenceTab;
      const result = await baseSelectResearchTarget(symbol);
      if (desiredTab && page === 'research') {
        researchTab = desiredTab;
        pendingResearchEvidenceTab = '';
        render();
      }
      return result;
    };
  }

  bindDynamic = function bindDynamicV186SurfaceConsolidation() {
    baseBindDynamic();
    bindCanonicalEvidenceRoutes();
  };

  window.__v186SurfaceConsolidation = Object.freeze({
    demoteProminentMarketActivity,
    marketActivitySupportingDrilldown,
    evidenceResearchTab,
    openCanonicalTickerEvidence,
    bindCanonicalEvidenceRoutes
  });
})();
