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
        const latency=completed>0?`median p50 ${Number(row.p50LatencyMs)} ms · p95 ${Number(row.p95LatencyMs)} ms · avg ${Number(row.averageLatencyMs)} ms`:'latency unmeasured';
        detail.textContent=`Transport reliability · ${reliability} · ${latency} · ${Number(row.rateLimited||0)} rate limited · peak ${Number(row.peakInFlight||0)}`;
      }
      card.dataset.providerTransportReliability='true';
    }
  }

  function metric(value,suffix=''){
    return value===null||value===undefined?'UNKNOWN':`${value}${suffix}`;
  }

  function operationalScorecard(row){
    const card=document.createElement('div');
    card.className='provider-scorecard-card';
    card.dataset.providerScorecard=String(row.provider||'');
    const title=document.createElement('span');
    title.textContent=`${String(row.provider||'Provider')} Operational Scorecard`;
    const headline=document.createElement('b');
    headline.textContent=String(row.state||'UNOBSERVED').toUpperCase();
    const detail=document.createElement('small');
    const transport=row.transportMeasurementState==='MEASURED'?`${metric(row.successPct,'%')} success · p50 ${metric(row.p50LatencyMs,'ms')} · p95 ${metric(row.p95LatencyMs,'ms')}`:'transport UNKNOWN';
    const freshness=(row.freshnessStates||[]).join(' · ')||'UNKNOWN';
    const headroom=`REST ${metric(row.requestBudgetRemaining)} · live ${metric(row.liveSubscriptionAvailable)}`;
    const usefulness=row.agreementPct===null||row.agreementPct===undefined?String(row.usefulnessMeasurementState||'UNKNOWN'):`${formatPercent(row.agreementPct)} agreement`;
    detail.textContent=`health ${String(row.healthMeasurementState||'UNKNOWN')} · freshness ${freshness} · ${transport} · headroom ${headroom} · rights ${String(row.rightsMeasurementState||'UNKNOWN')}/${String(row.rightsReviewState||'UNKNOWN')} · cost ${String(row.costMeasurementState||'UNKNOWN')} (${String(row.costClass||'UNKNOWN')}) · usefulness ${usefulness} · OBSERVABILITY ONLY · no routing effect`;
    card.append(title,headline,detail);
    return card;
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
    document.querySelectorAll('.provider-scorecard-card').forEach(card=>card.remove());
    if(!isPrivilegedProviderTelemetryRole())return;
    const grid=document.querySelector('.performance-observability .maintenance-health-grid');
    if(!grid)return;
    const load=runtime?.runtimeLoad||{};
    const providers=Array.isArray(load.providerRequests)?load.providerRequests:[];
    const usefulness=Array.isArray(load.providerUsefulness)?load.providerUsefulness:[];
    const scorecards=Array.isArray(load.providerScorecards)?load.providerScorecards:[];
    enrichTransportReliability(grid,providers);
    for(const row of scorecards)grid.appendChild(operationalScorecard(row));
    for(const row of usefulness)grid.appendChild(semanticUsefulnessCard(row));
  }

  afterRender=function providerTelemetryAfterRender(){
    originalAfterRender();
    applyProviderTelemetryProjection();
  };

  window.addEventListener('DOMContentLoaded',()=>requestAnimationFrame(applyProviderTelemetryProjection),{once:true});
})();
