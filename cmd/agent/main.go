package main

import (
	"log"

	"github.com/MakeItBright/go-metrics-devops/internal/agent"
)

func main() {

	config := agent.NewConfig()
	if err := agent.Start(config); err != nil {
		log.Fatal(err)
	}

}
