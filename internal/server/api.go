package server

import (
	"net/http"

	"github.com/MakeItBright/go-metrics-devops/internal/store/teststore"
)

// Start
func Start(config *Config) error {

	store := teststore.New()

	srv := newServer(store)

	return http.ListenAndServe(config.BindAddr, srv)
}
