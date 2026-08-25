'use strict';

const fs=require('fs');
const vm=require('vm');
const assert=require('assert');
const crypto=require('crypto');

const ownerPath='renderer/administration-ui.js';
const ownerSource=fs.readFileSync(ownerPath,'utf8');
const index=fs.readFileSync('renderer/index.html','utf8');
const nav={hidden:true};
const noop=()=>{};
const document={
  cookie:'depulse_csrf=test-token',
  hidden:false,
  querySelector(selector){return selector==='[data-page="administration"]'?nav:null},
  addEventListener:noop
};
const context=vm.createContext({
  console,
  globalThis:null,
  document,
  CSS:{escape:v=>String(v)},
  setInterval:noop,
  decodeURIComponent,
  fetch:async()=>({ok:true,json:async()=>({users:[],sessions:[]})})
});
context.globalThis=context;

function gitBlobToken(path){
  const data=fs.readFileSync(path);
  const header=Buffer.from(`blob ${data.length}\0`,'utf8');
  return crypto.createHash('sha1').update(header).update(data).digest('hex').slice(0,16);
}

vm.runInContext(`
let authPrincipal={role:'ADMIN'};
let page='settings';
let renderCalls=0;
const esc=v=>String(v??'').replace(/[&<>"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]));
function pageHead(title,sub){return '<header><h1>'+esc(title)+'</h1><p>'+esc(sub)+'</p></header>'}
function render(){renderCalls++}
function toast(){}
function renderSettings(){return '<div class="page settings"><h1>Settings</h1><section class="application-settings">Application</section></div>'}
function pageRenderer(target){return target==='settings'?renderSettings:()=>'<div>'+esc(target)+'</div>'}
function applyRoleSurfaceVisibility(){}
`,context);

vm.runInContext(ownerSource,context,{filename:ownerPath});
const owner=vm.runInContext('__DE_PULSE_ADMINISTRATION_UI__',context);
assert.strictEqual(owner.owner,'renderer/administration-ui.js');
assert.strictEqual(owner.registry().state,'ACTIVE_CAPABILITY_SCOPED_OWNER');
assert.strictEqual(owner.registry().settingsInjection,false,'Administration must not be injected into Settings');

vm.runInContext('v182AdminIdentity={users:[],sessions:[]}',context);
const adminHTML=vm.runInContext("pageRenderer('administration')()",context);
assert(adminHTML.includes('administration-page'),'authorized Administration must have a dedicated page');
assert(adminHTML.includes('Users, Presence & Sessions'),'dedicated Administration page must retain the canonical identity/session workflow');
assert(!vm.runInContext('renderSettings()',context).includes('v182-admin'),'Settings must no longer own Administration markup');

vm.runInContext('applyRoleSurfaceVisibility()',context);
assert.strictEqual(nav.hidden,false,'ADMIN must see the Administration navigation destination');
for(const role of ['OWNER','SUPER_OWNER']){
  vm.runInContext(`authPrincipal={role:'${role}'};applyRoleSurfaceVisibility()`,context);
  assert.strictEqual(nav.hidden,false,`${role} must see Administration navigation`);
  assert(vm.runInContext("pageRenderer('administration')()",context).includes('administration-page'),`${role} must render Administration`);
}
for(const role of ['USER','DEMO']){
  vm.runInContext(`authPrincipal={role:'${role}'};page='settings';applyRoleSurfaceVisibility()`,context);
  assert.strictEqual(nav.hidden,true,`${role} must not see Administration navigation`);
  const denied=vm.runInContext("pageRenderer('administration')()",context);
  assert(denied.includes('Access restricted'),`${role} direct renderer access must fail closed`);
}
vm.runInContext("authPrincipal={role:'USER'};page='administration';applyRoleSurfaceVisibility()",context);
assert.strictEqual(vm.runInContext('page',context),'settings','role loss while on Administration must leave the restricted page');

const navMarkup=index.match(/<button class="nav" data-page="administration"[\s\S]*?<\/button>/)?.[0]||'';
assert(navMarkup.includes('hidden'),'Administration nav must default hidden until authenticated role projection runs');
const ownerTag=`<script src="administration-ui.js?v=${gitBlobToken(ownerPath)}"></script>`;
assert(index.includes(ownerTag),'index must load Administration owner with content-derived cache identity');
assert(!ownerSource.includes('const v182BaseRenderSettings=renderSettings'),'legacy Settings injection wrapper must be removed');
assert(ownerSource.includes("page==='administration'"),'Administration refresh lifecycle must be scoped to its own page');

console.log('Administration capability-scoped surface regression PASS');
