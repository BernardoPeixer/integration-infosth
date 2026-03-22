package service_infosth

import "time"

type Config struct {
	ServiceName   string
	Headers       []ConfigHeader
	BaseUrl       string
	ErrorPath     string
	MetricsPath   string
	FlushInterval time.Duration
}

type ConfigHeader struct {
	HeaderName  string
	HeaderValue string
}
