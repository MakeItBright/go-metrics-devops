package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
		body string
	}
	tests := []struct {
		name    string
		method  string
		args    string
		useJSON bool
		body    string
		want    want
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
		{
			name:   "Get gauge",
			method: "GET",
			args:   "/value/gauge/Alloc",
			want: want{
				code: 200,
				body: "2.128506e+06",
			},
		},
		{
			name:   "counter no name 1",
			method: "POST",
			args:   "/update/counter/",
			want:   want{code: 404},
		},
		{
			name:   "bad type",
			method: "POST",
			args:   "/update/integer/x/1",
			want:   want{code: 501},
		},
		{
			name:    "update json counter",
			method:  "POST",
			args:    "/update/",
			useJSON: true,
			body:    `{"id":"xyz","type":"counter","delta":10}`,
			want: want{
				code: 200,
				body: `{"id":"xyz","type":"counter","delta":10}`,
			},
		},
		{
			name:    "update json counter",
			method:  "POST",
			args:    "/update/",
			useJSON: true,
			body:    `{"id":"xyz","type":"counter","delta":10}`,
			want: want{
				code: 200,
				body: `{"id":"xyz","type":"counter","delta":20}`,
			},
		},
		{
			name:    "get json counter",
			method:  "POST",
			args:    "/value/",
			useJSON: true,
			body:    `{"id":"xyz","type":"counter"}`,
			want: want{
				code: 200,
				body: `{"id":"xyz","type":"counter","delta":20}`,
			},
		},

		{
			name:    "update json gauge",
			method:  "POST",
			args:    "/update/",
			useJSON: true,
			body:    `{"id":"xyz","type":"gauge","value":10}`,
			want:    want{code: 200},
		},
		{
			name:    "get json gauge",
			method:  "POST",
			args:    "/value/",
			useJSON: true,
			body:    `{"id":"xyz","type":"gauge"}`,
			want: want{
				code: 200,
				body: `{"id":"xyz","type":"gauge","value":10}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rec := httptest.NewRecorder()

			r := strings.NewReader("")
			if tt.useJSON {
				r = strings.NewReader(tt.body)
			}

			req, _ := http.NewRequest(tt.method, tt.args, r)

			if tt.useJSON {
				req.Header.Add("Content-Type", "application/json")
			}
			s.ServeHTTP(rec, req)
			assert.Equal(t, tt.want.code, rec.Code)

			respBody, _ := io.ReadAll(rec.Body)
			fmt.Printf("out %+v", string(respBody))
			assert.Contains(t, string(respBody), tt.want.body)

		})
	}
}
