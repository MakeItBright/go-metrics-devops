package server

import (
	"net/http"

	"github.com/MakeItBright/go-metrics-devops/internal/logger"
	"github.com/MakeItBright/go-metrics-devops/internal/storage"
	"go.uber.org/zap"
)

// Start запуск сервера с переданной конфигурацией.
func Start(cfg Config) error {
	s := storage.NewMemStorage()

	srv := newServer(s)

	logger.Log.Info("Running server", zap.String("address", cfg.BindAddr))
	return http.ListenAndServe(cfg.BindAddr, srv)

}
