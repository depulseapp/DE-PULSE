package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// providerConfigurationObservations is process-local only. Values are one-way
// fingerprints produced by the canonical registration contract; credentials
// never enter Router snapshots, persistence, telemetry, logs or diagnostics.
var providerConfigurationObservations sync.Map // map[*Engine]map[string][32]byte

func providerConfigurationSnapshot(settings Settings, secrets Secrets) map[string][32]byte {
	out := map[string][32]byte{}
	registry := adaptiveProviderRegistrySnapshot(settings, secrets)
	registered := make(map[string]bool, len(registry.Providers))
	for _, manifest := range registry.Providers {
		registered[manifest.ProviderID] = true
	}
	for _, reg := range providerRegistrations() {
		if !registered[providerKey(reg.Name)] || reg.ConfigurationFingerprint == nil {
			continue
		}
		out[reg.Name] = reg.ConfigurationFingerprint(settings, secrets)
	}
	return out
}

func entitlementConfigurationState(state string) bool {
	return state == providerCapabilityNotEntitled || state == providerCapabilityNotConfigured
}

func (e *Engine) reopenProviderEntitlementStates(changed map[string]bool, reason string) []string {
	if e == nil || len(changed) == 0 {
		return nil
	}
	now := time.Now().UnixMilli()
	persist := make([]ProviderCapabilityStateRecord, 0)
	changedProviders := make([]string, 0, len(changed))

	e.mu.Lock()
	for key, record := range e.providerCapabilityStates {
		if !changed[providerKey(record.Provider)] || !entitlementConfigurationState(record.State) {
			continue
		}
		record.State = providerCapabilityUnknown
		record.Reason = reason
		record.LastObservedAt = now
		record.RevalidateAt = 0
		record.PolicyVersion = smartRouterPolicyVersion
		e.providerCapabilityStates[key] = record
		persist = append(persist, record)

		circuitKey := providerCapabilityCircuitKey(record.Provider, record.Dataset)
		circuit := e.providerCapabilityCircuits[circuitKey]
		circuit.Failures = 0
		circuit.OpenUntil = 0
		circuit.RateLimitedUntil = 0
		circuit.LastError = ""
		e.providerCapabilityCircuits[circuitKey] = circuit
	}

	for _, reg := range providerRegistrations() {
		if !changed[providerKey(reg.Name)] {
			continue
		}
		changedProviders = append(changedProviders, reg.Name)
		globalKey := providerKey(reg.Name)
		circuit := e.providerCircuits[globalKey]
		if strings.TrimSpace(circuit.LastError) == "" {
			continue
		}
		state, _, _ := classifyProviderCapabilityFailure(&providerConfigurationError{message: circuit.LastError})
		if state == providerCapabilityNotEntitled {
			circuit.Failures = 0
			circuit.OpenUntil = 0
			circuit.RateLimitedUntil = 0
			circuit.LastError = ""
			e.providerCircuits[globalKey] = circuit
		}
	}
	e.mu.Unlock()

	for _, record := range persist {
		e.persistProviderCapabilityState(record)
	}
	return changedProviders
}

// refreshProviderConfigurationEntitlements makes credential/configuration
// changes visible to Smart Provider Router v2 without polling. It is called at
// canonical route-decision time before ranking. Only stale entitlement/config
// suppression is reopened; SUPPORTED, NOT_SUPPORTED and transient health states
// are retained because they represent different evidence.
//
// On the first route decision after process start, resettable persisted records
// are reopened once. This safely handles an offline credential change between
// runs while bounding the extra provider work to at most one fresh probe per
// affected capability after restart.
func (e *Engine) refreshProviderConfigurationEntitlements(settings Settings, secrets Secrets) []string {
	if e == nil {
		return nil
	}
	current := providerConfigurationSnapshot(settings, secrets)
	previousAny, loaded := providerConfigurationObservations.Load(e)
	previous, _ := previousAny.(map[string][32]byte)
	providerConfigurationObservations.Store(e, current)

	changed := map[string]bool{}
	for provider, fingerprint := range current {
		if !loaded || previous[provider] != fingerprint {
			changed[providerKey(provider)] = true
		}
	}
	return e.reopenProviderEntitlementStates(changed, "provider configuration changed; entitlement requires fresh evidence")
}

// forceProviderEntitlementRevalidation is the bounded same-key plan-upgrade
// path used by the explicit capability recheck action. It does not poll and it
// does not erase health evidence: only configured providers with resettable
// NOT_ENTITLED/NOT_CONFIGURED state are reopened for one fresh canonical probe.
func (e *Engine) forceProviderEntitlementRevalidation(settings Settings, secrets Secrets) []string {
	if e == nil {
		return nil
	}
	providerConfigurationObservations.Store(e, providerConfigurationSnapshot(settings, secrets))
	changed := map[string]bool{}
	registry := adaptiveProviderRegistrySnapshot(settings, secrets)
	for _, manifest := range registry.Providers {
		if !manifest.Configured {
			continue
		}
		consistent := true
		for _, route := range manifest.Routes {
			if adaptiveProviderRegistryCapabilityState(manifest.Name, route.Dataset, settings, secrets) == providerCapabilityNotConfigured {
				consistent = false
				break
			}
		}
		if consistent {
			changed[manifest.ProviderID] = true
		}
	}
	return e.reopenProviderEntitlementStates(changed, "manual capability recheck requested; entitlement requires fresh evidence")
}

type providerConfigurationError struct{ message string }

func (e *providerConfigurationError) Error() string { return e.message }

const (
	providerCredentialContractVersion = "provider-credential-v19.1.0"
	marketDataProviderName             = "Market Data"
	marketDataTokenEnv                 = "MARKETDATA_TOKEN"

	providerCredentialPreserve = "PRESERVE"
	providerCredentialReplace  = "REPLACE"
	providerCredentialClear    = "CLEAR"

	providerCredentialSourceStored         = "stored"
	providerCredentialSourceEnvironment    = "environment"
	providerCredentialSourceManagedMounted = "managed-mounted"
	providerCredentialSourceUnconfigured   = "unconfigured"

	providerCredentialStatePath    = "/api/provider-credentials"
	providerCredentialMutationPath = "/api/provider-credentials/mutate"
)

// ProviderCredentialFieldContract is metadata only. It identifies the existing
// canonical Secrets field used by one provider credential and, where governed,
// an environment fallback for CI/headless/developer operation. It is not a
// credential store and never carries secret material in public projections.
type ProviderCredentialFieldContract struct {
	FieldID             string
	Label               string
	SecretReference     string
	Required            bool
	EnvironmentFallback string
	stored               func(Secrets) string
	replace              func(*Secrets, string)
}

// ProviderCredentialCardContract is the reusable Settings contract for a
// credential-bearing provider. Routing, lifecycle admission, health, freshness,
// rights and serving-provider selection remain outside this contract.
type ProviderCredentialCardContract struct {
	ProviderID   string
	ProviderName string
	DisplayName  string
	Description  string
	Transport    string
	Lifecycle    string
	TestProvider string
	Fields       []ProviderCredentialFieldContract
}

type ProviderCredentialFieldState struct {
	FieldID             string `json:"fieldId"`
	Label               string `json:"label"`
	SecretReference     string `json:"secretReference"`
	Required            bool   `json:"required"`
	Configured          bool   `json:"configured"`
	Source              string `json:"source"`
	Hint                string `json:"hint,omitempty"`
	EnvironmentFallback string `json:"environmentFallback,omitempty"`
}

type ProviderCredentialCardState struct {
	ContractVersion string                         `json:"contractVersion"`
	ProviderID      string                         `json:"providerId"`
	ProviderName    string                         `json:"providerName"`
	DisplayName     string                         `json:"displayName"`
	Description     string                         `json:"description,omitempty"`
	Transport       string                         `json:"transport,omitempty"`
	Lifecycle       string                         `json:"lifecycle,omitempty"`
	TestProvider    string                         `json:"testProvider"`
	Configured      bool                           `json:"configured"`
	Configuration   string                         `json:"configuration"`
	Fields          []ProviderCredentialFieldState `json:"fields"`
}

type ProviderCredentialMutation struct {
	Provider string `json:"provider"`
	FieldID  string `json:"fieldId"`
	Action   string `json:"action"`
	Value    string `json:"value,omitempty"`
}

func providerCredentialContracts() []ProviderCredentialCardContract {
	field := func(id, label, ref string, required bool, get func(Secrets) string, set func(*Secrets, string)) ProviderCredentialFieldContract {
		return ProviderCredentialFieldContract{FieldID: id, Label: label, SecretReference: ref, Required: required, stored: get, replace: set}
	}
	return []ProviderCredentialCardContract{
		{ProviderID: providerKey("Finnhub"), ProviderName: "Finnhub", DisplayName: "Finnhub · Live Recovery & Intelligence", Transport: "WebSocket + REST", Lifecycle: "PRODUCTION", TestProvider: "finnhub", Fields: []ProviderCredentialFieldContract{
			field("api-key", "API Key", "secrets.finnhub", true, func(s Secrets) string { return s.Finnhub }, func(s *Secrets, v string) { s.Finnhub = v }),
		}},
		{ProviderID: providerKey(tradeInsightProviderName), ProviderName: tradeInsightProviderName, DisplayName: "TradeInsight · Corporate & Market Intelligence", Transport: "REST", Lifecycle: "PRODUCTION", TestProvider: "tradeinsight", Fields: []ProviderCredentialFieldContract{
			field("api-key", "API Key", "secrets.tradeInsight", true, func(s Secrets) string { return s.TradeInsight }, func(s *Secrets, v string) { s.TradeInsight = v }),
		}},
		{ProviderID: providerKey("Alpaca"), ProviderName: "Alpaca", DisplayName: "Alpaca · Preferred U.S. Equities & History", Transport: "IEX + Overnight", Lifecycle: "PRODUCTION", TestProvider: "alpaca", Fields: []ProviderCredentialFieldContract{
			field("key-id", "Key ID", "secrets.alpacaKey", true, func(s Secrets) string { return s.AlpacaKey }, func(s *Secrets, v string) { s.AlpacaKey = v }),
			field("secret-key", "Secret Key", "secrets.alpacaSecret", true, func(s Secrets) string { return s.AlpacaSecret }, func(s *Secrets, v string) { s.AlpacaSecret = v }),
		}},
		{ProviderID: providerKey("Groq"), ProviderName: "Groq", DisplayName: "Groq AI", Transport: "HTTPS", Lifecycle: "PRODUCTION", TestProvider: "groq", Fields: []ProviderCredentialFieldContract{
			field("api-key", "API Key", "secrets.groq", true, func(s Secrets) string { return s.Groq }, func(s *Secrets, v string) { s.Groq = v }),
		}},
		{ProviderID: providerKey("OpenRouter"), ProviderName: "OpenRouter", DisplayName: "OpenRouter", Transport: "HTTPS", Lifecycle: "PRODUCTION", TestProvider: "openrouter", Fields: []ProviderCredentialFieldContract{
			field("api-key", "API Key", "secrets.openrouter", true, func(s Secrets) string { return s.OpenRouter }, func(s *Secrets, v string) { s.OpenRouter = v }),
		}},
		{ProviderID: providerKey("Gemini"), ProviderName: "Gemini", DisplayName: "Google Gemini", Transport: "HTTPS", Lifecycle: "PRODUCTION", TestProvider: "gemini", Fields: []ProviderCredentialFieldContract{
			field("api-key", "API Key", "secrets.gemini", true, func(s Secrets) string { return s.Gemini }, func(s *Secrets, v string) { s.Gemini = v }),
		}},
		{ProviderID: providerKey("FRED"), ProviderName: "FRED", DisplayName: "FRED · Rates & Credit", Transport: "REST", Lifecycle: "PRODUCTION", TestProvider: "fred", Fields: []ProviderCredentialFieldContract{
			field("api-key", "API Key", "secrets.fred", true, func(s Secrets) string { return s.FRED }, func(s *Secrets, v string) { s.FRED = v }),
		}},
		{ProviderID: providerKey("BLS"), ProviderName: "BLS", DisplayName: "BLS · CPI & Labor", Transport: "REST", Lifecycle: "PRODUCTION", TestProvider: "bls", Fields: []ProviderCredentialFieldContract{
			field("api-key", "Registration Key", "secrets.bls", false, func(s Secrets) string { return s.BLS }, func(s *Secrets, v string) { s.BLS = v }),
		}},
		{ProviderID: providerKey("EIA"), ProviderName: "EIA", DisplayName: "EIA · Energy", Transport: "REST", Lifecycle: "PRODUCTION", TestProvider: "eia", Fields: []ProviderCredentialFieldContract{
			field("api-key", "API Key", "secrets.eia", true, func(s Secrets) string { return s.EIA }, func(s *Secrets, v string) { s.EIA = v }),
		}},
		{ProviderID: providerKey("Twelve Data"), ProviderName: "Twelve Data", DisplayName: "Twelve Data · VIX / Global / History", Transport: "REST", Lifecycle: "PRODUCTION", TestProvider: "twelvedata", Fields: []ProviderCredentialFieldContract{
			field("api-key", "API Key", "secrets.twelveData", true, func(s Secrets) string { return s.TwelveData }, func(s *Secrets, v string) { s.TwelveData = v }),
		}},
		{ProviderID: providerKey("Marketaux"), ProviderName: "Marketaux", DisplayName: "Marketaux · News Fallback", Transport: "REST", Lifecycle: "PRODUCTION", TestProvider: "marketaux", Fields: []ProviderCredentialFieldContract{
			field("api-key", "API Key", "secrets.marketaux", false, func(s Secrets) string { return s.Marketaux }, func(s *Secrets, v string) { s.Marketaux = v }),
		}},
		{
			ProviderID: providerKey(marketDataProviderName), ProviderName: marketDataProviderName,
			DisplayName: "Market Data · Validation Provider", Description: "SHADOW-first Market Data credential contract. Transport and capability truth are implemented in APR-03.",
			Transport: "Bearer", Lifecycle: "SHADOW", TestProvider: "marketdata",
			Fields: []ProviderCredentialFieldContract{{
				FieldID: "token", Label: "API Token", SecretReference: "secrets.marketData", Required: true,
				EnvironmentFallback: marketDataTokenEnv,
			}},
		},
	}
}

func providerCredentialContract(provider string) (ProviderCredentialCardContract, bool) {
	want := providerKey(provider)
	for _, contract := range providerCredentialContracts() {
		if contract.ProviderID == want || providerKey(contract.TestProvider) == want {
			return contract, true
		}
	}
	return ProviderCredentialCardContract{}, false
}

func providerCredentialEffectiveValue(field ProviderCredentialFieldContract, secrets Secrets) (string, string) {
	if field.stored != nil {
		if value := cleanCredential(field.stored(secrets)); value != "" {
			if isHostedRuntime() {
				return value, providerCredentialSourceManagedMounted
			}
			return value, providerCredentialSourceStored
		}
	}
	// Environment fallback is intentionally a desktop/headless/developer path.
	// Hosted provider credentials are server-managed by the existing mounted
	// secret lifecycle and must never be sourced from ordinary process env UX.
	if !isHostedRuntime() && strings.TrimSpace(field.EnvironmentFallback) != "" {
		if value := cleanCredential(os.Getenv(field.EnvironmentFallback)); value != "" {
			return value, providerCredentialSourceEnvironment
		}
	}
	return "", providerCredentialSourceUnconfigured
}

func providerCredentialFieldState(field ProviderCredentialFieldContract, secrets Secrets) ProviderCredentialFieldState {
	value, source := providerCredentialEffectiveValue(field, secrets)
	state := ProviderCredentialFieldState{
		FieldID: field.FieldID, Label: field.Label, SecretReference: field.SecretReference,
		Required: field.Required, Configured: value != "", Source: source,
		EnvironmentFallback: field.EnvironmentFallback,
	}
	if source == providerCredentialSourceStored {
		state.Hint = keyHint(value)
	}
	return state
}

func providerCredentialCardState(contract ProviderCredentialCardContract, settings Settings, secrets Secrets) ProviderCredentialCardState {
	state := ProviderCredentialCardState{
		ContractVersion: providerCredentialContractVersion,
		ProviderID: contract.ProviderID, ProviderName: contract.ProviderName, DisplayName: contract.DisplayName,
		Description: contract.Description, Transport: contract.Transport, Lifecycle: contract.Lifecycle, TestProvider: contract.TestProvider,
		Configuration: providerCapabilityNotConfigured,
	}
	allRequired := true
	for _, field := range contract.Fields {
		fieldState := providerCredentialFieldState(field, secrets)
		state.Fields = append(state.Fields, fieldState)
		if field.Required && !fieldState.Configured {
			allRequired = false
		}
	}
	if reg, ok := providerRegistrationLookup(contract.ProviderName); ok && reg.Configured != nil {
		state.Configured = reg.Configured(settings, secrets)
	} else {
		state.Configured = allRequired
	}
	if state.Configured {
		state.Configuration = "CONFIGURED"
	}
	return state
}

func providerCredentialStateSnapshot(settings Settings, secrets Secrets) []ProviderCredentialCardState {
	contracts := providerCredentialContracts()
	out := make([]ProviderCredentialCardState, 0, len(contracts))
	for _, contract := range contracts {
		out = append(out, providerCredentialCardState(contract, settings, secrets))
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].ProviderID < out[j].ProviderID })
	return out
}

func applyProviderCredentialMutation(secrets *Secrets, mutation ProviderCredentialMutation) error {
	if secrets == nil {
		return fmt.Errorf("credential mutation requires canonical Secrets owner")
	}
	contract, ok := providerCredentialContract(mutation.Provider)
	if !ok {
		return fmt.Errorf("unknown provider credential contract: %s", strings.TrimSpace(mutation.Provider))
	}
	var field *ProviderCredentialFieldContract
	for i := range contract.Fields {
		if strings.EqualFold(strings.TrimSpace(contract.Fields[i].FieldID), strings.TrimSpace(mutation.FieldID)) {
			field = &contract.Fields[i]
			break
		}
	}
	if field == nil {
		return fmt.Errorf("unknown credential field %s for %s", strings.TrimSpace(mutation.FieldID), contract.ProviderName)
	}
	action := strings.ToUpper(strings.TrimSpace(mutation.Action))
	switch action {
	case providerCredentialPreserve:
		return nil
	case providerCredentialReplace:
		value := cleanCredential(mutation.Value)
		if value == "" {
			return fmt.Errorf("replacement credential must be non-empty")
		}
		if field.replace == nil {
			return fmt.Errorf("%s persistent credential slot is not active yet", contract.ProviderName)
		}
		field.replace(secrets, value)
		return nil
	case providerCredentialClear:
		if field.replace == nil {
			return fmt.Errorf("%s persistent credential slot is not active yet", contract.ProviderName)
		}
		field.replace(secrets, "")
		return nil
	default:
		return fmt.Errorf("unsupported credential mutation action: %s", mutation.Action)
	}
}

// providerCredentialBoundary exposes only metadata/redacted state and canonical
// credential mutation. It intentionally sits in front of the existing routes so
// future provider cards can use one backend contract without adding a second
// Settings or secret store. Hosted mutation remains fail-closed/server-managed.
func providerCredentialBoundary(app *Application, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app == nil {
			next.ServeHTTP(w, r)
			return
		}
		switch {
		case r.URL.Path == providerCredentialStatePath && r.Method == http.MethodGet:
			app.mu.Lock()
			state := providerCredentialStateSnapshot(app.state.Settings, app.secrets)
			app.mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]any{"providers": state}); err != nil {
				http.Error(w, "unable to encode provider credential state", http.StatusInternalServerError)
			}
			return
		case r.URL.Path == providerCredentialMutationPath && r.Method == http.MethodPost:
			if isHostedRuntime() {
				http.Error(w, "hosted provider credentials are server-managed", http.StatusConflict)
				return
			}
			var mutation ProviderCredentialMutation
			decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&mutation); err != nil {
				http.Error(w, "invalid provider credential mutation", http.StatusBadRequest)
				return
			}
			app.mu.Lock()
			err := applyProviderCredentialMutation(&app.secrets, mutation)
			if err == nil {
				err = app.saveLocked()
			}
			state := providerCredentialStateSnapshot(app.state.Settings, app.secrets)
			app.mu.Unlock()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]any{"providers": state}); err != nil {
				http.Error(w, "unable to encode provider credential state", http.StatusInternalServerError)
			}
			return
		default:
			next.ServeHTTP(w, r)
		}
	})
}
