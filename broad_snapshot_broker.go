package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

const broadSnapshotBrokerMaxEntries = 1200

type BroadSnapshotBrokerDiagnostics struct {
	Requests            int64 `json:"requests"`
	ProviderFetches     int64 `json:"providerFetches"`
	CacheOnlyRequests   int64 `json:"cacheOnlyRequests"`
	CoalescedWaiters    int64 `json:"coalescedWaiters"`
	SymbolsReused       int64 `json:"symbolsReused"`
	SymbolsFetched      int64 `json:"symbolsFetched"`
	Evictions           int64 `json:"evictions"`
	Entries             int   `json:"entries"`
	MaxEntries          int   `json:"maxEntries"`
	LastProviderFetchAt int64 `json:"lastProviderFetchAt,omitempty"`
	LastReuseAt         int64 `json:"lastReuseAt,omitempty"`
}

type broadSnapshotCacheEntry struct {
	Snapshot  alpacaLiveSnapshot
	FetchedAt time.Time
	LastUsed  time.Time
}

type broadSnapshotFlight struct {
	done chan struct{}
}

type BroadSnapshotBroker struct {
	mu         sync.Mutex
	entries    map[string]broadSnapshotCacheEntry
	flights    map[string]*broadSnapshotFlight
	maxEntries int
	diag       BroadSnapshotBrokerDiagnostics
}

type BroadSnapshotAcquireResult struct {
	ProviderCallAvoided bool
	SymbolsReused       int
	SymbolsFetched      int
	Coalesced           bool
}

func NewBroadSnapshotBroker(maxEntries int) *BroadSnapshotBroker {
	if maxEntries <= 0 {
		maxEntries = broadSnapshotBrokerMaxEntries
	}
	return &BroadSnapshotBroker{
		entries:    map[string]broadSnapshotCacheEntry{},
		flights:    map[string]*broadSnapshotFlight{},
		maxEntries: maxEntries,
		diag:       BroadSnapshotBrokerDiagnostics{MaxEntries: maxEntries},
	}
}

func broadSnapshotCacheKey(feed, symbol string) string {
	return strings.ToLower(strings.TrimSpace(feed)) + "|" + normalizeSymbol(symbol)
}

func broadSnapshotTTL(feed string, tier WorkTier, now time.Time) time.Duration {
	session := marketSessionET(now)
	var ttl time.Duration
	switch session {
	case "regular":
		ttl = 30 * time.Second
	case "pre-market", "after-hours":
		ttl = 60 * time.Second
	case "overnight":
		ttl = 90 * time.Second
	default:
		ttl = 5 * time.Minute
	}
	// Radar remains more freshness-sensitive than a user-triggered broad scan.
	// The broker can still reuse a very recent observation, but it must not turn
	// the radar cadence into a stale cache cadence.
	if tier == WorkTierRadarPromoted && ttl > 30*time.Second {
		ttl = 30 * time.Second
	}
	if strings.EqualFold(feed, "overnight") && ttl < 45*time.Second {
		ttl = 45 * time.Second
	}
	return ttl
}

func parseBroadSnapshotRequest(rawURL string) (feed string, symbols []string, ok bool) {
	u, err := url.Parse(rawURL)
	if err != nil || !strings.HasSuffix(strings.TrimRight(u.Path, "/"), "/v2/stocks/snapshots") {
		return "", nil, false
	}
	feed = strings.ToLower(strings.TrimSpace(u.Query().Get("feed")))
	if feed == "" {
		feed = "iex"
	}
	for _, raw := range strings.Split(u.Query().Get("symbols"), ",") {
		if sym := normalizeSymbol(raw); sym != "" {
			symbols = append(symbols, sym)
		}
	}
	symbols = uniqueSymbols(symbols)
	if len(symbols) == 0 {
		return "", nil, false
	}
	sort.Strings(symbols)
	return feed, symbols, true
}

func broadSnapshotURLForSymbols(rawURL string, symbols []string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if len(symbols) == 0 {
		return "", fmt.Errorf("broad snapshot request has no symbols")
	}
	rows := uniqueSymbols(symbols)
	sort.Strings(rows)
	q := u.Query()
	q.Set("symbols", strings.Join(rows, ","))
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func copyBroadSnapshots(in map[string]alpacaLiveSnapshot) map[string]alpacaLiveSnapshot {
	out := make(map[string]alpacaLiveSnapshot, len(in))
	for symbol, snapshot := range in {
		out[normalizeSymbol(symbol)] = snapshot
	}
	return out
}

func (b *BroadSnapshotBroker) freshLocked(feed string, symbols []string, ttl time.Duration, now time.Time, reused map[string]bool) (map[string]alpacaLiveSnapshot, []string) {
	fresh := make(map[string]alpacaLiveSnapshot, len(symbols))
	missing := make([]string, 0, len(symbols))
	for _, symbol := range symbols {
		key := broadSnapshotCacheKey(feed, symbol)
		entry, exists := b.entries[key]
		if exists && ttl > 0 && now.Sub(entry.FetchedAt) <= ttl {
			entry.LastUsed = now
			b.entries[key] = entry
			fresh[symbol] = entry.Snapshot
			reused[symbol] = true
			continue
		}
		if exists {
			delete(b.entries, key)
		}
		missing = append(missing, symbol)
	}
	return fresh, missing
}

func (b *BroadSnapshotBroker) evictLocked() {
	for len(b.entries) > b.maxEntries {
		oldestKey := ""
		var oldest time.Time
		for key, entry := range b.entries {
			if oldestKey == "" || entry.LastUsed.Before(oldest) {
				oldestKey = key
				oldest = entry.LastUsed
			}
		}
		if oldestKey == "" {
			return
		}
		delete(b.entries, oldestKey)
		b.diag.Evictions++
	}
}

func (b *BroadSnapshotBroker) Acquire(
	ctx context.Context,
	feed string,
	symbols []string,
	ttl time.Duration,
	now time.Time,
	fetch func(context.Context, []string) (map[string]alpacaLiveSnapshot, error),
) (map[string]alpacaLiveSnapshot, BroadSnapshotAcquireResult, error) {
	if b == nil {
		return nil, BroadSnapshotAcquireResult{}, fmt.Errorf("broad snapshot broker unavailable")
	}
	feed = strings.ToLower(strings.TrimSpace(feed))
	if feed == "" {
		feed = "iex"
	}
	symbols = uniqueSymbols(symbols)
	sort.Strings(symbols)
	if len(symbols) == 0 {
		return map[string]alpacaLiveSnapshot{}, BroadSnapshotAcquireResult{}, nil
	}
	if now.IsZero() {
		now = time.Now()
	}
	if ttl <= 0 {
		ttl = 30 * time.Second
	}

	b.mu.Lock()
	b.diag.Requests++
	b.mu.Unlock()
	reused := map[string]bool{}
	coalesced := false

	for {
		b.mu.Lock()
		fresh, missing := b.freshLocked(feed, symbols, ttl, now, reused)
		if len(missing) == 0 {
			b.diag.CacheOnlyRequests++
			b.diag.SymbolsReused += int64(len(reused))
			b.diag.LastReuseAt = now.UnixMilli()
			b.diag.Entries = len(b.entries)
			b.mu.Unlock()
			return copyBroadSnapshots(fresh), BroadSnapshotAcquireResult{ProviderCallAvoided: true, SymbolsReused: len(reused), Coalesced: coalesced}, nil
		}
		if flight := b.flights[feed]; flight != nil {
			b.diag.CoalescedWaiters++
			coalesced = true
			done := flight.done
			b.mu.Unlock()
			select {
			case <-ctx.Done():
				return nil, BroadSnapshotAcquireResult{SymbolsReused: len(reused), Coalesced: true}, ctx.Err()
			case <-done:
				continue
			}
		}

		flight := &broadSnapshotFlight{done: make(chan struct{})}
		b.flights[feed] = flight
		b.mu.Unlock()

		payload, err := fetch(ctx, append([]string{}, missing...))
		completedAt := time.Now()
		if completedAt.Before(now) {
			completedAt = now
		}

		b.mu.Lock()
		b.diag.ProviderFetches++
		b.diag.LastProviderFetchAt = completedAt.UnixMilli()
		if err == nil {
			for symbol, snapshot := range payload {
				sym := normalizeSymbol(symbol)
				if sym == "" {
					continue
				}
				b.entries[broadSnapshotCacheKey(feed, sym)] = broadSnapshotCacheEntry{Snapshot: snapshot, FetchedAt: completedAt, LastUsed: completedAt}
			}
			b.diag.SymbolsFetched += int64(len(payload))
			b.evictLocked()
		}
		delete(b.flights, feed)
		close(flight.done)
		b.diag.SymbolsReused += int64(len(reused))
		if len(reused) > 0 {
			b.diag.LastReuseAt = completedAt.UnixMilli()
		}
		b.diag.Entries = len(b.entries)
		result := BroadSnapshotAcquireResult{SymbolsReused: len(reused), SymbolsFetched: len(payload), Coalesced: coalesced}
		if err != nil {
			b.mu.Unlock()
			return copyBroadSnapshots(fresh), result, err
		}
		for symbol, snapshot := range payload {
			fresh[normalizeSymbol(symbol)] = snapshot
		}
		b.mu.Unlock()
		return copyBroadSnapshots(fresh), result, nil
	}
}

func (b *BroadSnapshotBroker) Diagnostics() BroadSnapshotBrokerDiagnostics {
	if b == nil {
		return BroadSnapshotBrokerDiagnostics{}
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	out := b.diag
	out.Entries = len(b.entries)
	out.MaxEntries = b.maxEntries
	return out
}

var engineBroadSnapshotBrokers sync.Map

func broadSnapshotBrokerFor(e *Engine) *BroadSnapshotBroker {
	if e == nil {
		return nil
	}
	if existing, ok := engineBroadSnapshotBrokers.Load(e); ok {
		return existing.(*BroadSnapshotBroker)
	}
	created := NewBroadSnapshotBroker(broadSnapshotBrokerMaxEntries)
	actual, _ := engineBroadSnapshotBrokers.LoadOrStore(e, created)
	return actual.(*BroadSnapshotBroker)
}

func (e *Engine) broadSnapshotDiagnostics() BroadSnapshotBrokerDiagnostics {
	return broadSnapshotBrokerFor(e).Diagnostics()
}

func (e *Engine) acquireSharedBroadSnapshots(ctx context.Context, provider string, tier WorkTier, client *http.Client, rawURL string, headers map[string]string, out any) (bool, error) {
	if e == nil || !strings.EqualFold(strings.TrimSpace(provider), "Alpaca") {
		return false, nil
	}
	tier = workTierFromContext(ctx, tier)
	if tier != WorkTierBroadDiscovery && tier != WorkTierRadarPromoted {
		return false, nil
	}
	target, typed := out.(*map[string]alpacaLiveSnapshot)
	if !typed {
		return false, nil
	}
	feed, symbols, ok := parseBroadSnapshotRequest(rawURL)
	if !ok {
		return false, nil
	}
	broker := broadSnapshotBrokerFor(e)
	ttl := broadSnapshotTTL(feed, tier, time.Now())
	payload, result, err := broker.Acquire(ctx, feed, symbols, ttl, time.Now(), func(fetchCtx context.Context, missing []string) (map[string]alpacaLiveSnapshot, error) {
		fetchURL, urlErr := broadSnapshotURLForSymbols(rawURL, missing)
		if urlErr != nil {
			return nil, urlErr
		}
		var rows map[string]alpacaLiveSnapshot
		fetchErr := e.providerGetJSONTierUncached(fetchCtx, provider, tier, client, fetchURL, headers, &rows)
		return rows, fetchErr
	})
	if result.ProviderCallAvoided {
		e.mu.Lock()
		e.providerCallsAvoided++
		e.lastUpdated["broad-snapshot-broker"] = time.Now().UnixMilli()
		e.mu.Unlock()
	}
	*target = payload
	return true, err
}

func (e *Engine) providerGetJSONTierUncached(ctx context.Context, provider string, tier WorkTier, client *http.Client, rawURL string, headers map[string]string, out any) error {
	if e == nil {
		return getJSON(ctx, client, rawURL, headers, out)
	}
	tier = workTierFromContext(ctx, tier)
	if ok, reason := e.providerTelemetry.Allow(provider, tier); !ok {
		return fmt.Errorf("%s deferred: %s", workTierLabel(tier), reason)
	}
	release, ok := e.workload.AcquireTier(ctx, "provider-rest", tier)
	if !ok {
		if err := ctx.Err(); err != nil {
			return err
		}
		return fmt.Errorf("provider request rejected by bounded workload budget")
	}
	defer release()
	done := e.providerTelemetry.begin(provider)
	err := getJSON(ctx, client, rawURL, headers, out)
	done(err)
	return err
}
