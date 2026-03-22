package service_infosth

import "time"

type Config struct {
	ServiceName   string
	BaseUrl       string
	FlushInterval time.Duration
}
