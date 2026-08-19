'use strict';

const fs=require('fs');
const vm=require('vm');
const assert=require('assert');

const ownerSource=fs.readFileSync('renderer/documentation-ui.js','utf8');
const accessSource=fs.readFileSync('renderer/documentation-access-v18.6.js','utf8');
const index=fs.readFileSync('renderer/index.html','utf8');
const releaseIdentity=JSON.parse(fs.readFileSync('release_identity.json','utf8'));
const context=vm.createContext({console});
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
function architectureDiagram(kind){return '<div data-diagram="'+esc(kind)+'"></div>'}
function render(){renderCalls++}
function renderMarkdown(){return 'legacy markdown'}
async function hydrateDocumentation(){hydrateCalls++;return documentationTab}
function renderDocumentation(){return '<article>legacy '+documentationTab+'</article>'}
async function fetch(url){hydrateCalls++;return {ok:true,status:200,text:async()=>url}}
`,context);
vm.runInContext(ownerSource,context,{filename:'documentation-ui.js'});

const owner=vm.runInContext('__DE_PULSE_DOCUMENTATION_UI__',context);
assert.strictEqual(owner.owner,'renderer/documentation-ui.js');
assert.strictEqual(owner.registry().state,'ACTIVE_OWNER_WITH_LEGACY_FALLBACK');
assert(vm.runInContext('renderDocumentation()',context).includes('data-render-owner="documentation-ui"'));

vm.runInContext(accessSource,context,{filename:'documentation-access-v18.6.js'});

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
  const ownerTag=`<script src="documentation-ui.js?v=${releaseIdentity.version}"></script>`;
  const accessTag=`<script src="documentation-access-v18.6.js?v=${releaseIdentity.version}"></script>`;
  assert(index.includes(ownerTag),'index must load capability-oriented Documentation owner with canonical release cache identity');
  assert(index.includes(accessTag),'index must load v18.6 documentation access extension with canonical release cache identity');
  assert(index.indexOf(ownerTag)<index.indexOf(accessTag),'owner must load before access decorator');
  assert(!accessSource.includes('/api/'),'composition extension must not create a parallel authorization API');
  console.log('v18.6 Documentation owner + role access regression PASS');
}).catch(err=>{console.error(err);process.exitCode=1});
