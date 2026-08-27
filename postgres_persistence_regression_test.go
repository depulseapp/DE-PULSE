//go:build postgres

package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

const host012ManagedRecoveryAck = "HOST012_MANAGED_PITR_OPERATOR_DRILL"

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

func host012ManagedRecoveryEnv(t *testing.T, key string) string {
	t.Helper()
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		t.Fatalf("%s is required for the managed recovery drill", key)
	}
	return value
}

func host012ManagedRecoveryDatabaseURL(t *testing.T, key string) string {
	t.Helper()
	raw := host012ManagedRecoveryEnv(t, key)
	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("%s is not a valid PostgreSQL URL: %v", key, err)
	}
	if parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" {
		t.Fatalf("%s must use postgres/postgresql scheme", key)
	}
	host := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	if host == "" || host == "localhost" || host == "host.docker.internal" {
		t.Fatalf("%s must identify a remote managed PostgreSQL host", key)
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
		t.Fatalf("%s loopback PostgreSQL is not managed recovery evidence", key)
	}
	if strings.TrimSpace(parsed.EscapedPath()) == "" || parsed.EscapedPath() == "/" {
		t.Fatalf("%s must identify an explicit database", key)
	}
	if strings.ToLower(strings.TrimSpace(parsed.Query().Get("sslmode"))) != "verify-full" {
		t.Fatalf("%s must use sslmode=verify-full", key)
	}
	return raw
}

func host012ManagedRecoveryDatabaseIdentity(t *testing.T, raw string) string {
	t.Helper()
	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	return strings.ToLower(parsed.Scheme) + "://" + strings.ToLower(parsed.Host) + parsed.EscapedPath()
}

func host012ManagedRecoveryBackend(t *testing.T, databaseURL string) *postgresPersistenceBackend {
	t.Helper()
	backend := newPostgresPersistenceBackend(postgresPersistenceConfig{
		DatabaseURL:     databaseURL,
		MaxOpenConns:    2,
		MaxIdleConns:    1,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: time.Minute,
	}).(*postgresPersistenceBackend)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := backend.Init(ctx); err != nil {
		t.Fatalf("managed postgres init failed: %v", err)
	}
	return backend
}

func host012FindUser(state IdentityPersistentState, userID string) (UserRecord, bool) {
	for _, user := range state.Users {
		if user.ID == userID {
			return user, true
		}
	}
	return UserRecord{}, false
}

func host012WorkspaceContainsPersonalData(workspace UserWorkspace) bool {
	workspace = normalizeUserWorkspace(workspace)
	for _, watchlist := range workspace.Watchlists {
		if len(watchlist.Symbols) > 0 {
			return true
		}
	}
	return workspace.UI.ScopeType != "watchlist" || workspace.UI.WatchlistID != "swing" || normalizeSymbol(workspace.UI.SelectedTicker) != "SPY"
}

func host012WorkspaceForUser(workspaces []UserWorkspace, userID string) (UserWorkspace, bool) {
	for _, workspace := range workspaces {
		if strings.TrimSpace(workspace.UserID) == userID {
			return workspace, true
		}
	}
	return UserWorkspace{}, false
}

func host012UserHash(userID string) string {
	sum := sha256.Sum256([]byte(userID))
	return hex.EncodeToString(sum[:])
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

// TestHOST012ManagedRecoveryPrivacyReplayDrill is intentionally excluded from
// ordinary CI execution by the postgres build tag plus an explicit operator ack.
// It mutates only the throwaway PITR-restored target after first proving that the
// target contains the pre-deletion live account. The source database is read only
// and supplies the later authoritative deletion tombstone state. The canonical
// archive anti-resurrection owner performs the replay; no test-only deletion path
// is allowed to become a second privacy owner.
func TestHOST012ManagedRecoveryPrivacyReplayDrill(t *testing.T) {
	if strings.TrimSpace(os.Getenv("DEPULSE_MANAGED_RECOVERY_ACK")) != host012ManagedRecoveryAck {
		t.Skip("managed PITR operator acknowledgement not present")
	}
	environmentClass := strings.ToLower(host012ManagedRecoveryEnv(t, "DEPULSE_MANAGED_RECOVERY_ENV_CLASS"))
	if environmentClass == "production" {
		t.Fatal("managed recovery drill refuses production targets")
	}
	if environmentClass != "development" && environmentClass != "test" && environmentClass != "stage" {
		t.Fatal("DEPULSE_MANAGED_RECOVERY_ENV_CLASS must be development, test, or stage")
	}
	candidateSHA := strings.ToLower(host012ManagedRecoveryEnv(t, "DEPULSE_MANAGED_RECOVERY_CANDIDATE_SHA"))
	if len(candidateSHA) != 40 {
		t.Fatal("DEPULSE_MANAGED_RECOVERY_CANDIDATE_SHA must be a full Git SHA")
	}
	if _, err := hex.DecodeString(candidateSHA); err != nil {
		t.Fatal("DEPULSE_MANAGED_RECOVERY_CANDIDATE_SHA must be hexadecimal")
	}
	userID := host012ManagedRecoveryEnv(t, "DEPULSE_MANAGED_RECOVERY_USER_ID")
	artifactPath := host012ManagedRecoveryEnv(t, "DEPULSE_MANAGED_RECOVERY_PRIVACY_ARTIFACT_PATH")
	sourceURL := host012ManagedRecoveryDatabaseURL(t, "DEPULSE_MANAGED_RECOVERY_SOURCE_URL")
	restoredURL := host012ManagedRecoveryDatabaseURL(t, "DEPULSE_MANAGED_RECOVERY_RESTORED_URL")
	if host012ManagedRecoveryDatabaseIdentity(t, sourceURL) == host012ManagedRecoveryDatabaseIdentity(t, restoredURL) {
		t.Fatal("managed recovery source and restored target must be distinct databases")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	source := host012ManagedRecoveryBackend(t, sourceURL)
	defer source.Close()
	restored := host012ManagedRecoveryBackend(t, restoredURL)

	sourceIdentity, err := source.LoadIdentityState(ctx)
	if err != nil {
		restored.Close()
		t.Fatalf("load authoritative source identity: %v", err)
	}
	sourceTombstone, ok := accountDeletionTombstones(sourceIdentity)[userID]
	if !ok {
		restored.Close()
		t.Fatal("authoritative source does not contain the required deletion tombstone")
	}
	if sourceTombstone.DisplayName != accountDeletionReasonUserRequest || sourceTombstone.UpdatedAt <= 0 {
		restored.Close()
		t.Fatal("authoritative source tombstone is not the controlled user-request deletion state")
	}
	for _, device := range sourceIdentity.Devices {
		if device.UserID == userID {
			restored.Close()
			t.Fatal("authoritative source still contains a deleted-user device")
		}
	}
	for _, session := range sourceIdentity.Sessions {
		if session.UserID == userID {
			restored.Close()
			t.Fatal("authoritative source still contains a deleted-user session")
		}
	}
	sourceWorkspaces, err := source.LoadUserWorkspaces(ctx)
	if err != nil {
		restored.Close()
		t.Fatalf("load authoritative source workspaces: %v", err)
	}
	if sourceWorkspace, found := host012WorkspaceForUser(sourceWorkspaces, userID); !found || host012WorkspaceContainsPersonalData(sourceWorkspace) {
		restored.Close()
		t.Fatal("authoritative source does not contain the canonical privacy-blank workspace")
	}

	preReplayIdentity, err := restored.LoadIdentityState(ctx)
	if err != nil {
		restored.Close()
		t.Fatalf("load PITR-restored identity: %v", err)
	}
	preReplayUser, found := host012FindUser(preReplayIdentity, userID)
	if !found || isAccountDeletionTombstone(preReplayUser) || preReplayUser.Status != UserActive {
		restored.Close()
		t.Fatal("PITR target does not prove a point before deletion with the live account present")
	}
	if normalizedTenantID(preReplayUser.TenantID) != normalizedTenantID(sourceTombstone.TenantID) {
		restored.Close()
		t.Fatal("PITR-restored account tenant does not match the authoritative tombstone")
	}

	historicalArchive, err := restored.ExportPersistenceArchive(ctx)
	if err != nil {
		restored.Close()
		t.Fatalf("export PITR-restored canonical archive: %v", err)
	}
	if !historicalArchive.HasIdentity {
		restored.Close()
		t.Fatal("PITR-restored archive has no canonical identity state")
	}

	replayStarted := time.Now()
	replayedArchive := enforceArchiveAccountDeletionPrivacy(historicalArchive, sourceIdentity)
	if _, present := accountDeletionTombstones(replayedArchive.Identity)[userID]; !present {
		restored.Close()
		t.Fatal("canonical anti-resurrection owner did not carry the authoritative tombstone")
	}
	if replayWorkspace, found := host012WorkspaceForUser(replayedArchive.UserWorkspaces, userID); !found || host012WorkspaceContainsPersonalData(replayWorkspace) {
		restored.Close()
		t.Fatal("canonical anti-resurrection owner did not privacy-blank the restored workspace")
	}
	if err := restored.RestorePersistenceArchive(ctx, replayedArchive, persistenceRestoreModeReplace); err != nil {
		restored.Close()
		t.Fatalf("apply canonical tombstone replay to PITR target: %v", err)
	}
	replayDurationMillis := time.Since(replayStarted).Milliseconds()
	if err := restored.Close(); err != nil {
		t.Fatalf("close restored target before restart proof: %v", err)
	}

	reopened := host012ManagedRecoveryBackend(t, restoredURL)
	defer reopened.Close()
	if err := reopened.HealthCheck(ctx); err != nil {
		t.Fatalf("restored target health check after replay/restart: %v", err)
	}
	finalIdentity, err := reopened.LoadIdentityState(ctx)
	if err != nil {
		t.Fatalf("load final restored identity: %v", err)
	}
	finalWorkspaces, err := reopened.LoadUserWorkspaces(ctx)
	if err != nil {
		t.Fatalf("load final restored workspaces: %v", err)
	}

	tombstonesPresent := 0
	liveProfiles := 0
	for _, user := range finalIdentity.Users {
		if user.ID != userID {
			continue
		}
		if isAccountDeletionTombstone(user) {
			tombstonesPresent++
		} else {
			liveProfiles++
		}
	}
	activeDevices := 0
	for _, device := range finalIdentity.Devices {
		if device.UserID == userID {
			activeDevices++
		}
	}
	activeSessions := 0
	for _, session := range finalIdentity.Sessions {
		if session.UserID == userID {
			activeSessions++
		}
	}
	personalWorkspaceRows := 0
	if finalWorkspace, found := host012WorkspaceForUser(finalWorkspaces, userID); !found {
		t.Fatal("final restored target lost the canonical blank workspace row")
	} else if host012WorkspaceContainsPersonalData(finalWorkspace) {
		personalWorkspaceRows++
	}
	projection := &IdentityService{state: finalIdentity}
	_, _, projectionErr := projection.accountPrivacyExport(userID)
	accountProjectionChecked := projectionErr != nil

	if tombstonesPresent < 1 || liveProfiles != 0 || activeDevices != 0 || activeSessions != 0 || personalWorkspaceRows != 0 || !accountProjectionChecked {
		t.Fatalf("post-replay privacy verification failed: tombstones=%d liveProfiles=%d devices=%d sessions=%d personalWorkspaceRows=%d projectionDenied=%v", tombstonesPresent, liveProfiles, activeDevices, activeSessions, personalWorkspaceRows, accountProjectionChecked)
	}

	artifact := map[string]any{
		"schema":           "DE.PULSE-HOST012-MANAGED-RECOVERY-PRIVACY-VERIFICATION-1",
		"candidateSha":     candidateSHA,
		"verificationMode": "CANONICAL_ARCHIVE_TOMBSTONE_REPLAY",
		"environmentClass": environmentClass,
		"verifiedAt":       time.Now().UTC().Format(time.RFC3339Nano),
		"userIdSha256":     host012UserHash(userID),
		"sourceTombstone": map[string]any{
			"present":   true,
			"updatedAt": sourceTombstone.UpdatedAt,
			"reason":    sourceTombstone.DisplayName,
		},
		"beforeReplay": map[string]any{
			"liveUserPresent": true,
			"tombstonePresent": false,
		},
		"canonicalReplay": map[string]any{
			"owner":                 "enforceArchiveAccountDeletionPrivacy",
			"restoreMode":           persistenceRestoreModeReplace,
			"durationMilliseconds":  replayDurationMillis,
			"restartVerified":       true,
			"sourceReadOnlyByDrill": true,
		},
		"privacyAssertions": map[string]any{
			"deletedUsersResurrected":              0,
			"liveDeletedProfiles":                  liveProfiles,
			"personalWorkspaceRowsForDeletedUsers": personalWorkspaceRows,
			"activeSessionsForDeletedUsers":        activeSessions,
			"activeDevicesForDeletedUsers":         activeDevices,
			"tombstonesPresentAfterReplay":         tombstonesPresent,
			"canonicalHealthCheckPassed":           true,
			"accountDataProjectionChecked":         accountProjectionChecked,
		},
	}
	raw, err := json.MarshalIndent(artifact, "", "  ")
	if err != nil {
		t.Fatalf("encode managed recovery privacy artifact: %v", err)
	}
	if dir := filepath.Dir(artifactPath); dir != "." {
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatalf("create privacy artifact directory: %v", err)
		}
	}
	if err := os.WriteFile(artifactPath, append(raw, '\n'), 0600); err != nil {
		t.Fatalf("write managed recovery privacy artifact: %v", err)
	}
	t.Log("HOST-012 managed PITR privacy replay drill: PASS (connection details suppressed)")
}
