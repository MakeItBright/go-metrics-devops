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

// Не экспортированные переменные, flagAddress, flagPollInterval, flagReportInterval,
// содержат значения флагов командной строки.
var (
	flagAddress        string
	flagPollInterval   int
	flagReportInterval int
)

func init() {
	flag.StringVar(&flagAddress, "a", defaultAddress, "адрес эндпоинта HTTP-сервера")
	flag.IntVar(&flagPollInterval, "p", defaultPollInterval, "частота опроса метрик ")
	flag.IntVar(&flagReportInterval, "r", defaultReportInterval, "частота отправки метрик на сервер")

}

func main() {
	flag.Parse()
	envParse()

	if err := agent.Start(agent.Config{
		Scheme:         defaultScheme,
		Address:        flagAddress,
		PollInterval:   time.Duration(flagPollInterval) * time.Second,
		ReportInterval: time.Duration(flagReportInterval) * time.Second,
		Logger:         logrus.New(),
	}); err != nil {
		log.Fatal(err)
	}

	//TODO: add interrupt ctrl+c
}

func envParse() {
	if envAddress, ok := os.LookupEnv("ADDRESS"); ok && envAddress != "" {
		flagAddress = envAddress
	}

	if envPollInterval, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		if intPollInterval, err := strconv.Atoi(envPollInterval); err != nil {
			log.Fatal(err)
		} else {
			flagPollInterval = intPollInterval
		}
	}

	if envReportInterval, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		if intReportInterval, err := strconv.Atoi(envReportInterval); err != nil {
			log.Fatal(err)
		} else {
			flagReportInterval = intReportInterval
		}
	}

}
