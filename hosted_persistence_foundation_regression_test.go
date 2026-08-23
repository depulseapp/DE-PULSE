package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestV183DesktopPersistenceRemainsLocalByDefault(t *testing.T) {
	t.Setenv(persistenceBackendEnv, "")
	t.Setenv(postgresDatabaseURLEnv, "")
	backend := newPersistenceBackend(t.TempDir())
	if backend.Name() == "postgresql" || backend.Name() == "unavailable" {
		t.Fatalf("default desktop persistence changed unexpectedly: %s", backend.Name())
	}
	if err := backend.Init(context.Background()); err != nil {
		t.Fatalf("default local persistence unavailable: %v", err)
	}
	_ = backend.Close()
}

func TestV183PostgresSelectionFailsClosedWhenUnavailableOrUnconfigured(t *testing.T) {
	t.Setenv(persistenceBackendEnv, "postgres")
	t.Setenv(postgresDatabaseURLEnv, "")
	backend := newPersistenceBackend(t.TempDir())
	if backend.Name() == "sqlite" || backend.Name() == "file-fallback" {
		t.Fatalf("postgres selection silently fell back to local backend: %s", backend.Name())
	}
	if err := backend.Init(context.Background()); err == nil {
		t.Fatal("postgres selection without a configured hosted database should fail closed")
	}
}

func TestV183PostgresPoolConfigurationIsBounded(t *testing.T) {
	t.Setenv(postgresMaxOpenConnsEnv, "9999")
	t.Setenv(postgresMaxIdleConnsEnv, "9999")
	t.Setenv(postgresConnMaxLifetimeEnv, "45m")
	t.Setenv(postgresConnMaxIdleTimeEnv, "7m")
	cfg := postgresPersistenceConfigFromEnv()
	if cfg.MaxOpenConns != 128 || cfg.MaxIdleConns != 128 {
		t.Fatalf("pool bounds not enforced: %+v", cfg)
	}
	if cfg.ConnMaxLifetime != 45*time.Minute || cfg.ConnMaxIdleTime != 7*time.Minute {
		t.Fatalf("pool durations not parsed: %+v", cfg)
	}
}

func TestV183HostedRuntimeAddressAndConfigAreExplicit(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(runtimeModeEnv, "hosted")
	t.Setenv(hostedConfigDirEnv, dir)
	t.Setenv(hostedListenAddrEnv, "")
	t.Setenv("PORT", "9091")
	if runtimeMode() != "hosted" || hostedListenAddress() != ":9091" {
		t.Fatalf("hosted runtime selection incorrect: mode=%s listen=%s", runtimeMode(), hostedListenAddress())
	}
	resolved, err := resolveV18RuntimeConfig(filepath.Dir(dir))
	if err != nil {
		t.Fatal(err)
	}
	if resolved != dir {
		t.Fatalf("hosted config dir mismatch: got=%s want=%s", resolved, dir)
	}
}

func TestV183ReadinessRequiresCanonicalPersistenceAndIdentity(t *testing.T) {
	p := NewPersistenceManager(t.TempDir())
	defer p.Close()
	identity, err := NewIdentityService(p)
	if err != nil {
		t.Fatal(err)
	}
	app := &Application{persistence: p, identity: identity, hub: NewHub(), state: defaultState(), aiCache: map[string]aiCacheEntry{}, httpTelemetry: NewRequestTelemetry()}
	req := httptest.NewRequest(http.MethodGet, "/api/ready", nil)
	rr := httptest.NewRecorder()
	app.routes().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("ready endpoint rejected healthy application: status=%d body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"ok":true`) || !strings.Contains(rr.Body.String(), `"backend"`) {
		t.Fatalf("ready response missing truthful persistence state: %s", rr.Body.String())
	}

	app.identity = nil
	rr = httptest.NewRecorder()
	app.routes().ServeHTTP(rr, req)
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("readiness must fail when identity is unavailable: status=%d", rr.Code)
	}
}

type v183HealthBackend struct {
	PersistenceBackend
	mu           sync.Mutex
	unhealthy    bool
	healthChecks int
}

func (b *v183HealthBackend) setUnhealthy(v bool) {
	b.mu.Lock()
	b.unhealthy = v
	b.mu.Unlock()
}

func (b *v183HealthBackend) isUnhealthy() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.unhealthy
}

func (b *v183HealthBackend) SaveQuotes(ctx context.Context, quotes map[string]Quote) (int, error) {
	if b.isUnhealthy() {
		return 0, errors.New("simulated database unavailable")
	}
	return b.PersistenceBackend.SaveQuotes(ctx, quotes)
}

func (b *v183HealthBackend) HealthCheck(ctx context.Context) error {
	b.mu.Lock()
	b.healthChecks++
	unhealthy := b.unhealthy
	b.mu.Unlock()
	if unhealthy {
		return errors.New("simulated database unavailable")
	}
	if health, ok := b.PersistenceBackend.(persistenceHealthBackend); ok {
		return health.HealthCheck(ctx)
	}
	return nil
}

func TestV183PersistenceOutageReadinessRecoversWithBoundedHysteresis(t *testing.T) {
	inner := newLocalPersistenceBackend(t.TempDir())
	backend := &v183HealthBackend{PersistenceBackend: inner, unhealthy: true}
	p := newPersistenceManagerWithBackend(backend)
	defer p.Close()
	p.mu.Lock()
	p.retryBackoff = 10 * time.Millisecond
	p.mu.Unlock()

	p.EnqueueQuotes(map[string]Quote{"NVDA": {Symbol: "NVDA", Price: 201, Source: "test", DataState: "live", UpdatedAt: time.Now().UnixMilli()}})
	deadline := time.Now().Add(750 * time.Millisecond)
	for time.Now().Before(deadline) {
		d := p.Diagnostics()
		if !d.Ready && d.HealthState == "DEGRADED" && d.RetryScheduled && d.RetryBackoffMs >= 10 {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	d := p.Diagnostics()
	if d.Ready || d.HealthState != "DEGRADED" || !d.RetryScheduled {
		t.Fatalf("runtime database failure did not truthfully degrade readiness: %+v", d)
	}
	identity, err := NewIdentityService(p)
	if err != nil {
		t.Fatal(err)
	}
	app := &Application{persistence: p, identity: identity, hub: NewHub(), state: defaultState(), aiCache: map[string]aiCacheEntry{}, httpTelemetry: NewRequestTelemetry()}
	req := httptest.NewRequest(http.MethodGet, "/api/ready", nil)
	rr := httptest.NewRecorder()
	app.routes().ServeHTTP(rr, req)
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("hosted readiness remained healthy during database outage: status=%d body=%s", rr.Code, rr.Body.String())
	}
	if d.QueueDepth == 0 {
		t.Fatalf("failed canonical write was not retained for recovery: %+v", d)
	}

	backend.setUnhealthy(false)
	deadline = time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		d = p.Diagnostics()
		if d.Ready && d.HealthState == "HEALTHY" && d.WriteBatches >= 1 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	d = p.Diagnostics()
	if !d.Ready || d.HealthState != "HEALTHY" || d.WriteBatches < 1 || d.QueueDepth != 0 {
		t.Fatalf("persistence did not recover with queued work intact: %+v", d)
	}
	rr = httptest.NewRecorder()
	app.routes().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("hosted readiness did not recover after database health returned: status=%d body=%s", rr.Code, rr.Body.String())
	}
	backend.mu.Lock()
	healthChecks := backend.healthChecks
	backend.mu.Unlock()
	if healthChecks < 2 || healthChecks > 10 {
		t.Fatalf("health probing was not bounded/hysteretic: checks=%d diagnostics=%+v", healthChecks, d)
	}
}

func seedV183ArchiveBackend(t *testing.T, backend PersistenceBackend) {
	t.Helper()
	ctx := context.Background()
	now := time.Now().UnixMilli()
	if _, err := backend.UpsertSymbols(ctx, []SymbolRegistryRecord{{Symbol: "NVDA", FirstSeenAt: now - 10000, LastSeenAt: now, Active: true, Selected: true, ProcessingTier: 0, DeskMembership: `["swing"]`, ProviderEligible: true}}); err != nil {
		t.Fatal(err)
	}
	first := Quote{Symbol: "NVDA", Price: 200, Bid: 199.9, Ask: 200.1, Volume: 1000, Source: "alpaca-iex", FeedType: "websocket", DataState: "live", ProviderTimestamp: now - 600000, UpdatedAt: now - 599900}
	second := Quote{Symbol: "NVDA", Price: 201, Bid: 200.9, Ask: 201.1, Volume: 1500, Source: "alpaca-iex", FeedType: "websocket", DataState: "live", ProviderTimestamp: now, UpdatedAt: now}
	if _, err := backend.SaveQuotes(ctx, map[string]Quote{"NVDA": first}); err != nil {
		t.Fatal(err)
	}
	if _, err := backend.SaveQuotes(ctx, map[string]Quote{"NVDA": second}); err != nil {
		t.Fatal(err)
	}
	batch := PersistenceIntelligenceBatch{
		Evidence:  []EvidenceRecord{{ID: "ev-archive", Symbol: "NVDA", Kind: "research", ObservedAt: now, Source: "canonical", Provenance: "v18.3-test", FreshnessState: "CURRENT", Payload: json.RawMessage(`{"price":201}`)}},
		Decisions: []DecisionLineageRecord{{ID: "dec-archive", Symbol: "NVDA", Horizon: "swing", EvidenceID: "ev-archive", DecisionKind: "readiness", DecisionValue: "WATCH", FormulaVersion: "protected", CreatedAt: now, Payload: json.RawMessage(`{"reason":"test"}`)}},
		Outcomes:  []OutcomeHistoryRecord{{ID: "out-archive", DecisionID: "dec-archive", Symbol: "NVDA", Horizon: "swing", ObservedAt: now, OutcomeLabel: "PENDING", Payload: json.RawMessage(`{"window":"open"}`)}},
		Features:  []DerivedFeatureRecord{{Symbol: "NVDA", FeatureKey: "behavior", FeatureVersion: "v1", AsOf: now, SourceHash: "archive", Payload: json.RawMessage(`{"score":1}`)}},
	}
	if _, err := backend.SaveIntelligence(ctx, batch); err != nil {
		t.Fatal(err)
	}
	identity := IdentityPersistentState{Version: 1, Users: []UserRecord{{ID: "user-archive", Username: "archive", Role: RoleUser, Status: UserActive, PasswordHash: "argon2id$redacted-test-hash", CreatedAt: now, UpdatedAt: now}}, Sessions: []SessionRecord{{ID: "session-archive", TokenHash: "sha256-test-token", UserID: "user-archive", CreatedAt: now, LastSeenAt: now, IdleExpiresAt: now + 10000, AbsoluteExpiresAt: now + 20000}}, UpdatedAt: now}
	if err := backend.SaveIdentityState(ctx, identity); err != nil {
		t.Fatal(err)
	}
	workspace := defaultUserWorkspace("user-archive")
	workspace.Watchlists[0].Symbols = []string{"NVDA", "TSLA"}
	workspace.UpdatedAt = now
	if err := backend.SaveUserWorkspace(ctx, workspace); err != nil {
		t.Fatal(err)
	}
}

func TestV183PersistenceArchiveIntegrityAndRestoreSafety(t *testing.T) {
	source := newPersistenceManagerWithBackend(newLocalPersistenceBackend(t.TempDir()))
	defer source.Close()
	seedV183ArchiveBackend(t, source.backend)
	path := filepath.Join(t.TempDir(), "depulse-persistence-backup.json")
	archive, err := source.ExportArchiveFile(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	if archive.SchemaVersion != persistenceArchiveSchemaVersion || archive.SourceBackend == "" || len(archive.Symbols) != 1 || len(archive.CanonicalQuotes) != 1 || len(archive.Evidence) != 1 || len(archive.Decisions) != 1 || len(archive.Outcomes) != 1 || len(archive.Features) != 1 || !archive.HasIdentity || len(archive.UserWorkspaces) != 1 {
		t.Fatalf("archive omitted canonical state: %+v", archive)
	}
	if archive.SourceBackend == "sqlite" && len(archive.QuoteHistory) < 2 {
		t.Fatalf("sqlite archive lost quote history: %+v", archive.QuoteHistory)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0077 != 0 {
		t.Fatalf("security-sensitive archive permissions are too broad: %o", info.Mode().Perm())
	}

	target := newPersistenceManagerWithBackend(newLocalPersistenceBackend(t.TempDir()))
	defer target.Close()
	if _, err := target.RestoreArchiveFile(context.Background(), path, persistenceRestoreModeEmpty); err != nil {
		t.Fatal(err)
	}
	restored, err := target.backend.(persistenceArchiveBackend).ExportPersistenceArchive(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(restored.Symbols) != len(archive.Symbols) || len(restored.CanonicalQuotes) != len(archive.CanonicalQuotes) || len(restored.Evidence) != len(archive.Evidence) || len(restored.Decisions) != len(archive.Decisions) || len(restored.Outcomes) != len(archive.Outcomes) || len(restored.Features) != len(archive.Features) || !restored.HasIdentity || len(restored.UserWorkspaces) != len(archive.UserWorkspaces) {
		t.Fatalf("restore parity failure: source=%+v restored=%+v", archive, restored)
	}
	if _, err := target.RestoreArchiveFile(context.Background(), path, persistenceRestoreModeEmpty); err == nil {
		t.Fatal("empty-only restore must reject a non-empty target")
	}
	if _, err := target.RestoreArchiveFile(context.Background(), path, persistenceRestoreModeReplace); err != nil {
		t.Fatalf("explicit replace restore failed: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for i := len(raw) / 2; i < len(raw); i++ {
		if raw[i] >= '0' && raw[i] <= '9' {
			raw[i] = '9' - (raw[i] - '0')
			break
		}
	}
	tampered := filepath.Join(t.TempDir(), "tampered.json")
	if err := os.WriteFile(tampered, raw, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := readPersistenceArchiveFile(tampered); err == nil {
		t.Fatal("tampered persistence archive passed integrity validation")
	}
}

func TestV183PersistenceIntelligenceQueueDeduplicatesAndStaysBounded(t *testing.T) {
	batch := PersistenceIntelligenceBatch{}
	for i := 0; i < maxPendingIntelligenceRecords+100; i++ {
		batch.Features = append(batch.Features, DerivedFeatureRecord{Symbol: "NVDA", FeatureKey: "f", FeatureVersion: strings.Repeat("x", 1) + string(rune(i+1)), AsOf: int64(i)})
	}
	batch.Evidence = append(batch.Evidence, EvidenceRecord{ID: "same", Kind: "test"}, EvidenceRecord{ID: "same", Kind: "test-new"})
	compacted, shed := compactPersistenceIntelligence(batch, maxPendingIntelligenceRecords)
	if compacted.Len() > maxPendingIntelligenceRecords || shed <= 0 {
		t.Fatalf("persistence intelligence queue was not bounded: len=%d shed=%d", compacted.Len(), shed)
	}
	count := 0
	for _, ev := range compacted.Evidence {
		if ev.ID == "same" {
			count++
			if ev.Kind != "test-new" {
				t.Fatalf("dedupe did not retain latest evidence: %+v", ev)
			}
		}
	}
	if count != 1 {
		t.Fatalf("evidence dedupe failed: count=%d", count)
	}
}
