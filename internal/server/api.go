package server

import (
	"net/http"

	"github.com/MakeItBright/go-metrics-devops/internal/storage/inmem"
)

// Start
func Start(config *Config) error {

	store := inmem.New()

	srv := newServer(store)

	return http.ListenAndServe(config.BindAddr, srv)
}
