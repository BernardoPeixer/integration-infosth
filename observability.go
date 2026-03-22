package service_infosth

import (
	"context"
	"log/slog"
	"time"
)

type Observability struct {
	config     Config
	aggregator *Aggregator
	reporter   Reporter
}

func New(config Config) *Observability {
	aggregator := NewAggregator(config)
	reporter := NewHTTPReporter(config)

	return &Observability{
		config:     config,
		aggregator: aggregator,
		reporter:   reporter,
	}
}

func (o *Observability) RunFlusher(ctx context.Context) {
	ticker := time.NewTicker(o.config.FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			snapshot := o.aggregator.SnapshotAndReset()

			if len(snapshot.Items) == 0 {
				continue
			}

			err := o.reporter.ReportMetrics(ctx, snapshot)
			if err != nil {
				slog.Error("error in reportMetrics", "error", err)
			}
		case <-ctx.Done():
			snapshot := o.aggregator.SnapshotAndReset()

			if len(snapshot.Items) == 0 {
				return
			}

			// TODO: Implements sent metrics here too
			slog.Info("service shutdown, sending last metrics")

			tempCtx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)

			err := o.reporter.ReportMetrics(tempCtx, snapshot)
			if err != nil {
				slog.Error("error in reportMetrics", "error", err)
			}

			cancelCtx()

			return
		}
	}
}

func (o *Observability) Start(ctx context.Context) {
	go o.RunFlusher(ctx)
}
