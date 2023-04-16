package main

import (
	"log"

	"github.com/MakeItBright/go-metrics-devops/internal/agent"
	"github.com/MakeItBright/go-metrics-devops/internal/config"
)

func main() {

	cfg := config.NewAgentConfig()
	if err := agent.RunAgent(cfg); err != nil {
		log.Fatal(err)
	}

}
