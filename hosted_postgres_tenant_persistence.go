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
	"time"
)

const hostedTenantPostgresSchemaVersion = 5

type hostedTenantPostgresBackend struct {
	PersistenceBackend
	pg *postgresPersistenceBackend
}

func wrapHostedTenantPostgresBackend(inner PersistenceBackend) PersistenceBackend {
	if inner == nil || !isHostedRuntime() {
		return inner
	}
	pg, ok := inner.(*postgresPersistenceBackend)
	if !ok {
		return newUnavailablePersistenceBackend("hosted PostgreSQL tenancy requires the canonical PostgreSQL backend")
	}
	return &hostedTenantPostgresBackend{PersistenceBackend: inner, pg: pg}
}

func (b *hostedTenantPostgresBackend) Capabilities() []string {
	if b == nil || b.PersistenceBackend == nil {
		return nil
	}
	out := append([]string(nil), b.PersistenceBackend.Capabilities()...)
	return append(out, "tenant-scoped-identity", "tenant-scoped-user-workspaces", "expand-contract-tenancy")
}

func (b *hostedTenantPostgresBackend) Init(ctx context.Context) error {
	if b == nil || b.PersistenceBackend == nil || b.pg == nil {
		return errors.New("hosted PostgreSQL tenant persistence unavailable")
	}
	if err := b.PersistenceBackend.Init(ctx); err != nil {
		return err
	}
	migrationCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := b.applyHostedTenantMigration(migrationCtx); err != nil {
		_ = b.PersistenceBackend.Close()
		return err
	}
	if err := b.migrateLegacyHostedTenantState(migrationCtx); err != nil {
		_ = b.PersistenceBackend.Close()
		return err
	}
	return nil
}

func (b *hostedTenantPostgresBackend) applyHostedTenantMigration(ctx context.Context) (err error) {
	started := time.Now()
	defer func() { b.pg.observe("hosted-tenant-migration", started, err) }()
	if b.pg.db == nil {
		return errors.New("postgres database is not open")
	}
	conn, err := b.pg.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("hosted tenant migration connection: %w", err)
	}
	defer conn.Close()
	if _, err := conn.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations(version INTEGER PRIMARY KEY, applied_at_ms BIGINT NOT NULL)`); err != nil {
		return fmt.Errorf("hosted tenant migration table: %w", err)
	}
	if _, err := conn.ExecContext(ctx, `SELECT pg_advisory_lock($1)`, postgresSchemaLockKey); err != nil {
		return fmt.Errorf("hosted tenant migration lock: %w", err)
	}
	defer func() {
		unlockCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, _ = conn.ExecContext(unlockCtx, `SELECT pg_advisory_unlock($1)`, postgresSchemaLockKey)
	}()

	var applied bool
	if err := conn.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)`, hostedTenantPostgresSchemaVersion).Scan(&applied); err != nil {
		return fmt.Errorf("hosted tenant migration check: %w", err)
	}
	if applied {
		return nil
	}
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("hosted tenant migration begin: %w", err)
	}
	defer tx.Rollback()
	statements := []string{
		`CREATE TABLE IF NOT EXISTS tenant_identity_state(
 tenant_id TEXT PRIMARY KEY CHECK(length(btrim(tenant_id)) > 0),
 payload_json JSONB NOT NULL,
 updated_at_ms BIGINT NOT NULL
)`,
		`CREATE INDEX IF NOT EXISTS idx_tenant_identity_state_updated ON tenant_identity_state(updated_at_ms DESC)`,
		`CREATE TABLE IF NOT EXISTS tenant_user_workspaces(
 tenant_id TEXT NOT NULL CHECK(length(btrim(tenant_id)) > 0),
 user_id TEXT NOT NULL CHECK(length(btrim(user_id)) > 0),
 payload_json JSONB NOT NULL,
 updated_at_ms BIGINT NOT NULL,
 PRIMARY KEY(tenant_id,user_id)
)`,
		`CREATE INDEX IF NOT EXISTS idx_tenant_user_workspaces_user ON tenant_user_workspaces(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tenant_user_workspaces_updated ON tenant_user_workspaces(updated_at_ms DESC)`,
	}
	for _, statement := range statements {
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("hosted tenant migration v%d: %w", hostedTenantPostgresSchemaVersion, err)
		}
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations(version,applied_at_ms) VALUES($1,$2)`, hostedTenantPostgresSchemaVersion, time.Now().UnixMilli()); err != nil {
		return fmt.Errorf("hosted tenant migration record: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("hosted tenant migration commit: %w", err)
	}
	return nil
}

func hostedTenantStateHasData(state IdentityPersistentState) bool {
	return state.Version != 0 || state.UpdatedAt != 0 || len(state.Tenants) > 0 || len(state.Users) > 0 || len(state.Devices) > 0 || len(state.Sessions) > 0 || len(state.SecurityEvents) > 0 || len(state.ProductEntitlements) > 0
}

func (b *hostedTenantPostgresBackend) migrateLegacyHostedTenantState(ctx context.Context) error {
	if b.pg == nil || b.pg.db == nil {
		return errors.New("postgres database is not open")
	}
	var tenantRows int
	if err := b.pg.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tenant_identity_state`).Scan(&tenantRows); err != nil {
		return fmt.Errorf("count tenant identity rows: %w", err)
	}
	if tenantRows == 0 {
		legacy, err := b.pg.LoadIdentityState(ctx)
		if err != nil {
			return fmt.Errorf("load legacy identity state: %w", err)
		}
		if hostedTenantStateHasData(legacy) {
			if err := b.writeHostedTenantIdentityState(ctx, legacy, true); err != nil {
				return fmt.Errorf("migrate legacy identity state: %w", err)
			}
		}
	} else {
		if _, err := b.LoadIdentityState(ctx); err != nil {
			return fmt.Errorf("validate tenant identity state: %w", err)
		}
		if _, err := b.pg.db.ExecContext(ctx, `DELETE FROM identity_state`); err != nil {
			return fmt.Errorf("retire legacy identity aggregate: %w", err)
		}
	}

	legacyWorkspaces, err := b.pg.LoadUserWorkspaces(ctx)
	if err != nil {
		return fmt.Errorf("load legacy workspaces: %w", err)
	}
	if len(legacyWorkspaces) == 0 {
		return nil
	}
	state, err := b.LoadIdentityState(ctx)
	if err != nil {
		return err
	}
	userIndex, err := hostedTenantUserIndex(state)
	if err != nil {
		return err
	}
	for _, workspace := range legacyWorkspaces {
		userID := strings.TrimSpace(workspace.UserID)
		tenantID, ok := userIndex[userID]
		if !ok {
			return fmt.Errorf("legacy workspace %q has no canonical tenant owner", userID)
		}
		exists, err := b.tenantWorkspaceExists(ctx, tenantID, userID)
		if err != nil {
			return err
		}
		if !exists {
			if err := b.SaveUserWorkspace(ctx, workspace); err != nil {
				return fmt.Errorf("migrate legacy workspace %q: %w", userID, err)
			}
			continue
		}
		if _, err := b.pg.db.ExecContext(ctx, `DELETE FROM user_workspaces WHERE user_id=$1`, userID); err != nil {
			return fmt.Errorf("retire migrated legacy workspace %q: %w", userID, err)
		}
	}
	return nil
}

func (b *hostedTenantPostgresBackend) tenantWorkspaceExists(ctx context.Context, tenantID, userID string) (bool, error) {
	rows, err := b.pg.db.QueryContext(ctx, `SELECT tenant_id FROM tenant_user_workspaces WHERE user_id=$1`, userID)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	matches := 0
	for rows.Next() {
		var storedTenant string
		if err := rows.Scan(&storedTenant); err != nil {
			return false, err
		}
		matches++
		if strings.TrimSpace(storedTenant) != tenantID {
			return false, fmt.Errorf("workspace %q is persisted under conflicting tenant %q", userID, storedTenant)
		}
	}
	if err := rows.Err(); err != nil {
		return false, err
	}
	if matches > 1 {
		return false, fmt.Errorf("workspace %q has ambiguous tenant persistence", userID)
	}
	return matches == 1, nil
}

func (b *hostedTenantPostgresBackend) LoadIdentityState(ctx context.Context) (state IdentityPersistentState, err error) {
	started := time.Now()
	defer func() { b.pg.observe("load-tenant-identity", started, err) }()
	if b == nil || b.pg == nil || b.pg.db == nil {
		return IdentityPersistentState{}, errors.New("postgres database is not open")
	}
	rows, err := b.pg.db.QueryContext(ctx, `SELECT tenant_id,payload_json FROM tenant_identity_state ORDER BY tenant_id`)
	if err != nil {
		return IdentityPersistentState{}, err
	}
	defer rows.Close()
	partitions := map[string]IdentityPersistentState{}
	for rows.Next() {
		var tenantID string
		var raw []byte
		if err = rows.Scan(&tenantID, &raw); err != nil {
			return IdentityPersistentState{}, err
		}
		tenantID = strings.TrimSpace(tenantID)
		if tenantID == "" {
			return IdentityPersistentState{}, errors.New("tenant identity row has empty tenant id")
		}
		if _, exists := partitions[tenantID]; exists {
			return IdentityPersistentState{}, fmt.Errorf("duplicate tenant identity row %q", tenantID)
		}
		var part IdentityPersistentState
		if err = json.Unmarshal(raw, &part); err != nil {
			return IdentityPersistentState{}, fmt.Errorf("decode tenant identity %q: %w", tenantID, err)
		}
		partitions[tenantID] = part
	}
	if err = rows.Err(); err != nil {
		return IdentityPersistentState{}, err
	}
	if len(partitions) == 0 {
		return b.pg.LoadIdentityState(ctx)
	}
	return hostedTenantIdentityFromPartitions(partitions)
}

func (b *hostedTenantPostgresBackend) SaveIdentityState(ctx context.Context, state IdentityPersistentState) (err error) {
	started := time.Now()
	defer func() { b.pg.observe("save-tenant-identity", started, err) }()
	return b.writeHostedTenantIdentityState(ctx, state, false)
}

func (b *hostedTenantPostgresBackend) writeHostedTenantIdentityState(ctx context.Context, state IdentityPersistentState, allowLegacy bool) error {
	if b == nil || b.pg == nil || b.pg.db == nil {
		return errors.New("postgres database is not open")
	}
	partitions, err := hostedTenantIdentityPartitions(state, allowLegacy)
	if err != nil {
		return err
	}
	keys := make([]string, 0, len(partitions))
	for tenantID := range partitions {
		keys = append(keys, tenantID)
	}
	sort.Strings(keys)
	tx, err := b.pg.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM tenant_identity_state`); err != nil {
		return err
	}
	for _, tenantID := range keys {
		part := partitions[tenantID]
		raw, err := json.Marshal(part)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO tenant_identity_state(tenant_id,payload_json,updated_at_ms) VALUES($1,$2::jsonb,$3)`, tenantID, string(raw), part.UpdatedAt); err != nil {
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM identity_state`); err != nil {
		return err
	}
	return tx.Commit()
}

func (b *hostedTenantPostgresBackend) LoadUserWorkspaces(ctx context.Context) (out []UserWorkspace, err error) {
	started := time.Now()
	defer func() { b.pg.observe("load-tenant-workspaces", started, err) }()
	if b == nil || b.pg == nil || b.pg.db == nil {
		return nil, errors.New("postgres database is not open")
	}
	state, err := b.LoadIdentityState(ctx)
	if err != nil {
		return nil, err
	}
	userIndex, err := hostedTenantUserIndex(state)
	if err != nil {
		return nil, err
	}
	rows, err := b.pg.db.QueryContext(ctx, `SELECT tenant_id,user_id,payload_json FROM tenant_user_workspaces ORDER BY tenant_id,user_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	seen := map[string]struct{}{}
	for rows.Next() {
		var tenantID, userID string
		var raw []byte
		if err = rows.Scan(&tenantID, &userID, &raw); err != nil {
			return nil, err
		}
		tenantID = strings.TrimSpace(tenantID)
		userID = strings.TrimSpace(userID)
		expectedTenant, ok := userIndex[userID]
		if !ok {
			return nil, fmt.Errorf("workspace %q has no canonical identity owner", userID)
		}
		if expectedTenant != tenantID {
			return nil, fmt.Errorf("workspace %q crosses tenant boundary %q -> %q", userID, tenantID, expectedTenant)
		}
		if _, exists := seen[userID]; exists {
			return nil, fmt.Errorf("workspace %q has duplicate tenant rows", userID)
		}
		seen[userID] = struct{}{}
		var workspace UserWorkspace
		if err = json.Unmarshal(raw, &workspace); err != nil {
			return nil, err
		}
		if strings.TrimSpace(workspace.UserID) != userID {
			return nil, fmt.Errorf("workspace row %q payload user mismatch %q", userID, workspace.UserID)
		}
		out = append(out, workspace)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool { return strings.TrimSpace(out[i].UserID) < strings.TrimSpace(out[j].UserID) })
	return out, nil
}

func (b *hostedTenantPostgresBackend) SaveUserWorkspace(ctx context.Context, workspace UserWorkspace) (err error) {
	started := time.Now()
	defer func() { b.pg.observe("save-tenant-workspace", started, err) }()
	if b == nil || b.pg == nil || b.pg.db == nil {
		return errors.New("postgres database is not open")
	}
	workspace.UserID = strings.TrimSpace(workspace.UserID)
	if workspace.UserID == "" {
		return errors.New("workspace user id is required")
	}
	state, err := b.LoadIdentityState(ctx)
	if err != nil {
		return err
	}
	userIndex, err := hostedTenantUserIndex(state)
	if err != nil {
		return err
	}
	tenantID, ok := userIndex[workspace.UserID]
	if !ok {
		return fmt.Errorf("workspace user %q has no canonical tenant owner", workspace.UserID)
	}
	raw, err := json.Marshal(workspace)
	if err != nil {
		return err
	}
	tx, err := b.pg.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM tenant_user_workspaces WHERE user_id=$1 AND tenant_id<>$2`, workspace.UserID, tenantID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO tenant_user_workspaces(tenant_id,user_id,payload_json,updated_at_ms) VALUES($1,$2,$3::jsonb,$4)
ON CONFLICT(tenant_id,user_id) DO UPDATE SET payload_json=EXCLUDED.payload_json,updated_at_ms=EXCLUDED.updated_at_ms`, tenantID, workspace.UserID, string(raw), workspace.UpdatedAt); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM user_workspaces WHERE user_id=$1`, workspace.UserID); err != nil {
		return err
	}
	return tx.Commit()
}

func (b *hostedTenantPostgresBackend) Stats(ctx context.Context) (PersistenceStoreStats, error) {
	if b == nil || b.pg == nil {
		return PersistenceStoreStats{}, errors.New("postgres database is not open")
	}
	stats, err := b.pg.Stats(ctx)
	if err != nil {
		return PersistenceStoreStats{}, err
	}
	state, err := b.LoadIdentityState(ctx)
	if err != nil {
		return PersistenceStoreStats{}, err
	}
	stats.UserCount = len(state.Users)
	stats.SessionCount = len(state.Sessions)
	var tenantBytes int64
	if err := b.pg.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(pg_total_relation_size(c.oid)),0)
FROM pg_class c JOIN pg_namespace n ON n.oid=c.relnamespace
WHERE n.nspname=current_schema() AND c.relkind='r' AND c.relname IN ('tenant_identity_state','tenant_user_workspaces')`).Scan(&tenantBytes); err == nil {
		stats.StorageBytes += tenantBytes
	}
	return stats, nil
}

func (b *hostedTenantPostgresBackend) HealthCheck(ctx context.Context) error {
	if b == nil || b.pg == nil {
		return errors.New("postgres database is not open")
	}
	return b.pg.HealthCheck(ctx)
}

func (b *hostedTenantPostgresBackend) PoolDiagnostics() PersistencePoolDiagnostics {
	if b == nil || b.pg == nil {
		return PersistencePoolDiagnostics{}
	}
	return b.pg.PoolDiagnostics()
}

func (b *hostedTenantPostgresBackend) DatabaseDiagnostics() PersistenceDatabaseDiagnostics {
	if b == nil || b.pg == nil {
		return PersistenceDatabaseDiagnostics{}
	}
	return b.pg.DatabaseDiagnostics()
}
