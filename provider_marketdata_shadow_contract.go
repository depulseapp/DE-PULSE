package main

import "context"

// marketDataShadowQuoteFetcher is the bounded transport contract for the
// v19.1 Market Data SHADOW adopter. It exists only to make the production
// adapter contract explicit and compile-time checked; it does not grant Router
// eligibility, production serving authority, lifecycle promotion, or rights.
type marketDataShadowQuoteFetcher interface {
	fetchQuote(context.Context, string) (Quote, error)
}

// Keep the SHADOW adapter and constructor bound to a production contract
// without executing network I/O or introducing another routing/state owner.
var (
	_ marketDataShadowQuoteFetcher = marketDataQuoteAdapter{}
	_                              = newMarketDataQuoteAdapter
)
