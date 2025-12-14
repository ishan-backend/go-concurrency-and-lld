package main

import (
	"fmt"
	"sync"

	"github.com/ishanpandey/LowLevelDesign/go/cpu-scheduler/models"
	"github.com/ishanpandey/LowLevelDesign/go/cpu-scheduler/scheduler"
)

// CPUScheduler is the main scheduler that manages job execution
// It uses the Strategy Pattern to support different scheduling algorithms
type CPUScheduler struct {
	strategy scheduler.SchedulerStrategy
	mu       sync.RWMutex
}

// NewCPUScheduler creates a new CPU scheduler with the given strategy
func NewCPUScheduler(strategy scheduler.SchedulerStrategy) *CPUScheduler {
	if strategy == nil {
		// Default to Priority Scheduling
		strategy = scheduler.NewPriorityScheduler()
	}
	
	return &CPUScheduler{
		strategy: strategy,
	}
}

// SetStrategy changes the scheduling algorithm at runtime
// This demonstrates the flexibility of the Strategy Pattern
func (cs *CPUScheduler) SetStrategy(strategy scheduler.SchedulerStrategy) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	
	cs.strategy = strategy
}

// GetStrategy returns the current scheduling strategy
func (cs *CPUScheduler) GetStrategy() scheduler.SchedulerStrategy {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	
	return cs.strategy
}

// AddJob adds a new job to the scheduler
func (cs *CPUScheduler) AddJob(jobID string, priority int) error {
	if jobID == "" {
		return fmt.Errorf("job ID cannot be empty")
	}
	
	cs.mu.RLock()
	strategy := cs.strategy
	cs.mu.RUnlock()
	
	job := models.NewJob(jobID, priority)
	strategy.AddJob(job)
	
	return nil
}

// AddJobWithBurstTime adds a job with burst time (for Round Robin)
func (cs *CPUScheduler) AddJobWithBurstTime(jobID string, priority int, burstTime int) error {
	if jobID == "" {
		return fmt.Errorf("job ID cannot be empty")
	}
	
	if burstTime <= 0 {
		return fmt.Errorf("burst time must be positive")
	}
	
	cs.mu.RLock()
	strategy := cs.strategy
	cs.mu.RUnlock()
	
	job := models.NewJobWithBurstTime(jobID, priority, burstTime)
	strategy.AddJob(job)
	
	return nil
}

// ExecuteNext executes the next job according to the current scheduling strategy
// Returns the job ID if a job was executed, or an error if no jobs are available
func (cs *CPUScheduler) ExecuteNext() (string, error) {
	cs.mu.RLock()
	strategy := cs.strategy
	cs.mu.RUnlock()
	
	if !strategy.HasJobs() {
		return "", fmt.Errorf("no jobs available to execute")
	}
	
	job := strategy.GetNextJob()
	if job == nil {
		return "", fmt.Errorf("failed to get next job")
	}
	
	// Job execution is instantaneous (as per requirements)
	// In a real system, this would involve actual CPU execution
	return job.ID, nil
}

// HasJobs checks if there are any jobs remaining
func (cs *CPUScheduler) HasJobs() bool {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	
	return cs.strategy.HasJobs()
}

// GetSchedulerName returns the name of the current scheduling algorithm
func (cs *CPUScheduler) GetSchedulerName() string {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	
	return cs.strategy.GetName()
}

