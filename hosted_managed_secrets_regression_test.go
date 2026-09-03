package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func writeManagedSecretFixture(t *testing.T, dir, name, value string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(value+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
}

func configureHostedManagedSecretContract(t *testing.T, aliases, generation string) {
	t.Helper()
	t.Setenv(hostedManagedSecretsRequiredEnv, aliases)
	t.Setenv(hostedManagedSecretGenerationEnv, generation)
}

func writeManagedSecretReferenceManifest(t *testing.T, environment string, references map[string]string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "managed-secret-references.json")
	payload := map[string]any{
		"schema":      "DE.PULSE-HOSTED-MANAGED-SECRET-REFERENCES-1",
		"environment": environment,
		"references":  references,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func renderManagedSecretFixture(t *testing.T, manifest string) string {
	t.Helper()
	cmd := exec.Command(
		"python3", "tools/hosted/render_managed_secrets.py",
		"--environment", "dev",
		"--key-vault-name", "kvdepulsedev1234",
		"--tenant-id", "11111111-1111-1111-1111-111111111111",
		"--workload-identity-client-id", "22222222-2222-2222-2222-222222222222",
		"--reference-manifest", manifest,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("managed-secret renderer failed: %v\n%s", err, out)
	}
	return string(out)
}

func managedSecretProviderClassName(t *testing.T, rendered string) string {
	t.Helper()
	re := regexp.MustCompile(`(?m)^  name: (depulse-managed-secrets-[0-9a-f]{12})$`)
	match := re.FindStringSubmatch(rendered)
	if len(match) != 2 {
		t.Fatalf("managed-secret provider class name missing from render:\n%s", rendered)
	}
	return match[1]
}

func TestHostedManagedSecretsResolveRotateRollbackAndRevoke(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	dir := t.TempDir()
	t.Setenv(hostedManagedSecretsDirEnv, dir)
	oldGeneration := strings.Repeat("a", 64)
	newGeneration := strings.Repeat("b", 64)
	configureHostedManagedSecretContract(t, "finnhub,alpaca-key", oldGeneration)
	writeManagedSecretFixture(t, dir, "finnhub", "fixture-old")
	writeManagedSecretFixture(t, dir, "alpaca-key", "fixture-ak")

	app := &Application{}
	health := app.refreshHostedManagedSecretsLocked()
	if health.ContractVersion != hostedManagedSecretsContractVersion || health.Status != "ready" || health.Generation != oldGeneration {
		t.Fatalf("unexpected health: %+v", health)
	}
	if app.secrets.Finnhub != "fixture-old" || app.secrets.AlpacaKey != "fixture-ak" {
		t.Fatalf("managed secret resolution failed")
	}

	writeManagedSecretFixture(t, dir, "finnhub", "fixture-new")
	t.Setenv(hostedManagedSecretGenerationEnv, newGeneration)
	health = app.refreshHostedManagedSecretsLocked()
	if health.Status != "ready" || health.Generation != newGeneration || app.secrets.Finnhub != "fixture-new" {
		t.Fatalf("rotation was not observed: health=%+v", health)
	}

	writeManagedSecretFixture(t, dir, "finnhub", "fixture-old")
	t.Setenv(hostedManagedSecretGenerationEnv, oldGeneration)
	health = app.refreshHostedManagedSecretsLocked()
	if health.Status != "ready" || health.Generation != oldGeneration || app.secrets.Finnhub != "fixture-old" {
		t.Fatalf("rollback was not observed: health=%+v", health)
	}

	if err := os.Remove(filepath.Join(dir, "finnhub")); err != nil {
		t.Fatal(err)
	}
	health = app.refreshHostedManagedSecretsLocked()
	if health.Status == "ready" || app.secrets.Finnhub != "" || app.secrets.AlpacaKey != "" {
		t.Fatalf("incomplete required projection must revoke the whole staged secret set: health=%+v", health)
	}
}

func TestHostedManagedSecretsRequireExplicitProjectionContract(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	dir := t.TempDir()
	t.Setenv(hostedManagedSecretsDirEnv, dir)
	writeManagedSecretFixture(t, dir, "finnhub", "fixture-value")

	if _, health := readHostedManagedSecrets(); health.Status == "ready" {
		t.Fatalf("hosted secret health must fail closed without required aliases and generation")
	}
	configureHostedManagedSecretContract(t, "not-a-provider", strings.Repeat("a", 64))
	if _, health := readHostedManagedSecrets(); health.Status == "ready" {
		t.Fatalf("unknown required alias must fail closed")
	}
	configureHostedManagedSecretContract(t, "finnhub", "not-a-generation")
	if _, health := readHostedManagedSecrets(); health.Status == "ready" {
		t.Fatalf("invalid generation must fail closed")
	}
	configureHostedManagedSecretContract(t, "finnhub", strings.Repeat("a", 64))
	writeManagedSecretFixture(t, dir, "finnhub", "")
	if _, health := readHostedManagedSecrets(); health.Status == "ready" {
		t.Fatalf("empty required secret projection must fail closed")
	}
}

func TestHostedSaveDoesNotPersistManagedSecrets(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	secretDir := t.TempDir()
	t.Setenv(hostedManagedSecretsDirEnv, secretDir)
	configureHostedManagedSecretContract(t, "groq", strings.Repeat("c", 64))
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

func TestHostedManagedSecretHTTPBoundaryRejectsRawCredentialMutationAndGatesReadiness(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	dir := t.TempDir()
	t.Setenv(hostedManagedSecretsDirEnv, dir)
	generation := strings.Repeat("d", 64)
	configureHostedManagedSecretContract(t, "finnhub", generation)
	writeManagedSecretFixture(t, dir, "finnhub", "fixture-mounted-secret")

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

	called = false
	req = httptest.NewRequest(http.MethodGet, "/health/managed-secrets", nil)
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	body := resp.Body.String()
	if resp.Code != http.StatusOK || called || !strings.Contains(body, generation) || strings.Contains(body, "fixture-mounted-secret") {
		t.Fatalf("managed-secret health must be sanitized and ready: status=%d called=%v body=%s", resp.Code, called, body)
	}

	if err := os.Remove(filepath.Join(dir, "finnhub")); err != nil {
		t.Fatal(err)
	}
	called = false
	req = httptest.NewRequest(http.MethodGet, "/api/ready", nil)
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusServiceUnavailable || called || !strings.Contains(resp.Body.String(), "managedSecrets") {
		t.Fatalf("hosted readiness must fail before downstream readiness when managed secrets are incomplete: status=%d called=%v body=%s", resp.Code, called, resp.Body.String())
	}

	writeManagedSecretFixture(t, dir, "finnhub", "fixture-mounted-secret")
	called = false
	req = httptest.NewRequest(http.MethodGet, "/api/ready", nil)
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusNoContent || !called {
		t.Fatalf("healthy managed-secret projection must continue to canonical readiness, status=%d called=%v", resp.Code, called)
	}
}

func TestHostedManagedSecretAzureCSIAndWorkloadContracts(t *testing.T) {
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

	oldManifest := writeManagedSecretReferenceManifest(t, "dev", map[string]string{
		"finnhub":       strings.Repeat("1", 32),
		"alpaca-secret": strings.Repeat("2", 32),
	})
	oldRendered := renderManagedSecretFixture(t, oldManifest)
	oldClass := managedSecretProviderClassName(t, oldRendered)
	for _, required := range []string{
		"kind: SecretProviderClass", "provider: azure", "clientID:", "keyvaultName:",
		"objectAlias: finnhub", "objectAlias: alpaca-secret",
		"objectVersion: \"11111111111111111111111111111111\"",
		"objectVersion: \"22222222222222222222222222222222\"",
	} {
		if !strings.Contains(oldRendered, required) {
			t.Fatalf("managed-secret renderer missing %q", required)
		}
	}
	if strings.Contains(oldRendered, "objectVersion: \"\"") {
		t.Fatalf("managed-secret renderer must never resolve mutable latest versions")
	}

	image := "ghcr.io/depulseapp/de-pulse@sha256:" + strings.Repeat("a", 64)
	cmd := exec.Command(
		"python3", "tools/hosted/render_kubernetes_trust.py",
		"--environment", "dev",
		"--mesh-profile", "aks-managed",
		"--istio-revision", "asm-1-30",
		"--workload-identity-client-id", "22222222-2222-2222-2222-222222222222",
		"--workload-image", image,
		"--public-origin", "https://dev.depulse.invalid",
		"--managed-secret-reference-manifest", oldManifest,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hosted workload renderer failed: %v\n%s", err, out)
	}
	workload := string(out)
	for _, required := range []string{
		"kind: Deployment", "app.kubernetes.io/component: hosted-web", "azure.workload.identity/use: \"true\"",
		"driver: secrets-store.csi.k8s.io", "mountPath: /var/run/depulse/secrets",
		"name: DEPULSE_HOSTED_REQUIRED_SECRETS", "value: \"finnhub,alpaca-secret\"",
		"name: DEPULSE_HOSTED_SECRET_GENERATION", "path: /api/ready", oldClass,
	} {
		if !strings.Contains(workload, required) {
			t.Fatalf("hosted workload contract missing %q", required)
		}
	}

	newManifest := writeManagedSecretReferenceManifest(t, "dev", map[string]string{
		"finnhub":       strings.Repeat("3", 32),
		"alpaca-secret": strings.Repeat("2", 32),
	})
	newClass := managedSecretProviderClassName(t, renderManagedSecretFixture(t, newManifest))
	if newClass == oldClass {
		t.Fatalf("versioned cutover must produce a new immutable SecretProviderClass identity")
	}
	rollbackClass := managedSecretProviderClassName(t, renderManagedSecretFixture(t, oldManifest))
	if rollbackClass != oldClass {
		t.Fatalf("rollback to the prior exact references must restore the prior provider-class identity")
	}

	invalidManifest := writeManagedSecretReferenceManifest(t, "dev", map[string]string{"finnhub": ""})
	invalid := exec.Command(
		"python3", "tools/hosted/render_managed_secrets.py",
		"--environment", "dev",
		"--key-vault-name", "kvdepulsedev1234",
		"--tenant-id", "11111111-1111-1111-1111-111111111111",
		"--workload-identity-client-id", "22222222-2222-2222-2222-222222222222",
		"--reference-manifest", invalidManifest,
	)
	if output, err := invalid.CombinedOutput(); err == nil {
		t.Fatalf("blank/mutable version reference must fail closed: %s", output)
	}

	rawManifest := filepath.Join(t.TempDir(), "raw-values-rejected.json")
	rawPayload := `{"schema":"DE.PULSE-HOSTED-MANAGED-SECRET-REFERENCES-1","environment":"dev","references":{"finnhub":"11111111111111111111111111111111"},"rawValues":{"finnhub":"fixture-secret"}}`
	if err := os.WriteFile(rawManifest, []byte(rawPayload), 0600); err != nil {
		t.Fatal(err)
	}
	raw := exec.Command(
		"python3", "tools/hosted/render_managed_secrets.py",
		"--environment", "dev",
		"--key-vault-name", "kvdepulsedev1234",
		"--tenant-id", "11111111-1111-1111-1111-111111111111",
		"--workload-identity-client-id", "22222222-2222-2222-2222-222222222222",
		"--reference-manifest", rawManifest,
	)
	if output, err := raw.CombinedOutput(); err == nil {
		t.Fatalf("managed-secret reference contract must reject raw-value fields: %s", output)
	}
}
