package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"depulse/internal/hostedenv"
)

const (
	runtimeModeEnv                      = "DEPULSE_RUNTIME_MODE"
	hostedListenAddrEnv                 = "DEPULSE_LISTEN_ADDR"
	hostedConfigDirEnv                  = "DEPULSE_CONFIG_DIR"
	hostedTrustProxyHeadersEnv          = "DEPULSE_TRUST_PROXY_HEADERS"
	hostedPublicOriginEnv               = "DEPULSE_PUBLIC_ORIGIN"
	hostedEnvironmentEnv                = "DEPULSE_HOSTED_ENVIRONMENT"
	hostedDesiredStateVersionEnv        = "DEPULSE_HOSTED_DESIRED_STATE_VERSION"
	hostedDesiredStateSHA256Env         = "DEPULSE_HOSTED_DESIRED_STATE_SHA256"
	hostedIsolationIDEnv                = "DEPULSE_HOSTED_ISOLATION_ID"
	hostedServiceIdentityEnv            = "DEPULSE_HOSTED_SERVICE_IDENTITY"
	hostedIngressPolicyEnv              = "DEPULSE_HOSTED_INGRESS_POLICY"
	hostedEgressPolicyEnv               = "DEPULSE_HOSTED_EGRESS_POLICY"
	hostedNetworkPolicyEnv              = "DEPULSE_HOSTED_NETWORK_POLICY"
	hostedTLSPolicyEnv                  = "DEPULSE_HOSTED_TLS_POLICY"
	hostedInternalMTLSEnv               = "DEPULSE_HOSTED_INTERNAL_MTLS"
	providerRightsEnforcementModeEnv    = "DEPULSE_PROVIDER_RIGHTS_ENFORCEMENT_MODE"
	providerRightsEnforcementPublicMode = "PUBLIC_PRODUCTION"
	hostedManagedSecretsContractVersion = "v1"
	hostedManagedSecretsDirEnv          = "DEPULSE_HOSTED_SECRETS_DIR"
	hostedManagedSecretsDefaultDir      = "/var/run/depulse/secrets"
)

func runtimeMode() string {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(runtimeModeEnv))) {
	case "hosted", "server", "web":
		return "hosted"
	default:
		return "desktop"
	}
}

func isHostedRuntime() bool { return runtimeMode() == "hosted" }

// providerRightsEnforcementActive intentionally stays independent from hosted
// runtime selection. During development and pre-public validation, provider
// rights are evaluated and surfaced as governance/audit truth without reducing
// configured provider capacity. Hard fail-closed routing, fanout, cache,
// persistence and serving activate only when PUBLIC_PRODUCTION is explicitly
// selected for a hosted runtime.
func providerRightsEnforcementActive() bool {
	return isHostedRuntime() && strings.EqualFold(
		strings.TrimSpace(os.Getenv(providerRightsEnforcementModeEnv)),
		providerRightsEnforcementPublicMode,
	)
}

func hostedListenAddress() string {
	if addr := strings.TrimSpace(os.Getenv(hostedListenAddrEnv)); addr != "" {
		return addr
	}
	if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
		if strings.Contains(port, ":") {
			return port
		}
		return ":" + port
	}
	return ":8080"
}

func trustHostedProxyHeaders() bool {
	if !isHostedRuntime() {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(os.Getenv(hostedTrustProxyHeadersEnv))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func hostedPublicOrigin() string {
	if !isHostedRuntime() {
		return ""
	}
	return strings.TrimSpace(os.Getenv(hostedPublicOriginEnv))
}

func hostedInternalMTLSEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(hostedInternalMTLSEnv))) {
	case "1", "true", "yes", "on", "required":
		return true
	default:
		return false
	}
}

func validateHostedEnvironment() error {
	if !isHostedRuntime() {
		return nil
	}
	return hostedenv.Validate(hostedenv.RuntimeDeclaration{
		Environment:         strings.TrimSpace(os.Getenv(hostedEnvironmentEnv)),
		DesiredStateVersion: strings.TrimSpace(os.Getenv(hostedDesiredStateVersionEnv)),
		DesiredStateSHA256:  strings.TrimSpace(os.Getenv(hostedDesiredStateSHA256Env)),
		IsolationID:         strings.TrimSpace(os.Getenv(hostedIsolationIDEnv)),
		ServiceIdentity:     strings.TrimSpace(os.Getenv(hostedServiceIdentityEnv)),
		IngressPolicy:       strings.TrimSpace(os.Getenv(hostedIngressPolicyEnv)),
		EgressPolicy:        strings.TrimSpace(os.Getenv(hostedEgressPolicyEnv)),
		NetworkPolicy:       strings.TrimSpace(os.Getenv(hostedNetworkPolicyEnv)),
		TLSPolicy:           strings.TrimSpace(os.Getenv(hostedTLSPolicyEnv)),
		InternalMTLS:        hostedInternalMTLSEnabled(),
		TrustedProxyHeaders: trustHostedProxyHeaders(),
		PublicOrigin:        hostedPublicOrigin(),
	})
}

type hostedManagedSecretHealth struct {
	ContractVersion string `json:"contractVersion"`
	Source          string `json:"source"`
	Status          string `json:"status"`
	Generation      string `json:"generation,omitempty"`
}

type hostedManagedSecretFile struct {
	logicalName string
	fileName    string
}

var hostedManagedSecretFiles = []hostedManagedSecretFile{
	{logicalName: "finnhub", fileName: "finnhub"},
	{logicalName: "tradeinsight", fileName: "tradeinsight"},
	{logicalName: "alpaca-key", fileName: "alpaca-key"},
	{logicalName: "alpaca-secret", fileName: "alpaca-secret"},
	{logicalName: "groq", fileName: "groq"},
	{logicalName: "openrouter", fileName: "openrouter"},
	{logicalName: "gemini", fileName: "gemini"},
	{logicalName: "fred", fileName: "fred"},
	{logicalName: "bls", fileName: "bls"},
	{logicalName: "eia", fileName: "eia"},
	{logicalName: "twelvedata", fileName: "twelvedata"},
	{logicalName: "marketaux", fileName: "marketaux"},
}

var hostedManagedSecretAudit sync.Map

func hostedManagedSecretsDir() string {
	if dir := strings.TrimSpace(os.Getenv(hostedManagedSecretsDirEnv)); dir != "" {
		return filepath.Clean(dir)
	}
	return hostedManagedSecretsDefaultDir
}

func normalizeHostedManagedSecret(data []byte) string {
	value := strings.TrimSpace(string(data))
	value = strings.ReplaceAll(value, "\r", "")
	value = strings.ReplaceAll(value, "\n", "")
	value = strings.ReplaceAll(value, "\t", "")
	return value
}

func assignHostedManagedSecret(secrets *Secrets, logicalName, value string) {
	switch logicalName {
	case "finnhub":
		secrets.Finnhub = value
	case "tradeinsight":
		secrets.TradeInsight = value
	case "alpaca-key":
		secrets.AlpacaKey = value
	case "alpaca-secret":
		secrets.AlpacaSecret = value
	case "groq":
		secrets.Groq = value
	case "openrouter":
		secrets.OpenRouter = value
	case "gemini":
		secrets.Gemini = value
	case "fred":
		secrets.FRED = value
	case "bls":
		secrets.BLS = value
	case "eia":
		secrets.EIA = value
	case "twelvedata":
		secrets.TwelveData = value
	case "marketaux":
		secrets.Marketaux = value
	}
}

func readHostedManagedSecrets() (Secrets, hostedManagedSecretHealth) {
	health := hostedManagedSecretHealth{
		ContractVersion: hostedManagedSecretsContractVersion,
		Source:          "managed-mounted",
		Status:          "ready",
	}
	dir := hostedManagedSecretsDir()
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		health.Status = "unavailable"
		return Secrets{}, health
	}

	var secrets Secrets
	var latest time.Time
	for _, spec := range hostedManagedSecretFiles {
		path := filepath.Join(dir, spec.fileName)
		data, readErr := os.ReadFile(path)
		if os.IsNotExist(readErr) {
			continue
		}
		if readErr != nil {
			health.Status = "degraded"
			continue
		}
		if value := normalizeHostedManagedSecret(data); value != "" {
			assignHostedManagedSecret(&secrets, spec.logicalName, value)
		}
		if stat, statErr := os.Stat(path); statErr == nil && stat.ModTime().After(latest) {
			latest = stat.ModTime()
		}
	}
	if latest.IsZero() {
		health.Generation = "empty"
	} else {
		health.Generation = latest.UTC().Format(time.RFC3339Nano)
	}
	return secrets, health
}

func (a *Application) refreshHostedManagedSecretsLocked() hostedManagedSecretHealth {
	secrets, health := readHostedManagedSecrets()
	a.secrets = secrets
	a.auditHostedManagedSecretHealth(health)
	return health
}

func (a *Application) auditHostedManagedSecretHealth(health hostedManagedSecretHealth) {
	marker := health.Status + "|" + health.Generation
	previous, loaded := hostedManagedSecretAudit.LoadOrStore(a, marker)
	if loaded && previous == marker {
		return
	}
	hostedManagedSecretAudit.Store(a, marker)
	log.Printf("hosted managed-secret lifecycle: contract=%s source=%s status=%s generation=%s", health.ContractVersion, health.Source, health.Status, health.Generation)
}

func (a *Application) startHostedManagedSecretRefresh() {
	if !isHostedRuntime() {
		return
	}
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			a.mu.Lock()
			a.refreshHostedManagedSecretsLocked()
			a.mu.Unlock()
		}
	}()
}

var hostedInlineSecretFields = []string{
	"finnhubKey", "tradeInsightKey", "alpacaKey", "alpacaSecret", "groqKey", "openRouterKey", "geminiKey", "fredKey", "blsKey", "eiaKey", "twelveDataKey", "marketauxKey",
}

func hostedRequestCarriesInlineSecret(r *http.Request) bool {
	if r.Body == nil {
		return false
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, (2<<20)+1))
	_ = r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(body))
	r.ContentLength = int64(len(body))
	if err != nil || len(body) > 2<<20 {
		return false
	}
	var payload map[string]any
	if json.Unmarshal(body, &payload) != nil {
		return false
	}
	for _, field := range hostedInlineSecretFields {
		if value, ok := payload[field].(string); ok && strings.TrimSpace(value) != "" {
			return true
		}
	}
	return false
}

func hostedManagedSecretBoundary(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isHostedRuntime() {
			next.ServeHTTP(w, r)
			return
		}
		if r.URL.Path == "/health/managed-secrets" {
			_, health := readHostedManagedSecrets()
			status := http.StatusOK
			if health.Status != "ready" {
				status = http.StatusServiceUnavailable
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_ = json.NewEncoder(w).Encode(health)
			return
		}
		if r.URL.Path == "/api/settings/clear-secret" {
			writeError(w, http.StatusConflict, "Hosted credentials are managed by the server secret lifecycle and cannot be cleared through product state.")
			return
		}
		if (r.URL.Path == "/api/settings/save" || r.URL.Path == "/api/provider/test") && hostedRequestCarriesInlineSecret(r) {
			writeError(w, http.StatusConflict, "Hosted credentials must be provisioned through the managed-secret lifecycle; raw credentials are not accepted by product APIs.")
			return
		}
		next.ServeHTTP(w, r)
	})
}
