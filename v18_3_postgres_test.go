//go:build postgres

package main

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"testing"
	"time"
)

func v183PostgresURL(t *testing.T) string {
	t.Helper()
	url := os.Getenv("DEPULSE_TEST_POSTGRES_URL")
	if url == "" {
		t.Skip("DEPULSE_TEST_POSTGRES_URL is required for PostgreSQL integration tests")
	}
	return url
}

func newV183PostgresBackend(t *testing.T) *postgresPersistenceBackend {
	t.Helper()
	backend := newPostgresPersistenceBackend(postgresPersistenceConfig{
		DatabaseURL:     v183PostgresURL(t),
		MaxOpenConns:    4,
		MaxIdleConns:    2,
		ConnMaxLifetime: 10 * time.Minute,
		ConnMaxIdleTime: 2 * time.Minute,
	}).(*postgresPersistenceBackend)
	if err := backend.Init(context.Background()); err != nil {
		t.Fatalf("postgres init: %v", err)
	}
	resetV183Postgres(t, backend)
	return backend
}

func resetV183Postgres(t *testing.T, backend *postgresPersistenceBackend) {
	t.Helper()
	_, err := backend.db.ExecContext(context.Background(), `TRUNCATE TABLE user_workspaces, identity_state, derived_features, outcome_history, decision_lineage, evidence_records, quote_history, canonical_quotes, symbol_registry`)
	if err != nil {
		t.Fatalf("reset postgres test database: %v", err)
	}
}

func TestV183PostgresRepositoryParityAndWarmStart(t *testing.T) {
	backend := newV183PostgresBackend(t)
	ctx := context.Background()
	now := time.Now().UnixMilli()

	if _, err := backend.UpsertSymbols(ctx, []SymbolRegistryRecord{{Symbol: "NVDA", FirstSeenAt: now - 1000, LastSeenAt: now, Active: true, Selected: true, ProcessingTier: 0, ProviderEligible: true}}); err != nil {
		t.Fatal(err)
	}
	quote := Quote{Symbol: "NVDA", Price: 201.5, Bid: 201.4, Ask: 201.6, Volume: 1000, Source: "alpaca-iex", FeedType: "websocket", DataState: "live", ProviderTimestamp: now, UpdatedAt: now}
	if _, err := backend.SaveQuotes(ctx, map[string]Quote{"NVDA": quote}); err != nil {
		t.Fatal(err)
	}
	batch := PersistenceIntelligenceBatch{
		Evidence:  []EvidenceRecord{{ID: "ev-v183", Symbol: "NVDA", Kind: "research", ObservedAt: now, Source: "canonical", FreshnessState: "CURRENT", Payload: json.RawMessage(`{"value":1}`)}},
		Decisions: []DecisionLineageRecord{{ID: "dec-v183", Symbol: "NVDA", Horizon: "day", EvidenceID: "ev-v183", DecisionKind: "readiness", CreatedAt: now, Payload: json.RawMessage(`{"state":"WATCH"}`)}},
		Outcomes:  []OutcomeHistoryRecord{{ID: "out-v183", DecisionID: "dec-v183", Symbol: "NVDA", Horizon: "day", ObservedAt: now, OutcomeLabel: "PENDING", Payload: json.RawMessage(`{}`)}},
		Features:  []DerivedFeatureRecord{{Symbol: "NVDA", FeatureKey: "behavior", FeatureVersion: "v1", AsOf: now, SourceHash: "abc", Payload: json.RawMessage(`{"score":1}`)}},
	}
	if _, err := backend.SaveIntelligence(ctx, batch); err != nil {
		t.Fatal(err)
	}
	identity := IdentityPersistentState{Version: 1, UpdatedAt: now}
	if err := backend.SaveIdentityState(ctx, identity); err != nil {
		t.Fatal(err)
	}
	workspace := defaultUserWorkspace("user-v183")
	workspace.UpdatedAt = now
	if err := backend.SaveUserWorkspace(ctx, workspace); err != nil {
		t.Fatal(err)
	}

	stats, err := backend.Stats(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if stats.SchemaVersion != 4 || stats.SymbolCount != 1 || stats.CanonicalQuotes != 1 || stats.EvidenceRows != 1 || stats.DecisionRows != 1 || stats.OutcomeRows != 1 || stats.FeatureRows != 1 {
		t.Fatalf("postgres stats parity failure: %+v", stats)
	}
	pool := backend.PoolDiagnostics()
	if pool.MaxOpenConnections != 4 || pool.OpenConnections < 1 {
		t.Fatalf("postgres pool diagnostics missing/bad: %+v", pool)
	}
	if err := backend.Close(); err != nil {
		t.Fatal(err)
	}

	reopened := newPostgresPersistenceBackend(postgresPersistenceConfig{DatabaseURL: v183PostgresURL(t), MaxOpenConns: 4, MaxIdleConns: 2, ConnMaxLifetime: 10 * time.Minute, ConnMaxIdleTime: 2 * time.Minute}).(*postgresPersistenceBackend)
	if err := reopened.Init(ctx); err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	quotes, err := reopened.LoadQuotes(ctx)
	if err != nil {
		t.Fatal(err)
	}
	q := quotes["NVDA"]
	if q.Price != quote.Price || q.DataState != "persisted" || q.FeedType != "persisted" {
		t.Fatalf("postgres warm-start retained incorrect truth: %+v", q)
	}
	workspaces, err := reopened.LoadUserWorkspaces(ctx)
	if err != nil || len(workspaces) != 1 || workspaces[0].UserID != "user-v183" {
		t.Fatalf("postgres workspace persistence failure: workspaces=%+v err=%v", workspaces, err)
	}
}

func TestV183PostgresConcurrentWorkspaceWritesStayBounded(t *testing.T) {
	backend := newV183PostgresBackend(t)
	defer backend.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const writers = 24
	var wg sync.WaitGroup
	errs := make(chan error, writers)
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			workspace := defaultUserWorkspace("concurrent-" + randomID("u"))
			workspace.UpdatedAt = time.Now().UnixMilli()
			errs <- backend.SaveUserWorkspace(ctx, workspace)
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent postgres write: %v", err)
		}
	}
	workspaces, err := backend.LoadUserWorkspaces(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(workspaces) != writers {
		t.Fatalf("lost concurrent workspaces: got=%d want=%d", len(workspaces), writers)
	}
	pool := backend.PoolDiagnostics()
	if pool.MaxOpenConnections != 4 || pool.OpenConnections > 4 {
		t.Fatalf("connection pool escaped configured bound: %+v", pool)
	}
	dbdiag := backend.DatabaseDiagnostics()
	if dbdiag.Operations < writers || dbdiag.Errors != 0 {
		t.Fatalf("database diagnostics missing concurrent activity: %+v", dbdiag)
	}
}

func TestV183PostgresConcurrentMigrationInitializationIsSerialized(t *testing.T) {
	url := v183PostgresURL(t)
	seed := newV183PostgresBackend(t)
	if err := seed.Close(); err != nil {
		t.Fatal(err)
	}

	const starters = 3
	var wg sync.WaitGroup
	errs := make(chan error, starters)
	for i := 0; i < starters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			backend := newPostgresPersistenceBackend(postgresPersistenceConfig{DatabaseURL: url, MaxOpenConns: 2, MaxIdleConns: 1, ConnMaxLifetime: 5 * time.Minute, ConnMaxIdleTime: time.Minute}).(*postgresPersistenceBackend)
			err := backend.Init(context.Background())
			if err == nil {
				err = backend.Close()
			}
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent migration init failed: %v", err)
		}
	}
}
