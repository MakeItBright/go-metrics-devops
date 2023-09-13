package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/MakeItBright/go-metrics-devops/internal/logger"
	"github.com/MakeItBright/go-metrics-devops/internal/server"
)

func main() {
	cfg := server.Config{}

	flagParse(&cfg)
	if err := envParse(&cfg); err != nil {
		log.Fatalf("cannot parse ENV variables: %s", err)
	}

	if err := logger.Initialize("info"); err != nil {
		log.Fatalf("cannot start server: %s", err)
	}

	if err := server.Start(cfg); err != nil {
		log.Fatalf("cannot start server: %s", err)
	}

}

func flagParse(cfg *server.Config) {
	flag.StringVar(&cfg.BindAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&cfg.StoreInterval, "i", 300, "Interval in seconds for storing metrics (default: 300)")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "Path to the file for storing metrics (default: /tmp/metrics-db.json)")
	flag.BoolVar(&cfg.Restore, "r", true, "Restore previously saved metrics from file (default: true)")
	flag.Parse()

}

func envParse(cfg *server.Config) error {
	if addressEnv, ok := os.LookupEnv("ADDRESS"); ok && addressEnv != "" {
		cfg.BindAddr = addressEnv
	}

	if storeIntervalEnv, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		if storeIntervalInt, err := strconv.Atoi(storeIntervalEnv); err != nil {
			return fmt.Errorf("cannot parse STORE_INTERVAL: %w", err)
		} else {
			cfg.StoreInterval = storeIntervalInt
		}
	}

	if fileStoragePathEnv, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok && fileStoragePathEnv != "" {
		cfg.FileStoragePath = fileStoragePathEnv
	}

	if restoreEnv, ok := os.LookupEnv("RESTORE"); ok && restoreEnv != "" {
		if restoreValueBool, err := strconv.ParseBool(restoreEnv); err != nil {
			return fmt.Errorf("cannot parse RESTORE: %w", err)
		} else {
			cfg.Restore = restoreValueBool
		}
	}

	return nil
}
