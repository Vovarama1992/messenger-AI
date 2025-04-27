package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Vovarama1992/go-ai-service/db"
	"github.com/Vovarama1992/go-ai-service/kafka"
)

var wg sync.WaitGroup

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		if r := recover(); r != nil {
			log.Printf("🔥 Panic recovered: %v", r)
			cancel()
			wg.Wait()
			log.Println("✅ Все горутины завершены после panic")
		}
	}()

	defer cancel()

	if err := db.InitDB(); err != nil {
		log.Fatalf("DB init failed: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	kafka.StartAdviceWorkers(ctx, &wg, 3)
	kafka.StartAdviceGPTWorker(ctx, &wg, 3)
	kafka.StartAdviceProducerWorker(ctx, &wg, 3)

	kafka.StartAutoreplyWorkers(ctx, &wg, 3)
	kafka.StartAutoreplyGPTWorkers(ctx, &wg, 3)
	kafka.StartAutoReplySenderWorkers(ctx, &wg, 3)

	<-stop
	log.Println("⛔ Завершение по сигналу...")

	cancel()
	wg.Wait()
	log.Println("✅ Все горутины завершены")
}
