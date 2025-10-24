package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// Task represents a unit of work
type Task struct {
	ID        int
	Data      string
	Timestamp time.Time
}

// Queue is a thread-safe fixed-size queue
type Queue struct {
	items    chan Task
	capacity int
	mu       sync.RWMutex
	size     int
}

// NewQueue creates a new fixed-size queue
func NewQueue(capacity int) *Queue {
	return &Queue{
		items:    make(chan Task, capacity),
		capacity: capacity,
	}
}

// Enqueue adds a task to the queue (blocks if full)
func (q *Queue) Enqueue(task Task) error {
	select {
	case q.items <- task:
		q.mu.Lock()
		q.size++
		q.mu.Unlock()
		return nil
	default:
		return fmt.Errorf("queue is full")
	}
}

// EnqueueBlocking adds a task to the queue (waits if full)
func (q *Queue) EnqueueBlocking(ctx context.Context, task Task) error {
	select {
	case q.items <- task:
		q.mu.Lock()
		q.size++
		q.mu.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Dequeue removes and returns a task from the queue
func (q *Queue) Dequeue(ctx context.Context) (Task, error) {
	select {
	case task := <-q.items:
		q.mu.Lock()
		q.size--
		q.mu.Unlock()
		return task, nil
	case <-ctx.Done():
		return Task{}, ctx.Err()
	}
}

// Size returns the current size of the queue
func (q *Queue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.size
}

// Close closes the queue
func (q *Queue) Close() {
	close(q.items)
}

// Producer generates tasks and adds them to the queue
type Producer struct {
	id    int
	queue *Queue
	wg    *sync.WaitGroup
}

// NewProducer creates a new producer
func NewProducer(id int, queue *Queue, wg *sync.WaitGroup) *Producer {
	return &Producer{
		id:    id,
		queue: queue,
		wg:    wg,
	}
}

// Start begins producing tasks
func (p *Producer) Start(ctx context.Context, taskCount int) {
	defer p.wg.Done()

	for i := 0; i < taskCount; i++ {
		select {
		case <-ctx.Done():
			log.Printf("Producer %d stopped due to context cancellation", p.id)
			return
		default:
			task := Task{
				ID:        i,
				Data:      fmt.Sprintf("Producer-%d-Task-%d", p.id, i),
				Timestamp: time.Now(),
			}

			// Try to enqueue with blocking
			err := p.queue.EnqueueBlocking(ctx, task)
			if err != nil {
				log.Printf("Producer %d failed to enqueue task %d: %v", p.id, i, err)
				return
			}

			log.Printf("Producer %d enqueued: %s", p.id, task.Data)

			// Simulate production time
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
		}
	}

	log.Printf("Producer %d finished producing %d tasks", p.id, taskCount)
}

// Worker processes tasks from the queue
type Worker struct {
	id    int
	queue *Queue
}

// NewWorker creates a new worker
func NewWorker(id int, queue *Queue) *Worker {
	return &Worker{
		id:    id,
		queue: queue,
	}
}

// Start begins processing tasks
func (w *Worker) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d stopped", w.id)
			return
		default:
			task, err := w.queue.Dequeue(ctx)
			if err != nil {
				// Context cancelled or queue closed
				return
			}

			w.processTask(task)
		}
	}
}

// processTask simulates task processing
func (w *Worker) processTask(task Task) {
	log.Printf("Worker %d processing: %s", w.id, task.Data)

	// Simulate processing time
	processingTime := time.Millisecond * time.Duration(rand.Intn(200)+100)
	time.Sleep(processingTime)

	log.Printf("Worker %d completed: %s (took %v)", w.id, task.Data, processingTime)
}

// WorkerPool manages a pool of workers
type WorkerPool struct {
	workers []*Worker
	queue   *Queue
	size    int
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(size int, queue *Queue) *WorkerPool {
	workers := make([]*Worker, size)
	for i := 0; i < size; i++ {
		workers[i] = NewWorker(i+1, queue)
	}

	return &WorkerPool{
		workers: workers,
		queue:   queue,
		size:    size,
	}
}

// Start begins all workers in the pool
func (wp *WorkerPool) Start(ctx context.Context) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(wp.size)

	for _, worker := range wp.workers {
		go worker.Start(ctx, &wg)
	}

	return &wg
}

func main() {
	// Configuration
	queueSize := 10
	numProducers := 3
	numWorkers := 5
	tasksPerProducer := 20

	// Create queue
	queue := NewQueue(queueSize)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start worker pool
	log.Println("Starting worker pool...")
	workerPool := NewWorkerPool(numWorkers, queue)
	workerWg := workerPool.Start(ctx)

	// Start producers
	log.Println("Starting producers...")
	var producerWg sync.WaitGroup
	producerWg.Add(numProducers)

	for i := 0; i < numProducers; i++ {
		producer := NewProducer(i+1, queue, &producerWg)
		go producer.Start(ctx, tasksPerProducer)
	}

	// Wait for all producers to finish
	producerWg.Wait()
	log.Println("All producers finished")

	// Give workers time to process remaining tasks
	time.Sleep(2 * time.Second)

	// Cancel context to stop workers
	cancel()

	// Wait for all workers to stop
	workerWg.Wait()
	log.Println("All workers finished")

	// Close queue
	queue.Close()
	log.Println("System shutdown complete")
}
