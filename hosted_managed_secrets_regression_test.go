package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func writeManagedSecretFixture(t *testing.T, dir, name, value string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(value+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestHostedManagedSecretsResolveRotateRollbackAndRevoke(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	dir := t.TempDir()
	t.Setenv(hostedManagedSecretsDirEnv, dir)
	writeManagedSecretFixture(t, dir, "finnhub", "fixture-old")
	writeManagedSecretFixture(t, dir, "alpaca-key", "fixture-ak")

	app := &Application{}
	health := app.refreshHostedManagedSecretsLocked()
	if health.ContractVersion != hostedManagedSecretsContractVersion || health.Status != "ready" {
		t.Fatalf("unexpected health: %+v", health)
	}
	if app.secrets.Finnhub != "fixture-old" || app.secrets.AlpacaKey != "fixture-ak" {
		t.Fatalf("managed secret resolution failed")
	}

	writeManagedSecretFixture(t, dir, "finnhub", "fixture-new")
	app.refreshHostedManagedSecretsLocked()
	if app.secrets.Finnhub != "fixture-new" {
		t.Fatalf("rotation was not observed")
	}

	writeManagedSecretFixture(t, dir, "finnhub", "fixture-old")
	app.refreshHostedManagedSecretsLocked()
	if app.secrets.Finnhub != "fixture-old" {
		t.Fatalf("rollback was not observed")
	}

	if err := os.Remove(filepath.Join(dir, "finnhub")); err != nil {
		t.Fatal(err)
	}
	app.refreshHostedManagedSecretsLocked()
	if app.secrets.Finnhub != "" {
		t.Fatalf("revoked/missing managed secret must fail closed")
	}
}

func TestHostedSaveDoesNotPersistManagedSecrets(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	secretDir := t.TempDir()
	t.Setenv(hostedManagedSecretsDirEnv, secretDir)
	writeManagedSecretFixture(t, secretDir, "groq", "fixture-managed")

	configDir := t.TempDir()
	app := &Application{configDir: configDir, state: defaultState()}
	app.refreshHostedManagedSecretsLocked()
	if err := app.saveLocked(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(app.secretsPath()); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("hosted save must not create secrets.json, err=%v", err)
	}
	if hint := keyHint(app.secrets.Groq); hint != "" {
		t.Fatalf("hosted public key hints must remain empty")
	}
}

func TestDesktopSecretPersistenceRemainsCompatible(t *testing.T) {
	t.Setenv(runtimeModeEnv, "desktop")
	configDir := t.TempDir()
	app := &Application{configDir: configDir, state: defaultState(), secrets: Secrets{Finnhub: "fixture-desktop-secret"}}
	if err := app.saveLocked(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(app.secretsPath())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "fixture-desktop-secret") {
		t.Fatalf("desktop local secret persistence regressed")
	}
}

func TestHostedManagedSecretHTTPBoundaryRejectsRawCredentialMutation(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})
	handler := hostedManagedSecretBoundary(next)

	req := httptest.NewRequest(http.MethodPost, "/api/settings/save", strings.NewReader("{\"settings\":{},\"finnhubKey\":\"fixture-inline\"}"))
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusConflict || called {
		t.Fatalf("inline hosted credential must be rejected, status=%d called=%v", resp.Code, called)
	}

	called = false
	req = httptest.NewRequest(http.MethodPost, "/api/settings/save", strings.NewReader("{\"settings\":{}}"))
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusNoContent || !called {
		t.Fatalf("non-secret hosted settings must pass through, status=%d called=%v", resp.Code, called)
	}

	called = false
	req = httptest.NewRequest(http.MethodPost, "/api/settings/clear-secret", strings.NewReader("{\"name\":\"finnhub\"}"))
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusConflict || called {
		t.Fatalf("hosted clear-secret must be managed externally")
	}
}

func TestHostedManagedSecretAzureAndCSIContracts(t *testing.T) {
	terraform, err := os.ReadFile("governance/hosted-infrastructure/azure/managed_secrets.tf")
	if err != nil {
		t.Fatal(err)
	}
	tf := string(terraform)
	for _, required := range []string{"azurerm_key_vault", "enable_rbac_authorization", "Key Vault Secrets User", "azurerm_private_endpoint", "purge_protection_enabled"} {
		if !strings.Contains(tf, required) {
			t.Fatalf("managed-secret terraform missing %q", required)
		}
	}
	if strings.Contains(tf, "azurerm_key_vault_secret") {
		t.Fatalf("terraform must never custody provider secret values in state")
	}

	cmd := exec.Command("python3", "tools/hosted/render_managed_secrets.py", "--environment", "dev", "--key-vault-name", "kvdepulsedev1234", "--tenant-id", "11111111-1111-1111-1111-111111111111", "--workload-identity-client-id", "22222222-2222-2222-2222-222222222222")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("managed-secret renderer failed: %v\n%s", err, out)
	}
	rendered := string(out)
	for _, required := range []string{"kind: SecretProviderClass", "provider: azure", "clientID:", "keyvaultName:", "objectAlias: finnhub", "objectAlias: alpaca-secret", "objectVersion: \"\""} {
		if !strings.Contains(rendered, required) {
			t.Fatalf("managed-secret renderer missing %q", required)
		}
	}
}
