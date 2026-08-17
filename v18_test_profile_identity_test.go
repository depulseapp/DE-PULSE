package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestV1851TestProfileIdentityTruth(t *testing.T) {
	if stableRuntimeProfileVersion != "v18.5.0" {
		t.Fatalf("stable migration source version = %q, want v18.5.0", stableRuntimeProfileVersion)
	}
	if v18TestRuntimeProfileVersion != "v18.5.1" {
		t.Fatalf("TEST target version = %q, want v18.5.1", v18TestRuntimeProfileVersion)
	}
	if v18TestRuntimeConfigDirName != "PersonalMarketTerminal-v18.5.1-TEST" {
		t.Fatalf("TEST config directory = %q", v18TestRuntimeConfigDirName)
	}
	if hostedRuntimeConfigDirName != "PersonalMarketTerminal-v18.5.1-HOSTED" {
		t.Fatalf("HOSTED config directory = %q", hostedRuntimeConfigDirName)
	}
	if v18TestProfileMigrationMarker != ".v18.5.1-test-profile-migration.json" {
		t.Fatalf("migration marker = %q", v18TestProfileMigrationMarker)
	}
}

func TestV1851TestProfileClonesV1850StableWithoutMutatingSource(t *testing.T) {
	base := t.TempDir()
	stable := filepath.Join(base, stableRuntimeConfigDirName)
	if err := os.MkdirAll(filepath.Join(stable, "nested"), 0700); err != nil {
		t.Fatal(err)
	}
	stableSentinel := filepath.Join(stable, "nested", "settings.json")
	original := []byte("{\"provider\":\"test-sentinel\"}\n")
	if err := os.WriteFile(stableSentinel, original, 0600); err != nil {
		t.Fatal(err)
	}
	// Transient runtime files must not be cloned into the isolated TEST profile.
	if err := os.WriteFile(filepath.Join(stable, "De-Pulse.log"), []byte("transient\n"), 0600); err != nil {
		t.Fatal(err)
	}

	target, err := prepareV18TestConfig(base)
	if err != nil {
		t.Fatalf("prepareV18TestConfig: %v", err)
	}
	if got, want := filepath.Base(target), v18TestRuntimeConfigDirName; got != want {
		t.Fatalf("target directory = %q, want %q", got, want)
	}
	cloned, err := os.ReadFile(filepath.Join(target, "nested", "settings.json"))
	if err != nil {
		t.Fatalf("read cloned sentinel: %v", err)
	}
	if string(cloned) != string(original) {
		t.Fatalf("cloned sentinel changed: %q", cloned)
	}
	if _, err := os.Stat(filepath.Join(target, "De-Pulse.log")); !os.IsNotExist(err) {
		t.Fatalf("transient log should not be cloned; stat err=%v", err)
	}

	markerPath := filepath.Join(target, v18TestProfileMigrationMarker)
	markerBytes, err := os.ReadFile(markerPath)
	if err != nil {
		t.Fatalf("read migration marker: %v", err)
	}
	var marker map[string]any
	if err := json.Unmarshal(markerBytes, &marker); err != nil {
		t.Fatalf("decode migration marker: %v", err)
	}
	checks := map[string]string{
		"source":        stableRuntimeConfigDirName,
		"sourceVersion": "v18.5.0",
		"target":        v18TestRuntimeConfigDirName,
		"targetVersion": "v18.5.1",
	}
	for key, want := range checks {
		if got, _ := marker[key].(string); got != want {
			t.Fatalf("marker %s = %q, want %q", key, got, want)
		}
	}
	if got, _ := marker["migratedAt"].(string); got == "" {
		t.Fatal("migration marker missing migratedAt")
	}

	// A second launch must reuse the already-isolated profile rather than reclone
	// or rewrite Stable state.
	second, err := prepareV18TestConfig(base)
	if err != nil {
		t.Fatalf("second prepareV18TestConfig: %v", err)
	}
	if second != target {
		t.Fatalf("second target = %q, want %q", second, target)
	}
	stableAfter, err := os.ReadFile(stableSentinel)
	if err != nil {
		t.Fatalf("read Stable sentinel after migration: %v", err)
	}
	if string(stableAfter) != string(original) {
		t.Fatalf("Stable source was modified: %q", stableAfter)
	}
}
