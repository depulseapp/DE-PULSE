package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	stableRuntimeConfigDirName     = "PersonalMarketTerminal"
	stableRuntimeProfileVersion   = "v18.5.0"
	v18TestRuntimeProfileVersion  = "v18.5.1"
	v18TestRuntimeConfigDirName    = "PersonalMarketTerminal-v18.5.1-TEST"
	hostedRuntimeConfigDirName     = "PersonalMarketTerminal-v18.5.1-HOSTED"
	v18TestProfileMigrationMarker = ".v18.5.1-test-profile-migration.json"
)

func resolveV18RuntimeConfig(base string) (string, error) {
	if isHostedRuntime() {
		target := strings.TrimSpace(os.Getenv(hostedConfigDirEnv))
		if target == "" {
			target = filepath.Join(base, hostedRuntimeConfigDirName)
		}
		if err := os.MkdirAll(target, 0700); err != nil {
			return "", err
		}
		return target, nil
	}
	if releaseChannel == "STABLE" {
		target := filepath.Join(base, stableRuntimeConfigDirName)
		if err := os.MkdirAll(target, 0700); err != nil {
			return "", err
		}
		return target, nil
	}
	return prepareV18TestConfig(base)
}

// prepareV18TestConfig creates an isolated v18.5.1 TEST profile. On first launch it
// clones the authoritative v18.5.0 Stable profile so settings, watchlists, API keys
// and persistent intelligence carry forward without allowing the TEST app to write
// into Stable's runtime directory.
func prepareV18TestConfig(base string) (string, error) {
	target := filepath.Join(base, v18TestRuntimeConfigDirName)
	if st, err := os.Stat(target); err == nil {
		if !st.IsDir() {
			return "", fmt.Errorf("%s TEST config path is not a directory: %s", v18TestRuntimeProfileVersion, target)
		}
		return target, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}

	source := filepath.Join(base, stableRuntimeConfigDirName)
	if st, err := os.Stat(source); err == nil && st.IsDir() {
		if stableInstanceIsLive(source) {
			return "", fmt.Errorf("%s TEST first-run migration requires DE.PULSE %s Stable to be closed once; Stable data was not modified", v18TestRuntimeProfileVersion, stableRuntimeProfileVersion)
		}
		tmp := target + ".migrating"
		_ = os.RemoveAll(tmp)
		if err := os.MkdirAll(tmp, 0700); err != nil {
			return "", err
		}
		if err := cloneStableProfile(source, tmp); err != nil {
			_ = os.RemoveAll(tmp)
			return "", fmt.Errorf("clone Stable profile: %w", err)
		}
		marker, _ := json.MarshalIndent(map[string]any{
			"source":        stableRuntimeConfigDirName,
			"sourceVersion": stableRuntimeProfileVersion,
			"target":        v18TestRuntimeConfigDirName,
			"targetVersion": v18TestRuntimeProfileVersion,
			"migratedAt":    time.Now().UTC().Format(time.RFC3339Nano),
		}, "", "  ")
		if err := os.WriteFile(filepath.Join(tmp, v18TestProfileMigrationMarker), append(marker, '\n'), 0600); err != nil {
			_ = os.RemoveAll(tmp)
			return "", err
		}
		if err := os.Rename(tmp, target); err != nil {
			_ = os.RemoveAll(tmp)
			return "", err
		}
		return target, nil
	}

	if err := os.MkdirAll(target, 0700); err != nil {
		return "", err
	}
	return target, nil
}

func stableInstanceIsLive(configDir string) bool {
	data, err := os.ReadFile(instancePath(configDir))
	if err != nil {
		return false
	}
	var in instanceInfo
	if json.Unmarshal(data, &in) != nil || strings.TrimSpace(in.URL) == "" {
		return false
	}
	client := http.Client{Timeout: 600 * time.Millisecond}
	resp, err := client.Get(strings.TrimRight(in.URL, "/") + "/api/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func cloneStableProfile(source, target string) error {
	return filepath.WalkDir(source, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		base := filepath.Base(path)
		if base == "instance.json" || base == "De-Pulse.log" || base == "native-window.log" || base == ".DS_Store" {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		dest := filepath.Join(target, rel)
		if d.IsDir() {
			return os.MkdirAll(dest, 0700)
		}
		if !d.Type().IsRegular() {
			return nil
		}
		return copyPrivateFile(path, dest)
	})
}

func copyPrivateFile(source, target string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
		return err
	}
	out, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}
