'use strict';

const fs=require('fs');
const vm=require('vm');
const assert=require('assert');
const crypto=require('crypto');
require('./market_header_owner_test.js');
require('./v18_8_1_trust_closure_test.js');

const architectureBytes=fs.readFileSync('renderer/documentation-architecture.js');
const ownerBytes=fs.readFileSync('renderer/documentation-ui.js');
const accessBytes=fs.readFileSync('renderer/documentation-access.js');
const architectureSource=architectureBytes.toString('utf8');
const ownerSource=ownerBytes.toString('utf8');
const accessSource=accessBytes.toString('utf8');
const index=fs.readFileSync('renderer/index.html','utf8');

function gitBlobToken(data){
  const bytes=Buffer.isBuffer(data)?data:Buffer.from(data,'utf8');
  return crypto.createHash('sha1')
    .update(Buffer.from(`blob ${bytes.length}\0`,'utf8'))
    .update(bytes)
    .digest('hex')
    .slice(0,16);
}

const context=vm.createContext({console});
vm.runInContext(`
let authPrincipal={role:'USER'};
let documentationTab='user';
let docCache={user:'# User\\nHello **DE.PULSE**',developer:'# Developer\\nSecret',limitations:'# Limits'};
let page='documentation';
let renderCalls=0;
const DOC_BRAND='<span>DE.PULSE</span>';
const esc=v=>String(v??'').replace(/[&<>\"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','\"':'&quot;',"'":'&#39;'}[c]));
function docInline(v){return esc(v).replace(/\\*\\*([^*]+)\\*\\*/g,'<strong>$1</strong>')}
function diagramNode(title,detail='',tone=''){return '<div class="arch-node '+tone+'"><b>'+esc(title)+'</b><small>'+esc(detail)+'</small></div>'}
function diagramArrow(label=''){return '<div class="arch-arrow">'+esc(label)+'</div>'}
function architectureDiagram(kind){return '<section data-legacy-diagram="'+esc(kind)+'"></section>'}
function render(){renderCalls++}
function renderMarkdown(){return 'LEGACY_MARKDOWN'}
async function hydrateDocumentation(){return 'LEGACY_HYDRATE'}
function renderDocumentation(){return 'LEGACY_RENDER'}
let fetchCalls=[];
async function fetch(url){fetchCalls.push(url);return {ok:true,status:200,text:async()=>url.includes('developer')?'# Developer\\nLoaded':'# Loaded'}}
`,context);

const legacy=vm.runInContext('({renderMarkdown,hydrateDocumentation,renderDocumentation,architectureDiagram})',context);
vm.runInContext(architectureSource,context,{filename:'documentation-architecture.js'});
const architectureOwner=vm.runInContext('__DE_PULSE_DOCUMENTATION_ARCHITECTURE__',context);
assert.strictEqual(architectureOwner.owner,'renderer/documentation-architecture.js');
assert.strictEqual(architectureOwner.registry().state,'ACTIVE_CANONICAL_OWNER');
assert.notStrictEqual(vm.runInContext('architectureDiagram',context),legacy.architectureDiagram);

vm.runInContext(ownerSource,context,{filename:'documentation-ui.js'});

const owner=vm.runInContext('__DE_PULSE_DOCUMENTATION_UI__',context);
assert.strictEqual(owner.owner,'renderer/documentation-ui.js');
assert.strictEqual(owner.version,2);
assert.strictEqual(owner.registry().state,'ACTIVE_OWNER_WITH_LEGACY_FALLBACK');
assert.strictEqual(owner.registry().architectureOwner,'renderer/documentation-architecture.js');
assert(owner.registry().dependencies.includes('documentation-architecture'));
assert.notStrictEqual(vm.runInContext('renderMarkdown',context),legacy.renderMarkdown);
assert.notStrictEqual(vm.runInContext('hydrateDocumentation',context),legacy.hydrateDocumentation);
assert.notStrictEqual(vm.runInContext('renderDocumentation',context),legacy.renderDocumentation);

const userHTML=vm.runInContext('renderDocumentation()',context);
assert(userHTML.includes('data-render-owner="documentation-ui"'));
assert(userHTML.includes('data-architecture-owner="renderer/documentation-architecture.js"'));
assert(userHTML.includes('User Documentation'));
assert(userHTML.includes('<strong>DE.PULSE</strong>'));
const markdown=vm.runInContext("renderMarkdown('## Heading\\n- One\\n[[diagram:overall]]')",context);
assert(markdown.includes('<h2>Heading</h2>'));
assert(markdown.includes('<li>One</li>'));
assert(markdown.includes('Current System Architecture'));
assert(markdown.includes('Smart Provider Router v2'));
assert(markdown.includes('direct filings / Form 4 authority'));
assert(markdown.includes('No Execution'));
assert(!markdown.includes('data-legacy-diagram="overall"'));

vm.runInContext(accessSource,context,{filename:'documentation-access.js'});
vm.runInContext("authPrincipal={role:'USER'};documentationTab='developer'",context);
const restricted=vm.runInContext('renderDocumentation()',context);
assert(!restricted.includes('data-doc-tab="developer"'),'USER must not receive developer tab after owner extraction');
assert.strictEqual(vm.runInContext('documentationTab',context),'user');

vm.runInContext("authPrincipal={role:'ADMIN'};documentationTab='developer'",context);
const admin=vm.runInContext('renderDocumentation()',context);
assert(admin.includes('data-doc-tab="developer"'),'ADMIN must retain developer documentation access');
assert(admin.includes('data-render-owner="documentation-ui"'),'access decorator must preserve canonical owner markup');

vm.runInContext("authPrincipal={role:'USER'};documentationTab='limitations';docCache.limitations=null",context);
vm.runInContext('hydrateDocumentation()',context).then(()=>{
  assert.deepStrictEqual(Array.from(vm.runInContext('fetchCalls',context)),['/docs/limitations.md']);
  assert.strictEqual(vm.runInContext('renderCalls',context),1);
  const architectureTag=`<script src="documentation-architecture.js?v=${gitBlobToken(architectureBytes)}"></script>`;
  const ownerTag=`<script src="documentation-ui.js?v=${gitBlobToken(ownerBytes)}"></script>`;
  const accessTag=`<script src="documentation-access.js?v=${gitBlobToken(accessBytes)}"></script>`;
  assert(index.includes(architectureTag),'index must load canonical Documentation architecture owner with canonical cache identity');
  assert(index.includes(ownerTag),'index must load capability-oriented Documentation owner with canonical cache identity');
  assert(index.includes(accessTag),'index must retain role-access decorator');
  assert(index.indexOf(architectureTag)<index.indexOf(ownerTag),'Documentation architecture owner must load before Documentation UI');
  assert(index.indexOf(ownerTag)<index.indexOf(accessTag),'Documentation owner must load before role-access decorator');
  console.log('Documentation architecture + capability owner regression PASS');
}).catch(err=>{console.error(err);process.exitCode=1});
