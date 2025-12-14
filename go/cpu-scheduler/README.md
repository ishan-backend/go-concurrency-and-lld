# CPU Scheduler - Low Level Design

A comprehensive, extensible CPU scheduling system implemented in Go that supports multiple scheduling algorithms using design patterns and clean architecture principles.

## 📋 Table of Contents

1. [Overview](#overview)
2. [Architecture & Design Patterns](#architecture--design-patterns)
3. [Project Structure](#project-structure)
4. [Core Components](#core-components)
5. [Scheduling Algorithms](#scheduling-algorithms)
6. [Usage Examples](#usage-examples)
7. [Interview Explanation Guide](#interview-explanation-guide)
8. [Design Decisions](#design-decisions)
9. [Extensibility](#extensibility)

---

## Overview

This system implements a CPU task scheduler that supports multiple scheduling algorithms:
- **Priority Scheduling**: Jobs with higher priority execute first (FIFO for same priority)
- **FCFS (First Come First Served)**: Jobs execute in arrival order
- **Round Robin**: Each job gets a fixed time quantum, then CPU switches to next job

### Key Features

✅ **Strategy Pattern**: Swap scheduling algorithms at runtime  
✅ **Thread-Safe**: All operations are concurrent-safe using mutexes  
✅ **Extensible**: Easy to add new scheduling algorithms  
✅ **Clean Architecture**: Separation of concerns with domain models  
✅ **Production-Ready**: Proper error handling and validation  

---

## Architecture & Design Patterns

### 1. **Strategy Pattern** ⭐ (Core Pattern)

**Purpose**: Define a family of algorithms (scheduling strategies), encapsulate each one, and make them interchangeable.

**Implementation**:
```go
type SchedulerStrategy interface {
    AddJob(job *models.Job)
    GetNextJob() *models.Job
    HasJobs() bool
    GetName() string
}
```

**Why it's powerful**:
- ✅ Algorithms can be swapped at runtime
- ✅ Open/Closed Principle: Add new algorithms without modifying existing code
- ✅ Single Responsibility: Each scheduler handles one algorithm
- ✅ Easy to test: Mock strategies for unit testing

**Real-world analogy**: Think of it like a navigation app - you can switch between "Fastest Route", "Shortest Route", or "Avoid Tolls" without changing the app's core functionality.

### 2. **Repository Pattern** (Implicit)

The `JobQueue` acts as a repository, abstracting data storage and access patterns.

### 3. **Factory Pattern** (Constructor Functions)

Each scheduler has a constructor function (`NewPriorityScheduler()`, `NewFCFSScheduler()`, etc.) that encapsulates object creation.

---

## Project Structure

```
go/cpu-scheduler/
├── models/
│   └── job.go              # Domain model: Job entity
├── queue/
│   └── job_queue.go        # Thread-safe job queue implementation
├── scheduler/
│   ├── strategy.go         # Strategy interface (core pattern)
│   ├── priority_scheduler.go    # Priority scheduling implementation
│   ├── fcfs_scheduler.go        # FCFS scheduling implementation
│   └── round_robin_scheduler.go # Round Robin implementation
├── cpu_scheduler.go        # Main scheduler manager (context)
├── main.go                 # Demo and examples
├── go.mod                  # Go module definition
└── README.md               # This file
```

### Package Responsibilities

- **`models`**: Domain entities (Job)
- **`queue`**: Data structure for job storage
- **`scheduler`**: Scheduling algorithm implementations
- **Root package**: Main scheduler manager and demo

---

## Core Components

### 1. Job Model (`models/job.go`)

Represents a task that needs CPU execution.

```go
type Job struct {
    ID            string    // Unique identifier
    Priority      int       // Priority level (higher = more important)
    SubmittedAt   time.Time // Arrival time
    BurstTime     int       // Total execution time needed
    RemainingTime int       // Remaining execution time
}
```

**Key Methods**:
- `NewJob(id, priority)`: Creates a simple job
- `NewJobWithBurstTime(id, priority, burstTime)`: Creates job with execution time
- `Execute(timeQuantum)`: Executes job for given time (Round Robin)

### 2. Job Queue (`queue/job_queue.go`)

Thread-safe queue for storing jobs.

**Features**:
- ✅ `sync.RWMutex` for concurrent access
- ✅ Add, Remove, Peek operations
- ✅ Size and IsEmpty checks

### 3. Scheduler Strategy Interface (`scheduler/strategy.go`)

The heart of the Strategy Pattern - defines the contract all schedulers must follow.

```go
type SchedulerStrategy interface {
    AddJob(job *models.Job)
    GetNextJob() *models.Job
    HasJobs() bool
    GetName() string
}
```

### 4. CPU Scheduler Manager (`cpu_scheduler.go`)

The context that uses the strategy. It:
- Manages the current scheduling strategy
- Provides a clean API for adding/executing jobs
- Allows runtime strategy switching

---

## Scheduling Algorithms

### 1. Priority Scheduling

**Algorithm**: Jobs with higher priority execute first. If priorities are equal, FIFO (First In First Out) is used.

**Implementation Logic**:
```go
// Find job with highest priority
// If same priority, select earliest arrival time
```

**Example**:
```
Jobs: Job1(priority=5), Job2(priority=3), Job3(priority=5)
Execution Order: Job1 → Job3 → Job2
Reason: Job1 and Job3 have highest priority (5), Job1 arrived first
```

**Time Complexity**: O(n) for GetNextJob (linear search)
**Space Complexity**: O(n) for storing jobs

**Use Cases**:
- Operating systems with priority-based task scheduling
- Real-time systems where some tasks are more critical
- Process managers that prioritize system processes

### 2. FCFS (First Come First Served)

**Algorithm**: Jobs execute in the order they arrive (FIFO queue).

**Implementation**: Simple queue - first job added is first job executed.

**Example**:
```
Jobs: JobA, JobB, JobC (added in this order)
Execution Order: JobA → JobB → JobC
```

**Time Complexity**: O(1) for both AddJob and GetNextJob
**Space Complexity**: O(n)

**Use Cases**:
- Simple batch processing systems
- Print queues
- When all jobs have equal importance

**Pros**: Simple, fair, no starvation  
**Cons**: Not optimal for varying job lengths (convoy effect)

### 3. Round Robin

**Algorithm**: Each job gets a fixed time quantum. After the quantum expires, CPU switches to the next job. Incomplete jobs go back to the queue.

**Key Parameters**:
- **Time Quantum**: Fixed time slice per job (e.g., 2 time units)

**Implementation Logic**:
```go
1. Get first job from queue
2. Execute for time quantum
3. If not completed, add back to end of queue
4. Return to step 1
```

**Example**:
```
Jobs: JobX(burst=5), JobY(burst=3), JobZ(burst=4)
Time Quantum: 2

Execution Timeline:
Turn 1: JobX (5→3 remaining)
Turn 2: JobY (3→1 remaining)
Turn 3: JobZ (4→2 remaining)
Turn 4: JobX (3→1 remaining)
Turn 5: JobY (1→0, completed)
Turn 6: JobZ (2→0, completed)
Turn 7: JobX (1→0, completed)
```

**Time Complexity**: O(1) per execution (amortized)
**Space Complexity**: O(n)

**Use Cases**:
- Time-sharing systems
- Interactive applications
- When fairness is important

**Pros**: Fair, no starvation, good response time  
**Cons**: Context switching overhead, not optimal for long jobs

---

## Usage Examples

### Basic Priority Scheduling

```go
// Create priority scheduler
priorityScheduler := scheduler.NewPriorityScheduler()
cpuScheduler := NewCPUScheduler(priorityScheduler)

// Add jobs
cpuScheduler.AddJob("Job1", 5)  // High priority
cpuScheduler.AddJob("Job2", 3)  // Medium priority
cpuScheduler.AddJob("Job3", 5)  // High priority

// Execute jobs
for cpuScheduler.HasJobs() {
    jobID, err := cpuScheduler.ExecuteNext()
    if err != nil {
        break
    }
    fmt.Printf("Executed: %s\n", jobID)
}
// Output: Job1, Job3, Job2
```

### FCFS Scheduling

```go
fcfsScheduler := scheduler.NewFCFSScheduler()
cpuScheduler := NewCPUScheduler(fcfsScheduler)

cpuScheduler.AddJob("JobA", 1)
cpuScheduler.AddJob("JobB", 5)  // Higher priority, but FCFS ignores it
cpuScheduler.AddJob("JobC", 3)

// Executes in order: JobA → JobB → JobC
```

### Round Robin Scheduling

```go
// Create Round Robin with time quantum of 2
rrScheduler := scheduler.NewRoundRobinScheduler(2)
cpuScheduler := NewCPUScheduler(rrScheduler)

// Add jobs with burst times
cpuScheduler.AddJobWithBurstTime("JobX", 1, 5)
cpuScheduler.AddJobWithBurstTime("JobY", 2, 3)
cpuScheduler.AddJobWithBurstTime("JobZ", 1, 4)

// Each job gets 2 time units per turn
```

### Runtime Strategy Switching

```go
// Start with Priority
priorityScheduler := scheduler.NewPriorityScheduler()
cpuScheduler := NewCPUScheduler(priorityScheduler)

cpuScheduler.AddJob("Task1", 3)
cpuScheduler.AddJob("Task2", 1)

// Execute one job
cpuScheduler.ExecuteNext()

// Switch to FCFS at runtime!
fcfsScheduler := scheduler.NewFCFSScheduler()
cpuScheduler.SetStrategy(fcfsScheduler)

// Remaining jobs now execute with FCFS
```

### Running the Demo

```bash
cd go/cpu-scheduler
go run .
```

---

## Interview Explanation Guide

### How to Explain This Design in an Interview

#### 1. **Start with the Problem Statement**

"I was asked to design a CPU scheduler that supports multiple scheduling algorithms. The key requirements were:
- Support Priority, FCFS, and Round Robin algorithms
- Make it extensible for future algorithms
- Ensure thread-safety
- Clean, maintainable code"

#### 2. **Explain the Core Design Pattern**

"I chose the **Strategy Pattern** as the foundation because:
- It allows me to encapsulate each scheduling algorithm
- Algorithms can be swapped at runtime
- Follows Open/Closed Principle - open for extension, closed for modification
- Each scheduler is independently testable"

**Show the interface**:
```go
type SchedulerStrategy interface {
    AddJob(job *models.Job)
    GetNextJob() *models.Job
    HasJobs() bool
    GetName() string
}
```

#### 3. **Walk Through the Architecture**

**Domain Models**:
- "I created a `Job` model to represent tasks. It includes ID, priority, submission time, and burst time for Round Robin."

**Queue Abstraction**:
- "I abstracted job storage into a `JobQueue` with thread-safe operations using `sync.RWMutex`."

**Strategy Implementations**:
- "Each algorithm (Priority, FCFS, Round Robin) implements the `SchedulerStrategy` interface."

**Context/Manager**:
- "The `CPUScheduler` acts as the context - it holds a reference to the current strategy and delegates operations to it."

#### 4. **Explain Each Algorithm**

**Priority Scheduling**:
- "Jobs with higher priority execute first. For same priority, I use FIFO based on submission time. This is O(n) because I search through all jobs to find the highest priority."

**FCFS**:
- "Simple FIFO queue - first job added is first executed. O(1) operations, very efficient."

**Round Robin**:
- "Each job gets a fixed time quantum. After execution, if not completed, it goes back to the end of the queue. This ensures fairness and prevents starvation."

#### 5. **Highlight Extensibility**

"To add a new algorithm (e.g., Shortest Job First), I just:
1. Create a new struct implementing `SchedulerStrategy`
2. Implement the four required methods
3. No changes needed to existing code!

This demonstrates the power of the Strategy Pattern."

#### 6. **Discuss Trade-offs**

**Priority Scheduling**:
- ✅ Prioritizes important tasks
- ❌ O(n) complexity, potential starvation for low-priority jobs

**FCFS**:
- ✅ Simple, fair, O(1) operations
- ❌ Not optimal for varying job lengths

**Round Robin**:
- ✅ Fair, no starvation, good response time
- ❌ Context switching overhead

#### 7. **Concurrency Considerations**

"I used `sync.RWMutex` for thread-safety:
- Read locks for `HasJobs()`, `Peek()` operations
- Write locks for `AddJob()`, `GetNextJob()` operations
- This allows concurrent reads while preventing race conditions"

#### 8. **Real-World Applications**

- **Priority**: Operating systems, real-time systems
- **FCFS**: Batch processing, print queues
- **Round Robin**: Time-sharing systems, interactive applications

---

## Design Decisions

### Why Strategy Pattern?

1. **Extensibility**: New algorithms can be added without modifying existing code
2. **Testability**: Each strategy can be tested independently
3. **Runtime Flexibility**: Algorithms can be swapped based on system conditions
4. **Separation of Concerns**: Each scheduler focuses on one algorithm

### Why Separate Queue Package?

1. **Reusability**: Queue can be used by any scheduler
2. **Testability**: Queue logic can be tested independently
3. **Abstraction**: Schedulers don't need to know queue implementation details

### Why Thread-Safe?

1. **Real-world Requirement**: CPU schedulers in OS handle concurrent job submissions
2. **Safety**: Prevents race conditions and data corruption
3. **Scalability**: Can handle multiple threads adding/executing jobs

### Why Burst Time for Round Robin?

1. **Realistic**: Real jobs have execution times
2. **Demonstration**: Shows how Round Robin handles incomplete jobs
3. **Completeness**: Makes the implementation more realistic

---

## Extensibility

### Adding a New Scheduling Algorithm

To add a new algorithm (e.g., **Shortest Job First**):

#### Step 1: Create the Scheduler

```go
// scheduler/sjf_scheduler.go
package scheduler

import (
    "sort"
    "sync"
    "github.com/ishanpandey/LowLevelDesign/go/cpu-scheduler/models"
    "github.com/ishanpandey/LowLevelDesign/go/cpu-scheduler/queue"
)

type SJFScheduler struct {
    queue *queue.JobQueue
    mu    sync.RWMutex
}

func NewSJFScheduler() *SJFScheduler {
    return &SJFScheduler{
        queue: queue.NewJobQueue(),
    }
}

func (sjf *SJFScheduler) AddJob(job *models.Job) {
    sjf.mu.Lock()
    defer sjf.mu.Unlock()
    sjf.queue.Add(job)
}

func (sjf *SJFScheduler) GetNextJob() *models.Job {
    sjf.mu.Lock()
    defer sjf.mu.Unlock()
    
    if sjf.queue.IsEmpty() {
        return nil
    }
    
    jobs := sjf.queue.GetAll()
    
    // Sort by burst time (shortest first)
    sort.Slice(jobs, func(i, j int) bool {
        return jobs[i].BurstTime < jobs[j].BurstTime
    })
    
    // Remove and return shortest job
    shortestJob := jobs[0]
    // Find index in original queue
    for i, job := range sjf.queue.GetAll() {
        if job.ID == shortestJob.ID {
            return sjf.queue.Remove(i)
        }
    }
    return nil
}

func (sjf *SJFScheduler) HasJobs() bool {
    sjf.mu.RLock()
    defer sjf.mu.RUnlock()
    return !sjf.queue.IsEmpty()
}

func (sjf *SJFScheduler) GetName() string {
    return "Shortest Job First (SJF)"
}
```

#### Step 2: Use It

```go
sjfScheduler := scheduler.NewSJFScheduler()
cpuScheduler := NewCPUScheduler(sjfScheduler)
// That's it! No other code changes needed.
```

### Key Points

✅ **No modifications** to existing code  
✅ **Follows Open/Closed Principle**  
✅ **Automatically works** with CPUScheduler  
✅ **Easy to test** independently  

---

## Key Takeaways for Interviews

1. **Design Pattern**: Strategy Pattern is perfect for interchangeable algorithms
2. **SOLID Principles**: 
   - Single Responsibility (each scheduler does one thing)
   - Open/Closed (open for extension, closed for modification)
   - Dependency Inversion (depend on abstractions, not concretions)
3. **Concurrency**: Thread-safety is crucial for real-world systems
4. **Extensibility**: Design should make adding features easy
5. **Clean Code**: Proper separation of concerns, meaningful names, error handling

---

## Future Enhancements

Potential improvements for production systems:

1. **Priority Queues**: Use heap data structure for O(log n) Priority Scheduling
2. **Preemption**: Support for preemptive scheduling
3. **Metrics**: Track waiting time, turnaround time, response time
4. **Multi-level Queues**: Support multiple priority levels
5. **Deadline Scheduling**: Support for real-time deadlines
6. **Load Balancing**: Distribute jobs across multiple CPUs

---

## Conclusion

This CPU scheduler demonstrates:
- ✅ **Clean Architecture**: Proper separation of concerns
- ✅ **Design Patterns**: Strategy Pattern for algorithm flexibility
- ✅ **Concurrency**: Thread-safe operations
- ✅ **Extensibility**: Easy to add new algorithms
- ✅ **Production-Ready**: Error handling, validation, documentation

Perfect for understanding low-level design principles and explaining design decisions in technical interviews!

---

## Author Notes

This implementation combines two requirements:
1. **Statement 1**: Priority-based scheduling with FIFO for same priority
2. **Statement 2**: Multiple scheduling algorithms (FCFS, Round Robin, Priority) with Strategy Pattern

The design emphasizes **readability** and **memorability** to help you explain your thought process clearly in interviews.

