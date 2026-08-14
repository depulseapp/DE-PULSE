package main

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (a *Application) handleCommunityEvidence(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Action        string `json:"action"`
		ID            string `json:"id"`
		Symbol        string `json:"symbol"`
		Source        string `json:"source"`
		Platform      string `json:"platform"`
		IngestionMode string `json:"ingestionMode"`
		RightsStatus  string `json:"rightsStatus"`
		AIEligibility string `json:"aiEligibility"`
		Stance        string `json:"stance"`
		Text          string `json:"text"`
		URL           string `json:"url"`
		ObservedAt    int64  `json:"observedAt"`
	}
	if decodeJSON(r, &in) != nil {
		writeError(w, http.StatusBadRequest, "Invalid community-evidence request")
		return
	}
	action := strings.ToLower(strings.TrimSpace(in.Action))
	if action == "delete" {
		id := strings.TrimSpace(in.ID)
		if id == "" {
			writeError(w, http.StatusBadRequest, "Evidence ID is required")
			return
		}
		a.mu.Lock()
		next := make([]CommunityEvidenceItem, 0, len(a.state.CommunityEvidence))
		found := false
		for _, item := range a.state.CommunityEvidence {
			if item.ID == id {
				found = true
				continue
			}
			next = append(next, item)
		}
		if !found {
			a.mu.Unlock()
			writeError(w, http.StatusNotFound, "Community evidence not found")
			return
		}
		a.state.CommunityEvidence = next
		err := a.saveLocked()
		a.mu.Unlock()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Could not save community evidence")
			return
		}
		a.broadcastSharedState()
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "runtime": a.engine.Snapshot()})
		return
	}
	if action != "add" && action != "" {
		writeError(w, http.StatusBadRequest, "Unknown community-evidence action")
		return
	}
	source := strings.TrimSpace(in.Source)
	text := strings.TrimSpace(in.Text)
	stance := strings.ToUpper(strings.TrimSpace(in.Stance))
	platform := canonicalCommunityPlatform(in.Platform, source)
	ingestionMode := strings.ToUpper(strings.TrimSpace(in.IngestionMode))
	if ingestionMode == "" {
		ingestionMode = "USER_AUTHORIZED_INPUT"
	}
	for _, banned := range []string{"SCRAPE", "SCRAPING", "BYPASS", "PRIVATE_SESSION_AUTOMATION", "UNAUTHORIZED"} {
		if ingestionMode == banned {
			writeError(w, http.StatusBadRequest, "Community ingestion must use an official, sanctioned, or user-authorized path")
			return
		}
	}
	if len(source) == 0 || len(source) > 80 || len(text) == 0 || len(text) > 1200 {
		writeError(w, http.StatusBadRequest, "Source and evidence text are required and must be within limits")
		return
	}
	switch stance {
	case "BULLISH", "BEARISH", "MIXED", "UNKNOWN":
	default:
		writeError(w, http.StatusBadRequest, "Stance must be BULLISH, BEARISH, MIXED, or UNKNOWN")
		return
	}
	symbol := ""
	if strings.TrimSpace(in.Symbol) != "" {
		var ok bool
		symbol, ok = parseUserTicker(in.Symbol)
		if !ok {
			writeError(w, http.StatusBadRequest, "Invalid ticker symbol")
			return
		}
	}
	link := strings.TrimSpace(in.URL)
	if len(link) > 600 {
		writeError(w, http.StatusBadRequest, "Source URL is too long")
		return
	}
	if link != "" {
		u, err := url.ParseRequestURI(link)
		if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" {
			writeError(w, http.StatusBadRequest, "Source URL must be an http(s) URL")
			return
		}
	}
	now := time.Now()
	observed := in.ObservedAt
	if observed <= 0 {
		observed = now.UnixMilli()
	}
	if observed > now.Add(5*time.Minute).UnixMilli() {
		writeError(w, http.StatusBadRequest, "Observed time cannot be in the future")
		return
	}
	item := normalizeCommunityEvidenceItem(CommunityEvidenceItem{ID: randomID("community"), Symbol: symbol, Source: source, Platform: platform, IngestionMode: ingestionMode, RightsStatus: strings.ToUpper(strings.TrimSpace(in.RightsStatus)), AIEligibility: strings.ToUpper(strings.TrimSpace(in.AIEligibility)), Stance: stance, Text: text, URL: link, ObservedAt: observed, SubmittedAt: now.UnixMilli()})
	if item.RetentionPolicy == "NO_RETENTION" {
		writeError(w, http.StatusBadRequest, "Source policy does not permit retaining this community evidence")
		return
	}
	a.mu.Lock()
	a.state.CommunityEvidence = append(a.state.CommunityEvidence, item)
	if len(a.state.CommunityEvidence) > 200 {
		a.state.CommunityEvidence = append([]CommunityEvidenceItem{}, a.state.CommunityEvidence[len(a.state.CommunityEvidence)-200:]...)
	}
	err := a.saveLocked()
	a.mu.Unlock()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Could not save community evidence")
		return
	}
	a.broadcastSharedState()
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "item": item, "runtime": a.engine.Snapshot()})
}
