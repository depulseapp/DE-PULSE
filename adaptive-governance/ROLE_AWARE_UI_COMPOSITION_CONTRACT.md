# DE.PULSE — Permanent Role-Aware UI Composition Contract

Status: **Permanent governing contract**  
Applies to: Adaptive Roadmap, Adaptive Build Plan, Adaptive Build Process, Adaptive Delivery Process  
Effective: **v18.2 and all later releases**

## Purpose

DE.PULSE must provide an intentional, seamless experience for each authorized role. Role-aware UX is not implemented by simply hiding controls. Every tab is composed for the authenticated role while preserving product hierarchy, information value, responsive fit, and security boundaries.

This contract operates inside canonical G0–G16, primarily G1/G2/G3/G9/G10/G12/G16. It creates no additional gate.

## Role classes

- **SUPER_OWNER / OWNER** — full governance and authorized diagnostics, including owner-only security/global controls.
- **ADMIN** — delegated operational administration only. Admin visibility/actions are capability-scoped and must not automatically inherit owner-only security, secrets, destructive controls, or governance authority.
- **USER** — market intelligence, personal symbols/workspaces, research and decision-support surfaces required for normal use; internal implementation machinery is suppressed.
- **DEMO** — USER-style composition with truthful demo/data limitations and no privileged operational controls.

UI visibility is never authorization. Backend/API authorization remains mandatory even when a control is absent from the rendered page.

## Hierarchy invariant

**Role changes composition, not information hierarchy.**

Primary market intelligence, decisions, material changes, risk, contradictions and user workflow remain above supporting context. A privileged diagnostic/control must not move to the top merely because the viewer is OWNER or ADMIN. A USER losing privileged content must not cause a lower-value card to become artificially prominent unless its product priority independently justifies that position.

Each tab follows this order where applicable:

1. primary decision/intelligence;
2. material risk/change/context;
3. user workflow/actions;
4. supporting evidence/context;
5. role-specific administration/operations;
6. deep diagnostics/governance/drill-down.

Owner/admin additions are inserted where they are contextually appropriate, not automatically prepended.

## Seamless composition and no-gap rule

Role-restricted content must be removed from layout composition, not merely made invisible while retaining reserved geometry.

Mandatory behavior:

- no `visibility:hidden`/transparent placeholders that preserve empty card/grid space for unauthorized content;
- no fixed-height/min-height shells whose only purpose is to preserve the OWNER layout for USER/DEMO;
- when a card is absent, neighboring content must reflow naturally using the canonical responsive grid/flex layout;
- when every child in a role-specific section is absent, remove the section heading/container/divider as well;
- no orphan separators, empty columns, blank cards, unexplained whitespace, uneven two-column leftovers, or awkward page endings;
- role-specific compositions must remain balanced at all supported desktop/tablet/narrow-browser widths;
- role changes must not create clipping, overlap, horizontal page scroll, unreadable compression, or focus/tab-order gaps.

## Canonical tab audit

G9 audits every current application tab for OWNER, delegated ADMIN, USER and DEMO composition:

1. Dashboard
2. Market Intelligence
3. Day Trade Desk
4. Swing Desk
5. Long-Term Desk
6. Discovery
7. Research
8. AI Copilot
9. Maintenance
10. Settings
11. Documentation

The navigation source of truth is the current renderer navigation; new tabs automatically enter this contract and cannot ship without role-composition ownership.

## Tab-level role direction

### Dashboard
All roles receive the same primary decision/intelligence hierarchy. Internal provider/runtime diagnostics do not enter the USER/DEMO composition. Owner/admin operational context, when useful, remains supporting/drill-down and never displaces Decision Queue, material market context or catalyst risk.

### Market Intelligence
All roles receive synthesized regime/driver/risk intelligence. Raw provider routing, circuits, queues, cache/database and implementation details remain privileged diagnostics and are not promoted above market conclusions.

### Day / Swing / Long-Term Desks
All roles receive the same protected deterministic Action/Score logic and horizon intelligence. Role never changes scoring formulas. Operational diagnostics are contextual/drill-down only. USER/DEMO retain personal workspace actions allowed by v18.1.

### Discovery
All roles retain candidate discovery, staging and Research flow. Provider/engine machinery stays privileged. Role-specific additions cannot outrank candidate qualification and action flow.

### Research
All roles retain the decision brief, evidence, freshness/materiality and Research workflow appropriate to their data access. Raw implementation/provider machinery is privileged. User-facing evidence labels remain truthful without exposing secrets or operational internals.

### AI Copilot
Normal decision-support use remains user-facing when enabled. Provider credentials, routing configuration, policy promotion and system-level AI operations belong in privileged Settings/Maintenance, not at the top of AI Copilot merely because the viewer is privileged.

### Maintenance
SUPER_OWNER/OWNER receive full authorized diagnostics. ADMIN receives only explicitly delegated operational controls/diagnostics. USER/DEMO do not receive the Maintenance navigation/page. Direct unauthorized navigation must be rejected/redirected rather than rendering privileged markup and relying on CSS hiding.

### Settings
Settings is composition-based. USER/DEMO see personal/account/user-facing preferences only. ADMIN receives delegated administration in a clearly placed Administration section. OWNER/SUPER_OWNER receive global/security/provider/runtime controls in their appropriate sections. Privileged sections do not automatically move above personal/general settings simply because they exist.

### Documentation
USER/DEMO default to user/capability documentation. Developer/operational/governance documentation is role-aware where sensitive/internal. Documentation navigation must not expose privileged operational detail merely because the underlying files are bundled.

## Navigation contract

Unauthorized tabs/controls are removed from navigation flow. Remaining groups/buttons reflow without blank rows or orphan headings. Direct URL/programmatic page selection is checked against the same role/capability policy so navigation hiding cannot be bypassed.

## Security + composition contract

Sensitive or unauthorized data should be filtered server-side whenever practical. Do not send secrets, token hashes, provider credentials, owner-only security state or privileged administrative payloads to unauthorized roles and rely on client hiding.

## G9 role-aware layout audit

G9 must exercise each audited tab under at least:

- OWNER-equivalent composition;
- delegated ADMIN composition;
- USER composition;
- DEMO composition;

For each supported viewport family, verify:

- correct allowed/forbidden surfaces;
- preserved top-level information hierarchy;
- no empty reserved geometry after suppression;
- no orphan headers/dividers/containers;
- no page-level overflow, overlap or clipping;
- sensible grid/card reflow;
- correct keyboard/focus order for remaining controls;
- direct navigation to unauthorized pages is rejected/redirected;
- role changes do not alter protected deterministic market logic.

## G16 closure

G16 records the role/tab audit result and any new UI-composition lesson. Any recurring role/layout failure is generalized into a regression or reusable composition primitive rather than fixed tab-by-tab forever.

## Permanent rule

**Do not design the OWNER page and subtract cards for everyone else. Compose each role from the same canonical information hierarchy, authorized capabilities and reusable layout primitives so every role receives a complete, intentional page.**
