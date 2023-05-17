package server

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MakeItBright/go-metrics-devops/internal/logger"
	"github.com/MakeItBright/go-metrics-devops/internal/storage"
	"go.uber.org/zap"
)

// Start запуск сервера с переданной конфигурацией.
func Start(cfg Config) error {
	s := storage.NewMemStorage()

	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	srv := newServer(s)

	consumer, err := storage.NewConsumer(cfg.FileStoragePath)
	if err != nil {
		logger.Log.Error("error", zap.Error(err))
	}

	if cfg.Restore {
		metricsFromFile, err := consumer.ReadMetrics()
		if err != nil {
			logger.Log.Error("Не смогли прочитать метрики")
		}

		logger.Log.Info("", zap.Any("metrics", metricsFromFile))
		//TODO
		// s.UpdateAll(metricsFromFile)

	}

	producer, err := storage.NewProducer(cfg.FileStoragePath)
	if err != nil {
		logger.Log.Error("Не смогли инициализировать продюсера")
	}

	storeIntervalTicker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
	defer storeIntervalTicker.Stop()

	go func() {
		for {
			select {
			case <-storeIntervalTicker.C:
				logger.Log.Info("Read and Write metrics")
				metrics := s.GetAllMetrics()

				producer.WriteMetrics(metrics)
			case <-osSigChan:
				logger.Log.Info("Read and Write metrics End")
				metrics := s.GetAllMetrics()

				producer.WriteMetrics(metrics)
				os.Exit(0)
			}
		}
	}()

	logger.Log.Info("Running server", zap.String("address", cfg.BindAddr))
	return http.ListenAndServe(cfg.BindAddr, srv)

}
