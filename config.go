package service_infosth

import "time"

type Config struct {
	ServiceName     string
	AuthHeaderName  string
	AuthHeaderValue string
	BaseUrl         string
	ErrorPath       string
	MetricsPath     string
	FlushInterval   time.Duration
}
