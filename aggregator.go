package service_infosth

import (
	"sync"
	"time"
)

type MetricKey struct {
	Service string
	Route   string
	Method  string
}

type MetricAgg struct {
	RequestCount uint64
	SuccessCount uint64
	ErrorCount   uint64

	DurationSumMs uint64
	MaxLatencyMs  uint64

	Status2xx uint64
	Status4xx uint64
	Status5xx uint64
}

type Aggregator struct {
	mu          sync.Mutex
	service     string
	windowStart time.Time
	data        map[MetricKey]*MetricAgg
}

func NewAggregator(config Config) *Aggregator {
	return &Aggregator{
		service:     config.ServiceName,
		data:        make(map[MetricKey]*MetricAgg),
		windowStart: time.Now(),
	}
}

func (a *Aggregator) SnapshotAndReset() Snapshot {
	a.mu.Lock()
	defer a.mu.Unlock()

	snapshotItems := make([]SnapshotAggregated, 0, len(a.data))

	for key, data := range a.data {
		snapshotAggregated := SnapshotAggregated{
			Route:         key.Route,
			Method:        key.Method,
			RequestCount:  data.RequestCount,
			SuccessCount:  data.SuccessCount,
			ErrorCount:    data.ErrorCount,
			DurationSumMs: data.DurationSumMs,
			MaxLatencyMs:  data.MaxLatencyMs,
			Status2xx:     data.Status2xx,
			Status4xx:     data.Status4xx,
			Status5xx:     data.Status5xx,
		}

		snapshotItems = append(snapshotItems, snapshotAggregated)
	}

	snapshot := Snapshot{
		Service:     a.service,
		WindowStart: a.windowStart,
		WindowEnd:   time.Now(),
		Items:       snapshotItems,
	}

	a.data = make(map[MetricKey]*MetricAgg)
	a.windowStart = time.Now()

	return snapshot
}

func (a *Aggregator) Observe(
	route string,
	method string,
	statusCode int,
	latency time.Duration,
) {
	a.mu.Lock()
	defer a.mu.Unlock()

	key := MetricKey{
		Service: a.service,
		Route:   route,
		Method:  method,
	}

	data, ok := a.data[key]
	if !ok {
		a.data[key] = &MetricAgg{}
		data = a.data[key]
	}

	data.RequestCount++

	if statusCode >= 500 {
		data.Status5xx++
		data.ErrorCount++
	} else if statusCode >= 400 {
		data.Status4xx++
		data.ErrorCount++
	} else {
		data.Status2xx++
		data.SuccessCount++
	}

	latencyMS := uint64(latency / time.Millisecond)

	if data.MaxLatencyMs < latencyMS {
		data.MaxLatencyMs = latencyMS
	}

	data.DurationSumMs += latencyMS
}
