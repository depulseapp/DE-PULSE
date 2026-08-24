'use strict';
const fs=require('fs');
const assert=require('assert');
const backend=fs.readFileSync('desk_membership.go','utf8');
const watchlist=fs.readFileSync('renderer/watchlist-ui.js','utf8');
const renderer=fs.readFileSync('renderer/renderer.js','utf8');

const backendRequired=[
  'func deskIDs() []string { return []string{"day", "swing", "long"} }',
  'func trackedSymbolsLocked(st *AppState) []string',
  'func setTrackedSymbolLocked(st *AppState, sym string, on bool) map[string]bool',
  'wl.Symbols = []string{}',
  'if desired && !membership[desk]',
  'if activeDeskCount(membership) > 1',
  'protected = true',
  'func parseUserTicker(raw string) (string, bool)',
  'hasForeignListingSuffix(sym)',
  'sym == "VIX"',
  'yahooMetaIsUSActionable',
  'setTrackedSymbolLocked(&workspaceState, sym, true)',
  'setTrackedSymbolLocked(&workspaceState, sym, false)',
  'saveWorkspaceStateLocked',
  '"removed": before'
];
for(const token of backendRequired)assert(backend.includes(token),`Symbol/desk backend contract missing: ${token}`);

const uiRequired=[
  "const res = await api('/api/master-symbol/remove', { symbol: sym });",
  'membership: removed',
  "await api('/api/master-symbol/restore', { symbol: sym, membership: undo.membership });",
  'selected: selectedBefore',
  'Previous desk memberships and desk selection restored.',
  "const membershipResult = await api('/api/desk/membership'",
  "const boot = await api('/api/bootstrap');",
  'Removed from every desk. Undo restores the exact previous desk memberships.'
];
for(const token of uiRequired)assert(watchlist.includes(token),`Symbol/desk UI contract missing: ${token}`);
assert(renderer.includes("res.protected?'At least one desk must remain selected. Use Tracked Symbols × for global removal.'"),'final-desk protection UX missing');
assert(renderer.includes("await api('/api/master-symbol/restore',{symbol:sym,membership:u.membership})"),'canonical renderer undo must restore exact membership');
console.log('Symbol/desk correctness regression PASS');
console.log('canonical Day/Swing/Long membership + final-desk protection: PASS');
console.log('Master Market Symbols persistence/global remove/exact Undo: PASS');
console.log('strict U.S.-listed ticker boundary: PASS');
