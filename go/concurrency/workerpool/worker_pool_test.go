package workerpool

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// SimpleTask is a simple task for benchmarking
type SimpleTask struct {
	work func()
}

func (t *SimpleTask) Run() {
	if t.work != nil {
		t.work()
	}
}

// newWorkerPoolForTesting creates a new worker pool for testing purposes.
// This bypasses the singleton pattern to allow multiple pools in benchmarks.
func newWorkerPoolForTesting(workersCount int, bufferSize int) *WorkerPool {
	if bufferSize <= 0 {
		bufferSize = 1000
	}

	pool := &WorkerPool{
		count:     workersCount,
		workQueue: make(chan Runnable, bufferSize),
	}

	for i := 0; i < pool.count; i++ {
		pool.wg.Add(1)
		go func(i int) {
			defer pool.wg.Done()
			for work := range pool.workQueue {
				work.Run()
			}
		}(i)
	}

	return pool
}

// BenchmarkWorkerPool_AddAndWait benchmarks adding tasks and waiting for completion
func BenchmarkWorkerPool_AddAndWait(b *testing.B) {
	pool := newWorkerPoolForTesting(4, 1000)
	defer pool.Wait()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task := &SimpleTask{
			work: func() {
				// Simulate some work
				time.Sleep(1 * time.Microsecond)
			},
		}
		pool.Add(task)
	}
}

// BenchmarkWorkerPool_DifferentWorkerCounts benchmarks with different worker counts
func BenchmarkWorkerPool_DifferentWorkerCounts(b *testing.B) {
	workerCounts := []int{1, 2, 4, 8, 16, 32}

	for _, workers := range workerCounts {
		b.Run(fmt.Sprintf("Workers_%d", workers), func(b *testing.B) {
			pool := newWorkerPoolForTesting(workers, 1000)
			var wg sync.WaitGroup
			wg.Add(b.N)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				task := &SimpleTask{
					work: func() {
						time.Sleep(10 * time.Microsecond)
						wg.Done()
					},
				}
				pool.Add(task)
			}
			wg.Wait()
			pool.Wait()
		})
	}
}

// BenchmarkWorkerPool_DifferentBufferSizes benchmarks with different buffer sizes
func BenchmarkWorkerPool_DifferentBufferSizes(b *testing.B) {
	bufferSizes := []int{0, 100, 500, 1000, 5000, 10000}

	for _, bufferSize := range bufferSizes {
		b.Run(fmt.Sprintf("Buffer_%d", bufferSize), func(b *testing.B) {
			pool := newWorkerPoolForTesting(4, bufferSize)
			var wg sync.WaitGroup
			wg.Add(b.N)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				task := &SimpleTask{
					work: func() {
						time.Sleep(10 * time.Microsecond)
						wg.Done()
					},
				}
				pool.Add(task)
			}
			wg.Wait()
			pool.Wait()
		})
	}
}

// BenchmarkWorkerPool_Throughput measures tasks processed per second
func BenchmarkWorkerPool_Throughput(b *testing.B) {
	pool := newWorkerPoolForTesting(8, 5000)
	var completed int64

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			task := &SimpleTask{
				work: func() {
					atomic.AddInt64(&completed, 1)
				},
			}
			pool.Add(task)
		}
	})

	pool.Wait()
	b.ReportMetric(float64(completed)/b.Elapsed().Seconds(), "tasks/sec")
}

// BenchmarkWorkerPool_HeavyWorkload benchmarks with CPU-intensive tasks
func BenchmarkWorkerPool_HeavyWorkload(b *testing.B) {
	pool := newWorkerPoolForTesting(8, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task := &SimpleTask{
			work: func() {
				// Simulate CPU-intensive work
				sum := 0
				for j := 0; j < 1000; j++ {
					sum += j
				}
				_ = sum
			},
		}
		pool.Add(task)
	}
	pool.Wait()
}

// BenchmarkWorkerPool_IOBoundWorkload benchmarks with I/O-bound tasks
func BenchmarkWorkerPool_IOBoundWorkload(b *testing.B) {
	pool := newWorkerPoolForTesting(8, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task := &SimpleTask{
			work: func() {
				// Simulate I/O-bound work
				time.Sleep(100 * time.Microsecond)
			},
		}
		pool.Add(task)
	}
	pool.Wait()
}

// BenchmarkWorkerPool_ConcurrentAdds benchmarks concurrent task additions
func BenchmarkWorkerPool_ConcurrentAdds(b *testing.B) {
	pool := newWorkerPoolForTesting(8, 10000)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			task := &SimpleTask{
				work: func() {
					time.Sleep(10 * time.Microsecond)
				},
			}
			pool.Add(task)
		}
	})

	pool.Wait()
}

// BenchmarkStringComparison_WorkerPool benchmarks string comparison using worker pool
func BenchmarkStringComparison_WorkerPool(b *testing.B) {
	// Generate test strings
	strings := make([]string, 100)
	for i := 0; i < 100; i++ {
		strings[i] = "test_string_" + string(rune(i%10+'0'))
	}

	workerCounts := []int{2, 4, 8, 16}

	for _, workers := range workerCounts {
		b.Run(fmt.Sprintf("Workers_%d", workers), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				comparator := NewStringComparatorWithBuffer(workers, 5000, strings)
				comparator.GenerateComparisonTasks()
				comparator.WaitForCompletion()
			}
		})
	}
}

// BenchmarkStringComparison_DifferentSizes benchmarks with different string array sizes
func BenchmarkStringComparison_DifferentSizes(b *testing.B) {
	sizes := []int{10, 50, 100, 500}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			strings := make([]string, size)
			for i := 0; i < size; i++ {
				strings[i] = "string_" + string(rune(i%20+'0'))
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				comparator := NewStringComparatorWithBuffer(8, 5000, strings)
				comparator.GenerateComparisonTasks()
				comparator.WaitForCompletion()
			}
		})
	}
}

// BenchmarkWorkerPool_MemoryUsage measures memory allocations
func BenchmarkWorkerPool_MemoryUsage(b *testing.B) {
	pool := newWorkerPoolForTesting(4, 1000)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		task := &SimpleTask{
			work: func() {
				time.Sleep(1 * time.Microsecond)
			},
		}
		pool.Add(task)
	}
	pool.Wait()
}

// BenchmarkWorkerPool_Scaling measures performance scaling with task count
func BenchmarkWorkerPool_Scaling(b *testing.B) {
	taskCounts := []int{100, 1000, 10000, 100000}

	for _, taskCount := range taskCounts {
		b.Run(fmt.Sprintf("Tasks_%d", taskCount), func(b *testing.B) {
			pool := newWorkerPoolForTesting(8, 10000)
			var wg sync.WaitGroup
			wg.Add(taskCount)

			b.ResetTimer()
			for i := 0; i < taskCount; i++ {
				task := &SimpleTask{
					work: func() {
						time.Sleep(1 * time.Microsecond)
						wg.Done()
					},
				}
				pool.Add(task)
			}
			wg.Wait()
			pool.Wait()
		})
	}
}

// BenchmarkWorkerPool_BufferVsUnbuffered compares buffered vs unbuffered performance
func BenchmarkWorkerPool_BufferVsUnbuffered(b *testing.B) {
	// Note: We can't truly test unbuffered with current implementation,
	// but we can compare small vs large buffers
	bufferSizes := []int{1, 10, 100, 1000, 10000}

	for _, bufferSize := range bufferSizes {
		b.Run(fmt.Sprintf("Buffer_%d", bufferSize), func(b *testing.B) {
			pool := newWorkerPoolForTesting(8, bufferSize)
			var wg sync.WaitGroup
			wg.Add(b.N)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				task := &SimpleTask{
					work: func() {
						time.Sleep(10 * time.Microsecond)
						wg.Done()
					},
				}
				pool.Add(task)
			}
			wg.Wait()
			pool.Wait()
		})
	}
}
