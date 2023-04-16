package main

import (
	"log"

	"github.com/MakeItBright/go-metrics-devops/internal/agent"
	"github.com/MakeItBright/go-metrics-devops/internal/config"
)

func main() {
	// создаем новую структуру конфигурации агента
	cfg := config.NewAgentConfig()
	// запускаем агента с заданной конфигурацией
	if err := agent.RunAgent(cfg); err != nil {
		log.Fatal(err)
	}

}
