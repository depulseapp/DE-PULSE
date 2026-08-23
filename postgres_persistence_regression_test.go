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

func TestV183PostgresSQLiteArchiveMigratesWithoutLineageLoss(t *testing.T) {
	source := newPersistenceManagerWithBackend(newLocalPersistenceBackend(t.TempDir()))
	defer source.Close()
	if source.backend.Name() != "sqlite" {
		t.Skip("SQLite source backend required for SQLite to PostgreSQL migration proof")
	}
	seedV183ArchiveBackend(t, source.backend)
	path := t.TempDir() + "/sqlite-to-postgres.json"
	archive, err := source.ExportArchiveFile(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	if len(archive.QuoteHistory) < 2 {
		t.Fatalf("source archive did not preserve quote history: %+v", archive.QuoteHistory)
	}

	target := newV183PostgresBackend(t)
	defer target.Close()
	if err := target.RestorePersistenceArchive(context.Background(), archive, persistenceRestoreModeEmpty); err != nil {
		t.Fatalf("SQLite to PostgreSQL restore: %v", err)
	}
	restored, err := target.ExportPersistenceArchive(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(restored.Symbols) != len(archive.Symbols) || len(restored.CanonicalQuotes) != len(archive.CanonicalQuotes) || len(restored.QuoteHistory) != len(archive.QuoteHistory) || len(restored.Evidence) != len(archive.Evidence) || len(restored.Decisions) != len(archive.Decisions) || len(restored.Outcomes) != len(archive.Outcomes) || len(restored.Features) != len(archive.Features) || !restored.HasIdentity || len(restored.UserWorkspaces) != len(archive.UserWorkspaces) {
		t.Fatalf("SQLite to PostgreSQL lineage parity failed: source=%+v target=%+v", archive, restored)
	}
	if restored.Identity.Users[0].PasswordHash != archive.Identity.Users[0].PasswordHash || restored.Identity.Sessions[0].TokenHash != archive.Identity.Sessions[0].TokenHash {
		t.Fatal("identity credential/session hash continuity was not preserved")
	}
	quotes, err := target.LoadQuotes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if q := quotes["NVDA"]; q.Price != 201 || q.DataState != "persisted" || q.FeedType != "persisted" {
		t.Fatalf("migrated warm-start quote is not truthful persisted state: %+v", q)
	}
}

func TestV183PostgresPoolContentionWaitsWithinConfiguredBound(t *testing.T) {
	url := v183PostgresURL(t)
	backend := newPostgresPersistenceBackend(postgresPersistenceConfig{DatabaseURL: url, MaxOpenConns: 2, MaxIdleConns: 2, ConnMaxLifetime: 10 * time.Minute, ConnMaxIdleTime: 2 * time.Minute}).(*postgresPersistenceBackend)
	if err := backend.Init(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer backend.Close()
	resetV183Postgres(t, backend)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn1, err := backend.db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	conn2, err := backend.db.Conn(ctx)
	if err != nil {
		conn1.Close()
		t.Fatal(err)
	}
	defer conn2.Close()

	done := make(chan error, 1)
	go func() {
		workspace := defaultUserWorkspace("pool-wait")
		workspace.UpdatedAt = time.Now().UnixMilli()
		done <- backend.SaveUserWorkspace(ctx, workspace)
	}()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) && backend.PoolDiagnostics().WaitCount == 0 {
		time.Sleep(5 * time.Millisecond)
	}
	pool := backend.PoolDiagnostics()
	if pool.MaxOpenConnections != 2 || pool.OpenConnections > 2 || pool.WaitCount == 0 {
		conn1.Close()
		t.Fatalf("PostgreSQL contention did not stay inside pool bounds: %+v", pool)
	}
	if err := conn1.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("queued database work failed after capacity returned: %v", err)
		}
	case <-ctx.Done():
		t.Fatal("queued database work did not recover after pool capacity returned")
	}
	pool = backend.PoolDiagnostics()
	if pool.OpenConnections > 2 || pool.WaitCount == 0 || pool.WaitDurationMs <= 0 {
		t.Fatalf("PostgreSQL contention observability missing: %+v", pool)
	}
}

func TestV183PostgresHealthCheckReflectsClosedDatabase(t *testing.T) {
	backend := newV183PostgresBackend(t)
	if err := backend.HealthCheck(context.Background()); err != nil {
		t.Fatalf("healthy database probe failed: %v", err)
	}
	if err := backend.Close(); err != nil {
		t.Fatal(err)
	}
	if err := backend.HealthCheck(context.Background()); err == nil {
		t.Fatal("closed PostgreSQL backend reported healthy")
	}
}
