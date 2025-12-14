package scheduler

import "github.com/ishanpandey/LowLevelDesign/go/cpu-scheduler/models"

// SchedulerStrategy defines the interface for different CPU scheduling algorithms
// This is the Strategy Pattern - allows us to swap algorithms at runtime
type SchedulerStrategy interface {
	// AddJob adds a job to the scheduler's internal queue
	AddJob(job *models.Job)

	// GetNextJob returns the next job to be executed according to the scheduling algorithm
	// Returns nil if no jobs are available
	GetNextJob() *models.Job

	// HasJobs checks if there are any jobs remaining
	HasJobs() bool

	// GetName returns the name of the scheduling algorithm
	GetName() string
}
