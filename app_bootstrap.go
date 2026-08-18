package main

import (
	"embed"
	"time"
)

//go:embed renderer/*
var rendererFiles embed.FS

var masterMarketSymbols = []string{"SPY", "QQQ", "DIA", "IWM", "GLD", "SLV", "TLT", "USO", "XLK", "XLC", "XLY", "XLP", "XLE", "XLF", "XLV", "XLI", "XLB", "XLRE", "XLU", "SMH", "EWY", "EWT", "EWJ", "EWH", "MCHI", "VGK", "FEZ", "UUP", "FXY", "CYB", "FXE", "HYG", "LQD", "BNO", "CPER", "UNG"}
var specialIndexSymbols = []string{"VIX"}
var generalSymbols = append(append([]string{}, masterMarketSymbols...), specialIndexSymbols...)
var vixCandidateSymbols = []string{"^VIX", "VIX", "CBOE:VIX", "INDEX:VIX"}
var alpacaDataBaseURL = "https://data.alpaca.markets"
var finnhubAPIBaseURL = "https://finnhub.io/api/v1"
var finnhubMinRequestInterval = 1100 * time.Millisecond

const appName = "DE.PULSE"
const appVersion = "18.5.2"
const releaseChannel = "STABLE"
const buildID = "v18.5.2-stable-hotfix-20260817"
