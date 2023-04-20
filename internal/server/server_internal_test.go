package server

import (
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
