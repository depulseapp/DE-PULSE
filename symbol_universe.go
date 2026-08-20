package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strings"
	"time"
)

const canonicalUSSymbolUniverseRetryTTL = 5 * time.Minute

// loadUSSymbolUniverse performs only the provider acquisition and eligibility
// normalization. Cache ownership lives in canonicalUSSymbolUniverse so Scanner
// and Opportunity Radar cannot create competing universe caches.
func (e *Engine) loadUSSymbolUniverse(ctx context.Context, key, secret string) ([]string, bool) {
	out := append([]string{}, discoverySeedUniverse...)
	client := &http.Client{Timeout: 15 * time.Second}
	headers := map[string]string{"APCA-API-KEY-ID": key, "APCA-API-SECRET-KEY": secret}
	var assets []alpacaAsset
	assetURLs := []string{
		"https://paper-api.alpaca.markets/v2/assets?status=active&asset_class=us_equity&attributes=has_options",
		"https://api.alpaca.markets/v2/assets?status=active&asset_class=us_equity&attributes=has_options",
	}
	loaded := false
	for _, raw := range assetURLs {
		assets = nil
		if err := e.providerGetJSONTier(ctx, "Alpaca", WorkTierBroadDiscovery, client, raw, headers, &assets); err == nil && len(assets) > 0 {
			loaded = true
			break
		}
	}
	if !loaded {
		return nil, false
	}

	eligible := make([]string, 0, len(assets))
	for _, a := range assets {
		sym := normalizeSymbol(a.Symbol)
		if !a.Tradable || sym == "" || strings.ContainsAny(sym, "/ .") || len(sym) > 5 {
			continue
		}
		switch strings.ToUpper(a.Exchange) {
		case "NASDAQ", "NYSE", "ARCA", "AMEX", "BATS", "NYSEARCA":
			eligible = append(eligible, sym)
		}
	}
	sort.Strings(eligible)

	// Preserve the established deterministic cross-section and provider-work cap.
	const capN = 500
	if len(eligible) <= capN {
		out = append(out, eligible...)
	} else {
		step := float64(len(eligible)-1) / float64(capN-1)
		for i := 0; i < capN; i++ {
			out = append(out, eligible[int(math.Round(float64(i)*step))])
		}
	}
	return uniqueSymbols(out), true
}

// canonicalUSSymbolUniverse is the single neutral owner of the broad eligible
// U.S.-equity universe shared by Discovery Scanner and Opportunity Radar. It
// owns only universe freshness/coalescing; BroadSnapshotBroker and Smart
// Provider Router v2 remain the existing snapshot and provider-route owners.
func (e *Engine) canonicalUSSymbolUniverse(ctx context.Context, key, secret string, now time.Time) []string {
	nowms := now.UnixMilli()
	ttlms := int64(opportunityUniverseTTL / time.Millisecond)

	for {
		e.mu.Lock()
		cached := append([]string{}, e.canonicalUSUniverse...)
		cachedAt := e.canonicalUSUniverseAt
		retryAt := e.canonicalUSUniverseRetryAt
		refresh := e.canonicalUSUniverseRefresh
		loader := e.canonicalUSUniverseLoader

		if len(cached) > 0 && cachedAt > 0 && nowms-cachedAt < ttlms {
			e.mu.Unlock()
			return cached
		}
		if len(cached) > 0 && retryAt > nowms {
			e.mu.Unlock()
			return cached
		}
		if refresh != nil {
			e.mu.Unlock()
			select {
			case <-refresh:
				continue
			case <-ctx.Done():
				if len(cached) > 0 {
					return cached
				}
				return uniqueSymbols(discoverySeedUniverse)
			}
		}

		refresh = make(chan struct{})
		e.canonicalUSUniverseRefresh = refresh
		e.mu.Unlock()

		var rows []string
		var providerLoaded bool
		if loader != nil {
			rows, providerLoaded = loader(ctx, key, secret)
		} else {
			rows, providerLoaded = e.loadUSSymbolUniverse(ctx, key, secret)
		}

		e.mu.Lock()
		switch {
		case providerLoaded && len(rows) > 0:
			e.canonicalUSUniverse = append([]string{}, rows...)
			e.canonicalUSUniverseAt = nowms
			e.canonicalUSUniverseSource = "alpaca-assets"
			e.canonicalUSUniverseRetryAt = 0
			if e.health != nil {
				e.health["scanner-universe"] = fmt.Sprintf("fresh · Alpaca assets · %d symbols", len(rows))
			}
		case len(e.canonicalUSUniverse) > 0 && e.canonicalUSUniverseAt > 0:
			// Preserve the original provider-backed timestamp. A failed refresh
			// must never make stale evidence appear fresh.
			e.canonicalUSUniverseSource = "stale-cache"
			e.canonicalUSUniverseRetryAt = now.Add(canonicalUSSymbolUniverseRetryTTL).UnixMilli()
			if e.health != nil {
				e.health["scanner-universe"] = fmt.Sprintf("stale cache · refresh failed · %d symbols", len(e.canonicalUSUniverse))
			}
		default:
			fallback := uniqueSymbols(discoverySeedUniverse)
			e.canonicalUSUniverse = append([]string{}, fallback...)
			e.canonicalUSUniverseAt = 0
			e.canonicalUSUniverseSource = "seed-fallback"
			e.canonicalUSUniverseRetryAt = now.Add(canonicalUSSymbolUniverseRetryTTL).UnixMilli()
			if e.health != nil {
				e.health["scanner-universe"] = fmt.Sprintf("seed fallback · refresh failed · %d symbols", len(fallback))
			}
		}
		e.canonicalUSUniverseRefresh = nil
		close(refresh)
		out := append([]string{}, e.canonicalUSUniverse...)
		e.mu.Unlock()
		return out
	}
}
