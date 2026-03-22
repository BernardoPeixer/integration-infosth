package service_infosth

import (
	"log/slog"
	"net/http"
	"time"
)

func (o *Observability) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     0,
		}

		next.ServeHTTP(rw, r)

		if rw.statusCode == 0 {
			rw.statusCode = 200
		}

		latency := time.Since(start)
		path := r.URL.Path
		method := r.Method
		statusCode := rw.statusCode

		o.aggregator.Observe(path, method, statusCode, latency)

		if statusCode >= 500 {
			latencyMs := uint64(latency / time.Millisecond)

			errorEvent := ErrorEvent{
				Service:    o.config.ServiceName,
				Route:      path,
				Method:     method,
				StatusCode: statusCode,
				DurationMs: latencyMs,
				Timestamp:  time.Now(),
			}

			// TODO: Create channel to listen this, and when channel receive a error, makes the report
			err := o.reporter.ReportError(r.Context(), errorEvent)
			if err != nil {
				slog.Error("error in reportError", "error", err)
			}
		}
	})
}
