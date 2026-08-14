package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// signalSnapshotPersistenceBatch turns the already-existing frozen validation
// snapshot into the v17 evidence -> decision -> outcome learning path. This
// deliberately reuses SignalSnapshot as the canonical decision lineage owner;
// persistence does not create a second scoring/readiness truth path.
func signalSnapshotPersistenceBatch(s SignalSnapshot) PersistenceIntelligenceBatch {
	if strings.TrimSpace(s.ID) == "" || normalizeSymbol(s.Symbol) == "" {
		return PersistenceIntelligenceBatch{}
	}
	s.Symbol = normalizeSymbol(s.Symbol)
	decisionPayload := struct {
		ID                  string             `json:"id"`
		Symbol              string             `json:"symbol"`
		Horizon             string             `json:"horizon"`
		Timestamp           int64              `json:"timestamp"`
		Price               float64            `json:"price"`
		Score               float64            `json:"score"`
		Action              string             `json:"action"`
		Readiness           string             `json:"readiness"`
		EvidenceSnapshotID  string             `json:"evidenceSnapshotId,omitempty"`
		FormulaVersion      string             `json:"formulaVersion,omitempty"`
		SettingsFingerprint string             `json:"settingsFingerprint,omitempty"`
		SignalProfile       string             `json:"signalProfile,omitempty"`
		EarningsPenalty     float64            `json:"earningsPenalty,omitempty"`
		EarningsDays        *int               `json:"earningsDays,omitempty"`
		EntryLow            float64            `json:"entryLow,omitempty"`
		EntryHigh           float64            `json:"entryHigh,omitempty"`
		TargetLow           float64            `json:"targetLow,omitempty"`
		TargetHigh          float64            `json:"targetHigh,omitempty"`
		Invalidation        float64            `json:"invalidation,omitempty"`
		FamilyScores        map[string]float64 `json:"familyScores,omitempty"`
	}{ID: s.ID, Symbol: s.Symbol, Horizon: s.Horizon, Timestamp: s.Timestamp, Price: s.Price, Score: s.Score, Action: s.Action, Readiness: s.Readiness, EvidenceSnapshotID: s.EvidenceSnapshotID, FormulaVersion: s.FormulaVersion, SettingsFingerprint: s.SettingsFingerprint, SignalProfile: s.SignalProfile, EarningsPenalty: s.EarningsPenalty, EarningsDays: s.EarningsDays, EntryLow: s.EntryLow, EntryHigh: s.EntryHigh, TargetLow: s.TargetLow, TargetHigh: s.TargetHigh, Invalidation: s.Invalidation, FamilyScores: clone(s.FamilyScores)}
	decisionRaw, err := json.Marshal(decisionPayload)
	if err != nil {
		return PersistenceIntelligenceBatch{}
	}
	evidenceID := strings.TrimSpace(s.EvidenceSnapshotID)
	if evidenceID == "" {
		evidenceID = "signal-evidence:" + s.ID
	}
	evidencePayload := struct {
		EvidenceSnapshotID  string             `json:"evidenceSnapshotId,omitempty"`
		Symbol              string             `json:"symbol"`
		Horizon             string             `json:"horizon"`
		Timestamp           int64              `json:"timestamp"`
		Price               float64            `json:"price"`
		Score               float64            `json:"score"`
		Readiness           string             `json:"readiness"`
		FormulaVersion      string             `json:"formulaVersion,omitempty"`
		SettingsFingerprint string             `json:"settingsFingerprint,omitempty"`
		FamilyScores        map[string]float64 `json:"familyScores,omitempty"`
		MarketRegime        string             `json:"marketRegime,omitempty"`
		MarketStructure     string             `json:"marketStructure,omitempty"`
		MarketTradeability  string             `json:"marketTradeability,omitempty"`
		RelativeStrength    string             `json:"relativeStrength,omitempty"`
		SectorRegime        string             `json:"sectorRegime,omitempty"`
		LiquidityState      string             `json:"liquidityState,omitempty"`
		GlobalContext       string             `json:"globalContext,omitempty"`
		OptionsBias         string             `json:"optionsBias,omitempty"`
		EventRisk           string             `json:"eventRisk,omitempty"`
		ResearchState       string             `json:"researchState,omitempty"`
		QueuePriority       string             `json:"queuePriority,omitempty"`
		KeyDriver           string             `json:"keyDriver,omitempty"`
		Contradictions      []string           `json:"contradictions,omitempty"`
	}{
		EvidenceSnapshotID: s.EvidenceSnapshotID, Symbol: s.Symbol, Horizon: s.Horizon,
		Timestamp: s.Timestamp, Price: s.Price, Score: s.Score, Readiness: s.Readiness,
		FormulaVersion: s.FormulaVersion, SettingsFingerprint: s.SettingsFingerprint,
		FamilyScores: clone(s.FamilyScores), MarketRegime: s.MarketRegime,
		MarketStructure: s.MarketStructure, MarketTradeability: s.MarketTradeability,
		RelativeStrength: s.RelativeStrength, SectorRegime: s.SectorRegime,
		LiquidityState: s.LiquidityState, GlobalContext: s.GlobalContext, OptionsBias: s.OptionsBias,
		EventRisk: s.EventRisk, ResearchState: s.ResearchState, QueuePriority: s.QueuePriority,
		KeyDriver: s.KeyDriver, Contradictions: append([]string(nil), s.Contradictions...),
	}
	evidenceRaw, err := json.Marshal(evidencePayload)
	if err != nil {
		return PersistenceIntelligenceBatch{}
	}
	batch := PersistenceIntelligenceBatch{
		Evidence: []EvidenceRecord{{
			ID: evidenceID, Symbol: s.Symbol, Kind: "frozen-signal-evidence", ObservedAt: s.Timestamp,
			Source: "canonical-signal-validation", Provenance: s.FormulaVersion, FreshnessState: "FROZEN",
			Payload: evidenceRaw,
		}},
		Decisions: []DecisionLineageRecord{{
			ID: s.ID, Symbol: s.Symbol, Horizon: s.Horizon, EvidenceID: evidenceID,
			DecisionKind: "deterministic-signal", DecisionValue: strings.TrimSpace(s.Action),
			FormulaVersion: s.FormulaVersion, CreatedAt: s.Timestamp, Payload: decisionRaw,
		}},
	}

	if s.OutcomeState != "" || len(s.Outcomes) > 0 || s.OutcomeUpdatedAt > 0 {
		observed := s.OutcomeUpdatedAt
		if observed <= 0 {
			observed = s.Timestamp
		}
		outcomeRaw, _ := json.Marshal(struct {
			State       string             `json:"state,omitempty"`
			Detail      string             `json:"detail,omitempty"`
			Outcomes    map[string]float64 `json:"outcomes,omitempty"`
			MFE         float64            `json:"mfe,omitempty"`
			MAE         float64            `json:"mae,omitempty"`
			EntryAt     int64              `json:"entryTouchedAt,omitempty"`
			TargetAt    int64              `json:"targetTouchedAt,omitempty"`
			InvalidAt   int64              `json:"invalidationAt,omitempty"`
			ElapsedMins int64              `json:"elapsedMinutes,omitempty"`
		}{s.OutcomeState, s.OutcomeDetail, clone(s.Outcomes), s.MFE, s.MAE, s.EntryTouchedAt, s.TargetTouchedAt, s.InvalidationAt, s.ElapsedMinutes})
		batch.Outcomes = append(batch.Outcomes, OutcomeHistoryRecord{
			ID: fmt.Sprintf("%s:%d", s.ID, observed), DecisionID: s.ID, Symbol: s.Symbol,
			Horizon: s.Horizon, ObservedAt: observed, OutcomeLabel: s.OutcomeState, Payload: outcomeRaw,
		})
	}

	if len(s.FamilyScores) > 0 {
		keys := make([]string, 0, len(s.FamilyScores))
		for key := range s.FamilyScores {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		featureRaw, _ := json.Marshal(struct {
			Scores       map[string]float64 `json:"scores"`
			Readiness    string             `json:"readiness,omitempty"`
			Regime       string             `json:"marketRegime,omitempty"`
			Tradeability string             `json:"marketTradeability,omitempty"`
		}{clone(s.FamilyScores), s.Readiness, s.MarketRegime, s.MarketTradeability})
		h := sha256.New()
		for _, key := range keys {
			fmt.Fprintf(h, "%s=%g;", key, s.FamilyScores[key])
		}
		fmt.Fprintf(h, "readiness=%s;regime=%s;tradeability=%s", s.Readiness, s.MarketRegime, s.MarketTradeability)
		batch.Features = append(batch.Features, DerivedFeatureRecord{
			Symbol: s.Symbol, FeatureKey: "signal-families:" + strings.ToLower(strings.TrimSpace(s.Horizon)),
			FeatureVersion: s.FormulaVersion, AsOf: s.Timestamp, SourceHash: hex.EncodeToString(h.Sum(nil))[:24], Payload: featureRaw,
		})
	}
	return batch
}

func (e *Engine) enqueueSignalSnapshotPersistence(s SignalSnapshot) {
	if e == nil || e.app == nil || e.app.persistence == nil {
		return
	}
	batch := signalSnapshotPersistenceBatch(s)
	if batch.Len() > 0 {
		e.app.persistence.EnqueueIntelligence(batch)
	}
}

func payloadOrEmpty(raw json.RawMessage) []byte {
	if len(raw) == 0 {
		return []byte("{}")
	}
	return raw
}
