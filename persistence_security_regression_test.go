package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

const t6FileFallbackHelperEnv = "DEPULSE_T6_FILE_FALLBACK_HELPER"

// TestV18T6FileFallbackSecurityPersistence makes the existing Qualified
// persistence lane execute the non-cgo local backend as well as the normal
// cgo/SQLite backend. The backend is selected by the same production build
// constraints; this test does not create a second persistence implementation.
func TestV18T6FileFallbackSecurityPersistence(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file-fallback backend is intentionally non-Windows")
	}
	if job := os.Getenv("GITHUB_JOB"); job != "" && job != "db-integration" {
		t.Skip("nested no-cgo fallback proof is owned by Qualified db-integration")
	}
	cmd := exec.Command("go", "test", "-count=1", ".", "-run", "^TestV18T6FileFallbackSecurityHelper$")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0", t6FileFallbackHelperEnv+"=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("no-cgo file-fallback security regression failed: %v\n%s", err, out)
	}
}

func TestV18T6FileFallbackSecurityHelper(t *testing.T) {
	if os.Getenv(t6FileFallbackHelperEnv) != "1" {
		t.Skip("helper executes only from the governed no-cgo fallback proof")
	}

	dir := t.TempDir()
	p := NewPersistenceManager(dir)
	defer p.Close()
	if d := p.Diagnostics(); !d.Ready || d.Backend != "file-fallback" {
		t.Fatalf("expected governed file-fallback backend, got %+v", d)
	}

	identity, err := NewIdentityService(p)
	if err != nil {
		t.Fatal(err)
	}
	bootstrapToken, principal, err := identity.bootstrapOwnerSession()
	if err != nil {
		t.Fatal(err)
	}
	password := "t6 fallback owner passphrase"
	rotatedToken, _, err := identity.setPassword(principal.UserID, password)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := identity.resolve(bootstrapToken, false); err == nil {
		t.Fatal("pre-credential bootstrap token remained valid on file-fallback")
	}
	loginToken, loggedIn, err := identity.authenticate("owner", password)
	if err != nil {
		t.Fatal(err)
	}
	if loggedIn.UserID != principal.UserID {
		t.Fatalf("file-fallback identity changed user across persistence: got=%s want=%s", loggedIn.UserID, principal.UserID)
	}

	workspaceA := v181WorkspaceWithSymbols("user-a", "NVDA")
	workspaceB := v181WorkspaceWithSymbols("user-b", "MSFT")
	if err := p.SaveUserWorkspace(context.Background(), workspaceA); err != nil {
		t.Fatal(err)
	}
	if err := p.SaveUserWorkspace(context.Background(), workspaceB); err != nil {
		t.Fatal(err)
	}
	loaded, err := p.LoadUserWorkspaces(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	byUser := map[string]UserWorkspace{}
	for _, workspace := range loaded {
		byUser[workspace.UserID] = workspace
	}
	if got := trackedSymbolsFromWatchlists(byUser["user-a"].Watchlists); !contains(got, "NVDA") || contains(got, "MSFT") {
		t.Fatalf("file-fallback user-a workspace leaked or lost data: %v", got)
	}
	if got := trackedSymbolsFromWatchlists(byUser["user-b"].Watchlists); !contains(got, "MSFT") || contains(got, "NVDA") {
		t.Fatalf("file-fallback user-b workspace leaked or lost data: %v", got)
	}

	raw, err := os.ReadFile(filepath.Join(dir, "persistent-intelligence-fallback.json"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(raw, []byte(password)) {
		t.Fatal("plaintext password found in file-fallback persistence")
	}
	for _, token := range []string{bootstrapToken, rotatedToken, loginToken} {
		if token != "" && bytes.Contains(raw, []byte(token)) {
			t.Fatal("raw session token found in file-fallback persistence")
		}
	}
}
