package server

import (
	"net/http"

	"github.com/MakeItBright/go-metrics-devops/internal/storage"
)

// Start запуск сервера с переданной конфигурацией.
func Start(cfg Config) error {
	s := storage.NewMemStorage()

	srv := newServer(s)

	// logger.Log.Info("Running server", cfg.BindAddr)
	return http.ListenAndServe(cfg.BindAddr, srv)

}
