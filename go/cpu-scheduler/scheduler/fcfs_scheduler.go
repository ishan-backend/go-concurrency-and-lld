package scheduler

import (
	"sync"

	"github.com/ishanpandey/LowLevelDesign/go/cpu-scheduler/models"
	"github.com/ishanpandey/LowLevelDesign/go/cpu-scheduler/queue"
)

// FCFSScheduler implements First Come First Served scheduling
// Jobs are executed in the order they arrive (FIFO queue)
type FCFSScheduler struct {
	queue *queue.JobQueue
	mu    sync.RWMutex
}

// NewFCFSScheduler creates a new FCFS scheduler
func NewFCFSScheduler() *FCFSScheduler {
	return &FCFSScheduler{
		queue: queue.NewJobQueue(),
	}
}

// AddJob adds a job to the end of the queue
func (fs *FCFSScheduler) AddJob(job *models.Job) {
	if job == nil {
		return
	}
	
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	fs.queue.Add(job)
}

// GetNextJob returns the job that arrived first (FIFO)
func (fs *FCFSScheduler) GetNextJob() *models.Job {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	if fs.queue.IsEmpty() {
		return nil
	}
	
	// Remove and return the first job (FIFO)
	return fs.queue.Remove(0)
}

// HasJobs checks if there are any jobs remaining
func (fs *FCFSScheduler) HasJobs() bool {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	return !fs.queue.IsEmpty()
}

// GetName returns the name of this scheduler
func (fs *FCFSScheduler) GetName() string {
	return "First Come First Served (FCFS)"
}

