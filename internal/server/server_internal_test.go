package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MakeItBright/go-metrics-devops/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestServer_HandleHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	s := newServer(storage.NewMemStorage())
	s.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)
}
func TestServer_handleJsonPostUpdateMetric(t *testing.T) {
	s := newServer(storage.NewMemStorage())

	testCases := []struct {
		name      string
		request   interface{}
		expectErr bool
	}{
		{
			name: "add_gauge_metric",
			request: map[string]interface{}{
				"type":  "gauge",
				"name":  "metric_1",
				"value": 42.0,
			},
			expectErr: false,
		},
		{
			name: "add_counter_metric",
			request: map[string]interface{}{
				"type":  "counter",
				"name":  "metric_2",
				"delta": 1,
			},
			expectErr: false,
		},
		{
			name: "invalid_request_body",
			request: map[string]interface{}{
				"type": "invalid",
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tc.request)
			req, err := http.NewRequest("POST", "/update", bytes.NewReader(reqBody))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(s.handleJSONPostUpdateMetric)
			handler.ServeHTTP(rr, req)

			if tc.expectErr {
				assert.NotEqual(t, rr.Code, http.StatusOK)
			} else {
				assert.Equal(t, rr.Code, http.StatusOK)
			}
		})
	}
}

func Test_server_handlePostUpdateMetric(t *testing.T) {
	s := newServer(storage.NewMemStorage())

	type want struct {
		code int
		body string
	}
	tests := []struct {
		name   string
		method string
		args   string
		want   want
	}{
		// TODO: Add test cases.
		{
			name:   "gauge OK",
			method: "POST",
			args:   "/update/gauge/Alloc/2128506",
			want:   want{code: 200},
		},
		{
			name:   "counter OK",
			method: "POST",
			args:   "/update/counter/RandomValue/100500",
			want:   want{code: 200},
		},
		{
			name:   "Get counter",
			method: "GET",
			args:   "/value/counter/RandomValue",
			want: want{
				code: 200,
				body: "100500",
			},
		},
		{
			name:   "counter OK",
			method: "POST",
			args:   "/update/counter/RandomValue/1",
			want:   want{code: 200},
		},
		{
			name:   "Get counter",
			method: "GET",
			args:   "/value/counter/RandomValue",
			want: want{
				code: 200,
				body: "100501",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.args, nil)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tt.want.code, rec.Code)

			respBody, _ := io.ReadAll(rec.Body)
			assert.Contains(t, string(respBody), tt.want.body)
			// for _, s := range tt.want.body {

			// }
		})
	}
}
