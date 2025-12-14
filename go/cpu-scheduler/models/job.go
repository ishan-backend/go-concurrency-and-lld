package models

import (
	"fmt"
	"time"
)

// Job represents a task that needs to be executed by the CPU
type Job struct {
	ID            string    // Unique identifier for the job
	Priority      int       // Priority level (higher number = higher priority)
	SubmittedAt   time.Time // Time when the job was submitted
	BurstTime     int       // Time required to execute (for Round Robin)
	RemainingTime int       // Remaining execution time (for Round Robin)
}

// NewJob creates a new job with the given ID and priority
func NewJob(id string, priority int) *Job {
	return &Job{
		ID:            id,
		Priority:      priority,
		SubmittedAt:   time.Now(),
		BurstTime:     0,
		RemainingTime: 0,
	}
}

// NewJobWithBurstTime creates a new job with burst time (for Round Robin)
func NewJobWithBurstTime(id string, priority int, burstTime int) *Job {
	return &Job{
		ID:            id,
		Priority:      priority,
		SubmittedAt:   time.Now(),
		BurstTime:     burstTime,
		RemainingTime: burstTime,
	}
}

// String returns a string representation of the job
func (j *Job) String() string {
	return fmt.Sprintf("Job{ID: %s, Priority: %d, SubmittedAt: %s}",
		j.ID, j.Priority, j.SubmittedAt.Format(time.RFC3339Nano))
}

// IsCompleted checks if the job has completed execution (for Round Robin)
func (j *Job) IsCompleted() bool {
	return j.RemainingTime <= 0
}

// Execute executes the job for the given time quantum (for Round Robin)
// Returns true if the job is completed after this execution
func (j *Job) Execute(timeQuantum int) bool {
	if j.RemainingTime <= 0 {
		return true
	}

	if j.RemainingTime <= timeQuantum {
		j.RemainingTime = 0
		return true
	}

	j.RemainingTime -= timeQuantum
	return false
}
