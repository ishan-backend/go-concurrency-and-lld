package queue

import (
	"sync"

	"github.com/ishanpandey/LowLevelDesign/go/cpu-scheduler/models"
)

// JobQueue is a thread-safe queue for storing jobs
type JobQueue struct {
	jobs []*models.Job
	mu   sync.RWMutex
}

// NewJobQueue creates a new empty job queue
func NewJobQueue() *JobQueue {
	return &JobQueue{
		jobs: make([]*models.Job, 0),
	}
}

// Add adds a job to the queue
func (q *JobQueue) Add(job *models.Job) {
	if job == nil {
		return
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	q.jobs = append(q.jobs, job)
}

// Remove removes and returns the job at the given index
func (q *JobQueue) Remove(index int) *models.Job {
	q.mu.Lock()
	defer q.mu.Unlock()

	if index < 0 || index >= len(q.jobs) {
		return nil
	}

	job := q.jobs[index]
	q.jobs = append(q.jobs[:index], q.jobs[index+1:]...)
	return job
}

// Peek returns the job at the given index without removing it
func (q *JobQueue) Peek(index int) *models.Job {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if index < 0 || index >= len(q.jobs) {
		return nil
	}

	return q.jobs[index]
}

// Size returns the number of jobs in the queue
func (q *JobQueue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.jobs)
}

// IsEmpty checks if the queue is empty
func (q *JobQueue) IsEmpty() bool {
	return q.Size() == 0
}

// GetAll returns a copy of all jobs in the queue
func (q *JobQueue) GetAll() []*models.Job {
	q.mu.RLock()
	defer q.mu.RUnlock()

	jobs := make([]*models.Job, len(q.jobs))
	copy(jobs, q.jobs)
	return jobs
}

// Clear removes all jobs from the queue
func (q *JobQueue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.jobs = make([]*models.Job, 0)
}
