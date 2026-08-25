'use strict';

(function providerTelemetryUIExtension(){
  const privilegedRoles=new Set(['SUPER_OWNER','OWNER','ADMIN']);
  const originalAfterRender=afterRender;

  function isPrivilegedProviderTelemetryRole(){
    const role=String(authPrincipal?.role||'').trim().toUpperCase();
    return privilegedRoles.has(role);
  }

  function formatPercent(value){
    const n=Number(value);
    return Number.isFinite(n)?`${n.toFixed(1)}%`:'—';
  }

  function providerRequestCard(grid,provider){
    const expected=`${String(provider||'').trim()} Requests`;
    return [...grid.children].find(card=>String(card.querySelector('span')?.textContent||'').trim()===expected)||null;
  }

  function enrichTransportReliability(grid,providers){
    for(const row of providers){
      const card=providerRequestCard(grid,row.provider);
      if(!card)continue;
      const completed=Number(row.successes||0)+Number(row.errors||0);
      const headline=card.querySelector('b');
      const detail=card.querySelector('small');
      if(headline)headline.textContent=`${Number(row.requestsLastMinute||0)}/min · ${Number(row.inFlight||0)} active · ${completed} completed`;
      if(detail){
        const reliability=completed>0?`${formatPercent(row.successPct)} success · ${Number(row.successes||0)}/${Number(row.errors||0)} success/failure`:'No completed requests yet';
        detail.textContent=`Transport reliability · ${reliability} · median p50 ${Number(row.p50LatencyMs||0)} ms · p95 ${Number(row.p95LatencyMs||0)} ms · avg ${Number(row.averageLatencyMs||0)} ms · ${Number(row.rateLimited||0)} rate limited · peak ${Number(row.peakInFlight||0)}`;
      }
      card.dataset.providerTransportReliability='true';
    }
  }

  function semanticUsefulnessCard(row){
    const card=document.createElement('div');
    card.className='provider-usefulness-card';
    card.dataset.providerUsefulness=String(row.provider||'');
    const title=document.createElement('span');
    title.textContent=`${String(row.provider||'Provider')} Semantic Evidence`;
    const headline=document.createElement('b');
    const detail=document.createElement('small');
    const state=String(row.state||'INSUFFICIENT').toUpperCase();
    if(state==='OBSERVING'&&Number.isFinite(Number(row.agreementPct))){
      headline.textContent=`${formatPercent(row.agreementPct)} cross-source agreement`;
    }else{
      headline.textContent='INSUFFICIENT EVIDENCE';
    }
    detail.textContent=`${Number(row.eligibleSamples||0)} eligible · ${Number(row.crossSourceSamples||0)} cross-source · ${Number(row.agreementSamples||0)} agreed · ${Number(row.conflictSamples||0)} conflict · ${Number(row.singleSourceSamples||0)} single-source · ${Number(row.canonicalSelections||0)} canonical selections · ${Number(row.excludedSamples||0)} excluded · ADVISORY ONLY · no routing effect`;
    card.append(title,headline,detail);
    return card;
  }

  function applyProviderTelemetryProjection(){
    document.querySelectorAll('.provider-usefulness-card').forEach(card=>card.remove());
    if(!isPrivilegedProviderTelemetryRole())return;
    const grid=document.querySelector('.performance-observability .maintenance-health-grid');
    if(!grid)return;
    const load=runtime?.runtimeLoad||{};
    const providers=Array.isArray(load.providerRequests)?load.providerRequests:[];
    const usefulness=Array.isArray(load.providerUsefulness)?load.providerUsefulness:[];
    enrichTransportReliability(grid,providers);
    for(const row of usefulness)grid.appendChild(semanticUsefulnessCard(row));
  }

  afterRender=function providerTelemetryAfterRender(){
    originalAfterRender();
    applyProviderTelemetryProjection();
  };

  window.addEventListener('DOMContentLoaded',()=>requestAnimationFrame(applyProviderTelemetryProjection),{once:true});
})();
