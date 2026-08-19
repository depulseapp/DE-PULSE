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
let documentationTab='user';
let docCache={user:'# User\\nHello **DE.PULSE**',developer:'# Developer\\nSecret',limitations:'# Limits'};
let page='documentation';
let renderCalls=0;
const DOC_BRAND='<span>DE.PULSE</span>';
const esc=v=>String(v??'').replace(/[&<>\"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','\"':'&quot;',"'":'&#39;'}[c]));
function docInline(v){return esc(v).replace(/\\*\\*([^*]+)\\*\\*/g,'<strong>$1</strong>')}
function architectureDiagram(kind){return '<section data-legacy-diagram="'+esc(kind)+'"></section>'}
function render(){renderCalls++}
function renderMarkdown(){return 'LEGACY_MARKDOWN'}
async function hydrateDocumentation(){return 'LEGACY_HYDRATE'}
function renderDocumentation(){return 'LEGACY_RENDER'}
let fetchCalls=[];
async function fetch(url){fetchCalls.push(url);return {ok:true,status:200,text:async()=>url.includes('developer')?'# Developer\\nLoaded':'# Loaded'}}
`,context);

const legacy=vm.runInContext('({renderMarkdown,hydrateDocumentation,renderDocumentation})',context);
vm.runInContext(ownerSource,context,{filename:'documentation-ui.js'});

const owner=vm.runInContext('__DE_PULSE_DOCUMENTATION_UI__',context);
assert.strictEqual(owner.owner,'renderer/documentation-ui.js');
assert.strictEqual(owner.version,1);
assert.strictEqual(owner.registry().state,'ACTIVE_OWNER_WITH_LEGACY_FALLBACK');
assert(owner.registry().dependencies.includes('legacy-architecture-diagram'));
assert.notStrictEqual(vm.runInContext('renderMarkdown',context),legacy.renderMarkdown);
assert.notStrictEqual(vm.runInContext('hydrateDocumentation',context),legacy.hydrateDocumentation);
assert.notStrictEqual(vm.runInContext('renderDocumentation',context),legacy.renderDocumentation);

const userHTML=vm.runInContext('renderDocumentation()',context);
assert(userHTML.includes('data-render-owner="documentation-ui"'));
assert(userHTML.includes('User Documentation'));
assert(userHTML.includes('<strong>DE.PULSE</strong>'));
const markdown=vm.runInContext("renderMarkdown('## Heading\\n- One\\n[[diagram:overall]]')",context);
assert(markdown.includes('<h2>Heading</h2>'));
assert(markdown.includes('<li>One</li>'));
assert(markdown.includes('data-legacy-diagram="overall"'));

vm.runInContext(accessSource,context,{filename:'documentation-access-v18.6.js'});
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
  const ownerTag=`<script src="documentation-ui.js?v=${releaseIdentity.version}"></script>`;
  const accessTag=`<script src="documentation-access-v18.6.js?v=${releaseIdentity.version}"></script>`;
  assert(index.includes(ownerTag),'index must load capability-oriented Documentation owner with canonical cache identity');
  assert(index.includes(accessTag),'index must retain role-access decorator');
  assert(index.indexOf(ownerTag)<index.indexOf(accessTag),'Documentation owner must load before role-access decorator');
  console.log('Documentation capability owner regression PASS');
}).catch(err=>{console.error(err);process.exitCode=1});
