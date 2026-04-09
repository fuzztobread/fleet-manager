package vehicle

import "context"

// Store is the interface this domain requires for persistence.
// It lives HERE (in vehicle/) not in storage/ — the domain
// defines what it needs, the adapter fulfills it.
//
// context.Context is passed into every method — Go convention.
// It carries deadlines and cancellation signals (e.g. if the
// HTTP request is cancelled, the DB query should stop too).
type Store interface {
	Create(ctx context.Context, v *Vehicle) error
	GetByID(ctx context.Context, id string) (*Vehicle, error)
	UpdateStatus(ctx context.Context, id string, status Status) error
	List(ctx context.Context) ([]*Vehicle, error)
}
