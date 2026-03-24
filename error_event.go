package service_infosth

import "time"

type ErrorEvent struct {
	Service    string    `json:"service"`
	Route      string    `json:"route"`
	Method     string    `json:"method"`
	StatusCode int       `json:"status_code"`
	DurationMs uint64    `json:"duration_ms"`
	Timestamp  time.Time `json:"timestamp"`
}
