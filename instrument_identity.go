package main

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"
)

const instrumentIdentitySchemaVersion = 5

// InstrumentIdentityRecord is slow-changing canonical company/instrument
// identity. It is deliberately separate from quotes/current market state and
// from SymbolRegistryRecord, whose persistence semantics are a complete active
// trading-registry snapshot.
type InstrumentIdentityRecord struct {
	Symbol          string `json:"symbol"`
	Name            string `json:"name,omitempty"`
	Exchange        string `json:"exchange,omitempty"`
	AssetClass      string `json:"assetClass,omitempty"`
	ProviderAssetID string `json:"providerAssetId,omitempty"`
	Source          string `json:"source,omitempty"`
	ObservedAt      int64  `json:"observedAt"`
}

// instrumentIdentityPersistence is a narrow capability extension of the
// canonical PersistenceBackend. Keeping this capability separate avoids
// widening unrelated test doubles while all production persistence owners use
// the same Save/Load contract.
type instrumentIdentityPersistence interface {
	SaveInstrumentIdentities(context.Context, []InstrumentIdentityRecord) (int, error)
	LoadInstrumentIdentities(context.Context) ([]InstrumentIdentityRecord, error)
}

func normalizeInstrumentIdentity(record InstrumentIdentityRecord) (InstrumentIdentityRecord, bool) {
	record.Symbol = normalizeSymbol(record.Symbol)
	if record.Symbol == "" {
		return InstrumentIdentityRecord{}, false
	}
	record.Name = strings.TrimSpace(record.Name)
	record.Exchange = strings.ToUpper(strings.TrimSpace(record.Exchange))
	record.AssetClass = strings.ToLower(strings.TrimSpace(record.AssetClass))
	record.ProviderAssetID = strings.TrimSpace(record.ProviderAssetID)
	record.Source = strings.ToLower(strings.TrimSpace(record.Source))
	if record.ObservedAt < 0 {
		record.ObservedAt = 0
	}
	return record, true
}

func canonicalInstrumentIdentities(records []InstrumentIdentityRecord) []InstrumentIdentityRecord {
	bySymbol := make(map[string]InstrumentIdentityRecord, len(records))
	for _, record := range records {
		normalized, ok := normalizeInstrumentIdentity(record)
		if !ok {
			continue
		}
		if previous, exists := bySymbol[normalized.Symbol]; exists && previous.ObservedAt > normalized.ObservedAt {
			continue
		}
		bySymbol[normalized.Symbol] = normalized
	}
	out := make([]InstrumentIdentityRecord, 0, len(bySymbol))
	for _, record := range bySymbol {
		out = append(out, record)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Symbol < out[j].Symbol })
	return out
}

func (p *PersistenceManager) SaveInstrumentIdentities(ctx context.Context, records []InstrumentIdentityRecord) (int, error) {
	if p == nil || p.backend == nil {
		return 0, errors.New("persistence unavailable")
	}
	backend, ok := p.backend.(instrumentIdentityPersistence)
	if !ok {
		return 0, errors.New("instrument identity persistence unavailable")
	}
	records = canonicalInstrumentIdentities(records)
	if len(records) == 0 {
		return 0, nil
	}
	written, err := backend.SaveInstrumentIdentities(ctx, records)
	if err != nil {
		p.recordPersistenceFailure(err)
		return written, err
	}
	p.refreshStoreStats()
	return written, nil
}

func (p *PersistenceManager) LoadInstrumentIdentities(ctx context.Context) ([]InstrumentIdentityRecord, error) {
	if p == nil || p.backend == nil {
		return nil, errors.New("persistence unavailable")
	}
	backend, ok := p.backend.(instrumentIdentityPersistence)
	if !ok {
		return nil, errors.New("instrument identity persistence unavailable")
	}
	records, err := backend.LoadInstrumentIdentities(ctx)
	if err != nil {
		p.recordPersistenceFailure(err)
		return nil, err
	}
	return canonicalInstrumentIdentities(records), nil
}

// warmInstrumentIdentities restores slow-changing identity before any consumer
// needs provider acquisition. A failed restore is retryable and never changes
// universe/provider health semantics.
func (e *Engine) warmInstrumentIdentities() int {
	if e == nil {
		return 0
	}
	e.mu.Lock()
	if e.instrumentIdentitiesLoaded {
		count := len(e.instrumentIdentities)
		e.mu.Unlock()
		return count
	}
	e.instrumentIdentitiesLoaded = true
	if e.instrumentIdentities == nil {
		e.instrumentIdentities = map[string]InstrumentIdentityRecord{}
	}
	e.mu.Unlock()

	if e.app == nil || e.app.persistence == nil {
		return 0
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	records, err := e.app.persistence.LoadInstrumentIdentities(ctx)
	cancel()
	if err != nil {
		e.mu.Lock()
		e.instrumentIdentitiesLoaded = false
		e.mu.Unlock()
		return 0
	}
	e.mu.Lock()
	for _, record := range records {
		e.instrumentIdentities[record.Symbol] = record
	}
	count := len(e.instrumentIdentities)
	e.mu.Unlock()
	return count
}

// acceptInstrumentIdentities updates canonical in-memory identity and persists
// the same normalized records through the existing persistence owner. It never
// changes SymbolRegistryRecord active/selected state and persistence failure is
// scoped to identity rather than turning a successful universe route into a
// provider failure.
func (e *Engine) acceptInstrumentIdentities(records []InstrumentIdentityRecord) int {
	if e == nil {
		return 0
	}
	records = canonicalInstrumentIdentities(records)
	if len(records) == 0 {
		return 0
	}
	e.warmInstrumentIdentities()
	e.mu.Lock()
	if e.instrumentIdentities == nil {
		e.instrumentIdentities = map[string]InstrumentIdentityRecord{}
	}
	for _, record := range records {
		if previous, exists := e.instrumentIdentities[record.Symbol]; exists && previous.ObservedAt > record.ObservedAt {
			continue
		}
		e.instrumentIdentities[record.Symbol] = record
	}
	e.mu.Unlock()

	if e.app == nil || e.app.persistence == nil {
		return 0
	}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	written, err := e.app.persistence.SaveInstrumentIdentities(ctx, records)
	cancel()
	e.mu.Lock()
	if e.health != nil {
		if err != nil {
			e.health["instrument-identity"] = "persist degraded · canonical in-memory identity retained"
		} else {
			e.health["instrument-identity"] = "ready · persisted canonical identity"
		}
	}
	e.mu.Unlock()
	return written
}
