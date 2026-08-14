package main

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}
func normalizeSymbol(v string) string {
	v = strings.ToUpper(strings.TrimSpace(v))
	var b strings.Builder
	for _, r := range v {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '^' || r == '-' {
			b.WriteRune(r)
		}
	}
	s := b.String()
	if len(s) > 15 {
		s = s[:15]
	}
	return s
}
func uniqueSymbols(items []string) []string {
	if items == nil {
		return nil
	}
	seen := map[string]bool{}
	out := make([]string, 0, len(items))
	for _, item := range items {
		s := normalizeSymbol(item)
		if s != "" && !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

func userTradingSymbols(items []string) []string {
	if items == nil {
		return nil
	}
	seen := map[string]bool{}
	out := make([]string, 0, len(items))
	for _, item := range items {
		s, ok := parseUserTicker(item)
		if ok && !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
func randomID(prefix string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return prefix + "_" + base64.RawURLEncoding.EncodeToString(b)
}
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
