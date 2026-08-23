package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	aiContextPolicyVersion      = "v18.6-ai-context-v1"
	aiSchemaPolicyVersion       = "v18.6-ai-schema-v1"
	aiCachePolicyVersion        = "v18.6-ai-cache-v1"
	aiEvalPolicyVersion         = "v18.6-ai-eval-v1"
	aiRightsEgressPolicyVersion = "v18.6.0-ai-egress-rights-1"
	aiMaxPromptBytes            = 24000
	aiMaxPromptTokenUpperBound  = 24000
	aiMaxTaskBytes              = 4000
	aiMaxEvidenceForPrompt      = 32
	aiCacheMaxEntries           = 256
)

const aiInferenceCacheTTL = 15 * time.Minute

//go:embed provider_dataset_ai_rights_registry.json
var providerDatasetAIRightsRegistryJSON []byte

type AIContextDiagnostics struct {
	PolicyVersion      string `json:"policyVersion"`
	OriginalEvidence   int    `json:"originalEvidence"`
	SentEvidence       int    `json:"sentEvidence"`
	Compacted          bool   `json:"compacted"`
	PromptBytes        int    `json:"promptBytes"`
	TokenUpperBound    int    `json:"tokenUpperBound"`
	MaxBytes           int    `json:"maxBytes"`
	MaxTokenUpperBound int    `json:"maxTokenUpperBound"`
}

type AIDatasetRightsRecord struct {
	Provider       string `json:"provider"`
	Dataset        string `json:"dataset"`
	CommercialUse  string `json:"commercialUse"`
	Redistribution string `json:"redistribution"`
	AIUse          string `json:"aiUse"`
	EvidenceBound  bool   `json:"evidenceBound"`
	Decision       string `json:"decision"`
}

type aiDatasetRightsRegistryDocument struct {
	Schema          string                  `json:"schema"`
	PolicyVersion   string                  `json:"policyVersion"`
	DefaultDecision string                  `json:"defaultDecision"`
	DefaultReason   string                  `json:"defaultReason"`
	Records         []AIDatasetRightsRecord `json:"records"`
}

type aiRightsCoordinate struct {
	Provider string
	Dataset  string
}

type AIEgressRightsDiagnostics struct {
	PolicyVersion string `json:"policyVersion"`
	AllowedItems  int    `json:"allowedItems"`
	WithheldItems int    `json:"withheldItems"`
	UnknownItems  int    `json:"unknownItems"`
	DeniedItems   int    `json:"deniedItems"`
	Summary       string `json:"summary"`
}

type AIInferenceTelemetry struct {
	PolicyVersion       string  `json:"policyVersion"`
	Requests            uint64  `json:"requests"`
	CacheHits           uint64  `json:"cacheHits"`
	SchemaFailures      uint64  `json:"schemaFailures"`
	ProviderFailures    uint64  `json:"providerFailures"`
	RightsWithheldItems uint64  `json:"rightsWithheldItems"`
	ContextCompactions  uint64  `json:"contextCompactions"`
	TotalInputTokens    uint64  `json:"totalInputTokens"`
	TotalOutputTokens   uint64  `json:"totalOutputTokens"`
	TotalCostUSD        float64 `json:"totalCostUsd"`
	LastLatencyMs       int64   `json:"lastLatencyMs"`
	MaxLatencyMs        int64   `json:"maxLatencyMs"`
	LastOutcome         string  `json:"lastOutcome"`
	LastUpdated         int64   `json:"lastUpdated"`
}

type aiTelemetryState struct {
	mu   sync.Mutex
	diag AIInferenceTelemetry
}

var (
	aiRightsRegistryOnce sync.Once
	aiRightsRegistry     map[string]AIDatasetRightsRecord
	aiRightsRegistryErr  error
	aiTelemetryByApp     sync.Map
)

func aiRightsKey(provider, dataset string) string {
	return strings.ToLower(strings.TrimSpace(provider)) + "|" + strings.ToLower(strings.TrimSpace(dataset))
}

func loadAIDatasetRightsRegistry() (map[string]AIDatasetRightsRecord, error) {
	aiRightsRegistryOnce.Do(func() {
		var doc aiDatasetRightsRegistryDocument
		if err := json.Unmarshal(providerDatasetAIRightsRegistryJSON, &doc); err != nil {
			aiRightsRegistryErr = fmt.Errorf("AI rights registry decode failed: %w", err)
			return
		}
		if doc.Schema != "DE.PULSE-PROVIDER-DATASET-AI-RIGHTS-1" || doc.PolicyVersion != aiRightsEgressPolicyVersion || strings.ToUpper(strings.TrimSpace(doc.DefaultDecision)) != "BLOCK" {
			aiRightsRegistryErr = fmt.Errorf("AI rights registry policy mismatch")
			return
		}
		rows := make(map[string]AIDatasetRightsRecord, len(doc.Records))
		for _, record := range doc.Records {
			key := aiRightsKey(record.Provider, record.Dataset)
			if key == "|" {
				aiRightsRegistryErr = fmt.Errorf("AI rights registry contains an empty provider/dataset key")
				return
			}
			if _, exists := rows[key]; exists {
				aiRightsRegistryErr = fmt.Errorf("AI rights registry contains duplicate key %s", key)
				return
			}
			rows[key] = record
		}
		aiRightsRegistry = rows
	})
	return aiRightsRegistry, aiRightsRegistryErr
}

func aiRightsRecordAllows(record AIDatasetRightsRecord) bool {
	return record.EvidenceBound &&
		strings.EqualFold(strings.TrimSpace(record.Decision), "ALLOW") &&
		strings.EqualFold(strings.TrimSpace(record.CommercialUse), providerRightsApproved) &&
		strings.EqualFold(strings.TrimSpace(record.Redistribution), providerRightsApproved) &&
		strings.EqualFold(strings.TrimSpace(record.AIUse), providerRightsApproved)
}

func aiEvidenceDataset(provider string, item AIEvidenceItem) string {
	kind := strings.ToLower(strings.TrimSpace(item.Kind))
	switch provider {
	case "finnhub":
		switch kind {
		case "fundamentals":
			return "fundamentals"
		case "earnings":
			return "earnings"
		case "news":
			return "news"
		default:
			return "derived"
		}
	case "alpaca":
		if strings.Contains(strings.ToLower(item.Label), "option") {
			return "options"
		}
		if kind == "quote" || strings.Contains(strings.ToLower(item.Label), "quote") || strings.Contains(strings.ToLower(item.Label), "liquidity") {
			return "quotes"
		}
		return "derived"
	case "fred", "bls":
		return "macro"
	case "eia":
		return "energy"
	case "twelve-data", "yfinance":
		return "market-context"
	case "cboe":
		return "vix"
	case "marketaux":
		return "news"
	case "sec-edgar":
		return "filings"
	case "community":
		return "community"
	case "depulse-derived":
		return "derived"
	default:
		return "unknown"
	}
}

func aiEvidenceRightsCoordinates(item AIEvidenceItem) []aiRightsCoordinate {
	source := strings.ToLower(strings.TrimSpace(item.Source + " " + item.Label))
	providers := make([]string, 0, 3)
	add := func(provider string) {
		for _, existing := range providers {
			if existing == provider {
				return
			}
		}
		providers = append(providers, provider)
	}
	for _, candidate := range []struct {
		needle   string
		provider string
	}{
		{"finnhub", "finnhub"},
		{"alpaca", "alpaca"},
		{"fred", "fred"},
		{"bureau of labor", "bls"},
		{"bls", "bls"},
		{"eia", "eia"},
		{"twelve data", "twelve-data"},
		{"twelvedata", "twelve-data"},
		{"yfinance", "yfinance"},
		{"yahoo", "yfinance"},
		{"cboe", "cboe"},
		{"marketaux", "marketaux"},
		{"sec/edgar", "sec-edgar"},
		{"sec edgar", "sec-edgar"},
	} {
		if strings.Contains(source, candidate.needle) {
			add(candidate.provider)
		}
	}
	if strings.EqualFold(strings.TrimSpace(item.Kind), "community") {
		add("community")
	}
	if len(providers) == 0 && (strings.Contains(source, "de.pulse") || strings.Contains(source, "canonical")) {
		add("depulse-derived")
	}
	if len(providers) == 0 {
		add("unknown")
	}
	coordinates := make([]aiRightsCoordinate, 0, len(providers))
	for _, provider := range providers {
		coordinates = append(coordinates, aiRightsCoordinate{Provider: provider, Dataset: aiEvidenceDataset(provider, item)})
	}
	return coordinates
}

func filterAIResearchPackageForEgressWithRegistry(pkg AIResearchPackage, rights map[string]AIDatasetRightsRecord) (AIResearchPackage, AIEgressRightsDiagnostics) {
	diag := AIEgressRightsDiagnostics{PolicyVersion: aiRightsEgressPolicyVersion}
	kept := make([]AIEvidenceItem, 0, len(pkg.Evidence))
	for _, item := range pkg.Evidence {
		coords := aiEvidenceRightsCoordinates(item)
		allowed := len(coords) > 0
		unknown := false
		denied := false
		for _, coord := range coords {
			record, ok := rights[aiRightsKey(coord.Provider, coord.Dataset)]
			if !ok {
				allowed = false
				unknown = true
				continue
			}
			if !aiRightsRecordAllows(record) {
				allowed = false
				if strings.EqualFold(strings.TrimSpace(record.Decision), "DENY") {
					denied = true
				} else {
					unknown = true
				}
			}
		}
		if allowed {
			kept = append(kept, item)
			diag.AllowedItems++
			continue
		}
		diag.WithheldItems++
		if denied {
			diag.DeniedItems++
		} else if unknown {
			diag.UnknownItems++
		}
	}
	pkg.Evidence = kept
	filterAIResearchPackageEvidenceIDs(&pkg)
	if diag.WithheldItems > 0 {
		warning := fmt.Sprintf("AI egress rights policy withheld %d evidence item(s); only explicitly evidence-bound provider/dataset content may leave DE.PULSE.", diag.WithheldItems)
		pkg.SafetyWarnings = appendUniqueString(pkg.SafetyWarnings, warning)
		pkg.MissingEvidence = appendUniqueString(pkg.MissingEvidence, "Provider/dataset evidence was withheld from external AI because required rights are not explicitly approved.")
	}
	diag.Summary = fmt.Sprintf("allowed=%d withheld=%d unknown=%d denied=%d; provider names, URLs and source text are intentionally omitted from diagnostics", diag.AllowedItems, diag.WithheldItems, diag.UnknownItems, diag.DeniedItems)
	refreshAIEgressPackageID(&pkg)
	return pkg, diag
}

func filterAIResearchPackageForEgress(pkg AIResearchPackage) (AIResearchPackage, AIEgressRightsDiagnostics, error) {
	rights, err := loadAIDatasetRightsRegistry()
	if err != nil {
		// Registry failure is fail-closed: return an evidence-empty package and a redacted diagnostic.
		withheld := len(pkg.Evidence)
		pkg.Evidence = nil
		filterAIResearchPackageEvidenceIDs(&pkg)
		pkg.SafetyWarnings = appendUniqueString(pkg.SafetyWarnings, "AI egress rights policy is unavailable; all provider/dataset evidence was withheld.")
		pkg.MissingEvidence = appendUniqueString(pkg.MissingEvidence, "Provider/dataset evidence was withheld because the AI egress rights policy could not be verified.")
		refreshAIEgressPackageID(&pkg)
		return pkg, AIEgressRightsDiagnostics{
			PolicyVersion: aiRightsEgressPolicyVersion,
			WithheldItems: withheld,
			UnknownItems:  withheld,
			Summary:       fmt.Sprintf("allowed=0 withheld=%d; rights registry unavailable and diagnostics are redacted", withheld),
		}, err
	}
	filtered, diag := filterAIResearchPackageForEgressWithRegistry(pkg, rights)
	return filtered, diag, nil
}

func filterAIResearchPackageEvidenceIDs(pkg *AIResearchPackage) {
	allowed := allowedEvidenceIDs(*pkg)
	filter := func(values []string) []string {
		out := make([]string, 0, len(values))
		for _, value := range values {
			if allowed[value] {
				out = appendUniqueString(out, value)
			}
		}
		return out
	}
	pkg.BullEvidenceIDs = filter(pkg.BullEvidenceIDs)
	pkg.BaseEvidenceIDs = filter(pkg.BaseEvidenceIDs)
	pkg.BearEvidenceIDs = filter(pkg.BearEvidenceIDs)
	pkg.RiskEvidenceIDs = filter(pkg.RiskEvidenceIDs)
	pkg.CatalystEvidenceIDs = filter(pkg.CatalystEvidenceIDs)
	pkg.MarketEvidenceIDs = filter(pkg.MarketEvidenceIDs)
	pkg.EventEvidenceIDs = filter(pkg.EventEvidenceIDs)
}

func refreshAIEgressPackageID(pkg *AIResearchPackage) {
	ids := make([]string, 0, len(pkg.Evidence))
	for _, item := range pkg.Evidence {
		ids = append(ids, item.ID)
	}
	sort.Strings(ids)
	pkg.PackageID = "aix-" + shortAIHash(aiContextPolicyVersion, aiRightsEgressPolicyVersion, pkg.PackageID, strings.Join(ids, "|"), strings.Join(pkg.MissingEvidence, "|"))
}

func aiEvidencePriority(item AIEvidenceItem) int {
	score := 0
	for _, role := range item.Roles {
		switch strings.ToLower(strings.TrimSpace(role)) {
		case "risk":
			score += 60
		case "event", "catalyst":
			score += 50
		case "bear", "bull":
			score += 30
		case "market":
			score += 20
		case "base":
			score += 10
		}
	}
	switch strings.ToLower(strings.TrimSpace(item.Kind)) {
	case "news", "sec", "earnings", "reaction":
		score += 35
	case "market", "fundamentals":
		score += 20
	}
	if item.Untrusted {
		score += 5 // Retain material adversarial/external evidence when it is rights-approved.
	}
	return score
}

func semanticCompactAIPackage(pkg AIResearchPackage, maxEvidence int) AIResearchPackage {
	if maxEvidence < 0 {
		maxEvidence = 0
	}
	items := append([]AIEvidenceItem(nil), pkg.Evidence...)
	sort.SliceStable(items, func(i, j int) bool {
		pi, pj := aiEvidencePriority(items[i]), aiEvidencePriority(items[j])
		if pi != pj {
			return pi > pj
		}
		if items[i].Timestamp != items[j].Timestamp {
			return items[i].Timestamp > items[j].Timestamp
		}
		return items[i].ID < items[j].ID
	})
	if len(items) > maxEvidence {
		items = items[:maxEvidence]
	}
	pkg.Evidence = items
	filterAIResearchPackageEvidenceIDs(&pkg)
	pkg.Contradictions = trimAIStringsToBytes(pkg.Contradictions, 6, 600)
	pkg.MissingEvidence = trimAIStringsToBytes(pkg.MissingEvidence, 6, 600)
	pkg.SafetyWarnings = trimAIStringsToBytes(pkg.SafetyWarnings, 6, 600)
	refreshAIEgressPackageID(&pkg)
	return pkg
}

func trimAIStringsToBytes(values []string, maxItems, maxBytes int) []string {
	out := make([]string, 0, minInt(len(values), maxItems))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		value = truncateUTF8Bytes(value, maxBytes)
		out = append(out, value)
		if len(out) >= maxItems {
			break
		}
	}
	return out
}

func truncateUTF8Bytes(value string, maxBytes int) string {
	if maxBytes <= 0 {
		return ""
	}
	if len([]byte(value)) <= maxBytes {
		return value
	}
	b := []byte(value)
	b = b[:maxBytes]
	for len(b) > 0 && !json.Valid(append([]byte{'"'}, append(bytes.ReplaceAll(b, []byte{'"'}, []byte{'\\', '"'}), '"')...)) {
		b = b[:len(b)-1]
	}
	return strings.TrimSpace(string(b))
}

func conservativeAITokenUpperBound(text string) int {
	// Every non-empty token consumes at least one encoded byte. Byte length is
	// therefore a deliberately conservative provider-independent upper bound.
	return len([]byte(text))
}

func aiStructuredOutputSchema(strictAdditionalProperties bool) map[string]any {
	stringArray := func(maxItems int) map[string]any {
		return map[string]any{"type": "array", "maxItems": maxItems, "items": map[string]any{"type": "string", "maxLength": 600}}
	}
	properties := map[string]any{
		"verdict":         map[string]any{"type": "string", "enum": []string{"FAVORABLE", "MIXED", "CAUTION", "INFORMATIONAL"}},
		"confidence":      map[string]any{"type": "integer", "minimum": 0, "maximum": 100},
		"bullCase":        stringArray(3),
		"baseCase":        stringArray(3),
		"bearCase":        stringArray(3),
		"reasons":         stringArray(3),
		"risks":           stringArray(3),
		"contradictions":  stringArray(3),
		"missingEvidence": stringArray(3),
		"evidenceIds":     map[string]any{"type": "array", "maxItems": 8, "items": map[string]any{"type": "string", "maxLength": 96}},
		"catalyst":        map[string]any{"type": "string", "maxLength": 600},
		"bestFitHorizon":  map[string]any{"type": "string", "enum": []string{"day", "swing", "long", "none"}},
		"nextAction":      map[string]any{"type": "string", "maxLength": 600},
		"summary":         map[string]any{"type": "string", "maxLength": 600},
		"details":         map[string]any{"type": "string", "maxLength": 1200},
	}
	schema := map[string]any{
		"type":       "object",
		"properties": properties,
		"required": []string{
			"verdict", "confidence", "bullCase", "baseCase", "bearCase", "reasons", "risks",
			"contradictions", "missingEvidence", "evidenceIds", "catalyst", "bestFitHorizon",
			"nextAction", "summary", "details",
		},
	}
	if strictAdditionalProperties {
		schema["additionalProperties"] = false
	}
	return schema
}

func aiSchemaFingerprint() string {
	b, _ := json.Marshal(aiStructuredOutputSchema(true))
	return shortAIHash(aiSchemaPolicyVersion, string(b))
}

func buildBoundedAIUserPrompt(task string, pkg AIResearchPackage) (string, AIResearchPackage, AIContextDiagnostics, error) {
	task = truncateUTF8Bytes(strings.TrimSpace(task), aiMaxTaskBytes)
	originalEvidence := len(pkg.Evidence)
	maxEvidence := minInt(len(pkg.Evidence), aiMaxEvidenceForPrompt)
	for maxEvidence >= 0 {
		candidate := semanticCompactAIPackage(pkg, maxEvidence)
		envelope := map[string]any{
			"task":            task,
			"evidencePackage": candidate,
			"responseSchema":  aiStructuredOutputSchema(true),
		}
		b, err := json.Marshal(envelope)
		if err != nil {
			return "", pkg, AIContextDiagnostics{}, err
		}
		tokenUpper := conservativeAITokenUpperBound(string(b))
		if len(b) <= aiMaxPromptBytes && tokenUpper <= aiMaxPromptTokenUpperBound {
			return string(b), candidate, AIContextDiagnostics{
				PolicyVersion:      aiContextPolicyVersion,
				OriginalEvidence:   originalEvidence,
				SentEvidence:       len(candidate.Evidence),
				Compacted:          len(candidate.Evidence) < originalEvidence,
				PromptBytes:        len(b),
				TokenUpperBound:    tokenUpper,
				MaxBytes:           aiMaxPromptBytes,
				MaxTokenUpperBound: aiMaxPromptTokenUpperBound,
			}, nil
		}
		maxEvidence--
	}
	return "", pkg, AIContextDiagnostics{}, fmt.Errorf("AI context cannot fit the hard byte/token budget")
}

func aiRoutingIdentity(settings Settings, routing AIRoutingDecision) string {
	parts := make([]string, 0, len(routing.Candidates)+2)
	parts = append(parts, strings.ToLower(strings.TrimSpace(routing.Policy)), strings.ToLower(strings.TrimSpace(settings.AIProvider)))
	for _, candidate := range routing.Candidates {
		provider := strings.ToLower(strings.TrimSpace(candidate.Provider))
		modelIdentity := ""
		switch provider {
		case "groq":
			modelIdentity = settings.GroqModel
		case "gemini":
			modelIdentity = settings.GeminiModel
		case "openrouter":
			cfg := openRouterConfig(candidate.Mode, settings.OpenRouterSpecificModel)
			modelIdentity = cfg.Primary + ">" + strings.Join(cfg.Fallback, ",")
		}
		parts = append(parts, fmt.Sprintf("%s:%s:%s:%d", provider, strings.ToLower(strings.TrimSpace(candidate.Mode)), strings.ToLower(strings.TrimSpace(modelIdentity)), candidate.MaxOutputTokens))
	}
	return strings.Join(parts, "|")
}

func aiInferenceCacheKey(req AIRequest, task string, pkg AIResearchPackage, routing AIRoutingDecision, settings Settings) string {
	return shortAIHash(
		aiCachePolicyVersion,
		aiContextPolicyVersion,
		aiSafetyPolicyVersion,
		aiSchemaPolicyVersion,
		aiRightsEgressPolicyVersion,
		aiSchemaFingerprint(),
		shortAIHash(aiSystemPrompt()),
		pkg.PackageID,
		pkg.EvidenceSnapshotID,
		strings.ToLower(strings.TrimSpace(req.Kind)),
		strings.TrimSpace(task),
		strings.ToLower(strings.TrimSpace(req.ScopeType)),
		strings.TrimSpace(req.WatchlistID),
		normalizeSymbol(req.Ticker),
		aiRoutingIdentity(settings, routing),
	)
}

func (a *Application) loadAICacheV186(key string, now time.Time) (AIResponse, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	entry, ok := a.aiCache[key]
	if !ok {
		return AIResponse{}, false
	}
	age := now.Sub(time.UnixMilli(entry.StoredAt))
	if age < 0 || age > aiInferenceCacheTTL {
		delete(a.aiCache, key)
		return AIResponse{}, false
	}
	return entry.Response, true
}

func parseAIStructuredPayloadStrict(text string, pkg AIResearchPackage) (aiStructuredPayload, error) {
	raw := strings.TrimSpace(text)
	if raw == "" || strings.HasPrefix(raw, "```") {
		return aiStructuredPayload{}, fmt.Errorf("AI structured output is empty or markdown-wrapped")
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(raw), &fields); err != nil {
		return aiStructuredPayload{}, fmt.Errorf("AI structured output is not valid JSON: %w", err)
	}
	required := []string{
		"verdict", "confidence", "bullCase", "baseCase", "bearCase", "reasons", "risks",
		"contradictions", "missingEvidence", "evidenceIds", "catalyst", "bestFitHorizon",
		"nextAction", "summary", "details",
	}
	for _, field := range required {
		if _, ok := fields[field]; !ok {
			return aiStructuredPayload{}, fmt.Errorf("AI structured output missing required field %s", field)
		}
	}
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.DisallowUnknownFields()
	var out aiStructuredPayload
	if err := decoder.Decode(&out); err != nil {
		return aiStructuredPayload{}, fmt.Errorf("AI structured output schema mismatch: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return aiStructuredPayload{}, fmt.Errorf("AI structured output contains trailing content")
	}
	out.Verdict = strings.ToUpper(strings.TrimSpace(out.Verdict))
	switch out.Verdict {
	case "FAVORABLE", "MIXED", "CAUTION", "INFORMATIONAL":
	default:
		return aiStructuredPayload{}, fmt.Errorf("AI structured output has invalid verdict")
	}
	if out.Confidence < 0 || out.Confidence > 100 {
		return aiStructuredPayload{}, fmt.Errorf("AI structured output confidence is outside 0-100")
	}
	out.BestFitHorizon = strings.ToLower(strings.TrimSpace(out.BestFitHorizon))
	if out.BestFitHorizon != "day" && out.BestFitHorizon != "swing" && out.BestFitHorizon != "long" && out.BestFitHorizon != "none" {
		return aiStructuredPayload{}, fmt.Errorf("AI structured output has invalid bestFitHorizon")
	}
	for name, values := range map[string][]string{
		"bullCase": out.BullCase, "baseCase": out.BaseCase, "bearCase": out.BearCase,
		"reasons": out.Reasons, "risks": out.Risks, "contradictions": out.Contradictions,
		"missingEvidence": out.MissingEvidence,
	} {
		if len(values) > 3 {
			return aiStructuredPayload{}, fmt.Errorf("AI structured output %s exceeds max items", name)
		}
		for _, value := range values {
			if strings.TrimSpace(value) == "" || len([]byte(value)) > 600 {
				return aiStructuredPayload{}, fmt.Errorf("AI structured output %s contains an invalid item", name)
			}
		}
	}
	if len(out.EvidenceIDs) > 8 {
		return aiStructuredPayload{}, fmt.Errorf("AI structured output evidenceIds exceeds max items")
	}
	allowed := allowedEvidenceIDs(pkg)
	seen := map[string]bool{}
	for _, id := range out.EvidenceIDs {
		id = strings.TrimSpace(id)
		if id == "" || !allowed[id] || seen[id] {
			return aiStructuredPayload{}, fmt.Errorf("AI structured output contains an unknown or duplicate evidence ID")
		}
		seen[id] = true
	}
	for name, value := range map[string]string{
		"catalyst": out.Catalyst, "nextAction": out.NextAction, "summary": out.Summary,
	} {
		if len([]byte(value)) > 600 {
			return aiStructuredPayload{}, fmt.Errorf("AI structured output %s exceeds max length", name)
		}
	}
	if strings.TrimSpace(out.Summary) == "" {
		return aiStructuredPayload{}, fmt.Errorf("AI structured output summary is empty")
	}
	if len([]byte(out.Details)) > 1200 {
		return aiStructuredPayload{}, fmt.Errorf("AI structured output details exceeds max length")
	}
	out.NextAction = sanitizeAINextAction(out.NextAction)
	return out, nil
}

func aiTelemetryStateFor(a *Application) *aiTelemetryState {
	if value, ok := aiTelemetryByApp.Load(a); ok {
		return value.(*aiTelemetryState)
	}
	state := &aiTelemetryState{diag: AIInferenceTelemetry{PolicyVersion: aiEvalPolicyVersion}}
	actual, _ := aiTelemetryByApp.LoadOrStore(a, state)
	return actual.(*aiTelemetryState)
}

func (a *Application) recordAIInferenceTelemetry(outcome string, latencyMs int64, inputTokens, outputTokens int, cost float64, contextDiag AIContextDiagnostics, rightsDiag AIEgressRightsDiagnostics, cacheHit bool, schemaFailure bool, providerFailure bool) {
	state := aiTelemetryStateFor(a)
	state.mu.Lock()
	defer state.mu.Unlock()
	d := &state.diag
	d.Requests++
	if cacheHit {
		d.CacheHits++
	}
	if schemaFailure {
		d.SchemaFailures++
	}
	if providerFailure {
		d.ProviderFailures++
	}
	if rightsDiag.WithheldItems > 0 {
		d.RightsWithheldItems += uint64(rightsDiag.WithheldItems)
	}
	if contextDiag.Compacted {
		d.ContextCompactions++
	}
	if inputTokens > 0 {
		d.TotalInputTokens += uint64(inputTokens)
	}
	if outputTokens > 0 {
		d.TotalOutputTokens += uint64(outputTokens)
	}
	if cost > 0 {
		d.TotalCostUSD += cost
	}
	d.LastLatencyMs = latencyMs
	if latencyMs > d.MaxLatencyMs {
		d.MaxLatencyMs = latencyMs
	}
	d.LastOutcome = strings.TrimSpace(outcome)
	d.LastUpdated = time.Now().UnixMilli()
}
