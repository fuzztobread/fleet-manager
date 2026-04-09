package vehicle

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Service struct {
	store Store // interface, not a concrete type
}

func NewService(s Store) *Service {
	return &Service{store: s}
}

// Register creates a new vehicle.
func (s *Service) Register(ctx context.Context, name string, capacity int) (*Vehicle, error) {
	if name == "" {
		return nil, errors.New("vehicle name is required")
	}

	if capacity <= 0 {
		return nil, errors.New("capacity must be greater than zero")
	}

	v := &Vehicle{
		ID:       uuid.NewString(), // generate a unique ID
		Name:     name,
		Status:   StatusAvailable, // all new vehicles start available
		Capacity: capacity,
	}

	if err := s.store.Create(ctx, v); err != nil {
		return nil, fmt.Errorf("register vehicle: %w", err)
	}

	return v, nil
}

// UpdateStatus changes a vehicle's status.
func (s *Service) UpdateStatus(ctx context.Context, id string, status Status) error {
	switch status {
	case StatusAvailable, StatusEnRoute, StatusOffline:
	default:
		return fmt.Errorf("unknown status: %s", status)
	}

	return s.store.UpdateStatus(ctx, id, status)
}

// Get retrieves a single vehicle by ID.
func (s *Service) Get(ctx context.Context, id string) (*Vehicle, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.store.GetByID(ctx, id)
}

// List returns all vehicles.
func (s *Service) List(ctx context.Context) ([]*Vehicle, error) {
	return s.store.List(ctx)
}
