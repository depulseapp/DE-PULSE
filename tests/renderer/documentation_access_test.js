'use strict';

const fs=require('fs');
const vm=require('vm');
const assert=require('assert');
const crypto=require('crypto');

const architecturePath='renderer/documentation-architecture.js';
const ownerPath='renderer/documentation-ui.js';
const accessPath='renderer/documentation-access.js';
const architectureSource=fs.readFileSync(architecturePath,'utf8');
const ownerSource=fs.readFileSync(ownerPath,'utf8');
const accessSource=fs.readFileSync(accessPath,'utf8');
const index=fs.readFileSync('renderer/index.html','utf8');
const context=vm.createContext({console});

function gitBlobToken(path){
  const data=fs.readFileSync(path);
  const header=Buffer.from(`blob ${data.length}\0`,'utf8');
  return crypto.createHash('sha1').update(header).update(data).digest('hex').slice(0,16);
}

vm.runInContext(`
let authPrincipal={role:'USER'};
let documentationTab='developer';
let docCache={user:'User body',developer:'Developer body',limitations:'Limitations body'};
let page='documentation';
let hydrateCalls=0;
let renderCalls=0;
const DOC_BRAND='DE.PULSE';
const esc=v=>String(v??'').replace(/[&<>"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]));
function docInline(v){return esc(v)}
function diagramNode(title,detail='',tone=''){return '<div class="arch-node '+tone+'"><b>'+esc(title)+'</b><small>'+esc(detail)+'</small></div>'}
function diagramArrow(label=''){return '<div class="arch-arrow">'+esc(label)+'</div>'}
function architectureDiagram(kind){return '<div data-legacy-diagram="'+esc(kind)+'"></div>'}
function render(){renderCalls++}
function renderMarkdown(){return 'legacy markdown'}
async function hydrateDocumentation(){hydrateCalls++;return documentationTab}
function renderDocumentation(){return '<article>legacy '+documentationTab+'</article>'}
async function fetch(url){hydrateCalls++;return {ok:true,status:200,text:async()=>url}}
`,context);

vm.runInContext(architectureSource,context,{filename:'documentation-architecture.js'});
const architectureOwner=vm.runInContext('__DE_PULSE_DOCUMENTATION_ARCHITECTURE__',context);
assert.strictEqual(architectureOwner.owner,'renderer/documentation-architecture.js');
assert.strictEqual(architectureOwner.registry().state,'ACTIVE_CANONICAL_OWNER');
assert.strictEqual(architectureOwner.registry().legacyMonolithFallbackPresent,true,'extracted owner should explicitly record the remaining monolith fallback until deletion evidence exists');
const feedDiagram=vm.runInContext("architectureDiagram('feeds')",context);
const overallDiagram=vm.runInContext("architectureDiagram('overall')",context);
assert(feedDiagram.includes('Dynamic Multi-Feed Allocation'),'live-feed documentation must name the canonical multi-feed allocator');
assert(feedDiagram.includes('Smart Provider Router v2'),'live-feed documentation must preserve Router v2 routing authority');
assert(feedDiagram.includes('preferred live streaming pool'),'current Alpaca preferred live allocation must be explicit');
assert(feedDiagram.includes('overflow/secondary live pool'),'current Finnhub overflow/secondary live allocation must be explicit');
assert(feedDiagram.includes('GLD, SLV and USO'),'actionable GLD/SLV/USO exceptions must remain documented');
assert(!feedDiagram.includes('Finnhub WebSocket</b><small>primary live trades')&&!feedDiagram.includes('Alpaca IEX</b><small>controlled fallback pool'),'reversed legacy Finnhub-primary / Alpaca-fallback wording must not reappear');
assert(overallDiagram.includes('direct filings / Form 4 authority'),'direct SEC/EDGAR authority must remain explicit');
assert(overallDiagram.includes('No Execution'),'permanent No Execution boundary must remain explicit');

vm.runInContext(ownerSource,context,{filename:'documentation-ui.js'});
const owner=vm.runInContext('__DE_PULSE_DOCUMENTATION_UI__',context);
assert.strictEqual(owner.owner,'renderer/documentation-ui.js');
assert.strictEqual(owner.registry().state,'ACTIVE_CANONICAL_OWNER_WITH_LEGACY_RENDER_FALLBACK');
assert.strictEqual(owner.registry().architectureOwner,'renderer/documentation-architecture.js');
assert(vm.runInContext('renderDocumentation()',context).includes('data-render-owner="documentation-ui"'));
assert(vm.runInContext('renderDocumentation()',context).includes('data-architecture-owner="renderer/documentation-architecture.js"'));

vm.runInContext(accessSource,context,{filename:'documentation-access.js'});

const userHTML=vm.runInContext('renderDocumentation()',context);
assert(!userHTML.includes('data-doc-tab="developer"'),'USER must not receive developer tab');
assert.strictEqual(vm.runInContext('documentationTab',context),'user','forbidden programmatic developer selection must normalize to user docs');
assert.strictEqual(vm.runInContext('__v186DocumentationAccess.canViewDeveloper()',context),false);
assert(userHTML.includes('data-render-owner="documentation-ui"'),'role decorator must preserve capability owner markup');

vm.runInContext(`authPrincipal={role:'DEMO'};documentationTab='developer'`,context);
assert(!vm.runInContext('renderDocumentation()',context).includes('Developer Documentation'),'DEMO must not receive developer docs');

for(const role of ['ADMIN','OWNER','SUPER_OWNER']){
  vm.runInContext(`authPrincipal={role:'${role}'};documentationTab='developer'`,context);
  const html=vm.runInContext('renderDocumentation()',context);
  assert(html.includes('data-doc-tab="developer"'),`${role} should receive developer tab`);
  assert(html.includes('Developer body'),`${role} developer selection should remain intact`);
  assert(html.includes('data-render-owner="documentation-ui"'),`${role} must retain canonical owner markup`);
}

vm.runInContext(`authPrincipal={role:'USER'};documentationTab='developer';docCache.user=null`,context);
vm.runInContext('hydrateDocumentation()',context).then(()=>{
  assert.strictEqual(vm.runInContext('documentationTab',context),'user','hydrate must not fetch forbidden developer docs');
  assert.strictEqual(vm.runInContext('docCache.user',context),'/docs/user.md','normalized owner hydrate must fetch user docs');
  const architectureTag=`<script src="documentation-architecture.js?v=${gitBlobToken(architecturePath)}"></script>`;
  const ownerTag=`<script src="documentation-ui.js?v=${gitBlobToken(ownerPath)}"></script>`;
  const accessTag=`<script src="documentation-access.js?v=${gitBlobToken(accessPath)}"></script>`;
  assert(index.includes(architectureTag),'index must load dedicated Documentation architecture owner with content-derived cache identity');
  assert(index.includes(ownerTag),'index must load capability-oriented Documentation owner with content-derived cache identity');
  assert(index.includes(accessTag),'index must load documentation access extension with content-derived cache identity');
  assert(index.indexOf(architectureTag)<index.indexOf(ownerTag),'architecture owner must load before Documentation UI');
  assert(index.indexOf(ownerTag)<index.indexOf(accessTag),'Documentation UI owner must load before access decorator');
  assert(!accessSource.includes('/api/'),'composition extension must not create a parallel authorization API');
  console.log('Documentation architecture + owner + role access regression PASS · current provider truth locked');
}).catch(err=>{console.error(err);process.exitCode=1});
