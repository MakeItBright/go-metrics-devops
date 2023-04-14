package server

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/MakeItBright/go-metrics-devops/internal/storage/inmem"
)

var (
	flagRunAddr string // неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
)

func init() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением localhost:8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
}

// Start
func Start(config *Config) error {
	flag.Parse()
	config.BindAddr = flagRunAddr
	store := inmem.New()

	srv := newServer(store)
	fmt.Println("Running server on", flagRunAddr)
	return http.ListenAndServe(config.BindAddr, srv)
}
