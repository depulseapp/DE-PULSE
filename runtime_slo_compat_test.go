package main

import "time"

func buildRuntimeSLOAssessment(status, mode string, feed FeedDiagnostics, freshness []FreshnessDiagnostic, load RuntimeLoadDiagnostics, scanner ScannerState) RuntimeSLOAssessment {
	return buildRuntimeSLOAssessmentWithContext(status, mode, feed, freshness, load, scanner, AppState{}, nil, time.Now())
}
