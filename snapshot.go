package service_infosth

import "time"

type Snapshot struct {
	Service     string               `json:"service"`
	WindowStart time.Time            `json:"window_start"`
	WindowEnd   time.Time            `json:"window_end"`
	Items       []SnapshotAggregated `json:"items"`
}

type SnapshotAggregated struct {
	Route         string `json:"route"`
	Method        string `json:"method"`
	RequestCount  uint64 `json:"request_count"`
	SuccessCount  uint64 `json:"success_count"`
	ErrorCount    uint64 `json:"error_count"`
	DurationSumMs uint64 `json:"duration_sum_ms"`
	MaxLatencyMs  uint64 `json:"max_latency_ms"`
	Status2xx     uint64 `json:"status_2xx"`
	Status4xx     uint64 `json:"status_4xx"`
	Status5xx     uint64 `json:"status_5xx"`
}
