package scheduler

import (
	"sync"

	"github.com/ishanpandey/LowLevelDesign/go/cpu-scheduler/models"
	"github.com/ishanpandey/LowLevelDesign/go/cpu-scheduler/queue"
)

// RoundRobinScheduler implements Round Robin scheduling
// Each job gets a fixed time quantum, then CPU switches to the next job
// Jobs that haven't completed are added back to the queue
type RoundRobinScheduler struct {
	queue       *queue.JobQueue
	timeQuantum int // Time slice allocated to each job
	mu          sync.RWMutex
}

// NewRoundRobinScheduler creates a new Round Robin scheduler with the given time quantum
func NewRoundRobinScheduler(timeQuantum int) *RoundRobinScheduler {
	if timeQuantum <= 0 {
		timeQuantum = 2 // Default time quantum
	}
	
	return &RoundRobinScheduler{
		queue:       queue.NewJobQueue(),
		timeQuantum: timeQuantum,
	}
}

// AddJob adds a job to the round robin queue
func (rr *RoundRobinScheduler) AddJob(job *models.Job) {
	if job == nil {
		return
	}
	
	rr.mu.Lock()
	defer rr.mu.Unlock()
	
	rr.queue.Add(job)
}

// GetNextJob returns the next job in the round robin queue
func (rr *RoundRobinScheduler) GetNextJob() *models.Job {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	
	if rr.queue.IsEmpty() {
		return nil
	}
	
	// Get the first job from the queue
	job := rr.queue.Remove(0)
	
	// Execute the job for the time quantum
	completed := job.Execute(rr.timeQuantum)
	
	// If job is not completed, add it back to the end of the queue
	if !completed {
		rr.queue.Add(job)
	}
	
	// Return the job (caller can check if it's completed)
	return job
}

// HasJobs checks if there are any jobs remaining
func (rr *RoundRobinScheduler) HasJobs() bool {
	rr.mu.RLock()
	defer rr.mu.RUnlock()
	
	return !rr.queue.IsEmpty()
}

// GetName returns the name of this scheduler
func (rr *RoundRobinScheduler) GetName() string {
	return "Round Robin"
}

// GetTimeQuantum returns the time quantum for this scheduler
func (rr *RoundRobinScheduler) GetTimeQuantum() int {
	rr.mu.RLock()
	defer rr.mu.RUnlock()
	
	return rr.timeQuantum
}

