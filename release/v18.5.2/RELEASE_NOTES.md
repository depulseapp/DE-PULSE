# DE.PULSE v18.5.2 Hotfix Candidate

Target build: `v18.5.2-stable-hotfix-20260817`  
Previous Stable: `v18.5.1`  
Status: **candidate — not promoted**

## Recovered behavior

- Day, Swing, and Long-Term desks no longer depend on an undeclared global symbol normalizer.
- ET/PT clocks are preserved at full readability inside the compact Market Pulse Ribbon, ordered session → coverage → clocks → data control.
- Research uses the full page hierarchy and removes duplicate top-level AI actions.
- Live quote reconciliation preserves hover, focus, selection, scroll, and tracked-symbol drafts.
- Ticker input, Add Symbol, and Remove All align beside the Tracked Symbols heading on desktop and stack responsively.
- Display name and sign-in username are configurable with recent-password verification; OWNER remains a separate role.
- The complete Save Settings row stays inside the visible application viewport, including the reported 402 px window height.

## Preserved boundaries

- Deterministic Day/Swing/Long formulas are unchanged.
- Actionable instruments remain U.S.-listed only.
- Provider intelligence and persistence owners are unchanged.
- No Execution remains permanent.

## Certification

Automatic GitHub Actions triggers are paused because the Actions budget is exhausted. Use `run_free_certification.sh` from this folder for the exact-source local lane. Stable promotion remains blocked until that lane and native macOS/Windows G13/G14 package/runtime audits pass.
