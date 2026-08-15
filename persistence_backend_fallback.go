//go:build !cgo && !windows

package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type filePersistenceBackend struct {
	mu   sync.Mutex
	path string
	data filePersistenceData
}

type filePersistenceData struct {
	Symbols    map[string]SymbolRegistryRecord  `json:"symbols"`
	Quotes     map[string]Quote                 `json:"quotes"`
	Evidence   map[string]EvidenceRecord        `json:"evidence,omitempty"`
	Decisions  map[string]DecisionLineageRecord `json:"decisions,omitempty"`
	Outcomes   map[string]OutcomeHistoryRecord  `json:"outcomes,omitempty"`
	Features   map[string]DerivedFeatureRecord  `json:"features,omitempty"`
	Identity   IdentityPersistentState          `json:"identity,omitempty"`
	Workspaces map[string]UserWorkspace         `json:"workspaces,omitempty"`
}

func newLocalPersistenceBackend(configDir string) PersistenceBackend {
	return &filePersistenceBackend{path: filepath.Join(configDir, "persistent-intelligence-fallback.json")}
}
func (b *filePersistenceBackend) Name() string { return "file-fallback" }
func (b *filePersistenceBackend) Capabilities() []string {
	return []string{"global-symbol-registry", "canonical-quotes", "evidence-records", "decision-lineage", "outcome-history", "derived-feature-store", "user-workspaces"}
}
func (b *filePersistenceBackend) Init(context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data.Symbols = map[string]SymbolRegistryRecord{}
	b.data.Quotes = map[string]Quote{}
	b.data.Evidence = map[string]EvidenceRecord{}
	b.data.Decisions = map[string]DecisionLineageRecord{}
	b.data.Outcomes = map[string]OutcomeHistoryRecord{}
	b.data.Features = map[string]DerivedFeatureRecord{}
	b.data.Workspaces = map[string]UserWorkspace{}
	if raw, err := os.ReadFile(b.path); err == nil {
		_ = json.Unmarshal(raw, &b.data)
	}
	if b.data.Symbols == nil {
		b.data.Symbols = map[string]SymbolRegistryRecord{}
	}
	if b.data.Quotes == nil {
		b.data.Quotes = map[string]Quote{}
	}
	if b.data.Evidence == nil {
		b.data.Evidence = map[string]EvidenceRecord{}
	}
	if b.data.Decisions == nil {
		b.data.Decisions = map[string]DecisionLineageRecord{}
	}
	if b.data.Outcomes == nil {
		b.data.Outcomes = map[string]OutcomeHistoryRecord{}
	}
	if b.data.Features == nil {
		b.data.Features = map[string]DerivedFeatureRecord{}
	}
	if b.data.Workspaces == nil {
		b.data.Workspaces = map[string]UserWorkspace{}
	}
	return nil
}
func (b *filePersistenceBackend) persistLocked() error {
	raw, err := json.Marshal(b.data)
	if err != nil {
		return err
	}
	return atomicWrite(b.path, raw, 0600)
}
func (b *filePersistenceBackend) UpsertSymbols(_ context.Context, records []SymbolRegistryRecord) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for sym, r := range b.data.Symbols {
		r.Active = false
		r.Selected = false
		b.data.Symbols[sym] = r
	}
	for _, r := range records {
		r.Symbol = normalizeSymbol(r.Symbol)
		if r.Symbol == "" {
			continue
		}
		if old, ok := b.data.Symbols[r.Symbol]; ok {
			if old.FirstSeenAt > 0 {
				r.FirstSeenAt = old.FirstSeenAt
			}
		}
		if r.FirstSeenAt <= 0 {
			r.FirstSeenAt = r.LastSeenAt
		}
		b.data.Symbols[r.Symbol] = r
	}
	return len(records), b.persistLocked()
}
func (b *filePersistenceBackend) LoadSymbols(context.Context) ([]SymbolRegistryRecord, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]SymbolRegistryRecord, 0, len(b.data.Symbols))
	for _, r := range b.data.Symbols {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Symbol < out[j].Symbol })
	return out, nil
}

func (b *filePersistenceBackend) SaveQuotes(_ context.Context, quotes map[string]Quote) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for sym, q := range quotes {
		b.data.Quotes[sym] = q
	}
	return len(quotes), b.persistLocked()
}
func (b *filePersistenceBackend) SaveIntelligence(_ context.Context, batch PersistenceIntelligenceBatch) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	written := 0
	for _, r := range batch.Evidence {
		if r.ID != "" {
			if _, exists := b.data.Evidence[r.ID]; !exists {
				b.data.Evidence[r.ID] = r
				written++
			}
		}
	}
	for _, r := range batch.Decisions {
		if r.ID != "" {
			if _, exists := b.data.Decisions[r.ID]; !exists {
				b.data.Decisions[r.ID] = r
				written++
			}
		}
	}
	for _, r := range batch.Outcomes {
		if r.ID != "" {
			b.data.Outcomes[r.ID] = r
			written++
		}
	}
	for _, r := range batch.Features {
		if r.Symbol == "" || r.FeatureKey == "" || r.FeatureVersion == "" {
			continue
		}
		key := normalizeSymbol(r.Symbol) + "|" + r.FeatureKey + "|" + r.FeatureVersion
		old, ok := b.data.Features[key]
		if !ok || old.SourceHash != r.SourceHash {
			b.data.Features[key] = r
			written++
		}
	}
	return written, b.persistLocked()
}

func (b *filePersistenceBackend) LoadQuotes(context.Context) (map[string]Quote, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := clone(b.data.Quotes)
	for sym, q := range out {
		q.Symbol = normalizeSymbol(sym)
		q.DataState = "persisted"
		q.FeedType = "persisted"
		out[sym] = q
	}
	return out, nil
}
func (b *filePersistenceBackend) LoadIdentityState(context.Context) (IdentityPersistentState, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	raw, err := json.Marshal(b.data.Identity)
	if err != nil {
		return IdentityPersistentState{}, err
	}
	var out IdentityPersistentState
	if err := json.Unmarshal(raw, &out); err != nil {
		return IdentityPersistentState{}, err
	}
	return out, nil
}
func (b *filePersistenceBackend) SaveIdentityState(_ context.Context, state IdentityPersistentState) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data.Identity = state
	return b.persistLocked()
}

func (b *filePersistenceBackend) LoadUserWorkspaces(context.Context) ([]UserWorkspace, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]UserWorkspace, 0, len(b.data.Workspaces))
	for _, workspace := range b.data.Workspaces {
		out = append(out, clone(workspace))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].UserID < out[j].UserID })
	return out, nil
}
func (b *filePersistenceBackend) SaveUserWorkspace(_ context.Context, workspace UserWorkspace) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if workspace.UserID == "" {
		return errors.New("workspace user id is required")
	}
	b.data.Workspaces[workspace.UserID] = clone(workspace)
	return b.persistLocked()
}

func (b *filePersistenceBackend) Stats(context.Context) (PersistenceStoreStats, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	active := 0
	for _, r := range b.data.Symbols {
		if r.Active {
			active++
		}
	}
	var size int64
	if info, err := os.Stat(b.path); err == nil {
		size = info.Size()
	}
	return PersistenceStoreStats{
		SchemaVersion: 4, SymbolCount: len(b.data.Symbols), ActiveSymbolCount: active,
		CanonicalQuotes: len(b.data.Quotes), EvidenceRows: len(b.data.Evidence), DecisionRows: len(b.data.Decisions),
		OutcomeRows: len(b.data.Outcomes), FeatureRows: len(b.data.Features), UserCount: len(b.data.Identity.Users), SessionCount: len(b.data.Identity.Sessions), StorageBytes: size,
	}, nil
}

func (b *filePersistenceBackend) Close() error { return nil }
