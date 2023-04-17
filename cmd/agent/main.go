package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/MakeItBright/go-metrics-devops/internal/agent"
	"github.com/sirupsen/logrus"
)

const (
	defaultAddress        = "localhost:8080"
	defaultScheme         = "http"
	defaultPollInterval   = 2 * time.Second  // интервал обновления метрик по умолчанию
	defaultReportInterval = 10 * time.Second // интервал отправки метрик на сервер по умолчанию
)

// неэкспортированные переменные flagAddress, flagPollInterval, flagReportInterval,
// содержат значения флагов командной строки
var (
	flagAddress        string
	flagPollInterval   time.Duration
	flagReportInterval time.Duration
)

func init() {
	// устанавливаем значения флагов командной строки
	flag.StringVar(&flagAddress, "a", defaultAddress, "адрес эндпоинта HTTP-сервера")
	flag.DurationVar(&flagPollInterval, "p", defaultPollInterval, "частота опроса метрик ")
	flag.DurationVar(&flagReportInterval, "r", defaultReportInterval, "частоту отправки метрик на сервер")
}

func main() {
	flag.Parse()
	envParse()

	// запускаем агента с заданной конфигурацией
	if err := agent.Run(agent.Config{
		Address:        flagAddress,
		PollInterval:   flagPollInterval,
		ReportInterval: flagReportInterval,
		Logger:         logrus.New(),
	}); err != nil {
		log.Fatal(err)
	}

	//TODO: add interrapt ctrl+c
}

func envParse() {
	// заполняем поля структуры значениями из переменных окружения
	if envAddress, ok := os.LookupEnv("ADDRESS"); ok && envAddress != "" {
		flagAddress = envAddress
	}

	if envPollInterval, err := time.ParseDuration(os.Getenv("POLL_INTERVAL")); err == nil {
		flagPollInterval = envPollInterval
	}

	if envReportInterval, err := time.ParseDuration(os.Getenv("REPORT_INTERVAL")); err == nil {
		flagReportInterval = envReportInterval
	}
}
