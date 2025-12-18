package models

import (
	"container/heap"
)

// TaskQueue is a priority queue for tasks (min-heap based on ExecuteAt time)
type TaskQueue []*Task

// Len returns the number of tasks in the queue
func (tq TaskQueue) Len() int {
	return len(tq)
}

// Less compares two tasks by their execution time
func (tq TaskQueue) Less(i, j int) bool {
	return tq[i].ExecuteAt.Before(tq[j].ExecuteAt)
}

// Swap swaps two tasks in the queue
func (tq TaskQueue) Swap(i, j int) {
	tq[i], tq[j] = tq[j], tq[i]
}

// Push adds a task to the queue
func (tq *TaskQueue) Push(x interface{}) {
	*tq = append(*tq, x.(*Task))
}

// Pop removes and returns the task with the earliest execution time
func (tq *TaskQueue) Pop() interface{} {
	old := *tq
	n := len(old)
	task := old[n-1]
	*tq = old[0 : n-1]
	return task
}

// Peek returns the task with the earliest execution time without removing it
func (tq TaskQueue) Peek() *Task {
	if len(tq) == 0 {
		return nil
	}
	return tq[0]
}

// NewTaskQueue creates a new task queue and initializes it as a heap
func NewTaskQueue() *TaskQueue {
	tq := &TaskQueue{}
	heap.Init(tq)
	return tq
}

