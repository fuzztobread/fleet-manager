package vehicle

type Status string

const (
	StatusAvailable Status = "available" // ready to be dispatched
	StatusEnRoute   Status = "en_route"  // currently on a job
	StatusOffline   Status = "offline"   // not in service
)

// Vehicle is the core domain model.
// Tags like `json:"id"` tell encoding/json what key name to use
// when this struct is serialized to/from JSON.
type Vehicle struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Status   Status `json:"status"`
	Capacity int    `json:"capacity"` // max load units this vehicle carries
}
