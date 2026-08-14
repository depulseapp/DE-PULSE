//go:build cgo && !windows

package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestV170NativeDesktopPersistenceIsSQLite(t *testing.T) {
	dir := t.TempDir()
	p := NewPersistenceManager(dir)
	if d := p.Diagnostics(); !d.Ready || d.Backend != "sqlite" {
		t.Fatalf("native v17 persistence must be SQLite: %+v", d)
	}
	p.EnqueueQuotes(map[string]Quote{"SPY": {Symbol: "SPY", Price: 600, Source: "test", ProviderTimestamp: 1000, UpdatedAt: 1010}})
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(dir, "depulse-v17.db"))
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) < 16 || string(raw[:16]) != "SQLite format 3\x00" {
		t.Fatalf("persistence file is not SQLite format: %q", raw[:min(16, len(raw))])
	}
}
