package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type MaintenanceReport struct {
	Version     string             `json:"version"`
	BuildID     string             `json:"buildId"`
	GeneratedAt int64              `json:"generatedAt"`
	Platform    string             `json:"platform"`
	Pass        int                `json:"pass"`
	Warnings    int                `json:"warnings"`
	Failures    int                `json:"failures"`
	Checks      []MaintenanceCheck `json:"checks"`
}

func maintenanceHealthStatus(value string) string {
	v := strings.ToLower(strings.TrimSpace(value))
	switch {
	case v == "":
		return "warning"
	case strings.Contains(v, "error"), strings.Contains(v, "failed"), strings.Contains(v, "corrupt"):
		return "fail"
	case strings.Contains(v, "offline"), strings.Contains(v, "stopped"), strings.Contains(v, "unavailable"), strings.Contains(v, "setup required"), strings.Contains(v, "reconnecting"), strings.Contains(v, "waiting"):
		return "warning"
	default:
		return "pass"
	}
}

func (a *Application) handleMaintenanceRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	report := MaintenanceReport{Version: appVersion, BuildID: buildID, GeneratedAt: time.Now().UnixMilli(), Platform: runtime.GOOS + "/" + runtime.GOARCH, Checks: []MaintenanceCheck{}}
	add := func(id, label, status, detail string, started time.Time) {
		report.Checks = append(report.Checks, MaintenanceCheck{ID: id, Label: label, Status: status, Detail: detail, DurationMs: time.Since(started).Milliseconds()})
		switch status {
		case "pass":
			report.Pass++
		case "fail":
			report.Failures++
		default:
			report.Warnings++
		}
	}

	started := time.Now()
	if data, err := os.ReadFile(a.statePath()); err == nil {
		if json.Valid(data) {
			add("profile", "Profile / settings integrity", "pass", "state.json is readable and valid JSON", started)
		} else {
			add("profile", "Profile / settings integrity", "fail", "state.json is not valid JSON", started)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		add("profile", "Profile / settings integrity", "warning", "Profile has not been written to disk yet", started)
	} else {
		add("profile", "Profile / settings integrity", "fail", err.Error(), started)
	}

	started = time.Now()
	a.mu.RLock()
	stateCopy := clone(a.state)
	a.mu.RUnlock()
	seen := map[string]bool{}
	watchStatus, watchDetail := "pass", "Dedicated watchlists are present and symbol lists are normalized"
	for _, id := range []string{"day", "swing", "long", "discovery"} {
		found := false
		for _, wl := range stateCopy.Watchlists {
			if wl.ID == id {
				found = true
				break
			}
		}
		if !found {
			watchStatus, watchDetail = "fail", "Missing required watchlist: "+id
			break
		}
	}
	for _, wl := range stateCopy.Watchlists {
		for _, sym := range wl.Symbols {
			key := wl.ID + ":" + sym
			if seen[key] {
				watchStatus, watchDetail = "fail", "Duplicate symbol found in "+wl.ID
				break
			}
			seen[key] = true
		}
	}
	add("watchlists", "Watchlist integrity", watchStatus, watchDetail, started)

	started = time.Now()
	ci := func() CacheInfo { a.mu.RLock(); defer a.mu.RUnlock(); return a.cacheInfoLocked() }()
	if data, err := os.ReadFile(a.cachePath()); err == nil {
		var cache MarketCache
		if json.Unmarshal(data, &cache) == nil {
			add("cache", "Market cache integrity", "pass", fmt.Sprintf("%d cached symbols · %d bytes", ci.CachedSymbols, ci.SizeBytes), started)
		} else {
			add("cache", "Market cache integrity", "fail", "market-cache.json could not be decoded", started)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		add("cache", "Market cache integrity", "warning", "No market cache has been written yet", started)
	} else {
		add("cache", "Market cache integrity", "fail", err.Error(), started)
	}

	snap := a.engine.Snapshot()
	started = time.Now()
	runtimeStatus := maintenanceHealthStatus(snap.Status)
	if snap.Status == "running" {
		runtimeStatus = "pass"
	}
	add("runtime", "Runtime", runtimeStatus, fmt.Sprintf("%s · %s", snap.Status, snap.Message), started)

	started = time.Now()
	feedStatus := maintenanceHealthStatus(snap.Feed.FeedState)
	if snap.Feed.FeedState == "streaming" || snap.Feed.FeedState == "alpaca-streaming" || snap.Feed.FeedState == "snapshot-live" || snap.Feed.FeedState == "overnight-streaming" || snap.Feed.FeedState == "connected-market-closed" {
		feedStatus = "pass"
	}
	add("market-feed", "Market feed", feedStatus, fmt.Sprintf("%s · %s", snap.Feed.FeedState, snap.Feed.MarketSession), started)

	for _, key := range []string{"quotes", "quotes-rest", "alpaca-stream", "alpaca-live", "history", "fundamentals", "news", "earnings", "filings", "scanner", "vix", "ai"} {
		started = time.Now()
		value := strings.TrimSpace(snap.Health[key])
		if value == "" {
			value = "Not initialized"
		}
		add("health-"+key, "Data Engine · "+key, maintenanceHealthStatus(value), value, started)
	}

	started = time.Now()
	providerConfigured := 0
	a.mu.RLock()
	if strings.TrimSpace(a.secrets.Finnhub) != "" {
		providerConfigured++
	}
	if strings.TrimSpace(a.secrets.AlpacaKey) != "" && strings.TrimSpace(a.secrets.AlpacaSecret) != "" {
		providerConfigured++
	}
	if strings.TrimSpace(a.secrets.Groq) != "" || strings.TrimSpace(a.secrets.OpenRouter) != "" || strings.TrimSpace(a.secrets.Gemini) != "" {
		providerConfigured++
	}
	a.mu.RUnlock()
	providerStatus := "pass"
	if providerConfigured == 0 {
		providerStatus = "warning"
	}
	add("credentials", "Configured provider groups", providerStatus, fmt.Sprintf("%d of 3 provider groups configured", providerConfigured), started)

	started = time.Now()
	add("build", "Build identity", "pass", appVersion+" · "+buildID+" · "+runtime.GOOS+"/"+runtime.GOARCH, started)

	a.mu.Lock()
	a.state.MaintenanceLastRun = report.GeneratedAt
	_ = a.saveLocked()
	a.mu.Unlock()
	a.broadcastSharedState()
	writeJSON(w, 200, report)
}

func (a *Application) handleAI(w http.ResponseWriter, r *http.Request) {
	var in AIRequest
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid AI request")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()
	out, err := a.GenerateAIForUser(ctx, requestUserID(r.Context()), in)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}
	writeJSON(w, 200, out)
}

func (a *Application) handleExport(w http.ResponseWriter, r *http.Request) {
	userID := requestUserID(r.Context())
	a.mu.RLock()
	exportSettings := clone(a.state.Settings)
	workspaceState := a.workspaceStateLocked(userID)

	exportSettings.DataMode = "demo"
	exportSettings.AutoStart = false
	profile := map[string]any{"exportedAt": time.Now().UTC().Format(time.RFC3339), "version": a.state.Version, "settings": exportSettings, "watchlists": workspaceState.Watchlists, "ui": workspaceState.UI, "note": "API keys are intentionally excluded. Personal market state belongs only to the exporting user."}
	a.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", `attachment; filename="de-pulse-profile.json"`)
	_ = json.NewEncoder(w).Encode(profile)
}
func (a *Application) handleImport(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Profile json.RawMessage `json:"profile"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, 400, "Invalid profile")
		return
	}
	var p struct {
		Settings   Settings    `json:"settings"`
		Watchlists []Watchlist `json:"watchlists"`
		UI         UIState     `json:"ui"`
	}
	if json.Unmarshal(in.Profile, &p) != nil || len(p.Watchlists) == 0 {
		writeError(w, 400, "This is not a valid DE.PULSE profile")
		return
	}
	userID := requestUserID(r.Context())
	a.mu.Lock()
	currentMode := a.state.Settings.DataMode
	a.state.Settings = p.Settings
	a.state.Settings.DataMode = currentMode
	workspaceState := a.workspaceStateLocked(userID)
	workspaceState.Watchlists = clone(p.Watchlists)
	workspaceState.UI = p.UI
	workspaceState = mergeState(workspaceState)
	if err := a.saveWorkspaceStateLocked(userID, workspaceState); err != nil {
		a.mu.Unlock()
		writeError(w, 500, err.Error())
		return
	}
	if err := a.saveLocked(); err != nil {
		a.mu.Unlock()
		writeError(w, 500, err.Error())
		return
	}
	state := a.publicStateLockedForUser(userID)
	a.mu.Unlock()
	a.broadcastSharedState()
	writeJSON(w, 200, state)
}
func openExternalURL(raw string) error {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u == nil || (u.Scheme != "https" && u.Scheme != "http") || strings.TrimSpace(u.Host) == "" {
		return errors.New("only valid http/https links can be opened")
	}
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("/usr/bin/open", u.String()).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", u.String()).Start()
	default:
		return exec.Command("xdg-open", u.String()).Start()
	}
}

func (a *Application) handleOpenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	var in struct {
		URL string `json:"url"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, http.StatusBadRequest, "Invalid link")
		return
	}
	if err := openExternalURL(in.URL); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (a *Application) terminateAppWindow() {
	data, err := os.ReadFile(instancePath(a.configDir))
	if err == nil {
		var in instanceInfo
		if json.Unmarshal(data, &in) == nil && in.WindowPID > 0 && in.WindowPID != os.Getpid() {
			if p, findErr := os.FindProcess(in.WindowPID); findErr == nil {
				_ = p.Kill()
			}
		}
	}
	_ = os.Remove(instancePath(a.configDir))
}

func (a *Application) handleQuit(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"ok": true, "message": "DE.PULSE is closing."})
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	go func() {

		time.Sleep(650 * time.Millisecond)
		a.engine.Stop()
		a.terminateAppWindow()
		if a.server != nil {
			_ = a.server.Close()
		}
		time.Sleep(100 * time.Millisecond)
		os.Exit(0)
	}()
}

func toFloat(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case string:
		f, _ := strconv.ParseFloat(strings.TrimSpace(x), 64)
		return f
	case json.Number:
		f, _ := x.Float64()
		return f
	default:
		f, _ := strconv.ParseFloat(fmt.Sprint(v), 64)
		return f
	}
}
