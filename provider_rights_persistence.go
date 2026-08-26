package main

import (
	"context"
	"strings"
	"time"
)

// hostedRightsPersistenceBackend is a narrow policy decorator around the
// existing canonical persistence backend. It does not own storage, migrations,
// tenancy, retries or recovery. It prevents provider data from crossing the
// hosted persistence boundary when current legal/data-rights evidence does not
// approve production serving and retention.
type hostedRightsPersistenceBackend struct {
	inner PersistenceBackend
}

func wrapHostedRightsPersistenceBackend(inner PersistenceBackend) PersistenceBackend {
	if inner == nil || !isHostedRuntime() {
		return inner
	}
	return &hostedRightsPersistenceBackend{inner: inner}
}

func (b *hostedRightsPersistenceBackend) Name() string {
	if b == nil || b.inner == nil {
		return "unavailable"
	}
	return b.inner.Name()
}

func (b *hostedRightsPersistenceBackend) Capabilities() []string {
	if b == nil || b.inner == nil {
		return nil
	}
	out := append([]string(nil), b.inner.Capabilities()...)
	return append(out, "hosted-provider-rights-filter")
}

func (b *hostedRightsPersistenceBackend) Init(ctx context.Context) error {
	return b.inner.Init(ctx)
}

func (b *hostedRightsPersistenceBackend) UpsertSymbols(ctx context.Context, records []SymbolRegistryRecord) (int, error) {
	return b.inner.UpsertSymbols(ctx, records)
}

func (b *hostedRightsPersistenceBackend) LoadSymbols(ctx context.Context) ([]SymbolRegistryRecord, error) {
	return b.inner.LoadSymbols(ctx)
}

func hostedRightsFilterQuotes(quotes map[string]Quote, now time.Time) map[string]Quote {
	if !isHostedRuntime() {
		return quotes
	}
	out := make(map[string]Quote, len(quotes))
	for symbol, q := range quotes {
		if providerQuoteHostedRightsAllowed(q, providerHostedUseProductionServing, now) {
			out[symbol] = q
		}
	}
	return out
}

func (b *hostedRightsPersistenceBackend) SaveQuotes(ctx context.Context, quotes map[string]Quote) (int, error) {
	filtered := hostedRightsFilterQuotes(quotes, time.Now())
	if len(filtered) == 0 {
		return 0, nil
	}
	return b.inner.SaveQuotes(ctx, filtered)
}

func (b *hostedRightsPersistenceBackend) LoadQuotes(ctx context.Context) (map[string]Quote, error) {
	quotes, err := b.inner.LoadQuotes(ctx)
	if err != nil {
		return nil, err
	}
	return hostedRightsFilterQuotes(quotes, time.Now()), nil
}

func hostedRightsExternalEvidenceAllowed(record EvidenceRecord, now time.Time) bool {
	if !isHostedRuntime() {
		return true
	}
	source := strings.TrimSpace(record.Source)
	if source == "" {
		return true
	}
	provider := sourceProvider(source)
	if provider == "—" {
		// Some canonical/internal evidence uses semantic source labels rather
		// than external-provider identities. HOST-022 owns full point-in-time
		// provenance conservation; recognized external providers are enforced
		// here without misclassifying internal evidence as provider data.
		return true
	}
	return hostedProviderRightsAllowed(provider, providerHostedUseProductionServing, now)
}

func (b *hostedRightsPersistenceBackend) SaveIntelligence(ctx context.Context, batch PersistenceIntelligenceBatch) (int, error) {
	if !isHostedRuntime() {
		return b.inner.SaveIntelligence(ctx, batch)
	}
	now := time.Now()
	filtered := batch
	filtered.Evidence = filtered.Evidence[:0]
	for _, record := range batch.Evidence {
		if hostedRightsExternalEvidenceAllowed(record, now) {
			filtered.Evidence = append(filtered.Evidence, record)
		}
	}
	if filtered.Len() == 0 {
		return 0, nil
	}
	return b.inner.SaveIntelligence(ctx, filtered)
}

func (b *hostedRightsPersistenceBackend) LoadIdentityState(ctx context.Context) (IdentityPersistentState, error) {
	return b.inner.LoadIdentityState(ctx)
}

func (b *hostedRightsPersistenceBackend) SaveIdentityState(ctx context.Context, state IdentityPersistentState) error {
	return b.inner.SaveIdentityState(ctx, state)
}

func (b *hostedRightsPersistenceBackend) LoadUserWorkspaces(ctx context.Context) ([]UserWorkspace, error) {
	return b.inner.LoadUserWorkspaces(ctx)
}

func (b *hostedRightsPersistenceBackend) SaveUserWorkspace(ctx context.Context, workspace UserWorkspace) error {
	return b.inner.SaveUserWorkspace(ctx, workspace)
}

func (b *hostedRightsPersistenceBackend) Stats(ctx context.Context) (PersistenceStoreStats, error) {
	return b.inner.Stats(ctx)
}

func (b *hostedRightsPersistenceBackend) Close() error {
	if b == nil || b.inner == nil {
		return nil
	}
	return b.inner.Close()
}

func (b *hostedRightsPersistenceBackend) HealthCheck(ctx context.Context) error {
	if health, ok := b.inner.(persistenceHealthBackend); ok {
		return health.HealthCheck(ctx)
	}
	return nil
}

func (b *hostedRightsPersistenceBackend) PoolDiagnostics() PersistencePoolDiagnostics {
	if observer, ok := b.inner.(persistencePoolObserver); ok {
		return observer.PoolDiagnostics()
	}
	return PersistencePoolDiagnostics{}
}

func (b *hostedRightsPersistenceBackend) DatabaseDiagnostics() PersistenceDatabaseDiagnostics {
	if observer, ok := b.inner.(persistenceDatabaseObserver); ok {
		return observer.DatabaseDiagnostics()
	}
	return PersistenceDatabaseDiagnostics{}
}
