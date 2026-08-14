package main

import (
	"context"
	"encoding/json"
	"errors"
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
	ID             string          `json:"id"`
	Symbol         string          `json:"symbol,omitempty"`
	Kind           string          `json:"kind"`
	ObservedAt     int64           `json:"observedAt"`
	Source         string          `json:"source,omitempty"`
	Provenance     string          `json:"provenance,omitempty"`
	FreshnessState string          `json:"freshnessState,omitempty"`
	Payload        json.RawMessage `json:"payload"`
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

type PersistenceDiagnostics struct {
	Backend                  string                `json:"backend"`
	Capabilities             []string              `json:"capabilities,omitempty"`
	Ready                    bool                  `json:"ready"`
	QueueDepth               int                   `json:"queueDepth"`
	OldestJobAgeMs           int64                 `json:"oldestJobAgeMs,omitempty"`
	WriteBatches             int64                 `json:"writeBatches"`
	RowsWritten              int64                 `json:"rowsWritten"`
	WriteBatchesLastMin      int                   `json:"writeBatchesLastMinute"`
	RowsWrittenLastMin       int64                 `json:"rowsWrittenLastMinute"`
	CoalescedRequests        int64                 `json:"coalescedRequests"`
	MaterialWritesSuppressed int64                 `json:"materialWritesSuppressed"`
	RetryBatches             int64                 `json:"retryBatches"`
	DroppedBatches           int64                 `json:"droppedBatches"`
	Errors                   int64                 `json:"errors"`
	LastError                string                `json:"lastError,omitempty"`
	LastWriteAt              int64                 `json:"lastWriteAt,omitempty"`
	WarmStartedQuotes        int                   `json:"warmStartedQuotes"`
	Store                    PersistenceStoreStats `json:"store"`
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
}

func NewPersistenceManager(configDir string) *PersistenceManager {
	backend := newPersistenceBackend(configDir)
	p := &PersistenceManager{backend: backend, wake: make(chan struct{}, 1), stop: make(chan struct{}), done: make(chan struct{}), pendingQuotes: map[string]Quote{}, lastAcceptedQuotes: map[string]Quote{}}
	p.diag.Backend = backend.Name()
	p.diag.Capabilities = append([]string(nil), backend.Capabilities()...)
	if err := backend.Init(context.Background()); err != nil {
		p.diag.LastError = err.Error()
		p.diag.Errors++
		close(p.done)
		return p
	}
	p.diag.Ready = true
	p.refreshStoreStats()
	go p.worker()
	return p
}

func (p *PersistenceManager) LoadIdentityState(ctx context.Context) (IdentityPersistentState, error) {
	if p == nil || p.backend == nil {
		return IdentityPersistentState{}, errors.New("persistence unavailable")
	}
	return p.backend.LoadIdentityState(ctx)
}

func (p *PersistenceManager) SaveIdentityState(ctx context.Context, state IdentityPersistentState) error {
	if p == nil || p.backend == nil {
		return errors.New("persistence unavailable")
	}
	if err := p.backend.SaveIdentityState(ctx, state); err != nil {
		return err
	}
	p.refreshStoreStats()
	return nil
}

func (p *PersistenceManager) LoadUserWorkspaces(ctx context.Context) ([]UserWorkspace, error) {
	if p == nil || p.backend == nil {
		return nil, errors.New("persistence unavailable")
	}
	return p.backend.LoadUserWorkspaces(ctx)
}

func (p *PersistenceManager) SaveUserWorkspace(ctx context.Context, workspace UserWorkspace) error {
	if p == nil || p.backend == nil {
		return errors.New("persistence unavailable")
	}
	if err := p.backend.SaveUserWorkspace(ctx, workspace); err != nil {
		return err
	}
	p.refreshStoreStats()
	return nil
}

func (p *PersistenceManager) worker() {
	defer close(p.done)
	for {
		select {
		case <-p.wake:
			p.flushPending()
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
			p.diag.Errors++
			p.diag.LastError = errors.Join(errs...).Error()
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
			if p.oldestPending.IsZero() {
				p.oldestPending = time.Now()
			}
			p.diag.RetryBatches++
			p.mu.Unlock()
			time.AfterFunc(250*time.Millisecond, p.signal)
			return
		}

		now := time.Now()
		p.diag.WriteBatches++
		p.diag.RowsWritten += int64(rows)
		p.diag.LastWriteAt = now.UnixMilli()
		p.writeEvents = append(p.writeEvents, persistenceWriteEvent{at: now, rows: int64(rows)})
		p.diag.LastError = ""
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
	if !p.diag.Ready || p.closed {
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
	if !p.diag.Ready || p.closed {
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
	if !p.diag.Ready || p.closed {
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
		p.diag.Errors++
		p.diag.LastError = err.Error()
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
		p.diag.Errors++
		p.diag.LastError = err.Error()
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
		p.mu.Lock()
		p.diag.Errors++
		p.diag.LastError = err.Error()
		p.mu.Unlock()
		return nil
	}
	return records
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
	if !p.oldestPending.IsZero() {
		d.OldestJobAgeMs = time.Since(p.oldestPending).Milliseconds()
	}
	return d
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
	ready := p.diag.Ready
	p.mu.Unlock()
	if ready {
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
