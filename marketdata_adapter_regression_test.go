package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type marketDataRoundTripFunc func(*http.Request) (*http.Response, error)

func (f marketDataRoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func marketDataQuoteFixtureServer(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v1/stocks/quotes/AAPL/" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("authorization = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

func TestMarketDataQuoteAccepts200And203WithoutClaimingLive(t *testing.T) {
	fixture := `{"s":"ok","symbol":["AAPL"],"ask":[248.10],"bid":[248.00],"last":[248.05],"volume":[1234567],"updated":[1788552000]}`
	for _, status := range []int{http.StatusOK, http.StatusNonAuthoritativeInfo} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := marketDataQuoteFixtureServer(t, status, fixture)
			defer server.Close()
			adapter := newMarketDataQuoteAdapter("test-token")
			adapter.baseURL = server.URL
			q, err := adapter.fetchQuote(context.Background(), "aapl")
			if err != nil {
				t.Fatalf("fetchQuote: %v", err)
			}
			if q.Symbol != "AAPL" || q.Price != 248.05 || q.Bid != 248.00 || q.Ask != 248.10 || q.Volume != 1234567 {
				t.Fatalf("unexpected quote: %+v", q)
			}
			if q.ProviderTimestamp != 1788552000 || q.Source != marketDataProviderName {
				t.Fatalf("unexpected source/timestamp: %+v", q)
			}
			if q.DataState != "delayed" || q.FeedType != "rest-shadow-delayed" {
				t.Fatalf("HTTP %d must remain delayed truth, got state=%q feed=%q", status, q.DataState, q.FeedType)
			}
			if q.Spread <= 0 || q.SpreadPct <= 0 {
				t.Fatalf("canonical spread normalization missing: %+v", q)
			}
		})
	}
}

func TestMarketDataQuoteFailsClosedOnHTTPFailures(t *testing.T) {
	for _, status := range []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusBadGateway} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := marketDataQuoteFixtureServer(t, status, `{"s":"error","errmsg":"fixture"}`)
			defer server.Close()
			adapter := newMarketDataQuoteAdapter("test-token")
			adapter.baseURL = server.URL
			_, err := adapter.fetchQuote(context.Background(), "AAPL")
			if err == nil || !strings.Contains(err.Error(), "HTTP ") {
				t.Fatalf("status %d err = %v", status, err)
			}
		})
	}
}

func TestMarketDataQuoteFailsClosedOnTimeout(t *testing.T) {
	adapter := newMarketDataQuoteAdapter("test-token")
	adapter.client = &http.Client{Transport: marketDataRoundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, context.DeadlineExceeded
	})}
	_, err := adapter.fetchQuote(context.Background(), "AAPL")
	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("timeout err = %v", err)
	}
}

func TestMarketDataQuoteFailsClosedOnMalformedAndSchemaDrift(t *testing.T) {
	cases := map[string]string{
		"malformed":    `{"s":"ok",`,
		"schema-drift": `{"s":"ok","symbol":"AAPL","ask":248.10,"bid":248.00,"last":248.05,"volume":1234567,"updated":1788552000}`,
		"missing-row":  `{"s":"ok","symbol":["AAPL"],"ask":[],"bid":[248.00],"last":[248.05],"volume":[1234567],"updated":[1788552000]}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			server := marketDataQuoteFixtureServer(t, http.StatusOK, body)
			defer server.Close()
			adapter := newMarketDataQuoteAdapter("test-token")
			adapter.baseURL = server.URL
			if _, err := adapter.fetchQuote(context.Background(), "AAPL"); err == nil {
				t.Fatal("expected fail-closed error")
			}
		})
	}
}

func TestMarketDataQuoteRejectsMissingCredentialBeforeTransport(t *testing.T) {
	adapter := newMarketDataQuoteAdapter(" ")
	if _, err := adapter.fetchQuote(context.Background(), "AAPL"); err == nil {
		t.Fatal("expected missing-token error")
	}
}

func TestMarketDataShadowRegistrationCannotEnterProductionRouteChain(t *testing.T) {
	reg := marketDataShadowRegistration()
	if reg.Name != marketDataProviderName || len(reg.Routes) != 1 {
		t.Fatalf("unexpected registration: %+v", reg)
	}
	route := reg.Routes[0]
	if route.Lifecycle != "SHADOW" {
		t.Fatalf("lifecycle = %q, want SHADOW", route.Lifecycle)
	}
	if providerDatasetContractProductionReady(route) {
		t.Fatal("SHADOW Market Data route must not be production ready")
	}
	chains := routeChainsFromProviderRegistrations([]ProviderRegistration{reg})
	if len(chains) != 0 {
		t.Fatalf("SHADOW Market Data unexpectedly entered executable route chain: %+v", chains)
	}
	state := adaptiveProviderRegistryCapabilityStateFromRegistrations([]ProviderRegistration{reg}, marketDataProviderName, "US Live Equities", Settings{}, Secrets{MarketData: "configured"})
	if state != providerRegistryCapabilityRegisteredNotProduction {
		t.Fatalf("registry state = %q, want %q", state, providerRegistryCapabilityRegisteredNotProduction)
	}
}
