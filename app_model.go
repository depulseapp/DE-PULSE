package main

type Settings struct {
	DataMode                string  `json:"dataMode"`
	AIProvider              string  `json:"aiProvider"`
	AIRoutingMode           string  `json:"aiRoutingMode"`
	GroqModel               string  `json:"groqModel"`
	OpenRouterMode          string  `json:"openRouterMode"`
	OpenRouterSpecificModel string  `json:"openRouterSpecificModel"`
	GeminiModel             string  `json:"geminiModel"`
	SECEmail                string  `json:"secEmail"`
	AutoStart               bool    `json:"autoStart"`
	SignalProfile           string  `json:"signalProfile"`
	MarketContext           float64 `json:"marketContext"`
	EarningsPenalty         float64 `json:"earningsPenalty"`
	SwingEnabled            bool    `json:"swingEnabled"`
	DayEnabled              bool    `json:"dayEnabled"`
	LongEnabled             bool    `json:"longEnabled"`
	OvernightDataMode       string  `json:"overnightDataMode"`
	SwingWatchlistID        string  `json:"swingWatchlistId"`
	DayWatchlistID          string  `json:"dayWatchlistId"`
	LongWatchlistID         string  `json:"longWatchlistId"`
	DiscoveryWatchlistID    string  `json:"discoveryWatchlistId"`
	GlobalProviderMode      string  `json:"globalProviderMode"`
	OptionsDataMode         string  `json:"optionsDataMode"`
	MacroEventModeEnabled   bool    `json:"macroEventModeEnabled"`
	ResearchAIMode          string  `json:"researchAiMode"` // manual (default) or automatic
}

type Secrets struct {
	Finnhub      string `json:"finnhub"`
	TradeInsight string `json:"tradeInsight,omitempty"`
	AlpacaKey    string `json:"alpacaKey"`
	AlpacaSecret string `json:"alpacaSecret"`
	Groq         string `json:"groq"`
	OpenRouter   string `json:"openrouter"`
	Gemini       string `json:"gemini"`
	FRED         string `json:"fred"`
	BLS          string `json:"bls,omitempty"`
	EIA          string `json:"eia"`
	TwelveData   string `json:"twelveData"`
	Marketaux    string `json:"marketaux,omitempty"`
	MarketData   string `json:"marketData,omitempty"`
	// Legacy only: retained so an older secrets.json can be read without data loss.
	OpenAI string `json:"openai,omitempty"`
}

type Watchlist struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Symbols []string `json:"symbols"`
}

type UIState struct {
	ScopeType      string `json:"scopeType"`
	WatchlistID    string `json:"watchlistId"`
	SelectedTicker string `json:"selectedTicker"`
}

// UserWorkspace owns only personal market state. Provider credentials, data
// routing, deterministic scoring policy, runtime lifecycle and canonical market
// intelligence remain shared application responsibilities. This keeps v18.1
// multi-user isolation from multiplying provider calls or scoring owners.
type UserWorkspace struct {
	Version    int         `json:"version"`
	UserID     string      `json:"userId"`
	Watchlists []Watchlist `json:"watchlists"`
	UI         UIState     `json:"ui"`
	UpdatedAt  int64       `json:"updatedAt"`
}

type AppState struct {
	Version            int                           `json:"version"`
	Settings           Settings                      `json:"settings"`
	Watchlists         []Watchlist                   `json:"watchlists"`
	UI                 UIState                       `json:"ui"`
	CommunityEvidence  []CommunityEvidenceItem       `json:"communityEvidence,omitempty"`
	ProviderStatus     map[string]ProviderTestResult `json:"providerStatus,omitempty"`
	SettingsSavedAt    int64                         `json:"settingsSavedAt,omitempty"`
	MaintenanceLastRun int64                         `json:"maintenanceLastRun,omitempty"`
	LastCacheCleared   int64                         `json:"lastCacheCleared,omitempty"`
}

type PublicState struct {
	Version            string                        `json:"version"`
	BuildID            string                        `json:"buildId"`
	Settings           Settings                      `json:"settings"`
	HasFinnhubKey      bool                          `json:"hasFinnhubKey"`
	HasTradeInsightKey bool                          `json:"hasTradeInsightKey"`
	HasAlpacaKey       bool                          `json:"hasAlpacaKey"`
	HasAlpacaSecret    bool                          `json:"hasAlpacaSecret"`
	HasGroqKey         bool                          `json:"hasGroqKey"`
	HasGeminiKey       bool                          `json:"hasGeminiKey"`
	HasOpenRouterKey   bool                          `json:"hasOpenRouterKey"`
	HasFREDKey         bool                          `json:"hasFREDKey"`
	HasBLSKey          bool                          `json:"hasBLSKey"`
	HasEIAKey          bool                          `json:"hasEIAKey"`
	HasTwelveDataKey   bool                          `json:"hasTwelveDataKey"`
	HasMarketauxKey    bool                          `json:"hasMarketauxKey"`
	FinnhubKeyHint     string                        `json:"finnhubKeyHint,omitempty"`
	AlpacaKeyHint      string                        `json:"alpacaKeyHint,omitempty"`
	GroqKeyHint        string                        `json:"groqKeyHint,omitempty"`
	GeminiKeyHint      string                        `json:"geminiKeyHint,omitempty"`
	OpenRouterKeyHint  string                        `json:"openRouterKeyHint,omitempty"`
	FREDKeyHint        string                        `json:"fredKeyHint,omitempty"`
	BLSKeyHint         string                        `json:"blsKeyHint,omitempty"`
	EIAKeyHint         string                        `json:"eiaKeyHint,omitempty"`
	TwelveDataKeyHint  string                        `json:"twelveDataKeyHint,omitempty"`
	MarketauxKeyHint   string                        `json:"marketauxKeyHint,omitempty"`
	Watchlists         []Watchlist                   `json:"watchlists"`
	UI                 UIState                       `json:"ui"`
	ProviderStatus     map[string]ProviderTestResult `json:"providerStatus,omitempty"`
	SettingsSavedAt    int64                         `json:"settingsSavedAt,omitempty"`
	MaintenanceLastRun int64                         `json:"maintenanceLastRun,omitempty"`
	LastCacheCleared   int64                         `json:"lastCacheCleared,omitempty"`
	CacheInfo          CacheInfo                     `json:"cacheInfo"`
	ConfigDir          string                        `json:"configDir,omitempty"`
	CachePath          string                        `json:"cachePath,omitempty"`
}

type Quote struct {
	Symbol            string  `json:"symbol"`
	Price             float64 `json:"price,omitempty"`
	Change            float64 `json:"change,omitempty"`
	ChangePercent     float64 `json:"changePercent,omitempty"`
	Open              float64 `json:"open,omitempty"`
	High              float64 `json:"high,omitempty"`
	Low               float64 `json:"low,omitempty"`
	PreviousClose     float64 `json:"previousClose,omitempty"`
	SessionClose      float64 `json:"sessionClose,omitempty"`
	SessionCloseAt    int64   `json:"sessionCloseAt,omitempty"`
	PriorSessionClose float64 `json:"priorSessionClose,omitempty"`
	Volume            float64 `json:"volume,omitempty"`
	Bid               float64 `json:"bid,omitempty"`
	Ask               float64 `json:"ask,omitempty"`
	BidSize           float64 `json:"bidSize,omitempty"`
	AskSize           float64 `json:"askSize,omitempty"`
	ProviderTimestamp int64   `json:"providerTimestamp,omitempty"`
	UpdatedAt         int64   `json:"updatedAt,omitempty"`
	Source            string  `json:"source,omitempty"`
	FeedType          string  `json:"feedType,omitempty"`
	DataState         string  `json:"dataState,omitempty"`
}

type Bar struct {
	T int64   `json:"t"`
	O float64 `json:"o"`
	H float64 `json:"h"`
	L float64 `json:"l"`
	C float64 `json:"c"`
	V float64 `json:"v"`
}

type FundamentalSnapshot struct {
	Symbol           string  `json:"symbol"`
	MarketCap        float64 `json:"marketCap,omitempty"`
	PERatio          float64 `json:"peRatio,omitempty"`
	ForwardPERatio   float64 `json:"forwardPeRatio,omitempty"`
	PSRatio          float64 `json:"psRatio,omitempty"`
	PEGRatio         float64 `json:"pegRatio,omitempty"`
	RevenueGrowth    float64 `json:"revenueGrowth,omitempty"`
	EPSGrowth        float64 `json:"epsGrowth,omitempty"`
	GrossMargin      float64 `json:"grossMargin,omitempty"`
	OperatingMargin  float64 `json:"operatingMargin,omitempty"`
	ROE              float64 `json:"roe,omitempty"`
	NetMargin        float64 `json:"netMargin,omitempty"`
	DebtToEquity     float64 `json:"debtToEquity,omitempty"`
	CurrentRatio     float64 `json:"currentRatio,omitempty"`
	FreeCashFlow     float64 `json:"freeCashFlow,omitempty"`
	DividendYield    float64 `json:"dividendYield,omitempty"`
	FiftyTwoWeekHigh float64 `json:"fiftyTwoWeekHigh,omitempty"`
	FiftyTwoWeekLow  float64 `json:"fiftyTwoWeekLow,omitempty"`
	UpdatedAt        int64   `json:"updatedAt,omitempty"`
	Source           string  `json:"source,omitempty"`
}

type ScannerResult struct {
	Symbol                string   `json:"symbol"`
	Mode                  string   `json:"mode"`
	Price                 float64  `json:"price"`
	ChangePercent         float64  `json:"changePercent"`
	GapPercent            float64  `json:"gapPercent,omitempty"`
	RelativeVolume        float64  `json:"relativeVolume,omitempty"`
	SessionRelativeVolume float64  `json:"sessionRelativeVolume,omitempty"`
	DollarVolume          float64  `json:"dollarVolume,omitempty"`
	SpreadPercent         float64  `json:"spreadPercent,omitempty"`
	RangeExpansion        float64  `json:"rangeExpansion,omitempty"`
	UnusualVolumeScore    float64  `json:"unusualVolumeScore,omitempty"`
	VolatilityScore       float64  `json:"volatilityScore,omitempty"`
	OpportunityScore      float64  `json:"opportunityScore,omitempty"`
	PriceConfirmation     string   `json:"priceConfirmation,omitempty"`
	RSI                   float64  `json:"rsi,omitempty"`
	TrendScore            float64  `json:"trendScore,omitempty"`
	MomentumScore         float64  `json:"momentumScore,omitempty"`
	RelativeStrength      float64  `json:"relativeStrengthPct,omitempty"`
	RSBenchmark           string   `json:"relativeStrengthBenchmark,omitempty"`
	FundamentalScore      float64  `json:"fundamentalScore,omitempty"`
	Score                 float64  `json:"score"`
	Reasons               []string `json:"reasons,omitempty"`
	Provider              string   `json:"provider"`
	UpdatedAt             int64    `json:"updatedAt"`
}

type OpportunityPromotion struct {
	Symbol           string   `json:"symbol"`
	Score            float64  `json:"score"`
	State            string   `json:"state"`
	Reasons          []string `json:"reasons,omitempty"`
	PromotedAt       int64    `json:"promotedAt"`
	LastConfirmedAt  int64    `json:"lastConfirmedAt"`
	ExpiresAt        int64    `json:"expiresAt"`
	ShadowWouldMatch bool     `json:"shadowWouldMatch,omitempty"`
}

type OpportunityRadarState struct {
	Status          string                 `json:"status"`
	Session         string                 `json:"session"`
	Message         string                 `json:"message"`
	Candidates      []ScannerResult        `json:"candidates"`
	Promotions      []OpportunityPromotion `json:"promotions"`
	Scanned         int                    `json:"scanned"`
	DurationMs      int64                  `json:"durationMs"`
	CadenceMs       int64                  `json:"cadenceMs"`
	LastRun         int64                  `json:"lastRun"`
	NextRun         int64                  `json:"nextRun"`
	Provider        string                 `json:"provider"`
	ProductionFloor float64                `json:"productionFloor"`
	ShadowFloor     float64                `json:"shadowFloor"`
	ShadowOnly      bool                   `json:"shadowOnly"`
}

type ScannerState struct {
	Mode       string                `json:"mode"`
	Status     string                `json:"status"`
	Message    string                `json:"message"`
	Results    []ScannerResult       `json:"results"`
	Scanned    int                   `json:"scanned"`
	DurationMs int64                 `json:"durationMs"`
	UpdatedAt  int64                 `json:"updatedAt"`
	Radar      OpportunityRadarState `json:"radar"`
}

type CacheInfo struct {
	SizeBytes     int64 `json:"sizeBytes"`
	CachedSymbols int   `json:"cachedSymbols"`
	LastUpdated   int64 `json:"lastUpdated"`
}

type ProviderTestResult struct {
	Provider  string   `json:"provider"`
	OK        bool     `json:"ok"`
	Status    string   `json:"status"`
	Message   string   `json:"message"`
	Details   []string `json:"details,omitempty"`
	CheckedAt string   `json:"checkedAt"`
}

type CapabilityStatus struct {
	Capability string   `json:"capability"`
	Source     string   `json:"source"`
	Mode       string   `json:"mode"`
	Status     string   `json:"status"`
	Freshness  string   `json:"freshness,omitempty"`
	UpdatedAt  int64    `json:"updatedAt,omitempty"`
	Details    []string `json:"details,omitempty"`
}

type GlobalDriver struct {
	Key            string  `json:"key"`
	Label          string  `json:"label"`
	State          string  `json:"state"`
	Value          float64 `json:"value,omitempty"`
	ChangePercent  float64 `json:"changePercent,omitempty"`
	Source         string  `json:"source"`
	Provenance     string  `json:"provenance"`
	UpdatedAt      int64   `json:"updatedAt,omitempty"`
	Confidence     int     `json:"confidence"`
	Detail         string  `json:"detail,omitempty"`
	Session        string  `json:"session,omitempty"`
	Underlying     string  `json:"underlying,omitempty"`
	ProviderSymbol string  `json:"providerSymbol,omitempty"`
	IsProxy        bool    `json:"isProxy,omitempty"`
}

type GlobalMarketContext struct {
	Tone       string                  `json:"tone"`
	Confidence int                     `json:"confidence"`
	Drivers    map[string]GlobalDriver `json:"drivers"`
	UpdatedAt  int64                   `json:"updatedAt,omitempty"`
	Mode       string                  `json:"mode"`
	Summary    string                  `json:"summary"`
}

type MacroMetric struct {
	Key        string  `json:"key"`
	Label      string  `json:"label"`
	Value      float64 `json:"value"`
	Unit       string  `json:"unit,omitempty"`
	Source     string  `json:"source"`
	Provenance string  `json:"provenance"`
	UpdatedAt  int64   `json:"updatedAt,omitempty"`
	Status     string  `json:"status"`
	Previous   float64 `json:"previous,omitempty"`
	Change5D   float64 `json:"change5d,omitempty"`
	Change20D  float64 `json:"change20d,omitempty"`
	Change1M   float64 `json:"change1m,omitempty"`
	Change3M   float64 `json:"change3m,omitempty"`
	Period     string  `json:"period,omitempty"`
}

type MacroEvent struct {
	ID              string   `json:"id"`
	Region          string   `json:"region"`
	Name            string   `json:"name"`
	Impact          string   `json:"impact"`
	ProcessingClass string   `json:"processingClass,omitempty"`
	Lifecycle       string   `json:"lifecycle"`
	StartsAt        int64    `json:"startsAt,omitempty"`
	Date            string   `json:"date"`
	TimeKnown       bool     `json:"timeKnown"`
	Expected        *float64 `json:"expected,omitempty"`
	Actual          *float64 `json:"actual,omitempty"`
	Previous        *float64 `json:"previous,omitempty"`
	Unit            string   `json:"unit,omitempty"`
	Source          string   `json:"source"`
	SourceURL       string   `json:"sourceUrl,omitempty"`
	UpdatedAt       int64    `json:"updatedAt,omitempty"`
}

type EventReaction struct {
	EventID           string             `json:"eventId"`
	OffsetSec         int                `json:"offsetSec"`
	CapturedAt        int64              `json:"capturedAt"`
	BaselineAt        int64              `json:"baselineAt,omitempty"`
	Moves             map[string]float64 `json:"moves"`
	Baseline          map[string]float64 `json:"baseline,omitempty"`
	InternalLatencyMs int64              `json:"internalLatencyMs,omitempty"`
}

type EventModeState struct {
	Active            bool     `json:"active"`
	EventID           string   `json:"eventId,omitempty"`
	Name              string   `json:"name,omitempty"`
	StartsAt          int64    `json:"startsAt,omitempty"`
	CountdownS        int64    `json:"countdownS,omitempty"`
	Phase             string   `json:"phase,omitempty"`
	AffectedSymbols   []string `json:"affectedSymbols,omitempty"`
	AffectedSectors   []string `json:"affectedSectors,omitempty"`
	Prepared          bool     `json:"prepared,omitempty"`
	QueuePrepared     bool     `json:"queuePrepared,omitempty"`
	MetadataReady     bool     `json:"metadataReady,omitempty"`
	ExpectationStatus string   `json:"expectationStatus,omitempty"`
	DataPriority      string   `json:"dataPriority,omitempty"`
}

type GEXStrikeLevel struct {
	Strike       float64 `json:"strike"`
	CallGEX      float64 `json:"callGex,omitempty"`
	PutGEX       float64 `json:"putGex,omitempty"`
	NetGEX       float64 `json:"netGex"`
	AbsoluteGEX  float64 `json:"absoluteGex"`
	OpenInterest float64 `json:"openInterest,omitempty"`
	Contracts    int     `json:"contracts"`
}

type GEXExpirationLevel struct {
	Expiration   string  `json:"expiration"`
	CallGEX      float64 `json:"callGex,omitempty"`
	PutGEX       float64 `json:"putGex,omitempty"`
	NetGEX       float64 `json:"netGex"`
	AbsoluteGEX  float64 `json:"absoluteGex"`
	OpenInterest float64 `json:"openInterest,omitempty"`
	Contracts    int     `json:"contracts"`
}

type GEXConcentrationZone struct {
	LowStrike   float64 `json:"lowStrike"`
	HighStrike  float64 `json:"highStrike"`
	NetGEX      float64 `json:"netGex"`
	AbsoluteGEX float64 `json:"absoluteGex"`
	SharePct    float64 `json:"sharePct,omitempty"`
	StrikeCount int     `json:"strikeCount"`
}

type OptionsContext struct {
	Symbol                string                 `json:"symbol"`
	Provider              string                 `json:"provider"`
	Feed                  string                 `json:"feed"`
	State                 string                 `json:"state"`
	Bias                  string                 `json:"bias"`
	CallContracts         int                    `json:"callContracts"`
	PutContracts          int                    `json:"putContracts"`
	CallVolume            float64                `json:"callVolume,omitempty"`
	PutVolume             float64                `json:"putVolume,omitempty"`
	PutCallVolume         float64                `json:"putCallVolume,omitempty"`
	AverageIV             float64                `json:"averageIv,omitempty"`
	IVChange              float64                `json:"ivChange,omitempty"` // percentage-point change vs prior real refresh
	ExpectedMove          float64                `json:"expectedMove,omitempty"`
	NearestExpiration     string                 `json:"nearestExpiration,omitempty"`
	UpdatedAt             int64                  `json:"updatedAt,omitempty"`
	Provenance            string                 `json:"provenance"`
	Limitations           []string               `json:"limitations,omitempty"`
	GammaContracts        int                    `json:"gammaContracts,omitempty"`
	OpenInterestContracts int                    `json:"openInterestContracts,omitempty"`
	GammaOIContracts      int                    `json:"gammaOiContracts,omitempty"`
	GammaOICoveragePct    float64                `json:"gammaOiCoveragePct,omitempty"`
	OpenInterestDate      string                 `json:"openInterestDate,omitempty"`
	CallGEX               float64                `json:"callGex,omitempty"`
	PutGEX                float64                `json:"putGex,omitempty"`
	NetGEX                float64                `json:"netGex,omitempty"`
	GEXState              string                 `json:"gexState,omitempty"`
	GEXQuality            string                 `json:"gexQuality,omitempty"`
	UnderlyingPrice       float64                `json:"underlyingPrice,omitempty"`
	MajorGammaStrikes     []GEXStrikeLevel       `json:"majorGammaStrikes,omitempty"`
	GammaZones            []GEXConcentrationZone `json:"gammaZones,omitempty"`
	ExpirationGEX         []GEXExpirationLevel   `json:"expirationGex,omitempty"`
	GammaFlip             *float64               `json:"gammaFlip,omitempty"`
	GammaFlipMethod       string                 `json:"gammaFlipMethod,omitempty"`
}

type SignalSnapshot struct {
	ID                      string             `json:"id"`
	Symbol                  string             `json:"symbol"`
	Horizon                 string             `json:"horizon"`
	Timestamp               int64              `json:"timestamp"`
	Price                   float64            `json:"price"`
	Score                   float64            `json:"score"`
	Action                  string             `json:"action"`
	Readiness               string             `json:"readiness"`
	EvidenceSnapshotID      string             `json:"evidenceSnapshotId,omitempty"`
	FormulaVersion          string             `json:"formulaVersion,omitempty"`
	SettingsFingerprint     string             `json:"settingsFingerprint,omitempty"`
	EarningsPenalty         float64            `json:"earningsPenalty,omitempty"`
	SignalProfile           string             `json:"signalProfile,omitempty"`
	FamilyScores            map[string]float64 `json:"familyScores,omitempty"`
	EarningsDays            *int               `json:"earningsDays,omitempty"`
	EntryLow                float64            `json:"entryLow,omitempty"`
	EntryHigh               float64            `json:"entryHigh,omitempty"`
	TargetLow               float64            `json:"targetLow,omitempty"`
	TargetHigh              float64            `json:"targetHigh,omitempty"`
	Invalidation            float64            `json:"invalidation,omitempty"`
	MarketRegime            string             `json:"marketRegime,omitempty"`
	MarketStructure         string             `json:"marketStructure,omitempty"`
	MarketTradeability      string             `json:"marketTradeability,omitempty"`
	RelativeStrength        string             `json:"relativeStrength,omitempty"`
	SectorRegime            string             `json:"sectorRegime,omitempty"`
	LiquidityState          string             `json:"liquidityState,omitempty"`
	MarketIntelAt           int64              `json:"marketIntelAt,omitempty"`
	EventIntelAt            int64              `json:"eventIntelAt,omitempty"`
	GlobalContext           string             `json:"globalContext,omitempty"`
	OptionsBias             string             `json:"optionsBias,omitempty"`
	EventRisk               string             `json:"eventRisk,omitempty"`
	ResearchState           string             `json:"researchState,omitempty"`
	QueuePriority           string             `json:"queuePriority,omitempty"`
	KeyDriver               string             `json:"keyDriver,omitempty"`
	Contradictions          []string           `json:"contradictions,omitempty"`
	Outcomes                map[string]float64 `json:"outcomes,omitempty"`
	MFE                     float64            `json:"mfe,omitempty"`
	MAE                     float64            `json:"mae,omitempty"`
	OutcomeState            string             `json:"outcomeState,omitempty"`
	OutcomeDetail           string             `json:"outcomeDetail,omitempty"`
	OutcomeUpdatedAt        int64              `json:"outcomeUpdatedAt,omitempty"`
	EntryTouched            bool               `json:"entryTouched,omitempty"`
	EntryTouchedAt          int64              `json:"entryTouchedAt,omitempty"`
	TargetTouchedAt         int64              `json:"targetTouchedAt,omitempty"`
	InvalidationAt          int64              `json:"invalidationAt,omitempty"`
	ElapsedMinutes          int64              `json:"elapsedMinutes,omitempty"`
	OutcomeAdjustmentFactor float64            `json:"outcomeAdjustmentFactor,omitempty"`
	OutcomeAdjustmentDetail string             `json:"outcomeAdjustmentDetail,omitempty"`
	LegacyPartial           bool               `json:"legacyPartial,omitempty"`
}

type SignalValidationState struct {
	Snapshots []SignalSnapshot `json:"snapshots"`
	UpdatedAt int64            `json:"updatedAt,omitempty"`
	Message   string           `json:"message,omitempty"`
}

type SeasonalityMetric struct {
	Key                    string   `json:"key"`
	Label                  string   `json:"label"`
	State                  string   `json:"state"`
	SampleCount            int      `json:"sampleCount"`
	HistoricalYears        int      `json:"historicalYears,omitempty"`
	AverageReturnPct       *float64 `json:"averageReturnPct,omitempty"`
	MedianReturnPct        *float64 `json:"medianReturnPct,omitempty"`
	PositiveFrequencyPct   *float64 `json:"positiveFrequencyPct,omitempty"`
	BestReturnPct          *float64 `json:"bestReturnPct,omitempty"`
	WorstReturnPct         *float64 `json:"worstReturnPct,omitempty"`
	CurrentYearReturnPct   *float64 `json:"currentYearReturnPct,omitempty"`
	CurrentYearObservation string   `json:"currentYearObservation,omitempty"`
	DateFrom               string   `json:"dateFrom,omitempty"`
	DateTo                 string   `json:"dateTo,omitempty"`
	Detail                 string   `json:"detail,omitempty"`
}

type SeasonalitySymbolState struct {
	Symbol      string              `json:"symbol"`
	State       string              `json:"state"`
	Source      string              `json:"source"`
	SampleFrom  string              `json:"sampleFrom,omitempty"`
	SampleTo    string              `json:"sampleTo,omitempty"`
	DailyBars   int                 `json:"dailyBars"`
	Monthly     []SeasonalityMetric `json:"monthly"`
	DayOfWeek   []SeasonalityMetric `json:"dayOfWeek"`
	Limitations []string            `json:"limitations,omitempty"`
}

type SeasonalitySnapshot struct {
	Symbols   map[string]SeasonalitySymbolState `json:"symbols"`
	UpdatedAt int64                             `json:"updatedAt,omitempty"`
	Message   string                            `json:"message,omitempty"`
}

type CalibrationGroup struct {
	Horizon          string   `json:"horizon"`
	ScoreBand        string   `json:"scoreBand"`
	Action           string   `json:"action"`
	State            string   `json:"state"`
	SampleCount      int      `json:"sampleCount"`
	TargetReached    int      `json:"targetReached"`
	Invalidated      int      `json:"invalidated"`
	NoEntry          int      `json:"noEntry"`
	Ambiguous        int      `json:"ambiguous"`
	TargetReachedPct *float64 `json:"targetReachedPct,omitempty"`
	Average5DReturn  *float64 `json:"average5DReturnPct,omitempty"`
	AverageMFE       *float64 `json:"averageMfePct,omitempty"`
	AverageMAE       *float64 `json:"averageMaePct,omitempty"`
	DateFrom         int64    `json:"dateFrom,omitempty"`
	DateTo           int64    `json:"dateTo,omitempty"`
	Detail           string   `json:"detail,omitempty"`
}

type CalibrationSnapshot struct {
	Groups                     []CalibrationGroup `json:"groups"`
	EligibleSnapshots          int                `json:"eligibleSnapshots"`
	SetupScoreIsWinProbability bool               `json:"setupScoreIsWinProbability"`
	UpdatedAt                  int64              `json:"updatedAt,omitempty"`
	Message                    string             `json:"message,omitempty"`
}

type CorrelationPair struct {
	SymbolA     string  `json:"symbolA"`
	SymbolB     string  `json:"symbolB"`
	State       string  `json:"state"`
	Correlation float64 `json:"correlation,omitempty"`
	SampleCount int     `json:"sampleCount"`
	WindowFrom  string  `json:"windowFrom,omitempty"`
	WindowTo    string  `json:"windowTo,omitempty"`
	Detail      string  `json:"detail,omitempty"`
}

type ConcentrationGroup struct {
	Key     string   `json:"key"`
	Kind    string   `json:"kind"`
	Symbols []string `json:"symbols"`
	State   string   `json:"state"`
	Detail  string   `json:"detail,omitempty"`
}

type CorrelationConcentrationSnapshot struct {
	Pairs          []CorrelationPair    `json:"pairs"`
	Concentrations []ConcentrationGroup `json:"concentrations"`
	CandidateCount int                  `json:"candidateCount"`
	UpdatedAt      int64                `json:"updatedAt,omitempty"`
	Message        string               `json:"message,omitempty"`
}

type ReplayScenarioDescriptor struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Label    string `json:"label"`
	Symbol   string `json:"symbol"`
	Horizon  string `json:"horizon"`
	Cutoff   int64  `json:"cutoff"`
	State    string `json:"state"`
	Source   string `json:"source"`
	Evidence string `json:"evidence"`
	Detail   string `json:"detail"`
}

type ReplayScenarioCatalog struct {
	Scenarios []ReplayScenarioDescriptor `json:"scenarios"`
	Available int                        `json:"available"`
	Kinds     []string                   `json:"kinds"`
	UpdatedAt int64                      `json:"updatedAt,omitempty"`
	Message   string                     `json:"message"`
}

type ValidationLearningSnapshot struct {
	Seasonality   SeasonalitySnapshot              `json:"seasonality"`
	Calibration   CalibrationSnapshot              `json:"calibration"`
	Concentration CorrelationConcentrationSnapshot `json:"concentration"`
	ReplayState   string                           `json:"replayState"`
	ReplayDetail  string                           `json:"replayDetail"`
	ReplayCatalog ReplayScenarioCatalog            `json:"replayCatalog"`
	UpdatedAt     int64                            `json:"updatedAt,omitempty"`
}

type HistoryPoint struct {
	T int64   `json:"t"`
	P float64 `json:"p"`
}

type NewsItem struct {
	ID       any      `json:"id"`
	Datetime int64    `json:"datetime"`
	Headline string   `json:"headline"`
	Summary  string   `json:"summary"`
	Source   string   `json:"source"`
	URL      string   `json:"url"`
	Related  string   `json:"related,omitempty"`
	Symbols  []string `json:"symbols,omitempty"`
	Scope    string   `json:"scope,omitempty"`
}

type EarningsItem struct {
	Symbol                  string   `json:"symbol"`
	Date                    string   `json:"date"`
	Hour                    string   `json:"hour"`
	EPSActual               *float64 `json:"epsActual,omitempty"`
	EPSEstimate             *float64 `json:"epsEstimate,omitempty"`
	RevenueActual           *float64 `json:"revenueActual,omitempty"`
	RevenueEstimate         *float64 `json:"revenueEstimate,omitempty"`
	Quarter                 int      `json:"quarter,omitempty"`
	Year                    int      `json:"year,omitempty"`
	GuidanceRevenuePrevLow  *float64 `json:"guidanceRevenuePrevLow,omitempty"`
	GuidanceRevenuePrevHigh *float64 `json:"guidanceRevenuePrevHigh,omitempty"`
	GuidanceRevenueNewLow   *float64 `json:"guidanceRevenueNewLow,omitempty"`
	GuidanceRevenueNewHigh  *float64 `json:"guidanceRevenueNewHigh,omitempty"`
	GuidanceEPSPrevLow      *float64 `json:"guidanceEpsPrevLow,omitempty"`
	GuidanceEPSPrevHigh     *float64 `json:"guidanceEpsPrevHigh,omitempty"`
	GuidanceEPSNewLow       *float64 `json:"guidanceEpsNewLow,omitempty"`
	GuidanceEPSNewHigh      *float64 `json:"guidanceEpsNewHigh,omitempty"`
	GuidanceSource          string   `json:"guidanceSource,omitempty"`
}

type SECInsiderTransaction struct {
	Actor           string  `json:"actor,omitempty"`
	Role            string  `json:"role,omitempty"`
	Classification  string  `json:"classification"`
	Code            string  `json:"code,omitempty"`
	Meaning         string  `json:"meaning,omitempty"`
	TransactionDate string  `json:"transactionDate,omitempty"`
	Shares          float64 `json:"shares,omitempty"`
	Price           float64 `json:"price,omitempty"`
	Value           float64 `json:"value,omitempty"`
	OwnershipAfter  float64 `json:"ownershipAfter,omitempty"`
	FiledAt         string  `json:"filedAt,omitempty"`
	URL             string  `json:"url,omitempty"`
}

type FilingItem struct {
	ID              string                  `json:"id"`
	Symbol          string                  `json:"symbol"`
	Company         string                  `json:"company"`
	Form            string                  `json:"form"`
	FiledAt         string                  `json:"filedAt"`
	ReportDate      string                  `json:"reportDate,omitempty"`
	Description     string                  `json:"description"`
	Meaning         string                  `json:"meaning,omitempty"`
	Category        string                  `json:"category,omitempty"`
	Signal          string                  `json:"signal,omitempty"`
	Actor           string                  `json:"actor,omitempty"`
	Role            string                  `json:"role,omitempty"`
	TransactionDate string                  `json:"transactionDate,omitempty"`
	TransactionCode string                  `json:"transactionCode,omitempty"`
	TransactionType string                  `json:"transactionType,omitempty"`
	Shares          float64                 `json:"shares,omitempty"`
	Price           float64                 `json:"price,omitempty"`
	Value           float64                 `json:"value,omitempty"`
	OwnershipAfter  float64                 `json:"ownershipAfter,omitempty"`
	Items           string                  `json:"items,omitempty"`
	URL             string                  `json:"url"`
	Transactions    []SECInsiderTransaction `json:"transactions,omitempty"`
}

type SECSignal struct {
	Label  string `json:"label"`
	Detail string `json:"detail"`
	Tone   string `json:"tone"`
	Date   string `json:"date,omitempty"`
	URL    string `json:"url,omitempty"`
}

type SECIntelligenceSummary struct {
	Symbol                    string                  `json:"symbol"`
	InsiderBuys               int                     `json:"insiderBuys"`
	InsiderSells              int                     `json:"insiderSells"`
	InsiderOthers             int                     `json:"insiderOthers"`
	InsiderBuyValue           float64                 `json:"insiderBuyValue,omitempty"`
	InsiderSellValue          float64                 `json:"insiderSellValue,omitempty"`
	OwnershipChanges          int                     `json:"ownershipChanges"`
	MaterialEvents            int                     `json:"materialEvents"`
	OfferingRisk              string                  `json:"offeringRisk"`
	Institutional             string                  `json:"institutional"`
	LatestForm                string                  `json:"latestForm,omitempty"`
	LatestFiledAt             string                  `json:"latestFiledAt,omitempty"`
	UpdatedAt                 int64                   `json:"updatedAt,omitempty"`
	Signals                   []SECSignal             `json:"signals,omitempty"`
	RecentTransactions        []FilingItem            `json:"recentTransactions,omitempty"`
	RecentInsiderTransactions []SECInsiderTransaction `json:"recentInsiderTransactions,omitempty"`
}

type FeedDiagnostics struct {
	WebSocketConnected       bool     `json:"webSocketConnected"`
	SubscribedSymbols        []string `json:"subscribedSymbols"`
	LastMessageAt            int64    `json:"lastMessageAt,omitempty"`
	LastTradeAt              int64    `json:"lastTradeAt,omitempty"`
	LastTradeSymbol          string   `json:"lastTradeSymbol,omitempty"`
	AlpacaWebSocketConnected bool     `json:"alpacaWebSocketConnected"`
	AlpacaSubscribedSymbols  []string `json:"alpacaSubscribedSymbols"`
	LastAlpacaStreamAt       int64    `json:"lastAlpacaStreamAt,omitempty"`
	LastAlpacaStreamSymbol   string   `json:"lastAlpacaStreamSymbol,omitempty"`
	LastAlpacaAt             int64    `json:"lastAlpacaAt,omitempty"`
	LastAlpacaSymbol         string   `json:"lastAlpacaSymbol,omitempty"`
	AlpacaLiveFeed           string   `json:"alpacaLiveFeed,omitempty"`
	OvernightDataMode        string   `json:"overnightDataMode,omitempty"`
	OvernightLiveAvailable   bool     `json:"overnightLiveAvailable"`
	MarketSession            string   `json:"marketSession"`
	SessionBoundaryAt        int64    `json:"sessionBoundaryAt,omitempty"`
	SessionBoundaryAction    string   `json:"sessionBoundaryAction,omitempty"`
	FeedState                string   `json:"feedState"`
	FinnhubMaxSymbols        int      `json:"finnhubMaxSymbols"`
	FinnhubReserveSlots      int      `json:"finnhubReserveSlots"`
	AlpacaMaxSymbols         int      `json:"alpacaMaxSymbols"`
	AlpacaReserveSlots       int      `json:"alpacaReserveSlots"`
	SnapshotSymbols          []string `json:"snapshotSymbols"`
	TrackedSymbols           int      `json:"trackedSymbols"`
	LiveSymbols              int      `json:"liveSymbols"`
}

type RuntimeSnapshot struct {
	Status                  string                                 `json:"status"`
	Mode                    string                                 `json:"mode"`
	Message                 string                                 `json:"message"`
	StartedAt               string                                 `json:"startedAt,omitempty"`
	Quotes                  map[string]Quote                       `json:"quotes"`
	History                 map[string][]HistoryPoint              `json:"history"`
	Bars                    map[string]map[string][]Bar            `json:"bars"`
	Fundamentals            map[string]FundamentalSnapshot         `json:"fundamentals"`
	News                    []NewsItem                             `json:"news"`
	Earnings                []EarningsItem                         `json:"earnings"`
	Filings                 []FilingItem                           `json:"filings"`
	SECIntelligence         map[string]SECIntelligenceSummary      `json:"secIntelligence"`
	Scanner                 ScannerState                           `json:"scanner"`
	Health                  map[string]string                      `json:"health"`
	LastUpdated             map[string]int64                       `json:"lastUpdated"`
	Feed                    FeedDiagnostics                        `json:"feed"`
	Global                  GlobalMarketContext                    `json:"global"`
	MacroMetrics            map[string]MacroMetric                 `json:"macroMetrics"`
	MacroEvents             []MacroEvent                           `json:"macroEvents"`
	EventMode               EventModeState                         `json:"eventMode"`
	EventReactions          []EventReaction                        `json:"eventReactions,omitempty"`
	Options                 map[string]OptionsContext              `json:"options"`
	Capabilities            []CapabilityStatus                     `json:"capabilities"`
	SignalValidation        SignalValidationState                  `json:"signalValidation"`
	ValidationLearning      ValidationLearningSnapshot             `json:"validationLearning"`
	Preparations            map[string]PreparationJobStatus        `json:"preparations"`
	Liquidity               map[string]LiquidityState              `json:"liquidity"`
	Intelligence            map[string]DerivedIntelligenceState    `json:"intelligence"`
	ProviderRegistry        []ProviderCapabilityEntry              `json:"providerRegistry"`
	SymbolIntelligence      map[string]SymbolIntelligence          `json:"symbolIntelligence"`
	CatalystReactions       map[string]CatalystReactionState       `json:"catalystReactions"`
	MarketOpenFlags         map[string][]string                    `json:"marketOpenFlags"`
	MarketOpenCheckpoint    MarketOpenCheckpoint                   `json:"marketOpenCheckpoint"`
	MarketActivity          MarketActivityState                    `json:"marketActivity"`
	CorporateActions        []CorporateAction                      `json:"corporateActions"`
	LiveCoverage            map[string]LiveCoverageState           `json:"liveCoverage"`
	ManualActions           map[string]ManualActionStatus          `json:"manualActions"`
	ProviderRouter          ProviderRouterSnapshot                 `json:"providerRouter"`
	RapidMove               RapidMoveState                         `json:"rapidMove"`
	ProviderReconciliation  []ProviderReconciliationDecision       `json:"providerReconciliation"`
	ResearchPackage         ResearchPackageTruth                   `json:"researchPackage"`
	EvidenceSnapshot        EvidenceSnapshot                       `json:"evidenceSnapshot"`
	CorporateActionTruth    CorporateActionTruth                   `json:"corporateActionTruth"`
	Freshness               []FreshnessDiagnostic                  `json:"freshness"`
	FreshnessSummary        FreshnessSummary                       `json:"freshnessSummary"`
	MarketIntelligence      MarketIntelligenceSnapshot             `json:"marketIntelligence"`
	EventIntelligence       EventIntelligenceSnapshot              `json:"eventIntelligence"`
	AlternativeIntelligence ContextAlternativeIntelligenceSnapshot `json:"alternativeIntelligence"`
	AdaptiveDataPolicy      AdaptiveDataPolicyState                `json:"adaptiveDataPolicy"`
	ShadowControl           ShadowControlState                     `json:"shadowControl"`
	RuntimeLoad             RuntimeLoadDiagnostics                 `json:"runtimeLoad"`
	Degradation             RuntimeDegradationState                `json:"degradation"`
	RuntimeSLO              RuntimeSLOAssessment                   `json:"runtimeSlo"`
	LastError               string                                 `json:"lastError,omitempty"`
}
