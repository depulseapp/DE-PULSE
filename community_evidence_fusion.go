package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
	"unicode"
)

var communityPlatforms = []string{"X", "REDDIT", "TELEGRAM", "DISCORD", "WHATSAPP", "MANUAL", "OTHER"}

func canonicalCommunityPlatform(platform, source string) string {
	p := strings.ToUpper(strings.TrimSpace(platform))
	if p == "TWITTER" {
		p = "X"
	}
	if p == "" {
		s := strings.ToLower(source)
		switch {
		case strings.Contains(s, "telegram"), strings.Contains(s, "watcher guru"):
			p = "TELEGRAM"
		case strings.Contains(s, "reddit"), strings.Contains(s, "wallstreetbets"), strings.Contains(s, "wsb"):
			p = "REDDIT"
		case strings.Contains(s, "discord"):
			p = "DISCORD"
		case strings.Contains(s, "whatsapp"):
			p = "WHATSAPP"
		case strings.Contains(s, "twitter"), strings.HasPrefix(s, "x "), strings.Contains(s, " x.com"):
			p = "X"
		case strings.Contains(s, "note"), strings.Contains(s, "manual"):
			p = "MANUAL"
		default:
			p = "OTHER"
		}
	}
	for _, allowed := range communityPlatforms {
		if p == allowed {
			return p
		}
	}
	return "OTHER"
}

func communitySourcePolicy(platform, rights string) CommunitySourcePolicy {
	p := canonicalCommunityPlatform(platform, "")
	r := strings.ToUpper(strings.TrimSpace(rights))
	if r == "" {
		r = "USER_AUTHORIZED"
	}
	policy := CommunitySourcePolicy{Platform: p, RightsStatus: r, AIEligibility: "CONTEXT_ONLY", RetentionPolicy: "MATERIAL_STRUCTURED_ONLY", Status: "SUPPORTED WITH POLICY", AccessMode: "OFFICIAL_OR_USER_AUTHORIZED_ONLY"}
	switch p {
	case "TELEGRAM":
		policy.Detail = "Telegram evidence is rapid context only by default. Use sanctioned Bot/API/user-authorized paths; no scraping/private-session automation. API-derived content is not sent to AI unless an explicit permitted rights path is recorded."
	case "REDDIT":
		policy.Detail = "Reddit evidence is contextual by default. Official/user-authorized access only; commercial/AI/retention rights must be explicit before broader use."
	case "X":
		policy.Detail = "X evidence is contextual by default. Official API/user-authorized access only; usage economics and AI rights remain source-policy controlled."
	case "DISCORD":
		policy.Detail = "Discord evidence requires a sanctioned bot/application or explicit authorized input and server/channel permissions."
	case "WHATSAPP":
		policy.Detail = "WhatsApp evidence requires a sanctioned supported integration or explicit user-authorized forwarding/export path; no private-session automation."
	case "MANUAL":
		policy.AccessMode = "USER_AUTHORIZED_INPUT"
		policy.AIEligibility = "AI_ALLOWED"
		policy.RetentionPolicy = "MATERIAL_STRUCTURED_PLUS_USER_INPUT"
		policy.Detail = "User-authored/manual evidence may be analyzed by AI; provenance and untrusted-content controls still apply."
	default:
		policy.Detail = "Unknown/community source remains context-only until access, retention and AI rights are explicitly classified."
	}
	if r == "EXPLICIT_RIGHTS" || r == "AI_AND_RETENTION_RIGHTS" {
		policy.AIEligibility = "AI_ALLOWED"
		policy.RetentionPolicy = "MATERIAL_STRUCTURED_PLUS_PERMITTED_RAW"
		policy.Detail += " Explicit rights recorded for AI/retention use."
	}
	if r == "METADATA_ONLY" {
		policy.AIEligibility = "METADATA_ONLY"
		policy.RetentionPolicy = "METADATA_ONLY"
	}
	if r == "DISABLED" || r == "PROHIBITED" {
		policy.AIEligibility = "DISABLED"
		policy.RetentionPolicy = "NO_RETENTION"
		policy.Status = "DISABLED"
	}
	return policy
}

func allCommunitySourcePolicies() []CommunitySourcePolicy {
	out := make([]CommunitySourcePolicy, 0, len(communityPlatforms))
	for _, p := range communityPlatforms {
		out = append(out, communitySourcePolicy(p, ""))
	}
	return out
}

func communityCanonicalText(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	space := false
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			space = false
		} else if !space {
			b.WriteByte(' ')
			space = true
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

func communityFingerprint(text, link string) string {
	basis := communityCanonicalText(text)
	if u, err := url.Parse(strings.TrimSpace(link)); err == nil && u.Host != "" {
		basis += "|" + strings.ToLower(u.Host+u.Path)
	}
	sum := sha256.Sum256([]byte(basis))
	return hex.EncodeToString(sum[:12])
}

func communityTokenSimilarity(a, b string) float64 {
	aa, bb := map[string]bool{}, map[string]bool{}
	for _, x := range strings.Fields(communityCanonicalText(a)) {
		if len(x) > 2 {
			aa[x] = true
		}
	}
	for _, x := range strings.Fields(communityCanonicalText(b)) {
		if len(x) > 2 {
			bb[x] = true
		}
	}
	if len(aa) == 0 || len(bb) == 0 {
		return 0
	}
	inter, union := 0, len(aa)
	for x := range bb {
		if aa[x] {
			inter++
		} else {
			union++
		}
	}
	if union == 0 {
		return 0
	}
	return float64(inter) / float64(union)
}

func normalizeCommunityEvidenceItem(x CommunityEvidenceItem) CommunityEvidenceItem {
	x.Symbol = normalizeSymbol(x.Symbol)
	x.Platform = canonicalCommunityPlatform(x.Platform, x.Source)
	x.IngestionMode = strings.ToUpper(strings.TrimSpace(x.IngestionMode))
	if x.IngestionMode == "" {
		x.IngestionMode = "USER_AUTHORIZED_INPUT"
	}
	policy := communitySourcePolicy(x.Platform, x.RightsStatus)
	x.RightsStatus = policy.RightsStatus
	requestedAI := strings.ToUpper(strings.TrimSpace(x.AIEligibility))
	x.AIEligibility = policy.AIEligibility
	if requestedAI == "METADATA_ONLY" || requestedAI == "DISABLED" {
		x.AIEligibility = requestedAI
	}
	if requestedAI == "CONTEXT_ONLY" {
		x.AIEligibility = "CONTEXT_ONLY"
	}
	x.RetentionPolicy = policy.RetentionPolicy
	x.Fingerprint = communityFingerprint(x.Text, x.URL)
	return x
}

func corroborateCommunityCluster(c *CommunityEvidenceCluster, news []NewsItem, filings []FilingItem, latest int64) {
	if c == nil || c.Symbol == "" {
		return
	}
	window := int64(24 * time.Hour / time.Millisecond)
	for _, n := range news {
		nt := n.Datetime
		if nt > 0 && nt < 1_000_000_000_000 {
			nt *= 1000
		}
		if latest > 0 && nt > 0 && absInt64(latest-nt) > window {
			continue
		}
		hit := false
		for _, s := range n.Symbols {
			if normalizeSymbol(s) == c.Symbol {
				hit = true
				break
			}
		}
		if !hit && strings.Contains(strings.ToUpper(n.Related), c.Symbol) {
			hit = true
		}
		if hit {
			c.Corroborated = true
			c.Corroboration = appendUniqueString(c.Corroboration, "NEWS · "+defaultString(n.Source, "canonical news"))
		}
	}
	for _, f := range filings {
		if normalizeSymbol(f.Symbol) == c.Symbol {
			c.Corroborated = true
			c.Corroboration = appendUniqueString(c.Corroboration, "SEC · "+defaultString(f.Form, "filing"))
		}
	}
}

func absInt64(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

func buildCommunityEvidenceFusion(items []CommunityEvidenceItem, news []NewsItem, filings []FilingItem, now time.Time) CommunityIntelligenceState {
	out := CommunityIntelligenceState{Label: "UNTRUSTED COMMUNITY INTELLIGENCE · EVIDENCE FUSION", State: "UNAVAILABLE", Items: []CommunityEvidenceItem{}, Clusters: []CommunityEvidenceCluster{}, Policies: allCommunitySourcePolicies(), Untrusted: true, DeterministicImpact: "NONE", Detail: "No permitted/user-authorized community evidence has been recorded."}
	cutoff := now.Add(-30 * 24 * time.Hour).UnixMilli()
	platforms := map[string]bool{}
	recent1h := now.Add(-time.Hour).UnixMilli()
	for _, raw := range items {
		if raw.SubmittedAt <= 0 || raw.SubmittedAt < cutoff {
			continue
		}
		x := normalizeCommunityEvidenceItem(raw)
		if x.RetentionPolicy == "NO_RETENTION" {
			continue
		}
		out.Items = append(out.Items, x)
		platforms[x.Platform] = true
		switch strings.ToUpper(x.Stance) {
		case "BULLISH":
			out.Bullish++
		case "BEARISH":
			out.Bearish++
		case "MIXED":
			out.Mixed++
		default:
			out.Unknown++
		}
		if x.AIEligibility == "AI_ALLOWED" {
			out.AIEligible++
		}
		if x.ObservedAt >= recent1h || x.SubmittedAt >= recent1h {
			out.MentionVelocity1H++
		}
		if x.SubmittedAt > out.UpdatedAt {
			out.UpdatedAt = x.SubmittedAt
		}
	}
	sort.Slice(out.Items, func(i, j int) bool { return out.Items[i].SubmittedAt > out.Items[j].SubmittedAt })
	out.Total = len(out.Items)
	out.SourceDiversity = len(platforms)
	// Cluster exact and near-duplicate narratives within a six-hour window, so reposts do not become independent evidence.
	for _, x := range out.Items {
		at := x.ObservedAt
		if at <= 0 {
			at = x.SubmittedAt
		}
		idx := -1
		for i := range out.Clusters {
			c := &out.Clusters[i]
			if c.Symbol != "" && x.Symbol != "" && c.Symbol != x.Symbol {
				continue
			}
			if absInt64(c.LatestAt-at) > int64(6*time.Hour/time.Millisecond) {
				continue
			}
			sameURL := false
			if x.URL != "" {
				for _, u := range c.URLs {
					if strings.EqualFold(strings.TrimSpace(u), strings.TrimSpace(x.URL)) {
						sameURL = true
						break
					}
				}
			}
			if c.ID == "cluster-"+x.Fingerprint || communityTokenSimilarity(c.Representative, x.Text) >= 0.72 || sameURL {
				idx = i
				break
			}
		}
		if idx < 0 {
			out.Clusters = append(out.Clusters, CommunityEvidenceCluster{ID: "cluster-" + x.Fingerprint, Symbol: x.Symbol, Representative: x.Text, LatestAt: at})
			idx = len(out.Clusters) - 1
		}
		c := &out.Clusters[idx]
		c.Mentions++
		c.Sources = appendUniqueString(c.Sources, x.Source)
		c.Platforms = appendUniqueString(c.Platforms, x.Platform)
		if x.URL != "" {
			c.URLs = appendUniqueString(c.URLs, x.URL)
		}
		c.SourceDiversity = len(c.Platforms)
		if at >= recent1h {
			c.MentionVelocity1H++
		}
		if at > c.LatestAt {
			c.LatestAt = at
		}
		if x.AIEligibility == "AI_ALLOWED" {
			c.AIEligibleItems++
		}
		switch strings.ToUpper(x.Stance) {
		case "BULLISH":
			c.Bullish++
		case "BEARISH":
			c.Bearish++
		case "MIXED":
			c.Mixed++
		default:
			c.Unknown++
		}
	}
	for i := range out.Clusters {
		c := &out.Clusters[i]
		corroborateCommunityCluster(c, news, filings, c.LatestAt)
		switch {
		case c.Corroborated && c.SourceDiversity >= 3:
			c.Materiality = "HIGH"
		case c.Corroborated || c.SourceDiversity >= 2 || c.MentionVelocity1H >= 3:
			c.Materiality = "ELEVATED"
		default:
			c.Materiality = "NORMAL"
		}
		c.Detail = fmt.Sprintf("%d mention(s) across %d platform/source group(s); 1h velocity %d; corroborated %t. Community evidence remains untrusted and contextual.", c.Mentions, c.SourceDiversity, c.MentionVelocity1H, c.Corroborated)
	}
	sort.Slice(out.Clusters, func(i, j int) bool {
		rank := func(s string) int {
			switch s {
			case "HIGH":
				return 3
			case "ELEVATED":
				return 2
			default:
				return 1
			}
		}
		ri, rj := rank(out.Clusters[i].Materiality), rank(out.Clusters[j].Materiality)
		if ri != rj {
			return ri > rj
		}
		return out.Clusters[i].LatestAt > out.Clusters[j].LatestAt
	})
	if len(out.Clusters) > 30 {
		out.Clusters = out.Clusters[:30]
	}
	out.UniqueEvents = len(out.Clusters)
	if len(out.Items) > 50 {
		out.Items = out.Items[:50]
	}
	if out.Total > 0 {
		out.State = "AVAILABLE · UNTRUSTED"
		out.Detail = fmt.Sprintf("%d permitted/user-authorized observations fused into %d canonical event cluster(s) across %d source platform(s); %d item(s) are AI-eligible under source policy. Reposts are not counted as independent confirmation.", out.Total, out.UniqueEvents, out.SourceDiversity, out.AIEligible)
	}
	return out
}
