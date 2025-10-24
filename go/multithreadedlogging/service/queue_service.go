package service

import (
	"context"
	"time"

	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/configuration"
	"github.com/ishanpandey/lowleveldesign/multithreadedlogging/model"
)

// QueueService interface defines the contract for message queue operations
type QueueService interface {
	Publish(ctx context.Context, messages []*model.WritableMessage) error
	Consume(ctx context.Context) (*model.WritableMessage, error)
}

// queueServiceImpl implements QueueService using Go channels
// This is Go's equivalent to Java's ArrayBlockingQueue
type queueServiceImpl struct {
	queue    chan *model.WritableMessage
	capacity int
	config   *configuration.Config
}

// NewQueueService creates a new QueueService instance
func NewQueueService(capacity int) QueueService {
	config := configuration.DefaultConfig()
	if capacity <= 0 {
		capacity = config.DefaultCapacity
	}

	return &queueServiceImpl{
		queue:    make(chan *model.WritableMessage, capacity),
		capacity: capacity,
		config:   config,
	}
}

// Publish adds messages to the queue
// In Go, this is equivalent to ArrayBlockingQueue.offer() with timeout
func (q *queueServiceImpl) Publish(ctx context.Context, messages []*model.WritableMessage) error {
	if messages == nil {
		return nil
	}

	for _, message := range messages {
		select {
		case q.queue <- message:
			// Successfully added to queue
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(q.config.DefaultPublishTimeoutMs) * time.Millisecond):
			// Timeout - equivalent to Java's offer() returning false
			// In a production system, you might want to handle this differently
			continue
		}
	}
	return nil
}

// Consume retrieves a message from the queue
// In Go, this is equivalent to ArrayBlockingQueue.poll() with timeout
func (q *queueServiceImpl) Consume(ctx context.Context) (*model.WritableMessage, error) {
	select {
	case message := <-q.queue:
		return message, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(time.Duration(q.config.DefaultConsumptionTimeoutMs) * time.Millisecond):
		// Timeout - equivalent to Java's poll() returning null
		return nil, nil
	}
}
