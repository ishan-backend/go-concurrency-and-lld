package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/service"
	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/worker"
)

func main() {
	test1()
}

func test1() {
	// Create queue service with capacity 100
	queueService := service.NewQueueService(100)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Start producer goroutines
	wg.Add(2)
	go func() {
		defer wg.Done()
		producer := worker.NewProducer(queueService)
		producer.Run(ctx)
	}()

	go func() {
		defer wg.Done()
		producer := worker.NewProducer(queueService)
		producer.Run(ctx)
	}()

	// Start consumer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer := worker.NewConsumer(queueService)
		consumer.Run(ctx)
	}()

	// Run for 30 seconds then shutdown
	time.Sleep(30 * time.Second)
	cancel()

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("Application shutdown complete")
}
