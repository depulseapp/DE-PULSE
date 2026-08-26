package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"depulse/internal/hostedpersistence"
)

const (
	persistenceBackendEnv        = "DEPULSE_PERSISTENCE_BACKEND"
	postgresDatabaseURLEnv       = "DEPULSE_DATABASE_URL"
	postgresMaxOpenConnsEnv      = "DEPULSE_DB_MAX_OPEN_CONNS"
	postgresMaxIdleConnsEnv      = "DEPULSE_DB_MAX_IDLE_CONNS"
	postgresConnMaxLifetimeEnv   = "DEPULSE_DB_CONN_MAX_LIFETIME"
	postgresConnMaxIdleTimeEnv   = "DEPULSE_DB_CONN_MAX_IDLE_TIME"
	postgresPolicyVersionEnv     = "DEPULSE_DB_POLICY_VERSION"
	postgresRLSDispositionEnv    = "DEPULSE_DB_RLS_DISPOSITION"
	postgresMigrationStrategyEnv = "DEPULSE_DB_MIGRATION_STRATEGY"
	postgresHAReadyEnv           = "DEPULSE_DB_HA_READY"
	postgresBackupEncryptedEnv   = "DEPULSE_DB_BACKUP_ENCRYPTED"
	postgresPITRReadyEnv         = "DEPULSE_DB_PITR_READY"
	postgresRPOEnv               = "DEPULSE_DB_RPO"
	postgresRTOEnv               = "DEPULSE_DB_RTO"
	postgresRestoreDrillAtEnv    = "DEPULSE_DB_RESTORE_DRILL_AT"
	postgresRestoreDrillIDEnv    = "DEPULSE_DB_RESTORE_DRILL_ID"
	postgresRollbackPolicyEnv    = "DEPULSE_DB_ROLLBACK_POLICY"
	postgresRollbackVerifiedEnv  = "DEPULSE_DB_ROLLBACK_VERIFIED"
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
	var backend PersistenceBackend
	switch mode {
	case "", "local", "sqlite":
		backend = newLocalPersistenceBackend(configDir)
	case "postgres", "postgresql":
		config := postgresPersistenceConfigFromEnv()
		if isHostedRuntime() {
			if err := hostedpersistence.Validate(hostedPostgresRuntimeDeclaration(config), time.Now()); err != nil {
				return newUnavailablePersistenceBackend("hosted PostgreSQL policy: " + err.Error())
			}
		}
		backend = newPostgresPersistenceBackend(config)
	default:
		return newUnavailablePersistenceBackend(fmt.Sprintf("unsupported persistence backend %q", mode))
	}
	return wrapHostedRightsPersistenceBackend(backend)
}

func hostedPostgresRuntimeDeclaration(config postgresPersistenceConfig) hostedpersistence.RuntimeDeclaration {
	return hostedpersistence.RuntimeDeclaration{
		Environment:        strings.TrimSpace(os.Getenv(hostedEnvironmentEnv)),
		DatabaseURL:        config.DatabaseURL,
		MaxOpenConnections: config.MaxOpenConns,
		PolicyVersion:      strings.TrimSpace(os.Getenv(postgresPolicyVersionEnv)),
		RLSDisposition:     strings.TrimSpace(os.Getenv(postgresRLSDispositionEnv)),
		MigrationStrategy:  strings.TrimSpace(os.Getenv(postgresMigrationStrategyEnv)),
		HAReady:            strictBoolEnv(postgresHAReadyEnv),
		BackupEncrypted:    strictBoolEnv(postgresBackupEncryptedEnv),
		PITRReady:          strictBoolEnv(postgresPITRReadyEnv),
		RPO:                strings.TrimSpace(os.Getenv(postgresRPOEnv)),
		RTO:                strings.TrimSpace(os.Getenv(postgresRTOEnv)),
		RestoreDrillAt:     strings.TrimSpace(os.Getenv(postgresRestoreDrillAtEnv)),
		RestoreDrillID:     strings.TrimSpace(os.Getenv(postgresRestoreDrillIDEnv)),
		RollbackPolicy:     strings.TrimSpace(os.Getenv(postgresRollbackPolicyEnv)),
		RollbackVerified:   strictBoolEnv(postgresRollbackVerifiedEnv),
	}
}

func strictBoolEnv(name string) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(name))) {
	case "1", "true", "yes", "on", "required", "verified":
		return true
	default:
		return false
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
