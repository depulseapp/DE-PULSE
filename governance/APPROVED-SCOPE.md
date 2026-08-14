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
Contextual intelligence may combine insider/congressional activity, institutional/13F positioning, unusual market activity, catalysts/news, technical evidence, market regime, and options context where useful. Correlated observations must not be counted as independent confirmations.

## INT-005 — Options Context Only
User trades stocks. Options data may be used where useful as contextual intelligence (e.g. unusual activity, put/call, IV, expected move) but must not silently change protected deterministic Day/Swing/Long formulas or create an options execution product.

## INT-006 — Institutional Holdings / 13F Intelligence
Track Form 13F as **lag-aware institutional-positioning intelligence** inside the existing Smart-Money / Institutional Intelligence architecture. It is not a standalone scanner and must not be treated as live ownership truth or as `institution bought = bullish`.

Use 13F evidence to understand disclosed institutional positioning, persistence, concentration, consensus/crowding, manager behavior, sector/thematic rotation, and historical outcome usefulness.

## INT-007 — Canonical 13F Source / Provenance
**Direct SEC EDGAR is canonical filing truth** for 13F. Other providers may normalize, enrich, map, backfill, or corroborate only through the canonical Provider Router and must preserve source/provenance.

Support applicable public filing types and amendments, including 13F-HR, 13F-HR/A, 13F-NT, and 13F-NT/A, with manager identity, CIK/13F file identity where available, accession, report period, filing/acceptance timestamp, amendment/restatement semantics, and provenance.

## INT-008 — Point-in-Time / Filing-Lag Truthfulness
13F reports are quarterly disclosures of report-period holdings, not current positions. Public filings can arrive up to the applicable filing deadline after quarter end; DE.PULSE must preserve both **report-period time** and **public filing/acceptance time** and never leak a filing into historical analysis before it became public.

The UI/intelligence layer must visibly distinguish:
- `POSITION AS OF <quarter end>`;
- `FILED / PUBLIC <timestamp>`;
- current age/staleness.

Never infer an exact purchase/sale date from quarter-over-quarter changes.

## INT-009 — Coverage / Disclosure Limitations
13F intelligence must explicitly account for known disclosure limits. Among them:
- only Section 13(f) securities are reportable;
- short positions are not reported and must not be inferred as absent exposure;
- certain small positions may be omitted under Form 13F instructions;
- confidential treatment can delay public visibility of qualifying holdings;
- a public filing is a disclosed snapshot, not the manager's complete real-time portfolio.

Therefore `REPORTED EXIT` means absent from the applicable public disclosed holdings after reconciliation—not proof of a contemporaneous full economic exit.

## INT-010 — Canonical 13F Holdings Record / Reconciliation
Normalize applicable fields such as manager identity, filing/accession, report period, issuer, security class, CUSIP, FIGI where present, reported value, shares/principal amount, share/principal type, put/call designation where applicable, investment discretion, other manager, voting authority, amendment/restatement state, and source timestamps.

Reconcile amendments, restatements, additions of omitted/confidential holdings when later public, security-identifier changes, splits/corporate actions, mergers/spinoffs, duplicate manager reporting, combination/notice reports, and mapping into the Global Symbol Registry before deriving position changes.

## INT-011 — Position Change / Institutional Conviction
Derive truthful quarter-over-quarter states such as:

**NEW / INCREASED / REDUCED / REPORTED EXIT / UNCHANGED / NOT COMPARABLE / INCOMPLETE**

Institutional Conviction may consider applicable evidence such as disclosed position size, change magnitude, disclosed-13F concentration, persistence across quarters, new-position significance, manager/cohort behavior, breadth of independent accumulation/reduction, and security liquidity/market context.

Do not imply that disclosed-13F concentration equals the manager's complete portfolio weight when non-13F assets are unknown.

## INT-012 — Manager Behavioral Fingerprints
Where sufficient history exists, learn a **Manager / Institutional Behavioral Fingerprint** such as:
- typical holding horizon/persistence;
- tendency to initiate vs scale positions;
- concentration behavior;
- momentum vs mean-reversion tendencies;
- sector/thematic specialization;
- tendency to add into weakness or trim into strength;
- post-earnings/catalyst positioning patterns;
- historical outcome usefulness after public disclosure;
- reliability by security type, sector, regime, and market-cap/liquidity cohort.

Cold-start managers/cohorts must use broader evidence transparently and may ABSTAIN when history is insufficient.

## INT-013 — Consensus, Crowding & Independence
Measure institutional breadth and change without naïvely counting every filer as independent conviction.

Where possible distinguish active discretionary conviction from passive/index/benchmark effects, related managers/sub-advisers, duplicate reporting relationships, mechanically similar strategies, and common-factor exposure.

Potential outputs include accumulation/reduction breadth, consensus persistence, crowding/concentration risk, divergence between manager cohorts, and sector/thematic rotation.

## INT-014 — Adaptive Correlation / Integration
13F is contextual evidence that must correlate with existing canonical intelligence rather than live in isolation. Applicable integrations include:

`13F institutional positioning + insider activity + congressional activity + ASBI + Rapid Move + Opportunity Radar + earnings/guidance + SEC/news catalysts + market/sector regime + price/volume/relative strength + options context where useful`

The system must surface both convergence and contradiction. Old 13F accumulation must not override newer contradictory evidence merely because the manager is prominent.

Validated 13F features may influence adaptive ranking/context only through normal governance; protected deterministic Day/Swing/Long formulas remain unchanged unless separately approved.

## INT-015 — Adaptive Outcome Learning / Accountability
Preserve point-in-time institutional evidence and measure what happened after the information was actually public. Learn which managers, cohorts, position-change patterns, consensus states, sectors, regimes, catalysts, and 13F-derived features add useful decision value.

Evaluation should include applicable calibration, false-positive/miss cost, lead/lag, return/outcome distributions, MFE/MAE where meaningful, regime robustness, crowding risk, evidence independence, manager/cohort usefulness, stale-data penalties, amendment/reconciliation accuracy, and drift.

Adaptive influence remains **SHADOW → VALIDATED → APPROVED → PRODUCTION** with Champion/Challenger evaluation where models/rankings are learned.

## INT-016 — 13F UX / Performance / Delivery
13F is naturally quarterly/event-driven. Prefer SEC filing detection, incremental filing/reconciliation work, background historical computation, bounded concurrency, cached normalized holdings/deltas, and material-change propagation rather than frequent wasteful polling or full-history recomputation.

Normal UI should show concise institutional intelligence with report period, filed/public date, freshness/lag warning, change/conviction/consensus, key managers or cohorts where material, contradiction, and concise why. Raw filing tables belong in drill-down/research views.

v18/v19 should begin/continue collecting trustworthy point-in-time 13F filing/history/outcome evidence; v19 hardens data quality/rights/reconciliation; major adaptive manager/consensus learning belongs in v20 Adaptive Intelligence & Decision Research.

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

It composes existing canonical capabilities including Global Symbol Registry, Shared Symbol Intelligence, Provider Router, Rapid Move, Opportunity Radar, Historical Replay/point-in-time evidence, Market Regime/sector context, catalyst/news/SEC/earnings intelligence, Institutional Holdings/13F Intelligence, options context where useful, and Smart-Money/TradeInsight contextual evidence.

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

# F. 10/10 Two-Sided Directional Thesis & Trade Plan Intelligence (TDTI)

## TDTI-001 — Purpose / AI-Style Decision Intelligence
For each eligible ticker and relevant Day / Swing / Long horizon, DE.PULSE must reason about **both** a Long Thesis and a Short Thesis from the same canonical point-in-time evidence and return the strongest truthful decision support possible.

The goal is not to mirror `BUY` into `SELL`. The goal is an AI/LLM-style research process that understands structure, context, competing scenarios, evidence conflict, behavior state, opportunity quality, timing, and what would change the view.

Canonical outcomes are:

**LONG THESIS / SHORT THESIS / BOTH DEVELOPING / NO RELIABLE EDGE**

A ticker is never required to have an actionable side.

## TDTI-002 — One Canonical Evidence Snapshot / No Duplicate Engines
Long and Short theses must consume the **same frozen canonical evidence snapshot** for a ticker/horizon. Do not create independent long-data and short-data silos, duplicate provider fetchers, or unrelated scoring engines.

Reuse existing Day/Swing/Long canonical plan/evidence ownership, ASBI, Market/sector regime, Event Intelligence, SEC/news/earnings, liquidity, relative strength, options context, Smart-Money/13F and Historical Validation where relevant.

If implementation needs directional components, they remain coordinated children of the same canonical horizon intelligence owner.

## TDTI-003 — Competing Thesis Model
Evaluate the Long and Short theses independently but comparatively. Each side may have:
- thesis state;
- directional probability where validated;
- thesis strength;
- confidence/data sufficiency;
- opportunity quality / decision utility;
- readiness;
- supporting evidence;
- contradictory evidence;
- trigger / confirmation;
- structural invalidation;
- expected path(s);
- expected magnitude/distribution;
- time-to-resolution;
- horizon relevance;
- freshness / provenance.

A stronger directional thesis does **not** automatically imply an attractive trade plan.

## TDTI-004 — Probability, Strength, Confidence, Quality & Readiness Are Different
Keep these concepts separate:

1. **Direction Probability** — estimated chance of the relevant directional/path outcome when validated.
2. **Thesis Strength** — coherence and materiality of supporting vs opposing evidence.
3. **Confidence** — data sufficiency, freshness, sample/analogue quality, independence, provider agreement and calibration.
4. **Opportunity Quality / Decision Utility** — whether expected reward justifies adverse excursion, spread/liquidity, volatility, event risk, uncertainty, squeeze/gap risk and opportunity cost.
5. **Readiness** — whether the current price/structure has reached a valid actionable research zone and confirmation state.

Never present Setup Score as win probability.

## TDTI-005 — Long Trade Plan Contract
Where a Long Thesis is supportable, provide applicable:
- **Entry Zone**;
- Entry trigger / confirmation;
- **Trim / Target Zone**;
- optional Target 1 / Target 2 / extension target where defensible;
- **Long Invalidation** based on structural thesis failure, not an arbitrary mirrored percentage;
- Long Risk / Reward;
- Long Readiness;
- extension / chase risk;
- expected favorable/adverse excursion;
- key catalyst/regime/liquidity risks;
- concise `WHY LONG`, `WHAT CONFIRMS`, and `WHAT CHANGES THE VIEW`.

## TDTI-006 — Short Trade Plan Contract
Where a Short Thesis is supportable, provide applicable:
- **Short Entry Zone**;
- short trigger / confirmation;
- **Cover / Trim Zone**;
- optional Cover 1 / Cover 2 / downside targets where defensible;
- **Short Invalidation** based on structural bearish-thesis failure;
- Short Risk / Reward;
- Short Readiness;
- downside extension / short-chase risk;
- expected favorable/adverse excursion from the short perspective;
- squeeze/gap/catalyst/liquidity risks;
- concise `WHY SHORT`, `WHAT CONFIRMS`, and `WHAT CHANGES THE VIEW`.

Do not infer borrow availability or execution feasibility from price structure alone.

## TDTI-007 — Directional Readiness Lifecycle
Use a truthful lifecycle per side, for example:

**NOT READY → APPROACHING → IN ZONE → CONFIRMING → CONFIRMED → EXTENDED → RESOLVED / INVALIDATED**

Readiness transitions should be material-change driven. A price entering a zone without required confirmation can remain `IN ZONE / NOT CONFIRMED`.

## TDTI-008 — Structural Confirmation & Invalidation
Levels must express the thesis, not just arithmetic.

Examples include support/resistance failure or reclaim, VWAP/opening-range behavior where horizon-appropriate, failed breakout/breakdown, lower-high/higher-low structure, volume participation/absorption, relative-strength change, catalyst reaction, moving-average/trend structure, or other validated evidence.

A short thesis should generally invalidate when bearish structure is meaningfully reclaimed; a long thesis should generally invalidate when bullish structure fails. Handle gap/corporate-action/volatility exceptions explicitly.

## TDTI-009 — ASBI State / Path Integration
ASBI behavior states and competing paths must inform thesis quality and timing.

Examples:
- `SELLING ACCELERATION → SELLER EXHAUSTION → BOUNCE ATTEMPT` may keep the broader Short Thesis bearish while reducing immediate Short Opportunity Quality.
- `BOUNCE FAILURE → CONTINUATION LOWER` may strengthen Short Readiness.
- `FAILED BREAKDOWN → RECLAIM → STABILIZATION` may weaken the Short Thesis and strengthen a Long rebound thesis.

Do not let a static trend label override a material behavior-state transition.

## TDTI-010 — Thesis Probability Momentum / View Change
Track how Long/Short scenario probability, strength, confidence, quality and readiness change over time.

A rapid change such as `SHORT 72% → 54%` while rebound probability rises can itself be material intelligence.

Prefer alerts for **meaningful thesis/readiness change** over noisy simple price thresholds.

## TDTI-011 — Expected Paths, Magnitude & Timing
Where evidence supports it, model multiple directional paths rather than a single target:
- breakout/retest/continuation;
- rejection/continuation lower;
- flush/stabilize/rebound;
- failed breakdown/reclaim;
- squeeze/failed squeeze;
- catalyst repricing/mean reversion;
- consolidation before resolution.

Track expected move distributions, target/level probabilities, MFE/MAE, tail risk and time-to-resolution rather than only binary success/failure.

## TDTI-012 — Long-Specific Risk Intelligence
Long-side opportunity quality should consider applicable:
- bull-trap / failed-breakout risk;
- overhead supply/resistance;
- gap-down/catalyst/earnings risk;
- deteriorating relative strength;
- distribution/weak participation;
- excessive upside extension/chase risk;
- liquidity/spread/volatility;
- valuation/fundamental risk where horizon-relevant;
- market/sector regime contradiction;
- ASBI transition toward distribution/continuation lower.

## TDTI-013 — Short-Specific Risk Intelligence
Short-side opportunity quality should consider applicable reliable evidence such as:
- short-squeeze / violent-reversal risk;
- already-extended downside / short-chase risk;
- positive catalyst / gap-up risk;
- failed breakdown / reclaim risk;
- abnormal upside volume or relative-strength improvement;
- market/sector reversal against the short;
- earnings/event risk;
- liquidity/spread/volatility;
- short interest / days-to-cover / crowding where trustworthy;
- shortability / borrow availability/cost where lawfully and reliably available;
- SSR or other market-structure restriction context where useful.

Unavailable borrow/short-interest data must be `UNKNOWN`, never guessed.

## TDTI-014 — Cross-Horizon Intelligence
Day, Swing and Long can legitimately disagree because their horizons differ.

DE.PULSE should explain relationships such as:
- Day Long rebound inside a Swing Short trend;
- Swing Long pullback opportunity inside a Long-Term bullish thesis;
- Day Short setup while Long-Term remains bullish;
- no actionable Day edge despite a strong Long-Term thesis.

Do not flatten cross-horizon disagreement into one opaque direction.

## TDTI-015 — Evidence Conflict / Independence / Cause
Preserve bullish and bearish evidence separately. Do not average away contradictions or count correlated sources as independent confirmation.

Interpret **why** the move exists: technical/liquidity dislocation, earnings/guidance, SEC filing, macro/regime shock, sector sympathy, institutional/insider context, or other material cause.

Old evidence such as lagged 13F accumulation cannot automatically overrule fresher contradictory price/catalyst behavior.

## TDTI-016 — AI/LLM Synthesis Boundary
DE.PULSE should feel AI/LLM-style by synthesizing structured evidence into concise, context-aware reasoning:
- strongest thesis now;
- competing thesis;
- why now;
- what is already priced/extended;
- what confirms;
- what invalidates;
- what contradicts;
- what would change the view;
- what to watch next;
- uncertainty / missing evidence.

However, an LLM is not the canonical owner of market truth and must not silently invent unsupported price levels, overwrite frozen deterministic outputs, or self-modify protected production formulas. Learned/AI-derived directional influence follows normal SHADOW → VALIDATED → APPROVED → PRODUCTION governance.

## TDTI-017 — ABSTAIN / NO RELIABLE EDGE
`NO RELIABLE EDGE` is a first-class successful output when:
- Long and Short evidence is balanced/contradictory;
- expected magnitude is insufficient;
- price is too extended on both sides;
- catalyst uncertainty dominates;
- liquidity/data quality is inadequate;
- analogue/sample quality is weak;
- risk/reward is unattractive;
- thesis confidence or independence is insufficient.

The system must not manufacture actionability to fill a UI field.

## TDTI-018 — Immutable Thesis / Trade-Plan Ledger
Before outcomes are known, preserve applicable point-in-time records for both sides:
- ticker/horizon;
- evidence snapshot ID/time;
- canonical source fingerprint / model/rule version;
- Long and Short thesis states/probabilities/strength/confidence/quality/readiness;
- Entry or Short Entry Zone;
- Trim/Target or Cover/Downside Targets;
- structural invalidation;
- R:R;
- ASBI state/scenarios;
- relevant regime/catalyst/liquidity/Smart-Money evidence;
- contradictions/missing evidence;
- explanation fingerprint.

Never rewrite an old thesis after seeing the outcome.

## TDTI-019 — Side-Aware Outcome Learning / Calibration
Evaluate Long and Short separately by symbol, horizon, setup/state, regime, sector, catalyst, liquidity and relevant behavior fingerprint.

Track applicable:
- did Entry / Short Entry occur;
- did confirmation occur;
- target/trim/cover vs invalidation ordering;
- first material outcome;
- MFE/MAE anchored to the actual eligible entry event where measurable;
- elapsed time / time-to-resolution;
- 1D/3D/5D/10D and horizon-appropriate outcome distributions;
- false positives/misses;
- short-squeeze/bounce failure and bull-trap/failed-breakout outcomes;
- calibration / Decision Utility;
- ABSTAIN quality;
- drift and regime robustness.

Target-before-entry or cover-before-short-entry must never count as a valid successful plan outcome.

## TDTI-020 — Champion / Challenger & Adaptive Promotion
Production thesis logic is the Champion. New directional models/rules/weights/prompts are Challengers in SHADOW on the same future observations.

Compare calibration, decision utility, false-positive/miss costs, MFE/MAE, lead time, regime robustness, side-specific failure modes, latency/stability, explainability and Professional Trader/Investor usefulness.

Promotion remains:

**SHADOW → VALIDATED → APPROVED → PRODUCTION**

No silent production self-modification.

## TDTI-021 — UI / Surface Integration
Do not create a separate Short Desk or duplicate Long/Short app architecture by default.

Use the existing Day/Swing/Long master-detail experience and contextual surfaces. Prefer:
- stronger thesis + concise competing thesis;
- clear Long vs Short labels;
- Entry / Short Entry;
- Trim/Target / Cover/Downside Target;
- Invalidation;
- R:R;
- readiness;
- thesis probability/strength/confidence/quality only when validated and useful;
- ASBI behavior state/probability momentum;
- concise evidence/contradiction/risk/why.

Research can expose deeper Bull/Base/Bear or Long/Short evidence. Opportunity Radar and Decision Queue should rank material directional opportunities without duplicating plan computation.

## TDTI-022 — No Execution Boundary
Long/Short Trade Plans are **research plans**, not orders.

This scope does not add:
- short-sale order entry;
- share borrowing;
- broker routing;
- position sizing/position management;
- live/paper execution;
- portfolio/P&L;
- automated trading.

Any future execution capability would require a separate explicit product-boundary decision.

## TDTI-023 — Roadmap Placement
The older approved direction already included `entry / short-entry zone context when supportable`. This 10/10 contract hardens that direction into a full two-sided intelligent thesis system.

v18/v19 should preserve/build the point-in-time evidence, trade-plan/outcome lineage, short-relevant data quality/provenance and historical depth needed for validation without forcing the mature adaptive engine into v18.2–v18.5.

Major adaptive two-sided thesis implementation/validation belongs in **v20 Adaptive Intelligence & Decision Research**, integrated with ASBI rather than as a competing intelligence silo.

---

# G. 10/10 Adaptive Opportunity Discovery & Recommendations (AODR)

## AODR-001 — Purpose / AI-Style Opportunity Prioritization
DE.PULSE must convert the market intelligence it already lawfully and efficiently understands into concise, continuously refreshed opportunity discovery for users rather than leaving useful cross-market intelligence hidden.

The goal is not to produce a generic stock leaderboard or force recommendations. The goal is an AI/LLM-style research assistant that identifies **what is worth the user's attention now, what is worth watching next, and when there is no reliable opportunity**.

## AODR-002 — Canonical Composition / No Duplicate Recommendation Engine
AODR is a productization and ranking layer over existing canonical capabilities, not a new scanner, symbol store, chart engine, or independent recommendation silo.

Canonical responsibilities remain:
- **Global Symbol Registry** — identity, eligibility and instrument truth;
- **Reliable Actionable Universe** — promotion eligibility/quality boundary;
- **Shared Symbol Intelligence** — canonical reusable symbol evidence/intelligence;
- **Opportunity Radar** — broad observation, candidate discovery, PROMOTE/DEMOTE lifecycle;
- **ASBI** — behavior state/path/outlook intelligence;
- **TDTI** — Long/Short thesis, opportunity quality, readiness and research trade-plan intelligence;
- **ADR-GDI** — data/reliability truth and graceful degradation;
- **AODR** — cross-candidate prioritization, bucket placement, relevance, explanation and recommendation/outcome learning.

Reuse existing evidence and calculations. Do not refetch/recalculate equivalent symbol intelligence merely because a candidate is being ranked for recommendations.

## AODR-003 — Two User-Facing Opportunity Buckets
Provide two clearly separated opportunity views:

1. **My Market Opportunities** — strongest material opportunities among symbols already in the user's watchlists/My Market.
2. **Global Opportunities** — strongest material opportunities from the eligible broader market universe that are **not currently in that user's My Market/watchlists**.

A symbol must not appear in both buckets for the same user at the same time. When the user adds a Global Opportunity to My Market, subsequent ranking should place it under My Market rather than continue presenting it as outside discovery.

## AODR-004 — Shared Market Truth, Personalized Relevance
Underlying symbol evidence, ASBI/TDTI outputs, freshness/provenance and canonical Opportunity Quality remain shared market truth where semantically valid.

Per-user state may influence **relevance and presentation**, including applicable preferences such as horizon interest, volatility/liquidity tolerance, sectors/themes followed, Long/Short interest, symbols already followed, and prior dismissals/interactions.

Personalization must not fabricate a different price, evidence state, thesis, probability, or market truth for each user. The model is:

**shared canonical market truth → user-specific relevance/ranking/presentation**

## AODR-005 — Eligibility / Quality Before Ranking
Only symbols that satisfy the applicable Global Symbol Registry and Reliable Actionable Universe requirements may be promoted into normal Global Opportunities.

Large percentage movement alone is insufficient. Exclude or heavily penalize unreliable, illiquid, bad-tick/corporate-action-distorted, stale, unmapped, rights-constrained, or otherwise low-quality candidates.

## AODR-006 — Intelligent Opportunity Ranking
Ranking should use validated, material, horizon-aware evidence rather than one opaque score. Applicable factors include:
- TDTI Long/Short thesis and competing-thesis strength;
- Opportunity Quality / Decision Utility;
- Readiness and distance to a valid Entry/Short Entry condition;
- ASBI state, path probabilities and probability momentum;
- expected magnitude/distribution and time-to-resolution;
- market/sector regime;
- relative strength/weakness;
- liquidity/spread/volatility;
- catalyst/news/earnings/SEC context;
- Smart-Money/13F/insider/congressional context where point-in-time valid and useful;
- evidence independence/contradiction;
- freshness/provenance/data sufficiency;
- extension/chase or short-chase/squeeze risk;
- historical setup/symbol/regime usefulness;
- opportunity cost and materiality.

A strong directional view with poor entry location, poor R:R, excessive extension, weak data or low decision utility must not rank as a top opportunity merely because its raw thesis score is high.

## AODR-007 — NOW / WATCH / PASS / ABSTAIN Lifecycle
AODR must support truthful candidate states such as:

**NOW → WATCH → PASS**, with **ABSTAIN / NO RELIABLE EDGE** available whenever evidence or opportunity quality is insufficient.

Interpretation:
- **NOW** — material opportunity with sufficient quality/readiness for immediate research attention;
- **WATCH** — promising but waiting for price, confirmation, catalyst resolution, data improvement or other condition;
- **PASS** — evaluated but currently not worth surfacing prominently;
- **ABSTAIN / NO RELIABLE EDGE** — evidence is insufficient, unreliable or too conflicted to rank responsibly.

No UI quota requires DE.PULSE to manufacture a minimum number of recommendations.

## AODR-008 — Readiness Over Chasing
AODR must distinguish a good company/trend/thesis from a good **current opportunity**.

If a candidate is already materially extended, outside its Entry/Short Entry Zone, exposed to poor R:R, or showing rising reversal/squeeze risk, ranking should fall or move to WATCH/PASS even when the broader thesis remains strong.

Prefer `wait for <condition>` over `chase because score is high`.

## AODR-009 — Scalable Broad-to-Deep Processing
AODR must preserve the staged Opportunity Radar architecture:

**broad low-cost observation → PROMOTE promising candidates → deeper shared analysis → rank/surface → DEMOTE when no longer material**

Do not run expensive deep AI/LLM/provider analysis continuously across the entire Global Symbol Registry. Use bounded concurrency, shared caches/state, material-change propagation, background work, provider budgets and ADR-GDI workload protection.

## AODR-010 — Correlation / Diversity-Aware Discovery
The recommendation set should avoid presenting many highly correlated candidates as if they were independent opportunities.

Where useful, recognize sector/thematic/common-catalyst/common-factor clustering and prefer the strongest representative or a small number of genuinely differentiated candidates unless several correlated names are independently exceptional.

Diversity logic is a recommendation-quality feature, not portfolio construction or position management.

## AODR-011 — Long / Short / Horizon Awareness
AODR may surface Long or Short research opportunities across Day, Swing and Long horizons using TDTI truth.

Different horizons may legitimately disagree. Ranking must preserve horizon labels and avoid flattening a Day rebound, Swing short trend and Long-Term bullish thesis into one opaque direction.

## AODR-012 — AI/LLM-Style Synthesis Boundary
Normal recommendation UX should explain concisely:
- why this candidate matters now;
- Long/Short/horizon state;
- opportunity quality/readiness;
- what confirms;
- what invalidates;
- key contradiction/risk;
- whether the candidate is extended or still near a useful zone;
- what to watch next;
- why it outranks other available candidates when useful.

LLMs may synthesize grounded structured evidence, but they are not canonical market-truth owners and must not invent unsupported recommendations, prices, levels, probabilities or catalysts.

## AODR-013 — Immutable Recommendation / Outcome Ledger
Record material surfaced recommendations point-in-time before outcomes are known, including applicable:
- user-independent candidate/ranking snapshot ID;
- symbol/horizon/direction;
- My Market vs Global bucket;
- NOW/WATCH state;
- rank/relevance components and model/rule version;
- canonical evidence/TDTI/ASBI fingerprints;
- freshness/data sufficiency;
- reason surfaced;
- relevant Entry/Short Entry, invalidation and expected-path context;
- correlation/diversity context;
- subsequent state transitions and outcomes.

Evaluate whether surfaced opportunities were timely, useful, overextended, missed, noisy, redundant, stale, or inferior to candidates that were not surfaced. Never rewrite old recommendation history after seeing outcomes.

## AODR-014 — Adaptive Ranking / Personalization Learning
Where sufficient evidence exists, learn which ranking features, opportunity states, horizons, candidate types and presentation choices improve decision utility.

User-specific interaction learning may improve relevance, but must remain isolated from shared market truth and must not cause one user's preferences/outcomes to silently mutate another user's canonical analysis.

Learned ranking/personalization changes follow:

**SHADOW → VALIDATED → APPROVED → PRODUCTION**

Use Champion/Challenger evaluation for material learned ranking changes and measure calibration/utility, false positives/misses, alert/recommendation usefulness, redundancy, diversity, staleness, lead time and user decision value.

## AODR-015 — UI / Surface Integration
Primary user-facing labels are:

**My Market Opportunities**  
**Global Opportunities**

Opportunity Radar is the natural full discovery surface. Dashboard may show a concise summary such as actionable and developing counts plus the highest-value items. Decision Queue should receive only material candidates that become sufficiently ready rather than duplicating the recommendation engine.

Prefer concise cards/synthesis over raw ranking tables. Deeper Research can expose evidence, scenarios, TDTI/ASBI detail and historical context.

## AODR-016 — Reliability / Graceful Degradation
AODR must consume ADR-GDI capability health and dependency truth.

Optional missing context may reduce confidence/relevance without invalidating an otherwise well-supported opportunity. Missing/stale/unreliable decision-critical evidence must reduce rank, move the candidate to WATCH/PASS, or ABSTAIN as appropriate.

Do not surface a candidate as `NOW` when required current evidence is materially degraded. Explain the smallest truthful impact rather than hiding or broadly overstating degradation.

## AODR-017 — No Execution Boundary
AODR recommendations are **research and decision-support prioritization**, not trade instructions or execution.

AODR does not add order entry, broker routing, automated/semi-automated execution, share borrowing, position sizing/management, portfolio/P&L or a trading journal.

## AODR-018 — Roadmap Placement
Foundational AODR productization may begin during v18/v19 where dependency-compatible by reusing existing Global Symbol Registry, Shared Symbol Intelligence and Opportunity Radar outputs, implementing My Market vs Global bucketing, bounded ranking, point-in-time recommendation lineage and outcome collection without disrupting frozen v18.2 scope.

v19 should harden ranking data quality, point-in-time outcome history, provider/rights/reliability dependencies, cross-user shared-state efficiency, diversity/correlation handling and recommendation observability.

Major mature adaptive ranking, ASBI/TDTI-driven opportunity prioritization, personalized relevance learning and Champion/Challenger recommendation intelligence belong in **v20 Adaptive Intelligence & Decision Research**.

---

# H. Release / Governance Scope

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

# I. Roadmap Reference

Canonical placement is maintained in `governance/ROADMAP.md`.

Current approved major sequence:

**v18 Secure Multi-User + Smart Provider Intelligence**  
→ **v18.5 Mandatory Major Closure**  
→ **v19 Professional Data Infrastructure**  
→ **v19 Major Closure**  
→ **v20 Adaptive Intelligence & Decision Research + ASBI + Two-Sided Thesis Intelligence + Adaptive Opportunity Discovery & Recommendations**

---

# J. Scope Lookup Rule

When a new idea is discussed, do not ask only “is this exact phrase present?”

Compare purpose, canonical owner, inputs, consumers, outputs, acceptance criteria, and roadmap placement. Classify it as:

**ALREADY APPROVED / PARTIALLY COVERED / REFINEMENT / CONFLICT-SUPERSESSION / NEW SCOPE / REJECT-NO CHANGE**

Then update this file only after approval when a material delta exists.