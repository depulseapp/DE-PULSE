package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func watchlistValueByID(watchlists []Watchlist, id string) (Watchlist, bool) {
	for _, wl := range watchlists {
		if wl.ID == id {
			return wl, true
		}
	}
	return Watchlist{}, false
}

func ensureDedicatedDeskWatchlists(st *AppState, def AppState) {
	original := clone(st.Watchlists)
	if len(original) == 0 {
		original = clone(def.Watchlists)
	}
	assigned := map[string]string{
		"swing":     st.Settings.SwingWatchlistID,
		"day":       st.Settings.DayWatchlistID,
		"long":      st.Settings.LongWatchlistID,
		"discovery": st.Settings.DiscoveryWatchlistID,
	}
	defaults := map[string]Watchlist{}
	for _, wl := range def.Watchlists {
		defaults[wl.ID] = wl
	}
	canonicalNames := map[string]string{
		"swing":     "Swing Watchlist",
		"day":       "Day Trade Watchlist",
		"long":      "Long-Term Watchlist",
		"discovery": "Discovery Watchlist",
	}
	canonical := make([]Watchlist, 0, 4)
	usedSources := map[string]bool{}
	for _, desk := range []string{"swing", "day", "long", "discovery"} {
		var symbols []string
		if existing, ok := watchlistValueByID(original, desk); ok {
			symbols = clone(existing.Symbols)
			usedSources[desk] = true
		} else if sourceID := assigned[desk]; sourceID != "" {
			if source, ok := watchlistValueByID(original, sourceID); ok {
				symbols = clone(source.Symbols)
				usedSources[sourceID] = true
			}
		}
		if symbols == nil {
			symbols = clone(defaults[desk].Symbols)
		}
		normalizedSymbols := userTradingSymbols(symbols)
		if normalizedSymbols == nil {
			normalizedSymbols = []string{}
		}
		canonical = append(canonical, Watchlist{ID: desk, Name: canonicalNames[desk], Symbols: normalizedSymbols})
	}
	for _, wl := range original {
		if wl.ID == "swing" || wl.ID == "day" || wl.ID == "long" || wl.ID == "discovery" || usedSources[wl.ID] {
			continue
		}
		if wl.ID == "" {
			wl.ID = randomID("watchlist")
		}
		if strings.TrimSpace(wl.Name) == "" {
			wl.Name = "Archived Watchlist"
		}
		wl.Symbols = uniqueSymbols(wl.Symbols)
		canonical = append(canonical, wl)
	}
	st.Watchlists = canonical
	st.Settings.SwingWatchlistID = "swing"
	st.Settings.DayWatchlistID = "day"
	st.Settings.LongWatchlistID = "long"
	st.Settings.DiscoveryWatchlistID = "discovery"
	st.UI.ScopeType = "watchlist"
	if st.UI.WatchlistID == "" || (st.UI.WatchlistID != "swing" && st.UI.WatchlistID != "day" && st.UI.WatchlistID != "long" && st.UI.WatchlistID != "discovery") {
		st.UI.WatchlistID = "swing"
	}
}

func activeDeskWatchlistIDs(settings Settings) []string {
	var ids []string
	if settings.SwingEnabled {
		ids = append(ids, settings.SwingWatchlistID)
	}
	if settings.DayEnabled {
		ids = append(ids, settings.DayWatchlistID)
	}
	if settings.LongEnabled {
		ids = append(ids, settings.LongWatchlistID)
	}
	return ids
}

func activeDeskSymbolsFromState(st AppState) []string {
	var out []string
	for _, id := range activeDeskWatchlistIDs(st.Settings) {
		if wl, ok := watchlistValueByID(st.Watchlists, id); ok {
			out = append(out, wl.Symbols...)
		}
	}
	return uniqueSymbols(out)
}

func discoverySymbolsFromState(st AppState) []string {
	if wl, ok := watchlistValueByID(st.Watchlists, st.Settings.DiscoveryWatchlistID); ok {
		return uniqueSymbols(wl.Symbols)
	}
	return nil
}

func analysisSymbolsFromState(st AppState) []string {
	out := append([]string{}, activeDeskSymbolsFromState(st)...)
	out = append(out, discoverySymbolsFromState(st)...)
	out = append(out, st.UI.SelectedTicker)
	return uniqueSymbols(out)
}

// masterSymbolsFromState is the canonical equity/ETF universe shared by Market
// Instruments, the three trading desks, and staged Discovery symbols. VIX is
// deliberately excluded because it uses the dedicated true-index provider path.
func masterSymbolsFromState(st AppState) []string {
	out := append([]string{}, masterMarketSymbols...)
	out = append(out, activeDeskSymbolsFromState(st)...)
	out = append(out, discoverySymbolsFromState(st)...)
	out = append(out, st.UI.SelectedTicker)
	uniq := uniqueSymbols(out)
	filtered := make([]string, 0, len(uniq))
	for _, symbol := range uniq {
		if symbol != "VIX" {
			filtered = append(filtered, symbol)
		}
	}
	return filtered
}

func specialIndexSymbolsFromState(st AppState) []string {

	return append([]string{}, specialIndexSymbols...)
}

func requiredSymbolsFromState(st AppState) []string {
	out := append([]string{}, masterSymbolsFromState(st)...)
	out = append(out, specialIndexSymbolsFromState(st)...)
	return uniqueSymbols(out)
}

func defaultState() AppState {
	return AppState{
		Version: 17,
		Settings: Settings{
			DataMode: "demo", AIProvider: "groq", AIRoutingMode: "manual", GroqModel: "openai/gpt-oss-120b", OpenRouterMode: "fast", OpenRouterSpecificModel: "openai/gpt-5.6-sol", GeminiModel: "gemini-3.1-flash-lite", SignalProfile: "balanced", MarketContext: 15, EarningsPenalty: 10,
			SwingEnabled: true, DayEnabled: true, LongEnabled: true, OvernightDataMode: "auto",
			SwingWatchlistID: "swing", DayWatchlistID: "day", LongWatchlistID: "long", DiscoveryWatchlistID: "discovery", GlobalProviderMode: "auto", OptionsDataMode: "auto", MacroEventModeEnabled: true, ResearchAIMode: "manual",
		},
		Watchlists: []Watchlist{
			{ID: "swing", Name: "Swing Watchlist", Symbols: []string{"NVDA", "META", "ORCL", "PLTR", "TSLA", "SOFI"}},
			{ID: "day", Name: "Day Trade Watchlist", Symbols: []string{"NVDA", "TSLA", "PLTR"}},
			{ID: "long", Name: "Long-Term Watchlist", Symbols: []string{"SPY", "QQQ", "META", "NVDA"}},
			{ID: "discovery", Name: "Discovery Watchlist", Symbols: []string{}},
		},
		UI:             UIState{ScopeType: "watchlist", WatchlistID: "swing", SelectedTicker: "NVDA"},
		ProviderStatus: map[string]ProviderTestResult{},
	}
}

func NewApplication() (*Application, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	configDir, err := resolveV18RuntimeConfig(base)
	if err != nil {
		return nil, err
	}
	app := &Application{configDir: configDir, hub: NewHub(), aiCache: map[string]aiCacheEntry{}, httpTelemetry: NewRequestTelemetry()}
	app.load()
	app.persistence = NewPersistenceManager(configDir)
	if err := configureStartupPersistenceRestore(app.persistence); err != nil {
		_ = app.persistence.Close()
		return nil, fmt.Errorf("persistence restore: %w", err)
	}
	identity, err := NewIdentityService(app.persistence)
	if err != nil {
		_ = app.persistence.Close()
		return nil, fmt.Errorf("identity initialization: %w", err)
	}
	app.identity = identity
	if err := app.initializeUserWorkspaces(); err != nil {
		_ = app.persistence.Close()
		return nil, fmt.Errorf("workspace initialization: %w", err)
	}
	app.persistence.EnqueueSymbols(symbolRegistryRecords(app.processingStateLocked(), time.Now()))
	if err := configureStartupPersistenceExport(app.persistence); err != nil {
		_ = app.persistence.Close()
		return nil, fmt.Errorf("persistence export: %w", err)
	}
	app.engine = NewEngine(app)
	return app, nil
}

func (a *Application) statePath() string   { return filepath.Join(a.configDir, "state.json") }
func (a *Application) secretsPath() string { return filepath.Join(a.configDir, "secrets.json") }
func (a *Application) cachePath() string   { return filepath.Join(a.configDir, "market-cache.json") }

func (a *Application) load() {
	a.state = defaultState()
	legacyAIProviderMissing := false
	migratedProvider := false
	if data, err := os.ReadFile(a.statePath()); err == nil {
		var st AppState
		if json.Unmarshal(data, &st) == nil {
			legacyAIProviderMissing = strings.TrimSpace(st.Settings.AIProvider) == ""
			a.state = mergeState(st)
		}
	}
	if data, err := os.ReadFile(a.secretsPath()); err == nil {
		_ = json.Unmarshal(data, &a.secrets)
	}

	if strings.EqualFold(strings.TrimSpace(a.state.Settings.AIProvider), "openai") {
		switch {
		case strings.TrimSpace(a.secrets.OpenRouter) != "":
			a.state.Settings.AIProvider = "openrouter"
		case strings.TrimSpace(a.secrets.Groq) != "":
			a.state.Settings.AIProvider = "groq"
		case strings.TrimSpace(a.secrets.Gemini) != "":
			a.state.Settings.AIProvider = "gemini"
		default:
			a.state.Settings.AIProvider = "groq"
		}
		migratedProvider = true
	}

	if legacyAIProviderMissing {
		switch {
		case strings.TrimSpace(a.secrets.Groq) != "":
			a.state.Settings.AIProvider = "groq"
		case strings.TrimSpace(a.secrets.OpenRouter) != "":
			a.state.Settings.AIProvider = "openrouter"
		case strings.TrimSpace(a.secrets.Gemini) != "":
			a.state.Settings.AIProvider = "gemini"
		}
		migratedProvider = true
	}
	if migratedProvider {
		_ = a.saveLocked()
	}
}
func mergeState(st AppState) AppState {
	def := defaultState()
	originalVersion := st.Version
	if st.Version == 0 {
		st.Version = 1
	}
	if st.Settings.DataMode == "" {
		st.Settings.DataMode = def.Settings.DataMode
	}
	if st.Settings.AIProvider == "" {
		st.Settings.AIProvider = def.Settings.AIProvider
	}
	if st.Settings.AIRoutingMode == "" {
		st.Settings.AIRoutingMode = def.Settings.AIRoutingMode
	}
	if st.Settings.GroqModel == "" {
		st.Settings.GroqModel = def.Settings.GroqModel
	}
	if st.Settings.GeminiModel == "" {
		st.Settings.GeminiModel = def.Settings.GeminiModel
	}
	if st.Settings.OpenRouterMode == "" {
		st.Settings.OpenRouterMode = def.Settings.OpenRouterMode
	}
	if st.Settings.OpenRouterSpecificModel == "" {
		st.Settings.OpenRouterSpecificModel = def.Settings.OpenRouterSpecificModel
	}
	if strings.TrimSpace(st.Settings.ResearchAIMode) == "" {
		st.Settings.ResearchAIMode = def.Settings.ResearchAIMode
	}
	if len(st.Watchlists) == 0 {
		st.Watchlists = def.Watchlists
	}
	for i := range st.Watchlists {
		if st.Watchlists[i].ID == "" {
			st.Watchlists[i].ID = randomID("watchlist")
		}
		if strings.TrimSpace(st.Watchlists[i].Name) == "" {
			st.Watchlists[i].Name = "Watchlist"
		}
		st.Watchlists[i].Symbols = uniqueSymbols(st.Watchlists[i].Symbols)
	}
	if st.UI.ScopeType == "" {
		st.UI = def.UI
	}
	if st.ProviderStatus == nil {
		st.ProviderStatus = map[string]ProviderTestResult{}
	}
	if st.Settings.SignalProfile == "" {
		st.Settings.SignalProfile = def.Settings.SignalProfile
	}
	if st.Settings.MarketContext <= 0 {
		st.Settings.MarketContext = def.Settings.MarketContext
	}
	if st.Settings.EarningsPenalty <= 0 {
		st.Settings.EarningsPenalty = def.Settings.EarningsPenalty
	}
	if originalVersion < 12 {
		st.Settings.SwingEnabled = true
		st.Settings.DayEnabled = true
		st.Settings.LongEnabled = true
	}
	switch strings.ToLower(strings.TrimSpace(st.Settings.OvernightDataMode)) {
	case "auto", "indicative", "live":
	default:
		st.Settings.OvernightDataMode = "auto"
	}
	switch strings.ToLower(strings.TrimSpace(st.Settings.AIProvider)) {
	case "groq", "openrouter", "gemini":
		st.Settings.AIProvider = strings.ToLower(strings.TrimSpace(st.Settings.AIProvider))
	default:
		st.Settings.AIProvider = def.Settings.AIProvider
	}
	switch strings.ToLower(strings.TrimSpace(st.Settings.AIRoutingMode)) {
	case "manual", "efficient", "balanced", "deep":
		st.Settings.AIRoutingMode = strings.ToLower(strings.TrimSpace(st.Settings.AIRoutingMode))
	default:
		st.Settings.AIRoutingMode = def.Settings.AIRoutingMode
	}
	switch strings.ToLower(strings.TrimSpace(st.Settings.OpenRouterMode)) {
	case "fast", "balanced", "powerful", "specific":
		st.Settings.OpenRouterMode = strings.ToLower(strings.TrimSpace(st.Settings.OpenRouterMode))
	default:
		st.Settings.OpenRouterMode = "fast"
	}
	if !allowedOpenRouterModel(st.Settings.OpenRouterSpecificModel) {
		st.Settings.OpenRouterSpecificModel = def.Settings.OpenRouterSpecificModel
	}
	switch strings.ToLower(strings.TrimSpace(st.Settings.GlobalProviderMode)) {
	case "auto", "direct", "free-first", "proxy":
		st.Settings.GlobalProviderMode = strings.ToLower(strings.TrimSpace(st.Settings.GlobalProviderMode))
	default:
		st.Settings.GlobalProviderMode = "auto"
	}
	switch strings.ToLower(strings.TrimSpace(st.Settings.OptionsDataMode)) {
	case "auto", "opra", "indicative", "off":
		st.Settings.OptionsDataMode = strings.ToLower(strings.TrimSpace(st.Settings.OptionsDataMode))
	default:
		st.Settings.OptionsDataMode = "auto"
	}
	switch strings.ToLower(strings.TrimSpace(st.Settings.ResearchAIMode)) {
	case "manual", "automatic":
		st.Settings.ResearchAIMode = strings.ToLower(strings.TrimSpace(st.Settings.ResearchAIMode))
	default:
		st.Settings.ResearchAIMode = "manual"
	}
	if originalVersion < 17 {
		st.Settings.MacroEventModeEnabled = true
	}
	ensureDedicatedDeskWatchlists(&st, def)
	if selected, ok := parseSelectableTicker(st.UI.SelectedTicker); ok {
		st.UI.SelectedTicker = selected
	} else {
		st.UI.SelectedTicker = ""
		if wl, found := watchlistValueByID(st.Watchlists, st.UI.WatchlistID); found && len(wl.Symbols) > 0 {
			st.UI.SelectedTicker = wl.Symbols[0]
		}
		if st.UI.SelectedTicker == "" {
			st.UI.SelectedTicker = "SPY"
		}
	}
	st.Version = 17
	return st
}

func (a *Application) saveLocked() error {
	data, err := json.MarshalIndent(a.state, "", "  ")
	if err != nil {
		return err
	}
	if err := atomicWrite(a.statePath(), data, 0600); err != nil {
		return err
	}
	secretData, err := json.MarshalIndent(a.secrets, "", "  ")
	if err != nil {
		return err
	}
	if err := atomicWrite(a.secretsPath(), secretData, 0600); err != nil {
		return err
	}
	if a.persistence != nil {
		a.persistence.EnqueueSymbols(symbolRegistryRecords(a.processingStateLocked(), time.Now()))
	}
	return nil
}

func atomicWrite(path string, data []byte, mode fs.FileMode) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, mode); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func keyHint(v string) string {
	v = strings.TrimSpace(v)
	if len(v) < 7 {
		return ""
	}
	return v[:3] + "…" + v[len(v)-3:]
}
func (a *Application) publicStateForUser(userID string) PublicState {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.publicStateLockedForUser(userID)
}
func (a *Application) cacheInfoLocked() CacheInfo {
	info := CacheInfo{}
	if st, err := os.Stat(a.cachePath()); err == nil {
		info.SizeBytes = st.Size()
		info.LastUpdated = st.ModTime().UnixMilli()
	}
	if data, err := os.ReadFile(a.cachePath()); err == nil {
		var c MarketCache
		if json.Unmarshal(data, &c) == nil {
			info.CachedSymbols = len(c.Quotes)
			if c.SavedAt > 0 {
				info.LastUpdated = c.SavedAt
			}
		}
	}
	return info
}
func (a *Application) publicStateLockedForUser(userID string) PublicState {
	view := a.workspaceStateLocked(userID)
	ensureDedicatedDeskWatchlists(&view, defaultEmptyWorkspaceState())
	return PublicState{Version: appVersion, BuildID: buildID, Settings: view.Settings,
		HasFinnhubKey:      strings.TrimSpace(a.secrets.Finnhub) != "",
		HasTradeInsightKey: strings.TrimSpace(a.secrets.TradeInsight) != "",
		HasAlpacaKey:       strings.TrimSpace(a.secrets.AlpacaKey) != "", HasAlpacaSecret: strings.TrimSpace(a.secrets.AlpacaSecret) != "", HasGroqKey: strings.TrimSpace(a.secrets.Groq) != "", HasOpenRouterKey: strings.TrimSpace(a.secrets.OpenRouter) != "", HasGeminiKey: strings.TrimSpace(a.secrets.Gemini) != "", HasFREDKey: strings.TrimSpace(a.secrets.FRED) != "", HasBLSKey: strings.TrimSpace(a.secrets.BLS) != "", HasEIAKey: strings.TrimSpace(a.secrets.EIA) != "", HasTwelveDataKey: strings.TrimSpace(a.secrets.TwelveData) != "", HasMarketauxKey: strings.TrimSpace(a.secrets.Marketaux) != "", FinnhubKeyHint: keyHint(a.secrets.Finnhub), AlpacaKeyHint: keyHint(a.secrets.AlpacaKey), GroqKeyHint: keyHint(a.secrets.Groq), OpenRouterKeyHint: keyHint(a.secrets.OpenRouter), GeminiKeyHint: keyHint(a.secrets.Gemini), FREDKeyHint: keyHint(a.secrets.FRED), BLSKeyHint: keyHint(a.secrets.BLS), EIAKeyHint: keyHint(a.secrets.EIA), TwelveDataKeyHint: keyHint(a.secrets.TwelveData), MarketauxKeyHint: keyHint(a.secrets.Marketaux), Watchlists: view.Watchlists, UI: view.UI, ProviderStatus: clone(a.state.ProviderStatus), SettingsSavedAt: a.state.SettingsSavedAt, MaintenanceLastRun: a.state.MaintenanceLastRun, LastCacheCleared: a.state.LastCacheCleared, CacheInfo: a.cacheInfoLocked(), ConfigDir: a.configDir, CachePath: a.cachePath()}
}

func clone[T any](v T) T {
	data, _ := json.Marshal(v)
	var out T
	_ = json.Unmarshal(data, &out)
	return out
}
