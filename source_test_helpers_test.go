package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// productionGoSourceForTest makes source-contract tests refactor-safe: they
// validate the production source set rather than a historical monolith filename.
func productionGoSourceForTest(t *testing.T) string {
	t.Helper()
	entries, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(entries)
	var b strings.Builder
	for _, name := range entries {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		raw, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		b.Write(raw)
		b.WriteByte('\n')
	}
	return b.String()
}

// isolateUserConfig makes tests that construct a real Application portable across
// Linux, macOS, and Windows. os.UserConfigDir consults different environment
// variables on each platform, so setting only HOME/XDG_CONFIG_HOME is not enough.
func isolateUserConfig(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	t.Setenv("HOME", root)
	t.Setenv("XDG_CONFIG_HOME", root)
	t.Setenv("APPDATA", root)
	base, err := os.UserConfigDir()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(base, 0700); err != nil {
		t.Fatal(err)
	}
	return base
}

// registerApplicationCleanup closes the real persistence backend before t.TempDir
// cleanup. This is required on Windows, where an open SQLite handle prevents the
// temporary profile directory from being removed.
func registerApplicationCleanup(t *testing.T, app *Application) {
	t.Helper()
	if app == nil {
		return
	}
	t.Cleanup(func() {
		if app.persistence != nil {
			_ = app.persistence.Close()
		}
	})
}
