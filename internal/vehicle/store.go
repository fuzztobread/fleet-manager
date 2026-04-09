package vehicle

import "context"

type Store interface {
	Create(ctx context.Context, v *Vehicle) error
	GetByID(ctx context.Context, id string) (*Vehicle, error)
	UpdateStatus(ctx context.Context, id string, status Status) error
	List(ctx context.Context) ([]*Vehicle, error)
}
