package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MakeItBright/go-metrics-devops/internal/storage/inmem"
	"github.com/stretchr/testify/assert"
)

func TestServer_HandleHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	s := newServer(inmem.New())
	s.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)
}
