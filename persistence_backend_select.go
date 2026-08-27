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

// hostedRightsPersistenceBackend is a narrow policy decorator around the
// existing canonical persistence backend. It owns no storage, migration,
// tenancy, retry or recovery behavior; it only enforces current legal/data
// rights at the existing persistence boundary.
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

func normalizeProviderRightsSource(value string) string {
	replacer := strings.NewReplacer(
		"_", " ",
		"-", " ",
		"/", " ",
		":", " ",
		".", " ",
		"=", " ",
		"|", " ",
		"(", " ",
		")", " ",
	)
	return strings.Join(strings.Fields(strings.ToLower(replacer.Replace(strings.TrimSpace(value)))), " ")
}

func providerRightsSourceContains(source, candidate string) bool {
	candidate = normalizeProviderRightsSource(candidate)
	return candidate != "" && strings.Contains(" "+source+" ", " "+candidate+" ")
}

// providerRightsSourceProvider resolves external evidence through the canonical
// provider-registration inventory before hosted persistence applies legal/data
// rights. Longest registered-name wins so SEC EDGAR remains distinct from SEC.
// The aliases below exist only for known legacy/upstream source spellings; they
// do not create provider registrations or grant rights.
func providerRightsSourceProvider(source string) string {
	normalized := normalizeProviderRightsSource(source)
	if normalized == "" {
		return "—"
	}

	bestProvider := ""
	bestLength := 0
	for _, registration := range providerRegistrations() {
		candidate := normalizeProviderRightsSource(registration.Name)
		if candidate == "" || !providerRightsSourceContains(normalized, candidate) {
			continue
		}
		if len(candidate) > bestLength {
			bestProvider = registration.Name
			bestLength = len(candidate)
		}
	}
	if bestProvider != "" {
		return bestProvider
	}

	compact := strings.ReplaceAll(normalized, " ", "")
	switch {
	case strings.Contains(compact, "twelvedata"):
		return "Twelve Data"
	case providerRightsSourceContains(normalized, "yahoo"):
		return "yfinance"
	case providerRightsSourceContains(normalized, "edgar"):
		return "SEC EDGAR"
	case strings.Contains(normalized, "bureau of labor statistics"):
		return "BLS"
	case strings.Contains(normalized, "energy information administration"):
		return "EIA"
	default:
		return "—"
	}
}

func hostedRightsExternalEvidenceAllowed(record EvidenceRecord, now time.Time) bool {
	if !isHostedRuntime() {
		return true
	}
	source := strings.TrimSpace(record.Source)
	if source == "" {
		return true
	}
	provider := providerRightsSourceProvider(source)
	if provider == "—" {
		// Some canonical/internal evidence uses semantic source labels rather
		// than external-provider identities. HOST-022 owns full point-in-time
		// provenance conservation; every registered external provider is
		// resolved through providerRegistrations above and enforced here.
		return true
	}
	return hostedProviderRightsAllowed(provider, providerHostedUseProductionServing, now)
}

func (b *hostedRightsPersistenceBackend) SaveIntelligence(ctx context.Context, batch PersistenceIntelligenceBatch) (int, error) {
	if !isHostedRuntime() {
		return b.inner.SaveIntelligence(ctx, batch)
	}
	now := time.Now()
	filtered := PersistenceIntelligenceBatch{
		Decisions: append([]DecisionLineageRecord(nil), batch.Decisions...),
		Outcomes:  append([]OutcomeHistoryRecord(nil), batch.Outcomes...),
		Features:  append([]DerivedFeatureRecord(nil), batch.Features...),
	}
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
