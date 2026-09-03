package main

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"sync"
	"time"
)

// PersistenceBackend is deliberately storage-agnostic. v17 desktop uses SQLite
// where the platform can provide it, while these contracts remain suitable for
// a PostgreSQL implementation in a later hosted major without changing callers.
type PersistenceBackend interface {
	Name() string
	Capabilities() []string
	Init(context.Context) error
	UpsertSymbols(context.Context, []SymbolRegistryRecord) (int, error)
	LoadSymbols(context.Context) ([]SymbolRegistryRecord, error)
	SaveQuotes(context.Context, map[string]Quote) (int, error)
	LoadQuotes(context.Context) (map[string]Quote, error)
	SaveIntelligence(context.Context, PersistenceIntelligenceBatch) (int, error)
	LoadIdentityState(context.Context) (IdentityPersistentState, error)
	SaveIdentityState(context.Context, IdentityPersistentState) error
	LoadUserWorkspaces(context.Context) ([]UserWorkspace, error)
	SaveUserWorkspace(context.Context, UserWorkspace) error
	Stats(context.Context) (PersistenceStoreStats, error)
	Close() error
}

type EvidenceRecord struct {
	ID                string          `json:"id"`
	Symbol            string          `json:"symbol,omitempty"`
	Kind              string          `json:"kind"`
	TemporalSchema    string          `json:"temporalSchema,omitempty"`
	SourceAt          int64           `json:"sourceAt,omitempty"`
	ObservedAt        int64           `json:"observedAt"`
	IngestedAt        int64           `json:"ingestedAt,omitempty"`
	KnownAt           int64           `json:"knownAt,omitempty"`
	EffectiveFrom     int64           `json:"effectiveFrom,omitempty"`
	EffectiveTo       int64           `json:"effectiveTo,omitempty"`
	ReportPeriod      string          `json:"reportPeriod,omitempty"`
	RevisionID        string          `json:"revisionId,omitempty"`
	SupersedesID      string          `json:"supersedesId,omitempty"`
	AmendmentState    string          `json:"amendmentState,omitempty"`
	Source            string          `json:"source,omitempty"`
	Provenance        string          `json:"provenance,omitempty"`
	FreshnessState    string          `json:"freshnessState,omitempty"`
	RightsState       string          `json:"rightsState,omitempty"`
	RightsEvidenceRef string          `json:"rightsEvidenceRef,omitempty"`
	RetentionClass    string          `json:"retentionClass,omitempty"`
	Payload           json.RawMessage `json:"payload"`
}

type DecisionLineageRecord struct {
	ID             string          `json:"id"`
	Symbol         string          `json:"symbol,omitempty"`
	Horizon        string          `json:"horizon,omitempty"`
	EvidenceID     string          `json:"evidenceId,omitempty"`
	DecisionKind   string          `json:"decisionKind"`
	DecisionValue  string          `json:"decisionValue,omitempty"`
	FormulaVersion string          `json:"formulaVersion,omitempty"`
	CreatedAt      int64           `json:"createdAt"`
	Payload        json.RawMessage `json:"payload"`
}

type OutcomeHistoryRecord struct {
	ID           string          `json:"id"`
	DecisionID   string          `json:"decisionId,omitempty"`
	Symbol       string          `json:"symbol,omitempty"`
	Horizon      string          `json:"horizon,omitempty"`
	ObservedAt   int64           `json:"observedAt"`
	OutcomeLabel string          `json:"outcomeLabel,omitempty"`
	Payload      json.RawMessage `json:"payload"`
}

type DerivedFeatureRecord struct {
	Symbol         string          `json:"symbol"`
	FeatureKey     string          `json:"featureKey"`
	FeatureVersion string          `json:"featureVersion"`
	AsOf           int64           `json:"asOf"`
	SourceHash     string          `json:"sourceHash,omitempty"`
	Payload        json.RawMessage `json:"payload"`
}

type PersistenceIntelligenceBatch struct {
	Evidence  []EvidenceRecord        `json:"evidence,omitempty"`
	Decisions []DecisionLineageRecord `json:"decisions,omitempty"`
	Outcomes  []OutcomeHistoryRecord  `json:"outcomes,omitempty"`
	Features  []DerivedFeatureRecord  `json:"features,omitempty"`
}

func (b PersistenceIntelligenceBatch) Len() int {
	return len(b.Evidence) + len(b.Decisions) + len(b.Outcomes) + len(b.Features)
}

type SymbolRegistryRecord struct {
	Symbol           string `json:"symbol"`
	FirstSeenAt      int64  `json:"firstSeenAt,omitempty"`
	Active           bool   `json:"active"`
	Selected         bool   `json:"selected"`
	ProcessingTier   int    `json:"processingTier"`
	DeskMembership   string `json:"deskMembership,omitempty"`
	ProviderEligible bool   `json:"providerEligible"`
	LastSeenAt       int64  `json:"lastSeenAt"`
	LastSubscribedAt int64  `json:"lastSubscribedAt,omitempty"`
	LastProcessedAt  int64  `json:"lastProcessedAt,omitempty"`
}

type PersistenceStoreStats struct {
	SchemaVersion     int   `json:"schemaVersion"`
	SymbolCount       int   `json:"symbolCount"`
	ActiveSymbolCount int   `json:"activeSymbolCount"`
	CanonicalQuotes   int   `json:"canonicalQuotes"`
	QuoteHistoryRows  int   `json:"quoteHistoryRows"`
	EvidenceRows      int   `json:"evidenceRows"`
	DecisionRows      int   `json:"decisionRows"`
	OutcomeRows       int   `json:"outcomeRows"`
	FeatureRows       int   `json:"featureRows"`
	UserCount         int   `json:"userCount"`
	SessionCount      int   `json:"sessionCount"`
	StorageBytes      int64 `json:"storageBytes"`
}

type PersistencePoolDiagnostics struct {
	MaxOpenConnections int   `json:"maxOpenConnections,omitempty"`
	OpenConnections    int   `json:"openConnections,omitempty"`
	InUse              int   `json:"inUse,omitempty"`
	Idle               int   `json:"idle,omitempty"`
	WaitCount          int64 `json:"waitCount,omitempty"`
	WaitDurationMs     int64 `json:"waitDurationMs,omitempty"`
	MaxIdleClosed      int64 `json:"maxIdleClosed,omitempty"`
	MaxIdleTimeClosed  int64 `json:"maxIdleTimeClosed,omitempty"`
	MaxLifetimeClosed  int64 `json:"maxLifetimeClosed,omitempty"`
}

type PersistenceDatabaseDiagnostics struct {
	Operations              int64  `json:"operations,omitempty"`
	Errors                  int64  `json:"errors,omitempty"`
	SlowOperations          int64  `json:"slowOperations,omitempty"`
	LastOperation           string `json:"lastOperation,omitempty"`
	LastOperationDurationMs int64  `json:"lastOperationDurationMs,omitempty"`
	MaxOperationDurationMs  int64  `json:"maxOperationDurationMs,omitempty"`
}

type persistencePoolObserver interface {
	PoolDiagnostics() PersistencePoolDiagnostics
}

type persistenceDatabaseObserver interface {
	DatabaseDiagnostics() PersistenceDatabaseDiagnostics
}

type persistenceHealthBackend interface {
	HealthCheck(context.Context) error
}

type PersistenceDiagnostics struct {
	Backend                  string                         `json:"backend"`
	Capabilities             []string                       `json:"capabilities,omitempty"`
	Ready                    bool                           `json:"ready"`
	QueueDepth               int                            `json:"queueDepth"`
	OldestJobAgeMs           int64                          `json:"oldestJobAgeMs,omitempty"`
	WriteBatches             int64                          `json:"writeBatches"`
	RowsWritten              int64                          `json:"rowsWritten"`
	WriteBatchesLastMin      int                            `json:"writeBatchesLastMinute"`
	RowsWrittenLastMin       int64                          `json:"rowsWrittenLastMinute"`
	CoalescedRequests        int64                          `json:"coalescedRequests"`
	MaterialWritesSuppressed int64                          `json:"materialWritesSuppressed"`
	RetryBatches             int64                          `json:"retryBatches"`
	DroppedBatches           int64                          `json:"droppedBatches"`
	Errors                   int64                          `json:"errors"`
	LastError                string                         `json:"lastError,omitempty"`
	LastWriteAt              int64                          `json:"lastWriteAt,omitempty"`
	WarmStartedQuotes        int                            `json:"warmStartedQuotes"`
	Store                    PersistenceStoreStats          `json:"store"`
	Pool                     PersistencePoolDiagnostics     `json:"pool,omitempty"`
	Database                 PersistenceDatabaseDiagnostics `json:"database,omitempty"`
	HealthState              string                         `json:"healthState,omitempty"`
	ConsecutiveFailures      int                            `json:"consecutiveFailures,omitempty"`
	ConsecutiveSuccesses     int                            `json:"consecutiveSuccesses,omitempty"`
	LastHealthCheckAt        int64                          `json:"lastHealthCheckAt,omitempty"`
	LastHealthyAt            int64                          `json:"lastHealthyAt,omitempty"`
	RetryScheduled           bool                           `json:"retryScheduled,omitempty"`
	RetryBackoffMs           int64                          `json:"retryBackoffMs,omitempty"`
	PendingIntelligence      int                            `json:"pendingIntelligenceRecords,omitempty"`
	ShedIntelligenceRecords  int64                          `json:"shedIntelligenceRecords,omitempty"`
}

type persistenceWriteEvent struct {
	at   time.Time
	rows int64
}

type PersistenceManager struct {
	backend PersistenceBackend
	wake    chan struct{}
	stop    chan struct{}
	done    chan struct{}

	mu                  sync.Mutex
	pendingSymbols      []SymbolRegistryRecord
	pendingQuotes       map[string]Quote
	pendingIntelligence PersistenceIntelligenceBatch
	lastAcceptedQuotes  map[string]Quote
	oldestPending       time.Time
	diag                PersistenceDiagnostics
	writeEvents         []persistenceWriteEvent
	closed              bool
	initialized         bool
	retryScheduled      bool
	retryBackoff        time.Duration
}

func NewPersistenceManager(configDir string) *PersistenceManager {
	return newPersistenceManagerWithBackend(newPersistenceBackend(configDir))
}

func newPersistenceManagerWithBackend(backend PersistenceBackend) *PersistenceManager {
	if backend == nil {
		backend = newUnavailablePersistenceBackend("persistence backend is nil")
	}
	p := &PersistenceManager{backend: backend, wake: make(chan struct{}, 1), stop: make(chan struct{}), done: make(chan struct{}), pendingQuotes: map[string]Quote{}, lastAcceptedQuotes: map[string]Quote{}}
	p.diag.Backend = backend.Name()
	p.diag.Capabilities = append([]string(nil), backend.Capabilities()...)
	if err := backend.Init(context.Background()); err != nil {
		p.diag.LastError = err.Error()
		p.diag.Errors++
		close(p.done)
		return p
	}
	p.initialized = true
	p.diag.Ready = true
	p.diag.HealthState = "HEALTHY"
	p.diag.ConsecutiveSuccesses = 2
	p.diag.LastHealthCheckAt = time.Now().UnixMilli()
	p.diag.LastHealthyAt = p.diag.LastHealthCheckAt
	p.retryBackoff = 250 * time.Millisecond
	p.refreshStoreStats()
	go p.worker()
	return p
}

func (p *PersistenceManager) LoadIdentityState(ctx context.Context) (IdentityPersistentState, error) {
	if p == nil || p.backend == nil {
		return IdentityPersistentState{}, errors.New("persistence unavailable")
	}
	state, err := p.backend.LoadIdentityState(ctx)
	if err != nil {
		p.recordPersistenceFailure(err)
	}
	return state, err
}

func (p *PersistenceManager) SaveIdentityState(ctx context.Context, state IdentityPersistentState) error {
	if p == nil || p.backend == nil {
		return errors.New("persistence unavailable")
	}
	if err := p.backend.SaveIdentityState(ctx, state); err != nil {
		p.recordPersistenceFailure(err)
		return err
	}
	p.refreshStoreStats()
	return nil
}

func (p *PersistenceManager) LoadUserWorkspaces(ctx context.Context) ([]UserWorkspace, error) {
	if p == nil || p.backend == nil {
		return nil, errors.New("persistence unavailable")
	}
	workspaces, err := p.backend.LoadUserWorkspaces(ctx)
	if err != nil {
		p.recordPersistenceFailure(err)
	}
	return workspaces, err
}

func (p *PersistenceManager) SaveUserWorkspace(ctx context.Context, workspace UserWorkspace) error {
	if p == nil || p.backend == nil {
		return errors.New("persistence unavailable")
	}
	if err := p.backend.SaveUserWorkspace(ctx, workspace); err != nil {
		p.recordPersistenceFailure(err)
		return err
	}
	p.refreshStoreStats()
	return nil
}

func (p *PersistenceManager) worker() {
	p.mu.Lock()
	if !p.initialized && p.diag.Ready {
		p.initialized = true
		if p.diag.HealthState == "" {
			p.diag.HealthState = "HEALTHY"
			p.diag.ConsecutiveSuccesses = 2
			p.diag.LastHealthCheckAt = time.Now().UnixMilli()
			p.diag.LastHealthyAt = p.diag.LastHealthCheckAt
		}
		if p.retryBackoff <= 0 {
			p.retryBackoff = 250 * time.Millisecond
		}
	}
	p.mu.Unlock()
	defer close(p.done)
	for {
		select {
		case <-p.wake:
			p.mu.Lock()
			ready := p.diag.Ready
			p.mu.Unlock()
			if ready {
				p.flushPending()
			}
		case <-p.stop:
			p.flushPending()
			return
		}
	}
}

func (p *PersistenceManager) flushPending() {
	for {
		p.mu.Lock()
		symbols := p.pendingSymbols
		quotes := p.pendingQuotes
		intel := p.pendingIntelligence
		p.pendingSymbols = nil
		p.pendingQuotes = map[string]Quote{}
		p.pendingIntelligence = PersistenceIntelligenceBatch{}
		p.oldestPending = time.Time{}
		p.mu.Unlock()

		if len(symbols) == 0 && len(quotes) == 0 && intel.Len() == 0 {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		rows := 0
		var errs []error
		if len(symbols) > 0 {
			n, err := p.backend.UpsertSymbols(ctx, symbols)
			rows += n
			if err != nil {
				errs = append(errs, err)
			}
		}
		if len(quotes) > 0 {
			n, err := p.backend.SaveQuotes(ctx, quotes)
			rows += n
			if err != nil {
				errs = append(errs, err)
			}
		}
		if intel.Len() > 0 {
			n, err := p.backend.SaveIntelligence(ctx, intel)
			rows += n
			if err != nil {
				errs = append(errs, err)
			}
		}
		cancel()

		p.mu.Lock()
		closed := p.closed
		if len(errs) > 0 {
			p.markPersistenceFailureLocked(errors.Join(errs...))
			if closed {
				// Close is the final bounded flush opportunity. Never spin forever
				// during shutdown; surface any unrecoverable loss explicitly.
				p.diag.DroppedBatches++
				p.mu.Unlock()
				return
			}

			// Restore the failed batch so a transient DB lock/I/O failure cannot
			// silently discard evidence. Newer queued values win per symbol.
			if len(p.pendingSymbols) == 0 && len(symbols) > 0 {
				p.pendingSymbols = append([]SymbolRegistryRecord(nil), symbols...)
			}
			for sym, q := range quotes {
				if _, newer := p.pendingQuotes[sym]; !newer {
					p.pendingQuotes[sym] = q
				}
			}
			p.pendingIntelligence.Evidence = append(intel.Evidence, p.pendingIntelligence.Evidence...)
			p.pendingIntelligence.Decisions = append(intel.Decisions, p.pendingIntelligence.Decisions...)
			p.pendingIntelligence.Outcomes = append(intel.Outcomes, p.pendingIntelligence.Outcomes...)
			p.pendingIntelligence.Features = append(intel.Features, p.pendingIntelligence.Features...)
			var shed int
			p.pendingIntelligence, shed = compactPersistenceIntelligence(p.pendingIntelligence, maxPendingIntelligenceRecords)
			p.diag.ShedIntelligenceRecords += int64(shed)
			if p.oldestPending.IsZero() {
				p.oldestPending = time.Now()
			}
			p.diag.RetryBatches++
			p.scheduleRetryLocked()
			p.mu.Unlock()
			return
		}

		now := time.Now()
		p.diag.WriteBatches++
		p.diag.RowsWritten += int64(rows)
		p.diag.LastWriteAt = now.UnixMilli()
		p.writeEvents = append(p.writeEvents, persistenceWriteEvent{at: now, rows: int64(rows)})
		p.diag.LastError = ""
		p.diag.HealthState = "HEALTHY"
		p.diag.Ready = true
		p.diag.ConsecutiveFailures = 0
		p.diag.ConsecutiveSuccesses = 2
		p.diag.LastHealthyAt = now.UnixMilli()
		p.retryBackoff = 250 * time.Millisecond
		more := len(p.pendingSymbols) > 0 || len(p.pendingQuotes) > 0 || p.pendingIntelligence.Len() > 0
		p.mu.Unlock()
		p.refreshStoreStats()
		if !more {
			return
		}
	}
}

func (p *PersistenceManager) signal() {
	select {
	case p.wake <- struct{}{}:
	default:
		p.mu.Lock()
		p.diag.CoalescedRequests++
		p.mu.Unlock()
	}
}

func (p *PersistenceManager) EnqueueSymbols(records []SymbolRegistryRecord) {
	if p == nil || len(records) == 0 {
		return
	}
	copyRecords := append([]SymbolRegistryRecord(nil), records...)
	p.mu.Lock()
	if (!p.initialized && !p.diag.Ready) || p.closed {
		p.mu.Unlock()
		return
	}
	if len(p.pendingSymbols) > 0 {
		p.diag.CoalescedRequests++
		merged := make(map[string]SymbolRegistryRecord, len(p.pendingSymbols)+len(copyRecords))
		for _, old := range p.pendingSymbols {
			old.Symbol = normalizeSymbol(old.Symbol)
			if old.Symbol == "" {
				continue
			}
			old.Active = false
			old.Selected = false
			merged[old.Symbol] = old
		}
		for _, next := range copyRecords {
			next.Symbol = normalizeSymbol(next.Symbol)
			if next.Symbol == "" {
				continue
			}
			if old, ok := merged[next.Symbol]; ok && old.FirstSeenAt > 0 && (next.FirstSeenAt <= 0 || old.FirstSeenAt < next.FirstSeenAt) {
				next.FirstSeenAt = old.FirstSeenAt
			}
			merged[next.Symbol] = next
		}
		copyRecords = copyRecords[:0]
		for _, r := range merged {
			copyRecords = append(copyRecords, r)
		}
	}
	p.pendingSymbols = copyRecords
	if p.oldestPending.IsZero() {
		p.oldestPending = time.Now()
	}
	p.mu.Unlock()
	p.signal()
}

func quoteMateriallyChanged(prev, next Quote) bool {
	if prev.Symbol == "" || prev.Price <= 0 {
		return true
	}
	if prev.Source != next.Source || prev.FeedType != next.FeedType || prev.DataState != next.DataState {
		return true
	}
	prevStamp := maxInt64(prev.ProviderTimestamp, prev.UpdatedAt)
	nextStamp := maxInt64(next.ProviderTimestamp, next.UpdatedAt)
	if nextStamp > prevStamp && nextStamp-prevStamp >= 90_000 {
		return true
	}
	pct := func(a, b float64) float64 {
		if a == 0 {
			if b != 0 {
				return 1
			}
			return 0
		}
		d := a - b
		if d < 0 {
			d = -d
		}
		return d / a
	}
	if pct(prev.Price, next.Price) >= 0.0002 || pct(prev.Bid, next.Bid) >= 0.0005 || pct(prev.Ask, next.Ask) >= 0.0005 {
		return true
	}
	if prev.PreviousClose != next.PreviousClose || prev.SessionClose != next.SessionClose {
		return true
	}
	return false
}

func (p *PersistenceManager) EnqueueQuotes(quotes map[string]Quote) {
	if p == nil || len(quotes) == 0 {
		return
	}
	p.mu.Lock()
	if (!p.initialized && !p.diag.Ready) || p.closed {
		p.mu.Unlock()
		return
	}
	accepted := 0
	for sym, q := range quotes {
		sym = normalizeSymbol(sym)
		if sym == "" || q.Price <= 0 {
			continue
		}
		q.Symbol = sym
		if prev, ok := p.lastAcceptedQuotes[sym]; ok && !quoteMateriallyChanged(prev, q) {
			p.diag.MaterialWritesSuppressed++
			continue
		}
		if _, queued := p.pendingQuotes[sym]; queued {
			p.diag.CoalescedRequests++
		}
		p.pendingQuotes[sym] = q
		p.lastAcceptedQuotes[sym] = q
		accepted++
	}
	if accepted == 0 {
		p.mu.Unlock()
		return
	}
	if p.oldestPending.IsZero() {
		p.oldestPending = time.Now()
	}
	p.mu.Unlock()
	p.signal()
}

func (p *PersistenceManager) EnqueueIntelligence(batch PersistenceIntelligenceBatch) {
	if p == nil || batch.Len() == 0 {
		return
	}
	p.mu.Lock()
	if (!p.initialized && !p.diag.Ready) || p.closed {
		p.mu.Unlock()
		return
	}
	if p.pendingIntelligence.Len() > 0 {
		p.diag.CoalescedRequests++
	}
	p.pendingIntelligence.Evidence = append(p.pendingIntelligence.Evidence, batch.Evidence...)
	p.pendingIntelligence.Decisions = append(p.pendingIntelligence.Decisions, batch.Decisions...)
	p.pendingIntelligence.Outcomes = append(p.pendingIntelligence.Outcomes, batch.Outcomes...)
	p.pendingIntelligence.Features = append(p.pendingIntelligence.Features, batch.Features...)
	var shed int
	p.pendingIntelligence, shed = compactPersistenceIntelligence(p.pendingIntelligence, maxPendingIntelligenceRecords)
	p.diag.ShedIntelligenceRecords += int64(shed)
	if p.oldestPending.IsZero() {
		p.oldestPending = time.Now()
	}
	p.mu.Unlock()
	p.signal()
}

func (p *PersistenceManager) LoadQuotes() map[string]Quote {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	ready := p.diag.Ready
	p.mu.Unlock()
	if !ready {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	quotes, err := p.backend.LoadQuotes(ctx)
	p.mu.Lock()
	defer p.mu.Unlock()
	if err != nil {
		p.markPersistenceFailureLocked(err)
		return nil
	}
	p.diag.WarmStartedQuotes = len(quotes)
	for sym, q := range quotes {
		p.lastAcceptedQuotes[sym] = q
	}
	return quotes
}

func (p *PersistenceManager) refreshStoreStats() {
	if p == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	stats, err := p.backend.Stats(ctx)
	p.mu.Lock()
	defer p.mu.Unlock()
	if err != nil {
		p.markPersistenceFailureLocked(err)
		return
	}
	p.diag.Store = stats
}

func (p *PersistenceManager) LoadSymbols() []SymbolRegistryRecord {
	if p == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	records, err := p.backend.LoadSymbols(ctx)
	if err != nil {
		p.recordPersistenceFailure(err)
		return nil
	}
	return records
}

func (p *PersistenceManager) ProbeReady(ctx context.Context) bool {
	if p == nil || p.backend == nil {
		return false
	}
	p.mu.Lock()
	initialized := p.initialized
	ready := p.diag.Ready
	retryScheduled := p.retryScheduled
	lastCheck := p.diag.LastHealthCheckAt
	p.mu.Unlock()
	if !initialized {
		return false
	}
	if !ready && retryScheduled {
		return false
	}
	if ready && lastCheck > 0 && time.Since(time.UnixMilli(lastCheck)) < 2*time.Second {
		return true
	}
	health, ok := p.backend.(persistenceHealthBackend)
	if !ok {
		return ready
	}
	if err := health.HealthCheck(ctx); err != nil {
		p.recordPersistenceFailure(err)
		return false
	}
	p.recordPersistenceHealthy()
	p.mu.Lock()
	ready = p.diag.Ready
	if !ready {
		p.scheduleRetryLocked()
	}
	p.mu.Unlock()
	return ready
}

func (p *PersistenceManager) Diagnostics() PersistenceDiagnostics {
	if p == nil {
		return PersistenceDiagnostics{Backend: "disabled"}
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-time.Minute)
	keep := p.writeEvents[:0]
	var recentRows int64
	for _, ev := range p.writeEvents {
		if ev.at.After(cutoff) {
			keep = append(keep, ev)
			recentRows += ev.rows
		}
	}
	p.writeEvents = keep
	d := p.diag
	d.WriteBatchesLastMin = len(keep)
	d.RowsWrittenLastMin = recentRows
	if len(p.pendingSymbols) > 0 {
		d.QueueDepth++
	}
	if len(p.pendingQuotes) > 0 {
		d.QueueDepth++
	}
	if p.pendingIntelligence.Len() > 0 {
		d.QueueDepth++
	}
	d.PendingIntelligence = p.pendingIntelligence.Len()
	d.RetryScheduled = p.retryScheduled
	if !p.oldestPending.IsZero() {
		d.OldestJobAgeMs = time.Since(p.oldestPending).Milliseconds()
	}
	if observer, ok := p.backend.(persistencePoolObserver); ok {
		d.Pool = observer.PoolDiagnostics()
	}
	if observer, ok := p.backend.(persistenceDatabaseObserver); ok {
		d.Database = observer.DatabaseDiagnostics()
	}
	return d
}

const maxPendingIntelligenceRecords = 50000

func compactPersistenceIntelligence(batch PersistenceIntelligenceBatch, limit int) (PersistenceIntelligenceBatch, int) {
	if batch.Len() == 0 {
		return batch, 0
	}
	evidence := map[string]EvidenceRecord{}
	for _, record := range batch.Evidence {
		if record.ID != "" {
			evidence[record.ID] = record
		}
	}
	decisions := map[string]DecisionLineageRecord{}
	for _, record := range batch.Decisions {
		if record.ID != "" {
			decisions[record.ID] = record
		}
	}
	outcomes := map[string]OutcomeHistoryRecord{}
	for _, record := range batch.Outcomes {
		if record.ID != "" {
			outcomes[record.ID] = record
		}
	}
	features := map[string]DerivedFeatureRecord{}
	for _, record := range batch.Features {
		key := normalizeSymbol(record.Symbol) + "|" + record.FeatureKey + "|" + record.FeatureVersion
		if key != "||" {
			features[key] = record
		}
	}
	out := PersistenceIntelligenceBatch{}
	for _, record := range evidence {
		out.Evidence = append(out.Evidence, record)
	}
	for _, record := range decisions {
		out.Decisions = append(out.Decisions, record)
	}
	for _, record := range outcomes {
		out.Outcomes = append(out.Outcomes, record)
	}
	for _, record := range features {
		out.Features = append(out.Features, record)
	}
	sort.Slice(out.Evidence, func(i, j int) bool { return out.Evidence[i].ID < out.Evidence[j].ID })
	sort.Slice(out.Decisions, func(i, j int) bool { return out.Decisions[i].ID < out.Decisions[j].ID })
	sort.Slice(out.Outcomes, func(i, j int) bool { return out.Outcomes[i].ID < out.Outcomes[j].ID })
	sort.Slice(out.Features, func(i, j int) bool {
		a := normalizeSymbol(out.Features[i].Symbol) + "|" + out.Features[i].FeatureKey + "|" + out.Features[i].FeatureVersion
		b := normalizeSymbol(out.Features[j].Symbol) + "|" + out.Features[j].FeatureKey + "|" + out.Features[j].FeatureVersion
		return a < b
	})
	if limit <= 0 || out.Len() <= limit {
		return out, 0
	}
	before := out.Len()
	// Derived features are reproducible and are shed before immutable evidence/decision/outcome lineage.
	for out.Len() > limit && len(out.Features) > 0 {
		out.Features = out.Features[1:]
	}
	// A very high hard ceiling still bounds memory during a prolonged outage after shedding reproducible features.
	for out.Len() > limit && len(out.Evidence) > 0 {
		out.Evidence = out.Evidence[1:]
	}
	for out.Len() > limit && len(out.Decisions) > 0 {
		out.Decisions = out.Decisions[1:]
	}
	for out.Len() > limit && len(out.Outcomes) > 0 {
		out.Outcomes = out.Outcomes[1:]
	}
	return out, before - out.Len()
}

func (p *PersistenceManager) markPersistenceFailureLocked(err error) {
	if err == nil {
		return
	}
	p.diag.Errors++
	p.diag.LastError = err.Error()
	p.diag.HealthState = "DEGRADED"
	p.diag.Ready = false
	p.diag.ConsecutiveFailures++
	p.diag.ConsecutiveSuccesses = 0
	p.diag.LastHealthCheckAt = time.Now().UnixMilli()
	p.scheduleRetryLocked()
}

func (p *PersistenceManager) recordPersistenceFailure(err error) {
	if p == nil || err == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.markPersistenceFailureLocked(err)
}

func (p *PersistenceManager) recordPersistenceHealthy() {
	if p == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	now := time.Now().UnixMilli()
	p.diag.LastHealthCheckAt = now
	p.diag.LastHealthyAt = now
	p.diag.ConsecutiveSuccesses++
	if p.diag.ConsecutiveSuccesses >= 2 {
		p.diag.Ready = true
		p.diag.HealthState = "HEALTHY"
		p.diag.ConsecutiveFailures = 0
		p.diag.LastError = ""
		p.retryBackoff = 250 * time.Millisecond
	}
}

func (p *PersistenceManager) scheduleRetryLocked() {
	if p == nil || p.closed || !p.initialized || p.retryScheduled {
		return
	}
	if p.retryBackoff <= 0 {
		p.retryBackoff = 250 * time.Millisecond
	}
	delay := p.retryBackoff
	if delay > 5*time.Second {
		delay = 5 * time.Second
	}
	p.retryScheduled = true
	p.diag.RetryScheduled = true
	p.diag.RetryBackoffMs = delay.Milliseconds()
	next := delay * 2
	if next > 5*time.Second {
		next = 5 * time.Second
	}
	p.retryBackoff = next
	time.AfterFunc(delay, p.retryProbe)
}

func (p *PersistenceManager) retryProbe() {
	if p == nil {
		return
	}
	p.mu.Lock()
	if p.closed || !p.initialized {
		p.retryScheduled = false
		p.mu.Unlock()
		return
	}
	p.retryScheduled = false
	p.diag.RetryScheduled = false
	p.mu.Unlock()
	health, ok := p.backend.(persistenceHealthBackend)
	if !ok {
		p.recordPersistenceHealthy()
		p.recordPersistenceHealthy()
		p.signal()
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	err := health.HealthCheck(ctx)
	cancel()
	if err != nil {
		p.recordPersistenceFailure(err)
		return
	}
	p.recordPersistenceHealthy()
	p.mu.Lock()
	ready := p.diag.Ready
	if !ready {
		p.scheduleRetryLocked()
	}
	p.mu.Unlock()
	if ready {
		p.signal()
	}
}

func (p *PersistenceManager) Close() error {
	if p == nil {
		return nil
	}
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	started := p.initialized || p.diag.Ready
	p.mu.Unlock()
	if started {
		close(p.stop)
		<-p.done
	}
	return p.backend.Close()
}

func symbolRegistryRecords(st AppState, now time.Time) []SymbolRegistryRecord {
	selected := normalizeSymbol(st.UI.SelectedTicker)
	activeDesk := map[string]bool{}
	desksBySymbol := map[string][]string{}
	for _, id := range deskIDs() {
		if wl, ok := watchlistValueByID(st.Watchlists, id); ok {
			for _, sym := range userTradingSymbols(wl.Symbols) {
				activeDesk[sym] = true
				desksBySymbol[sym] = append(desksBySymbol[sym], id)
			}
		}
	}
	discovery := map[string]bool{}
	for _, sym := range discoverySymbolsFromState(st) {
		discovery[sym] = true
	}
	out := make([]SymbolRegistryRecord, 0, len(masterSymbolsFromState(st)))
	for _, sym := range masterSymbolsFromState(st) {
		tier := 3
		if sym == "SPY" || sym == "QQQ" || sym == selected {
			tier = 0
		} else if activeDesk[sym] || sym == "GLD" || sym == "SLV" || sym == "USO" {
			tier = 1
		} else if discovery[sym] {
			tier = 2
		}
		memberships, _ := json.Marshal(desksBySymbol[sym])
		// Membership persistence must not fabricate provider/runtime events.
		// LastSubscribedAt and LastProcessedAt are advanced only by owners that
		// actually observe those events; zero preserves any prior DB value.
		out = append(out, SymbolRegistryRecord{Symbol: sym, FirstSeenAt: now.UnixMilli(), Active: true, Selected: sym == selected, ProcessingTier: tier, DeskMembership: string(memberships), ProviderEligible: validUserTicker(sym), LastSeenAt: now.UnixMilli()})
	}
	return out
}
