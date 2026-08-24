//go:build postgres

package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func (b *postgresPersistenceBackend) ensureInstrumentIdentitySchema(ctx context.Context) error {
	if b == nil || b.db == nil {
		return errors.New("postgres database is not open")
	}
	conn, err := b.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("postgres instrument identity migration connection: %w", err)
	}
	defer conn.Close()
	if _, err := conn.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations(version INTEGER PRIMARY KEY, applied_at_ms BIGINT NOT NULL)`); err != nil {
		return fmt.Errorf("postgres instrument identity migration table: %w", err)
	}
	if _, err := conn.ExecContext(ctx, `SELECT pg_advisory_lock($1)`, postgresSchemaLockKey); err != nil {
		return fmt.Errorf("postgres instrument identity migration lock: %w", err)
	}
	defer func() {
		unlockCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, _ = conn.ExecContext(unlockCtx, `SELECT pg_advisory_unlock($1)`, postgresSchemaLockKey)
	}()

	var applied bool
	if err := conn.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)`, instrumentIdentitySchemaVersion).Scan(&applied); err != nil {
		return fmt.Errorf("postgres instrument identity migration check: %w", err)
	}
	if applied {
		return nil
	}
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("postgres instrument identity migration begin: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS instrument_identities(
 symbol TEXT PRIMARY KEY,
 display_name TEXT NOT NULL DEFAULT '',
 exchange TEXT NOT NULL DEFAULT '',
 asset_class TEXT NOT NULL DEFAULT '',
 provider_asset_id TEXT NOT NULL DEFAULT '',
 source TEXT NOT NULL DEFAULT '',
 observed_at_ms BIGINT NOT NULL
)`); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("postgres instrument identity migration: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_instrument_identities_observed ON instrument_identities(observed_at_ms DESC)`); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("postgres instrument identity index: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations(version,applied_at_ms) VALUES($1,$2)`, instrumentIdentitySchemaVersion, time.Now().UnixMilli()); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("postgres instrument identity migration record: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("postgres instrument identity migration commit: %w", err)
	}
	return nil
}

func (b *postgresPersistenceBackend) SaveInstrumentIdentities(ctx context.Context, records []InstrumentIdentityRecord) (written int, err error) {
	started := time.Now()
	defer func() { b.observe("save-instrument-identities", started, err) }()
	if err = b.ensureInstrumentIdentitySchema(ctx); err != nil {
		return 0, err
	}
	records = canonicalInstrumentIdentities(records)
	if len(records) == 0 {
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
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO instrument_identities(symbol,display_name,exchange,asset_class,provider_asset_id,source,observed_at_ms)
VALUES($1,$2,$3,$4,$5,$6,$7)
ON CONFLICT(symbol) DO UPDATE SET display_name=EXCLUDED.display_name,exchange=EXCLUDED.exchange,asset_class=EXCLUDED.asset_class,provider_asset_id=EXCLUDED.provider_asset_id,source=EXCLUDED.source,observed_at_ms=EXCLUDED.observed_at_ms
WHERE EXCLUDED.observed_at_ms >= instrument_identities.observed_at_ms`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	for _, record := range records {
		if err = ctx.Err(); err != nil {
			return written, err
		}
		if _, err = stmt.ExecContext(ctx, record.Symbol, record.Name, record.Exchange, record.AssetClass, record.ProviderAssetID, record.Source, record.ObservedAt); err != nil {
			return written, err
		}
		written++
	}
	if err = tx.Commit(); err != nil {
		return written, err
	}
	return written, nil
}

func (b *postgresPersistenceBackend) LoadInstrumentIdentities(ctx context.Context) (out []InstrumentIdentityRecord, err error) {
	started := time.Now()
	defer func() { b.observe("load-instrument-identities", started, err) }()
	if err = b.ensureInstrumentIdentitySchema(ctx); err != nil {
		return nil, err
	}
	rows, err := b.db.QueryContext(ctx, `SELECT symbol,display_name,exchange,asset_class,provider_asset_id,source,observed_at_ms FROM instrument_identities ORDER BY symbol`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record InstrumentIdentityRecord
		if err = rows.Scan(&record.Symbol, &record.Name, &record.Exchange, &record.AssetClass, &record.ProviderAssetID, &record.Source, &record.ObservedAt); err != nil {
			return nil, err
		}
		if normalized, ok := normalizeInstrumentIdentity(record); ok {
			out = append(out, normalized)
		}
	}
	return canonicalInstrumentIdentities(out), rows.Err()
}
