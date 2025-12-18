package models

import (
	"time"
)

// TaskType represents the type of task scheduling
type TaskType int

const (
	// OneTimeTask executes once at a specific time
	OneTimeTask TaskType = iota
	// RecurringTask executes repeatedly at fixed intervals
	RecurringTask
)

// Task represents a task to be executed by the scheduler
type Task struct {
	ID          string
	Name        string
	TaskType    TaskType
	ExecuteAt   time.Time      // When to execute the task
	Interval    time.Duration  // For recurring tasks: interval between executions
	Handler     func()         // The function to execute
	CreatedAt   time.Time
	LastRun     *time.Time     // Last execution time (nil if never run)
	RunCount    int            // Number of times executed
}

// NewOneTimeTask creates a new one-time task
func NewOneTimeTask(id, name string, executeAt time.Time, handler func()) *Task {
	return &Task{
		ID:        id,
		Name:      name,
		TaskType:  OneTimeTask,
		ExecuteAt: executeAt,
		Handler:   handler,
		CreatedAt: time.Now(),
		RunCount:  0,
	}
}

// NewRecurringTask creates a new recurring task
func NewRecurringTask(id, name string, executeAt time.Time, interval time.Duration, handler func()) *Task {
	return &Task{
		ID:        id,
		Name:      name,
		TaskType:  RecurringTask,
		ExecuteAt: executeAt,
		Interval:  interval,
		Handler:   handler,
		CreatedAt: time.Now(),
		RunCount:  0,
	}
}

// NewDelayedTask creates a task that executes after a delay
func NewDelayedTask(id, name string, delay time.Duration, handler func()) *Task {
	return &Task{
		ID:        id,
		Name:      name,
		TaskType:  OneTimeTask,
		ExecuteAt: time.Now().Add(delay),
		Handler:   handler,
		CreatedAt: time.Now(),
		RunCount:  0,
	}
}

// NewRecurringTaskWithDelay creates a recurring task that starts after a delay
func NewRecurringTaskWithDelay(id, name string, initialDelay time.Duration, interval time.Duration, handler func()) *Task {
	return &Task{
		ID:        id,
		Name:      name,
		TaskType:  RecurringTask,
		ExecuteAt: time.Now().Add(initialDelay),
		Interval:  interval,
		Handler:   handler,
		CreatedAt: time.Now(),
		RunCount:  0,
	}
}

// Execute runs the task handler
func (t *Task) Execute() {
	now := time.Now()
	t.LastRun = &now
	t.RunCount++
	if t.Handler != nil {
		t.Handler()
	}
}

// IsRecurring returns true if the task is recurring
func (t *Task) IsRecurring() bool {
	return t.TaskType == RecurringTask
}

// NextExecutionTime calculates the next execution time for recurring tasks
func (t *Task) NextExecutionTime() time.Time {
	if t.TaskType == RecurringTask {
		if t.LastRun != nil {
			return t.LastRun.Add(t.Interval)
		}
		return t.ExecuteAt
	}
	return t.ExecuteAt
}

// ShouldReschedule returns true if the task should be rescheduled
func (t *Task) ShouldReschedule() bool {
	return t.TaskType == RecurringTask
}

