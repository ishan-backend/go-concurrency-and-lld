package main

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
	"vcs/models"
)

// TaskScheduler manages task scheduling and execution
type TaskScheduler struct {
	taskQueue    *models.TaskQueue
	tasks        map[string]*models.Task // taskID -> Task for quick lookup
	mu           sync.RWMutex            // Protects taskQueue and tasks
	workers      int                     // Number of worker goroutines
	workerPool   chan struct{}           // Semaphore for limiting concurrent executions
	stopChan     chan struct{}           // Signal to stop the scheduler
	wg           sync.WaitGroup          // Wait group for graceful shutdown
	isRunning    bool
	runningMu    sync.RWMutex
}

// NewTaskScheduler creates a new task scheduler
func NewTaskScheduler(workers int) *TaskScheduler {
	if workers <= 0 {
		workers = 1 // Default to 1 worker
	}

	return &TaskScheduler{
		taskQueue:  models.NewTaskQueue(),
		tasks:      make(map[string]*models.Task),
		workers:    workers,
		workerPool: make(chan struct{}, workers),
		stopChan:   make(chan struct{}),
		isRunning:  false,
	}
}

// Start starts the scheduler
func (ts *TaskScheduler) Start() error {
	ts.runningMu.Lock()
	if ts.isRunning {
		ts.runningMu.Unlock()
		return fmt.Errorf("scheduler is already running")
	}
	ts.isRunning = true
	ts.runningMu.Unlock()

	ts.wg.Add(1)
	go ts.run()
	return nil
}

// Stop stops the scheduler gracefully
func (ts *TaskScheduler) Stop() {
	ts.runningMu.Lock()
	if !ts.isRunning {
		ts.runningMu.Unlock()
		return
	}
	ts.isRunning = false
	ts.runningMu.Unlock()

	close(ts.stopChan)
	ts.wg.Wait()
}

// ScheduleTask schedules a task for execution
func (ts *TaskScheduler) ScheduleTask(task *models.Task) error {
	if task == nil {
		return fmt.Errorf("task cannot be nil")
	}

	if task.ID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}

	if task.Handler == nil {
		return fmt.Errorf("task handler cannot be nil")
	}

	ts.mu.Lock()
	defer ts.mu.Unlock()

	// Check if task already exists
	if _, exists := ts.tasks[task.ID]; exists {
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}

	// Add to map
	ts.tasks[task.ID] = task

	// Add to priority queue
	heap.Push(ts.taskQueue, task)

	return nil
}

// CancelTask cancels a scheduled task
func (ts *TaskScheduler) CancelTask(taskID string) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if _, exists := ts.tasks[taskID]; !exists {
		return fmt.Errorf("task with ID %s not found", taskID)
	}

	// Remove from map
	delete(ts.tasks, taskID)

	// Rebuild queue without the cancelled task
	newQueue := models.NewTaskQueue()
	for _, task := range *ts.taskQueue {
		if task.ID != taskID {
			heap.Push(newQueue, task)
		}
	}
	ts.taskQueue = newQueue

	return nil
}

// GetTaskStatus returns the status of a task
func (ts *TaskScheduler) GetTaskStatus(taskID string) (*models.Task, error) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	task, exists := ts.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task with ID %s not found", taskID)
	}

	// Return a copy to avoid race conditions
	return &models.Task{
		ID:        task.ID,
		Name:      task.Name,
		TaskType:  task.TaskType,
		ExecuteAt: task.ExecuteAt,
		Interval:  task.Interval,
		CreatedAt: task.CreatedAt,
		LastRun:   task.LastRun,
		RunCount:  task.RunCount,
	}, nil
}

// GetPendingTasksCount returns the number of pending tasks
func (ts *TaskScheduler) GetPendingTasksCount() int {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return len(*ts.taskQueue)
}

// run is the main scheduler loop
func (ts *TaskScheduler) run() {
	defer ts.wg.Done()

	for {
		select {
		case <-ts.stopChan:
			return
		default:
			ts.processTasks()
		}
	}
}

// processTasks processes tasks that are ready to execute
func (ts *TaskScheduler) processTasks() {
	ts.mu.Lock()

	// Check if there are any tasks
	if ts.taskQueue.Len() == 0 {
		ts.mu.Unlock()
		time.Sleep(100 * time.Millisecond) // Sleep when no tasks
		return
	}

	now := time.Now()
	nextTask := ts.taskQueue.Peek()

	// If the next task is not ready, wait
	if nextTask.ExecuteAt.After(now) {
		waitDuration := nextTask.ExecuteAt.Sub(now)
		ts.mu.Unlock()

		// Sleep until the next task is ready, but check stopChan periodically
		select {
		case <-ts.stopChan:
			return
		case <-time.After(waitDuration):
			// Continue to process tasks
		}
		return
	}

	// Pop the task that's ready to execute
	task := heap.Pop(ts.taskQueue).(*models.Task)

	// Remove from tasks map if it's a one-time task
	if !task.ShouldReschedule() {
		delete(ts.tasks, task.ID)
	}

	ts.mu.Unlock()

	// Execute task in a worker goroutine (non-blocking)
	ts.wg.Add(1)
	go ts.executeTask(task)
}

// executeTask executes a task
func (ts *TaskScheduler) executeTask(task *models.Task) {
	defer ts.wg.Done()

	// Acquire worker slot
	ts.workerPool <- struct{}{}
	defer func() { <-ts.workerPool }()

	// Execute the task
	task.Execute()

	// Reschedule if it's a recurring task
	if task.ShouldReschedule() {
		ts.mu.Lock()
		// Check if task still exists (might have been cancelled)
		if _, exists := ts.tasks[task.ID]; exists {
			// Update execution time for next run
			task.ExecuteAt = task.NextExecutionTime()
			// Add back to queue
			heap.Push(ts.taskQueue, task)
		}
		ts.mu.Unlock()
	}
}

// IsRunning returns whether the scheduler is running
func (ts *TaskScheduler) IsRunning() bool {
	ts.runningMu.RLock()
	defer ts.runningMu.RUnlock()
	return ts.isRunning
}

