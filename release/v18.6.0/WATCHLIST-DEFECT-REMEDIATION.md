# v18.6 Watchlist Defect Remediation

Status: **IMPLEMENTED — qualification evidence bound; release qualification still required**

This note records user-reported watchlist defects repaired during v18.6 development without expanding the immutable G1 parent scope. The fixes preserve canonical per-user workspace state and existing No Execution / U.S. Equities boundaries.

## DEF-18.6-WL-001 — Deprecated CURRENT desk marker

**Observed:** symbol rows displayed `DAY / SWING / CURRENT / LONG`, making the current page look like a second membership concept.

**Required behavior:** Day, Swing and Long are the membership controls. A button is pressed/selected if and only if the symbol belongs to that desk. There is no separate CURRENT badge, `aria-current`, or `current-desk` styling.

**Implementation:** `renderer/watchlist-v18.5.1.js` derives each button's `active` class and `aria-pressed` value directly from canonical desk membership. `renderer/watchlist-v18.5.1.css` styles only active/`aria-pressed=true` membership state.

## DEF-18.6-WL-002 — Desk symbol add regression

**Observed:** add-symbol actions from trading-desk surfaces could fail to survive the subsequent per-user bootstrap/reconciliation.

**Root cause:** desk add, table add, and AI-to-desk add paths still called the legacy `/api/watchlists/add-symbol` mutation after desk membership had moved to the canonical per-user `/api/desk/membership` state.

**Implementation:** all three desk-add surfaces now call one `addSymbolToDesk()` path using `/api/desk/membership` with `active:true`, re-bootstrap canonical state, verify durable membership, and only then clear the input/draft and report success. Discovery remains a distinct non-desk watchlist and retains its own watchlist mutation path in the base renderer.

## DEF-18.6-WL-003 — Qualification expected obsolete behavior

**Observed:** the qualified browser lane still required `CURRENT`, `aria-current`, and `current-desk`, so a correct product fix would fail certification.

**Implementation:** historical v18.5.1 release proof remains untouched. v18.6 adds `release/v18.6.0/browser_watchlist_membership_test.py`, and the canonical qualified browser workflow now targets that current proof.

## Regression protection

- `tools/ci/watchlist_membership_contract.py` fails if the v18.6 override reintroduces `CURRENT`, `current-desk`, or legacy desk-add routing.
- CI Fast runs the watchlist membership regression contract on every development-branch push.
- CI Qualified browser lane is bound to the v18.6 browser behavior proof.
- Existing global tracked-symbol removal, exact Undo restoration, last-desk protection, and prior desk selection restoration remain covered.

## Fresh development proof performed during remediation

- JavaScript syntax check for the production watchlist extension: **PASS**.
- Static membership/add routing contract: **PASS**.
- Headless Chromium: Swing-only symbol renders `Swing` pressed, `Day`/`Long` unpressed, no CURRENT marker; desk, table and AI adds use only `/api/desk/membership`: **PASS**.
- Headless Chromium: all 7 legal Day/Swing/Long membership combinations, global remove + exact Undo, pressed-state-only rendering, and final-desk protection: **PASS**.

These development proofs do not substitute for later G5+ / G12+ qualification, packaging, or actual-artifact certification.

## Documentation impact disposition

- **User documentation:** behavior correction only; no new feature or workflow concept requiring a new user-doc section. Existing wording must not describe a CURRENT desk marker.
- **Developer/certification documentation:** impacted; this remediation note, CI Fast contract, and v18.6 browser proof are the required current evidence.
- **Adaptive process:** no new top-level gate; evidence is consumed by the existing G0–G16 process.
