package workerpool

import (
	"fmt"
	"sync"
)

var (
	once sync.Once
	pool *WorkerPool
)

type WorkerPool struct {
	count     int
	wg        sync.WaitGroup
	workQueue chan Runnable
}

// NewWorkerPool creates a new worker pool with the specified number of workers.
// The workQueue is buffered to improve throughput and reduce blocking when adding tasks.
// Default buffer size is 1000, which allows queuing tasks without blocking while workers process them.
func NewWorkerPool(workersCount int) *WorkerPool {
	return NewWorkerPoolWithBuffer(workersCount, 1000)
}

// NewWorkerPoolWithBuffer creates a new worker pool with specified worker count and buffer size.
// bufferSize determines how many tasks can be queued without blocking the Add() method.
// For millions of tasks, a buffer size of 1000-10000 is recommended to balance memory usage and throughput.
func NewWorkerPoolWithBuffer(workersCount int, bufferSize int) *WorkerPool {
	once.Do(func() {
		if bufferSize <= 0 {
			bufferSize = 1000 // Default buffer size
		}
		pool = &WorkerPool{
			count:     workersCount,
			workQueue: make(chan Runnable, bufferSize),
		}

		for i := 0; i < pool.count; i++ {
			pool.wg.Add(1)
			fmt.Printf("worker initiated with index %v\n", i)

			go func(i int) {
				defer pool.wg.Done()

				for work := range pool.workQueue {
					work.Run()
				}
			}(i)
		}
	})

	return pool
}

func (w *WorkerPool) Add(work Runnable) {
	w.workQueue <- work
}

func (w *WorkerPool) Wait() {
	close(w.workQueue)
	w.wg.Wait()
}
