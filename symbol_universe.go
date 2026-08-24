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
const canonicalUSSymbolUniverseProviderQuery = "status=active&asset_class=us_equity"
const canonicalUSAssetUniverseDataset = "US Asset Universe"

// Keep the established health-map key for compatibility with existing
// consumers, but all emitted diagnostics identify the neutral shared owner.
const canonicalUSSymbolUniverseHealthKey = "scanner-universe"
const canonicalUSSymbolUniverseDiagnosticLabel = "shared U.S.-equity universe"
const canonicalUSSymbolUniverseEvidenceTimeUnknown = "UNKNOWN"

var canonicalUSSymbolUniverseExchanges = map[string]struct{}{
	"NASDAQ":   {},
	"NYSE":     {},
	"ARCA":     {},
	"AMEX":     {},
	"BATS":     {},
	"NYSEARCA": {},
}

type canonicalUSUniverseTiming struct {
	RetrievedAtMS     int64
	EvidenceAtMS      int64
	EvidenceTimeState string
}

// canonicalUSUniverseAsset decodes the useful identity fields already present
// in Alpaca /v2/assets. It does not introduce a company-profile request or a
// second provider path.
type canonicalUSUniverseAsset struct {
	ID       string `json:"id"`
	Class    string `json:"class"`
	Exchange string `json:"exchange"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Tradable bool   `json:"tradable"`
}

func (a canonicalUSUniverseAsset) eligibilityAsset() alpacaAsset {
	return alpacaAsset{Symbol: a.Symbol, Status: a.Status, Tradable: a.Tradable, Exchange: a.Exchange}
}

func (a canonicalUSUniverseAsset) identityRecord(observedAt int64) InstrumentIdentityRecord {
	return InstrumentIdentityRecord{
		Symbol:          a.Symbol,
		Name:            a.Name,
		Exchange:        a.Exchange,
		AssetClass:      a.Class,
		ProviderAssetID: a.ID,
		Source:          "alpaca-assets",
		ObservedAt:      observedAt,
	}
}

// canonicalUSUniverseTimingSnapshot separates transport/cache time from source
// evidence time. Alpaca's /v2/assets payload does not carry an authoritative
// market/evidence timestamp, so EvidenceAtMS must remain 0/UNKNOWN rather than
// being fabricated from wall-clock retrieval time.
func (e *Engine) canonicalUSUniverseTimingSnapshot() canonicalUSUniverseTiming {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return canonicalUSUniverseTiming{
		RetrievedAtMS:     e.canonicalUSUniverseAt,
		EvidenceAtMS:      0,
		EvidenceTimeState: canonicalUSSymbolUniverseEvidenceTimeUnknown,
	}
}

// canonicalUSUniverseAssetEligible is the explicit discovery-universe boundary.
// The shared Scanner/Radar universe is broad U.S.-equity discovery, not an
// options-only universe: provider attributes such as has_options must not
// silently remove otherwise eligible symbols. Options capability remains a
// downstream consumer concern.
func canonicalUSUniverseAssetEligible(a alpacaAsset) bool {
	sym := normalizeSymbol(a.Symbol)
	if !strings.EqualFold(strings.TrimSpace(a.Status), "active") || !a.Tradable || sym == "" {
		return false
	}
	if strings.ContainsAny(sym, "/ .") || len(sym) > 5 {
		return false
	}
	_, ok := canonicalUSSymbolUniverseExchanges[strings.ToUpper(strings.TrimSpace(a.Exchange))]
	return ok
}

func canonicalUSSymbolUniverseAssetURLs() []string {
	return []string{
		"https://paper-api.alpaca.markets/v2/assets?" + canonicalUSSymbolUniverseProviderQuery,
		"https://api.alpaca.markets/v2/assets?" + canonicalUSSymbolUniverseProviderQuery,
	}
}

// loadUSSymbolUniverse performs only provider acquisition and explicit
// eligibility normalization. Cache ownership lives in
// canonicalUSSymbolUniverse so Scanner and Opportunity Radar cannot create
// competing universe caches.
func (e *Engine) loadUSSymbolUniverse(ctx context.Context, key, secret string) ([]string, bool) {
	return e.loadUSSymbolUniverseWithClient(ctx, key, secret, &http.Client{Timeout: 15 * time.Second})
}

// loadUSSymbolUniverseWithClient is the production Router v2 acquisition path
// with an explicit client seam for deterministic transport regression tests.
// Paper/live URLs are same-provider endpoint fallback inside one Alpaca
// capability attempt; they are never ranked as separate providers.
func (e *Engine) loadUSSymbolUniverseWithClient(ctx context.Context, key, secret string, client *http.Client) ([]string, bool) {
	out := append([]string{}, discoverySeedUniverse...)
	headers := map[string]string{"APCA-API-KEY-ID": key, "APCA-API-SECRET-KEY": secret}
	var assets []canonicalUSUniverseAsset
	attempts := map[string]providerRouteAttempt{
		"Alpaca": func(routeCtx context.Context) bool {
			var terminalErr error
			for _, raw := range canonicalUSSymbolUniverseAssetURLs() {
				assets = nil
				err := e.providerGetJSONTier(routeCtx, "Alpaca", WorkTierBroadDiscovery, client, raw, headers, &assets)
				if err == nil && len(assets) > 0 {
					e.recordProviderSuccess("Alpaca")
					return true
				}
				if err == nil {
					err = fmt.Errorf("Alpaca assets returned an empty payload")
				}
				if providerRequestFailureIsLocalNeutral(routeCtx, err) {
					return false
				}
				terminalErr = err
			}
			if terminalErr != nil {
				reportProviderRouteFailure(routeCtx, terminalErr)
			}
			return false
		},
	}
	if _, loaded := e.executeProviderRoute(ctx, canonicalUSAssetUniverseDataset, attempts); !loaded {
		return nil, false
	}

	observedAt := time.Now().UnixMilli()
	eligible := make([]string, 0, len(assets))
	identities := make([]InstrumentIdentityRecord, 0, len(assets))
	for _, a := range assets {
		if canonicalUSUniverseAssetEligible(a.eligibilityAsset()) {
			eligible = append(eligible, normalizeSymbol(a.Symbol))
			identities = append(identities, a.identityRecord(observedAt))
		}
	}
	// Identity is captured from this exact routed response. Persistence failure
	// is scoped to the slow-changing identity capability and never converts a
	// successful universe acquisition into provider failure.
	e.acceptInstrumentIdentities(identities)
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

// finishCanonicalUSUniverseRefresh is deliberately deferred immediately after
// a caller becomes refresh owner. It guarantees single-flight waiters are
// released on success, provider cancellation, or loader panic. A panic is not
// swallowed; cleanup happens first and the panic continues to the caller.
func (e *Engine) finishCanonicalUSUniverseRefresh(refresh chan struct{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.canonicalUSUniverseRefresh == refresh {
		e.canonicalUSUniverseRefresh = nil
		close(refresh)
	}
}

func (e *Engine) refreshCanonicalUSSymbolUniverse(ctx context.Context, key, secret string, now time.Time, refresh chan struct{}, loader func(context.Context, string, string) ([]string, bool)) []string {
	defer e.finishCanonicalUSUniverseRefresh(refresh)

	var rows []string
	var providerLoaded bool
	if loader != nil {
		rows, providerLoaded = loader(ctx, key, secret)
	} else {
		rows, providerLoaded = e.loadUSSymbolUniverse(ctx, key, secret)
	}

	// A caller cancellation is not provider evidence. Do not stamp retry
	// suppression or degradation from a canceled request; simply return the last
	// known cache (or deterministic seed) and let a later live caller refresh.
	if ctx.Err() != nil && !providerLoaded {
		e.mu.RLock()
		cached := append([]string{}, e.canonicalUSUniverse...)
		e.mu.RUnlock()
		if len(cached) > 0 {
			return cached
		}
		return uniqueSymbols(discoverySeedUniverse)
	}

	nowms := now.UnixMilli()
	e.mu.Lock()
	switch {
	case providerLoaded && len(rows) > 0:
		e.canonicalUSUniverse = append([]string{}, rows...)
		// This is retrieval/cache time only. It must never be presented as
		// provider evidence time because /v2/assets supplies no such timestamp.
		e.canonicalUSUniverseAt = nowms
		e.canonicalUSUniverseSource = "alpaca-assets"
		e.canonicalUSUniverseRetryAt = 0
		if e.health != nil {
			e.health[canonicalUSSymbolUniverseHealthKey] = fmt.Sprintf("%s · fresh retrieval · evidence time %s · Alpaca assets · %d symbols", canonicalUSSymbolUniverseDiagnosticLabel, canonicalUSSymbolUniverseEvidenceTimeUnknown, len(rows))
		}
	case len(e.canonicalUSUniverse) > 0 && e.canonicalUSUniverseAt > 0:
		// Preserve the original provider-backed retrieval timestamp. A failed
		// refresh must never make stale cache evidence appear fresh.
		e.canonicalUSUniverseSource = "stale-cache"
		e.canonicalUSUniverseRetryAt = now.Add(canonicalUSSymbolUniverseRetryTTL).UnixMilli()
		if e.health != nil {
			e.health[canonicalUSSymbolUniverseHealthKey] = fmt.Sprintf("%s · stale cache · evidence time %s · refresh failed · %d symbols", canonicalUSSymbolUniverseDiagnosticLabel, canonicalUSSymbolUniverseEvidenceTimeUnknown, len(e.canonicalUSUniverse))
		}
	default:
		fallback := uniqueSymbols(discoverySeedUniverse)
		e.canonicalUSUniverse = append([]string{}, fallback...)
		e.canonicalUSUniverseAt = 0
		e.canonicalUSUniverseSource = "seed-fallback"
		e.canonicalUSUniverseRetryAt = now.Add(canonicalUSSymbolUniverseRetryTTL).UnixMilli()
		if e.health != nil {
			e.health[canonicalUSSymbolUniverseHealthKey] = fmt.Sprintf("%s · seed fallback · evidence time %s · refresh failed · %d symbols", canonicalUSSymbolUniverseDiagnosticLabel, canonicalUSSymbolUniverseEvidenceTimeUnknown, len(fallback))
		}
	}
	out := append([]string{}, e.canonicalUSUniverse...)
	e.mu.Unlock()
	return out
}

// canonicalUSSymbolUniverse is the single neutral owner of the broad eligible
// U.S.-equity universe shared by Discovery Scanner and Opportunity Radar. It
// owns only universe freshness/coalescing; BroadSnapshotBroker and Smart
// Provider Router v2 remain the existing snapshot and provider-route owners.
func (e *Engine) canonicalUSSymbolUniverse(ctx context.Context, key, secret string, now time.Time) []string {
	// Slow-changing identity is persistence-first and available even when this
	// universe call is served from cache or a later provider refresh fails.
	e.warmInstrumentIdentities()
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
		return e.refreshCanonicalUSSymbolUniverse(ctx, key, secret, now, refresh, loader)
	}
}
