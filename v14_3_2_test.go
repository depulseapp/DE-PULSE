package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func runManualActionHandler(t *testing.T, fn http.HandlerFunc) map[string]any {
	t.Helper()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	fn(rr, req)
	if rr.Code < 200 || rr.Code >= 300 {
		t.Fatalf("manual action returned %d: %s", rr.Code, rr.Body.String())
	}
	var out map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode manual action response: %v", err)
	}
	if ok, _ := out["ok"].(bool); !ok {
		t.Fatalf("manual action did not report ok: %#v", out)
	}
	return out
}

func TestV1432DataEngineManualActionsReuseProductionPipelinesInDemo(t *testing.T) {
	app := newTestApplication(t)
	app.engine.mu.Lock()
	app.engine.status = "running"
	app.engine.mode = "demo"
	app.engine.mu.Unlock()

	catalyst := runManualActionHandler(t, app.handleCatalystEvaluate)
	if msg, _ := catalyst["message"].(string); !strings.Contains(strings.ToLower(msg), "catalyst") && !strings.Contains(strings.ToLower(msg), "earnings") {
		t.Fatalf("unexpected catalyst message: %#v", catalyst)
	}
	global := runManualActionHandler(t, app.handleGlobalRefresh)
	if msg, _ := global["message"].(string); !strings.Contains(msg, "Demo") {
		t.Fatalf("demo global refresh should not call providers: %#v", global)
	}
	capabilities := runManualActionHandler(t, app.handleCapabilityRecheck)
	if msg, _ := capabilities["message"].(string); !strings.Contains(msg, "Demo") {
		t.Fatalf("demo capability recheck should be explicit: %#v", capabilities)
	}
	app.engine.mu.RLock()
	capHealth := app.engine.health["provider-capabilities"]
	app.engine.mu.RUnlock()
	if !strings.Contains(strings.ToLower(capHealth), "demo") {
		t.Fatalf("capability recheck did not update health: %q", capHealth)
	}
	vix := runManualActionHandler(t, app.handleVIXRefresh)
	if msg, _ := vix["message"].(string); !strings.Contains(msg, "Demo") {
		t.Fatalf("demo VIX refresh should not call providers: %#v", vix)
	}
	stream := runManualActionHandler(t, app.handleStreamReconnect)
	if msg, _ := stream["message"].(string); !strings.Contains(msg, "Demo") {
		t.Fatalf("demo stream reconnect should not create provider loops: %#v", stream)
	}
}

func TestV1432ManualActionsRequireActiveRuntime(t *testing.T) {
	app := newTestApplication(t)
	handlers := []http.HandlerFunc{app.handleCatalystEvaluate, app.handleGlobalRefresh, app.handleCapabilityRecheck, app.handleVIXRefresh, app.handleStreamReconnect}
	for i, fn := range handlers {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{}"))
		fn(rr, req)
		if rr.Code != http.StatusConflict {
			t.Fatalf("handler %d should reject stopped runtime, got %d: %s", i, rr.Code, rr.Body.String())
		}
	}
}

func TestV1432StreamReconnectReusesExistingLoopWhenDisconnected(t *testing.T) {
	app := newTestApplication(t)
	app.engine.mu.Lock()
	app.engine.status = "degraded"
	app.engine.mode = "live"
	app.engine.webSocketConnected = false
	app.engine.ws = nil
	app.engine.mu.Unlock()
	out := runManualActionHandler(t, app.handleStreamReconnect)
	msg, _ := out["message"].(string)
	if !strings.Contains(strings.ToLower(msg), "already") || !strings.Contains(strings.ToLower(msg), "reconnecting") {
		t.Fatalf("expected existing reconnect-loop message, got %#v", out)
	}
}
