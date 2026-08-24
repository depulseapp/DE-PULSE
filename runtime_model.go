package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Hub struct {
	mu      sync.Mutex
	clients map[chan []byte]string
}

func NewHub() *Hub { return &Hub{clients: make(map[chan []byte]string)} }

// Subscribe keeps the historical all-events behavior for internal/tests.
// Browser sessions use SubscribeForUser so personal state events cannot cross
// authenticated workspace boundaries.
func (h *Hub) Subscribe() chan []byte { return h.SubscribeForUser("") }
func (h *Hub) SubscribeForUser(userID string) chan []byte {
	ch := make(chan []byte, 128)
	h.mu.Lock()
	h.clients[ch] = userID
	h.mu.Unlock()
	return ch
}

func (h *Hub) Unsubscribe(ch chan []byte) {
	h.mu.Lock()
	delete(h.clients, ch)
	close(ch)
	h.mu.Unlock()
}

func (h *Hub) Broadcast(v any) {

	h.mu.Lock()
	hasClients := len(h.clients) > 0
	h.mu.Unlock()
	if !hasClients {
		return
	}
	payload, err := json.Marshal(v)
	if err != nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.clients {
		select {
		case ch <- payload:
		default:

		}
	}
}

func (h *Hub) BroadcastUser(userID string, v any) {
	payload, err := json.Marshal(v)
	if err != nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch, clientUserID := range h.clients {
		if clientUserID != "" && clientUserID != userID {
			continue
		}
		select {
		case ch <- payload:
		default:
		}
	}
}

type Application struct {
	mu         sync.RWMutex
	state      AppState
	secrets    Secrets
	configDir  string
	hub        *Hub
	engine     *Engine
	identity   *IdentityService
	workspaces map[string]UserWorkspace
	// sessionKey is retained only for legacy in-package test fixtures that construct Application directly.
	// NewApplication never sets it; production authentication always uses IdentityService.
	sessionKey    string
	server        *http.Server
	aiCache       map[string]aiCacheEntry
	persistence   *PersistenceManager
	httpTelemetry *RequestTelemetry
}

type Engine struct {
	app                        *Application
	mu                         sync.RWMutex
	status                     string
	mode                       string
	message                    string
	startedAt                  string
	lastError                  string
	health                     map[string]string
	lastUpdated                map[string]int64
	quotes                     map[string]Quote
	history                    map[string][]HistoryPoint
	bars                       map[string]map[string][]Bar
	fundamentals               map[string]FundamentalSnapshot
	news                       []NewsItem
	earnings                   []EarningsItem
	filings                    []FilingItem
	secIntelligence            map[string]SECIntelligenceSummary
	scanner                    ScannerState
	globalDirect               map[string]GlobalDriver
	macroMetrics               map[string]MacroMetric
	macroEvents                []MacroEvent
	eventReactions             []EventReaction
	options                    map[string]OptionsContext
	capabilities               []CapabilityStatus
	signalValidation           SignalValidationState
	preparations               map[string]PreparationJobStatus
	symbolIntelligence         map[string]SymbolIntelligence
	catalystReactions          map[string]CatalystReactionState
	marketOpenFlags            map[string][]string
	marketOpenCheckpoint       MarketOpenCheckpoint
	alpacaCalendar             map[string]AlpacaCalendarDay
	marketActivity             MarketActivityState
	corporateActions           []CorporateAction
	manualActions              map[string]ManualActionStatus
	livePriorityHints          map[string]int64
	cancel                     context.CancelFunc
	wg                         sync.WaitGroup
	ws                         *WSClient
	alpacaWS                   *WSClient
	lastBroadcast              map[string]time.Time
	vixProviderSymbol          string
	vixHistoryAt               int64
	webSocketConnected         bool
	subscribedSymbols          map[string]bool
	alpacaWebSocketConnected   bool
	alpacaSubscribedSymbols    map[string]bool
	lastAlpacaStreamAt         int64
	lastAlpacaStreamSymbol     string
	lastMessageAt              int64
	lastTradeAt                int64
	lastTradeSymbol            string
	wsConnectedAt              int64
	lastAlpacaAt               int64
	lastAlpacaSymbol           string
	lastAlpacaFeed             string
	finnhubRateMu              sync.Mutex
	finnhubLastRequest         time.Time
	historyRefreshMu           sync.Mutex
	lastCacheHash              [32]byte
	lastCacheHashSet           bool
	providerCircuits           map[string]providerCircuit
	providerCapabilityCircuits map[string]providerCircuit
	providerCapabilityStates   map[string]ProviderCapabilityStateRecord
	smartRouterScorecard       SmartRouterScorecard
	rapidMoveEvents            map[string]RapidMoveEvent
	rapidMoveRecent            []RapidMoveEvent
	rapidMoveScorecard         RapidMoveScorecard
	providerQuotes             map[string]map[string]Quote
	rawHistoryCoverage         map[string]RawHistoryCoverage
	liquidityBaselines         map[string]LiquidityBaseline
	canonicalUSUniverse        []string
	canonicalUSUniverseAt      int64
	canonicalUSUniverseSource  string
	canonicalUSUniverseRetryAt int64
	canonicalUSUniverseRefresh chan struct{}
	canonicalUSUniverseLoader  func(context.Context, string, string) ([]string, bool)
	instrumentIdentities       map[string]InstrumentIdentityRecord
	instrumentIdentitiesLoaded bool
	radarCursor                int
	runtimeLoad                RuntimeLoadDiagnostics
	runtimeSLOTracker          *RuntimeSLOTracker
	startupProfile             RuntimeStartupDiagnostics
	lastRuntimeCPUTotal        float64
	lastRuntimeCPUSampledAt    time.Time
	initialPersistenceBytes    int64
	canonicalPipeline          CanonicalPipelineDiagnostics
	workload                   *WorkloadController
	providerTelemetry          *ProviderTelemetry
	providerCallsAvoided       int64
}

type MarketCache struct {
	Quotes                   map[string]Quote                         `json:"quotes"`
	History                  map[string][]HistoryPoint                `json:"history"`
	Bars                     map[string]map[string][]Bar              `json:"bars"`
	Fundamentals             map[string]FundamentalSnapshot           `json:"fundamentals"`
	News                     []NewsItem                               `json:"news,omitempty"`
	Earnings                 []EarningsItem                           `json:"earnings,omitempty"`
	Filings                  []FilingItem                             `json:"filings,omitempty"`
	SECIntelligence          map[string]SECIntelligenceSummary        `json:"secIntelligence,omitempty"`
	Scanner                  ScannerState                             `json:"scanner,omitempty"`
	GlobalDirect             map[string]GlobalDriver                  `json:"globalDirect,omitempty"`
	MacroMetrics             map[string]MacroMetric                   `json:"macroMetrics,omitempty"`
	MacroEvents              []MacroEvent                             `json:"macroEvents,omitempty"`
	EventReactions           []EventReaction                          `json:"eventReactions,omitempty"`
	Options                  map[string]OptionsContext                `json:"options,omitempty"`
	SignalValidation         SignalValidationState                    `json:"signalValidation,omitempty"`
	Preparations             map[string]PreparationJobStatus          `json:"preparations,omitempty"`
	CatalystReactions        map[string]CatalystReactionState         `json:"catalystReactions,omitempty"`
	MarketOpenFlags          map[string][]string                      `json:"marketOpenFlags,omitempty"`
	MarketOpenCheckpoint     MarketOpenCheckpoint                     `json:"marketOpenCheckpoint,omitempty"`
	CorporateActions         []CorporateAction                        `json:"corporateActions,omitempty"`
	RawHistoryCoverage       map[string]RawHistoryCoverage            `json:"rawHistoryCoverage,omitempty"`
	LiquidityBaselines       map[string]LiquidityBaseline             `json:"liquidityBaselines,omitempty"`
	ProviderCapabilityStates map[string]ProviderCapabilityStateRecord `json:"providerCapabilityStates,omitempty"`
	RapidMoveEvents          map[string]RapidMoveEvent                `json:"rapidMoveEvents,omitempty"`
	RapidMoveRecent          []RapidMoveEvent                         `json:"rapidMoveRecent,omitempty"`
	RapidMoveScorecard       RapidMoveScorecard                       `json:"rapidMoveScorecard,omitempty"`
	LastUpdated              map[string]int64                         `json:"lastUpdated,omitempty"`
	SavedAt                  int64                                    `json:"savedAt"`
}
