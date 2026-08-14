package main

import (
	"sort"
	"strings"
	"time"
)

const (
	finnhubPlanMaxSymbols = 50
	finnhubActiveTarget   = 45 // normal pool; five reserve slots are available for urgent promotions
	alpacaPlanMaxSymbols  = 30
	alpacaActiveTarget    = 25 // normal pool; five reserve slots are available for urgent promotions/failover
	livePriorityHintTTL   = 4 * time.Minute
)

var marketCriticalLiveSymbols = []string{"SPY", "QQQ"}
var pinnedTradableLiveSymbols = []string{"GLD", "SLV", "USO"}
var coreLiveContextSymbols = []string{"IWM", "DIA", "XLK", "SMH", "HYG", "TLT"}

type LiveCoverageState struct {
	Symbol    string `json:"symbol"`
	State     string `json:"state"`
	Provider  string `json:"provider"`
	Priority  int    `json:"priority"`
	Detail    string `json:"detail,omitempty"`
	UpdatedAt int64  `json:"updatedAt,omitempty"`
}

type ManualActionStatus struct {
	Key          string `json:"key"`
	Label        string `json:"label"`
	State        string `json:"state"`
	Message      string `json:"message,omitempty"`
	LastStarted  int64  `json:"lastStarted,omitempty"`
	LastComplete int64  `json:"lastComplete,omitempty"`
	LastSuccess  int64  `json:"lastSuccess,omitempty"`
	DurationMs   int64  `json:"durationMs,omitempty"`
}

type liveAllocation struct {
	Finnhub  []string
	Alpaca   []string
	Snapshot []string
	Priority map[string]int
	Urgent   map[string]bool
}

func initialManualActions() map[string]ManualActionStatus {
	rows := []struct{ key, label string }{
		{"refresh-due", "Refresh Due Data"},
		{"global-refresh", "Global / FX Context"},
		{"capability-recheck", "Provider Capabilities"},
		{"vix-refresh", "VIX"},
		{"stream-reconnect", "Primary Live Stream"},
		{"integrity-check", "Weekly Integrity"},
		{"catalyst-evaluate", "Earnings & Material Catalyst Watch"},
	}
	out := map[string]ManualActionStatus{}
	for _, r := range rows {
		out[r.key] = ManualActionStatus{Key: r.key, Label: r.label, State: "READY"}
	}
	return out
}

func (e *Engine) setManualAction(key, label, state, message string, success bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.manualActions == nil {
		e.manualActions = initialManualActions()
	}
	row := e.manualActions[key]
	if row.Key == "" {
		row.Key = key
	}
	if label != "" {
		row.Label = label
	}
	now := time.Now().UnixMilli()
	upper := strings.ToUpper(strings.TrimSpace(state))
	row.State = upper
	row.Message = message
	if upper == "RUNNING" {
		row.LastStarted = now
		row.DurationMs = 0
	} else if upper == "COMPLETE" || upper == "DEGRADED" || upper == "FAILED" {
		row.LastComplete = now
		if row.LastStarted > 0 {
			row.DurationMs = now - row.LastStarted
		}
		if success {
			row.LastSuccess = now
		}
	}
	e.manualActions[key] = row
}

func addPrioritySymbol(out *[]string, seen map[string]bool, priority map[string]int, symbol string, p int) {
	symbol = normalizeSymbol(symbol)
	if symbol == "" || symbol == "VIX" {
		return
	}
	if old, ok := priority[symbol]; !ok || p < old {
		priority[symbol] = p
	}
	if seen[symbol] {
		return
	}
	seen[symbol] = true
	*out = append(*out, symbol)
}

// baselineLiveCandidatesFrom creates the stable normal-capacity ordering. Urgent
// state is deliberately kept separate so a newly urgent Queue/catalyst/selected
// symbol can consume a reserved slot instead of immediately churning a normal
// subscription. If the symbol is already in the normal pool, no reserve is used.
func baselineLiveCandidatesFrom(st AppState) ([]string, map[string]int) {
	out := []string{}
	seen := map[string]bool{}
	priority := map[string]int{}
	add := func(symbol string, p int) { addPrioritySymbol(&out, seen, priority, symbol, p) }
	addWL := func(id string, enabled bool, p int) {
		if !enabled {
			return
		}
		if wl, ok := watchlistValueByID(st.Watchlists, id); ok {
			for _, s := range wl.Symbols {
				add(s, p)
			}
		}
	}
	// Tier 0: SPY/QQQ and the user's active research/selection always lead normal
	// live capacity. They are market-critical context, not afterthought overflow.
	for _, s := range marketCriticalLiveSymbols {
		add(s, 0)
	}
	if sym, ok := parseUserTicker(st.UI.SelectedTicker); ok {
		add(sym, 0)
	}
	// Tier 1: explicit tradables and active Day names.
	for _, s := range pinnedTradableLiveSymbols {
		add(s, 1)
	}
	addWL(st.Settings.DayWatchlistID, st.Settings.DayEnabled, 1)
	addWL(st.Settings.SwingWatchlistID, st.Settings.SwingEnabled, 2)
	addWL(st.Settings.LongWatchlistID, st.Settings.LongEnabled, 3)
	addWL(st.Settings.DiscoveryWatchlistID, true, 3)
	for _, s := range coreLiveContextSymbols {
		add(s, 4)
	}
	return out, priority
}

func activePriorityHints(hints map[string]int64, now time.Time) []string {
	if len(hints) == 0 {
		return nil
	}
	cutoff := now.Add(-livePriorityHintTTL).UnixMilli()
	out := make([]string, 0, len(hints))
	for symbol, stamp := range hints {
		if stamp >= cutoff {
			if symbol = normalizeSymbol(symbol); symbol != "" && symbol != "VIX" {
				out = append(out, symbol)
			}
		}
	}
	sort.Strings(out)
	return out
}

func urgentLiveSymbols(st AppState, catalyst map[string]CatalystReactionState, openFlags map[string][]string, hints map[string]int64, now time.Time) []string {
	out := []string{}
	seen := map[string]bool{}
	add := func(symbol string) {
		symbol = normalizeSymbol(symbol)
		if symbol == "" || symbol == "VIX" || seen[symbol] {
			return
		}
		seen[symbol] = true
		out = append(out, symbol)
	}
	// Current user focus and current Decision Queue hints deserve instant promotion.
	add(st.UI.SelectedTicker)
	for _, s := range activePriorityHints(hints, now) {
		add(s)
	}
	// Material catalysts and Market Open risk/readiness flags are event-driven urgent state.
	cats := make([]string, 0, len(catalyst))
	for s := range catalyst {
		cats = append(cats, s)
	}
	sort.Strings(cats)
	for _, s := range cats {
		add(s)
	}
	flags := make([]string, 0, len(openFlags))
	for s := range openFlags {
		flags = append(flags, s)
	}
	sort.Strings(flags)
	for _, s := range flags {
		add(s)
	}
	return out
}

func containsLiveSymbol(xs []string, symbol string) bool {
	for _, x := range xs {
		if x == symbol {
			return true
		}
	}
	return false
}

func multiFeedAllocationWithHints(st AppState, catalyst map[string]CatalystReactionState, openFlags map[string][]string, hints map[string]int64, now time.Time) liveAllocation {
	baseline, priority := baselineLiveCandidatesFrom(st)
	urgentList := urgentLiveSymbols(st, catalyst, openFlags, hints, now)
	urgent := map[string]bool{}
	for _, s := range urgentList {
		urgent[s] = true
		if _, ok := priority[s]; !ok {
			priority[s] = 1
		}
	}
	alloc := liveAllocation{Priority: priority, Urgent: urgent}
	// v15 Provider Router: Alpaca IEX is the preferred U.S. equity feed. Keep
	// its normal streaming pool below the free-plan ceiling, then place overflow
	// in Finnhub. Both providers retain reserve headroom for urgent/failover use.
	for _, s := range baseline {
		if len(alloc.Alpaca) < alpacaActiveTarget {
			alloc.Alpaca = append(alloc.Alpaca, s)
			continue
		}
		if len(alloc.Finnhub) < finnhubActiveTarget {
			alloc.Finnhub = append(alloc.Finnhub, s)
		}
	}
	assigned := map[string]bool{}
	for _, s := range alloc.Finnhub {
		assigned[s] = true
	}
	for _, s := range alloc.Alpaca {
		assigned[s] = true
	}
	// Urgent symbols use the preferred Alpaca reserve first, then Finnhub
	// reserve. This preserves the approved Alpaca → Finnhub provider order.
	for _, s := range urgentList {
		if assigned[s] {
			continue
		}
		if len(alloc.Alpaca) < alpacaPlanMaxSymbols {
			alloc.Alpaca = append(alloc.Alpaca, s)
			assigned[s] = true
			continue
		}
		if len(alloc.Finnhub) < finnhubPlanMaxSymbols {
			alloc.Finnhub = append(alloc.Finnhub, s)
			assigned[s] = true
		}
	}
	for _, s := range masterSymbolsFromState(st) {
		if !assigned[s] {
			alloc.Snapshot = append(alloc.Snapshot, s)
		}
		if _, ok := alloc.Priority[s]; !ok {
			alloc.Priority[s] = 5
		}
	}
	sort.Strings(alloc.Snapshot)
	return alloc
}

func (e *Engine) multiFeedAllocation() liveAllocation {
	e.app.mu.RLock()
	st := clone(e.app.state)
	e.app.mu.RUnlock()
	e.mu.RLock()
	catalyst := clone(e.catalystReactions)
	openFlags := clone(e.marketOpenFlags)
	hints := clone(e.livePriorityHints)
	e.mu.RUnlock()
	return multiFeedAllocationWithHints(st, catalyst, openFlags, hints, time.Now())
}

func quoteIsRecentFinnhubLive(q Quote, now int64) bool {
	if q.Source != "finnhub-websocket" || q.ProviderTimestamp <= 0 {
		return false
	}
	return now-q.ProviderTimestamp <= 30_000
}

// quoteIsRecentAlpacaLive identifies a current Alpaca IEX websocket observation.
func quoteIsRecentAlpacaLive(q Quote, now int64) bool {
	if !strings.HasPrefix(strings.ToLower(q.Source), "alpaca-iex-websocket") || q.ProviderTimestamp <= 0 {
		return false
	}
	return now-q.ProviderTimestamp <= 30_000
}

// effectiveAlpacaIEXSymbols keeps the v15 preferred Alpaca pool subscribed and
// uses its five reserve slots for highest-priority Finnhub overflow symbols when
// their secondary stream is unavailable/stale.
func (e *Engine) effectiveAlpacaIEXSymbols() []string {
	alloc := e.multiFeedAllocation()
	e.mu.RLock()
	finnhubConnected := e.webSocketConnected
	quotes := clone(e.quotes)
	e.mu.RUnlock()
	now := time.Now().UnixMilli()
	out := []string{}
	seen := map[string]bool{}
	add := func(s string) {
		if s == "" || s == "VIX" || seen[s] || len(out) >= alpacaPlanMaxSymbols {
			return
		}
		seen[s] = true
		out = append(out, s)
	}
	// Primary pool always wins capacity.
	for _, s := range alloc.Alpaca {
		add(s)
	}
	// Remaining reserve is true failover for the next-priority Finnhub pool.
	for _, s := range alloc.Finnhub {
		if !finnhubConnected || !quoteIsRecentFinnhubLive(quotes[s], now) {
			add(s)
		}
	}
	return out
}

// effectiveFinnhubSymbols keeps normal overflow coverage and dynamically uses
// Finnhub as the second provider for Alpaca-primary symbols whose IEX observation
// is stale/unavailable. On an Alpaca disconnect, high-priority primary symbols
// are promoted before lower-priority overflow, up to Finnhub's plan ceiling.
func (e *Engine) effectiveFinnhubSymbols() []string {
	alloc := e.multiFeedAllocation()
	e.mu.RLock()
	alpacaConnected := e.alpacaWebSocketConnected
	quotes := clone(e.quotes)
	e.mu.RUnlock()
	now := time.Now().UnixMilli()
	out := []string{}
	seen := map[string]bool{}
	add := func(s string) {
		if s == "" || s == "VIX" || seen[s] || len(out) >= finnhubPlanMaxSymbols {
			return
		}
		seen[s] = true
		out = append(out, s)
	}
	// Failover primary Alpaca symbols first so trading-priority names retain live
	// coverage if the preferred stream degrades.
	for _, s := range alloc.Alpaca {
		if !alpacaConnected || !quoteIsRecentAlpacaLive(quotes[s], now) {
			add(s)
		}
	}
	for _, s := range alloc.Finnhub {
		add(s)
	}
	return out
}

func buildLiveCoverageStatesFrom(alloc liveAllocation, quotes map[string]Quote, finnhubConnected, alpacaConnected bool, finnhubSubs, alpacaSubs map[string]bool) map[string]LiveCoverageState {
	out := map[string]LiveCoverageState{}
	all := map[string]bool{}
	for _, s := range alloc.Finnhub {
		all[s] = true
	}
	for _, s := range alloc.Alpaca {
		all[s] = true
	}
	for _, s := range alloc.Snapshot {
		all[s] = true
	}
	now := time.Now().UnixMilli()
	for symbol := range all {
		q := quotes[symbol]
		stamp := q.ProviderTimestamp
		if stamp == 0 {
			stamp = q.UpdatedAt
		}
		source := strings.ToLower(q.Source)
		age := int64(0)
		if stamp > 0 {
			age = now - stamp
		}
		state, provider, detail := "UNAVAILABLE", "", "Tracked without a current quote"
		switch {
		case alpacaConnected && alpacaSubs[symbol] && strings.HasPrefix(source, "alpaca-iex-websocket") && age <= 30_000:
			state, provider, detail = "ALPACA IEX LIVE", "Alpaca IEX", "Confirmed preferred live subscription"
		case strings.Contains(source, "alpaca") && q.Price > 0 && age <= 120_000:
			state, provider, detail = "ALPACA SNAPSHOT", "Alpaca", "Preferred current Alpaca snapshot"
		case finnhubConnected && finnhubSubs[symbol] && source == "finnhub-websocket" && age <= 30_000:
			state, provider, detail = "FINNHUB LIVE", "Finnhub", "Confirmed secondary live/failover subscription"
		case strings.Contains(source, "twelvedata") && q.Price > 0 && age <= 5*60_000:
			state, provider, detail = "TWELVE SNAPSHOT", "Twelve Data", "Tertiary provider recovery snapshot"
		case q.Price > 0 && (q.DataState == "cache" || q.FeedType == "cache" || age > int64(2*time.Minute/time.Millisecond)):
			state, provider, detail = "CACHED", "Cache", "Last valid real quote retained"
		case q.Price > 0:
			state, provider, detail = "SNAPSHOT", sourceProvider(q.Source), "Current non-streaming quote"
		}
		out[symbol] = LiveCoverageState{Symbol: symbol, State: state, Provider: provider, Priority: alloc.Priority[symbol], Detail: detail, UpdatedAt: stamp}
	}
	return out
}

func uniqueLiveSubscriptionCount(finnhubSubs, alpacaSubs map[string]bool) int {
	seen := map[string]bool{}
	for s := range finnhubSubs {
		seen[s] = true
	}
	for s := range alpacaSubs {
		seen[s] = true
	}
	return len(seen)
}
