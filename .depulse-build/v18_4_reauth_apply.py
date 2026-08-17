#!/usr/bin/env python3
from pathlib import Path

ROOT = Path('.')

def replace_exact(path, old, new, label):
    p = ROOT / path
    text = p.read_text()
    count = text.count(old)
    if count != 1:
        raise SystemExit(f'{label}: expected exactly one match, got {count}')
    p.write_text(text.replace(old, new, 1))

Path('identity_reauth.go').write_text(Path('.depulse-build/identity_reauth.go.payload').read_text())
Path('v18_4_reauth_test.go').write_text(Path('.depulse-build/v18_4_reauth_test.go.payload').read_text())
Path('v18_4_reauth_ui_gate.py').write_text(Path('.depulse-build/v18_4_reauth_ui_gate.py.payload').read_text())

http_auth = Path('http_auth.go')
text = http_auth.read_text()
marker = 'func (a *Application) handleReauthenticate('
if marker in text:
    raise SystemExit('http_auth reauth handlers already present')
text += r'''

func (a *Application) handleReauthenticate(w http.ResponseWriter, r *http.Request) {
	p, ok := principalFromContext(r.Context())
	if !ok || a.identity == nil {
		writeError(w, http.StatusUnauthorized, "Authentication required.")
		return
	}
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&req); err != nil || strings.TrimSpace(req.Password) == "" {
		writeError(w, http.StatusBadRequest, "Current password is required.")
		return
	}
	abuseKey := "reauth:" + loginAbuseKey(p.Username, r)
	if allowed, retryAfter := loginLimiter.Allow(abuseKey); !allowed {
		rejectThrottledLogin(w, retryAfter)
		return
	}
	verified, err := a.identity.reauthenticateSession(p.SessionID, req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Password verification unavailable.")
		return
	}
	if !verified {
		loginLimiter.RecordFailure(abuseKey)
		time.Sleep(120 * time.Millisecond)
		writeError(w, http.StatusForbidden, "Password verification failed.")
		return
	}
	loginLimiter.Reset(abuseKey)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "recentAuthentication": true})
}

func writeReauthRequired(w http.ResponseWriter) {
	writeJSON(w, http.StatusPreconditionRequired, map[string]any{
		"error": "Recent password authentication is required for this security-sensitive change.",
		"code":  "REAUTH_REQUIRED",
	})
}

func (a *Application) requireRecentAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if a.identity == nil {
			next(w, r)
			return
		}
		p, ok := principalFromContext(r.Context())
		if !ok {
			writeError(w, http.StatusUnauthorized, "Authentication required.")
			return
		}
		if !a.identity.sessionRecentlyAuthenticated(p.SessionID, defaultSensitiveReauthTTL) {
			writeReauthRequired(w)
			return
		}
		next(w, r)
	}
}
'''
http_auth.write_text(text)

replace_exact('http_api.go',
    '\tmux.HandleFunc("/api/auth/rotate", a.auth(postOnly(a.handleRotateSession)))\n',
    '\tmux.HandleFunc("/api/auth/rotate", a.auth(postOnly(a.handleRotateSession)))\n\tmux.HandleFunc("/api/auth/reauth", a.auth(postOnly(a.handleReauthenticate)))\n',
    'reauth route')

routes = [
    ('/api/admin/users/create', 'a.requireRole(RoleAdmin, postOnly(a.handleAdminUserCreate))', 'a.requireRole(RoleAdmin, a.requireRecentAuthentication(postOnly(a.handleAdminUserCreate)))'),
    ('/api/admin/users/role', 'a.requireRole(RoleAdmin, postOnly(a.handleAdminUserRole))', 'a.requireRole(RoleAdmin, a.requireRecentAuthentication(postOnly(a.handleAdminUserRole)))'),
    ('/api/admin/users/status', 'a.requireRole(RoleAdmin, postOnly(a.handleAdminUserStatus))', 'a.requireRole(RoleAdmin, a.requireRecentAuthentication(postOnly(a.handleAdminUserStatus)))'),
    ('/api/admin/users/reset-password', 'a.requireRole(RoleAdmin, postOnly(a.handleAdminPasswordReset))', 'a.requireRole(RoleAdmin, a.requireRecentAuthentication(postOnly(a.handleAdminPasswordReset)))'),
    ('/api/admin/sessions/revoke', 'a.requireRole(RoleAdmin, postOnly(a.handleAdminSessionRevoke))', 'a.requireRole(RoleAdmin, a.requireRecentAuthentication(postOnly(a.handleAdminSessionRevoke)))'),
    ('/api/settings/save', 'a.requireRole(RoleAdmin, postOnly(a.handleSettingsSave))', 'a.requireRole(RoleAdmin, a.requireRecentAuthentication(postOnly(a.handleSettingsSave)))'),
    ('/api/settings/clear-secret', 'a.requireRole(RoleAdmin, postOnly(a.handleClearSecret))', 'a.requireRole(RoleAdmin, a.requireRecentAuthentication(postOnly(a.handleClearSecret)))'),
]
for route, oldwrap, newwrap in routes:
    old = f'\tmux.HandleFunc("{route}", {oldwrap})\n'
    new = f'\tmux.HandleFunc("{route}", {newwrap})\n'
    replace_exact('http_api.go', old, new, 'recent-auth route '+route)

renderer = Path('renderer/renderer.js')
js = renderer.read_text()
old = "function cookieValue(name){const hit=document.cookie.split('; ').find(x=>x.startsWith(name+'='));return hit?decodeURIComponent(hit.slice(name.length+1)):''}async function api(path,body){const opt=body===undefined?{}:{method:'POST',headers:{'Content-Type':'application/json','X-DE-PULSE-CSRF':cookieValue('depulse_csrf')},body:JSON.stringify(body)};const r=await fetch(path,opt);const data=await r.json().catch(()=>({}));if(r.status===401||r.status===428){location.replace('/');throw new Error(data.error||'Authentication required.')}if(!r.ok)throw new Error(data.error||`Request failed (${r.status})`);return data}"
new = r'''function cookieValue(name){const hit=document.cookie.split('; ').find(x=>x.startsWith(name+'='));return hit?decodeURIComponent(hit.slice(name.length+1)):''}
let __reauthPromise=null;
function requestRecentAuthentication(){if(__reauthPromise)return __reauthPromise;__reauthPromise=new Promise(resolve=>{let style=document.getElementById('depulse-reauth-style');if(!style){style=document.createElement('style');style.id='depulse-reauth-style';style.textContent='.reauth-overlay{position:fixed;inset:0;z-index:10000;display:grid;place-items:center;background:rgba(4,8,14,.72);backdrop-filter:blur(5px)}.reauth-card{width:min(420px,calc(100vw - 32px));background:#111820;border:1px solid rgba(255,255,255,.14);border-radius:14px;padding:22px;box-shadow:0 24px 64px rgba(0,0,0,.5)}.reauth-card h2{margin:0 0 8px;font-size:18px}.reauth-card p{margin:0 0 16px;opacity:.8;line-height:1.45}.reauth-card label{display:block;margin:0 0 7px;font-size:12px;font-weight:700}.reauth-card input{box-sizing:border-box;width:100%;padding:11px 12px;border-radius:8px;border:1px solid rgba(255,255,255,.18);background:#0b1118;color:inherit}.reauth-error{min-height:20px;margin-top:8px;font-size:12px}.reauth-actions{display:flex;justify-content:flex-end;gap:8px;margin-top:16px}.reauth-actions button{min-width:92px;padding:9px 13px;border-radius:8px;border:1px solid rgba(255,255,255,.16);cursor:pointer}.reauth-actions button[type=submit]{font-weight:700}';document.head.appendChild(style)}const overlay=document.createElement('div');overlay.className='reauth-overlay';overlay.setAttribute('role','presentation');overlay.innerHTML='<form class="reauth-card" role="dialog" aria-modal="true" aria-labelledby="reauth-title"><h2 id="reauth-title">Confirm password</h2><p>For this security-sensitive change, enter your current password.</p><label for="reauth-password">Current password</label><input id="reauth-password" type="password" autocomplete="current-password" required><div class="reauth-error" role="alert" aria-live="polite"></div><div class="reauth-actions"><button type="button" data-reauth-cancel>Cancel</button><button type="submit">Continue</button></div></form>';document.body.appendChild(overlay);const form=overlay.querySelector('form'),input=overlay.querySelector('#reauth-password'),error=overlay.querySelector('.reauth-error'),submit=form.querySelector('button[type=submit]');let settled=false;const finish=value=>{if(settled)return;settled=true;input.value='';overlay.remove();resolve(value)};overlay.querySelector('[data-reauth-cancel]').addEventListener('click',()=>finish(false));overlay.addEventListener('click',e=>{if(e.target===overlay)finish(false)});overlay.addEventListener('keydown',e=>{if(e.key==='Escape'){e.preventDefault();finish(false)}});form.addEventListener('submit',async e=>{e.preventDefault();const password=input.value;if(!password)return;submit.disabled=true;error.textContent='';try{const r=await fetch('/api/auth/reauth',{method:'POST',headers:{'Content-Type':'application/json','X-DE-PULSE-CSRF':cookieValue('depulse_csrf')},body:JSON.stringify({password})});const data=await r.json().catch(()=>({}));input.value='';if(r.ok){finish(true);return}if(r.status===403||r.status===429){error.textContent=data.error||'Password verification failed.';input.focus();return}if(r.status===401){location.replace('/');finish(false);return}error.textContent=data.error||`Verification failed (${r.status}).`;input.focus()}catch(_){error.textContent='Password verification is temporarily unavailable.';input.focus()}finally{submit.disabled=false}});requestAnimationFrame(()=>input.focus())});return __reauthPromise.finally(()=>{__reauthPromise=null})}
async function api(path,body){const send=async()=>{const opt=body===undefined?{}:{method:'POST',headers:{'Content-Type':'application/json','X-DE-PULSE-CSRF':cookieValue('depulse_csrf')},body:JSON.stringify(body)};const r=await fetch(path,opt);const data=await r.json().catch(()=>({}));return {r,data}};let {r,data}=await send();if(r.status===428&&data.code==='REAUTH_REQUIRED'){if(await requestRecentAuthentication())({r,data}=await send());else throw new Error(data.error||'Recent password authentication required.')}if(r.status===401||r.status===428){location.replace('/');throw new Error(data.error||'Authentication required.')}if(!r.ok)throw new Error(data.error||`Request failed (${r.status})`);return data}'''
count = js.count(old)
if count != 1:
    raise SystemExit(f'renderer api helper: expected exactly one match, got {count}')
renderer.write_text(js.replace(old, new, 1))

replace_exact('.github/workflows/v18.4-dev-ci.yml',
'''      - name: v18.4 focused security tests
        run: go test -count=1 -run '^TestV184' ./...
''',
'''      - name: v18.4 focused security tests
        run: |
          python3 v18_4_reauth_ui_gate.py
          node --check renderer/renderer.js
          go test -count=1 -run '^TestV184' ./...
''',
'v18.4 CI focused security step')
