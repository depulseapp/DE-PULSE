//go:build cgo && !windows

package main

/*
#include <sqlite3.h>
*/
import "C"

import (
	"context"
	"fmt"
	"unsafe"
)

func (b *sqlitePersistenceBackend) ensureInstrumentIdentitySchemaLocked() error {
	applied, err := b.migrationApplied(instrumentIdentitySchemaVersion)
	if err != nil {
		return err
	}
	if applied {
		return nil
	}
	if err := b.exec("BEGIN IMMEDIATE"); err != nil {
		return err
	}
	ok := false
	defer func() {
		if !ok {
			_ = b.exec("ROLLBACK")
		}
	}()
	if err := b.exec(`CREATE TABLE IF NOT EXISTS instrument_identities(
 symbol TEXT PRIMARY KEY,
 display_name TEXT NOT NULL DEFAULT '',
 exchange TEXT NOT NULL DEFAULT '',
 asset_class TEXT NOT NULL DEFAULT '',
 provider_asset_id TEXT NOT NULL DEFAULT '',
 source TEXT NOT NULL DEFAULT '',
 observed_at_ms INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_instrument_identities_observed ON instrument_identities(observed_at_ms DESC);`); err != nil {
		return err
	}
	if err := b.exec(fmt.Sprintf("INSERT INTO schema_migrations(version,applied_at_ms) VALUES(%d, strftime('%%s','now')*1000)", instrumentIdentitySchemaVersion)); err != nil {
		return err
	}
	if err := b.exec("COMMIT"); err != nil {
		return err
	}
	ok = true
	return nil
}

func (b *sqlitePersistenceBackend) SaveInstrumentIdentities(ctx context.Context, records []InstrumentIdentityRecord) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if err := b.ensureInstrumentIdentitySchemaLocked(); err != nil {
		return 0, err
	}
	if len(records) == 0 {
		return 0, nil
	}
	if err := b.exec("BEGIN IMMEDIATE"); err != nil {
		return 0, err
	}
	ok := false
	defer func() {
		if !ok {
			_ = b.exec("ROLLBACK")
		}
	}()
	stmt, err := prepare(b.db, `INSERT INTO instrument_identities(symbol,display_name,exchange,asset_class,provider_asset_id,source,observed_at_ms)
VALUES(?,?,?,?,?,?,?)
ON CONFLICT(symbol) DO UPDATE SET display_name=excluded.display_name,exchange=excluded.exchange,asset_class=excluded.asset_class,provider_asset_id=excluded.provider_asset_id,source=excluded.source,observed_at_ms=excluded.observed_at_ms
WHERE excluded.observed_at_ms >= instrument_identities.observed_at_ms`)
	if err != nil {
		return 0, err
	}
	defer C.sqlite3_finalize(stmt)
	written := 0
	for _, record := range canonicalInstrumentIdentities(records) {
		if err := ctx.Err(); err != nil {
			return written, err
		}
		C.sqlite3_reset(stmt)
		C.sqlite3_clear_bindings(stmt)
		bindText(stmt, 1, record.Symbol)
		bindText(stmt, 2, record.Name)
		bindText(stmt, 3, record.Exchange)
		bindText(stmt, 4, record.AssetClass)
		bindText(stmt, 5, record.ProviderAssetID)
		bindText(stmt, 6, record.Source)
		C.sqlite3_bind_int64(stmt, 7, C.sqlite3_int64(record.ObservedAt))
		if rc := C.sqlite3_step(stmt); rc != C.SQLITE_DONE {
			return written, b.sqliteErr()
		}
		written++
	}
	if err := b.exec("COMMIT"); err != nil {
		return written, err
	}
	ok = true
	return written, nil
}

func sqliteIdentityColumnText(stmt *C.sqlite3_stmt, col int) string {
	ptr := C.sqlite3_column_text(stmt, C.int(col))
	if ptr == nil {
		return ""
	}
	return C.GoString((*C.char)(unsafe.Pointer(ptr)))
}

func (b *sqlitePersistenceBackend) LoadInstrumentIdentities(ctx context.Context) ([]InstrumentIdentityRecord, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := b.ensureInstrumentIdentitySchemaLocked(); err != nil {
		return nil, err
	}
	stmt, err := prepare(b.db, `SELECT symbol,display_name,exchange,asset_class,provider_asset_id,source,observed_at_ms FROM instrument_identities ORDER BY symbol`)
	if err != nil {
		return nil, err
	}
	defer C.sqlite3_finalize(stmt)
	out := []InstrumentIdentityRecord{}
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		switch rc := C.sqlite3_step(stmt); rc {
		case C.SQLITE_DONE:
			return out, nil
		case C.SQLITE_ROW:
			out = append(out, InstrumentIdentityRecord{
				Symbol:          normalizeSymbol(sqliteIdentityColumnText(stmt, 0)),
				Name:            sqliteIdentityColumnText(stmt, 1),
				Exchange:        sqliteIdentityColumnText(stmt, 2),
				AssetClass:      sqliteIdentityColumnText(stmt, 3),
				ProviderAssetID: sqliteIdentityColumnText(stmt, 4),
				Source:          sqliteIdentityColumnText(stmt, 5),
				ObservedAt:      int64(C.sqlite3_column_int64(stmt, 6)),
			})
		default:
			return nil, b.sqliteErr()
		}
	}
}
