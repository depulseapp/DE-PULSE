package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

const evidenceTemporalEnvelopeSchema = "DE.PULSE-EVIDENCE-TEMPORAL-1"

type evidenceTemporalStorageMetadata struct {
	TemporalSchema    string `json:"temporalSchema"`
	SourceAt          int64  `json:"sourceAt"`
	IngestedAt        int64  `json:"ingestedAt"`
	KnownAt           int64  `json:"knownAt"`
	EffectiveFrom     int64  `json:"effectiveFrom"`
	EffectiveTo       int64  `json:"effectiveTo,omitempty"`
	ReportPeriod      string `json:"reportPeriod,omitempty"`
	RevisionID        string `json:"revisionId"`
	SupersedesID      string `json:"supersedesId,omitempty"`
	AmendmentState    string `json:"amendmentState"`
	RightsState       string `json:"rightsState"`
	RightsEvidenceRef string `json:"rightsEvidenceRef,omitempty"`
	RetentionClass    string `json:"retentionClass"`
}

type evidenceTemporalStorageEnvelope struct {
	Schema   string                          `json:"_depulseSchema"`
	Temporal evidenceTemporalStorageMetadata `json:"temporal"`
	Payload  json.RawMessage                 `json:"payload"`
}

func normalizeEvidenceTemporalRecord(record EvidenceRecord, fallbackKnownAt int64) (EvidenceRecord, error) {
	record.ID = strings.TrimSpace(record.ID)
	record.Kind = strings.TrimSpace(record.Kind)
	record.Symbol = normalizeSymbol(record.Symbol)
	if record.ID == "" || record.Kind == "" {
		return record, fmt.Errorf("evidence id and kind are required")
	}
	if fallbackKnownAt <= 0 {
		fallbackKnownAt = time.Now().UnixMilli()
	}
	for name, value := range map[string]int64{
		"sourceAt": record.SourceAt, "observedAt": record.ObservedAt, "ingestedAt": record.IngestedAt,
		"knownAt": record.KnownAt, "effectiveFrom": record.EffectiveFrom, "effectiveTo": record.EffectiveTo,
	} {
		if value < 0 {
			return record, fmt.Errorf("evidence %s cannot be negative", name)
		}
	}
	if record.ObservedAt == 0 {
		record.ObservedAt = record.IngestedAt
	}
	if record.ObservedAt == 0 {
		record.ObservedAt = fallbackKnownAt
	}
	if record.IngestedAt == 0 {
		record.IngestedAt = record.ObservedAt
	}
	if record.KnownAt == 0 {
		record.KnownAt = maxInt64(record.ObservedAt, record.IngestedAt)
	}
	if record.SourceAt == 0 {
		record.SourceAt = record.ObservedAt
	}
	if record.EffectiveFrom == 0 {
		record.EffectiveFrom = record.SourceAt
	}
	if record.KnownAt < record.ObservedAt || record.KnownAt < record.IngestedAt {
		return record, fmt.Errorf("evidence knownAt must not precede observation or ingestion")
	}
	const maxSourceClockSkewMs = int64(30 * time.Second / time.Millisecond)
	if record.SourceAt > record.KnownAt+maxSourceClockSkewMs {
		return record, fmt.Errorf("evidence sourceAt is materially in the future of knownAt")
	}
	if record.EffectiveTo > 0 && record.EffectiveTo <= record.EffectiveFrom {
		return record, fmt.Errorf("evidence effectiveTo must follow effectiveFrom")
	}
	record.RevisionID = strings.TrimSpace(record.RevisionID)
	if record.RevisionID == "" {
		record.RevisionID = record.ID
	}
	record.SupersedesID = strings.TrimSpace(record.SupersedesID)
	if record.SupersedesID == record.ID {
		return record, fmt.Errorf("evidence cannot supersede itself")
	}
	record.AmendmentState = strings.ToUpper(strings.TrimSpace(record.AmendmentState))
	if record.AmendmentState == "" {
		if record.SupersedesID == "" {
			record.AmendmentState = "ORIGINAL"
		} else {
			record.AmendmentState = "REVISION"
		}
	}
	allowedAmendment := map[string]bool{"ORIGINAL": true, "REVISION": true, "AMENDMENT": true, "CORRECTION": true, "RESTATEMENT": true, "NOT_INFERRED_NO_VENDOR_ID": true, "UNKNOWN": true}
	if !allowedAmendment[record.AmendmentState] {
		return record, fmt.Errorf("unsupported evidence amendment state %q", record.AmendmentState)
	}
	if record.SupersedesID != "" && record.AmendmentState == "ORIGINAL" {
		return record, fmt.Errorf("original evidence cannot declare a superseded record")
	}
	record.RightsEvidenceRef = strings.TrimSpace(record.RightsEvidenceRef)
	if record.RightsEvidenceRef == "" {
		record.RightsState = "UNBOUND"
	} else {
		record.RightsState = "BOUND"
	}
	record.RetentionClass = strings.ToUpper(strings.TrimSpace(record.RetentionClass))
	if record.RetentionClass == "" {
		record.RetentionClass = "UNSPECIFIED"
	}
	record.TemporalSchema = evidenceTemporalEnvelopeSchema
	if !json.Valid(payloadOrEmpty(record.Payload)) {
		return record, fmt.Errorf("evidence payload is not valid JSON")
	}
	record.Payload = append(json.RawMessage(nil), payloadOrEmpty(record.Payload)...)
	return record, nil
}

func evidenceTemporalStoragePayload(record EvidenceRecord, fallbackKnownAt int64) (EvidenceRecord, []byte, error) {
	record, err := normalizeEvidenceTemporalRecord(record, fallbackKnownAt)
	if err != nil {
		return record, nil, err
	}
	envelope := evidenceTemporalStorageEnvelope{
		Schema: evidenceTemporalEnvelopeSchema,
		Temporal: evidenceTemporalStorageMetadata{
			TemporalSchema: record.TemporalSchema, SourceAt: record.SourceAt, IngestedAt: record.IngestedAt,
			KnownAt: record.KnownAt, EffectiveFrom: record.EffectiveFrom, EffectiveTo: record.EffectiveTo,
			ReportPeriod: record.ReportPeriod, RevisionID: record.RevisionID, SupersedesID: record.SupersedesID,
			AmendmentState: record.AmendmentState, RightsState: record.RightsState,
			RightsEvidenceRef: record.RightsEvidenceRef, RetentionClass: record.RetentionClass,
		},
		Payload: record.Payload,
	}
	raw, err := json.Marshal(envelope)
	if err != nil {
		return record, nil, err
	}
	return record, raw, nil
}

func evidenceRecordFromStorage(record EvidenceRecord, raw json.RawMessage) (EvidenceRecord, error) {
	var envelope evidenceTemporalStorageEnvelope
	if json.Unmarshal(raw, &envelope) == nil && envelope.Schema == evidenceTemporalEnvelopeSchema {
		record.TemporalSchema = envelope.Temporal.TemporalSchema
		record.SourceAt = envelope.Temporal.SourceAt
		record.IngestedAt = envelope.Temporal.IngestedAt
		record.KnownAt = envelope.Temporal.KnownAt
		record.EffectiveFrom = envelope.Temporal.EffectiveFrom
		record.EffectiveTo = envelope.Temporal.EffectiveTo
		record.ReportPeriod = envelope.Temporal.ReportPeriod
		record.RevisionID = envelope.Temporal.RevisionID
		record.SupersedesID = envelope.Temporal.SupersedesID
		record.AmendmentState = envelope.Temporal.AmendmentState
		record.RightsState = envelope.Temporal.RightsState
		record.RightsEvidenceRef = envelope.Temporal.RightsEvidenceRef
		record.RetentionClass = envelope.Temporal.RetentionClass
		record.Payload = append(json.RawMessage(nil), envelope.Payload...)
		return normalizeEvidenceTemporalRecord(record, record.ObservedAt)
	}
	record.Payload = append(json.RawMessage(nil), raw...)
	return normalizeEvidenceTemporalRecord(record, record.ObservedAt)
}

// EvidencePointInTime reconstructs the evidence knowable at knownAt. When
// effectiveAt is non-zero it also evaluates the valid/effective-time interval.
// A revision only supersedes an older record after that revision itself was
// knowable, which prevents later corrections from leaking into earlier replay.
func EvidencePointInTime(records []EvidenceRecord, knownAt, effectiveAt int64) ([]EvidenceRecord, error) {
	if knownAt <= 0 {
		return nil, fmt.Errorf("point-in-time knownAt cutoff is required")
	}
	all := make(map[string]EvidenceRecord, len(records))
	eligible := make(map[string]EvidenceRecord, len(records))
	for _, candidate := range records {
		record, err := normalizeEvidenceTemporalRecord(candidate, candidate.ObservedAt)
		if err != nil {
			return nil, fmt.Errorf("evidence %q: %w", candidate.ID, err)
		}
		if _, exists := all[record.ID]; exists {
			return nil, fmt.Errorf("duplicate evidence id %q", record.ID)
		}
		all[record.ID] = record
		if record.KnownAt > knownAt {
			continue
		}
		if effectiveAt > 0 && (record.EffectiveFrom > effectiveAt || (record.EffectiveTo > 0 && effectiveAt >= record.EffectiveTo)) {
			continue
		}
		eligible[record.ID] = record
	}
	for id := range all {
		seen := map[string]bool{id: true}
		current := all[id]
		for current.SupersedesID != "" {
			if seen[current.SupersedesID] {
				return nil, fmt.Errorf("evidence revision cycle involving %q", current.SupersedesID)
			}
			seen[current.SupersedesID] = true
			next, ok := all[current.SupersedesID]
			if !ok {
				break
			}
			current = next
		}
	}
	superseded := map[string]bool{}
	for _, record := range eligible {
		if record.SupersedesID != "" {
			if _, present := eligible[record.SupersedesID]; present {
				superseded[record.SupersedesID] = true
			}
		}
	}
	out := make([]EvidenceRecord, 0, len(eligible))
	for id, record := range eligible {
		if !superseded[id] {
			out = append(out, record)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		left := out[i].Symbol + "|" + out[i].Kind + "|" + fmt.Sprintf("%020d", out[i].EffectiveFrom) + "|" + out[i].RevisionID + "|" + out[i].ID
		right := out[j].Symbol + "|" + out[j].Kind + "|" + fmt.Sprintf("%020d", out[j].EffectiveFrom) + "|" + out[j].RevisionID + "|" + out[j].ID
		return left < right
	})
	return out, nil
}

func EvidenceAsKnownAt(records []EvidenceRecord, knownAt int64) ([]EvidenceRecord, error) {
	return EvidencePointInTime(records, knownAt, 0)
}

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
			SourceAt: s.Timestamp, IngestedAt: s.Timestamp, KnownAt: s.Timestamp, EffectiveFrom: s.Timestamp,
			RevisionID: evidenceID, AmendmentState: "ORIGINAL", RightsEvidenceRef: "internal:canonical-signal-validation", RetentionClass: "DECISION_LINEAGE",
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
