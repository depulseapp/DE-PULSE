//go:build windows

package main

import (
	"context"
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

var (
	winSQLiteIdentityBindText      = winSQLiteDLL.NewProc("sqlite3_bind_text")
	winSQLiteIdentityBindInt64     = winSQLiteDLL.NewProc("sqlite3_bind_int64")
	winSQLiteIdentityReset         = winSQLiteDLL.NewProc("sqlite3_reset")
	winSQLiteIdentityClearBindings = winSQLiteDLL.NewProc("sqlite3_clear_bindings")
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

func bindWinIdentityText(stmt uintptr, idx int, value string) error {
	buf := append([]byte(value), 0)
	rc, _, _ := winSQLiteIdentityBindText.Call(stmt, uintptr(idx), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(value)), ^uintptr(0))
	runtime.KeepAlive(buf)
	if int32(rc) != winSQLiteOK {
		return syscall.Errno(rc)
	}
	return nil
}

func bindWinIdentityInt64(stmt uintptr, idx int, value int64) error {
	rc, _, _ := winSQLiteIdentityBindInt64.Call(stmt, uintptr(idx), uintptr(value))
	if int32(rc) != winSQLiteOK {
		return syscall.Errno(rc)
	}
	return nil
}

func resetWinIdentityStatement(stmt uintptr) {
	_, _, _ = winSQLiteIdentityReset.Call(stmt)
	_, _, _ = winSQLiteIdentityClearBindings.Call(stmt)
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
	stmt, err := b.prepare(`INSERT INTO instrument_identities(symbol,display_name,exchange,asset_class,provider_asset_id,source,observed_at_ms)
VALUES(?,?,?,?,?,?,?)
ON CONFLICT(symbol) DO UPDATE SET display_name=excluded.display_name,exchange=excluded.exchange,asset_class=excluded.asset_class,provider_asset_id=excluded.provider_asset_id,source=excluded.source,observed_at_ms=excluded.observed_at_ms
WHERE excluded.observed_at_ms >= instrument_identities.observed_at_ms`)
	if err != nil {
		return 0, err
	}
	defer finalizeWinSQLite(stmt)
	written := 0
	for _, record := range canonicalInstrumentIdentities(records) {
		if err := ctx.Err(); err != nil {
			return written, err
		}
		resetWinIdentityStatement(stmt)
		if err := bindWinIdentityText(stmt, 1, record.Symbol); err != nil {
			return written, b.sqliteErr()
		}
		if err := bindWinIdentityText(stmt, 2, record.Name); err != nil {
			return written, b.sqliteErr()
		}
		if err := bindWinIdentityText(stmt, 3, record.Exchange); err != nil {
			return written, b.sqliteErr()
		}
		if err := bindWinIdentityText(stmt, 4, record.AssetClass); err != nil {
			return written, b.sqliteErr()
		}
		if err := bindWinIdentityText(stmt, 5, record.ProviderAssetID); err != nil {
			return written, b.sqliteErr()
		}
		if err := bindWinIdentityText(stmt, 6, record.Source); err != nil {
			return written, b.sqliteErr()
		}
		if err := bindWinIdentityInt64(stmt, 7, record.ObservedAt); err != nil {
			return written, b.sqliteErr()
		}
		if rc := stepWinSQLite(stmt); rc != winSQLiteDone {
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

func (b *sqlitePersistenceBackend) LoadInstrumentIdentities(ctx context.Context) ([]InstrumentIdentityRecord, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := b.ensureInstrumentIdentitySchemaLocked(); err != nil {
		return nil, err
	}
	stmt, err := b.prepare(`SELECT symbol,display_name,exchange,asset_class,provider_asset_id,source,observed_at_ms FROM instrument_identities ORDER BY symbol`)
	if err != nil {
		return nil, err
	}
	defer finalizeWinSQLite(stmt)
	out := []InstrumentIdentityRecord{}
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		switch stepWinSQLite(stmt) {
		case winSQLiteDone:
			return out, nil
		case winSQLiteRow:
			out = append(out, InstrumentIdentityRecord{
				Symbol:          normalizeSymbol(winSQLiteText(stmt, 0)),
				Name:            winSQLiteText(stmt, 1),
				Exchange:        winSQLiteText(stmt, 2),
				AssetClass:      winSQLiteText(stmt, 3),
				ProviderAssetID: winSQLiteText(stmt, 4),
				Source:          winSQLiteText(stmt, 5),
				ObservedAt:      winSQLiteInt64(stmt, 6),
			})
		default:
			return nil, b.sqliteErr()
		}
	}
}
