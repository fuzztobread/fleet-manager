// internal/dispatch/service.go
package dispatch

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"fleet-manager/internal/route"
	"fleet-manager/internal/vehicle"

	"github.com/google/uuid"
)

var ErrNoVehicleAvailable = errors.New("no available vehicle meets capacity requirement")

type Service struct {
	vehicles vehicle.Store
	routes   *route.Service
	queue    *PriorityQueue
}

func NewService(vehicles vehicle.Store, routes *route.Service) *Service {
	return &Service{
		vehicles: vehicles,
		routes:   routes,
		queue:    NewPriorityQueue(),
	}
}

// Enqueue validates and pushes a job onto the priority queue.
func (s *Service) Enqueue(from, to string, urgency Urgency, minCap int) (*Job, error) {
	if from == "" || to == "" {
		return nil, errors.New("from and to are required")
	}

	job := &Job{
		ID:        uuid.NewString(),
		From:      from,
		To:        to,
		Urgency:   urgency,
		MinCap:    minCap,
		CreatedAt: time.Now(),
	}

	s.queue.Push(job)
	log.Printf("[dispatch] enqueued job %s (%s → %s, urgency=%s)", job.ID, from, to, urgency)
	return job, nil
}

// Run starts the background worker. It blocks — call it in a goroutine.
// Shuts down cleanly when ctx is cancelled.
func (s *Service) Run(ctx context.Context) {
	log.Println("[dispatch] worker started")
	for {
		select {
		case <-ctx.Done():
			log.Println("[dispatch] worker shutting down")
			return
		default:
			job := s.queue.Pop()
			if job == nil {
				// queue empty — sleep briefly to avoid a busy-loop
				time.Sleep(500 * time.Millisecond)
				continue
			}
			if err := s.process(ctx, job); err != nil {
				log.Printf("[dispatch] job %s failed: %v", job.ID, err)
			}
		}
	}
}

// process handles a single job: pick vehicle, compute route, assign.
func (s *Service) process(ctx context.Context, job *Job) error {
	v, err := s.pickVehicle(ctx, job.MinCap)
	if err != nil {
		return fmt.Errorf("pick vehicle: %w", err)
	}

	_, err = s.routes.FindRoute(job.From, job.To)
	if err != nil {
		return fmt.Errorf("routing: %w", err)
	}

	if err := s.vehicles.UpdateStatus(ctx, v.ID, vehicle.StatusEnRoute); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	job.VehicleID = &v.ID
	log.Printf("[dispatch] job %s assigned to vehicle %s (%s)", job.ID, v.ID, v.Name)
	return nil
}

func (s *Service) pickVehicle(ctx context.Context, minCap int) (*vehicle.Vehicle, error) {
	all, err := s.vehicles.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list vehicles: %w", err)
	}
	for _, v := range all {
		if v.Status == vehicle.StatusAvailable && v.Capacity >= minCap {
			return v, nil
		}
	}
	return nil, ErrNoVehicleAvailable
}

// QueueLen exposes queue depth — useful for the ops endpoint.
func (s *Service) QueueLen() int {
	return s.queue.Len()
}

// PeekNext exposes the top job without removing it.
func (s *Service) PeekNext() *Job {
	return s.queue.Peek()
}
