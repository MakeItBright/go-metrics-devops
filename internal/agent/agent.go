package agent

import (
	"time"

	"github.com/MakeItBright/go-metrics-devops/internal/monitor"
	"github.com/MakeItBright/go-metrics-devops/internal/sender"
	"github.com/MakeItBright/go-metrics-devops/internal/storage"
	"github.com/sirupsen/logrus"
)

// Run запуск агента с переданной конфигурацией.
func Run(cfg Config) error {
	if cfg.Logger == nil {
		cfg.Logger = logrus.StandardLogger()
	}

	// устанавливаем интервал для периодической отправки HTTP-запросов
	pollInterval := cfg.PollInterval
	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()

	// устанавливаем интервал для отправки метрик
	reportInterval := cfg.ReportInterval
	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	m := monitor.NewMonitor(storage.NewMemStorage(), *sender.NewSender("http://" + cfg.Address))

	// запускаем бесконечный цикл для периодической отправки HTTP-запросов
	for {

		select {
		case <-pollTicker.C:
			// собираем метрики
			cfg.Logger.Infof(
				"agent is running, collect metrics every %v seconds",
				pollInterval.Seconds(),
			)
			m.CollectMetrics()

		case <-reportTicker.C:
			// отправляем HTTP-запросы на указанные адреса
			cfg.Logger.Infof(
				"agent is running, sending requests to %v every %v seconds",
				cfg.Address,
				reportInterval.Seconds(),
			)
			m.Dump()

		}
	}
}
