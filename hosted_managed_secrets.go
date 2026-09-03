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
)

const (
	hostedManagedSecretsContractVersion = "v1"
	hostedManagedSecretsDirEnv          = "DEPULSE_HOSTED_SECRETS_DIR"
	hostedManagedSecretsDefaultDir      = "/var/run/depulse/secrets"
)

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
