package main

import (
	"log"

	"github.com/MakeItBright/go-metrics-devops/internal/server"
)

func main() {
	config := server.NewConfig()

	if err := server.Start(config); err != nil {
		log.Fatal(err)
	}
}
