package service_infosth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testCall struct {
	route      string
	method     string
	statusCode int
	latency    time.Duration
}

type testCase struct {
	name              string
	calls             []testCall
	wantRequests      uint64
	wantSuccesses     uint64
	wantErrors        uint64
	wantDurationSumMs uint64
	wantMaxLatencyMs  uint64
	wantStatus2xx     uint64
	wantStatus4xx     uint64
	wantStatus5xx     uint64
	wantGroups        uint64
}

func TestAggregator_Observe(t *testing.T) {
	tests := []testCase{
		{
			name: "single call",
			calls: []testCall{
				{route: "/test", method: "GET", statusCode: 200, latency: 100 * time.Millisecond},
			},
			wantRequests:      1,
			wantSuccesses:     1,
			wantErrors:        0,
			wantDurationSumMs: 100,
			wantMaxLatencyMs:  100,
			wantStatus2xx:     1,
			wantStatus4xx:     0,
			wantStatus5xx:     0,
			wantGroups:        1,
		},
		{
			name: "multiple calls",
			calls: []testCall{
				{route: "/test", method: "GET", statusCode: 200, latency: 100 * time.Millisecond},
				{route: "/test", method: "GET", statusCode: 401, latency: 100 * time.Millisecond},
				{route: "/test", method: "GET", statusCode: 401, latency: 100 * time.Millisecond},
				{route: "/test", method: "GET", statusCode: 200, latency: 100 * time.Millisecond},
			},
			wantRequests:      4,
			wantSuccesses:     2,
			wantErrors:        2,
			wantDurationSumMs: 400,
			wantMaxLatencyMs:  100,
			wantStatus2xx:     2,
			wantStatus4xx:     2,
			wantStatus5xx:     0,
			wantGroups:        1,
		},
		{
			name: "multiple calls with different routes",
			calls: []testCall{
				{route: "/test1", method: "GET", statusCode: 200, latency: 100 * time.Millisecond},
				{route: "/test2", method: "GET", statusCode: 401, latency: 100 * time.Millisecond},
			},
			wantRequests:      2,
			wantSuccesses:     1,
			wantErrors:        1,
			wantDurationSumMs: 200,
			wantMaxLatencyMs:  100,
			wantStatus2xx:     1,
			wantStatus4xx:     1,
			wantStatus5xx:     0,
			wantGroups:        2,
		},
		{
			name: "multiple calls with different methods",
			calls: []testCall{
				{route: "/test", method: "GET", statusCode: 200, latency: 100 * time.Millisecond},
				{route: "/test", method: "POST", statusCode: 401, latency: 100 * time.Millisecond},
			},
			wantRequests:      2,
			wantSuccesses:     1,
			wantErrors:        1,
			wantDurationSumMs: 200,
			wantMaxLatencyMs:  100,
			wantStatus2xx:     1,
			wantStatus4xx:     1,
			wantStatus5xx:     0,
			wantGroups:        2,
		},
		{
			name:              "empty calls",
			calls:             []testCall{},
			wantRequests:      0,
			wantSuccesses:     0,
			wantErrors:        0,
			wantDurationSumMs: 0,
			wantMaxLatencyMs:  0,
			wantStatus2xx:     0,
			wantStatus4xx:     0,
			wantStatus5xx:     0,
			wantGroups:        0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAggregator(Config{
				ServiceName: "test",
			})

			for _, call := range tt.calls {
				a.Observe(call.route, call.method, call.statusCode, call.latency)
			}

			snapshot := a.SnapshotAndReset()

			totalRequests := uint64(0)
			totalSuccesses := uint64(0)
			totalErrors := uint64(0)
			totalDurationSumMs := uint64(0)
			totalStatus2xx := uint64(0)
			totalStatus4xx := uint64(0)
			totalStatus5xx := uint64(0)
			maxLatencyMs := uint64(0)

			for _, item := range snapshot.Items {
				if item.MaxLatencyMs > maxLatencyMs {
					maxLatencyMs = item.MaxLatencyMs
				}

				totalRequests += item.RequestCount
				totalSuccesses += item.SuccessCount
				totalErrors += item.ErrorCount
				totalDurationSumMs += item.DurationSumMs
				totalStatus2xx += item.Status2xx
				totalStatus4xx += item.Status4xx
				totalStatus5xx += item.Status5xx
			}

			assert.Equal(t, tt.wantRequests, totalRequests, "total requests")
			assert.Equal(t, tt.wantSuccesses, totalSuccesses, "total successes")
			assert.Equal(t, tt.wantErrors, totalErrors, "total errors")
			assert.Equal(t, tt.wantDurationSumMs, totalDurationSumMs, "total duration sum ms")
			assert.Equal(t, tt.wantStatus2xx, totalStatus2xx, "total status 2xx")
			assert.Equal(t, tt.wantStatus4xx, totalStatus4xx, "total status 4xx")
			assert.Equal(t, tt.wantStatus5xx, totalStatus5xx, "total status 5xx")
			assert.Equal(t, tt.wantMaxLatencyMs, maxLatencyMs, "max latency ms")
			assert.Equal(t, tt.wantGroups, uint64(len(snapshot.Items)), "total groups")
		})
	}
}

func TestAggregator_SnapshotAndReset(t *testing.T) {
	tests := []testCase{
		{
			name:              "empty snapshot",
			calls:             []testCall{},
			wantRequests:      0,
			wantSuccesses:     0,
			wantErrors:        0,
			wantDurationSumMs: 0,
			wantMaxLatencyMs:  0,
			wantStatus2xx:     0,
			wantStatus4xx:     0,
			wantStatus5xx:     0,
			wantGroups:        0,
		},
		{
			name: "single call",
			calls: []testCall{
				{route: "/test", method: "GET", statusCode: 200, latency: 100 * time.Millisecond},
			},
			wantRequests:      1,
			wantSuccesses:     1,
			wantErrors:        0,
			wantDurationSumMs: 100,
			wantMaxLatencyMs:  100,
			wantStatus2xx:     1,
			wantStatus4xx:     0,
			wantStatus5xx:     0,
			wantGroups:        1,
		},
		{
			name: "multiple calls with same route",
			calls: []testCall{
				{route: "/test", method: "GET", statusCode: 200, latency: 100 * time.Millisecond},
				{route: "/test", method: "GET", statusCode: 401, latency: 100 * time.Millisecond},
			},
			wantRequests:      2,
			wantSuccesses:     1,
			wantErrors:        1,
			wantDurationSumMs: 200,
			wantMaxLatencyMs:  100,
			wantStatus2xx:     1,
			wantStatus4xx:     1,
			wantStatus5xx:     0,
			wantGroups:        1,
		},
		{
			name: "multiple calls with different routes",
			calls: []testCall{
				{route: "/test1", method: "GET", statusCode: 200, latency: 100 * time.Millisecond},
				{route: "/test2", method: "GET", statusCode: 401, latency: 100 * time.Millisecond},
			},
			wantRequests:      2,
			wantSuccesses:     1,
			wantErrors:        1,
			wantDurationSumMs: 200,
			wantMaxLatencyMs:  100,
			wantStatus2xx:     1,
			wantStatus4xx:     1,
			wantStatus5xx:     0,
			wantGroups:        2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAggregator(Config{
				ServiceName: "test",
			})

			for _, call := range tt.calls {
				a.Observe(call.route, call.method, call.statusCode, call.latency)
			}

			snapshot := a.SnapshotAndReset()

			totalRequests := uint64(0)
			totalSuccesses := uint64(0)
			totalErrors := uint64(0)
			totalDurationSumMs := uint64(0)
			totalStatus2xx := uint64(0)
			totalStatus4xx := uint64(0)
			totalStatus5xx := uint64(0)
			maxLatencyMs := uint64(0)

			for _, item := range snapshot.Items {
				if item.MaxLatencyMs > maxLatencyMs {
					maxLatencyMs = item.MaxLatencyMs
				}

				totalRequests += item.RequestCount
				totalSuccesses += item.SuccessCount
				totalErrors += item.ErrorCount
				totalDurationSumMs += item.DurationSumMs
				totalStatus2xx += item.Status2xx
				totalStatus4xx += item.Status4xx
				totalStatus5xx += item.Status5xx
			}

			assert.Equal(t, tt.wantRequests, totalRequests, "total requests")
			assert.Equal(t, tt.wantSuccesses, totalSuccesses, "total successes")
			assert.Equal(t, tt.wantErrors, totalErrors, "total errors")
			assert.Equal(t, tt.wantDurationSumMs, totalDurationSumMs, "total duration sum ms")
			assert.Equal(t, tt.wantStatus2xx, totalStatus2xx, "total status 2xx")
			assert.Equal(t, tt.wantStatus4xx, totalStatus4xx, "total status 4xx")
			assert.Equal(t, tt.wantStatus5xx, totalStatus5xx, "total status 5xx")
			assert.Equal(t, tt.wantMaxLatencyMs, maxLatencyMs, "max latency ms")
			assert.Equal(t, tt.wantGroups, uint64(len(snapshot.Items)), "total groups")

			assert.Len(t, a.data, 0, "data length")
			assert.Equal(t, "test", snapshot.Service, "service")
		})
	}
}
