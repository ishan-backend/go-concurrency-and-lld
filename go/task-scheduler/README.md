# Task Scheduler

A thread-safe, concurrent Task Scheduler implementation in Go that supports one-time tasks, recurring tasks, and delayed execution with configurable worker pools.

## Overview

This Task Scheduler is designed to manage the execution of tasks at predefined times or intervals. It's commonly used in operating systems, distributed systems, and backend services to automate jobs like backups, notifications, report generation, and periodic cleanup tasks.

## Features

### Functional Requirements

✅ **One-Time Tasks**: Schedule a task to execute once at a specific future time  
✅ **Recurring Tasks**: Schedule tasks that run repeatedly at fixed intervals (e.g., every 5 seconds)  
✅ **Delayed Execution**: Schedule tasks to run after an initial delay  
✅ **Concurrent Execution**: Execute multiple tasks concurrently using configurable worker threads  
✅ **Task Cancellation**: Cancel scheduled tasks before execution  
✅ **Task Status**: Query the status of scheduled tasks  

### Non-Functional Requirements

✅ **Thread-Safety**: Safe to use in multi-threaded environments with concurrent scheduling  
✅ **Efficiency**: Uses efficient sleep/wait mechanisms to avoid CPU spinning  
✅ **Robustness**: Long-running tasks don't block the scheduler from executing other tasks  
✅ **Extensibility**: Design allows adding new scheduling logic (e.g., CRON-based schedules)  

## Architecture

### Core Components

#### Models

- **Task**: Represents a task with ID, name, type, execution time, handler function, and metadata
- **TaskType**: Enum for OneTimeTask or RecurringTask
- **TaskQueue**: Priority queue (min-heap) that orders tasks by execution time

#### Scheduler

The `TaskScheduler` class manages:
- Priority queue for efficient task scheduling
- Worker pool for concurrent task execution
- Thread-safe operations using mutexes
- Task lifecycle management

### Design Patterns

1. **Priority Queue Pattern**: Uses a min-heap to efficiently find the next task to execute
2. **Worker Pool Pattern**: Limits concurrent executions using a semaphore
3. **Producer-Consumer Pattern**: Scheduler produces tasks, workers consume and execute them

## Usage

### Basic Operations

#### Create a Scheduler

```go
// Create scheduler with 3 worker threads
scheduler := NewTaskScheduler(3)
scheduler.Start()
defer scheduler.Stop()
```

#### Schedule One-Time Task

```go
task := models.NewOneTimeTask(
    "task-1",
    "My Task",
    time.Now().Add(5*time.Second),
    func() {
        fmt.Println("Task executed!")
    },
)
scheduler.ScheduleTask(task)
```

#### Schedule Delayed Task

```go
task := models.NewDelayedTask(
    "task-2",
    "Delayed Task",
    3*time.Second, // Execute after 3 seconds
    func() {
        fmt.Println("Delayed task executed!")
    },
)
scheduler.ScheduleTask(task)
```

#### Schedule Recurring Task

```go
task := models.NewRecurringTask(
    "task-3",
    "Recurring Task",
    time.Now().Add(1*time.Second), // Start in 1 second
    5*time.Second,                  // Repeat every 5 seconds
    func() {
        fmt.Println("Recurring task executed!")
    },
)
scheduler.ScheduleTask(task)
```

#### Schedule Recurring Task with Initial Delay

```go
task := models.NewRecurringTaskWithDelay(
    "task-4",
    "Delayed Recurring Task",
    10*time.Second, // Start after 10 seconds
    30*time.Second,  // Then repeat every 30 seconds
    func() {
        fmt.Println("Delayed recurring task executed!")
    },
)
scheduler.ScheduleTask(task)
```

#### Cancel a Task

```go
err := scheduler.CancelTask("task-1")
if err != nil {
    fmt.Printf("Error cancelling task: %v\n", err)
}
```

#### Get Task Status

```go
status, err := scheduler.GetTaskStatus("task-1")
if err == nil {
    fmt.Printf("Task RunCount: %d\n", status.RunCount)
    fmt.Printf("Next Execution: %s\n", status.ExecuteAt)
}
```

#### Check Pending Tasks

```go
count := scheduler.GetPendingTasksCount()
fmt.Printf("Pending tasks: %d\n", count)
```

### Example: Complete Workflow

```go
package main

import (
    "fmt"
    "time"
    "vcs/models"
)

func main() {
    // Create scheduler with 5 workers
    scheduler := NewTaskScheduler(5)
    scheduler.Start()
    defer scheduler.Stop()

    // Schedule a one-time task
    task1 := models.NewOneTimeTask(
        "backup-1",
        "Database Backup",
        time.Now().Add(10*time.Second),
        func() {
            fmt.Println("Performing database backup...")
            // Backup logic here
            fmt.Println("Backup completed!")
        },
    )
    scheduler.ScheduleTask(task1)

    // Schedule a recurring cleanup task
    task2 := models.NewRecurringTask(
        "cleanup",
        "Cleanup Task",
        time.Now().Add(1*time.Minute),
        1*time.Hour, // Run every hour
        func() {
            fmt.Println("Cleaning up temporary files...")
            // Cleanup logic here
        },
    )
    scheduler.ScheduleTask(task2)

    // Keep the program running
    time.Sleep(2 * time.Hour)
}
```

### Example: Long-Running Tasks

The scheduler handles long-running tasks without blocking:

```go
// Long-running task (5 seconds)
longTask := models.NewDelayedTask(
    "long-task",
    "Long Task",
    1*time.Second,
    func() {
        fmt.Println("Starting long task...")
        time.Sleep(5 * time.Second)
        fmt.Println("Long task completed!")
    },
)
scheduler.ScheduleTask(longTask)

// Quick task that will run concurrently
quickTask := models.NewDelayedTask(
    "quick-task",
    "Quick Task",
    2*time.Second,
    func() {
        fmt.Println("Quick task executed!")
    },
)
scheduler.ScheduleTask(quickTask)
```

## Project Structure

```
task-scheduler/
├── models/
│   ├── task.go           # Task model and constructors
│   └── priority_queue.go # Priority queue implementation
├── scheduler.go          # Main scheduler implementation
├── main.go               # Example usage and demo
├── go.mod                # Go module file
└── README.md             # This file
```

## Implementation Details

### Priority Queue

The scheduler uses a min-heap (priority queue) to efficiently manage tasks ordered by execution time. This ensures:
- O(log n) insertion time
- O(log n) removal time
- O(1) peek at next task

### Worker Pool

The scheduler uses a semaphore-based worker pool to limit concurrent task executions:
- Configurable number of workers
- Non-blocking execution (tasks run in goroutines)
- Prevents resource exhaustion

### Thread Safety

All operations are protected with mutexes:
- `mu sync.RWMutex` protects the task queue and task map
- `runningMu sync.RWMutex` protects scheduler state
- Multiple goroutines can safely schedule tasks concurrently

### Efficient Waiting

The scheduler uses efficient sleep mechanisms:
- Sleeps when no tasks are ready
- Calculates exact wait time until next task
- Periodically checks for stop signal

### Task Rescheduling

Recurring tasks are automatically rescheduled:
- After execution, next execution time is calculated
- Task is re-inserted into the priority queue
- Continues until task is cancelled

## Running the Demo

```bash
cd task-scheduler
go run .
```

The demo program demonstrates:
1. One-time tasks scheduled for future execution
2. Delayed tasks
3. Recurring tasks with automatic rescheduling
4. Multiple concurrent tasks
5. Task status checking
6. Task cancellation
7. Long-running tasks that don't block the scheduler

## Performance Characteristics

- **Time Complexity**:
  - Schedule task: O(log n) where n is number of tasks
  - Cancel task: O(n) - requires rebuilding queue
  - Get next task: O(1) peek, O(log n) pop
  - Execute task: O(1) + handler execution time

- **Space Complexity**: O(n) where n is number of scheduled tasks

- **Concurrency**: Supports configurable number of concurrent executions

## Limitations

1. **In-Memory Only**: Tasks are not persisted across restarts
2. **No Task Dependencies**: Tasks are independent (no dependency graph)
3. **No Retries**: Failed tasks are not automatically retried
4. **No Distributed Support**: Single-process scheduler
5. **No CRON Expressions**: Only supports fixed intervals (can be extended)
6. **No Task Priorities**: All tasks are scheduled by time only

## Future Enhancements

- [ ] Add persistence to disk/database
- [ ] Support CRON expressions for complex schedules
- [ ] Implement task dependencies and DAG execution
- [ ] Add automatic retry mechanism for failed tasks
- [ ] Support task priorities
- [ ] Add distributed scheduling support
- [ ] Implement task result storage
- [ ] Add metrics and monitoring
- [ ] Support task timeouts
- [ ] Add webhook/HTTP callback support

## Design Decisions

### Why Priority Queue?

A priority queue (min-heap) is the most efficient data structure for this use case:
- Always know the next task to execute (O(1) peek)
- Efficient insertion and removal (O(log n))
- Natural ordering by execution time

### Why Worker Pool?

A worker pool prevents resource exhaustion:
- Limits concurrent executions
- Prevents too many goroutines
- Allows controlled concurrency

### Why Full Task Snapshots?

Tasks store full execution context:
- Easy to reschedule recurring tasks
- Simple status queries
- Extensible for future features

## Testing

Example test structure:

```go
func TestScheduler(t *testing.T) {
    scheduler := NewTaskScheduler(2)
    scheduler.Start()
    defer scheduler.Stop()

    executed := false
    task := models.NewDelayedTask(
        "test-task",
        "Test",
        100*time.Millisecond,
        func() {
            executed = true
        },
    )

    scheduler.ScheduleTask(task)
    time.Sleep(200 * time.Millisecond)

    if !executed {
        t.Error("Task was not executed")
    }
}
```

## References

Design inspired by: [Design a Task Scheduler - AlgoMaster](https://algomaster.io/learn/lld/design-task-scheduler)

## License

This project is part of the LowLevelDesign collection.

