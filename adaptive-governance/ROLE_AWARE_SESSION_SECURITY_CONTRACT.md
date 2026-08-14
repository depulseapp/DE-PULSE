# DE.PULSE — Role-Aware Session Security Contract

Status: PERMANENT
Effective: v18.2 onward
Canonical name: **Role-Aware Session Security**

## 1. Purpose

DE.PULSE uses persistent authentication for normal weekday use while applying role-aware session security. Weekend reset is a normal-user session-lifecycle control, not a blanket privileged-access shutdown.

The governing principle is:

**Persistent convenience for normal operation → deterministic weekend reset for USER/DEMO → stronger continuous controls for privileged roles.**

Session policy is enforced server-side. UI logout behavior alone is never sufficient security.

## 2. Canonical session profiles

### STANDARD_USER_SESSION

Applies by default to:
- `USER`
- `DEMO`

Requirements:
- persistent sign-in is permitted during normal weekday operation subject to normal expiry/revocation controls;
- weekend reset is mandatory;
- any active USER/DEMO session presented during the weekend window is invalidated server-side;
- after weekend invalidation, fresh authentication is required before access resumes;
- client-side cached authentication state must not restore a revoked/invalid weekend session;
- direct API/SSE/browser/native requests must receive the same invalid-session result.

### PRIVILEGED_SESSION

Applies by default to:
- `SUPER_OWNER`
- `OWNER`
- `ADMIN`, including both full-capability and limited/delegated ADMIN profiles

Requirements:
- the weekend-reset rule is disabled by default so legitimate weekend maintenance, security, release, provider and administrative work remains possible;
- privileged sessions remain subject to configured inactivity timeout and maximum session lifetime;
- sensitive actions require fresh authentication/re-authentication when the applicable security policy requires it;
- session creation, renewal, revocation, privilege/capability changes and sensitive administrative activity are auditable;
- revocation must take effect server-side for API, SSE and UI/native clients;
- `ADMIN` weekend exemption does not grant capabilities. Authorization remains capability-based.

## 3. Weekend definition

The canonical weekend policy timezone is **`America/New_York`**, aligned to the primary U.S. equities operating context.

Default weekend window:

**Saturday 00:00:00 America/New_York through Monday 00:00:00 America/New_York.**

Rules:
- USER/DEMO sessions are not valid inside this window;
- boundary evaluation must be server-side and independent of the client machine timezone;
- DST transitions use the canonical timezone database rather than fixed UTC offsets;
- weekends are distinct from exchange holidays. A market holiday does not automatically become a weekend session reset unless separately governed;
- the server should proactively revoke applicable sessions at the boundary where scheduling is available, and must also enforce the rule on every authenticated request/renewal so scheduler failure cannot bypass it.

## 4. Role/capability changes

Session policy is re-evaluated immediately when a user's role/status/capabilities materially change.

Examples:
- USER promoted to ADMIN: future authorization follows the delegated ADMIN capability set and the privileged session profile;
- ADMIN demoted to USER during a weekend: the session becomes invalid immediately;
- disabled/suspended account: all applicable sessions are revoked regardless of role or day;
- capability removal does not require waiting for session expiry; authorization reflects current canonical capability truth.

No stale session claim may preserve permissions that the canonical identity state has revoked.

## 5. Revocation and emergency controls

- `SUPER_OWNER` / authorized `OWNER` governance retains global force-sign-out capability.
- Individual session/user revocation is supported.
- ADMIN revocation authority is capability-scoped; `ADMIN` alone does not imply global session control.
- Security-driven revocation overrides persistence, weekend exemption and normal expiry.
- Revocation should propagate to long-lived SSE/streaming connections and not only the next page load.

## 6. Sensitive-action re-authentication

Privileged access remains available on weekends, but sensitive operations may require recent authentication according to the security policy. Candidate operations include:
- ownership transfer/security recovery changes;
- credential/provider-secret changes;
- role/capability elevation;
- global session revocation;
- security-policy changes;
- other high-impact administrative controls identified by the capability/security registry.

The exact inactivity/max-lifetime/re-authentication durations are configuration/security-policy values and must be centrally owned, testable and not duplicated in individual pages.

## 7. Cross-platform consistency

The same canonical server-side session truth applies to:
- macOS native delivery;
- Windows native delivery;
- hosted web/shared-server delivery.

A client may display a friendly reason such as `Weekend sign-in reset`, but it must not independently decide whether a session is valid.

## 8. Required evidence

Affected releases must test at minimum:
- USER weekday persistent session remains valid when otherwise eligible;
- USER/DEMO session is invalidated at Saturday boundary;
- USER/DEMO cannot restore access from stale local/client state during weekend;
- fresh authentication is required after weekend reset;
- OWNER/SUPER_OWNER weekend access remains available subject to privileged controls;
- full-capability and limited ADMIN weekend access remains available but authorization is still capability-scoped;
- ADMIN-to-USER demotion during weekend invalidates the session;
- disabled account revokes regardless of weekend/role;
- global and individual revocation work for API/UI/SSE;
- America/New_York DST/timezone behavior is deterministic;
- macOS, Windows and hosted/server behavior is consistent.

## 9. Gate placement

This contract does not require a new canonical G-gate by itself.

- **G1:** freeze affected roles/session policies and release disposition.
- **G2:** define canonical identity/session-policy owner, timezone and server-side enforcement boundary.
- **G3:** define session state transitions, capability interactions, revocation and acceptance matrix.
- **G4:** implement server-side policy and clients that consume canonical session truth.
- **G7:** prove security, capability, revocation, audit and sensitive-data boundaries.
- **G9:** verify user-facing sign-out/login flow is seamless and role-appropriate.
- **G10:** block freeze on missing weekend/privileged session evidence for affected scope.
- **G12:** replay on immutable RC.
- **G14:** validate actual native artifacts where session behavior is platform-relevant.
- **G15:** require session-policy certification before promotion when affected.
- **G16:** record incidents, revocation/session-policy learnings and preventative regressions.

## 10. v18.2 disposition

Because v18.2 owns Admin / Presence / Sessions, Role-Aware Session Security is a `CURRENT_RELEASE_BLOCKER` for the applicable identity/session scope. Documentation alone does not close it. Stable promotion requires implemented and evidenced server-side weekend reset for USER/DEMO plus privileged-role exemption/security behavior defined by this contract.
