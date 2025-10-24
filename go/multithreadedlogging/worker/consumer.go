package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/logger/factory"
	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/service"
)

// Consumer represents a message consumer worker
type Consumer struct {
	queueService  service.QueueService
	loggerFactory *factory.InternalLoggerFactory
}

// NewConsumer creates a new Consumer instance
func NewConsumer(queueService service.QueueService) *Consumer {
	return &Consumer{
		queueService:  queueService,
		loggerFactory: factory.NewInternalLoggerFactory(),
	}
}

// Run starts the consumer worker
func (c *Consumer) Run(ctx context.Context) {
	fmt.Println("Consumer started")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Consumer shutting down")
			return
		default:
			// Consume message from queue
			message, err := c.queueService.Consume(ctx)
			if err != nil {
				if err == context.Canceled {
					return
				}
				fmt.Printf("Error consuming message: %v\n", err)
				continue
			}

			if message != nil {
				// Get appropriate logger for the medium
				logger := c.loggerFactory.GetLoggerInstance(message.Medium)
				if logger != nil {
					logger.Log(message)
				}

				// Simulate processing time
				time.Sleep(500 * time.Millisecond)
			}
		}
	}
}
