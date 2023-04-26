package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MakeItBright/go-metrics-devops/internal/model"
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

func TestHandleJsonPostUpdateMetric(t *testing.T) {
	sm := storage.NewMemStorage()
	srv := newServer(sm)

	m := model.Metric{
		Name:  "metric_name",
		Type:  model.MetricTypeCounter,
		Delta: 10,
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r, err := http.NewRequest("POST", "/update", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	w := httptest.NewRecorder()
	srv.handleJSONPostUpdateMetric(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("unexpected status code: want %d, got %d", http.StatusOK, w.Code)
	}

	expectedBody := "Metric updated"
	if body := w.Body.String(); body != expectedBody {
		t.Errorf("unexpected response body: want %q, got %q", expectedBody, body)
	}

	// check if the metric was added
	if _, ok := sm.GetCounter("metric_name"); !ok {
		t.Errorf("metric was not added to the storage")
	}
}
func TestServer(t *testing.T) {
	// create a new instance of the in-memory storage
	sm := storage.NewMemStorage()

	// create a new instance of the server
	s := newServer(sm)

	// define test cases
	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus int
		expectedBody   string
	}{

		{
			name:           "post gauge metric",
			method:         "POST",
			path:           "/update",
			body:           `{"id":"test_gauge_metric","type":"gauge","value":2.5}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"test_gauge_metric","type":"gauge","value":2.5}`,
		},
		{
			name:           "post counter metric",
			method:         "POST",
			path:           "/update",
			body:           `{"name":"test_counter_metric","type":"counter","delta":1}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `Metric updated`,
		},
		{
			name:           "get gauge metric",
			method:         "POST",
			path:           "/value",
			body:           `{"name":"test_gauge_metric","type":"gauge"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"name":"test_gauge_metric","type":"gauge","value":2.5}`,
		},
		{
			name:           "get counter metric",
			method:         "POST",
			path:           "/value",
			body:           `{"name":"test_counter_metric","type":"counter"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"name":"test_counter_metric","type":"counter","delta":1}`,
		},
		{
			name:           "get unknown metric",
			method:         "POST",
			path:           "/value",
			body:           `{"name":"unknown_metric","type":"gauge"}`,
			expectedStatus: http.StatusNotFound,
			expectedBody:   ``,
		},
	}

	// run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create request
			reqBody := strings.NewReader(tt.body)
			req := httptest.NewRequest(tt.method, tt.path, reqBody)

			// create response recorder
			recorder := httptest.NewRecorder()

			// call server handler
			s.ServeHTTP(recorder, req)

			// check response
			// assert.Equal(t, tt.expectedStatus, recorder.Code)

			if recorder.Body.Len() > 0 {
				var response map[string]interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal response body: %v", err)
				}
				fmt.Println(response)
				assert.Equal(t, tt.expectedBody, response)
			} else {
				assert.Equal(t, tt.expectedBody, "")
			}
		})
	}
}

func Test_server_handlePostUpdateMetric(t *testing.T) {
	s := newServer(storage.NewMemStorage())

	type want struct {
		code int
		body []string
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
				body: []string{"100500"},
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
				body: []string{"100501"},
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

			for _, s := range tt.want.body {
				assert.Contains(t, string(respBody), s)
			}
		})
	}
}
