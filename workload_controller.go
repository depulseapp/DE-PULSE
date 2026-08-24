package main

import (
	"context"
	"sync"
	"time"
)

type WorkTier int

const (
	WorkTierMarketCritical WorkTier = iota
	WorkTierUserActionable
	WorkTierRadarPromoted
	WorkTierBroadDiscovery
	WorkTierBackground
)

type workTierContextKey struct{}

func withWorkTier(ctx context.Context, tier WorkTier) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, workTierContextKey{}, normalizeWorkTier(tier))
}

func workTierFromContext(ctx context.Context, fallback WorkTier) WorkTier {
	if ctx != nil {
		if tier, ok := ctx.Value(workTierContextKey{}).(WorkTier); ok {
			return normalizeWorkTier(tier)
		}
	}
	return normalizeWorkTier(fallback)
}

func normalizeWorkTier(t WorkTier) WorkTier {
	if t < WorkTierMarketCritical {
		return WorkTierMarketCritical
	}
	if t > WorkTierBackground {
		return WorkTierBackground
	}
	return t
}

type WorkClassDiagnostics struct {
	Class                 string   `json:"class"`
	Capacity              int      `json:"capacity"`
	MaxQueue              int      `json:"maxQueue"`
	ReservedCritical      int      `json:"reservedCritical"`
	ReservedCriticalQueue int      `json:"reservedCriticalQueue,omitempty"`
	InFlight              int      `json:"inFlight"`
	Queued                int      `json:"queued"`
	OldestQueueAgeMs      int64    `json:"oldestQueueAgeMs,omitempty"`
	Completed             int64    `json:"completed"`
	Canceled              int64    `json:"canceled"`
	Rejected              int64    `json:"rejected"`
	Shed                  int64    `json:"shed"`
	QueuedByTier          [5]int   `json:"queuedByTier"`
	InFlightByTier        [5]int   `json:"inFlightByTier"`
	RejectedByTier        [5]int64 `json:"rejectedByTier"`
}

type workClassState struct {
	name             string
	capacity         int
	maxQueue         int
	reservedCritical int
	inFlight         int
	queued           [5]int
	inFlightByTier   [5]int
	oldestQueue      [5]time.Time
	completed        int64
	canceled         int64
	rejected         int64
	rejectedByTier   [5]int64
	shed             int64
	notify           chan struct{}
}

type WorkloadController struct {
	mu      sync.Mutex
	classes map[string]*workClassState
}

func newWorkClass(name string, capacity, maxQueue, reservedCritical int) *workClassState {
	if capacity < 1 {
		capacity = 1
	}
	if maxQueue < 0 {
		maxQueue = 0
	}
	if reservedCritical < 0 {
		reservedCritical = 0
	}
	if reservedCritical >= capacity {
		reservedCritical = capacity - 1
	}
	return &workClassState{name: name, capacity: capacity, maxQueue: maxQueue, reservedCritical: reservedCritical, notify: make(chan struct{})}
}

func NewWorkloadController() *WorkloadController {
	return &WorkloadController{classes: map[string]*workClassState{
		"provider-rest": newWorkClass("provider-rest", 6, 24, 2),
		"scanner":       newWorkClass("scanner", 1, 2, 0),
		"background":    newWorkClass("background", 1, 1, 0),
		"ai":            newWorkClass("ai", 1, 2, 0),
	}}
}

func (w *WorkloadController) classLocked(class string) *workClassState {
	st := w.classes[class]
	if st == nil {
		st = newWorkClass(class, 1, 4, 0)
		w.classes[class] = st
	}
	return st
}

func totalQueued(st *workClassState) int {
	n := 0
	for _, v := range st.queued {
		n += v
	}
	return n
}

func reservedCriticalQueue(st *workClassState) int {
	if st == nil || st.maxQueue <= 0 || st.reservedCritical <= 0 {
		return 0
	}
	if st.reservedCritical > st.maxQueue {
		return st.maxQueue
	}
	return st.reservedCritical
}

// queueLimitForTier prevents optional/radar work from consuming the entire
// bounded queue. The queue remains globally capped at maxQueue, but a small
// portion stays available for MarketCritical/UserActionable recovery work.
func queueLimitForTier(st *workClassState, tier WorkTier) int {
	limit := st.maxQueue
	if tier >= WorkTierRadarPromoted {
		limit -= reservedCriticalQueue(st)
	}
	if limit < 0 {
		return 0
	}
	return limit
}

func canAdmit(st *workClassState, tier WorkTier) bool {
	if st.inFlight >= st.capacity {
		return false
	}
	if tier >= WorkTierRadarPromoted {
		if st.inFlight >= st.capacity-st.reservedCritical {
			return false
		}
		if st.queued[WorkTierMarketCritical]+st.queued[WorkTierUserActionable] > 0 {
			return false
		}
	}
	if tier >= WorkTierBroadDiscovery && st.queued[WorkTierRadarPromoted] > 0 {
		return false
	}
	return true
}

func signalWaitersLocked(st *workClassState) {
	close(st.notify)
	st.notify = make(chan struct{})
}

func (w *WorkloadController) Acquire(ctx context.Context, class string) (func(), bool) {
	return w.AcquireTier(ctx, class, WorkTierUserActionable)
}

func (w *WorkloadController) AcquireTier(ctx context.Context, class string, tier WorkTier) (func(), bool) {
	if w == nil {
		return func() {}, true
	}
	tier = normalizeWorkTier(tier)
	registered := false
	queuedAt := time.Now()
	for {
		w.mu.Lock()
		st := w.classLocked(class)
		if canAdmit(st, tier) {
			if registered {
				st.queued[tier]--
				if st.queued[tier] == 0 {
					st.oldestQueue[tier] = time.Time{}
				}
			}
			st.inFlight++
			st.inFlightByTier[tier]++
			w.mu.Unlock()
			var once sync.Once
			return func() {
				once.Do(func() {
					w.mu.Lock()
					st.inFlight--
					st.inFlightByTier[tier]--
					st.completed++
					signalWaitersLocked(st)
					w.mu.Unlock()
				})
			}, true
		}
		if !registered {
			queueLimit := queueLimitForTier(st, tier)
			if totalQueued(st) >= queueLimit {
				st.rejected++
				st.rejectedByTier[tier]++
				if tier >= WorkTierRadarPromoted {
					st.shed++
				}
				w.mu.Unlock()
				return nil, false
			}
			registered = true
			st.queued[tier]++
			if st.oldestQueue[tier].IsZero() {
				st.oldestQueue[tier] = queuedAt
			}
		}
		ch := st.notify
		w.mu.Unlock()
		select {
		case <-ch:
			continue
		case <-ctx.Done():
			w.mu.Lock()
			st = w.classLocked(class)
			if registered && st.queued[tier] > 0 {
				st.queued[tier]--
				if st.queued[tier] == 0 {
					st.oldestQueue[tier] = time.Time{}
				}
			}
			st.canceled++
			signalWaitersLocked(st)
			w.mu.Unlock()
			return nil, false
		}
	}
}

func (w *WorkloadController) TryAcquireTier(class string, tier WorkTier) (func(), bool) {
	if w == nil {
		return func() {}, true
	}
	tier = normalizeWorkTier(tier)
	w.mu.Lock()
	st := w.classLocked(class)
	if !canAdmit(st, tier) {
		st.rejected++
		st.rejectedByTier[tier]++
		if tier >= WorkTierBroadDiscovery {
			st.shed++
		}
		w.mu.Unlock()
		return nil, false
	}
	st.inFlight++
	st.inFlightByTier[tier]++
	w.mu.Unlock()
	var once sync.Once
	return func() {
		once.Do(func() {
			w.mu.Lock()
			st.inFlight--
			st.inFlightByTier[tier]--
			st.completed++
			signalWaitersLocked(st)
			w.mu.Unlock()
		})
	}, true
}

func (w *WorkloadController) Pressured() bool {
	if w == nil {
		return false
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	provider := w.classes["provider-rest"]
	if provider != nil {
		if provider.queued[WorkTierMarketCritical]+provider.queued[WorkTierUserActionable] > 0 {
			return true
		}
		if provider.inFlight >= provider.capacity-provider.reservedCritical {
			return true
		}
	}
	scanner := w.classes["scanner"]
	return scanner != nil && totalQueued(scanner) > 0
}

func (w *WorkloadController) ShouldShed(tier WorkTier) bool {
	if tier <= WorkTierUserActionable {
		return false
	}
	return w.Pressured()
}

func (w *WorkloadController) Diagnostics() []WorkClassDiagnostics {
	if w == nil {
		return nil
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	order := []string{"provider-rest", "scanner", "background", "ai"}
	out := make([]WorkClassDiagnostics, 0, len(w.classes))
	seen := map[string]bool{}
	appendState := func(st *workClassState) {
		if st == nil || seen[st.name] {
			return
		}
		seen[st.name] = true
		d := WorkClassDiagnostics{Class: st.name, Capacity: st.capacity, MaxQueue: st.maxQueue, ReservedCritical: st.reservedCritical, ReservedCriticalQueue: reservedCriticalQueue(st), InFlight: st.inFlight, Queued: totalQueued(st), Completed: st.completed, Canceled: st.canceled, Rejected: st.rejected, Shed: st.shed, QueuedByTier: st.queued, InFlightByTier: st.inFlightByTier, RejectedByTier: st.rejectedByTier}
		oldest := time.Time{}
		for _, at := range st.oldestQueue {
			if !at.IsZero() && (oldest.IsZero() || at.Before(oldest)) {
				oldest = at
			}
		}
		if !oldest.IsZero() {
			d.OldestQueueAgeMs = time.Since(oldest).Milliseconds()
		}
		out = append(out, d)
	}
	for _, name := range order {
		appendState(w.classes[name])
	}
	for _, st := range w.classes {
		appendState(st)
	}
	return out
}
