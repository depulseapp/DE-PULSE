package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	persistenceBackendEnv      = "DEPULSE_PERSISTENCE_BACKEND"
	postgresDatabaseURLEnv     = "DEPULSE_DATABASE_URL"
	postgresMaxOpenConnsEnv    = "DEPULSE_DB_MAX_OPEN_CONNS"
	postgresMaxIdleConnsEnv    = "DEPULSE_DB_MAX_IDLE_CONNS"
	postgresConnMaxLifetimeEnv = "DEPULSE_DB_CONN_MAX_LIFETIME"
	postgresConnMaxIdleTimeEnv = "DEPULSE_DB_CONN_MAX_IDLE_TIME"
)

type postgresPersistenceConfig struct {
	DatabaseURL     string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func newPersistenceBackend(configDir string) PersistenceBackend {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv(persistenceBackendEnv)))
	switch mode {
	case "", "local", "sqlite":
		return newLocalPersistenceBackend(configDir)
	case "postgres", "postgresql":
		return newPostgresPersistenceBackend(postgresPersistenceConfigFromEnv())
	default:
		return newUnavailablePersistenceBackend(fmt.Sprintf("unsupported persistence backend %q", mode))
	}
}

func postgresPersistenceConfigFromEnv() postgresPersistenceConfig {
	maxOpen := boundedEnvInt(postgresMaxOpenConnsEnv, 16, 2, 128)
	maxIdle := boundedEnvInt(postgresMaxIdleConnsEnv, 8, 1, maxOpen)
	return postgresPersistenceConfig{
		DatabaseURL:     strings.TrimSpace(os.Getenv(postgresDatabaseURLEnv)),
		MaxOpenConns:    maxOpen,
		MaxIdleConns:    maxIdle,
		ConnMaxLifetime: durationEnv(postgresConnMaxLifetimeEnv, 30*time.Minute),
		ConnMaxIdleTime: durationEnv(postgresConnMaxIdleTimeEnv, 5*time.Minute),
	}
}

func boundedEnvInt(name string, fallback, minValue, maxValue int) int {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	if v < minValue {
		return minValue
	}
	if v > maxValue {
		return maxValue
	}
	return v
}

func durationEnv(name string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback
	}
	v, err := time.ParseDuration(raw)
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}

type unavailablePersistenceBackend struct {
	reason string
}

func newUnavailablePersistenceBackend(reason string) PersistenceBackend {
	return &unavailablePersistenceBackend{reason: strings.TrimSpace(reason)}
}

func (b *unavailablePersistenceBackend) err() error {
	if b == nil || b.reason == "" {
		return errors.New("persistence backend unavailable")
	}
	return errors.New(b.reason)
}
func (b *unavailablePersistenceBackend) Name() string               { return "unavailable" }
func (b *unavailablePersistenceBackend) Capabilities() []string     { return nil }
func (b *unavailablePersistenceBackend) Init(context.Context) error { return b.err() }
func (b *unavailablePersistenceBackend) UpsertSymbols(context.Context, []SymbolRegistryRecord) (int, error) {
	return 0, b.err()
}
func (b *unavailablePersistenceBackend) LoadSymbols(context.Context) ([]SymbolRegistryRecord, error) {
	return nil, b.err()
}
func (b *unavailablePersistenceBackend) SaveQuotes(context.Context, map[string]Quote) (int, error) {
	return 0, b.err()
}
func (b *unavailablePersistenceBackend) LoadQuotes(context.Context) (map[string]Quote, error) {
	return nil, b.err()
}
func (b *unavailablePersistenceBackend) SaveIntelligence(context.Context, PersistenceIntelligenceBatch) (int, error) {
	return 0, b.err()
}
func (b *unavailablePersistenceBackend) LoadIdentityState(context.Context) (IdentityPersistentState, error) {
	return IdentityPersistentState{}, b.err()
}
func (b *unavailablePersistenceBackend) SaveIdentityState(context.Context, IdentityPersistentState) error {
	return b.err()
}
func (b *unavailablePersistenceBackend) LoadUserWorkspaces(context.Context) ([]UserWorkspace, error) {
	return nil, b.err()
}
func (b *unavailablePersistenceBackend) SaveUserWorkspace(context.Context, UserWorkspace) error {
	return b.err()
}
func (b *unavailablePersistenceBackend) Stats(context.Context) (PersistenceStoreStats, error) {
	return PersistenceStoreStats{}, b.err()
}
func (b *unavailablePersistenceBackend) Close() error { return nil }

func (b *unavailablePersistenceBackend) HealthCheck(context.Context) error { return b.err() }
