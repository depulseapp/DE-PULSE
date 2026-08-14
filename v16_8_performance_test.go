package main

import (
	"fmt"
	"testing"
	"time"
)

func v168PerfFixture(n int) (AppState, map[string]Quote, map[string]map[string][]Bar) {
	st := defaultState()
	symbols := make([]string, 0, n)
	quotes := make(map[string]Quote, n)
	bars := make(map[string]map[string][]Bar, n)
	now := time.Now()
	for i := 0; i < n; i++ {
		s := fmt.Sprintf("X%03d", i)
		symbols = append(symbols, s)
		px := 50.0 + float64(i%200)
		quotes[s] = Quote{Symbol: s, Price: px, Bid: px - .02, Ask: px + .02, Volume: 1_000_000, ChangePercent: float64((i%11)-5) / 10, ProviderTimestamp: now.UnixMilli(), UpdatedAt: now.UnixMilli()}
		bars[s] = map[string][]Bar{"daily": v167DailyBars(s, 25, px-3, .12, now)}
	}
	st.Watchlists = []Watchlist{{ID: "swing", Name: "Swing Watchlist", Symbols: symbols}}
	st.Settings.SwingWatchlistID = "swing"
	st.Settings.DayWatchlistID = ""
	st.Settings.LongWatchlistID = ""
	st.Settings.DiscoveryWatchlistID = ""
	return st, quotes, bars
}

func BenchmarkV168HeatMap500Symbols(b *testing.B) {
	st, q, _ := v168PerfFixture(500)
	now := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v165HeatMapForState(st, q, now)
	}
}

func BenchmarkV168Liquidity500Symbols(b *testing.B) {
	_, q, bars := v168PerfFixture(500)
	now := time.Now()
	base := map[string]LiquidityBaseline{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = deriveLiquidityStatesWithContext(q, bars, base, now)
	}
}
