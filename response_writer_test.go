package service_infosth

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseWriter_WriteHeader(t *testing.T) {
	rw := httptest.NewRecorder()
	w := &responseWriter{
		ResponseWriter: rw,
		statusCode:     0,
	}

	w.WriteHeader(200)

	assert.Equal(t, 200, w.statusCode)
	assert.Equal(t, 200, rw.Code)
}

func TestResponseWriter_Write(t *testing.T) {
	tests := []struct {
		name              string
		initialStatusCode int
		bodyWrites        []string
		wantStatusCode    int
		wantBody          string
	}{
		{
			name:              "without write header",
			initialStatusCode: 0,
			bodyWrites:        []string{"test"},
			wantStatusCode:    200,
			wantBody:          "test",
		},
		{
			name:              "with write header 201",
			initialStatusCode: 201,
			bodyWrites:        []string{"test"},
			wantStatusCode:    201,
			wantBody:          "test",
		},
		{
			name:              "with write header 404",
			initialStatusCode: 404,
			bodyWrites:        []string{"test"},
			wantStatusCode:    404,
			wantBody:          "test",
		},
		{
			name:              "with empty body",
			initialStatusCode: 0,
			bodyWrites:        []string{""},
			wantStatusCode:    200,
			wantBody:          "",
		},
		{
			name:              "multiple writes without write header",
			initialStatusCode: 0,
			bodyWrites:        []string{"te", "st"},
			wantStatusCode:    200,
			wantBody:          "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()

			w := &responseWriter{
				ResponseWriter: rw,
				statusCode:     tt.initialStatusCode,
			}

			if tt.initialStatusCode != 0 {
				w.WriteHeader(tt.initialStatusCode)
			}

			for _, part := range tt.bodyWrites {
				w.Write([]byte(part))
			}
			
			assert.Equal(t, tt.wantStatusCode, w.statusCode)
			assert.Equal(t, tt.wantBody, rw.Body.String())
		})
	}

}
