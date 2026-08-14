//go:build windows

package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

// Windows uses the OS-provided winsqlite3.dll so desktop persistence remains
// real SQLite even when DE.PULSE is cross-compiled with CGO_ENABLED=0.
// The application currently packages Windows x64; on amd64 the Windows ABI is
// compatible with syscall.Proc.Call for the SQLite C entry points used here.
const (
	winSQLiteOK            = 0
	winSQLiteRow           = 100
	winSQLiteDone          = 101
	winSQLiteOpenReadWrite = 0x00000002
	winSQLiteOpenCreate    = 0x00000004
	winSQLiteOpenFullMutex = 0x00010000
)

var (
	winSQLiteDLL         = syscall.NewLazyDLL("winsqlite3.dll")
	winSQLiteOpenV2      = winSQLiteDLL.NewProc("sqlite3_open_v2")
	winSQLiteCloseV2     = winSQLiteDLL.NewProc("sqlite3_close_v2")
	winSQLiteErrmsg      = winSQLiteDLL.NewProc("sqlite3_errmsg")
	winSQLiteExec        = winSQLiteDLL.NewProc("sqlite3_exec")
	winSQLiteBusyTimeout = winSQLiteDLL.NewProc("sqlite3_busy_timeout")
	winSQLitePrepareV2   = winSQLiteDLL.NewProc("sqlite3_prepare_v2")
	winSQLiteStep        = winSQLiteDLL.NewProc("sqlite3_step")
	winSQLiteFinalize    = winSQLiteDLL.NewProc("sqlite3_finalize")
	winSQLiteColumnText  = winSQLiteDLL.NewProc("sqlite3_column_text")
	winSQLiteColumnInt   = winSQLiteDLL.NewProc("sqlite3_column_int")
	winSQLiteColumnInt64 = winSQLiteDLL.NewProc("sqlite3_column_int64")
	winSQLiteColumnBlob  = winSQLiteDLL.NewProc("sqlite3_column_blob")
	winSQLiteColumnBytes = winSQLiteDLL.NewProc("sqlite3_column_bytes")
)

type sqlitePersistenceBackend struct {
	mu   sync.Mutex
	path string
	db   uintptr
}

func newPersistenceBackend(configDir string) PersistenceBackend {
	return &sqlitePersistenceBackend{path: filepath.Join(configDir, "depulse-v17.db")}
}

func (b *sqlitePersistenceBackend) Name() string { return "sqlite" }

func winCString(s string) []byte { return append([]byte(s), 0) }

func winGoString(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}
	const maxCString = 16 << 20
	for n := 0; n < maxCString; n++ {
		if *(*byte)(unsafe.Pointer(ptr + uintptr(n))) == 0 {
			return string(unsafe.Slice((*byte)(unsafe.Pointer(ptr)), n))
		}
	}
	return ""
}

func (b *sqlitePersistenceBackend) sqliteErr() error {
	if b.db == 0 {
		return errors.New("sqlite database is not open")
	}
	ptr, _, _ := winSQLiteErrmsg.Call(b.db)
	msg := winGoString(ptr)
	if msg == "" {
		msg = "unknown SQLite error"
	}
	return errors.New(msg)
}

func (b *sqlitePersistenceBackend) exec(sqlText string) error {
	if b.db == 0 {
		return errors.New("sqlite database is not open")
	}
	buf := winCString(sqlText)
	rc, _, _ := winSQLiteExec.Call(b.db, uintptr(unsafe.Pointer(&buf[0])), 0, 0, 0)
	if int32(rc) != winSQLiteOK {
		return b.sqliteErr()
	}
	return nil
}

func (b *sqlitePersistenceBackend) Init(context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := winSQLiteDLL.Load(); err != nil {
		return fmt.Errorf("load Windows system SQLite (winsqlite3.dll): %w", err)
	}
	for name, proc := range map[string]*syscall.LazyProc{
		"sqlite3_open_v2": winSQLiteOpenV2, "sqlite3_close_v2": winSQLiteCloseV2,
		"sqlite3_errmsg": winSQLiteErrmsg, "sqlite3_exec": winSQLiteExec,
		"sqlite3_busy_timeout": winSQLiteBusyTimeout, "sqlite3_prepare_v2": winSQLitePrepareV2,
		"sqlite3_step": winSQLiteStep, "sqlite3_finalize": winSQLiteFinalize,
		"sqlite3_column_text": winSQLiteColumnText, "sqlite3_column_int": winSQLiteColumnInt,
		"sqlite3_column_int64": winSQLiteColumnInt64, "sqlite3_column_blob": winSQLiteColumnBlob,
		"sqlite3_column_bytes": winSQLiteColumnBytes,
	} {
		if err := proc.Find(); err != nil {
			return fmt.Errorf("Windows system SQLite missing %s: %w", name, err)
		}
	}
	if err := os.MkdirAll(filepath.Dir(b.path), 0700); err != nil {
		return err
	}
	pathBuf := winCString(b.path)
	flags := uintptr(winSQLiteOpenReadWrite | winSQLiteOpenCreate | winSQLiteOpenFullMutex)
	rc, _, _ := winSQLiteOpenV2.Call(uintptr(unsafe.Pointer(&pathBuf[0])), uintptr(unsafe.Pointer(&b.db)), flags, 0)
	if int32(rc) != winSQLiteOK {
		err := b.sqliteErr()
		if b.db != 0 {
			_, _, _ = winSQLiteCloseV2.Call(b.db)
			b.db = 0
		}
		return err
	}
	_, _, _ = winSQLiteBusyTimeout.Call(b.db, 2500)
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

func (b *sqlitePersistenceBackend) prepare(sqlText string) (uintptr, error) {
	buf := winCString(sqlText)
	var stmt uintptr
	rc, _, _ := winSQLitePrepareV2.Call(b.db, uintptr(unsafe.Pointer(&buf[0])), ^uintptr(0), uintptr(unsafe.Pointer(&stmt)), 0)
	if int32(rc) != winSQLiteOK {
		return 0, b.sqliteErr()
	}
	return stmt, nil
}

func finalizeWinSQLite(stmt uintptr) {
	if stmt != 0 {
		_, _, _ = winSQLiteFinalize.Call(stmt)
	}
}

func stepWinSQLite(stmt uintptr) int32 {
	rc, _, _ := winSQLiteStep.Call(stmt)
	return int32(rc)
}

func winSQLiteText(stmt uintptr, col int) string {
	ptr, _, _ := winSQLiteColumnText.Call(stmt, uintptr(col))
	return winGoString(ptr)
}

func winSQLiteInt64(stmt uintptr, col int) int64 {
	v, _, _ := winSQLiteColumnInt64.Call(stmt, uintptr(col))
	return int64(v)
}

func winSQLiteInt(stmt uintptr, col int) int {
	v, _, _ := winSQLiteColumnInt.Call(stmt, uintptr(col))
	return int(int32(v))
}

func winSQLiteBlob(stmt uintptr, col int) []byte {
	ptr, _, _ := winSQLiteColumnBlob.Call(stmt, uintptr(col))
	n, _, _ := winSQLiteColumnBytes.Call(stmt, uintptr(col))
	if ptr == 0 || int32(n) <= 0 {
		return nil
	}
	src := unsafe.Slice((*byte)(unsafe.Pointer(ptr)), int(int32(n)))
	return append([]byte(nil), src...)
}

func (b *sqlitePersistenceBackend) migrationApplied(version int) (bool, error) {
	stmt, err := b.prepare(fmt.Sprintf(`SELECT 1 FROM schema_migrations WHERE version=%d LIMIT 1`, version))
	if err != nil {
		return false, err
	}
	defer finalizeWinSQLite(stmt)
	switch stepWinSQLite(stmt) {
	case winSQLiteRow:
		return true, nil
	case winSQLiteDone:
		return false, nil
	default:
		return false, b.sqliteErr()
	}
}

func sqlText(v string) string { return "'" + strings.ReplaceAll(v, "'", "''") + "'" }
func sqlBlob(v []byte) string { return "X'" + hex.EncodeToString(v) + "'" }
func sqlBool(v bool) string {
	if v {
		return "1"
	}
	return "0"
}
func sqlInt64(v int64) string   { return strconv.FormatInt(v, 10) }
func sqlInt(v int) string       { return strconv.Itoa(v) }
func sqlFloat(v float64) string { return strconv.FormatFloat(v, 'g', -1, 64) }

func (b *sqlitePersistenceBackend) UpsertSymbols(ctx context.Context, records []SymbolRegistryRecord) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	var q strings.Builder
	q.WriteString("BEGIN IMMEDIATE; UPDATE symbol_registry SET active=0, selected=0;")
	written := 0
	for _, r := range records {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		r.Symbol = normalizeSymbol(r.Symbol)
		if r.Symbol == "" {
			continue
		}
		firstSeen := r.FirstSeenAt
		if firstSeen <= 0 {
			firstSeen = r.LastSeenAt
		}
		q.WriteString("INSERT INTO symbol_registry(symbol,first_seen_ms,last_seen_ms,active,selected,processing_tier,desk_membership,provider_eligible,last_subscribed_ms,last_processed_ms) VALUES(")
		q.WriteString(strings.Join([]string{sqlText(r.Symbol), sqlInt64(firstSeen), sqlInt64(r.LastSeenAt), sqlBool(r.Active), sqlBool(r.Selected), sqlInt(r.ProcessingTier), sqlText(r.DeskMembership), sqlBool(r.ProviderEligible), sqlInt64(r.LastSubscribedAt), sqlInt64(r.LastProcessedAt)}, ","))
		q.WriteString(") ON CONFLICT(symbol) DO UPDATE SET last_seen_ms=excluded.last_seen_ms,active=excluded.active,selected=excluded.selected,processing_tier=excluded.processing_tier,desk_membership=excluded.desk_membership,provider_eligible=excluded.provider_eligible,last_subscribed_ms=MAX(symbol_registry.last_subscribed_ms,excluded.last_subscribed_ms),last_processed_ms=MAX(symbol_registry.last_processed_ms,excluded.last_processed_ms);")
		written++
	}
	q.WriteString("COMMIT;")
	if err := b.exec(q.String()); err != nil {
		_ = b.exec("ROLLBACK")
		return 0, err
	}
	return written, nil
}

func (b *sqlitePersistenceBackend) LoadSymbols(ctx context.Context) ([]SymbolRegistryRecord, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	stmt, err := b.prepare(`SELECT symbol,first_seen_ms,last_seen_ms,active,selected,processing_tier,desk_membership,provider_eligible,last_subscribed_ms,last_processed_ms FROM symbol_registry ORDER BY symbol`)
	if err != nil {
		return nil, err
	}
	defer finalizeWinSQLite(stmt)
	out := []SymbolRegistryRecord{}
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		rc := stepWinSQLite(stmt)
		if rc == winSQLiteDone {
			break
		}
		if rc != winSQLiteRow {
			return nil, b.sqliteErr()
		}
		out = append(out, SymbolRegistryRecord{
			Symbol: normalizeSymbol(winSQLiteText(stmt, 0)), FirstSeenAt: winSQLiteInt64(stmt, 1), LastSeenAt: winSQLiteInt64(stmt, 2),
			Active: winSQLiteInt(stmt, 3) != 0, Selected: winSQLiteInt(stmt, 4) != 0, ProcessingTier: winSQLiteInt(stmt, 5), DeskMembership: winSQLiteText(stmt, 6),
			ProviderEligible: winSQLiteInt(stmt, 7) != 0, LastSubscribedAt: winSQLiteInt64(stmt, 8), LastProcessedAt: winSQLiteInt64(stmt, 9),
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
	now := time.Now().UnixMilli()
	var q strings.Builder
	q.WriteString("BEGIN IMMEDIATE;")
	written := 0
	for sym, quote := range quotes {
		if err := ctx.Err(); err != nil {
			return written, err
		}
		raw, err := json.Marshal(quote)
		if err != nil {
			return written, err
		}
		sym = normalizeSymbol(sym)
		if sym == "" {
			continue
		}
		q.WriteString("INSERT INTO canonical_quotes(symbol,payload_json,provider_timestamp_ms,received_timestamp_ms,persisted_at_ms,source,data_state) VALUES(")
		q.WriteString(strings.Join([]string{sqlText(sym), sqlBlob(raw), sqlInt64(quote.ProviderTimestamp), sqlInt64(quote.UpdatedAt), sqlInt64(now), sqlText(quote.Source), sqlText(quote.DataState)}, ","))
		q.WriteString(") ON CONFLICT(symbol) DO UPDATE SET payload_json=excluded.payload_json,provider_timestamp_ms=excluded.provider_timestamp_ms,received_timestamp_ms=excluded.received_timestamp_ms,persisted_at_ms=excluded.persisted_at_ms,source=excluded.source,data_state=excluded.data_state;")
		stamp := quote.ProviderTimestamp
		if stamp <= 0 {
			stamp = quote.UpdatedAt
		}
		if stamp <= 0 {
			stamp = now
		}
		bucket := (stamp / 300000) * 300000
		q.WriteString("INSERT INTO quote_history(symbol,bucket_ms,provider_timestamp_ms,received_timestamp_ms,price,bid,ask,volume,source,data_state) VALUES(")
		q.WriteString(strings.Join([]string{sqlText(sym), sqlInt64(bucket), sqlInt64(quote.ProviderTimestamp), sqlInt64(quote.UpdatedAt), sqlFloat(quote.Price), sqlFloat(quote.Bid), sqlFloat(quote.Ask), sqlFloat(quote.Volume), sqlText(quote.Source), sqlText(quote.DataState)}, ","))
		q.WriteString(") ON CONFLICT(symbol,bucket_ms) DO UPDATE SET provider_timestamp_ms=excluded.provider_timestamp_ms,received_timestamp_ms=excluded.received_timestamp_ms,price=excluded.price,bid=excluded.bid,ask=excluded.ask,volume=excluded.volume,source=excluded.source,data_state=excluded.data_state;")
		written++
	}
	q.WriteString("DELETE FROM quote_history WHERE bucket_ms < (strftime('%s','now')*1000 - 2592000000); COMMIT;")
	if err := b.exec(q.String()); err != nil {
		_ = b.exec("ROLLBACK")
		return written, err
	}
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
	var q strings.Builder
	q.WriteString("BEGIN IMMEDIATE;")
	written := 0
	for _, r := range batch.Evidence {
		if err := ctx.Err(); err != nil {
			return written, err
		}
		if r.ID == "" || r.Kind == "" {
			continue
		}
		vals := []string{sqlText(r.ID), sqlText(normalizeSymbol(r.Symbol)), sqlText(r.Kind), sqlInt64(r.ObservedAt), sqlText(r.Source), sqlText(r.Provenance), sqlText(r.FreshnessState), sqlBlob(payloadOrEmpty(r.Payload))}
		q.WriteString("INSERT INTO evidence_records(evidence_id,symbol,evidence_kind,observed_at_ms,source,provenance,freshness_state,payload_json) VALUES(" + strings.Join(vals, ",") + ") ON CONFLICT(evidence_id) DO NOTHING;")
		written++
	}
	for _, r := range batch.Decisions {
		if err := ctx.Err(); err != nil {
			return written, err
		}
		if r.ID == "" || r.DecisionKind == "" {
			continue
		}
		vals := []string{sqlText(r.ID), sqlText(normalizeSymbol(r.Symbol)), sqlText(r.Horizon), sqlText(r.EvidenceID), sqlText(r.DecisionKind), sqlText(r.DecisionValue), sqlText(r.FormulaVersion), sqlInt64(r.CreatedAt), sqlBlob(payloadOrEmpty(r.Payload))}
		q.WriteString("INSERT INTO decision_lineage(decision_id,symbol,horizon,evidence_id,decision_kind,decision_value,formula_version,created_at_ms,payload_json) VALUES(" + strings.Join(vals, ",") + ") ON CONFLICT(decision_id) DO NOTHING;")
		written++
	}
	for _, r := range batch.Outcomes {
		if err := ctx.Err(); err != nil {
			return written, err
		}
		if r.ID == "" {
			continue
		}
		vals := []string{sqlText(r.ID), sqlText(r.DecisionID), sqlText(normalizeSymbol(r.Symbol)), sqlText(r.Horizon), sqlInt64(r.ObservedAt), sqlText(r.OutcomeLabel), sqlBlob(payloadOrEmpty(r.Payload))}
		q.WriteString("INSERT INTO outcome_history(outcome_id,decision_id,symbol,horizon,observed_at_ms,outcome_label,payload_json) VALUES(" + strings.Join(vals, ",") + ") ON CONFLICT(outcome_id) DO UPDATE SET decision_id=excluded.decision_id,symbol=excluded.symbol,horizon=excluded.horizon,observed_at_ms=excluded.observed_at_ms,outcome_label=excluded.outcome_label,payload_json=excluded.payload_json;")
		written++
	}
	for _, r := range batch.Features {
		if err := ctx.Err(); err != nil {
			return written, err
		}
		sym := normalizeSymbol(r.Symbol)
		if sym == "" || r.FeatureKey == "" || r.FeatureVersion == "" {
			continue
		}
		vals := []string{sqlText(sym), sqlText(r.FeatureKey), sqlText(r.FeatureVersion), sqlInt64(r.AsOf), sqlText(r.SourceHash), sqlBlob(payloadOrEmpty(r.Payload))}
		q.WriteString("INSERT INTO derived_features(symbol,feature_key,feature_version,as_of_ms,source_hash,payload_json) VALUES(" + strings.Join(vals, ",") + ") ON CONFLICT(symbol,feature_key,feature_version) DO UPDATE SET as_of_ms=excluded.as_of_ms,source_hash=excluded.source_hash,payload_json=excluded.payload_json WHERE derived_features.source_hash<>excluded.source_hash;")
		written++
	}
	q.WriteString("COMMIT;")
	if err := b.exec(q.String()); err != nil {
		_ = b.exec("ROLLBACK")
		return written, err
	}
	return written, nil
}

func (b *sqlitePersistenceBackend) LoadQuotes(ctx context.Context) (map[string]Quote, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	stmt, err := b.prepare(`SELECT symbol,payload_json FROM canonical_quotes ORDER BY persisted_at_ms DESC LIMIT 1000`)
	if err != nil {
		return nil, err
	}
	defer finalizeWinSQLite(stmt)
	out := map[string]Quote{}
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		rc := stepWinSQLite(stmt)
		if rc == winSQLiteDone {
			break
		}
		if rc != winSQLiteRow {
			return nil, b.sqliteErr()
		}
		sym := normalizeSymbol(winSQLiteText(stmt, 0))
		raw := winSQLiteBlob(stmt, 1)
		if len(raw) == 0 {
			continue
		}
		var quote Quote
		if err := json.Unmarshal(raw, &quote); err != nil {
			continue
		}
		quote.Symbol = sym
		quote.DataState = "persisted"
		quote.FeedType = "persisted"
		out[sym] = quote
	}
	return out, nil
}

func (b *sqlitePersistenceBackend) scalarInt64(sqlText string) (int64, error) {
	stmt, err := b.prepare(sqlText)
	if err != nil {
		return 0, err
	}
	defer finalizeWinSQLite(stmt)
	if rc := stepWinSQLite(stmt); rc != winSQLiteRow {
		return 0, b.sqliteErr()
	}
	return winSQLiteInt64(stmt, 0), nil
}

func (b *sqlitePersistenceBackend) LoadIdentityState(ctx context.Context) (IdentityPersistentState, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return IdentityPersistentState{}, err
	}
	stmt, err := b.prepare(`SELECT payload_json FROM identity_state WHERE id=1`)
	if err != nil {
		return IdentityPersistentState{}, err
	}
	defer finalizeWinSQLite(stmt)
	rc := stepWinSQLite(stmt)
	if rc == winSQLiteDone {
		return IdentityPersistentState{}, nil
	}
	if rc != winSQLiteRow {
		return IdentityPersistentState{}, b.sqliteErr()
	}
	raw := winSQLiteBlob(stmt, 0)
	var state IdentityPersistentState
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &state); err != nil {
			return IdentityPersistentState{}, err
		}
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
	q := `INSERT INTO identity_state(id,payload_json,updated_at_ms) VALUES(1,` + sqlBlob(raw) + `,` + sqlInt64(state.UpdatedAt) + `) ON CONFLICT(id) DO UPDATE SET payload_json=excluded.payload_json,updated_at_ms=excluded.updated_at_ms;`
	return b.exec(q)
}

func (b *sqlitePersistenceBackend) LoadUserWorkspaces(ctx context.Context) ([]UserWorkspace, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	stmt, err := b.prepare(`SELECT payload_json FROM user_workspaces ORDER BY user_id`)
	if err != nil {
		return nil, err
	}
	defer finalizeWinSQLite(stmt)
	out := []UserWorkspace{}
	for {
		rc := stepWinSQLite(stmt)
		if rc == winSQLiteDone {
			break
		}
		if rc != winSQLiteRow {
			return nil, b.sqliteErr()
		}
		raw := winSQLiteBlob(stmt, 0)
		if len(raw) == 0 {
			continue
		}
		var workspace UserWorkspace
		if err := json.Unmarshal(raw, &workspace); err != nil {
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
	q := `INSERT INTO user_workspaces(user_id,payload_json,updated_at_ms) VALUES(` + sqlText(workspace.UserID) + `,` + sqlBlob(raw) + `,` + sqlInt64(workspace.UpdatedAt) + `) ON CONFLICT(user_id) DO UPDATE SET payload_json=excluded.payload_json,updated_at_ms=excluded.updated_at_ms;`
	return b.exec(q)
}

func (b *sqlitePersistenceBackend) identityCountsLocked() (int, int, error) {
	stmt, err := b.prepare(`SELECT payload_json FROM identity_state WHERE id=1`)
	if err != nil {
		return 0, 0, err
	}
	defer finalizeWinSQLite(stmt)
	rc := stepWinSQLite(stmt)
	if rc == winSQLiteDone {
		return 0, 0, nil
	}
	if rc != winSQLiteRow {
		return 0, 0, b.sqliteErr()
	}
	var state IdentityPersistentState
	if raw := winSQLiteBlob(stmt, 0); len(raw) > 0 {
		if err := json.Unmarshal(raw, &state); err != nil {
			return 0, 0, err
		}
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
	if b.db == 0 {
		return nil
	}
	rc, _, _ := winSQLiteCloseV2.Call(b.db)
	if int32(rc) != winSQLiteOK {
		return fmt.Errorf("sqlite close: %w", b.sqliteErr())
	}
	b.db = 0
	return nil
}
