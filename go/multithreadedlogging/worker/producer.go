package worker

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/constant"
	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/logger/impl"
	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/model"
	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/service"
)

// Producer represents a message producer worker
type Producer struct {
	queueService service.QueueService
}

// NewProducer creates a new Producer instance
func NewProducer(queueService service.QueueService) *Producer {
	return &Producer{
		queueService: queueService,
	}
}

// Run starts the producer worker
func (p *Producer) Run(ctx context.Context) {
	fmt.Println("Producer started")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Producer shutting down")
			return
		default:
			// Create application-aware logger
			appLogger := impl.NewApplicationAwareLogger(p.queueService)

			// Generate messages for different log levels
			for i := 0; i < 100; i++ {
				select {
				case <-ctx.Done():
					return
				default:
					// Create messages for different mediums
					mediums := []constant.SupportedLogMedium{
						constant.FILE,
						constant.DB,
					}

					// Generate random messages
					rand.Seed(time.Now().UnixNano())

					// Info message
					infoMsg := model.NewApplicationAwareMessage(
						fmt.Sprintf("info - %d", rand.Intn(10000)),
						mediums,
					)
					appLogger.Info(ctx, infoMsg)

					// Debug message
					debugMsg := model.NewApplicationAwareMessage(
						fmt.Sprintf("debug - %d", rand.Intn(10000)),
						mediums,
					)
					appLogger.Debug(ctx, debugMsg)

					// Warning message
					warnMsg := model.NewApplicationAwareMessage(
						fmt.Sprintf("warn - %d", rand.Intn(10000)),
						mediums,
					)
					appLogger.Warn(ctx, warnMsg)

					// Error message
					errorMsg := model.NewApplicationAwareMessage(
						fmt.Sprintf("error - %d", rand.Intn(10000)),
						mediums,
					)
					appLogger.Error(ctx, errorMsg)

					// Sleep for 1 second between batches
					time.Sleep(1 * time.Second)
				}
			}
		}
	}
}
