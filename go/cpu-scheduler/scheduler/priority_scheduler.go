package scheduler

import (
	"sync"

	"github.com/ishanpandey/LowLevelDesign/go/cpu-scheduler/models"
	"github.com/ishanpandey/LowLevelDesign/go/cpu-scheduler/queue"
)

// PriorityScheduler implements priority-based scheduling
// Jobs with higher priority are executed first
// For jobs with the same priority, FIFO (First In First Out) is used
type PriorityScheduler struct {
	queue *queue.JobQueue
	mu    sync.RWMutex
}

// NewPriorityScheduler creates a new priority-based scheduler
func NewPriorityScheduler() *PriorityScheduler {
	return &PriorityScheduler{
		queue: queue.NewJobQueue(),
	}
}

// AddJob adds a job to the priority queue
func (ps *PriorityScheduler) AddJob(job *models.Job) {
	if job == nil {
		return
	}
	
	ps.mu.Lock()
	defer ps.mu.Unlock()
	
	ps.queue.Add(job)
}

// GetNextJob returns the job with the highest priority
// If multiple jobs have the same priority, returns the one that arrived first (FIFO)
func (ps *PriorityScheduler) GetNextJob() *models.Job {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	
	if ps.queue.IsEmpty() {
		return nil
	}
	
	jobs := ps.queue.GetAll()
	
	// Find the job with highest priority
	// If priorities are equal, select the one that arrived first (earliest SubmittedAt)
	bestIndex := 0
	bestJob := jobs[0]
	
	for i := 1; i < len(jobs); i++ {
		currentJob := jobs[i]
		
		// Higher priority wins
		if currentJob.Priority > bestJob.Priority {
			bestIndex = i
			bestJob = currentJob
		} else if currentJob.Priority == bestJob.Priority {
			// Same priority: FIFO - earlier submission time wins
			if currentJob.SubmittedAt.Before(bestJob.SubmittedAt) {
				bestIndex = i
				bestJob = currentJob
			}
		}
	}
	
	// Remove and return the selected job
	return ps.queue.Remove(bestIndex)
}

// HasJobs checks if there are any jobs remaining
func (ps *PriorityScheduler) HasJobs() bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	
	return !ps.queue.IsEmpty()
}

// GetName returns the name of this scheduler
func (ps *PriorityScheduler) GetName() string {
	return "Priority Scheduling"
}

