package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MakeItBright/go-metrics-devops/internal/logger"
	"github.com/MakeItBright/go-metrics-devops/internal/model"
	"github.com/MakeItBright/go-metrics-devops/internal/storage"
	"go.uber.org/zap"
)

// Start запуск сервера с переданной конфигурацией.
func Start(cfg Config) error {
	s := storage.NewMemStorage()

	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	srv := newServer(s)

	reader, err := storage.NewReader(cfg.FileStoragePath)
	if err != nil {
		return fmt.Errorf("cannot create reader: %w", err)
	}

	if cfg.Restore {
		metricsFromFile, err := reader.ReadMetrics()
		if err != nil {
			return fmt.Errorf("cannot read metrics: %w", err)
		}

		for _, metricValue := range metricsFromFile {
			switch model.MetricType(metricValue.Type) {
			case model.MetricTypeGauge:
				s.AddGauge(string(metricValue.Name), metricValue.Value)
			case model.MetricTypeCounter:
				s.AddCounter(string(metricValue.Name), metricValue.Delta)
			default:
				return fmt.Errorf("cannot read metrics: %w", err)
			}

		}

	}

	writer, err := storage.NewWriter(cfg.FileStoragePath)
	if err != nil {
		return fmt.Errorf("cannot init writer: %w", err)
	}

	storeIntervalTicker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
	defer storeIntervalTicker.Stop()

	go func() {
		for {
			select {
			case <-storeIntervalTicker.C:
				metrics := s.GetAllMetrics()
				writer.WriteMetrics(metrics)
			case <-osSigChan:

				os.Exit(0)
			}
		}
	}()

	logger.Log.Info("Running server", zap.String("address", cfg.BindAddr))
	return http.ListenAndServe(cfg.BindAddr, srv)

}
