package main

import (
	"context"
	"fmt"
	"time"
)

type sessionCheckpointPlan struct {
	Run    bool
	Reason string
	Late   bool
}

func preMarketCheckpointPlan(now time.Time, status PreparationJobStatus) sessionCheckpointPlan {
	if preparationRanForTradingDay(status, now) {
		return sessionCheckpointPlan{}
	}
	if preMarketPrepWindow(now) {
		if preparationRetryDue(status, now, 10*time.Minute) {
			return sessionCheckpointPlan{Run: true, Reason: "scheduled"}
		}
		return sessionCheckpointPlan{}
	}
	et := now.In(easternLocation())
	mins := et.Hour()*60 + et.Minute()
	if mins > 3*60+50 && mins < 9*60+20 && preparationRetryDue(status, now, 20*time.Minute) {
		return sessionCheckpointPlan{Run: true, Reason: "missed-window catch-up", Late: true}
	}
	return sessionCheckpointPlan{}
}

func marketOpenCheckpointPlan(now time.Time, status PreparationJobStatus) sessionCheckpointPlan {
	if preparationRanForTradingDay(status, now) {
		return sessionCheckpointPlan{}
	}
	if marketOpenPrepWindow(now) {
		if preparationRetryDue(status, now, 2*time.Minute) {
			return sessionCheckpointPlan{Run: true, Reason: "scheduled"}
		}
		return sessionCheckpointPlan{}
	}
	et := now.In(easternLocation())
	mins := et.Hour()*60 + et.Minute()
	if mins > 9*60+25 && mins <= 10*60+15 && preparationRetryDue(status, now, 5*time.Minute) {
		return sessionCheckpointPlan{Run: true, Reason: "missed-window catch-up", Late: true}
	}
	return sessionCheckpointPlan{}
}

// evaluateSessionIntelligence is the single scheduling owner for the two
// session-readiness checkpoints. The checkpoints keep distinct windows,
// retry/catch-up semantics and outputs, while execution is serialized through
// one coordinator so they cannot create parallel preparation pipelines.
func (e *Engine) evaluateSessionIntelligence(ctx context.Context, finnhubKey, alpacaKey, alpacaSecret string, now time.Time) {
	if e == nil || !e.isTradingDay(now) {
		return
	}

	e.mu.RLock()
	pre := e.preparations["pre-market-prep"]
	open := e.preparations["market-open-prep"]
	e.mu.RUnlock()

	prePlan := preMarketCheckpointPlan(now, pre)
	openPlan := marketOpenCheckpointPlan(now, open)

	runs := []string{}
	if prePlan.Run {
		e.runPreMarketPrep(ctx, finnhubKey, alpacaKey, alpacaSecret, prePlan.Reason, prePlan.Late)
		runs = append(runs, "Pre-Market Prep · "+prePlan.Reason)
	}
	if openPlan.Run {
		e.runMarketOpenPrep(openPlan.Reason)
		runs = append(runs, "Market Open Prep · "+openPlan.Reason)
	}

	e.mu.Lock()
	e.lastUpdated["session-intelligence-coordinator"] = now.UnixMilli()
	if len(runs) > 0 {
		e.health["session-intelligence-coordinator"] = fmt.Sprintf("healthy · %d checkpoint(s) executed serially", len(runs))
	} else {
		e.health["session-intelligence-coordinator"] = "healthy · scheduler active"
	}
	e.mu.Unlock()
}

// sessionIntelligenceCoordinatorLoop owns automatic Pre-Market Prep and Market
// Open Prep scheduling. It evaluates immediately on runtime start for restart / late
// catch-up, then once per minute so the narrow 9:20–9:25 market-open window is
// preserved. Manual maintenance actions remain explicit one-off checkpoint runs.
func (e *Engine) sessionIntelligenceCoordinatorLoop(ctx context.Context, finnhubKey, alpacaKey, alpacaSecret string) {
	e.evaluateSessionIntelligence(ctx, finnhubKey, alpacaKey, alpacaSecret, time.Now())
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			e.evaluateSessionIntelligence(ctx, finnhubKey, alpacaKey, alpacaSecret, now)
		}
	}
}
