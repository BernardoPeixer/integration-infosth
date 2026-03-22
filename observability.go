package service_infosth

import "context"

type Observability struct {
	config     Config
	aggregator *Aggregator
	reporter   Reporter
}

func New(config Config) *Observability {
	aggregator := NewAggregator()

	reporter := NewHTTPReporter(config.BaseUrl)

	return &Observability{
		config:     config,
		aggregator: aggregator,
		reporter:   reporter,
	}
}

// Start esse cara aqui vai iniciar o flush
func (o *Observability) Start(ctx context.Context) {

}

func (o *Observability) Shutdown() {

}
