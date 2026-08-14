# DE.PULSE — Canonical Approved Scope

**Status:** AUTHORITATIVE FOR APPROVED PRODUCT / ARCHITECTURE INTENT  
**Update rule:** Material additions/supersessions require explicit approval and a matching `governance/DECISION-LOG.md` entry.

This file answers: **What have we already agreed DE.PULSE should be, build, preserve, or avoid?**

It does not prove implementation status. For implementation truth, inspect source, release evidence, G1 scope, and G16 handoff.

---

# A. Permanent Product Scope

## SCOPE-001 — Product Identity
DE.PULSE is a **research + intelligence + decision-support system** designed to synthesize material market evidence into concise, explainable, context-aware intelligence.

## SCOPE-002 — No Execution Boundary
Permanent unless explicitly superseded:
- no order execution;
- no broker routing;
- no automated/semi-automated trading;
- no live/paper trading product;
- no OMS/blotter;
- no portfolio/P&L or journal product scope.

## SCOPE-003 — Smart Product / AI-Style North Star
DE.PULSE should behave like a continuously improving intelligent research system, not an information-dump terminal.

Preferred transformation:

**raw data → normalized evidence → correlation → materiality → interpretation → prioritization → explanation → outcome measurement → learning**

Raw availability never justifies prominent display.

## SCOPE-004 — UI North Star
**Complex inside → intelligent synthesis → simple outside.**

Normal users should see conclusions, material intelligence, confidence, freshness, risks/contradictions, and concise why—not internal provider/cache/queue/database machinery.

---

# B. Canonical Data / Architecture Scope

## ARCH-001 — Canonical Ownership
One canonical owner per dataset/intelligence responsibility where feasible. No duplicate routers, duplicate data silos, duplicate equivalent fetch paths, or duplicate equivalent calculations.

Engineering order:

**REUSE → CONSOLIDATE → REFACTOR → DELETE/REPLACE → ADD**

## ARCH-002 — Provider Router
One canonical Provider Router. Provider selection/routing is dataset/capability/entitlement aware rather than a single opaque provider-health flag.

Required concepts include Preferred vs Serving, fallback/routing reason, capability, entitlement, circuit/cooldown/rate-limit state, disagreement handling, budget/headroom, calls avoided, and provider usefulness.

## ARCH-003 — Shared Symbol Intelligence
A unique symbol should have one canonical equivalent intelligence/data-processing path where lawful and semantically valid. Multiple authorized users watching the same symbol should reuse/fan out shared canonical intelligence rather than independently causing duplicate provider/calculation work.

Per-user membership/preferences remain isolated.

## ARCH-004 — Global Symbol Registry
The Global Symbol Registry is the canonical symbol identity/eligibility authority and is distinct from per-user market membership.

## ARCH-005 — Performance / Scalability
Permanent concern across all roadmap/build/delivery work: provider efficiency, bounded concurrency/queues/caches, CPU/memory/GC/contention, storage/indexing, async persistence, background isolation, material-change propagation, workload priority/load shedding, UI/API latency, long-run stability, and restart/recovery.

## ARCH-006 — Data Utility / Evidence Value
No data is fetched/stored/computed/displayed/retained merely because it exists. Every datum needs a purpose, consumer, freshness/materiality expectation, independence/correlation rationale, retention policy, cost, and stale/absent/contradiction behavior.

---

# C. Provider / Intelligence Scope

## INT-001 — Rapid Move Intelligence
Rapid Move remains the low-latency event-first detector, with context-aware short windows and lifecycle such as:

**EARLY → VALIDATING → CONFIRMED → EXTENDED → FADING**

It evaluates velocity, abnormal volume, liquidity/spread, price/mcap, session, volatility/regime, sector/market sympathy, halt/corporate-action/bad-tick risk, source agreement, catalysts, prior history, and actionability. It tracks continuation/fade/reversal/extension outcomes.

## INT-002 — Opportunity Radar
Opportunity Radar remains the broad opportunity discovery/promotion layer. Use staged broad observation → PROMOTE → deeper/live processing → DEMOTE rather than streaming/deep-processing the entire market continuously.

## INT-003 — TradeInsight
TradeInsight is an additional **SHADOW / SECONDARY** provider through the same canonical Smart Router. It does not replace Finnhub/Alpaca for live/current truth or Direct SEC for canonical filing truth.

Approved uses include:
- corporate insider conviction/clusters;
- opportunistic vs routine insider behavior;
- congressional activity with disclosure-delay/freshness treatment;
- insider + Congress / Smart-Money convergence;
- historical OHLCV fallback/backfill/reconciliation;
- dividends/splits reconciliation;
- Global Symbol Registry enrichment;
- Opportunity Radar Top Movers corroboration;
- optional future controlled AI/MCP research when measured value and rights permit.

## INT-004 — Smart-Money / Multi-Source Convergence
Contextual intelligence may combine insider/congressional activity, unusual market activity, catalysts/news, technical evidence, market regime, and options context where useful. Correlated observations must not be counted as independent confirmations.

## INT-005 — Options Context Only
User trades stocks. Options data may be used where useful as contextual intelligence (e.g. unusual activity, put/call, IV, expected move) but must not silently change protected deterministic Day/Swing/Long formulas or create an options execution product.

---

# D. Reliable Actionable Universe

## UNIVERSE-001 — Reliable Securities First
DE.PULSE should not promote penny/OTC/illiquid/noisy/unreliable securities simply because their percentage moves are large.

Maintain a dynamic **Reliable Actionable Universe** classification using applicable evidence such as:
- listing venue/quality;
- price;
- dollar liquidity/volume;
- spread;
- market cap/quality where useful;
- trading/data history;
- volatility/liquidity resilience;
- corporate-action/bad-tick risk;
- manipulation/noise risk;
- provider/data quality.

Numeric thresholds should be evidence-based/configurable rather than one naive price cutoff.

Explicit approved actionable instruments include `GLD`, `SLV`, and `USO`.

---

# E. 10/10 Adaptive Stock Behavior Intelligence (ASBI)

## ASBI-001 — Purpose
For every eligible/reliable actionable symbol, build a continuously updated behavioral understanding from point-in-time evidence and outcomes so DE.PULSE can identify emerging behavior changes early and estimate what is likely to happen next without pretending certainty.

ASBI is **not** a duplicate scanner, a candlestick-pattern library, or an `RSI oversold = buy` system.

It composes existing canonical capabilities including Global Symbol Registry, Shared Symbol Intelligence, Provider Router, Rapid Move, Opportunity Radar, Historical Replay/point-in-time evidence, Market Regime/sector context, catalyst/news/SEC/earnings intelligence, options context where useful, and Smart-Money/TradeInsight contextual evidence.

## ASBI-002 — Behavioral State Machine
Model state transitions and competing next states, not binary UP/DOWN predictions.

Illustrative transitions include:

`SELLING ACCELERATION`
→ `CAPITULATION RISK`
→ `SELLER EXHAUSTION`
→ `STABILIZATION`
→ `BOUNCE ATTEMPT`
→ `REVERSAL CONFIRMED` **or** `BOUNCE FAILURE`
→ `CONTINUATION LOWER`

Also support, where evidence justifies:
- multi-day exhaustion/rebound candidate;
- dead-cat-bounce candidate;
- failed bounce;
- capitulation/washout;
- failed breakdown;
- trend pullback;
- relative-strength divergence;
- accumulation/distribution transition;
- compression → expansion;
- event overreaction/underreaction.

## ASBI-003 — Behavior Outlooks
Produce probabilistic multi-horizon **Behavior Outlooks** rather than deterministic claims.

A Behavior Outlook may include:
- current behavioral state;
- next-state probabilities;
- competing scenario/path probabilities;
- applicable horizons (intraday, next session, 1–3 days, 5 days, intermediate where supported);
- Behavior Probability Momentum/change over time;
- expected return/move distribution;
- expected upside/downside;
- MFE/MAE distributions;
- time-to-resolution;
- confidence;
- sample/analogue quality;
- supporting evidence;
- contradictory evidence;
- catalyst/regime/sector context;
- trigger/confirmation;
- invalidation;
- freshness/data quality;
- uncertainty.

## ASBI-004 — Behavior Probability vs Opportunity Quality
Separate:

**What is likely to happen?** from **Is it worth caring about?**

A high-probability move can have poor decision utility if expected magnitude is small relative to adverse excursion, spread, liquidity, volatility, event risk, uncertainty, or opportunity cost.

Optimize for decision utility/materiality, not directional hit rate alone.

## ASBI-005 — Expected Outcome Distribution
Track more than success/failure, including applicable:
- median/quantile returns;
- expected favorable/adverse excursion;
- probability of reaching meaningful levels;
- time to peak/failure/resolution;
- tail outcomes.

## ASBI-006 — Behavioral Fingerprint
Maintain a per-symbol **Behavioral Fingerprint** where history supports it, including characteristics such as:
- trend persistence;
- mean-reversion tendency;
- selloff elasticity;
- bounce durability;
- gap behavior;
- event sensitivity;
- earnings-reaction persistence;
- relative-strength persistence;
- intraday reversal tendency;
- closing behavior;
- volatility persistence;
- liquidity resilience;
- sector dependence;
- market-beta sensitivity.

## ASBI-007 — Hierarchical Learning / Cold Start
Learn at multiple levels and adapt weights based on evidence sufficiency:
1. symbol-specific behavior;
2. similar Behavioral Fingerprint peers;
3. sector/industry;
4. market regime;
5. reliable actionable universe.

Do not overfit a symbol with too little history; use broader analogues transparently when needed.

## ASBI-008 — Sequence / Path Dependence
Model **how the symbol arrived at the current state**, not just the current snapshot.

Consider prior states, acceleration/deceleration, failed bounces, support tests, gaps, volume absorption, catalyst sequence/reaction, and repeated inability/ability to make new extremes.

## ASBI-009 — Early Detection / Microstructure
Where data quality/cost justify it, use short-horizon evidence such as spread behavior, trade/price velocity, volume acceleration, range contraction, VWAP/opening-range behavior, failed breakdown/reclaim, and repeated inability to make new lows/highs.

Apply staged broad observe → PROMOTE → deep analysis → DEMOTE so full-universe processing remains bounded.

## ASBI-010 — Regime / Sector / Catalyst Conditioning
A setup must be interpreted in context. Distinguish technical/liquidity-driven moves from information-driven repricing.

Condition analogue/model evidence on relevant:
- market regime;
- sector behavior;
- volatility environment;
- catalyst type/severity;
- earnings/guidance/SEC/macro/event context;
- relative strength/liquidity/trend.

## ASBI-011 — Historical Analogues
Historical analogue similarity may use:

`price path + volume + volatility + catalyst + regime + sector + relative strength + liquidity + trend + Behavioral Fingerprint`

Report analogue count, similarity/quality, and sample sufficiency. Never use future information in historical feature construction.

## ASBI-012 — Evidence Conflict / Independence
Do not flatten all evidence into an opaque average. Preserve supporting and opposing evidence, independence, source agreement/disagreement, and material contradiction.

Correlated observations must not masquerade as multiple independent confirmations.

## ASBI-013 — Probability vs Confidence
Probability and confidence are separate.

Confidence should reflect applicable:
- sample size;
- analogue similarity;
- evidence quality;
- freshness;
- regime coverage;
- provider/source agreement;
- model stability;
- historical calibration.

## ASBI-014 — ABSTAIN / NO RELIABLE EDGE
ASBI must be allowed to decline a forecast when evidence is insufficient or unreliable, including poor history, contradictory evidence, unresolved catalysts, provider disagreement, bad/stale data, unusual corporate actions, extreme unseen regimes, low analogue similarity, or inadequate liquidity.

A high-quality system is not required to predict every symbol at every moment.

## ASBI-015 — Early Warning vs Confirmation
Separate incomplete but material behavioral change from confirmed transition.

Possible lifecycle:

**EARLY WARNING → DEVELOPING → CONFIRMED → FAILED / RESOLVED**

Surface early intelligence without overstating certainty.

## ASBI-016 — Behavior Probability Momentum
Track how scenario/state probabilities change over time. A rapid change in the system's belief can itself be material intelligence.

Prefer meaningful **behavior change** alerts over simple price-threshold alerts when validated.

## ASBI-017 — Immutable Behavior Intelligence Ledger
Every eligible forecast/outlook used for evaluation must be recorded **before** outcome realization with enough information to reconstruct what DE.PULSE knew.

Record applicable:
- prediction ID;
- symbol;
- timestamp;
- model/rule version;
- available evidence;
- behavioral state;
- scenario probabilities;
- confidence;
- analogue set/quality;
- regime/catalyst context;
- trigger/invalidation;
- expected outcomes.

Afterward record:
- actual outcome;
- MFE/MAE;
- resolution time;
- realized transition;
- calibration error;
- evidence/provider usefulness.

Never rewrite a failed historical prediction after the fact.

## ASBI-018 — Champion / Challenger
Approved production behavior intelligence remains the **Champion**. New learned models/rules run as **Challengers** in SHADOW on the same future observations.

Compare not just accuracy but calibration, false positives/misses, decision utility, MFE/MAE, regime robustness, lead time, latency/stability, drift, and usefulness.

Promotion remains:

**SHADOW → VALIDATED → APPROVED → PRODUCTION**

## ASBI-019 — Adaptive Intelligence Scorecard
ASBI acceptance/evaluation must include applicable:
- no-lookahead / walk-forward proof;
- calibration;
- false-positive rate/cost;
- misses / missed-material-move cost;
- behavior-state precision/usefulness;
- signal lead time;
- outcome distributions and MFE/MAE;
- ABSTAIN quality;
- analogue quality;
- regime robustness;
- drift;
- provider/evidence usefulness;
- unnecessary alert/processing rate;
- decision utility;
- Champion/Challenger evidence;
- Professional Trader/Investor acceptance.

## ASBI-020 — Cross-Symbol Behavioral Context
Leader/follower, sector, thematic, and sympathy relationships may provide supporting evidence, but correlation must never automatically be treated as causation.

## ASBI-021 — UI Placement
Do not create a raw ML/pattern-data-dump dashboard by default.

Prefer concise contextual Behavior Outlook intelligence in:
- Research;
- Opportunity Radar;
- Decision Queue;
- horizon-relevant Day/Swing/Long surfaces;
- Dashboard only for materially important behavior changes.

Show state, probability/confidence, scenarios, key evidence/contradictions, trigger/invalidation, uncertainty, and concise why.

## ASBI-022 — Roadmap Placement
v18/v19 begin/continue collecting trustworthy point-in-time evidence, features, outcomes, eligibility metadata, provenance, rights, and historical depth.

Major implementation/validation belongs in **v20 Adaptive Intelligence & Decision Research**, not as scope creep into v18.2–v18.5.

---

# F. Release / Governance Scope

## GOV-001 — Canonical G0–G16 Only
The permanent release gate model is defined in `governance/ADAPTIVE-OPERATING-CONTRACT.md`. Do not invent G17+ gates; add checks inside G0–G16.

## GOV-002 — Actual Artifact Audit
A release is not certified solely because source tests pass. G14 audits the actual packaged application/runtime.

## GOV-003 — Final Hashes Last
Final release hashes/SHA manifest are generated only after the final packaged/provenance contents are immutable.

## GOV-004 — Major Closure
A Major Closure & Release Assurance build is mandatory before v18→v19, v19→v20, and equivalent future major transitions.

## GOV-005 — Adaptive Feedback
G16 feeds defects, runtime evidence, provider effectiveness, outcome quality, UI findings, data usefulness, and release-process learning into the next Adaptive Roadmap/Build Plan. No issue silently disappears.

## GOV-006 — Autonomous v18 Continuation
Autonomous work through v18.5 is authorized subject to the stop conditions in the Adaptive Operating Contract.

---

# G. Roadmap Reference

Canonical placement is maintained in `governance/ROADMAP.md`.

Current approved major sequence:

**v18 Secure Multi-User + Smart Provider Intelligence**  
→ **v18.5 Mandatory Major Closure**  
→ **v19 Professional Data Infrastructure**  
→ **v19 Major Closure**  
→ **v20 Adaptive Intelligence & Decision Research + ASBI**

---

# H. Scope Lookup Rule

When a new idea is discussed, do not ask only “is this exact phrase present?”

Compare purpose, canonical owner, inputs, consumers, outputs, acceptance criteria, and roadmap placement. Classify it as:

**ALREADY APPROVED / PARTIALLY COVERED / REFINEMENT / CONFLICT-SUPERSESSION / NEW SCOPE / REJECT-NO CHANGE**

Then update this file only after approval when a material delta exists.
