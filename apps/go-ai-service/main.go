package main

import (
	"log"

	"github.com/Vovarama1992/go-ai-service/db"
	"github.com/Vovarama1992/go-ai-service/kafka"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatalf("DB init failed: %v", err)
	}

	kafka.StartAdviceWorkers(3)
	kafka.StartAutoreplyWorkers(2)

	select {}
}
