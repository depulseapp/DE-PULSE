package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"depulse/internal/providerlifecycle"
)

const (
	marketDataDefaultBaseURL = "https://api.marketdata.app"
	marketDataQuoteBodyLimit = 1 << 20
	marketDataQuoteTimeout   = 8 * time.Second
)

// marketDataQuoteAdapter is deliberately transport-only. APR-03 keeps Market
// Data in SHADOW and does not attach this adapter to Smart Provider Router v2,
// serving state, subscriptions, caches, or canonical persistence.
type marketDataQuoteAdapter struct {
	token   string
	baseURL string
	client  *http.Client
}

type marketDataQuotePayload struct {
	Status  string    `json:"s"`
	Symbol  []string  `json:"symbol"`
	Ask     []float64 `json:"ask"`
	Bid     []float64 `json:"bid"`
	Last    []float64 `json:"last"`
	Volume  []float64 `json:"volume"`
	Updated []int64   `json:"updated"`
	Error   string    `json:"errmsg"`
}

func newMarketDataQuoteAdapter(token string) marketDataQuoteAdapter {
	return marketDataQuoteAdapter{
		token:   strings.TrimSpace(token),
		baseURL: marketDataDefaultBaseURL,
		client:  &http.Client{Timeout: marketDataQuoteTimeout},
	}
}

func (a marketDataQuoteAdapter) fetchQuote(ctx context.Context, symbol string) (Quote, error) {
	if strings.TrimSpace(a.token) == "" {
		return Quote{}, errors.New("market data quote token is not configured")
	}
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return Quote{}, errors.New("market data quote symbol is required")
	}
	base := strings.TrimRight(strings.TrimSpace(a.baseURL), "/")
	if base == "" {
		base = marketDataDefaultBaseURL
	}
	if _, err := url.ParseRequestURI(base); err != nil {
		return Quote{}, fmt.Errorf("market data quote base URL is invalid: %w", err)
	}
	client := a.client
	if client == nil {
		client = &http.Client{Timeout: marketDataQuoteTimeout}
	}

	endpoint := base + "/v1/stocks/quotes/" + url.PathEscape(symbol) + "/"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return Quote{}, fmt.Errorf("market data quote request build failed: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+a.token)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return Quote{}, fmt.Errorf("market data quote transport failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, marketDataQuoteBodyLimit+1))
	if err != nil {
		return Quote{}, fmt.Errorf("market data quote response read failed: %w", err)
	}
	if len(body) > marketDataQuoteBodyLimit {
		return Quote{}, errors.New("market data quote response exceeds body limit")
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNonAuthoritativeInfo {
		return Quote{}, fmt.Errorf("market data quote request failed: HTTP %d", resp.StatusCode)
	}

	var payload marketDataQuotePayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return Quote{}, fmt.Errorf("market data quote payload malformed: %w", err)
	}
	if !strings.EqualFold(strings.TrimSpace(payload.Status), "ok") {
		if strings.TrimSpace(payload.Error) != "" {
			return Quote{}, fmt.Errorf("market data quote provider error: %s", strings.TrimSpace(payload.Error))
		}
		return Quote{}, fmt.Errorf("market data quote provider status is %q", strings.TrimSpace(payload.Status))
	}
	if len(payload.Symbol) != 1 || len(payload.Ask) != 1 || len(payload.Bid) != 1 || len(payload.Last) != 1 || len(payload.Volume) != 1 || len(payload.Updated) != 1 {
		return Quote{}, errors.New("market data quote schema mismatch: expected one aligned quote row")
	}
	if !strings.EqualFold(strings.TrimSpace(payload.Symbol[0]), symbol) {
		return Quote{}, fmt.Errorf("market data quote symbol mismatch: requested %s got %s", symbol, strings.TrimSpace(payload.Symbol[0]))
	}
	if payload.Last[0] <= 0 || payload.Updated[0] <= 0 {
		return Quote{}, errors.New("market data quote payload missing positive last price or provider timestamp")
	}

	q := Quote{
		Symbol:            symbol,
		Price:             payload.Last[0],
		Bid:               payload.Bid[0],
		Ask:               payload.Ask[0],
		Volume:            payload.Volume[0],
		Source:            marketDataProviderName,
		PriceScale:        1,
		FeedType:          "rest-shadow-delayed",
		DataState:         "delayed",
		ProviderTimestamp: payload.Updated[0],
		IngestedAt:        time.Now().Unix(),
	}
	return finalizeQuote(q), nil
}

// marketDataShadowRegistration is a normal ProviderRegistration contract, but
// its only route is SHADOW. routeChainsFromProviderRegistrations therefore
// rejects it from executable Router membership until a separately governed
// lifecycle promotion supplies fresh evidence.
func marketDataShadowRegistration() ProviderRegistration {
	return ProviderRegistration{
		Name:       marketDataProviderName,
		QuotaLabel: "Provider entitlement / request quota dependent",
		CostClass:  "Provider plan dependent",
		Configured: func(_ Settings, s Secrets) bool { return strings.TrimSpace(s.MarketData) != "" },
		ConfigurationFingerprint: func(_ Settings, s Secrets) [32]byte {
			return providerConfigurationFingerprint(marketDataProviderName, strings.TrimSpace(s.MarketData))
		},
		Routes: []ProviderDatasetContract{{
			Dataset:           "US Live Equities",
			Capability:        "Delayed U.S. stock quote validation",
			Priority:          100,
			Lifecycle:         providerlifecycle.Shadow,
			ContractVersion:   "marketdata-shadow-quote-v19.1.0",
			CanonicalOwner:    "Quote",
			Consumer:          "Adaptive Provider Registry shadow evaluation",
			AdapterContract:   "Bearer-authenticated Market Data delayed-stock-quote adapter normalizes one stock quote into canonical Quote without serving or Router authority",
			SchemaContract:    "GET /v1/stocks/quotes/{symbol}/; aligned symbol/ask/bid/last/volume/updated arrays with s=ok are normalized fail-closed",
			TimestampContract: "updated is preserved as provider Unix-second timestamp; IngestedAt records local receive time",
			FreshnessContract: "stock quote endpoint is delayed market truth; HTTP 200 versus 203 identifies upstream-versus-cache transport only and never upgrades DataState to live",
			FailureContract:   "fail closed on missing credential, timeout/transport error, HTTP 401/403/429/5xx, malformed JSON, provider error status, schema drift, symbol mismatch, or invalid price/timestamp",
			RightsContract:    "SHADOW validation only; no production serving, Router eligibility, redistribution, lifecycle promotion, or live-data claim without separately governed rights/entitlement evidence",
			EvidenceRef:       "governance/work-slices/ADAPT-V19-1-CANONICAL-FOUNDATION-001/closure.json#APR-03",
			ApprovalRef:       "APR-03 exact-head qualification; lifecycle remains SHADOW",
			InvalidationRule:  "endpoint/auth/schema/timestamp/freshness/rights semantics or adapter behavior change invalidates APR-03 evidence and keeps the route non-production",
			ExpectedDelay:     "Delayed stock quote endpoint; exchange/entitlement freshness is not inferred from HTTP cache status",
			Uses:              []string{"Adaptive Provider Registry shadow evaluation", "canonical quote validation"},
		}},
	}
}
