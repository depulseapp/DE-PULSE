package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

const tradeInsightCongressSchemaProbeFlag = "--tradeinsight-congress-schema-probe"

type tradeInsightCongressSchemaProbeReport struct {
	Capability  string                        `json:"capability"`
	Lifecycle   string                        `json:"lifecycle"`
	Fingerprint tradeInsightSchemaFingerprint `json:"fingerprint"`
}

type tradeInsightCongressProbeFunc func(context.Context) (tradeInsightSchemaFingerprint, error)

// runDeveloperCommand is deliberately checked before NewApplication in main.
// When handled, the process exits without opening persistence, binding the app
// server, focusing an existing desktop instance, or starting market services.
func runDeveloperCommand(args []string, stdout, stderr io.Writer) (bool, int) {
	probe := func(ctx context.Context) (tradeInsightSchemaFingerprint, error) {
		return tradeInsightCongressSchemaProbeAt(ctx, nil, tradeInsightRESTBaseURL, tradeInsightAPIKey())
	}
	return runDeveloperCommandWithProbe(args, stdout, stderr, probe)
}

func runDeveloperCommandWithProbe(args []string, stdout, stderr io.Writer, probe tradeInsightCongressProbeFunc) (bool, int) {
	if len(args) == 0 || args[0] != tradeInsightCongressSchemaProbeFlag {
		return false, 0
	}
	if len(args) != 1 {
		_, _ = fmt.Fprintf(stderr, "%s does not accept additional arguments\n", tradeInsightCongressSchemaProbeFlag)
		return true, 2
	}
	if probe == nil {
		_, _ = fmt.Fprintln(stderr, "TradeInsight Congress schema probe unavailable")
		return true, 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	fingerprint, err := probe(ctx)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "TradeInsight Congress schema probe failed: %v\n", err)
		return true, 1
	}

	report := tradeInsightCongressSchemaProbeReport{
		Capability:  "congressional-trades",
		Lifecycle:   tradeInsightCapabilityLifecycleTruth("congressional-trades"),
		Fingerprint: fingerprint,
	}
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		_, _ = fmt.Fprintf(stderr, "TradeInsight Congress schema fingerprint output failed: %v\n", err)
		return true, 1
	}
	return true, 0
}
