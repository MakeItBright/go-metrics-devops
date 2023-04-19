package agent

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Config структура конфигурации агента.
type Config struct {
	Scheme         string
	Address        string        // адрес HTTP-сервера
	PollInterval   time.Duration // интервал обновления метрик
	ReportInterval time.Duration // интервал между отправками метрик
	Logger         *logrus.Logger
}
