'use strict';
const fs=require('fs');
function need(ok,msg){if(!ok)throw new Error(msg)}
const js=fs.readFileSync('renderer/admin-v18.2.js','utf8');
const css=fs.readFileSync('renderer/admin-v18.2.css','utf8');
const html=fs.readFileSync('renderer/index.html','utf8');

need(html.includes('admin-v18.2.css?v=18.2.0'),'v18.2 admin CSS is not loaded');
need(html.includes('admin-v18.2.js?v=18.2.0'),'v18.2 admin JS is not loaded');
// Release identity/version consistency gates own the current title. This inherited
// v18.2 regression harness verifies only that a canonical DE.PULSE release title remains.
need(/<title>DE\.PULSE v\d+\.\d+\.\d+<\/title>/.test(html),'canonical renderer title missing');
need(js.includes("['SUPER_OWNER','OWNER','ADMIN']"),'admin UI role allowlist drift');
need(js.includes("return '';")&&js.includes('if(!v182CanAdmin())'),'normal-user admin suppression missing');
need(js.includes("'X-DE-PULSE-CSRF':v182Cookie('depulse_csrf')"),'admin mutation CSRF protection missing');
need(js.includes("Credential hashes and session tokens are never exposed"),'credential-redaction explanation missing');
need(js.includes('SESSION TRUTH'),'presence source is not explained as canonical session truth');
need(js.includes("/api/admin/users/create")&&js.includes("/api/admin/users/role")&&js.includes("/api/admin/users/status")&&js.includes("/api/admin/users/reset-password")&&js.includes("/api/admin/sessions/revoke"),'admin lifecycle actions incomplete');
need(!/passwordHash|tokenHash/.test(js),'admin renderer references secret credential fields');
need(css.includes('@media(max-width:760px)')&&css.includes('@media(max-width:480px)'),'narrow responsive rules missing');
need(css.includes('.v182-user-actions{grid-template-columns:1fr 1fr}')&&css.includes('.v182-user-actions{grid-template-columns:1fr}'),'admin action layout does not progressively reflow');
console.log('v18.2 renderer admin / role visibility / responsive acceptance: PASS');
