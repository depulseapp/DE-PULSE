package main

import (
	"strings"
	"testing"
)

func TestV18_9_1MacOSWindowScriptAvoidsFormalDelegateProtocolLookup(t *testing.T) {
	script := macOSWindowScript("http://127.0.0.1:43123/o'hare", "/tmp/de'pulse.icns")

	if strings.Contains(script, "NSApplicationDelegate") {
		t.Fatalf("macOS JXA script must not formally resolve NSApplicationDelegate after ADAPT-RUNTIME-CRASH-001")
	}
	for _, required := range []string{
		"ObjC.import('Cocoa')",
		"ObjC.import('WebKit')",
		"ObjC.registerSubclass({",
		"superclass:'NSObject'",
		"applicationShouldTerminateAfterLastWindowClosed:",
		"applicationShouldHandleReopen:hasVisibleWindows:",
		"app.delegate=delegate",
		"WKWebViewConfiguration",
		"app.run",
	} {
		if !strings.Contains(script, required) {
			t.Fatalf("macOS JXA script lost required native lifecycle behavior %q", required)
		}
	}
	if !strings.Contains(script, "http://127.0.0.1:43123/o%27hare") {
		t.Fatalf("raw URL quote was not safely encoded in JXA script: %s", script)
	}
	if !strings.Contains(script, "/tmp/de\\'pulse.icns") {
		t.Fatalf("icon path quote was not safely escaped in JXA script: %s", script)
	}
}
