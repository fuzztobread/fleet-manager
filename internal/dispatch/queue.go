// internal/dispatch/queue.go
package dispatch

import (
	"container/heap"
	"sync"
)

// jobHeap is the internal heap — implements heap.Interface.
// Not exported — callers use PriorityQueue only.
type jobHeap []*Job

func (h jobHeap) Len() int { return len(h) }

// Higher urgency = higher priority, so we invert the comparison
// to make container/heap behave as a max-heap.
func (h jobHeap) Less(i, j int) bool {
	if h[i].Urgency != h[j].Urgency {
		return h[i].Urgency > h[j].Urgency
	}
	// tie-break: earlier job wins (FIFO within same urgency)
	return h[i].CreatedAt.Before(h[j].CreatedAt)
}

func (h jobHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *jobHeap) Push(x any) {
	*h = append(*h, x.(*Job))
}

func (h *jobHeap) Pop() any {
	old := *h
	n := len(old)
	job := old[n-1]
	old[n-1] = nil // avoid memory leak
	*h = old[:n-1]
	return job
}

// PriorityQueue is a thread-safe max-heap of Jobs.
// Critical jobs are dequeued before High, High before Medium, etc.
// Within the same urgency level, jobs are ordered FIFO by CreatedAt.
type PriorityQueue struct {
	mu   sync.Mutex
	data jobHeap
}

func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{}
	heap.Init(&pq.data)
	return pq
}

// Push adds a job to the queue.
func (pq *PriorityQueue) Push(job *Job) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	heap.Push(&pq.data, job)
}

// Pop removes and returns the highest-priority job.
// Returns nil if the queue is empty.
func (pq *PriorityQueue) Pop() *Job {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if pq.data.Len() == 0 {
		return nil
	}
	return heap.Pop(&pq.data).(*Job)
}

// Peek returns the highest-priority job without removing it.
// Returns nil if the queue is empty.
func (pq *PriorityQueue) Peek() *Job {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if pq.data.Len() == 0 {
		return nil
	}
	return pq.data[0]
}

// Len returns the number of jobs in the queue.
func (pq *PriorityQueue) Len() int {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	return pq.data.Len()
}
