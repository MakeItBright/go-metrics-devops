package main

import (
	"log"

	"github.com/MakeItBright/go-metrics-devops/internal/server"
)

func main() {
	config := server.NewConfig()
	s := server.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
