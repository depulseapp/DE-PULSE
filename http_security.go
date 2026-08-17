package main

import (
	"crypto/sha256"
	"encoding/base64"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const browserContentSecurityPolicy = "default-src 'self'; img-src 'self' https: data:; style-src 'self' 'unsafe-inline'; script-src 'self'; connect-src 'self' https://finnhub.io https://openrouter.ai https://api.groq.com https://generativelanguage.googleapis.com https://www.sec.gov https://data.sec.gov https://data.alpaca.markets https://api.stlouisfed.org https://api.bls.gov https://api.eia.gov https://api.twelvedata.com https://api.marketaux.com https://query1.finance.yahoo.com https://www.cboe.com https://home.treasury.gov https://www.bea.gov wss://ws.finnhub.io;"

func firstForwardedValue(v string) string {
	if i := strings.IndexByte(v, ','); i >= 0 {
		v = v[:i]
	}
	return strings.TrimSpace(v)
}

func effectiveRequestScheme(r *http.Request) string {
	if r == nil {
		return "http"
	}
	if r.TLS != nil {
		return "https"
	}
	if trustHostedProxyHeaders() {
		switch strings.ToLower(firstForwardedValue(r.Header.Get("X-Forwarded-Proto"))) {
		case "https":
			return "https"
		case "http":
			return "http"
		}
	}
	if r.URL != nil {
		switch strings.ToLower(strings.TrimSpace(r.URL.Scheme)) {
		case "https":
			return "https"
		case "http":
			return "http"
		}
	}
	return "http"
}

func validForwardedHost(v string) bool {
	v = strings.TrimSpace(v)
	return v != "" && !strings.ContainsAny(v, "/\\ 	\r\n")
}

func effectiveRequestHost(r *http.Request) string {
	if r == nil {
		return ""
	}
	host := strings.TrimSpace(r.Host)
	if trustHostedProxyHeaders() {
		forwarded := firstForwardedValue(r.Header.Get("X-Forwarded-Host"))
		if validForwardedHost(forwarded) {
			host = forwarded
		}
	}
	return strings.ToLower(host)
}

func requestIsSecure(r *http.Request) bool { return effectiveRequestScheme(r) == "https" }

func directClientAddress(r *http.Request) string {
	if r == nil {
		return "unknown"
	}
	remote := strings.TrimSpace(r.RemoteAddr)
	if host, _, err := net.SplitHostPort(remote); err == nil && host != "" {
		return strings.ToLower(host)
	}
	if ip := net.ParseIP(strings.Trim(remote, "[]")); ip != nil {
		return ip.String()
	}
	if remote == "" {
		return "unknown"
	}
	return strings.ToLower(remote)
}

func effectiveClientAddress(r *http.Request) string {
	if r != nil && trustHostedProxyHeaders() {
		if candidate := firstForwardedValue(r.Header.Get("X-Forwarded-For")); net.ParseIP(strings.TrimSpace(candidate)) != nil {
			return net.ParseIP(strings.TrimSpace(candidate)).String()
		}
	}
	return directClientAddress(r)
}

func normalizedOrigin(v string) (string, bool) {
	v = strings.TrimSpace(v)
	if v == "" || strings.EqualFold(v, "null") {
		return "", false
	}
	u, err := url.Parse(v)
	if err != nil || u.User != nil || u.Host == "" || (u.Path != "" && u.Path != "/") || u.RawQuery != "" || u.Fragment != "" {
		return "", false
	}
	scheme := strings.ToLower(strings.TrimSpace(u.Scheme))
	if scheme != "http" && scheme != "https" {
		return "", false
	}
	return scheme + "://" + strings.ToLower(u.Host), true
}

func sameOriginAllowed(r *http.Request) bool {
	if r == nil {
		return false
	}
	switch r.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	}
	originHeader := strings.TrimSpace(r.Header.Get("Origin"))
	if originHeader == "" {
		// Native/local clients do not necessarily emit Origin; CSRF still protects
		// authenticated mutations. Explicit foreign browser origins are fail-closed.
		return true
	}
	origin, ok := normalizedOrigin(originHeader)
	if !ok {
		return false
	}
	if configured := hostedPublicOrigin(); configured != "" {
		publicOrigin, valid := normalizedOrigin(configured)
		return valid && origin == publicOrigin
	}
	host := effectiveRequestHost(r)
	if host == "" {
		return false
	}
	return origin == effectiveRequestScheme(r)+"://"+host
}

func securityPerimeter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=(), usb=()")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set("Content-Security-Policy", browserContentSecurityPolicy)
		if isHostedRuntime() && requestIsSecure(r) {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000")
		}
		if !sameOriginAllowed(r) {
			writeError(w, http.StatusForbidden, "Cross-origin request rejected.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

type loginAbuseEntry struct {
	Failures     int
	WindowStart  time.Time
	LastSeen     time.Time
	BlockedUntil time.Time
}

type loginAbuseLimiter struct {
	mu         sync.Mutex
	entries    map[string]loginAbuseEntry
	maxEntries int
	limit      int
	window     time.Duration
	block      time.Duration
	now        func() time.Time
}

func newLoginAbuseLimiter(maxEntries, limit int, window, block time.Duration) *loginAbuseLimiter {
	return &loginAbuseLimiter{entries: map[string]loginAbuseEntry{}, maxEntries: maxEntries, limit: limit, window: window, block: block, now: time.Now}
}

func loginAbuseKey(username string, r *http.Request) string {
	raw := normalizeUsername(username) + "\x00" + effectiveClientAddress(r)
	sum := sha256.Sum256([]byte(raw))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func (l *loginAbuseLimiter) cleanupLocked(now time.Time) {
	for key, entry := range l.entries {
		if now.Sub(entry.LastSeen) > l.window+l.block {
			delete(l.entries, key)
		}
	}
	for len(l.entries) >= l.maxEntries && len(l.entries) > 0 {
		oldestKey := ""
		var oldest time.Time
		for key, entry := range l.entries {
			if oldestKey == "" || entry.LastSeen.Before(oldest) {
				oldestKey, oldest = key, entry.LastSeen
			}
		}
		delete(l.entries, oldestKey)
	}
}

func (l *loginAbuseLimiter) Allow(key string) (bool, int) {
	now := l.now()
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cleanupLocked(now)
	entry, ok := l.entries[key]
	if !ok {
		return true, 0
	}
	entry.LastSeen = now
	l.entries[key] = entry
	if now.Before(entry.BlockedUntil) {
		retry := int((entry.BlockedUntil.Sub(now) + time.Second - 1) / time.Second)
		if retry < 1 {
			retry = 1
		}
		return false, retry
	}
	if now.Sub(entry.WindowStart) >= l.window {
		delete(l.entries, key)
	}
	return true, 0
}

func (l *loginAbuseLimiter) RecordFailure(key string) {
	now := l.now()
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cleanupLocked(now)
	entry, ok := l.entries[key]
	if !ok || now.Sub(entry.WindowStart) >= l.window {
		entry = loginAbuseEntry{WindowStart: now}
	}
	entry.Failures++
	entry.LastSeen = now
	if entry.Failures >= l.limit {
		entry.BlockedUntil = now.Add(l.block)
	}
	l.entries[key] = entry
}

func (l *loginAbuseLimiter) Reset(key string) {
	l.mu.Lock()
	delete(l.entries, key)
	l.mu.Unlock()
}

var loginLimiter = newLoginAbuseLimiter(4096, 5, 5*time.Minute, 5*time.Minute)

func rejectThrottledLogin(w http.ResponseWriter, retryAfter int) {
	w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
	writeError(w, http.StatusTooManyRequests, "Too many login attempts. Try again later.")
}
