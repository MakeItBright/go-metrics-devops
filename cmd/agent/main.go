package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/MakeItBright/go-metrics-devops/internal/agent"
	"github.com/MakeItBright/go-metrics-devops/internal/logger"
)

const (
	addressDefault        = "localhost:8080"
	schemeDefault         = "http"
	pollIntervalDefault   = 2  // интервал обновления метрик по умолчанию
	reportIntervalDefault = 10 // интервал отправки метрик на сервер по умолчанию
)

func main() {
	cfg := agent.Config{
		Scheme:         schemeDefault,
		Address:        addressDefault,
		PollInterval:   pollIntervalDefault,
		ReportInterval: reportIntervalDefault,
	}

	flagParse(&cfg)
	if err := envParse(&cfg); err != nil {
		log.Fatalf("cannot parse ENV variables: %s", err)
	}

	if err := logger.Initialize("info"); err != nil {
		log.Fatalf("cannot start agent: %s", err)

	}

	if err := agent.Start(cfg); err != nil {
		log.Fatalf("cannot start agent: %s", err)
	}

	//TODO: add interrupt ctrl+c
}

func flagParse(cfg *agent.Config) {
	flag.StringVar(&cfg.Address, "a", addressDefault, "адрес эндпоинта HTTP-сервера")
	flag.IntVar(&cfg.PollInterval, "p", pollIntervalDefault, "частота опроса метрик ")
	flag.IntVar(&cfg.ReportInterval, "r", reportIntervalDefault, "частота отправки метрик на сервер")
	flag.Parse()
}

func envParse(cfg *agent.Config) error {
	if envAddress, ok := os.LookupEnv("ADDRESS"); ok && envAddress != "" {
		cfg.Address = envAddress
	}

	if pollIntervalEnv, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		if pollIntervalInt, err := strconv.Atoi(pollIntervalEnv); err != nil {
			return fmt.Errorf("cannot parse POLL_INTERVAL: %w", err)
		} else {
			cfg.PollInterval = pollIntervalInt
		}
	}

	if reportIntervalEnv, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		if reportIntervalInt, err := strconv.Atoi(reportIntervalEnv); err != nil {
			return fmt.Errorf("cannot parse REPORT_INTERVAL: %w", err)
		} else {
			cfg.ReportInterval = reportIntervalInt
		}
	}
	return nil
}
