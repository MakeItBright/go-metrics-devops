package server

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/MakeItBright/go-metrics-devops/internal/storage/inmem"
)

var (
	flagRunAddr string // неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
)

func init() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением localhost:8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
}

// Start
func Start(config *Config) error {
	// Приоритет параметров должен быть таким:
	// Если указана переменная окружения, то используется она.
	// Если нет переменной окружения, но есть аргумент командной строки (флаг), то используется он.
	// Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.
	flag.Parse()
	config.BindAddr = flagRunAddr

	if envBindAddr := os.Getenv("ADDRESS"); envBindAddr != "" {
		config.BindAddr = envBindAddr
	}
	store := inmem.New()

	srv := newServer(store)
	fmt.Println("Running server on", flagRunAddr)
	return http.ListenAndServe(config.BindAddr, srv)
}
