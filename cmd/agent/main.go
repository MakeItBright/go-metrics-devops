package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/MakeItBright/go-metrics-devops/internal/agent"
	"github.com/sirupsen/logrus"
)

const (
	defaultAddress        = "localhost:8080"
	defaultScheme         = "http"
	defaultPollInterval   = 2  // интервал обновления метрик по умолчанию
	defaultReportInterval = 10 // интервал отправки метрик на сервер по умолчанию
)

// неэкспортированные переменные flagAddress, flagPollInterval, flagReportInterval,
// содержат значения флагов командной строки
var (
	flagAddress        string
	flagPollInterval   int
	flagReportInterval int
)

func init() {
	// устанавливаем значения флагов командной строки
	flag.StringVar(&flagAddress, "a", defaultAddress, "адрес эндпоинта HTTP-сервера")
	flag.IntVar(&flagPollInterval, "p", defaultPollInterval, "частота опроса метрик ")
	flag.IntVar(&flagReportInterval, "r", defaultReportInterval, "частоту отправки метрик на сервер")
}

func main() {
	flag.Parse()
	envParse()

	// запускаем агента с заданной конфигурацией
	if err := agent.Run(agent.Config{
		Address:        flagAddress,
		PollInterval:   time.Duration(flagPollInterval) * time.Second,
		ReportInterval: time.Duration(flagReportInterval) * time.Second,
		Logger:         logrus.New(),
	}); err != nil {
		log.Fatal(err)
	}

	//TODO: add intercept ctrl+c
}

func envParse() {
	if envAddress, ok := os.LookupEnv("ADDRESS"); ok && envAddress != "" {
		flagAddress = envAddress
	}

	if envPollInterval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		intPollInterval, err := strconv.Atoi(envPollInterval)
		if err != nil {
			log.Fatal(err)
		}
		flagPollInterval = intPollInterval
	}

	if envReportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		intReportInterval, err := strconv.Atoi(envReportInterval)
		if err != nil {
			log.Fatal(err)
		}
		flagReportInterval = intReportInterval
	}

}
