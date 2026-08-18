package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func getJSON(ctx context.Context, client *http.Client, rawURL string, headers map[string]string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 300))
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
func mustInt(s string) int { n, _ := strconv.Atoi(s); return n }
func contains(items []string, v string) bool {
	for _, x := range items {
		if x == v {
			return true
		}
	}
	return false
}

func (a *Application) GenerateAIForUser(ctx context.Context, userID string, req AIRequest) (AIResponse, error) {
	a.mu.RLock()
	settings := a.state.Settings
	secrets := a.secrets
	ui := a.workspaceStateLocked(userID).UI
	a.mu.RUnlock()
	if settings.AIProvider == "" {
		settings.AIProvider = "groq"
	}
	if settings.AIRoutingMode == "" {
		settings.AIRoutingMode = "manual"
	}
	if settings.GroqModel == "" {
		settings.GroqModel = "openai/gpt-oss-120b"
	}
	if settings.GeminiModel == "" {
		settings.GeminiModel = "gemini-3.1-flash-lite"
	}
	if settings.OpenRouterMode == "" {
		settings.OpenRouterMode = "fast"
	}
	if !allowedOpenRouterModel(settings.OpenRouterSpecificModel) {
		settings.OpenRouterSpecificModel = "openai/gpt-5.6-sol"
	}
	if req.ScopeType == "" {
		req.ScopeType = ui.ScopeType
	}
	if req.WatchlistID == "" {
		req.WatchlistID = ui.WatchlistID
	}
	if req.Ticker == "" {
		req.Ticker = ui.SelectedTicker
	}
	if strings.TrimSpace(req.Ticker) != "" {
		var valid bool
		req.Ticker, valid = parseSelectableTicker(req.Ticker)
		if !valid {
			return AIResponse{}, errors.New("invalid ticker")
		}
	}

	snap := a.engine.Snapshot()
	pkg := buildAIResearchPackage(req, snap)
	task := strings.TrimSpace(req.Question)
	if task == "" {
		switch req.Kind {
		case "morning":
			task = "Create a compact morning briefing from the supplied evidence: market/watchlist mood, important moves, material events, risks, contradictions, and what deserves human attention today."
		case "ticker":
			task = "Analyze " + req.Ticker + " as a sourced second opinion. Build Bull/Base/Bear cases, identify contradictions and missing evidence, and explain the best-fit horizon without changing deterministic decisions."
		case "news":
			task = "Analyze the supplied material news/catalyst evidence. Distinguish what is new, the plausible business/valuation mechanism, whether the observed reaction is proportional, and what remains unproven."
		case "risk":
			task = "Review the supplied evidence for market, event, liquidity, freshness, contradiction, and missing-evidence risk. Do not recommend or execute trades."
		default:
			task = "Summarize what matters most in the supplied DE.PULSE evidence package."
		}
	}

	routing := resolveAIRouting(settings, secrets, req)
	if len(routing.Candidates) == 0 {
		if routing.Policy == "manual" {
			return AIResponse{}, fmt.Errorf("add a %s API key in Settings or choose an automatic routing policy", strings.Title(settings.AIProvider))
		}
		return AIResponse{}, errors.New("configure at least one AI provider key in Settings")
	}

	// Provider/dataset rights are evaluated before any external AI request. A
	// missing/invalid rights registry fails closed by producing an evidence-empty
	// package; the task itself can still be answered as an explicit abstention.
	pkg, rightsDiag, rightsErr := filterAIResearchPackageForEgress(pkg)
	userPrompt, boundedPkg, contextDiag, err := buildBoundedAIUserPrompt(task, pkg)
	if err != nil {
		a.recordAIInferenceTelemetry("context-budget-failure", 0, 0, 0, 0, contextDiag, rightsDiag, false, false, false)
		return AIResponse{}, err
	}
	pkg = boundedPkg

	cacheKey := aiInferenceCacheKey(req, task, pkg, routing, settings)
	if cached, ok := a.loadAICacheV186(cacheKey, time.Now()); ok {
		cached.CacheHit = true
		cached.RouteReason = "Inference identity unchanged across evidence, provider/model route, prompt, safety, rights and schema policies; reused a TTL-valid AI second opinion without another provider call."
		cached.EvidencePackageID = pkg.PackageID
		cached.EvidenceSnapshotID = pkg.EvidenceSnapshotID
		cached.SafetyWarnings = append([]string(nil), pkg.SafetyWarnings...)
		a.recordAIInferenceTelemetry("cache-hit", 0, cached.InputTokens, cached.OutputTokens, 0, contextDiag, rightsDiag, true, false, false)
		return cached, nil
	}

	systemPrompt := aiSystemPrompt()
	a.engine.setHealth("ai", "loading")
	startAt := time.Now()
	var text, model, requestedModel, providerLabel, providerID, mode string
	var inputTokens, outputTokens int
	var cost float64
	var fallback bool
	var lastErr error
	var routeReason string
	var structured aiStructuredPayload
	var sawSchemaFailure bool
	var sawProviderFailure bool

	for i, candidate := range routing.Candidates {
		providerID = candidate.Provider
		routeReason = candidate.Reason
		text, model, requestedModel, providerLabel, mode = "", "", "", "", candidate.Mode
		inputTokens, outputTokens, cost, fallback = 0, 0, 0, i > 0
		switch candidate.Provider {
		case "groq":
			model, requestedModel, providerLabel = settings.GroqModel, settings.GroqModel, "Groq"
			text, lastErr = generateOpenAICompatibleResponse(ctx, "https://api.groq.com/openai/v1/responses", secrets.Groq, model, systemPrompt, userPrompt, candidate.MaxOutputTokens, false)
		case "openrouter":
			providerLabel = "OpenRouter"
			orMode := candidate.Mode
			if orMode == "" {
				orMode = settings.OpenRouterMode
			}
			var result OpenRouterResult
			result, lastErr = generateOpenRouterResponse(ctx, secrets.OpenRouter, orMode, settings.OpenRouterSpecificModel, systemPrompt, userPrompt, candidate.MaxOutputTokens)
			if lastErr == nil {
				text, model, requestedModel, mode = result.Text, result.ActualModel, result.RequestedModel, result.Mode
				inputTokens, outputTokens, cost = result.InputTokens, result.OutputTokens, result.Cost
				fallback = fallback || result.Fallback
			}
		case "gemini":
			model, requestedModel, providerLabel = settings.GeminiModel, settings.GeminiModel, "Google Gemini"
			text, lastErr = generateGeminiResponse(ctx, secrets.Gemini, model, systemPrompt, userPrompt, candidate.MaxOutputTokens)
		default:
			lastErr = errors.New("unsupported AI route")
		}
		if lastErr == nil && strings.TrimSpace(text) != "" {
			var schemaErr error
			structured, schemaErr = parseAIStructuredPayloadStrict(text, pkg)
			if schemaErr == nil {
				break
			}
			sawSchemaFailure = true
			lastErr = fmt.Errorf("AI provider response failed strict schema/citation validation: %w", schemaErr)
		}
		if lastErr == nil {
			lastErr = errors.New("AI provider returned no text")
		}
		if !sawSchemaFailure {
			sawProviderFailure = true
		}
		if isContextLimitError(lastErr) && len(pkg.Evidence) > 0 {
			pkg = semanticCompactAIPackage(pkg, minInt(len(pkg.Evidence), 12))
			var compactErr error
			userPrompt, pkg, contextDiag, compactErr = buildBoundedAIUserPrompt(task, pkg)
			if compactErr != nil {
				lastErr = compactErr
				break
			}
			cacheKey = aiInferenceCacheKey(req, task, pkg, routing, settings)
			routeReason = candidate.Reason + " Provider context limit encountered; subsequent fallback uses a smaller materiality-ranked evidence package with a new cache identity."
		}
	}
	if lastErr != nil || strings.TrimSpace(text) == "" {
		latency := time.Since(startAt).Milliseconds()
		a.engine.setHealth("ai", "degraded")
		if lastErr == nil {
			lastErr = errors.New("AI response did not contain text")
		}
		a.recordAIInferenceTelemetry("safe-abstention", latency, inputTokens, outputTokens, cost, contextDiag, rightsDiag, false, sawSchemaFailure, sawProviderFailure)
		return AIResponse{}, lastErr
	}

	latency := time.Since(startAt).Milliseconds()
	structured = sanitizeAIResponse(structured, pkg)
	a.engine.setHealth("ai", "ready")
	if rightsDiag.WithheldItems > 0 {
		routeReason += fmt.Sprintf(" Rights-aware egress withheld %d unapproved provider/dataset evidence item(s).", rightsDiag.WithheldItems)
	}
	if rightsErr != nil {
		routeReason += " Rights registry verification failed closed; no unverified provider/dataset evidence was sent."
	}
	response := AIResponse{
		Text: strings.TrimSpace(text), Verdict: structured.Verdict, Confidence: structured.Confidence,
		BullCase: structured.BullCase, BaseCase: structured.BaseCase, BearCase: structured.BearCase,
		Reasons: structured.Reasons, Risks: structured.Risks, Contradictions: structured.Contradictions,
		MissingEvidence: structured.MissingEvidence, EvidenceIDs: structured.EvidenceIDs,
		Catalyst: structured.Catalyst, BestFitHorizon: structured.BestFitHorizon, NextAction: structured.NextAction,
		Summary: structured.Summary, Details: structured.Details, Model: model, RequestedModel: requestedModel,
		Provider: providerID, ProviderLabel: providerLabel, Mode: mode, LatencyMs: latency,
		InputTokens: inputTokens, OutputTokens: outputTokens, Cost: cost, Fallback: fallback,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339), EvidencePackageID: pkg.PackageID,
		EvidenceSnapshotID: pkg.EvidenceSnapshotID, RoutePolicy: routing.Policy,
		RouteReason: routeReason, CacheHit: false, SafetyWarnings: pkg.SafetyWarnings,
	}
	a.storeAICache(cacheKey, response)
	a.recordAIInferenceTelemetry("success", latency, inputTokens, outputTokens, cost, contextDiag, rightsDiag, false, sawSchemaFailure, sawProviderFailure)
	return response, nil
}

type OpenRouterConfig struct {
	Mode     string
	Primary  string
	Fallback []string
}

type OpenRouterResult struct {
	Text           string
	Mode           string
	RequestedModel string
	ActualModel    string
	InputTokens    int
	OutputTokens   int
	Cost           float64
	Fallback       bool
}

func allowedOpenRouterModel(model string) bool {
	m := strings.ToLower(strings.TrimSpace(model))
	if m == "" {
		return false
	}
	return strings.HasPrefix(m, "openai/") || strings.HasPrefix(m, "x-ai/") || strings.HasPrefix(m, "anthropic/") || strings.HasPrefix(m, "~anthropic/")
}

func openRouterConfig(mode, specific string) OpenRouterConfig {
	mode = strings.ToLower(strings.TrimSpace(mode))
	switch mode {
	case "balanced":
		return OpenRouterConfig{Mode: "balanced", Primary: "x-ai/grok-4.20", Fallback: []string{"openai/gpt-5.6-terra", "anthropic/claude-sonnet-5"}}
	case "powerful":
		return OpenRouterConfig{Mode: "powerful", Primary: "anthropic/claude-sonnet-5", Fallback: []string{"openai/gpt-5.6-sol", "x-ai/grok-4.5"}}
	case "specific":
		if !allowedOpenRouterModel(specific) {
			specific = "openai/gpt-5.6-sol"
		}
		return OpenRouterConfig{Mode: "specific", Primary: specific}
	default:
		return OpenRouterConfig{Mode: "fast", Primary: "openai/gpt-5.6-luna", Fallback: []string{"x-ai/grok-4.20", "anthropic/claude-haiku-4.5"}}
	}
}

func extractChatCompletionText(data map[string]any) string {
	choices, _ := data["choices"].([]any)
	if len(choices) == 0 {
		return ""
	}
	choice, _ := choices[0].(map[string]any)
	msg, _ := choice["message"].(map[string]any)
	if text, ok := msg["content"].(string); ok {
		return strings.TrimSpace(text)
	}
	return ""
}

func generateOpenRouterResponse(ctx context.Context, key, mode, specific, systemPrompt, userPrompt string, maxOutputTokens int) (OpenRouterResult, error) {
	config := openRouterConfig(mode, specific)
	if maxOutputTokens <= 0 {
		maxOutputTokens = 1400
	}
	payload := map[string]any{
		"model": config.Primary,
		"messages": []any{
			map[string]any{"role": "system", "content": systemPrompt},
			map[string]any{"role": "user", "content": userPrompt},
		},
		"provider": map[string]any{
			"sort":               "throughput",
			"allow_fallbacks":    true,
			"require_parameters": true,
		},
		"response_format": map[string]any{
			"type": "json_schema",
			"json_schema": map[string]any{
				"name":   "depulse_research",
				"strict": true,
				"schema": aiStructuredOutputSchema(true),
			},
		},
		"usage":      map[string]any{"include": true},
		"max_tokens": maxOutputTokens,
	}
	if len(config.Fallback) > 0 {
		payload["models"] = config.Fallback
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://openrouter.ai/api/v1/chat/completions", strings.NewReader(string(body)))
	if err != nil {
		return OpenRouterResult{}, err
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Title", appName)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return OpenRouterResult{}, err
	}
	defer resp.Body.Close()
	var data map[string]any
	if json.NewDecoder(resp.Body).Decode(&data) != nil {
		return OpenRouterResult{}, errors.New("invalid OpenRouter response")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return OpenRouterResult{}, providerAPIError(data, resp.StatusCode)
	}
	text := extractChatCompletionText(data)
	actual := strings.TrimSpace(fmt.Sprint(data["model"]))
	if actual == "<nil>" || actual == "" {
		actual = config.Primary
	}
	usage, _ := data["usage"].(map[string]any)
	inTok := int(toFloat(usage["prompt_tokens"]))
	outTok := int(toFloat(usage["completion_tokens"]))
	cost := toFloat(usage["cost"])
	return OpenRouterResult{Text: text, Mode: config.Mode, RequestedModel: config.Primary, ActualModel: actual, InputTokens: inTok, OutputTokens: outTok, Cost: cost, Fallback: !strings.EqualFold(actual, config.Primary)}, nil
}

func generateOpenAICompatibleResponse(ctx context.Context, endpoint, key, model, systemPrompt, userPrompt string, maxOutputTokens int, includeStore bool) (string, error) {
	if maxOutputTokens <= 0 {
		maxOutputTokens = 1400
	}
	payload := map[string]any{
		"model":             model,
		"instructions":      systemPrompt,
		"input":             userPrompt,
		"max_output_tokens": maxOutputTokens,
		"text": map[string]any{
			"format": map[string]any{
				"type":   "json_schema",
				"name":   "depulse_research",
				"strict": true,
				"schema": aiStructuredOutputSchema(true),
			},
		},
	}
	if includeStore {
		payload["store"] = false
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var data map[string]any
	if json.NewDecoder(resp.Body).Decode(&data) != nil {
		return "", errors.New("invalid AI response")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", providerAPIError(data, resp.StatusCode)
	}
	return extractAIText(data), nil
}

func generateGeminiResponse(ctx context.Context, key, model, systemPrompt, userPrompt string, maxOutputTokens int) (string, error) {
	if maxOutputTokens <= 0 {
		maxOutputTokens = 1400
	}
	body, _ := json.Marshal(map[string]any{
		"system_instruction": map[string]any{"parts": []any{map[string]any{"text": systemPrompt}}},
		"contents":           []any{map[string]any{"role": "user", "parts": []any{map[string]any{"text": userPrompt}}}},
		"generationConfig": map[string]any{
			"maxOutputTokens": maxOutputTokens,
			"responseMimeType": "application/json",
			"responseSchema":   aiStructuredOutputSchema(false),
		},
	})
	endpoint := "https://generativelanguage.googleapis.com/v1beta/models/" + url.PathEscape(model) + ":generateContent"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("x-goog-api-key", key)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var data map[string]any
	if json.NewDecoder(resp.Body).Decode(&data) != nil {
		return "", errors.New("invalid Gemini response")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", providerAPIError(data, resp.StatusCode)
	}
	return extractGeminiText(data), nil
}

func providerAPIError(data map[string]any, status int) error {
	if er, ok := data["error"].(map[string]any); ok {
		if m, ok := er["message"].(string); ok && strings.TrimSpace(m) != "" {
			return errors.New(strings.TrimSpace(m))
		}
	}
	return fmt.Errorf("AI provider error %d", status)
}

func extractAIText(data map[string]any) string {
	if s, ok := data["output_text"].(string); ok {
		return strings.TrimSpace(s)
	}
	var chunks []string
	if outputs, ok := data["output"].([]any); ok {
		for _, o := range outputs {
			om, _ := o.(map[string]any)
			contents, _ := om["content"].([]any)
			for _, c := range contents {
				cm, _ := c.(map[string]any)
				if s, ok := cm["text"].(string); ok {
					chunks = append(chunks, s)
				}
				if s, ok := cm["output_text"].(string); ok {
					chunks = append(chunks, s)
				}
			}
		}
	}
	return strings.TrimSpace(strings.Join(chunks, "\n"))
}

func extractGeminiText(data map[string]any) string {
	var chunks []string
	candidates, _ := data["candidates"].([]any)
	for _, candidate := range candidates {
		cm, _ := candidate.(map[string]any)
		content, _ := cm["content"].(map[string]any)
		parts, _ := content["parts"].([]any)
		for _, part := range parts {
			pm, _ := part.(map[string]any)
			if text, ok := pm["text"].(string); ok && strings.TrimSpace(text) != "" {
				chunks = append(chunks, strings.TrimSpace(text))
			}
		}
	}
	return strings.TrimSpace(strings.Join(chunks, "\n"))
}

type AIRequest struct {
	Kind          string         `json:"kind"`
	Question      string         `json:"question"`
	ScopeType     string         `json:"scopeType"`
	WatchlistID   string         `json:"watchlistId"`
	Ticker        string         `json:"ticker"`
	ClientContext map[string]any `json:"clientContext,omitempty"`
}
type AIResponse struct {
	Text               string   `json:"text"`
	Verdict            string   `json:"verdict,omitempty"`
	Confidence         int      `json:"confidence,omitempty"`
	BullCase           []string `json:"bullCase,omitempty"`
	BaseCase           []string `json:"baseCase,omitempty"`
	BearCase           []string `json:"bearCase,omitempty"`
	Reasons            []string `json:"reasons,omitempty"`
	Risks              []string `json:"risks,omitempty"`
	Contradictions     []string `json:"contradictions,omitempty"`
	MissingEvidence    []string `json:"missingEvidence,omitempty"`
	EvidenceIDs        []string `json:"evidenceIds,omitempty"`
	Catalyst           string   `json:"catalyst,omitempty"`
	BestFitHorizon     string   `json:"bestFitHorizon,omitempty"`
	NextAction         string   `json:"nextAction,omitempty"`
	Summary            string   `json:"summary,omitempty"`
	Details            string   `json:"details,omitempty"`
	Model              string   `json:"model"`
	RequestedModel     string   `json:"requestedModel,omitempty"`
	Provider           string   `json:"provider"`
	ProviderLabel      string   `json:"providerLabel"`
	Mode               string   `json:"mode,omitempty"`
	LatencyMs          int64    `json:"latencyMs,omitempty"`
	InputTokens        int      `json:"inputTokens,omitempty"`
	OutputTokens       int      `json:"outputTokens,omitempty"`
	Cost               float64  `json:"cost,omitempty"`
	Fallback           bool     `json:"fallback,omitempty"`
	GeneratedAt        string   `json:"generatedAt"`
	EvidencePackageID  string   `json:"evidencePackageId,omitempty"`
	EvidenceSnapshotID string   `json:"evidenceSnapshotId,omitempty"`
	RoutePolicy        string   `json:"routePolicy,omitempty"`
	RouteReason        string   `json:"routeReason,omitempty"`
	CacheHit           bool     `json:"cacheHit,omitempty"`
	SafetyWarnings     []string `json:"safetyWarnings,omitempty"`
}

type aiStructuredPayload struct {
	Verdict         string   `json:"verdict"`
	Confidence      int      `json:"confidence"`
	BullCase        []string `json:"bullCase"`
	BaseCase        []string `json:"baseCase"`
	BearCase        []string `json:"bearCase"`
	Reasons         []string `json:"reasons"`
	Risks           []string `json:"risks"`
	Contradictions  []string `json:"contradictions"`
	MissingEvidence []string `json:"missingEvidence"`
	EvidenceIDs     []string `json:"evidenceIds"`
	Catalyst        string   `json:"catalyst"`
	BestFitHorizon  string   `json:"bestFitHorizon"`
	NextAction      string   `json:"nextAction"`
	Summary         string   `json:"summary"`
	Details         string   `json:"details"`
}

// parseAIStructuredPayload remains for compatibility with historical tests and
// diagnostics. Production inference uses parseAIStructuredPayloadStrict and
// never accepts this lenient fallback as a successful AI result.
func parseAIStructuredPayload(text string) aiStructuredPayload {
	raw := strings.TrimSpace(text)
	if strings.HasPrefix(raw, "```") {
		if i := strings.Index(raw, "\n"); i >= 0 {
			raw = raw[i+1:]
		}
		raw = strings.TrimSpace(strings.TrimSuffix(raw, "```"))
	}
	var out aiStructuredPayload
	if json.Unmarshal([]byte(raw), &out) != nil {
		out.Verdict = "INFORMATIONAL"
		out.Confidence = 0
		out.Summary = "AI returned an unstructured response; review the details manually."
		out.Details = raw
		return out
	}
	out.Verdict = strings.ToUpper(strings.TrimSpace(out.Verdict))
	switch out.Verdict {
	case "FAVORABLE", "MIXED", "CAUTION", "INFORMATIONAL":
	default:
		out.Verdict = "INFORMATIONAL"
	}
	if out.Confidence < 0 {
		out.Confidence = 0
	}
	if out.Confidence > 100 {
		out.Confidence = 100
	}
	out.BestFitHorizon = strings.ToLower(strings.TrimSpace(out.BestFitHorizon))
	if out.BestFitHorizon != "day" && out.BestFitHorizon != "swing" && out.BestFitHorizon != "long" {
		out.BestFitHorizon = "none"
	}
	out.BullCase = trimAIStrings(out.BullCase, 3)
	out.BaseCase = trimAIStrings(out.BaseCase, 3)
	out.BearCase = trimAIStrings(out.BearCase, 3)
	out.Reasons = trimAIStrings(out.Reasons, 3)
	out.Risks = trimAIStrings(out.Risks, 3)
	out.Contradictions = trimAIStrings(out.Contradictions, 3)
	out.MissingEvidence = trimAIStrings(out.MissingEvidence, 3)
	if len(out.EvidenceIDs) > 8 {
		out.EvidenceIDs = out.EvidenceIDs[:8]
	}
	if len(out.Details) > 1400 {
		out.Details = out.Details[:1400] + "…"
	}
	return out
}
