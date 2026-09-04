package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProviderCredentialContractsAreUniqueAndMarketDataMetadataIsGoverned(t *testing.T) {
	seenProviders := map[string]bool{}
	for _, contract := range providerCredentialContracts() {
		if strings.TrimSpace(contract.ProviderID) == "" || strings.TrimSpace(contract.ProviderName) == "" {
			t.Fatalf("credential contract missing provider identity: %+v", contract)
		}
		if seenProviders[contract.ProviderID] {
			t.Fatalf("duplicate provider credential contract: %s", contract.ProviderID)
		}
		seenProviders[contract.ProviderID] = true
		seenFields := map[string]bool{}
		for _, field := range contract.Fields {
			if strings.TrimSpace(field.FieldID) == "" || strings.TrimSpace(field.SecretReference) == "" {
				t.Fatalf("credential field missing identity/reference for %s: %+v", contract.ProviderName, field)
			}
			if seenFields[field.FieldID] {
				t.Fatalf("duplicate credential field %s for %s", field.FieldID, contract.ProviderName)
			}
			seenFields[field.FieldID] = true
		}
	}

	marketData, ok := providerCredentialContract(marketDataProviderName)
	if !ok {
		t.Fatal("Market Data credential metadata missing")
	}
	if marketData.Lifecycle != "SHADOW" || marketData.Transport != "Bearer" || marketData.TestProvider != "marketdata" {
		t.Fatalf("Market Data credential metadata drift: %+v", marketData)
	}
	if marketData.SettingsSection != "global-data" || marketData.SettingsOrder <= 0 {
		t.Fatalf("Market Data Settings placement metadata drift: %+v", marketData)
	}
	if len(marketData.Fields) != 1 || marketData.Fields[0].FieldID != "token" || marketData.Fields[0].EnvironmentFallback != marketDataTokenEnv {
		t.Fatalf("Market Data token contract drift: %+v", marketData.Fields)
	}
	if marketData.Fields[0].SecretReference != "secrets.marketData" {
		t.Fatalf("Market Data must point at the canonical Secrets owner, got %q", marketData.Fields[0].SecretReference)
	}
	if marketData.Fields[0].stored == nil || marketData.Fields[0].replace == nil {
		t.Fatal("Market Data credential contract must use the canonical Secrets.MarketData slot")
	}
}

func TestProviderCredentialMutationPreserveReplaceAndClearUseCanonicalSecrets(t *testing.T) {
	secrets := Secrets{Finnhub: "old-finnhub", AlpacaKey: "old-key", AlpacaSecret: "old-secret"}
	if err := applyProviderCredentialMutation(&secrets, ProviderCredentialMutation{Provider: "finnhub", FieldID: "api-key", Action: providerCredentialPreserve}); err != nil {
		t.Fatal(err)
	}
	if secrets.Finnhub != "old-finnhub" {
		t.Fatalf("preserve mutated canonical secret: %q", secrets.Finnhub)
	}
	if err := applyProviderCredentialMutation(&secrets, ProviderCredentialMutation{Provider: "finnhub", FieldID: "api-key", Action: providerCredentialReplace, Value: "  new-finnhub\n"}); err != nil {
		t.Fatal(err)
	}
	if secrets.Finnhub != "new-finnhub" {
		t.Fatalf("replace did not normalize/write canonical secret: %q", secrets.Finnhub)
	}
	if err := applyProviderCredentialMutation(&secrets, ProviderCredentialMutation{Provider: "alpaca", FieldID: "key-id", Action: providerCredentialClear}); err != nil {
		t.Fatal(err)
	}
	if secrets.AlpacaKey != "" || secrets.AlpacaSecret != "old-secret" {
		t.Fatalf("field-scoped clear disturbed unrelated Alpaca secret: %+v", secrets)
	}
}

func TestProviderCredentialBlankReplaceFailsClosedWithoutErasingSecret(t *testing.T) {
	secrets := Secrets{Finnhub: "keep-me"}
	err := applyProviderCredentialMutation(&secrets, ProviderCredentialMutation{Provider: "finnhub", FieldID: "api-key", Action: providerCredentialReplace, Value: " \n\t "})
	if err == nil {
		t.Fatal("blank replacement must fail closed")
	}
	if secrets.Finnhub != "keep-me" {
		t.Fatalf("blank replacement erased existing secret: %q", secrets.Finnhub)
	}
}

func TestProviderCredentialPublicProjectionNeverContainsRawSecret(t *testing.T) {
	t.Setenv(runtimeModeEnv, "desktop")
	raw := "super-secret-provider-token"
	states := providerCredentialStateSnapshot(defaultState().Settings, Secrets{Finnhub: raw, AlpacaKey: "alpaca-key-secret", AlpacaSecret: "alpaca-secret-secret", MarketData: "market-data-secret"})
	data, err := json.Marshal(states)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	for _, secret := range []string{raw, "alpaca-key-secret", "alpaca-secret-secret", "market-data-secret"} {
		if strings.Contains(text, secret) {
			t.Fatalf("raw secret leaked through credential projection: %s", text)
		}
	}
	if !strings.Contains(text, "stored") || !strings.Contains(text, "secretReference") {
		t.Fatalf("redacted configured/source metadata missing: %s", text)
	}
}

func TestMarketDataEnvironmentFallbackIsDesktopHeadlessOnly(t *testing.T) {
	t.Setenv(marketDataTokenEnv, "fixture-marketdata-token")
	t.Setenv(runtimeModeEnv, "desktop")
	contract, ok := providerCredentialContract(marketDataProviderName)
	if !ok {
		t.Fatal("Market Data credential contract missing")
	}
	state := providerCredentialCardState(contract, defaultState().Settings, Secrets{})
	if !state.Configured || state.Configuration != "CONFIGURED" || len(state.Fields) != 1 || state.Fields[0].Source != providerCredentialSourceEnvironment {
		t.Fatalf("desktop MARKETDATA_TOKEN fallback not projected truthfully: %+v", state)
	}
	data, _ := json.Marshal(state)
	if strings.Contains(string(data), "fixture-marketdata-token") {
		t.Fatalf("environment token leaked through redacted state: %s", data)
	}

	t.Setenv(runtimeModeEnv, "hosted")
	hosted := providerCredentialCardState(contract, defaultState().Settings, Secrets{})
	if hosted.Configured || hosted.Fields[0].Source != providerCredentialSourceUnconfigured {
		t.Fatalf("hosted runtime must not accept ordinary MARKETDATA_TOKEN fallback: %+v", hosted)
	}
}

func TestProviderCredentialHostedStoredValueIsManagedMountedAndRedacted(t *testing.T) {
	t.Setenv(runtimeModeEnv, "hosted")
	contract, _ := providerCredentialContract("finnhub")
	state := providerCredentialCardState(contract, defaultState().Settings, Secrets{Finnhub: "fixture-mounted"})
	if !state.Configured || len(state.Fields) != 1 || state.Fields[0].Source != providerCredentialSourceManagedMounted || state.Fields[0].Hint != "" {
		t.Fatalf("hosted secret projection must be managed-mounted and hint-free: %+v", state)
	}
	data, _ := json.Marshal(state)
	if strings.Contains(string(data), "fixture-mounted") {
		t.Fatalf("hosted mounted secret leaked through projection: %s", data)
	}
}

func TestMarketDataPersistentMutationUsesCanonicalSlotAndRedactedProjection(t *testing.T) {
	t.Setenv(runtimeModeEnv, "desktop")
	secrets := Secrets{MarketData: "old-marketdata-token"}
	if err := applyProviderCredentialMutation(&secrets, ProviderCredentialMutation{Provider: marketDataProviderName, FieldID: "token", Action: providerCredentialPreserve}); err != nil {
		t.Fatal(err)
	}
	if secrets.MarketData != "old-marketdata-token" {
		t.Fatalf("Market Data preserve changed canonical slot: %q", secrets.MarketData)
	}
	if err := applyProviderCredentialMutation(&secrets, ProviderCredentialMutation{Provider: marketDataProviderName, FieldID: "token", Action: providerCredentialReplace, Value: "  new-marketdata-token\n"}); err != nil {
		t.Fatal(err)
	}
	if secrets.MarketData != "new-marketdata-token" {
		t.Fatalf("Market Data replace did not write canonical slot: %q", secrets.MarketData)
	}
	contract, _ := providerCredentialContract(marketDataProviderName)
	state := providerCredentialCardState(contract, defaultState().Settings, secrets)
	if !state.Configured || len(state.Fields) != 1 || state.Fields[0].Source != providerCredentialSourceStored {
		t.Fatalf("Market Data stored credential not projected truthfully: %+v", state)
	}
	data, _ := json.Marshal(state)
	if strings.Contains(string(data), secrets.MarketData) {
		t.Fatalf("Market Data raw token leaked through credential projection: %s", data)
	}
	if err := applyProviderCredentialMutation(&secrets, ProviderCredentialMutation{Provider: marketDataProviderName, FieldID: "token", Action: providerCredentialClear}); err != nil {
		t.Fatal(err)
	}
	if secrets.MarketData != "" {
		t.Fatalf("Market Data clear did not clear canonical slot: %q", secrets.MarketData)
	}
}

func localProviderCredentialRequest(method, path, body, sessionKey string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.AddCookie(&http.Cookie{Name: "pmt_session", Value: sessionKey})
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

func TestProviderCredentialRoutesRequireCanonicalAdminBoundaryAndNeverLeak(t *testing.T) {
	t.Setenv(runtimeModeEnv, "desktop")
	app := newTestApplication(t)
	handler := app.routes()

	blocked := httptest.NewRecorder()
	handler.ServeHTTP(blocked, httptest.NewRequest(http.MethodGet, providerCredentialStatePath, nil))
	if blocked.Code != http.StatusForbidden {
		t.Fatalf("credential metadata escaped local authentication boundary: code=%d body=%s", blocked.Code, blocked.Body.String())
	}

	stateResponse := httptest.NewRecorder()
	handler.ServeHTTP(stateResponse, localProviderCredentialRequest(http.MethodGet, providerCredentialStatePath, "", app.sessionKey))
	if stateResponse.Code != http.StatusOK || !strings.Contains(stateResponse.Body.String(), marketDataTokenEnv) {
		t.Fatalf("authenticated redacted metadata unavailable: code=%d body=%s", stateResponse.Code, stateResponse.Body.String())
	}

	raw := "route-fixture-marketdata-token"
	mutation := `{"provider":"market-data","fieldId":"token","action":"REPLACE","value":"` + raw + `"}`
	mutated := httptest.NewRecorder()
	handler.ServeHTTP(mutated, localProviderCredentialRequest(http.MethodPost, providerCredentialMutationPath, mutation, app.sessionKey))
	if mutated.Code != http.StatusOK || app.secrets.MarketData != raw {
		t.Fatalf("authenticated credential replacement failed: code=%d secrets=%+v body=%s", mutated.Code, app.secrets, mutated.Body.String())
	}
	if strings.Contains(mutated.Body.String(), raw) {
		t.Fatalf("credential mutation response leaked raw secret: %s", mutated.Body.String())
	}

	duplicateClear := httptest.NewRecorder()
	handler.ServeHTTP(duplicateClear, localProviderCredentialRequest(http.MethodPost, providerCredentialMutationPath, `{"provider":"market-data","fieldId":"token","action":"CLEAR"}`, app.sessionKey))
	if duplicateClear.Code != http.StatusBadRequest || app.secrets.MarketData != raw {
		t.Fatalf("generic mutation route became a competing clear owner: code=%d secret=%q body=%s", duplicateClear.Code, app.secrets.MarketData, duplicateClear.Body.String())
	}

	cleared := httptest.NewRecorder()
	handler.ServeHTTP(cleared, localProviderCredentialRequest(http.MethodPost, "/api/settings/clear-secret", `{"name":"market-data"}`, app.sessionKey))
	if cleared.Code != http.StatusOK || app.secrets.MarketData != "" {
		t.Fatalf("canonical clear-secret owner did not clear Market Data: code=%d secret=%q body=%s", cleared.Code, app.secrets.MarketData, cleared.Body.String())
	}
}

func TestMarketDataSharedProviderTestIsTruthfulAndUnknownNeverFallsThrough(t *testing.T) {
	t.Setenv(runtimeModeEnv, "desktop")
	t.Setenv(marketDataTokenEnv, "")
	app := newTestApplication(t)
	app.secrets.MarketData = "provider-test-marketdata-token"
	handler := app.routes()

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, localProviderCredentialRequest(http.MethodPost, "/api/provider/test", `{"provider":"marketdata"}`, app.sessionKey))
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"status":"pending"`) || !strings.Contains(response.Body.String(), "No network request was made") {
		t.Fatalf("Market Data test must report configured-but-unverified truth: code=%d body=%s", response.Code, response.Body.String())
	}
	if strings.Contains(response.Body.String(), app.secrets.MarketData) {
		t.Fatalf("Market Data test leaked raw token: %s", response.Body.String())
	}

	unknown := httptest.NewRecorder()
	handler.ServeHTTP(unknown, localProviderCredentialRequest(http.MethodPost, "/api/provider/test", `{"provider":"not-a-provider"}`, app.sessionKey))
	if unknown.Code != http.StatusBadRequest || !strings.Contains(unknown.Body.String(), "Unsupported provider test") {
		t.Fatalf("unknown provider test fell through instead of failing closed: code=%d body=%s", unknown.Code, unknown.Body.String())
	}
}
