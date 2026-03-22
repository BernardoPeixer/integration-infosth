package service_infosth

import "context"

type Reporter interface {
	ReportError(ctx context.Context, errorEvent ErrorEvent) error
	ReportMetrics(ctx context.Context, snapshot Snapshot) error
}
