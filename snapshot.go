package service_infosth

import "time"

type Snapshot struct {
	Service     string
	WindowStart time.Time
	WindowEnd   time.Time
	Items       []SnapshotAggregated
}

type SnapshotAggregated struct {
	Route         string
	Method        string
	RequestCount  uint64
	SuccessCount  uint64
	ErrorCount    uint64
	DurationSumMs uint64
	MaxLatencyMs  uint64
	Status2xx     uint64
	Status4xx     uint64
	Status5xx     uint64
}
