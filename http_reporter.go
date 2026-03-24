package service_infosth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type httpReporter struct {
	config Config
	client *http.Client
}

func NewHTTPReporter(config Config) Reporter {
	client := http.Client{Timeout: 10 * time.Second}

	return &httpReporter{
		config: config,
		client: &client,
	}
}

func (h *httpReporter) ReportError(ctx context.Context, errorEvent ErrorEvent) error {
	body, err := json.Marshal(errorEvent)
	if err != nil {
		return fmt.Errorf("error in marshal: %w", err)
	}

	reportErrorUrl := h.config.ErrorPath

	if reportErrorUrl == "" {
		reportErrorUrl = string(reportError)
	}

	path := buildURLPath(h.config.BaseUrl, reportErrorUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error in newRequest (path): %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	for _, value := range h.config.Headers {
		req.Header.Set(value.HeaderName, value.HeaderValue)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("error in do (req): %w", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	if statusCode >= 400 {
		return fmt.Errorf("error in do (req): status code equals/higher than 400 (%d)", statusCode)
	}

	return nil
}

func (h *httpReporter) ReportMetrics(ctx context.Context, snapshot Snapshot) error {
	body, err := marshalMetricsSnapshot(snapshot)
	if err != nil {
		return fmt.Errorf("error in marshal: %w", err)
	}

	reportMetricsUrl := h.config.MetricsPath

	if reportMetricsUrl == "" {
		reportMetricsUrl = string(reportMetrics)
	}

	path := buildURLPath(h.config.BaseUrl, reportMetricsUrl)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error in newRequest (path): %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	for _, value := range h.config.Headers {
		req.Header.Set(value.HeaderName, value.HeaderValue)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("error in do (req): %w", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	if statusCode >= 400 {
		return fmt.Errorf("error in do (req): status code equals/higher than 400 (%d)", statusCode)
	}

	return nil
}

func marshalMetricsSnapshot(snapshot Snapshot) ([]byte, error) {
	items := make([]map[string]any, 0, len(snapshot.Items))

	for _, item := range snapshot.Items {
		items = append(items, map[string]any{
			"route":           item.Route,
			"method":          item.Method,
			"request_count":   item.RequestCount,
			"success_count":   item.SuccessCount,
			"error_count":     item.ErrorCount,
			"duration_sum_ms": item.DurationSumMs,
			"max_latency_ms":  item.MaxLatencyMs,
			"status_2xx":      item.Status2xx,
			"status_4xx":      item.Status4xx,
			"status_5xx":      item.Status5xx,
		})
	}

	payload := map[string]any{
		"service":      snapshot.Service,
		"window_start": snapshot.WindowStart,
		"window_end":   snapshot.WindowEnd,
		"items":        items,
	}

	return json.Marshal(payload)
}

func buildURLPath(baseURL, endpointPath string) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	endpointPath = strings.TrimPrefix(endpointPath, "/")

	return baseURL + "/" + endpointPath
}
