# DE.PULSE v16.10 — Information Architecture & Data Efficiency Pass

**Status:** REQUIRED FIRST ACTIVITY · COMPLETE BEFORE v16.10 FEATURE FREEZE  
**Baseline:** v16.9.0 Stable · original professional roadmap 30 FULL / 0 PARTIAL / 0 MISSING

## Surface decisions

| Surface | Decision | v16.10 placement |
|---|---|---|
| Dashboard | KEEP | Remains attention summary. Opportunity Radar appears only through material Smart Notifications; no duplicate scanner table. |
| Market Intelligence | KEEP / REUSE MORE | Existing regime/event/liquidity/global owners remain canonical. Radar may consume their context; it does not create a competing regime signal. |
| Discovery | REUSE MORE / MERGE | Single canonical Scanner owner. Add Always-On Opportunity Radar, unusual volume/volatility, rapid dislocation, promotion/demotion here. Manual horizon scan remains available. |
| Research | KEEP | Radar candidates open existing Research. No second research surface. |
| Day / Swing / Long | KEEP | Deterministic Action/Score remains protected. Radar never auto-adds a candidate to a desk. |
| Decision Queue | KEEP | No automatic insertion from Radar in v16.10; material state remains handled by existing readiness/event owners. |
| Maintenance | REUSE MORE | Show adaptive cadence, hot-symbol set, provider state, cache policy and Shadow experiments. No trading decision logic lives here. |
| Settings | NO CHANGE REQUIRED | No new user-facing tuning knobs until Shadow/validation proves a need. |
| Documentation | UPDATE AFFECTED MODULES | Document Radar semantics, adaptive data policy and Shadow boundary only. |

## Data / processing decisions

| Data / owner | Decision | Reason |
|---|---|---|
| Alpaca U.S. asset universe | USE / CACHE | Existing canonical Discovery universe. Cache for 12h so always-on Radar does not refetch the asset list every cycle. |
| Alpaca stock snapshots | REUSE MORE | Broad low-cost observation layer. Batch <=50 symbols; rotate a bounded subset instead of permanent full-universe streaming. |
| Alpaca market activity | REUSE MORE | Most Active / Gainers / Losers become dynamic Radar seeds instead of UI-only information. |
| Existing Discovery score | KEEP | Manual horizon ranking remains intact. Opportunity Score is a separate Discovery context field, not a protected desk score. |
| Relative volume | REUSE / IMPROVE FOR RADAR | Radar adds regular-session-normalized RVOL; manual legacy fields remain compatible. |
| Daily bar / previous daily bar | REUSE MORE | Range expansion and session participation baseline; no second history store. |
| Quote spread / dollar volume | REUSE MORE | Liquidity gates prevent noisy low-quality promotions. |
| News / Catalyst Watch | REUSE MORE | Material recent events can strengthen Radar context. |
| Community Evidence Fusion | REUSE MORE | Recent HIGH/ELEVATED fused community evidence is contextual Radar evidence only; still UNTRUSTED and no deterministic mutation. |
| Provider health | REUSE MORE | Radar cadence slows when provider health degrades. |
| livePriorityHints | REUSE MORE | Radar promotions reuse existing bounded live-reserve allocation. No new subscription owner. |
| Intraday history | REUSE MORE | Newly promoted symbols get immediate targeted hydration; hot symbols can refresh between normal full cycles. |
| Cache persistence | ADAPT | 1m hot active-session / 2m normal active / 5m overnight / 10m closed, bounded by existing asynchronous cache owner. |
| Shadow experiments | ADD | Read-only evaluation metadata with explicit SHADOW → VALIDATED → APPROVED → PRODUCTION contract. Cannot mutate production. |

## Decision-quality placement review

- **Rapid Price Dislocation Watch:** MERGE into Opportunity Radar in v16.10; no separate tab/service.
- **Unusual Volume / Volatility:** MERGE into Opportunity Radar in v16.10.
- **Short Trade Plan:** DEFER implementation; later extend the existing horizon plan architecture after dedicated acceptance design. Do not create separate short desks.
- **Horizon-specific R:R:** KEEP existing plan ownership; perform dedicated coverage enhancement only when the plan acceptance scope is explicitly frozen.
- **SPY / QQQ drawdown from ATH:** MOVE LATER / Market Intelligence; useful but not required to establish Radar architecture.
- **High-Beta / Market Sensitivity:** REUSE existing correlation/relative-strength evidence first; add only after the v16.10 Radar baseline is measured.
- **Today's Key Events:** KEEP/MERGE existing Event Intelligence and Dashboard summaries; avoid another calendar surface.
- **ATR / Technical / Fundamental coverage:** KEEP current horizon-appropriate ownership; Radar uses only cheap broad evidence and defers deep enrichment until qualification.
- **Readiness / Catalyst placement:** NO CHANGE REQUIRED; Radar feeds investigation and Smart Notifications, not deterministic readiness.

## Architecture invariants

1. Provider Router remains the executable provider authority.
2. Discovery/Scanner remains the only scanner owner.
3. Opportunity Radar does not permanently stream the full U.S. universe.
4. Radar promotions reuse the existing live allocation/reserve-slot mechanism.
5. No new provider is required by default.
6. No v17 database, v18 auth/multi-user, Portfolio, Journal, Paper Trading or execution scope is pulled forward.
7. Shadow observations cannot mutate production behavior.
8. v16.9 original roadmap closure must remain 30 FULL / 0 PARTIAL / 0 MISSING.
