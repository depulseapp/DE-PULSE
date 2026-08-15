package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
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
