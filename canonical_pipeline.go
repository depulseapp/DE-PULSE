package main

import "time"

// CanonicalPipelineDiagnostics measures heavy downstream work caused by
// canonical quote changes. Every valid observation may update in-memory truth,
// while only material changes fan out into persistence/intelligence work.
type CanonicalPipelineDiagnostics struct {
	ReceivedQuoteChanges int64 `json:"receivedQuoteChanges"`
	MaterialQuoteChanges int64 `json:"materialQuoteChanges"`
	SuppressedDownstream int64 `json:"suppressedDownstream"`
	PersistenceEnqueues  int64 `json:"persistenceEnqueues"`
	CatalystEvaluations  int64 `json:"catalystEvaluations"`
}

func catalystQuoteReactionActiveForSymbol(reactions map[string]CatalystReactionState, symbol string) bool {
	r, ok := reactions[normalizeSymbol(symbol)]
	return ok && r.TriggerAt > 0 && r.CompletedAt == 0
}

// propagateCanonicalQuoteChange is the single heavy downstream owner for a
// canonical quote mutation. It does not own provider selection, deterministic
// scores, history retention, or lightweight UI broadcasting.
func (e *Engine) propagateCanonicalQuoteChange(symbol string, prev, next Quote) bool {
	if e == nil || normalizeSymbol(symbol) == "" || next.Price <= 0 {
		return false
	}
	material := quoteMateriallyChanged(prev, next)
	e.mu.Lock()
	e.canonicalPipeline.ReceivedQuoteChanges++
	if material {
		e.canonicalPipeline.MaterialQuoteChanges++
	} else {
		e.canonicalPipeline.SuppressedDownstream++
	}
	e.mu.Unlock()
	if !material {
		return false
	}
	if e.app != nil && e.app.persistence != nil {
		e.app.persistence.EnqueueQuotes(map[string]Quote{normalizeSymbol(symbol): next})
		e.mu.Lock()
		e.canonicalPipeline.PersistenceEnqueues++
		e.mu.Unlock()
	}
	e.mu.RLock()
	catalystActive := catalystQuoteReactionActiveForSymbol(e.catalystReactions, symbol)
	e.mu.RUnlock()
	if catalystActive {
		e.evaluateCatalystWatch(time.Now())
		e.mu.Lock()
		e.canonicalPipeline.CatalystEvaluations++
		e.mu.Unlock()
	}
	return true
}

// updateCanonicalSessionClose keeps history providers as the evidence owner for
// close metadata while canonical downstream propagation remains centralized.
func (e *Engine) updateCanonicalSessionClose(symbol string, sessionClose float64, sessionCloseAt int64, priorClose float64) {
	symbol = normalizeSymbol(symbol)
	if e == nil || symbol == "" || sessionClose <= 0 {
		return
	}
	e.mu.Lock()
	prev := e.quotes[symbol]
	next := prev
	next.Symbol = symbol
	next.SessionClose = sessionClose
	next.SessionCloseAt = sessionCloseAt
	if priorClose > 0 {
		next.PriorSessionClose = priorClose
		if next.PreviousClose == 0 {
			next.PreviousClose = priorClose
		}
	}
	e.quotes[symbol] = next
	e.mu.Unlock()
	e.propagateCanonicalQuoteChange(symbol, prev, next)
}
