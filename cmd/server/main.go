package main

import (
	"flag"
	"log"
	"os"

	"github.com/MakeItBright/go-metrics-devops/internal/server"
)

var (
	flagRunAddr string // неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
)

func init() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением localhost:8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
}

func main() {
	flag.Parse()
	envParse()

	if err := server.Start(server.Config{
		BindAddr: flagRunAddr,
	}); err != nil {
		log.Fatal(err)
	}
}

func envParse() {
	if envAddress, ok := os.LookupEnv("ADDRESS"); ok && envAddress != "" {
		flagRunAddr = envAddress
	}
}
