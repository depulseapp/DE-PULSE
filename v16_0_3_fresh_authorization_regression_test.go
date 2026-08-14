package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestV1603FreshAuthorizationFutureEvidenceAgeIsNotDisplayedAsZero(t *testing.T) {
	now := time.Date(2026, 8, 11, 14, 0, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health := completeResearchFixture(now)
	q := quotes["AAPL"]
	q.ProviderTimestamp = now.Add(5 * time.Minute).UnixMilli()
	q.UpdatedAt = now.UnixMilli()
	quotes["AAPL"] = q
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now)
	c, _ := researchComponent(pkg, "Quote")
	if c.State == "FRESH" || c.DataAgeMs == 0 {
		t.Fatalf("future quote truth misleading: %+v", c)
	}
	f := fundamentals["AAPL"]
	f.UpdatedAt = now.Add(5 * time.Minute).UnixMilli()
	fundamentals["AAPL"] = f
	pkg = buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now)
	c, _ = researchComponent(pkg, "Fundamentals")
	if c.State == "FRESH" || c.DataAgeMs == 0 {
		t.Fatalf("future fundamental truth misleading: %+v", c)
	}
}

func TestV1603FreshAuthorizationResearchUsesOwnFreshnessAfterSharedTimestampValidity(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC) // regular session ET
	st, quotes, bars, fundamentals, last, health := completeResearchFixture(now)
	q := quotes["AAPL"]
	q.ProviderTimestamp = now.Add(-150 * time.Second).UnixMilli()
	q.UpdatedAt = now.Add(-time.Second).UnixMilli()
	quotes["AAPL"] = q
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now)
	c, _ := researchComponent(pkg, "Quote")
	if c.State != "FRESH" {
		t.Fatalf("valid 2.5m Research quote was incorrectly forced to reconciliation cadence: %+v", c)
	}
	_, _, reconCurrent := quoteEvidenceAges(q, now.UnixMilli())
	if reconCurrent {
		t.Fatal("same 2.5m quote should be outside stricter regular-session reconciliation window")
	}
}

func TestV1603FreshAuthorizationFundamentalDataAgeBoundary(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health := completeResearchFixture(now)
	f := fundamentals["AAPL"]
	f.UpdatedAt = now.Add(-72 * time.Hour).UnixMilli()
	fundamentals["AAPL"] = f
	last["research-fundamentals:AAPL"] = now.Add(-time.Minute).UnixMilli()
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now)
	c, _ := researchComponent(pkg, "Fundamentals")
	if c.State != "FRESH" {
		t.Fatalf("exact professional data-age boundary should remain accepted: %+v", c)
	}
	f.UpdatedAt = now.Add(-72*time.Hour - time.Second).UnixMilli()
	fundamentals["AAPL"] = f
	pkg = buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now)
	c, _ = researchComponent(pkg, "Fundamentals")
	if c.State != "STALE" {
		t.Fatalf("fundamental data beyond boundary not stale: %+v", c)
	}
}

func TestV1603FreshAuthorizationEmptyNaturalCorporateBackfillCanClose(t *testing.T) {
	old := alpacaDataBaseURL
	defer func() { alpacaDataBaseURL = old }()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/corporate-actions" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()
	alpacaDataBaseURL = srv.URL
	root := t.TempDir()
	t.Setenv("HOME", root)
	t.Setenv("XDG_CONFIG_HOME", root)
	t.Setenv("APPDATA", root)
	app, err := NewApplication()
	if err != nil {
		t.Fatal(err)
	}
	registerApplicationCleanup(t, app)
	app.engine = NewEngine(app)
	app.mu.Lock()
	app.state.Watchlists = []Watchlist{{ID: "day", Symbols: []string{"AAPL"}}, {ID: "swing"}, {ID: "long"}, {ID: "discovery"}}
	app.mu.Unlock()
	app.engine.refreshAlpacaCorporateActions(context.Background(), "k", "s")
	if app.engine.lastUpdated["corporate-actions-backfill:AAPL"] == 0 {
		t.Fatal("natural empty historical query should prove backfill checked and complete")
	}
}

func TestV1603FreshAuthorizationFutureReceiptCannotAppearFresh(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 30, 0, 0, time.UTC)
	st, quotes, bars, fundamentals, last, health := completeResearchFixture(now)
	q := quotes["AAPL"]
	q.ProviderTimestamp = now.Add(-time.Second).UnixMilli()
	q.UpdatedAt = now.Add(5 * time.Minute).UnixMilli()
	quotes["AAPL"] = q
	pkg := buildResearchPackageTruth(st, quotes, bars, fundamentals, nil, nil, nil, last, health, agreedResearchRec(now), now)
	c, _ := researchComponent(pkg, "Quote")
	if c.State == "FRESH" || c.CheckAgeMs == 0 {
		t.Fatalf("future receipt timestamp appeared fresh: %+v", c)
	}
}
