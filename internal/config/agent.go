package config

import (
	"flag"
	"os"
	"time"
)

// неэкспортированные переменные flagAddress, flagPollInterval, flagReportInterval,
// содержат значения флагов командной строки
var (
	flagAddress        string
	flagPollInterval   time.Duration
	flagReportInterval time.Duration
)

// константы, задающие значения по умолчанию
const (
	defaultAddress        = "localhost:8080" // адрес сервера по умолчанию
	defaultScheme         = "http"
	defaultPollInterval   = 2 * time.Second  // интервал обновления метрик по умолчанию
	defaultReportInterval = 10 * time.Second // интервал отправки метрик на сервер по умолчанию
)

// AgentConfig - структура конфигурации агента
type AgentConfig struct {
	Address        string        `env:"ADDRESS"`         // адрес HTTP-сервера
	PollInterval   time.Duration `env:"POLL_INTERVAL"`   // интервал обновления метрик
	ReportInterval time.Duration `env:"REPORT_INTERVAL"` // интервал между отправками метрик
}

// NewAgentConfig - создает новую структуру конфигурации агента
func NewAgentConfig() *AgentConfig {
	// устанавливаем значения флагов командной строки
	flag.StringVar(&flagAddress, "a", defaultAddress, "ADDRESS=<ЗНАЧЕНИЕ> адрес эндпоинта HTTP-сервера")
	flag.DurationVar(&flagPollInterval, "p", defaultPollInterval, "POLL_INTERVAL=<ЗНАЧЕНИЕ> частота опроса метрик ")
	flag.DurationVar(&flagReportInterval, "r", defaultReportInterval, "REPORT_INTERVAL=<ЗНАЧЕНИЕ> частоту отправки метрик на сервер")
	flag.Parse()

	// создаем новую структуру конфигурации агента, заполняем ее значениями из флагов командной строки или из переменных окружения
	cfg := &AgentConfig{
		Address:        flagAddress,
		PollInterval:   flagPollInterval,
		ReportInterval: flagReportInterval,
	}

	// заполняем поля структуры значениями из переменных окружения
	if envAddress, ok := os.LookupEnv("ADDRESS"); ok && envAddress != "" {
		cfg.Address = envAddress
	}

	if envPollInterval, err := time.ParseDuration(os.Getenv("POLL_INTERVAL")); err == nil {
		cfg.PollInterval = envPollInterval
	}

	if envReportInterval, err := time.ParseDuration(os.Getenv("REPORT_INTERVAL")); err == nil {
		cfg.ReportInterval = envReportInterval
	}

	return cfg // возвращаем указатель на созданную структуру конфигурации агента
}
