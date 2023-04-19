package server

import (
	"fmt"
	"net/http"

	"github.com/MakeItBright/go-metrics-devops/internal/storage"
)

// Start запуск сервера с переданной конфигурацией.
func Start(cfg Config) error {
	s := storage.NewMemStorage()

	srv := newServer(s)

	fmt.Println("Running server on", cfg.BindAddr)

	return http.ListenAndServe(cfg.BindAddr, srv)

}
