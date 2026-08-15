//go:build cgo && !windows

package main

/*
#cgo linux LDFLAGS: -lsqlite3
#cgo darwin LDFLAGS: -lsqlite3
#include <stdlib.h>
#include <sqlite3.h>
*/
import "C"

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	"unsafe"
)

type sqlitePersistenceBackend struct {
	mu   sync.Mutex
	path string
	db   *C.sqlite3
}

func newLocalPersistenceBackend(configDir string) PersistenceBackend {
	return &sqlitePersistenceBackend{path: filepath.Join(configDir, "depulse-v17.db")}
}
func (b *sqlitePersistenceBackend) Name() string { return "sqlite" }

func (b *sqlitePersistenceBackend) sqliteErr() error {
	if b.db == nil {
		return errors.New("sqlite database is not open")
	}
	return errors.New(C.GoString(C.sqlite3_errmsg(b.db)))
}

func (b *sqlitePersistenceBackend) exec(sql string) error {
	csql := C.CString(sql)
	defer C.free(unsafe.Pointer(csql))
	var msg *C.char
	if rc := C.sqlite3_exec(b.db, csql, nil, nil, &msg); rc != C.SQLITE_OK {
		if msg != nil {
			defer C.sqlite3_free(unsafe.Pointer(msg))
			return errors.New(C.GoString(msg))
		}
		return b.sqliteErr()
	}
	return nil
}

func (b *sqlitePersistenceBackend) Init(context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	cpath := C.CString(b.path)
	defer C.free(unsafe.Pointer(cpath))
	flags := C.int(C.SQLITE_OPEN_READWRITE | C.SQLITE_OPEN_CREATE | C.SQLITE_OPEN_FULLMUTEX)
	if rc := C.sqlite3_open_v2(cpath, &b.db, flags, nil); rc != C.SQLITE_OK {
		return b.sqliteErr()
	}
	C.sqlite3_busy_timeout(b.db, 2500)
	if err := b.exec(`PRAGMA journal_mode=WAL; PRAGMA synchronous=NORMAL; PRAGMA temp_store=MEMORY; CREATE TABLE IF NOT EXISTS schema_migrations(version INTEGER PRIMARY KEY, applied_at_ms INTEGER NOT NULL);`); err != nil {
		return err
	}
	migrations := []struct {
		version int
		sql     string
	}{
		{1, `
CREATE TABLE IF NOT EXISTS symbol_registry(
 symbol TEXT PRIMARY KEY,
 first_seen_ms INTEGER NOT NULL,
 last_seen_ms INTEGER NOT NULL,
 active INTEGER NOT NULL,
 selected INTEGER NOT NULL,
 processing_tier INTEGER NOT NULL,
 desk_membership TEXT NOT NULL DEFAULT '[]',
 provider_eligible INTEGER NOT NULL,
 last_subscribed_ms INTEGER NOT NULL DEFAULT 0,
 last_processed_ms INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_symbol_registry_active_tier ON symbol_registry(active,processing_tier);
CREATE TABLE IF NOT EXISTS canonical_quotes(
 symbol TEXT PRIMARY KEY,
 payload_json BLOB NOT NULL,
 provider_timestamp_ms INTEGER NOT NULL DEFAULT 0,
 received_timestamp_ms INTEGER NOT NULL DEFAULT 0,
 persisted_at_ms INTEGER NOT NULL,
 source TEXT NOT NULL DEFAULT '',
 data_state TEXT NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_canonical_quotes_persisted ON canonical_quotes(persisted_at_ms);
`},
		{2, `
CREATE TABLE IF NOT EXISTS quote_history(
 symbol TEXT NOT NULL,
 bucket_ms INTEGER NOT NULL,
 provider_timestamp_ms INTEGER NOT NULL DEFAULT 0,
 received_timestamp_ms INTEGER NOT NULL DEFAULT 0,
 price REAL NOT NULL,
 bid REAL NOT NULL DEFAULT 0,
 ask REAL NOT NULL DEFAULT 0,
 volume REAL NOT NULL DEFAULT 0,
 source TEXT NOT NULL DEFAULT '',
 data_state TEXT NOT NULL DEFAULT '',
 PRIMARY KEY(symbol,bucket_ms)
);
CREATE INDEX IF NOT EXISTS idx_quote_history_bucket ON quote_history(bucket_ms);
CREATE TABLE IF NOT EXISTS evidence_records(
 evidence_id TEXT PRIMARY KEY,
 symbol TEXT NOT NULL DEFAULT '',
 evidence_kind TEXT NOT NULL,
 observed_at_ms INTEGER NOT NULL,
 source TEXT NOT NULL DEFAULT '',
 provenance TEXT NOT NULL DEFAULT '',
 freshness_state TEXT NOT NULL DEFAULT '',
 payload_json BLOB NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_evidence_symbol_time ON evidence_records(symbol,observed_at_ms DESC);
CREATE TABLE IF NOT EXISTS decision_lineage(
 decision_id TEXT PRIMARY KEY,
 symbol TEXT NOT NULL DEFAULT '',
 horizon TEXT NOT NULL DEFAULT '',
 evidence_id TEXT NOT NULL DEFAULT '',
 decision_kind TEXT NOT NULL,
 decision_value TEXT NOT NULL DEFAULT '',
 formula_version TEXT NOT NULL DEFAULT '',
 created_at_ms INTEGER NOT NULL,
 payload_json BLOB NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_decision_lineage_symbol_time ON decision_lineage(symbol,created_at_ms DESC);
CREATE TABLE IF NOT EXISTS outcome_history(
 outcome_id TEXT PRIMARY KEY,
 decision_id TEXT NOT NULL DEFAULT '',
 symbol TEXT NOT NULL DEFAULT '',
 horizon TEXT NOT NULL DEFAULT '',
 observed_at_ms INTEGER NOT NULL,
 outcome_label TEXT NOT NULL DEFAULT '',
 payload_json BLOB NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_outcome_symbol_time ON outcome_history(symbol,observed_at_ms DESC);
CREATE TABLE IF NOT EXISTS derived_features(
 symbol TEXT NOT NULL,
 feature_key TEXT NOT NULL,
 feature_version TEXT NOT NULL,
 as_of_ms INTEGER NOT NULL,
 source_hash TEXT NOT NULL DEFAULT '',
 payload_json BLOB NOT NULL,
 PRIMARY KEY(symbol,feature_key,feature_version)
);
CREATE INDEX IF NOT EXISTS idx_derived_features_asof ON derived_features(as_of_ms DESC);
`},
		{3, `
CREATE TABLE IF NOT EXISTS identity_state(
 id INTEGER PRIMARY KEY CHECK(id=1),
 payload_json BLOB NOT NULL,
 updated_at_ms INTEGER NOT NULL
);
`},
		{4, `
CREATE TABLE IF NOT EXISTS user_workspaces(
 user_id TEXT PRIMARY KEY,
 payload_json BLOB NOT NULL,
 updated_at_ms INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_user_workspaces_updated ON user_workspaces(updated_at_ms DESC);
`},
	}
	for _, migration := range migrations {
		applied, err := b.migrationApplied(migration.version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}
		if err := b.exec("BEGIN IMMEDIATE"); err != nil {
			return err
		}
		if err := b.exec(migration.sql); err != nil {
			_ = b.exec("ROLLBACK")
			return fmt.Errorf("sqlite migration %d: %w", migration.version, err)
		}
		if err := b.exec(fmt.Sprintf("INSERT INTO schema_migrations(version,applied_at_ms) VALUES(%d, strftime('%%s','now')*1000)", migration.version)); err != nil {
			_ = b.exec("ROLLBACK")
			return fmt.Errorf("sqlite migration %d record: %w", migration.version, err)
		}
		if err := b.exec("COMMIT"); err != nil {
			_ = b.exec("ROLLBACK")
			return fmt.Errorf("sqlite migration %d commit: %w", migration.version, err)
		}
	}
	return nil
}

func (b *sqlitePersistenceBackend) migrationApplied(version int) (bool, error) {
	stmt, err := prepare(b.db, `SELECT 1 FROM schema_migrations WHERE version=? LIMIT 1`)
	if err != nil {
		return false, err
	}
	defer C.sqlite3_finalize(stmt)
	C.sqlite3_bind_int(stmt, 1, C.int(version))
	rc := C.sqlite3_step(stmt)
	if rc == C.SQLITE_ROW {
		return true, nil
	}
	if rc == C.SQLITE_DONE {
		return false, nil
	}
	return false, b.sqliteErr()
}

func prepare(db *C.sqlite3, sql string) (*C.sqlite3_stmt, error) {
	csql := C.CString(sql)
	defer C.free(unsafe.Pointer(csql))
	var stmt *C.sqlite3_stmt
	if rc := C.sqlite3_prepare_v2(db, csql, -1, &stmt, nil); rc != C.SQLITE_OK {
		return nil, errors.New(C.GoString(C.sqlite3_errmsg(db)))
	}
	return stmt, nil
}
func bindText(stmt *C.sqlite3_stmt, idx int, value string) {
	c := C.CString(value)
	defer C.free(unsafe.Pointer(c))
	C.sqlite3_bind_text(stmt, C.int(idx), c, C.int(len(value)), C.SQLITE_TRANSIENT)
}
func bindBlob(stmt *C.sqlite3_stmt, idx int, value []byte) {
	if len(value) == 0 {
		C.sqlite3_bind_blob(stmt, C.int(idx), nil, 0, C.SQLITE_TRANSIENT)
		return
	}
	C.sqlite3_bind_blob(stmt, C.int(idx), unsafe.Pointer(&value[0]), C.int(len(value)), C.SQLITE_TRANSIENT)
}

func (b *sqlitePersistenceBackend) UpsertSymbols(ctx context.Context, records []SymbolRegistryRecord) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return 0, err
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
	if err := b.exec("UPDATE symbol_registry SET active=0, selected=0"); err != nil {
		return 0, err
	}
	stmt, err := prepare(b.db, `INSERT INTO symbol_registry(symbol,first_seen_ms,last_seen_ms,active,selected,processing_tier,desk_membership,provider_eligible,last_subscribed_ms,last_processed_ms)
VALUES(?,?,?,?,?,?,?,?,?,?)
ON CONFLICT(symbol) DO UPDATE SET last_seen_ms=excluded.last_seen_ms,active=excluded.active,selected=excluded.selected,processing_tier=excluded.processing_tier,desk_membership=excluded.desk_membership,provider_eligible=excluded.provider_eligible,last_subscribed_ms=MAX(symbol_registry.last_subscribed_ms,excluded.last_subscribed_ms),last_processed_ms=MAX(symbol_registry.last_processed_ms,excluded.last_processed_ms)`)
	if err != nil {
		return 0, err
	}
	defer C.sqlite3_finalize(stmt)
	for _, r := range records {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		C.sqlite3_reset(stmt)
		C.sqlite3_clear_bindings(stmt)
		r.Symbol = normalizeSymbol(r.Symbol)
		if r.Symbol == "" {
			continue
		}
		firstSeen := r.FirstSeenAt
		if firstSeen <= 0 {
			firstSeen = r.LastSeenAt
		}
		bindText(stmt, 1, r.Symbol)
		C.sqlite3_bind_int64(stmt, 2, C.sqlite3_int64(firstSeen))
		C.sqlite3_bind_int64(stmt, 3, C.sqlite3_int64(r.LastSeenAt))
		if r.Active {
			C.sqlite3_bind_int(stmt, 4, 1)
		} else {
			C.sqlite3_bind_int(stmt, 4, 0)
		}
		if r.Selected {
			C.sqlite3_bind_int(stmt, 5, 1)
		} else {
			C.sqlite3_bind_int(stmt, 5, 0)
		}
		C.sqlite3_bind_int(stmt, 6, C.int(r.ProcessingTier))
		bindText(stmt, 7, r.DeskMembership)
		if r.ProviderEligible {
			C.sqlite3_bind_int(stmt, 8, 1)
		} else {
			C.sqlite3_bind_int(stmt, 8, 0)
		}
		C.sqlite3_bind_int64(stmt, 9, C.sqlite3_int64(r.LastSubscribedAt))
		C.sqlite3_bind_int64(stmt, 10, C.sqlite3_int64(r.LastProcessedAt))
		if rc := C.sqlite3_step(stmt); rc != C.SQLITE_DONE {
			return 0, b.sqliteErr()
		}
	}
	if err := b.exec("COMMIT"); err != nil {
		return 0, err
	}
	ok = true
	return len(records), nil
}

func (b *sqlitePersistenceBackend) LoadSymbols(ctx context.Context) ([]SymbolRegistryRecord, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	stmt, err := prepare(b.db, `SELECT symbol,first_seen_ms,last_seen_ms,active,selected,processing_tier,desk_membership,provider_eligible,last_subscribed_ms,last_processed_ms FROM symbol_registry ORDER BY symbol`)
	if err != nil {
		return nil, err
	}
	defer C.sqlite3_finalize(stmt)
	out := []SymbolRegistryRecord{}
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		rc := C.sqlite3_step(stmt)
		if rc == C.SQLITE_DONE {
			break
		}
		if rc != C.SQLITE_ROW {
			return nil, b.sqliteErr()
		}
		text := func(col int) string {
			ptr := C.sqlite3_column_text(stmt, C.int(col))
			if ptr == nil {
				return ""
			}
			return C.GoString((*C.char)(unsafe.Pointer(ptr)))
		}
		out = append(out, SymbolRegistryRecord{
			Symbol: normalizeSymbol(text(0)), FirstSeenAt: int64(C.sqlite3_column_int64(stmt, 1)), LastSeenAt: int64(C.sqlite3_column_int64(stmt, 2)),
			Active: C.sqlite3_column_int(stmt, 3) != 0, Selected: C.sqlite3_column_int(stmt, 4) != 0, ProcessingTier: int(C.sqlite3_column_int(stmt, 5)),
			DeskMembership: text(6), ProviderEligible: C.sqlite3_column_int(stmt, 7) != 0, LastSubscribedAt: int64(C.sqlite3_column_int64(stmt, 8)), LastProcessedAt: int64(C.sqlite3_column_int64(stmt, 9)),
		})
	}
	return out, nil
}

func (b *sqlitePersistenceBackend) SaveQuotes(ctx context.Context, quotes map[string]Quote) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return 0, err
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
	stmt, err := prepare(b.db, `INSERT INTO canonical_quotes(symbol,payload_json,provider_timestamp_ms,received_timestamp_ms,persisted_at_ms,source,data_state)
VALUES(?,?,?,?,?,?,?)
ON CONFLICT(symbol) DO UPDATE SET payload_json=excluded.payload_json,provider_timestamp_ms=excluded.provider_timestamp_ms,received_timestamp_ms=excluded.received_timestamp_ms,persisted_at_ms=excluded.persisted_at_ms,source=excluded.source,data_state=excluded.data_state`)
	if err != nil {
		return 0, err
	}
	defer C.sqlite3_finalize(stmt)
	historyStmt, err := prepare(b.db, `INSERT INTO quote_history(symbol,bucket_ms,provider_timestamp_ms,received_timestamp_ms,price,bid,ask,volume,source,data_state)
VALUES(?,?,?,?,?,?,?,?,?,?)
ON CONFLICT(symbol,bucket_ms) DO UPDATE SET provider_timestamp_ms=excluded.provider_timestamp_ms,received_timestamp_ms=excluded.received_timestamp_ms,price=excluded.price,bid=excluded.bid,ask=excluded.ask,volume=excluded.volume,source=excluded.source,data_state=excluded.data_state`)
	if err != nil {
		return 0, err
	}
	defer C.sqlite3_finalize(historyStmt)
	now := time.Now().UnixMilli()
	written := 0
	for sym, q := range quotes {
		if err := ctx.Err(); err != nil {
			return written, err
		}
		raw, err := json.Marshal(q)
		if err != nil {
			return written, err
		}
		C.sqlite3_reset(stmt)
		C.sqlite3_clear_bindings(stmt)
		bindText(stmt, 1, normalizeSymbol(sym))
		bindBlob(stmt, 2, raw)
		C.sqlite3_bind_int64(stmt, 3, C.sqlite3_int64(q.ProviderTimestamp))
		C.sqlite3_bind_int64(stmt, 4, C.sqlite3_int64(q.UpdatedAt))
		C.sqlite3_bind_int64(stmt, 5, C.sqlite3_int64(now))
		bindText(stmt, 6, q.Source)
		bindText(stmt, 7, q.DataState)
		if rc := C.sqlite3_step(stmt); rc != C.SQLITE_DONE {
			return written, b.sqliteErr()
		}
		stamp := q.ProviderTimestamp
		if stamp <= 0 {
			stamp = q.UpdatedAt
		}
		if stamp <= 0 {
			stamp = now
		}
		bucket := (stamp / 300000) * 300000
		C.sqlite3_reset(historyStmt)
		C.sqlite3_clear_bindings(historyStmt)
		bindText(historyStmt, 1, normalizeSymbol(sym))
		C.sqlite3_bind_int64(historyStmt, 2, C.sqlite3_int64(bucket))
		C.sqlite3_bind_int64(historyStmt, 3, C.sqlite3_int64(q.ProviderTimestamp))
		C.sqlite3_bind_int64(historyStmt, 4, C.sqlite3_int64(q.UpdatedAt))
		C.sqlite3_bind_double(historyStmt, 5, C.double(q.Price))
		C.sqlite3_bind_double(historyStmt, 6, C.double(q.Bid))
		C.sqlite3_bind_double(historyStmt, 7, C.double(q.Ask))
		C.sqlite3_bind_double(historyStmt, 8, C.double(q.Volume))
		bindText(historyStmt, 9, q.Source)
		bindText(historyStmt, 10, q.DataState)
		if rc := C.sqlite3_step(historyStmt); rc != C.SQLITE_DONE {
			return written, b.sqliteErr()
		}
		written++
	}
	if err := b.exec("DELETE FROM quote_history WHERE bucket_ms < (strftime('%s','now')*1000 - 2592000000)"); err != nil {
		return written, err
	}
	if err := b.exec("COMMIT"); err != nil {
		return written, err
	}
	ok = true
	return written, nil
}

func (b *sqlitePersistenceBackend) SaveIntelligence(ctx context.Context, batch PersistenceIntelligenceBatch) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if batch.Len() == 0 {
		return 0, nil
	}
	if err := b.exec("BEGIN IMMEDIATE"); err != nil {
		return 0, err
	}
	committed := false
	defer func() {
		if !committed {
			_ = b.exec("ROLLBACK")
		}
	}()
	written := 0

	evidenceStmt, err := prepare(b.db, `INSERT INTO evidence_records(evidence_id,symbol,evidence_kind,observed_at_ms,source,provenance,freshness_state,payload_json)
VALUES(?,?,?,?,?,?,?,?) ON CONFLICT(evidence_id) DO NOTHING`)
	if err != nil {
		return 0, err
	}
	defer C.sqlite3_finalize(evidenceStmt)
	for _, r := range batch.Evidence {
		if err := ctx.Err(); err != nil {
			return written, err
		}
		if r.ID == "" || r.Kind == "" {
			continue
		}
		C.sqlite3_reset(evidenceStmt)
		C.sqlite3_clear_bindings(evidenceStmt)
		bindText(evidenceStmt, 1, r.ID)
		bindText(evidenceStmt, 2, normalizeSymbol(r.Symbol))
		bindText(evidenceStmt, 3, r.Kind)
		C.sqlite3_bind_int64(evidenceStmt, 4, C.sqlite3_int64(r.ObservedAt))
		bindText(evidenceStmt, 5, r.Source)
		bindText(evidenceStmt, 6, r.Provenance)
		bindText(evidenceStmt, 7, r.FreshnessState)
		bindBlob(evidenceStmt, 8, payloadOrEmpty(r.Payload))
		if rc := C.sqlite3_step(evidenceStmt); rc != C.SQLITE_DONE {
			return written, b.sqliteErr()
		}
		written++
	}

	decisionStmt, err := prepare(b.db, `INSERT INTO decision_lineage(decision_id,symbol,horizon,evidence_id,decision_kind,decision_value,formula_version,created_at_ms,payload_json)
VALUES(?,?,?,?,?,?,?,?,?) ON CONFLICT(decision_id) DO NOTHING`)
	if err != nil {
		return written, err
	}
	defer C.sqlite3_finalize(decisionStmt)
	for _, r := range batch.Decisions {
		if err := ctx.Err(); err != nil {
			return written, err
		}
		if r.ID == "" || r.DecisionKind == "" {
			continue
		}
		C.sqlite3_reset(decisionStmt)
		C.sqlite3_clear_bindings(decisionStmt)
		bindText(decisionStmt, 1, r.ID)
		bindText(decisionStmt, 2, normalizeSymbol(r.Symbol))
		bindText(decisionStmt, 3, r.Horizon)
		bindText(decisionStmt, 4, r.EvidenceID)
		bindText(decisionStmt, 5, r.DecisionKind)
		bindText(decisionStmt, 6, r.DecisionValue)
		bindText(decisionStmt, 7, r.FormulaVersion)
		C.sqlite3_bind_int64(decisionStmt, 8, C.sqlite3_int64(r.CreatedAt))
		bindBlob(decisionStmt, 9, payloadOrEmpty(r.Payload))
		if rc := C.sqlite3_step(decisionStmt); rc != C.SQLITE_DONE {
			return written, b.sqliteErr()
		}
		written++
	}

	outcomeStmt, err := prepare(b.db, `INSERT INTO outcome_history(outcome_id,decision_id,symbol,horizon,observed_at_ms,outcome_label,payload_json)
VALUES(?,?,?,?,?,?,?) ON CONFLICT(outcome_id) DO UPDATE SET decision_id=excluded.decision_id,symbol=excluded.symbol,horizon=excluded.horizon,observed_at_ms=excluded.observed_at_ms,outcome_label=excluded.outcome_label,payload_json=excluded.payload_json`)
	if err != nil {
		return written, err
	}
	defer C.sqlite3_finalize(outcomeStmt)
	for _, r := range batch.Outcomes {
		if err := ctx.Err(); err != nil {
			return written, err
		}
		if r.ID == "" {
			continue
		}
		C.sqlite3_reset(outcomeStmt)
		C.sqlite3_clear_bindings(outcomeStmt)
		bindText(outcomeStmt, 1, r.ID)
		bindText(outcomeStmt, 2, r.DecisionID)
		bindText(outcomeStmt, 3, normalizeSymbol(r.Symbol))
		bindText(outcomeStmt, 4, r.Horizon)
		C.sqlite3_bind_int64(outcomeStmt, 5, C.sqlite3_int64(r.ObservedAt))
		bindText(outcomeStmt, 6, r.OutcomeLabel)
		bindBlob(outcomeStmt, 7, payloadOrEmpty(r.Payload))
		if rc := C.sqlite3_step(outcomeStmt); rc != C.SQLITE_DONE {
			return written, b.sqliteErr()
		}
		written++
	}

	featureStmt, err := prepare(b.db, `INSERT INTO derived_features(symbol,feature_key,feature_version,as_of_ms,source_hash,payload_json)
VALUES(?,?,?,?,?,?) ON CONFLICT(symbol,feature_key,feature_version) DO UPDATE SET as_of_ms=excluded.as_of_ms,source_hash=excluded.source_hash,payload_json=excluded.payload_json WHERE derived_features.source_hash<>excluded.source_hash`)
	if err != nil {
		return written, err
	}
	defer C.sqlite3_finalize(featureStmt)
	for _, r := range batch.Features {
		if err := ctx.Err(); err != nil {
			return written, err
		}
		if normalizeSymbol(r.Symbol) == "" || r.FeatureKey == "" || r.FeatureVersion == "" {
			continue
		}
		C.sqlite3_reset(featureStmt)
		C.sqlite3_clear_bindings(featureStmt)
		bindText(featureStmt, 1, normalizeSymbol(r.Symbol))
		bindText(featureStmt, 2, r.FeatureKey)
		bindText(featureStmt, 3, r.FeatureVersion)
		C.sqlite3_bind_int64(featureStmt, 4, C.sqlite3_int64(r.AsOf))
		bindText(featureStmt, 5, r.SourceHash)
		bindBlob(featureStmt, 6, payloadOrEmpty(r.Payload))
		if rc := C.sqlite3_step(featureStmt); rc != C.SQLITE_DONE {
			return written, b.sqliteErr()
		}
		written++
	}
	if err := b.exec("COMMIT"); err != nil {
		return written, err
	}
	committed = true
	return written, nil
}

func (b *sqlitePersistenceBackend) LoadQuotes(ctx context.Context) (map[string]Quote, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	stmt, err := prepare(b.db, `SELECT symbol,payload_json FROM canonical_quotes ORDER BY persisted_at_ms DESC LIMIT 1000`)
	if err != nil {
		return nil, err
	}
	defer C.sqlite3_finalize(stmt)
	out := map[string]Quote{}
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		rc := C.sqlite3_step(stmt)
		if rc == C.SQLITE_DONE {
			break
		}
		if rc != C.SQLITE_ROW {
			return nil, b.sqliteErr()
		}
		sym := C.GoString((*C.char)(unsafe.Pointer(C.sqlite3_column_text(stmt, 0))))
		ptr := C.sqlite3_column_blob(stmt, 1)
		n := C.sqlite3_column_bytes(stmt, 1)
		if ptr == nil || n <= 0 {
			continue
		}
		raw := C.GoBytes(ptr, n)
		var q Quote
		if err := json.Unmarshal(raw, &q); err != nil {
			continue
		}
		q.Symbol = normalizeSymbol(sym)
		q.DataState = "persisted"
		q.FeedType = "persisted"
		out[q.Symbol] = q
	}
	return out, nil
}

func (b *sqlitePersistenceBackend) scalarInt64(sql string) (int64, error) {
	stmt, err := prepare(b.db, sql)
	if err != nil {
		return 0, err
	}
	defer C.sqlite3_finalize(stmt)
	if rc := C.sqlite3_step(stmt); rc != C.SQLITE_ROW {
		return 0, b.sqliteErr()
	}
	return int64(C.sqlite3_column_int64(stmt, 0)), nil
}

func (b *sqlitePersistenceBackend) LoadIdentityState(ctx context.Context) (IdentityPersistentState, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return IdentityPersistentState{}, err
	}
	stmt, err := prepare(b.db, `SELECT payload_json FROM identity_state WHERE id=1`)
	if err != nil {
		return IdentityPersistentState{}, err
	}
	defer C.sqlite3_finalize(stmt)
	rc := C.sqlite3_step(stmt)
	if rc == C.SQLITE_DONE {
		return IdentityPersistentState{}, nil
	}
	if rc != C.SQLITE_ROW {
		return IdentityPersistentState{}, b.sqliteErr()
	}
	ptr := C.sqlite3_column_blob(stmt, 0)
	n := int(C.sqlite3_column_bytes(stmt, 0))
	if ptr == nil || n == 0 {
		return IdentityPersistentState{}, nil
	}
	raw := C.GoBytes(ptr, C.int(n))
	var state IdentityPersistentState
	if err := json.Unmarshal(raw, &state); err != nil {
		return IdentityPersistentState{}, err
	}
	return state, nil
}

func (b *sqlitePersistenceBackend) SaveIdentityState(ctx context.Context, state IdentityPersistentState) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	raw, err := json.Marshal(state)
	if err != nil {
		return err
	}
	stmt, err := prepare(b.db, `INSERT INTO identity_state(id,payload_json,updated_at_ms) VALUES(1,?,?) ON CONFLICT(id) DO UPDATE SET payload_json=excluded.payload_json,updated_at_ms=excluded.updated_at_ms`)
	if err != nil {
		return err
	}
	defer C.sqlite3_finalize(stmt)
	bindBlob(stmt, 1, raw)
	C.sqlite3_bind_int64(stmt, 2, C.sqlite3_int64(state.UpdatedAt))
	if rc := C.sqlite3_step(stmt); rc != C.SQLITE_DONE {
		return b.sqliteErr()
	}
	return nil
}

func (b *sqlitePersistenceBackend) LoadUserWorkspaces(ctx context.Context) ([]UserWorkspace, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	stmt, err := prepare(b.db, `SELECT payload_json FROM user_workspaces ORDER BY user_id`)
	if err != nil {
		return nil, err
	}
	defer C.sqlite3_finalize(stmt)
	out := []UserWorkspace{}
	for {
		rc := C.sqlite3_step(stmt)
		if rc == C.SQLITE_DONE {
			break
		}
		if rc != C.SQLITE_ROW {
			return nil, b.sqliteErr()
		}
		ptr := C.sqlite3_column_blob(stmt, 0)
		n := int(C.sqlite3_column_bytes(stmt, 0))
		if ptr == nil || n == 0 {
			continue
		}
		var workspace UserWorkspace
		if err := json.Unmarshal(C.GoBytes(ptr, C.int(n)), &workspace); err != nil {
			return nil, err
		}
		out = append(out, workspace)
	}
	return out, nil
}

func (b *sqlitePersistenceBackend) SaveUserWorkspace(ctx context.Context, workspace UserWorkspace) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	if workspace.UserID == "" {
		return errors.New("workspace user id is required")
	}
	raw, err := json.Marshal(workspace)
	if err != nil {
		return err
	}
	stmt, err := prepare(b.db, `INSERT INTO user_workspaces(user_id,payload_json,updated_at_ms) VALUES(?,?,?) ON CONFLICT(user_id) DO UPDATE SET payload_json=excluded.payload_json,updated_at_ms=excluded.updated_at_ms`)
	if err != nil {
		return err
	}
	defer C.sqlite3_finalize(stmt)
	bindText(stmt, 1, workspace.UserID)
	bindBlob(stmt, 2, raw)
	C.sqlite3_bind_int64(stmt, 3, C.sqlite3_int64(workspace.UpdatedAt))
	if rc := C.sqlite3_step(stmt); rc != C.SQLITE_DONE {
		return b.sqliteErr()
	}
	return nil
}

func (b *sqlitePersistenceBackend) identityCountsLocked() (int, int, error) {
	stmt, err := prepare(b.db, `SELECT payload_json FROM identity_state WHERE id=1`)
	if err != nil {
		return 0, 0, err
	}
	defer C.sqlite3_finalize(stmt)
	rc := C.sqlite3_step(stmt)
	if rc == C.SQLITE_DONE {
		return 0, 0, nil
	}
	if rc != C.SQLITE_ROW {
		return 0, 0, b.sqliteErr()
	}
	ptr := C.sqlite3_column_blob(stmt, 0)
	n := int(C.sqlite3_column_bytes(stmt, 0))
	if ptr == nil || n == 0 {
		return 0, 0, nil
	}
	var state IdentityPersistentState
	if err := json.Unmarshal(C.GoBytes(ptr, C.int(n)), &state); err != nil {
		return 0, 0, err
	}
	return len(state.Users), len(state.Sessions), nil
}

func (b *sqlitePersistenceBackend) Stats(ctx context.Context) (PersistenceStoreStats, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return PersistenceStoreStats{}, err
	}
	queries := []string{
		"SELECT COALESCE(MAX(version),0) FROM schema_migrations",
		"SELECT COUNT(*) FROM symbol_registry",
		"SELECT COUNT(*) FROM symbol_registry WHERE active=1",
		"SELECT COUNT(*) FROM canonical_quotes",
		"SELECT COUNT(*) FROM quote_history",
		"SELECT COUNT(*) FROM evidence_records",
		"SELECT COUNT(*) FROM decision_lineage",
		"SELECT COUNT(*) FROM outcome_history",
		"SELECT COUNT(*) FROM derived_features",
	}
	vals := make([]int64, len(queries))
	for i, q := range queries {
		v, err := b.scalarInt64(q)
		if err != nil {
			return PersistenceStoreStats{}, err
		}
		vals[i] = v
	}
	users, sessions, err := b.identityCountsLocked()
	if err != nil {
		return PersistenceStoreStats{}, err
	}
	var size int64
	for _, path := range []string{b.path, b.path + "-wal", b.path + "-shm"} {
		if info, err := os.Stat(path); err == nil {
			size += info.Size()
		}
	}
	return PersistenceStoreStats{SchemaVersion: int(vals[0]), SymbolCount: int(vals[1]), ActiveSymbolCount: int(vals[2]), CanonicalQuotes: int(vals[3]), QuoteHistoryRows: int(vals[4]), EvidenceRows: int(vals[5]), DecisionRows: int(vals[6]), OutcomeRows: int(vals[7]), FeatureRows: int(vals[8]), UserCount: users, SessionCount: sessions, StorageBytes: size}, nil
}

func (b *sqlitePersistenceBackend) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.db == nil {
		return nil
	}
	if rc := C.sqlite3_close_v2(b.db); rc != C.SQLITE_OK {
		return fmt.Errorf("sqlite close: %s", C.GoString(C.sqlite3_errmsg(b.db)))
	}
	b.db = nil
	return nil
}
