package main

import (
	"flag"
	"log"
	"os"

	"github.com/MakeItBright/go-metrics-devops/internal/server"
)

func main() {
	cfg := server.Config{}

	flagParse(&cfg)
	envParse(&cfg)

	if err := server.Start(cfg); err != nil {
		log.Fatalf("cannot start server: %s", err)
	}

}

func flagParse(cfg *server.Config) {
	flag.StringVar(&cfg.BindAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

}

func envParse(cfg *server.Config) {
	if envAddress, ok := os.LookupEnv("ADDRESS"); ok && envAddress != "" {
		cfg.BindAddr = envAddress
	}

}
