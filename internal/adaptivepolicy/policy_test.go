package adaptivepolicy

import (
	"testing"
	"time"
)

func TestIntradayHistoryCadence(t *testing.T) {
	cases := []struct {
		name    string
		session string
		hot     bool
		want    time.Duration
	}{
		{"regular hot", "regular", true, 2 * time.Minute},
		{"regular normal", "regular", false, 5 * time.Minute},
		{"premarket hot", "pre-market", true, 3 * time.Minute},
		{"after hours normal", "after-hours", false, 5 * time.Minute},
		{"overnight", "overnight", true, 15 * time.Minute},
		{"unknown", "closed", false, 30 * time.Minute},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IntradayHistoryCadence(tc.session, tc.hot); got != tc.want {
				t.Fatalf("got %s, want %s", got, tc.want)
			}
		})
	}
}

func TestCachePersistCadence(t *testing.T) {
	cases := []struct {
		session string
		hot     bool
		want    time.Duration
	}{
		{"regular", true, time.Minute},
		{"regular", false, 2 * time.Minute},
		{"pre-market", true, time.Minute},
		{"after-hours", false, 2 * time.Minute},
		{"overnight", true, 5 * time.Minute},
		{"closed", false, 10 * time.Minute},
	}
	for _, tc := range cases {
		if got := CachePersistCadence(tc.session, tc.hot); got != tc.want {
			t.Fatalf("session=%s hot=%v: got %s, want %s", tc.session, tc.hot, got, tc.want)
		}
	}
}

func TestProviderDegraded(t *testing.T) {
	if ProviderDegraded(map[string]string{"alpaca-live": "healthy", "alpaca-stream": "connected", "quotes": "fresh"}) {
		t.Fatal("healthy provider state classified as degraded")
	}
	for _, state := range []string{"ERROR", "failed", "Reconnecting", "degraded"} {
		if !ProviderDegraded(map[string]string{"quotes": state}) {
			t.Fatalf("state %q must classify as degraded", state)
		}
	}
}
