package main

import (
	"os"
	"strings"
	"testing"
)

func TestV189ResearchRefreshKeepsCongressShadowOptionalAfterSEC(t *testing.T) {
	source, err := os.ReadFile("routed_refresh.go")
	if err != nil {
		t.Fatal(err)
	}
	text := string(source)
	secCall := "a.engine.refreshSECResearchSymbol(ctx, sym)"
	congressCall := "_ = a.engine.refreshTradeInsightCongressResearchSymbol(ctx, sym)"
	readinessCall := "ready, issues := a.engine.researchPackageReadiness(sym)"
	secAt := strings.Index(text, secCall)
	congressAt := strings.Index(text, congressCall)
	readinessAt := strings.Index(text, readinessCall)
	if secAt < 0 || congressAt < 0 || readinessAt < 0 {
		t.Fatalf("research refresh contract missing: sec=%d congress=%d readiness=%d", secAt, congressAt, readinessAt)
	}
	if !(secAt < congressAt && congressAt < readinessAt) {
		t.Fatalf("Congress SHADOW must run after authoritative SEC and before readiness evaluation: sec=%d congress=%d readiness=%d", secAt, congressAt, readinessAt)
	}

	readinessStart := strings.Index(text, "func (e *Engine) researchPackageReadinessAt")
	handlerStart := strings.Index(text, "func (a *Application) handleResearchRefresh")
	if readinessStart < 0 || handlerStart <= readinessStart {
		t.Fatal("could not isolate Research readiness implementation")
	}
	readinessSource := text[readinessStart:handlerStart]
	if strings.Contains(readinessSource, "research-congress:") || strings.Contains(readinessSource, "TradeInsight") {
		t.Fatal("Congress SHADOW alternative evidence must not become a critical Research readiness requirement")
	}
}
