//go:build postgres

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const postgresSchemaLockKey int64 = 18030001

type postgresPersistenceBackend struct {
	config postgresPersistenceConfig
	db     *sql.DB
	// Set once by the hosted tenant wrapper so archive operations use tenant
	// tables inside the same database transaction as all other canonical data.
	tenantScopedArchive bool
	mu                  sync.Mutex
	diag                PersistenceDatabaseDiagnostics
}

func newPostgresPersistenceBackend(config postgresPersistenceConfig) PersistenceBackend {
	return &postgresPersistenceBackend{config: config}
}

func (b *postgresPersistenceBackend) Name() string { return "postgresql" }
func (b *postgresPersistenceBackend) Capabilities() []string {
	return []string{
		"global-symbol-registry",
		"canonical-quotes",
		"quote-history",
		"evidence-records",
		"decision-lineage",
		"outcome-history",
		"derived-feature-store",
		"user-workspaces",
		"hosted-shared-state",
		"transactions",
		"connection-pooling",
	}
}

func (b *postgresPersistenceBackend) observe(operation string, started time.Time, err error) {
	duration := time.Since(started).Milliseconds()
	b.mu.Lock()
	defer b.mu.Unlock()
	b.diag.Operations++
	b.diag.LastOperation = operation
	b.diag.LastOperationDurationMs = duration
	if duration > b.diag.MaxOperationDurationMs {
		b.diag.MaxOperationDurationMs = duration
	}
	if duration >= 250 {
		b.diag.SlowOperations++
	}
	if err != nil {
		b.diag.Errors++
	}
}

func (b *postgresPersistenceBackend) DatabaseDiagnostics() PersistenceDatabaseDiagnostics {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.diag
}

func (b *postgresPersistenceBackend) PoolDiagnostics() PersistencePoolDiagnostics {
	if b == nil || b.db == nil {
		return PersistencePoolDiagnostics{}
	}
	s := b.db.Stats()
	return PersistencePoolDiagnostics{
		MaxOpenConnections: s.MaxOpenConnections,
		OpenConnections:    s.OpenConnections,
		InUse:              s.InUse,
		Idle:               s.Idle,
		WaitCount:          s.WaitCount,
		WaitDurationMs:     s.WaitDuration.Milliseconds(),
		MaxIdleClosed:      s.MaxIdleClosed,
		MaxIdleTimeClosed:  s.MaxIdleTimeClosed,
		MaxLifetimeClosed:  s.MaxLifetimeClosed,
	}
}

func (b *postgresPersistenceBackend) Init(ctx context.Context) (err error) {
	started := time.Now()
	defer func() { b.observe("init", started, err) }()
	if strings.TrimSpace(b.config.DatabaseURL) == "" {
		return errors.New("DEPULSE_DATABASE_URL is required when PostgreSQL persistence is selected")
	}
	db, err := sql.Open("pgx", b.config.DatabaseURL)
	if err != nil {
		return err
	}
	b.db = db
	b.db.SetMaxOpenConns(b.config.MaxOpenConns)
	b.db.SetMaxIdleConns(b.config.MaxIdleConns)
	b.db.SetConnMaxLifetime(b.config.ConnMaxLifetime)
	b.db.SetConnMaxIdleTime(b.config.ConnMaxIdleTime)

	pingCtx, pingCancel := context.WithTimeout(ctx, 8*time.Second)
	if err := b.db.PingContext(pingCtx); err != nil {
		pingCancel()
		_ = b.db.Close()
		b.db = nil
		return fmt.Errorf("postgres ping: %w", err)
	}
	pingCancel()
	migrationCtx, migrationCancel := context.WithTimeout(ctx, 30*time.Second)
	defer migrationCancel()
	if err := b.applyMigrations(migrationCtx); err != nil {
		_ = b.db.Close()
		b.db = nil
		return err
	}
	return nil
}

type postgresMigration struct {
	version    int
	statements []string
}

func postgresMigrations() []postgresMigration {
	return []postgresMigration{
		{1, []string{
			`CREATE TABLE IF NOT EXISTS symbol_registry(
 symbol TEXT PRIMARY KEY,
 first_seen_ms BIGINT NOT NULL,
 last_seen_ms BIGINT NOT NULL,
 active BOOLEAN NOT NULL,
 selected BOOLEAN NOT NULL,
 processing_tier INTEGER NOT NULL,
 desk_membership TEXT NOT NULL DEFAULT '[]',
 provider_eligible BOOLEAN NOT NULL,
 last_subscribed_ms BIGINT NOT NULL DEFAULT 0,
 last_processed_ms BIGINT NOT NULL DEFAULT 0
)`,
			`CREATE INDEX IF NOT EXISTS idx_symbol_registry_active_tier ON symbol_registry(active,processing_tier)`,
			`CREATE TABLE IF NOT EXISTS canonical_quotes(
 symbol TEXT PRIMARY KEY,
 payload_json JSONB NOT NULL,
 provider_timestamp_ms BIGINT NOT NULL DEFAULT 0,
 received_timestamp_ms BIGINT NOT NULL DEFAULT 0,
 persisted_at_ms BIGINT NOT NULL,
 source TEXT NOT NULL DEFAULT '',
 data_state TEXT NOT NULL DEFAULT ''
)`,
			`CREATE INDEX IF NOT EXISTS idx_canonical_quotes_persisted ON canonical_quotes(persisted_at_ms)`,
		}},
		{2, []string{
			`CREATE TABLE IF NOT EXISTS quote_history(
 symbol TEXT NOT NULL,
 bucket_ms BIGINT NOT NULL,
 provider_timestamp_ms BIGINT NOT NULL DEFAULT 0,
 received_timestamp_ms BIGINT NOT NULL DEFAULT 0,
 price DOUBLE PRECISION NOT NULL,
 bid DOUBLE PRECISION NOT NULL DEFAULT 0,
 ask DOUBLE PRECISION NOT NULL DEFAULT 0,
 volume DOUBLE PRECISION NOT NULL DEFAULT 0,
 source TEXT NOT NULL DEFAULT '',
 data_state TEXT NOT NULL DEFAULT '',
 PRIMARY KEY(symbol,bucket_ms)
)`,
			`CREATE INDEX IF NOT EXISTS idx_quote_history_bucket ON quote_history(bucket_ms)`,
			`CREATE TABLE IF NOT EXISTS evidence_records(
 evidence_id TEXT PRIMARY KEY,
 symbol TEXT NOT NULL DEFAULT '',
 evidence_kind TEXT NOT NULL,
 observed_at_ms BIGINT NOT NULL,
 source TEXT NOT NULL DEFAULT '',
 provenance TEXT NOT NULL DEFAULT '',
 freshness_state TEXT NOT NULL DEFAULT '',
 payload_json JSONB NOT NULL
)`,
			`CREATE INDEX IF NOT EXISTS idx_evidence_symbol_time ON evidence_records(symbol,observed_at_ms DESC)`,
			`CREATE TABLE IF NOT EXISTS decision_lineage(
 decision_id TEXT PRIMARY KEY,
 symbol TEXT NOT NULL DEFAULT '',
 horizon TEXT NOT NULL DEFAULT '',
 evidence_id TEXT NOT NULL DEFAULT '',
 decision_kind TEXT NOT NULL,
 decision_value TEXT NOT NULL DEFAULT '',
 formula_version TEXT NOT NULL DEFAULT '',
 created_at_ms BIGINT NOT NULL,
 payload_json JSONB NOT NULL
)`,
			`CREATE INDEX IF NOT EXISTS idx_decision_lineage_symbol_time ON decision_lineage(symbol,created_at_ms DESC)`,
			`CREATE TABLE IF NOT EXISTS outcome_history(
 outcome_id TEXT PRIMARY KEY,
 decision_id TEXT NOT NULL DEFAULT '',
 symbol TEXT NOT NULL DEFAULT '',
 horizon TEXT NOT NULL DEFAULT '',
 observed_at_ms BIGINT NOT NULL,
 outcome_label TEXT NOT NULL DEFAULT '',
 payload_json JSONB NOT NULL
)`,
			`CREATE INDEX IF NOT EXISTS idx_outcome_symbol_time ON outcome_history(symbol,observed_at_ms DESC)`,
			`CREATE TABLE IF NOT EXISTS derived_features(
 symbol TEXT NOT NULL,
 feature_key TEXT NOT NULL,
 feature_version TEXT NOT NULL,
 as_of_ms BIGINT NOT NULL,
 source_hash TEXT NOT NULL DEFAULT '',
 payload_json JSONB NOT NULL,
 PRIMARY KEY(symbol,feature_key,feature_version)
)`,
			`CREATE INDEX IF NOT EXISTS idx_derived_features_asof ON derived_features(as_of_ms DESC)`,
		}},
		{3, []string{
			`CREATE TABLE IF NOT EXISTS identity_state(
 id SMALLINT PRIMARY KEY CHECK(id=1),
 payload_json JSONB NOT NULL,
 updated_at_ms BIGINT NOT NULL
)`,
		}},
		{4, []string{
			`CREATE TABLE IF NOT EXISTS user_workspaces(
 user_id TEXT PRIMARY KEY,
 payload_json JSONB NOT NULL,
 updated_at_ms BIGINT NOT NULL
)`,
			`CREATE INDEX IF NOT EXISTS idx_user_workspaces_updated ON user_workspaces(updated_at_ms DESC)`,
		}},
	}
}

func (b *postgresPersistenceBackend) applyMigrations(ctx context.Context) error {
	if b.db == nil {
		return errors.New("postgres database is not open")
	}
	conn, err := b.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("postgres migration connection: %w", err)
	}
	defer conn.Close()
	if _, err := conn.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations(version INTEGER PRIMARY KEY, applied_at_ms BIGINT NOT NULL)`); err != nil {
		return fmt.Errorf("postgres migration table: %w", err)
	}
	if _, err := conn.ExecContext(ctx, `SELECT pg_advisory_lock($1)`, postgresSchemaLockKey); err != nil {
		return fmt.Errorf("postgres migration lock: %w", err)
	}
	defer func() {
		unlockCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, _ = conn.ExecContext(unlockCtx, `SELECT pg_advisory_unlock($1)`, postgresSchemaLockKey)
	}()

	for _, migration := range postgresMigrations() {
		var applied bool
		err := conn.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)`, migration.version).Scan(&applied)
		if err != nil {
			return fmt.Errorf("postgres migration %d check: %w", migration.version, err)
		}
		if applied {
			continue
		}
		tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			return fmt.Errorf("postgres migration %d begin: %w", migration.version, err)
		}
		for _, statement := range migration.statements {
			if _, err := tx.ExecContext(ctx, statement); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("postgres migration %d: %w", migration.version, err)
			}
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations(version,applied_at_ms) VALUES($1,$2)`, migration.version, time.Now().UnixMilli()); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("postgres migration %d record: %w", migration.version, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("postgres migration %d commit: %w", migration.version, err)
		}
	}
	return nil
}

func (b *postgresPersistenceBackend) UpsertSymbols(ctx context.Context, records []SymbolRegistryRecord) (written int, err error) {
	started := time.Now()
	defer func() { b.observe("upsert-symbols", started, err) }()
	if b.db == nil {
		return 0, errors.New("postgres database is not open")
	}
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	if _, err = tx.ExecContext(ctx, `UPDATE symbol_registry SET active=FALSE, selected=FALSE`); err != nil {
		return 0, err
	}
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO symbol_registry(symbol,first_seen_ms,last_seen_ms,active,selected,processing_tier,desk_membership,provider_eligible,last_subscribed_ms,last_processed_ms)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
ON CONFLICT(symbol) DO UPDATE SET last_seen_ms=EXCLUDED.last_seen_ms,active=EXCLUDED.active,selected=EXCLUDED.selected,processing_tier=EXCLUDED.processing_tier,desk_membership=EXCLUDED.desk_membership,provider_eligible=EXCLUDED.provider_eligible,last_subscribed_ms=GREATEST(symbol_registry.last_subscribed_ms,EXCLUDED.last_subscribed_ms),last_processed_ms=GREATEST(symbol_registry.last_processed_ms,EXCLUDED.last_processed_ms)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	for _, r := range records {
		if err = ctx.Err(); err != nil {
			return written, err
		}
		r.Symbol = normalizeSymbol(r.Symbol)
		if r.Symbol == "" {
			continue
		}
		firstSeen := r.FirstSeenAt
		if firstSeen <= 0 {
			firstSeen = r.LastSeenAt
		}
		if _, err = stmt.ExecContext(ctx, r.Symbol, firstSeen, r.LastSeenAt, r.Active, r.Selected, r.ProcessingTier, r.DeskMembership, r.ProviderEligible, r.LastSubscribedAt, r.LastProcessedAt); err != nil {
			return written, err
		}
		written++
	}
	if err = tx.Commit(); err != nil {
		return written, err
	}
	return written, nil
}

func (b *postgresPersistenceBackend) LoadSymbols(ctx context.Context) (out []SymbolRegistryRecord, err error) {
	started := time.Now()
	defer func() { b.observe("load-symbols", started, err) }()
	rows, err := b.db.QueryContext(ctx, `SELECT symbol,first_seen_ms,last_seen_ms,active,selected,processing_tier,desk_membership,provider_eligible,last_subscribed_ms,last_processed_ms FROM symbol_registry ORDER BY symbol`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var r SymbolRegistryRecord
		if err = rows.Scan(&r.Symbol, &r.FirstSeenAt, &r.LastSeenAt, &r.Active, &r.Selected, &r.ProcessingTier, &r.DeskMembership, &r.ProviderEligible, &r.LastSubscribedAt, &r.LastProcessedAt); err != nil {
			return nil, err
		}
		r.Symbol = normalizeSymbol(r.Symbol)
		out = append(out, r)
	}
	return out, rows.Err()
}

func (b *postgresPersistenceBackend) SaveQuotes(ctx context.Context, quotes map[string]Quote) (written int, err error) {
	started := time.Now()
	defer func() { b.observe("save-quotes", started, err) }()
	if len(quotes) == 0 {
		return 0, nil
	}
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	quoteStmt, err := tx.PrepareContext(ctx, `INSERT INTO canonical_quotes(symbol,payload_json,provider_timestamp_ms,received_timestamp_ms,persisted_at_ms,source,data_state)
VALUES($1,$2::jsonb,$3,$4,$5,$6,$7)
ON CONFLICT(symbol) DO UPDATE SET payload_json=EXCLUDED.payload_json,provider_timestamp_ms=EXCLUDED.provider_timestamp_ms,received_timestamp_ms=EXCLUDED.received_timestamp_ms,persisted_at_ms=EXCLUDED.persisted_at_ms,source=EXCLUDED.source,data_state=EXCLUDED.data_state`)
	if err != nil {
		return 0, err
	}
	defer quoteStmt.Close()
	historyStmt, err := tx.PrepareContext(ctx, `INSERT INTO quote_history(symbol,bucket_ms,provider_timestamp_ms,received_timestamp_ms,price,bid,ask,volume,source,data_state)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
ON CONFLICT(symbol,bucket_ms) DO UPDATE SET provider_timestamp_ms=EXCLUDED.provider_timestamp_ms,received_timestamp_ms=EXCLUDED.received_timestamp_ms,price=EXCLUDED.price,bid=EXCLUDED.bid,ask=EXCLUDED.ask,volume=EXCLUDED.volume,source=EXCLUDED.source,data_state=EXCLUDED.data_state`)
	if err != nil {
		return 0, err
	}
	defer historyStmt.Close()
	now := time.Now().UnixMilli()
	for sym, q := range quotes {
		if err = ctx.Err(); err != nil {
			return written, err
		}
		sym = normalizeSymbol(sym)
		if sym == "" {
			continue
		}
		raw, marshalErr := json.Marshal(q)
		if marshalErr != nil {
			err = marshalErr
			return written, err
		}
		if _, err = quoteStmt.ExecContext(ctx, sym, string(raw), q.ProviderTimestamp, q.UpdatedAt, now, q.Source, q.DataState); err != nil {
			return written, err
		}
		stamp := q.ProviderTimestamp
		if stamp <= 0 {
			stamp = q.UpdatedAt
		}
		if stamp <= 0 {
			stamp = now
		}
		bucket := (stamp / 300000) * 300000
		if _, err = historyStmt.ExecContext(ctx, sym, bucket, q.ProviderTimestamp, q.UpdatedAt, q.Price, q.Bid, q.Ask, q.Volume, q.Source, q.DataState); err != nil {
			return written, err
		}
		written++
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM quote_history WHERE bucket_ms < $1`, now-2_592_000_000); err != nil {
		return written, err
	}
	if err = tx.Commit(); err != nil {
		return written, err
	}
	return written, nil
}

func (b *postgresPersistenceBackend) SaveIntelligence(ctx context.Context, batch PersistenceIntelligenceBatch) (written int, err error) {
	started := time.Now()
	defer func() { b.observe("save-intelligence", started, err) }()
	if batch.Len() == 0 {
		return 0, nil
	}
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	for _, r := range batch.Evidence {
		if r.ID == "" || r.Kind == "" {
			continue
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO evidence_records(evidence_id,symbol,evidence_kind,observed_at_ms,source,provenance,freshness_state,payload_json)
VALUES($1,$2,$3,$4,$5,$6,$7,$8::jsonb) ON CONFLICT(evidence_id) DO NOTHING`, r.ID, normalizeSymbol(r.Symbol), r.Kind, r.ObservedAt, r.Source, r.Provenance, r.FreshnessState, string(payloadOrEmpty(r.Payload))); err != nil {
			return written, err
		}
		written++
	}
	for _, r := range batch.Decisions {
		if r.ID == "" || r.DecisionKind == "" {
			continue
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO decision_lineage(decision_id,symbol,horizon,evidence_id,decision_kind,decision_value,formula_version,created_at_ms,payload_json)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9::jsonb) ON CONFLICT(decision_id) DO NOTHING`, r.ID, normalizeSymbol(r.Symbol), r.Horizon, r.EvidenceID, r.DecisionKind, r.DecisionValue, r.FormulaVersion, r.CreatedAt, string(payloadOrEmpty(r.Payload))); err != nil {
			return written, err
		}
		written++
	}
	for _, r := range batch.Outcomes {
		if r.ID == "" {
			continue
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO outcome_history(outcome_id,decision_id,symbol,horizon,observed_at_ms,outcome_label,payload_json)
VALUES($1,$2,$3,$4,$5,$6,$7::jsonb) ON CONFLICT(outcome_id) DO UPDATE SET decision_id=EXCLUDED.decision_id,symbol=EXCLUDED.symbol,horizon=EXCLUDED.horizon,observed_at_ms=EXCLUDED.observed_at_ms,outcome_label=EXCLUDED.outcome_label,payload_json=EXCLUDED.payload_json`, r.ID, r.DecisionID, normalizeSymbol(r.Symbol), r.Horizon, r.ObservedAt, r.OutcomeLabel, string(payloadOrEmpty(r.Payload))); err != nil {
			return written, err
		}
		written++
	}
	for _, r := range batch.Features {
		if normalizeSymbol(r.Symbol) == "" || r.FeatureKey == "" || r.FeatureVersion == "" {
			continue
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO derived_features(symbol,feature_key,feature_version,as_of_ms,source_hash,payload_json)
VALUES($1,$2,$3,$4,$5,$6::jsonb) ON CONFLICT(symbol,feature_key,feature_version) DO UPDATE SET as_of_ms=EXCLUDED.as_of_ms,source_hash=EXCLUDED.source_hash,payload_json=EXCLUDED.payload_json WHERE derived_features.source_hash<>EXCLUDED.source_hash`, normalizeSymbol(r.Symbol), r.FeatureKey, r.FeatureVersion, r.AsOf, r.SourceHash, string(payloadOrEmpty(r.Payload))); err != nil {
			return written, err
		}
		written++
	}
	if err = tx.Commit(); err != nil {
		return written, err
	}
	return written, nil
}

func (b *postgresPersistenceBackend) LoadQuotes(ctx context.Context) (out map[string]Quote, err error) {
	started := time.Now()
	defer func() { b.observe("load-quotes", started, err) }()
	rows, err := b.db.QueryContext(ctx, `SELECT symbol,payload_json FROM canonical_quotes ORDER BY persisted_at_ms DESC LIMIT 1000`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out = map[string]Quote{}
	for rows.Next() {
		var sym string
		var raw []byte
		if err = rows.Scan(&sym, &raw); err != nil {
			return nil, err
		}
		var q Quote
		if unmarshalErr := json.Unmarshal(raw, &q); unmarshalErr != nil {
			continue
		}
		q.Symbol = normalizeSymbol(sym)
		q.DataState = "persisted"
		q.FeedType = "persisted"
		out[q.Symbol] = q
	}
	return out, rows.Err()
}

func (b *postgresPersistenceBackend) LoadIdentityState(ctx context.Context) (state IdentityPersistentState, err error) {
	started := time.Now()
	defer func() { b.observe("load-identity", started, err) }()
	var raw []byte
	err = b.db.QueryRowContext(ctx, `SELECT payload_json FROM identity_state WHERE id=1`).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return IdentityPersistentState{}, nil
	}
	if err != nil {
		return IdentityPersistentState{}, err
	}
	if len(raw) == 0 {
		return IdentityPersistentState{}, nil
	}
	if err = json.Unmarshal(raw, &state); err != nil {
		return IdentityPersistentState{}, err
	}
	return state, nil
}

func (b *postgresPersistenceBackend) SaveIdentityState(ctx context.Context, state IdentityPersistentState) (err error) {
	started := time.Now()
	defer func() { b.observe("save-identity", started, err) }()
	raw, err := json.Marshal(state)
	if err != nil {
		return err
	}
	_, err = b.db.ExecContext(ctx, `INSERT INTO identity_state(id,payload_json,updated_at_ms) VALUES(1,$1::jsonb,$2)
ON CONFLICT(id) DO UPDATE SET payload_json=EXCLUDED.payload_json,updated_at_ms=EXCLUDED.updated_at_ms`, string(raw), state.UpdatedAt)
	return err
}

func (b *postgresPersistenceBackend) LoadUserWorkspaces(ctx context.Context) (out []UserWorkspace, err error) {
	started := time.Now()
	defer func() { b.observe("load-workspaces", started, err) }()
	rows, err := b.db.QueryContext(ctx, `SELECT payload_json FROM user_workspaces ORDER BY user_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var raw []byte
		if err = rows.Scan(&raw); err != nil {
			return nil, err
		}
		var workspace UserWorkspace
		if err = json.Unmarshal(raw, &workspace); err != nil {
			return nil, err
		}
		out = append(out, workspace)
	}
	return out, rows.Err()
}

func (b *postgresPersistenceBackend) SaveUserWorkspace(ctx context.Context, workspace UserWorkspace) (err error) {
	started := time.Now()
	defer func() { b.observe("save-workspace", started, err) }()
	if strings.TrimSpace(workspace.UserID) == "" {
		return errors.New("workspace user id is required")
	}
	raw, err := json.Marshal(workspace)
	if err != nil {
		return err
	}
	_, err = b.db.ExecContext(ctx, `INSERT INTO user_workspaces(user_id,payload_json,updated_at_ms) VALUES($1,$2::jsonb,$3)
ON CONFLICT(user_id) DO UPDATE SET payload_json=EXCLUDED.payload_json,updated_at_ms=EXCLUDED.updated_at_ms`, workspace.UserID, string(raw), workspace.UpdatedAt)
	return err
}

func (b *postgresPersistenceBackend) identityCounts(ctx context.Context) (users int, sessions int, err error) {
	state, err := b.LoadIdentityState(ctx)
	if err != nil {
		return 0, 0, err
	}
	return len(state.Users), len(state.Sessions), nil
}

func (b *postgresPersistenceBackend) scalarInt64(ctx context.Context, query string) (int64, error) {
	var out int64
	if err := b.db.QueryRowContext(ctx, query).Scan(&out); err != nil {
		return 0, err
	}
	return out, nil
}

func (b *postgresPersistenceBackend) Stats(ctx context.Context) (stats PersistenceStoreStats, err error) {
	started := time.Now()
	defer func() { b.observe("stats", started, err) }()
	queries := []string{
		`SELECT COALESCE(MAX(version),0) FROM schema_migrations`,
		`SELECT COUNT(*) FROM symbol_registry`,
		`SELECT COUNT(*) FROM symbol_registry WHERE active=TRUE`,
		`SELECT COUNT(*) FROM canonical_quotes`,
		`SELECT COUNT(*) FROM quote_history`,
		`SELECT COUNT(*) FROM evidence_records`,
		`SELECT COUNT(*) FROM decision_lineage`,
		`SELECT COUNT(*) FROM outcome_history`,
		`SELECT COUNT(*) FROM derived_features`,
	}
	vals := make([]int64, len(queries))
	for i, query := range queries {
		vals[i], err = b.scalarInt64(ctx, query)
		if err != nil {
			return PersistenceStoreStats{}, err
		}
	}
	users, sessions, err := b.identityCounts(ctx)
	if err != nil {
		return PersistenceStoreStats{}, err
	}
	var storageBytes int64
	if err = b.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(pg_total_relation_size(c.oid)),0)
FROM pg_class c JOIN pg_namespace n ON n.oid=c.relnamespace
WHERE n.nspname=current_schema() AND c.relkind='r' AND c.relname IN ('schema_migrations','symbol_registry','canonical_quotes','quote_history','evidence_records','decision_lineage','outcome_history','derived_features','identity_state','user_workspaces')`).Scan(&storageBytes); err != nil {
		// Storage-size telemetry is useful but must not make the canonical store unavailable.
		storageBytes = 0
		err = nil
	}
	return PersistenceStoreStats{
		SchemaVersion: int(vals[0]), SymbolCount: int(vals[1]), ActiveSymbolCount: int(vals[2]), CanonicalQuotes: int(vals[3]), QuoteHistoryRows: int(vals[4]), EvidenceRows: int(vals[5]), DecisionRows: int(vals[6]), OutcomeRows: int(vals[7]), FeatureRows: int(vals[8]), UserCount: users, SessionCount: sessions, StorageBytes: storageBytes,
	}, nil
}

func (b *postgresPersistenceBackend) Close() (err error) {
	started := time.Now()
	defer func() { b.observe("close", started, err) }()
	if b == nil || b.db == nil {
		return nil
	}
	err = b.db.Close()
	b.db = nil
	return err
}

func (b *postgresPersistenceBackend) HealthCheck(ctx context.Context) (err error) {
	started := time.Now()
	defer func() { b.observe("health-check", started, err) }()
	if b == nil || b.db == nil {
		return errors.New("postgres database is not open")
	}
	return b.db.PingContext(ctx)
}

func (b *postgresPersistenceBackend) ExportPersistenceArchive(ctx context.Context) (archive PersistenceArchive, err error) {
	started := time.Now()
	defer func() { b.observe("archive-export", started, err) }()
	if b == nil || b.db == nil {
		return PersistenceArchive{}, errors.New("postgres database is not open")
	}
	tx, err := b.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})
	if err != nil {
		return PersistenceArchive{}, err
	}
	defer tx.Rollback()
	if err := tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(version),0) FROM schema_migrations`).Scan(&archive.SourceStoreSchema); err != nil {
		return PersistenceArchive{}, err
	}
	rows, err := tx.QueryContext(ctx, `SELECT symbol,first_seen_ms,last_seen_ms,active,selected,processing_tier,desk_membership,provider_eligible,last_subscribed_ms,last_processed_ms FROM symbol_registry ORDER BY symbol`)
	if err != nil {
		return PersistenceArchive{}, err
	}
	for rows.Next() {
		var r SymbolRegistryRecord
		if err := rows.Scan(&r.Symbol, &r.FirstSeenAt, &r.LastSeenAt, &r.Active, &r.Selected, &r.ProcessingTier, &r.DeskMembership, &r.ProviderEligible, &r.LastSubscribedAt, &r.LastProcessedAt); err != nil {
			rows.Close()
			return PersistenceArchive{}, err
		}
		r.Symbol = normalizeSymbol(r.Symbol)
		archive.Symbols = append(archive.Symbols, r)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return PersistenceArchive{}, err
	}
	rows.Close()
	rows, err = tx.QueryContext(ctx, `SELECT symbol,payload_json,provider_timestamp_ms,received_timestamp_ms,persisted_at_ms,source,data_state FROM canonical_quotes ORDER BY symbol`)
	if err != nil {
		return PersistenceArchive{}, err
	}
	for rows.Next() {
		var r PersistenceCanonicalQuoteRecord
		var raw []byte
		if err := rows.Scan(&r.Symbol, &raw, &r.ProviderTimestamp, &r.ReceivedTimestamp, &r.PersistedAt, &r.Source, &r.DataState); err != nil {
			rows.Close()
			return PersistenceArchive{}, err
		}
		if err := json.Unmarshal(raw, &r.Quote); err != nil {
			rows.Close()
			return PersistenceArchive{}, err
		}
		r.Symbol = normalizeSymbol(r.Symbol)
		r.Quote.Symbol = r.Symbol
		archive.CanonicalQuotes = append(archive.CanonicalQuotes, r)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return PersistenceArchive{}, err
	}
	rows.Close()
	rows, err = tx.QueryContext(ctx, `SELECT symbol,bucket_ms,provider_timestamp_ms,received_timestamp_ms,price,bid,ask,volume,source,data_state FROM quote_history ORDER BY symbol,bucket_ms`)
	if err != nil {
		return PersistenceArchive{}, err
	}
	for rows.Next() {
		var r PersistenceQuoteHistoryRecord
		if err := rows.Scan(&r.Symbol, &r.Bucket, &r.ProviderTimestamp, &r.ReceivedTimestamp, &r.Price, &r.Bid, &r.Ask, &r.Volume, &r.Source, &r.DataState); err != nil {
			rows.Close()
			return PersistenceArchive{}, err
		}
		r.Symbol = normalizeSymbol(r.Symbol)
		archive.QuoteHistory = append(archive.QuoteHistory, r)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return PersistenceArchive{}, err
	}
	rows.Close()
	rows, err = tx.QueryContext(ctx, `SELECT evidence_id,symbol,evidence_kind,observed_at_ms,source,provenance,freshness_state,payload_json FROM evidence_records ORDER BY evidence_id`)
	if err != nil {
		return PersistenceArchive{}, err
	}
	for rows.Next() {
		var r EvidenceRecord
		var raw []byte
		if err := rows.Scan(&r.ID, &r.Symbol, &r.Kind, &r.ObservedAt, &r.Source, &r.Provenance, &r.FreshnessState, &raw); err != nil {
			rows.Close()
			return PersistenceArchive{}, err
		}
		r.Symbol = normalizeSymbol(r.Symbol)
		r.Payload = append(json.RawMessage(nil), raw...)
		archive.Evidence = append(archive.Evidence, r)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return PersistenceArchive{}, err
	}
	rows.Close()
	rows, err = tx.QueryContext(ctx, `SELECT decision_id,symbol,horizon,evidence_id,decision_kind,decision_value,formula_version,created_at_ms,payload_json FROM decision_lineage ORDER BY decision_id`)
	if err != nil {
		return PersistenceArchive{}, err
	}
	for rows.Next() {
		var r DecisionLineageRecord
		var raw []byte
		if err := rows.Scan(&r.ID, &r.Symbol, &r.Horizon, &r.EvidenceID, &r.DecisionKind, &r.DecisionValue, &r.FormulaVersion, &r.CreatedAt, &raw); err != nil {
			rows.Close()
			return PersistenceArchive{}, err
		}
		r.Symbol = normalizeSymbol(r.Symbol)
		r.Payload = append(json.RawMessage(nil), raw...)
		archive.Decisions = append(archive.Decisions, r)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return PersistenceArchive{}, err
	}
	rows.Close()
	rows, err = tx.QueryContext(ctx, `SELECT outcome_id,decision_id,symbol,horizon,observed_at_ms,outcome_label,payload_json FROM outcome_history ORDER BY outcome_id`)
	if err != nil {
		return PersistenceArchive{}, err
	}
	for rows.Next() {
		var r OutcomeHistoryRecord
		var raw []byte
		if err := rows.Scan(&r.ID, &r.DecisionID, &r.Symbol, &r.Horizon, &r.ObservedAt, &r.OutcomeLabel, &raw); err != nil {
			rows.Close()
			return PersistenceArchive{}, err
		}
		r.Symbol = normalizeSymbol(r.Symbol)
		r.Payload = append(json.RawMessage(nil), raw...)
		archive.Outcomes = append(archive.Outcomes, r)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return PersistenceArchive{}, err
	}
	rows.Close()
	rows, err = tx.QueryContext(ctx, `SELECT symbol,feature_key,feature_version,as_of_ms,source_hash,payload_json FROM derived_features ORDER BY symbol,feature_key,feature_version`)
	if err != nil {
		return PersistenceArchive{}, err
	}
	for rows.Next() {
		var r DerivedFeatureRecord
		var raw []byte
		if err := rows.Scan(&r.Symbol, &r.FeatureKey, &r.FeatureVersion, &r.AsOf, &r.SourceHash, &raw); err != nil {
			rows.Close()
			return PersistenceArchive{}, err
		}
		r.Symbol = normalizeSymbol(r.Symbol)
		r.Payload = append(json.RawMessage(nil), raw...)
		archive.Features = append(archive.Features, r)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return PersistenceArchive{}, err
	}
	rows.Close()
	if b.tenantScopedArchive {
		var present bool
		archive.Identity, present, err = loadHostedTenantIdentityArchive(ctx, tx)
		if err != nil {
			return PersistenceArchive{}, err
		}
		if !present {
			var legacyCount int64
			if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM identity_state`).Scan(&legacyCount); err != nil {
				return PersistenceArchive{}, err
			}
			if legacyCount != 0 {
				return PersistenceArchive{}, errors.New("hosted tenant archive found unretired legacy identity authority")
			}
		}
		archive.HasIdentity = present
		archive.UserWorkspaces, err = loadHostedTenantWorkspaceArchive(ctx, tx, archive.Identity)
		if err != nil {
			return PersistenceArchive{}, err
		}
	} else {
		var identityRaw []byte
		err = tx.QueryRowContext(ctx, `SELECT payload_json FROM identity_state WHERE id=1`).Scan(&identityRaw)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			err = nil
		case err != nil:
			return PersistenceArchive{}, err
		default:
			if err := json.Unmarshal(identityRaw, &archive.Identity); err != nil {
				return PersistenceArchive{}, err
			}
			archive.HasIdentity = true
		}
		rows, err = tx.QueryContext(ctx, `SELECT payload_json FROM user_workspaces ORDER BY user_id`)
		if err != nil {
			return PersistenceArchive{}, err
		}
		for rows.Next() {
			var raw []byte
			var w UserWorkspace
			if err := rows.Scan(&raw); err != nil {
				rows.Close()
				return PersistenceArchive{}, err
			}
			if err := json.Unmarshal(raw, &w); err != nil {
				rows.Close()
				return PersistenceArchive{}, err
			}
			archive.UserWorkspaces = append(archive.UserWorkspaces, w)
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return PersistenceArchive{}, err
		}
		rows.Close()
	}
	if err := tx.Commit(); err != nil {
		return PersistenceArchive{}, err
	}
	return archive, nil
}

func (b *postgresPersistenceBackend) RestorePersistenceArchive(ctx context.Context, archive PersistenceArchive, mode string) (err error) {
	started := time.Now()
	defer func() { b.observe("archive-restore", started, err) }()
	if b == nil || b.db == nil {
		return errors.New("postgres database is not open")
	}
	if archive.SchemaVersion != persistenceArchiveSchemaVersion {
		return errors.New("unsupported persistence archive schema")
	}
	var tenantPartitions map[string]IdentityPersistentState
	workspaceOwners := map[string]string{}
	if b.tenantScopedArchive {
		if archive.HasIdentity {
			var validationErr error
			tenantPartitions, validationErr = hostedTenantIdentityPartitions(archive.Identity, false)
			if validationErr != nil {
				return fmt.Errorf("hosted tenant archive identity invalid: %w", validationErr)
			}
		} else if len(archive.UserWorkspaces) != 0 {
			return errors.New("hosted tenant archive has workspaces without canonical identity")
		}
		userIndex, validationErr := hostedTenantUserIndex(archive.Identity)
		if validationErr != nil {
			return fmt.Errorf("hosted tenant archive ownership invalid: %w", validationErr)
		}
		for _, workspace := range archive.UserWorkspaces {
			userID := strings.TrimSpace(workspace.UserID)
			tenantID, ok := userIndex[userID]
			if userID == "" || !ok {
				return fmt.Errorf("hosted tenant archive workspace %q has no canonical tenant owner", userID)
			}
			if _, duplicate := workspaceOwners[userID]; duplicate {
				return fmt.Errorf("hosted tenant archive workspace %q is duplicated", userID)
			}
			workspaceOwners[userID] = tenantID
		}
	}
	tx, err := b.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var count int64
	countQuery := `SELECT (SELECT COUNT(*) FROM symbol_registry)+(SELECT COUNT(*) FROM canonical_quotes)+(SELECT COUNT(*) FROM quote_history)+(SELECT COUNT(*) FROM evidence_records)+(SELECT COUNT(*) FROM decision_lineage)+(SELECT COUNT(*) FROM outcome_history)+(SELECT COUNT(*) FROM derived_features)+(SELECT COUNT(*) FROM identity_state)+(SELECT COUNT(*) FROM user_workspaces)`
	if b.tenantScopedArchive {
		countQuery += `+(SELECT COUNT(*) FROM tenant_identity_state)+(SELECT COUNT(*) FROM tenant_user_workspaces)`
	}
	if err := tx.QueryRowContext(ctx, countQuery).Scan(&count); err != nil {
		return err
	}
	if mode == persistenceRestoreModeEmpty && count > 0 {
		return errors.New("persistence restore target is not empty; use explicit replace mode")
	}
	if mode == persistenceRestoreModeReplace {
		tables := `user_workspaces, identity_state, derived_features, outcome_history, decision_lineage, evidence_records, quote_history, canonical_quotes, symbol_registry`
		if b.tenantScopedArchive {
			tables = `tenant_user_workspaces, tenant_identity_state, ` + tables
		}
		if _, err := tx.ExecContext(ctx, `TRUNCATE TABLE `+tables); err != nil {
			return err
		}
	}
	for _, r := range archive.Symbols {
		if _, err := tx.ExecContext(ctx, `INSERT INTO symbol_registry(symbol,first_seen_ms,last_seen_ms,active,selected,processing_tier,desk_membership,provider_eligible,last_subscribed_ms,last_processed_ms) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, normalizeSymbol(r.Symbol), r.FirstSeenAt, r.LastSeenAt, r.Active, r.Selected, r.ProcessingTier, r.DeskMembership, r.ProviderEligible, r.LastSubscribedAt, r.LastProcessedAt); err != nil {
			return err
		}
	}
	for _, r := range archive.CanonicalQuotes {
		raw, e := json.Marshal(r.Quote)
		if e != nil {
			return e
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO canonical_quotes(symbol,payload_json,provider_timestamp_ms,received_timestamp_ms,persisted_at_ms,source,data_state) VALUES($1,$2,$3,$4,$5,$6,$7)`, normalizeSymbol(r.Symbol), raw, r.ProviderTimestamp, r.ReceivedTimestamp, r.PersistedAt, r.Source, r.DataState); err != nil {
			return err
		}
	}
	for _, r := range archive.QuoteHistory {
		if _, err := tx.ExecContext(ctx, `INSERT INTO quote_history(symbol,bucket_ms,provider_timestamp_ms,received_timestamp_ms,price,bid,ask,volume,source,data_state) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, normalizeSymbol(r.Symbol), r.Bucket, r.ProviderTimestamp, r.ReceivedTimestamp, r.Price, r.Bid, r.Ask, r.Volume, r.Source, r.DataState); err != nil {
			return err
		}
	}
	for _, r := range archive.Evidence {
		if _, err := tx.ExecContext(ctx, `INSERT INTO evidence_records(evidence_id,symbol,evidence_kind,observed_at_ms,source,provenance,freshness_state,payload_json) VALUES($1,$2,$3,$4,$5,$6,$7,$8)`, r.ID, normalizeSymbol(r.Symbol), r.Kind, r.ObservedAt, r.Source, r.Provenance, r.FreshnessState, payloadOrEmpty(r.Payload)); err != nil {
			return err
		}
	}
	for _, r := range archive.Decisions {
		if _, err := tx.ExecContext(ctx, `INSERT INTO decision_lineage(decision_id,symbol,horizon,evidence_id,decision_kind,decision_value,formula_version,created_at_ms,payload_json) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`, r.ID, normalizeSymbol(r.Symbol), r.Horizon, r.EvidenceID, r.DecisionKind, r.DecisionValue, r.FormulaVersion, r.CreatedAt, payloadOrEmpty(r.Payload)); err != nil {
			return err
		}
	}
	for _, r := range archive.Outcomes {
		if _, err := tx.ExecContext(ctx, `INSERT INTO outcome_history(outcome_id,decision_id,symbol,horizon,observed_at_ms,outcome_label,payload_json) VALUES($1,$2,$3,$4,$5,$6,$7)`, r.ID, r.DecisionID, normalizeSymbol(r.Symbol), r.Horizon, r.ObservedAt, r.OutcomeLabel, payloadOrEmpty(r.Payload)); err != nil {
			return err
		}
	}
	for _, r := range archive.Features {
		if _, err := tx.ExecContext(ctx, `INSERT INTO derived_features(symbol,feature_key,feature_version,as_of_ms,source_hash,payload_json) VALUES($1,$2,$3,$4,$5,$6)`, normalizeSymbol(r.Symbol), r.FeatureKey, r.FeatureVersion, r.AsOf, r.SourceHash, payloadOrEmpty(r.Payload)); err != nil {
			return err
		}
	}
	if archive.HasIdentity {
		if b.tenantScopedArchive {
			keys := make([]string, 0, len(tenantPartitions))
			for tenantID := range tenantPartitions {
				keys = append(keys, tenantID)
			}
			sort.Strings(keys)
			for _, tenantID := range keys {
				part := tenantPartitions[tenantID]
				raw, e := json.Marshal(part)
				if e != nil {
					return e
				}
				if _, err := tx.ExecContext(ctx, `INSERT INTO tenant_identity_state(tenant_id,payload_json,updated_at_ms) VALUES($1,$2,$3)`, tenantID, raw, part.UpdatedAt); err != nil {
					return err
				}
			}
		} else {
			raw, e := json.Marshal(archive.Identity)
			if e != nil {
				return e
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO identity_state(id,payload_json,updated_at_ms) VALUES(1,$1,$2)`, raw, archive.Identity.UpdatedAt); err != nil {
				return err
			}
		}
	}
	for _, w := range archive.UserWorkspaces {
		raw, e := json.Marshal(w)
		if e != nil {
			return e
		}
		if b.tenantScopedArchive {
			if _, err := tx.ExecContext(ctx, `INSERT INTO tenant_user_workspaces(tenant_id,user_id,payload_json,updated_at_ms) VALUES($1,$2,$3,$4)`, workspaceOwners[strings.TrimSpace(w.UserID)], strings.TrimSpace(w.UserID), raw, w.UpdatedAt); err != nil {
				return err
			}
		} else {
			if _, err := tx.ExecContext(ctx, `INSERT INTO user_workspaces(user_id,payload_json,updated_at_ms) VALUES($1,$2,$3)`, w.UserID, raw, w.UpdatedAt); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}
