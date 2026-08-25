'use strict';

let v182AdminIdentity=null;
let v182AdminLoading=false;
let v182AdminError='';

function v182CanAdmin(){return ['SUPER_OWNER','OWNER','ADMIN'].includes(String(authPrincipal?.role||'').toUpperCase())}
function v182Cookie(name){const p=document.cookie.split(';').map(x=>x.trim()).find(x=>x.startsWith(name+'='));return p?decodeURIComponent(p.slice(name.length+1)):''}
async function v182AdminPost(url,body){const r=await fetch(url,{method:'POST',headers:{'Content-Type':'application/json','X-DE-PULSE-CSRF':v182Cookie('depulse_csrf')},body:JSON.stringify(body||{})});const data=await r.json().catch(()=>({}));if(!r.ok)throw new Error(data.error||`Request failed (${r.status})`);return data}
async function v182LoadAdminIdentity(force=false){if(!v182CanAdmin()||v182AdminLoading||(!force&&v182AdminIdentity))return;v182AdminLoading=true;v182AdminError='';try{const r=await fetch('/api/admin/identity',{cache:'no-store'});const data=await r.json().catch(()=>({}));if(!r.ok)throw new Error(data.error||`Admin identity unavailable (${r.status})`);v182AdminIdentity=data}catch(e){v182AdminError=e.message||String(e);v182AdminIdentity=null}finally{v182AdminLoading=false;if(page==='administration')render()}}
function v182Time(ms){if(!Number(ms))return '—';try{return new Date(Number(ms)).toLocaleString()}catch{return '—'}}
function v182Presence(p){const x=String(p||'OFFLINE').toUpperCase(),tone=x==='ACTIVE'?'positive':x==='IDLE'?'warning':'muted';return `<span class="count-pill ${tone}">${esc(x)}</span>`}
function v182AllowedRoles(){const role=String(authPrincipal?.role||'').toUpperCase();if(role==='SUPER_OWNER')return ['OWNER','ADMIN','USER','DEMO'];if(role==='OWNER')return ['ADMIN','USER','DEMO'];if(role==='ADMIN')return ['USER','DEMO'];return []}
function v182RoleOptions(current){const allowed=v182AllowedRoles();if(current&&!allowed.includes(current))allowed.unshift(current);return allowed.map(r=>`<option value="${esc(r)}" ${r===current?'selected':''}>${esc(r)}</option>`).join('')}
function v182AdminMarkup(){
 if(!v182CanAdmin())return '';
 if(!v182AdminIdentity&&!v182AdminLoading&&!v182AdminError)v182LoadAdminIdentity();
 if(v182AdminLoading&&!v182AdminIdentity)return `<section class="settings-section v182-admin"><div class="section-heading-row"><div><span class="eyebrow">Administration</span><h2>Users, Presence & Sessions</h2><p>Loading operational identity state…</p></div></div></section>`;
 if(v182AdminError)return `<section class="settings-section v182-admin"><div class="section-heading-row"><div><span class="eyebrow">Administration</span><h2>Users, Presence & Sessions</h2><p>${esc(v182AdminError)}</p></div><button class="btn" data-v182-admin-refresh>Retry</button></div></section>`;
 const users=v182AdminIdentity?.users||[],sessions=v182AdminIdentity?.sessions||[],roles=v182AllowedRoles();
 const create=roles.length?`<div class="v182-admin-create"><div class="section-heading-row compact"><div><h3>Create User</h3><p>New accounts receive a temporary password and must replace it on first sign-in.</p></div></div><div class="settings-grid-3"><label class="field"><span>Username</span><input data-v182-new-username maxlength="64" autocomplete="off" placeholder="analyst.name"></label><label class="field"><span>Display Name</span><input data-v182-new-display maxlength="80" autocomplete="off" placeholder="Analyst Name"></label><label class="field"><span>Role</span><select data-v182-new-role>${roles.map(r=>`<option value="${esc(r)}">${esc(r)}</option>`).join('')}</select></label></div><div class="settings-grid-2"><label class="field"><span>Temporary Password</span><input data-v182-new-password type="password" autocomplete="new-password" placeholder="Minimum 12 characters"></label><div class="field v182-create-action"><span>Account lifecycle</span><button class="btn primary" data-v182-create-user>Create User</button><small>Password is never returned by the API or persisted in plaintext.</small></div></div></div>`:'';
 const userRows=users.map(u=>{const managed=!!u.manageable,disabled=String(u.status)==='DISABLED';return `<div class="v182-user-row" data-v182-user="${esc(u.id)}"><div class="v182-user-main"><div>${v182Presence(u.presence)}<b>${esc(u.displayName||u.username)}</b><small>@${esc(u.username)} · ${esc(u.role)} · ${esc(u.status)}</small></div><div class="v182-user-meta"><span>${Number(u.activeSessionCount||0)} session${Number(u.activeSessionCount||0)===1?'':'s'}</span><span>Last seen ${esc(v182Time(u.lastSeenAt))}</span></div></div>${managed?`<div class="v182-user-actions"><select data-v182-role="${esc(u.id)}">${v182RoleOptions(String(u.role))}</select><button class="btn small ${disabled?'':'danger ghost'}" data-v182-status="${esc(u.id)}" data-status="${disabled?'ACTIVE':'DISABLED'}">${disabled?'Enable':'Disable'}</button><input data-v182-reset-input="${esc(u.id)}" type="password" autocomplete="new-password" placeholder="Temporary password"><button class="btn small ghost" data-v182-reset="${esc(u.id)}">Reset Password</button></div>`:`<div class="v182-user-actions"><span class="data-badge muted">Protected / current</span></div>`}</div>`}).join('')||'<div class="empty">No users found.</div>';
 const sessionRows=sessions.map(s=>`<div class="compact-row v182-session-row"><b>${v182Presence(s.presence)} ${esc(s.username)}${s.current?' · CURRENT':''}</b><span>Last seen ${esc(v182Time(s.lastSeenAt))} · expires ${esc(v182Time(s.absoluteExpiresAt))}${s.revokable?` <button class="btn tiny ghost danger" data-v182-revoke="${esc(s.id)}">Revoke</button>`:''}</span></div>`).join('')||'<div class="empty">No session history.</div>';
 return `<section class="settings-section v182-admin" id="v182-admin"><div class="section-heading-row"><div><span class="eyebrow">Administration</span><h2>Users, Presence & Sessions</h2><p>Presence is derived from authenticated sessions. Role/status/password changes revoke affected sessions immediately. Credential hashes and session tokens are never exposed.</p></div><button class="btn ghost" data-v182-admin-refresh>Refresh</button></div><div class="summary-grid"><article class="card"><small>Users</small><h3>${users.length}</h3><p>${users.filter(u=>u.presence==='ACTIVE').length} active · ${users.filter(u=>u.status==='DISABLED').length} disabled</p></article><article class="card"><small>Sessions</small><h3>${sessions.filter(s=>s.presence!=='OFFLINE').length}</h3><p>Current valid sessions; expired/revoked history is retained only for bounded operational context.</p></article><article class="card"><small>Presence Source</small><h3>SESSION TRUTH</h3><p>ACTIVE / IDLE / OFFLINE · no separate heartbeat database.</p></article></div>${create}<div class="section-heading-row compact"><div><h3>Users</h3><p>Manage only roles below your own authority. Your own account and higher/equal roles are protected.</p></div></div><div class="v182-user-list">${userRows}</div><div class="section-heading-row compact"><div><h3>Sessions</h3><p>Current session is protected from accidental self-revocation on this surface.</p></div></div><div class="compact-list">${sessionRows}</div></section>`;
}

function v182RenderAdministration(){
 if(!v182CanAdmin())return `<div class="page administration-page access-denied" data-render-owner="administration-ui"><header class="page-head no-regime"><div class="page-head-copy"><h1>Administration</h1><p>This security workflow is not available to your current role.</p></div></header><section class="settings-section"><h2>Access restricted</h2><p>Administration requires an authorized SUPER_OWNER, OWNER, or delegated ADMIN capability. Direct identity/session APIs enforce the same boundary.</p></section></div>`;
 return `<div class="page administration-page" data-render-owner="administration-ui">${pageHead('Administration','Capability-scoped user, role, presence and session administration.','none','','none')}${v182AdminMarkup()}</div>`;
}

const v182BasePageRenderer=pageRenderer;
pageRenderer=function(target){return target==='administration'?v182RenderAdministration:v182BasePageRenderer(target)};

const v182BaseApplyRoleSurfaceVisibility=applyRoleSurfaceVisibility;
applyRoleSurfaceVisibility=function(){
 v182BaseApplyRoleSurfaceVisibility();
 const nav=document.querySelector('[data-page="administration"]');
 if(nav)nav.hidden=!v182CanAdmin();
 if(page==='administration'&&!v182CanAdmin())page='settings';
};

async function v182RefreshAfter(action){try{await action();v182AdminIdentity=null;await v182LoadAdminIdentity(true);toast('Administration updated','Identity/session state refreshed.')}catch(e){toast('Administration action failed',e.message||String(e),'error')}}

document.addEventListener('change',e=>{const el=e.target.closest('[data-v182-role]');if(!el)return;const userId=el.getAttribute('data-v182-role'),role=el.value;v182RefreshAfter(()=>v182AdminPost('/api/admin/users/role',{userId,role}))});
document.addEventListener('click',e=>{
 const refresh=e.target.closest('[data-v182-admin-refresh]');if(refresh){v182AdminIdentity=null;v182LoadAdminIdentity(true);return}
 const create=e.target.closest('[data-v182-create-user]');if(create){const username=document.querySelector('[data-v182-new-username]')?.value||'',displayName=document.querySelector('[data-v182-new-display]')?.value||'',role=document.querySelector('[data-v182-new-role]')?.value||'USER',temporaryPassword=document.querySelector('[data-v182-new-password]')?.value||'';v182RefreshAfter(async()=>{await v182AdminPost('/api/admin/users/create',{username,displayName,role,temporaryPassword});const p=document.querySelector('[data-v182-new-password]');if(p)p.value='' });return}
 const status=e.target.closest('[data-v182-status]');if(status){v182RefreshAfter(()=>v182AdminPost('/api/admin/users/status',{userId:status.getAttribute('data-v182-status'),status:status.getAttribute('data-status')}));return}
 const reset=e.target.closest('[data-v182-reset]');if(reset){const userId=reset.getAttribute('data-v182-reset'),input=document.querySelector(`[data-v182-reset-input="${CSS.escape(userId)}"]`),temporaryPassword=input?.value||'';v182RefreshAfter(async()=>{await v182AdminPost('/api/admin/users/reset-password',{userId,temporaryPassword});if(input)input.value='' });return}
 const revoke=e.target.closest('[data-v182-revoke]');if(revoke){v182RefreshAfter(()=>v182AdminPost('/api/admin/sessions/revoke',{sessionId:revoke.getAttribute('data-v182-revoke')}));return}
});

setInterval(()=>{if(page==='administration'&&v182CanAdmin()&&!document.hidden){v182AdminIdentity=null;v182LoadAdminIdentity(true)}},30000);

const v182Registry=globalThis.__DE_PULSE_RENDERER_OWNERS__||(globalThis.__DE_PULSE_RENDERER_OWNERS__={});
v182Registry.administration={
 owner:'renderer/administration-ui.js',
 version:2,
 state:'ACTIVE_CAPABILITY_SCOPED_OWNER',
 responsibilities:['administration-page','role-scoped-navigation','user-management','session-management','presence-projection'],
 dependencies:['canonical-identity-session-api','capability-authorization','recent-authentication-backend'],
 settingsInjection:false
};
globalThis.__DE_PULSE_ADMINISTRATION_UI__={owner:'renderer/administration-ui.js',version:2,canView:v182CanAdmin,render:v182RenderAdministration,registry:()=>v182Registry.administration};
