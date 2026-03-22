package service_infosth

import "time"

type ErrorEvent struct {
	Service    string
	Route      string
	Method     string
	StatusCode int
	DurationMs uint64
	Timestamp  time.Time
}
