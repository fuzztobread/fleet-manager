// internal/dispatch/model.go
package dispatch

import "time"

type Urgency int

const (
	UrgencyLow      Urgency = 1
	UrgencyMedium   Urgency = 2
	UrgencyHigh     Urgency = 3
	UrgencyCritical Urgency = 4
)

func (u Urgency) String() string {
	switch u {
	case UrgencyLow:
		return "low"
	case UrgencyMedium:
		return "medium"
	case UrgencyHigh:
		return "high"
	case UrgencyCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// Job represents a pending dispatch request.
// MinCap is carried so the worker knows the vehicle requirement at pop time.
type Job struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Urgency   Urgency   `json:"urgency"`
	MinCap    int       `json:"min_cap"`
	CreatedAt time.Time `json:"created_at"`
	VehicleID *string   `json:"vehicle_id,omitempty"`
}
